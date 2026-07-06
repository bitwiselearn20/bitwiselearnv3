// Package code ports apps/python-server/routers/code_runner.py to Go.
//
// Structural change from the legacy router, per the rewrite plan's
// code-level optimizations: test cases run concurrently (bounded by a
// semaphore via errgroup.SetLimit) against the pooled Piston client instead
// of the original's sequential for-loop, so a submission with N test cases
// takes ~max(latency) instead of ~sum(latency).
package code

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/sync/errgroup"

	"github.com/bitwiselearn/go-server/internal/db"
	"github.com/bitwiselearn/go-server/internal/dbutil"
	"github.com/bitwiselearn/go-server/internal/middleware"
	"github.com/bitwiselearn/go-server/internal/models"
	"github.com/bitwiselearn/go-server/internal/response"
	"github.com/bitwiselearn/go-server/internal/services/piston"
)

const reqTimeout = 30 * time.Second

// maxConcurrentExecutions bounds how many test cases run against Piston at
// once per submission, so one large submission can't monopolize the shared
// execution pool.
const maxConcurrentExecutions = 8

// Deps holds the dependencies the code-runner handlers need.
type Deps struct {
	DB     *db.Client
	Auth   *middleware.Auth
	Piston *piston.Client
}

// Register mounts every code-runner route on g (expected prefix "/api/v1/code").
func Register(g *echo.Group, d Deps) {
	h := &handler{db: d.DB, piston: d.Piston}
	authed := d.Auth.Required
	studentTeacher := d.Auth.RequireRoles(models.RoleStudent, models.RoleTeacher)

	g.POST("/run", h.run, authed)
	g.POST("/compile", h.compile)
	g.POST("/submit", h.submit, authed, studentTeacher)
}

type handler struct {
	db     *db.Client
	piston *piston.Client
}

func (h *handler) coll(name string) *mongo.Collection { return h.db.Coll(name) }

func extractStdoutStderr(result piston.Result) (string, string) {
	stdout := result.Stdout()
	var stderr string
	if run, ok := result["run"].(map[string]any); ok {
		if s, ok := run["stderr"].(string); ok && s != "" {
			stderr = s
		}
	}
	if stderr == "" {
		if s, ok := result["stderr"].(string); ok && s != "" {
			stderr = s
		} else if e := result.Error(); e != "" {
			stderr = e
		}
	}
	return stdout, stderr
}

// findProblemTemplate ports _find_problem_template.
func (h *handler) findProblemTemplate(ctx context.Context, problemID primitive.ObjectID, language string) *models.ProblemTemplate {
	var t models.ProblemTemplate
	err := h.coll(models.CollProblemTemplates).FindOne(ctx, bson.M{
		"problem_id": problemID, "language": language,
	}).Decode(&t)
	if err == nil {
		return &t
	}
	normalized := piston.NormalizeLanguage(language)
	if normalized == language {
		return nil
	}
	err = h.coll(models.CollProblemTemplates).FindOne(ctx, bson.M{
		"problem_id": problemID, "language": normalized,
	}).Decode(&t)
	if err == nil {
		return &t
	}
	return nil
}

// composeFullCode ports _compose_full_code.
func composeFullCode(userCode string, template *models.ProblemTemplate) string {
	if template == nil || template.FunctionBody == "" {
		return userCode
	}
	if strings.Contains(template.FunctionBody, "_solution_") {
		return strings.ReplaceAll(template.FunctionBody, "_solution_", userCode)
	}
	return userCode + "\n" + template.FunctionBody
}

type testCaseResult struct {
	testCaseID     primitive.ObjectID
	input          string
	expectedOutput string
	actualOutput   string
	passed         bool
	stderr         string
	testType       string
	runtimeCode    any
}

// runTestCasesParallel executes fullCode against every test case concurrently
// (bounded by maxConcurrentExecutions) and returns results in the same order
// as testCases.
func (h *handler) runTestCasesParallel(ctx context.Context, language, fullCode string, testCases []models.ProblemTestCase) []testCaseResult {
	results := make([]testCaseResult, len(testCases))
	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(maxConcurrentExecutions)

	for i, tc := range testCases {
		i, tc := i, tc
		g.Go(func() error {
			result := h.piston.Execute(gctx, language, fullCode, tc.Input)
			actualRaw, stderr := extractStdoutStderr(result)
			actual := strings.TrimSpace(strings.ReplaceAll(actualRaw, "\r\n", "\n"))
			expected := strings.TrimSpace(tc.Output)
			var runtimeCode any
			if run, ok := result["run"].(map[string]any); ok {
				runtimeCode = run["code"]
			}
			results[i] = testCaseResult{
				testCaseID: tc.ID, input: tc.Input, expectedOutput: expected,
				actualOutput: actual, passed: actual == expected, stderr: stderr,
				testType: tc.TestType, runtimeCode: runtimeCode,
			}
			return nil
		})
	}
	_ = g.Wait()
	return results
}

func (h *handler) run(c echo.Context) error {
	var body runCodeRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	problemID, ok := dbutil.ParseID(body.ProblemID)
	if !ok {
		return response.Err(c, http.StatusNotFound, "Problem not found", "Not found")
	}
	user := middleware.UserFrom(c)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	if h.coll(models.CollProblems).FindOne(ctx, bson.M{"_id": problemID}).Err() != nil {
		return response.Err(c, http.StatusNotFound, "Problem not found", "Not found")
	}

	template := h.findProblemTemplate(ctx, problemID, body.Language)

	filter := bson.M{"problem_id": problemID}
	if user.Type != models.RoleSuperadmin && user.Type != models.RoleAdmin {
		filter["test_type"] = models.TestcaseExample
	}
	cur, err := h.coll(models.CollProblemTestCases).Find(ctx, filter)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch test cases", err.Error())
	}
	var testCases []models.ProblemTestCase
	if err := cur.All(ctx, &testCases); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch test cases", err.Error())
	}
	if len(testCases) == 0 {
		return response.Err(c, http.StatusBadRequest, "No test cases found", "No test cases")
	}

	fullCode := composeFullCode(body.Code, template)
	results := h.runTestCasesParallel(ctx, body.Language, fullCode, testCases)

	data := make([]map[string]any, 0, len(results))
	allPassed := true
	for _, r := range results {
		if !r.passed {
			allPassed = false
		}
		data = append(data, map[string]any{
			"test_case_id": r.testCaseID.Hex(), "input": r.input, "expected_output": r.expectedOutput,
			"actual_output": r.actualOutput, "passed": r.passed, "stderr": r.stderr, "test_type": r.testType,
		})
	}

	return response.OK(c, "Code executed", map[string]any{"results": data, "all_passed": allPassed})
}

func (h *handler) compile(c echo.Context) error {
	var body compileCodeRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	result := h.piston.Execute(ctx, body.Language, body.Code, body.Stdin)
	if errMsg := result.Error(); errMsg != "" {
		if details := result.Details(); details != "" {
			errMsg = errMsg + ": " + details
		}
		return response.Err(c, http.StatusBadRequest, "Compile failed", errMsg)
	}

	if compileOut, ok := result["compile"].(map[string]any); ok && len(compileOut) > 0 {
		if code, ok := compileOut["code"]; ok && code != nil && code != float64(0) {
			compileErr := "Compilation failed"
			if s, ok := compileOut["stderr"].(string); ok && s != "" {
				compileErr = s
			} else if s, ok := compileOut["output"].(string); ok && s != "" {
				compileErr = s
			} else if s, ok := compileOut["message"].(string); ok && s != "" {
				compileErr = s
			}
			return response.Err(c, http.StatusBadRequest, "Compile failed", compileErr)
		}
	}

	stdout, stderr := extractStdoutStderr(result)
	var runCode, runSignal any
	if run, ok := result["run"].(map[string]any); ok {
		runCode, runSignal = run["code"], run["signal"]
	}
	return response.OK(c, "Code compiled", map[string]any{
		"stdout": stdout, "stderr": stderr, "code": runCode, "signal": runSignal,
	})
}

func (h *handler) submit(c echo.Context) error {
	var body submitCodeRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	problemID, ok := dbutil.ParseID(body.ProblemID)
	if !ok {
		return response.Err(c, http.StatusNotFound, "Problem not found", "Not found")
	}
	user := middleware.UserFrom(c)
	studentID, _ := dbutil.ParseID(user.ID)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	if h.coll(models.CollProblems).FindOne(ctx, bson.M{"_id": problemID}).Err() != nil {
		return response.Err(c, http.StatusNotFound, "Problem not found", "Not found")
	}

	template := h.findProblemTemplate(ctx, problemID, body.Language)

	cur, err := h.coll(models.CollProblemTestCases).Find(ctx, bson.M{"problem_id": problemID})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch test cases", err.Error())
	}
	var testCases []models.ProblemTestCase
	if err := cur.All(ctx, &testCases); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch test cases", err.Error())
	}
	if len(testCases) == 0 {
		return response.Err(c, http.StatusBadRequest, "No test cases", "No test cases")
	}

	fullCode := composeFullCode(body.Code, template)
	results := h.runTestCasesParallel(ctx, body.Language, fullCode, testCases)

	allPassed := true
	var failedTC *string
	passedCount := 0
	for _, r := range results {
		if r.passed {
			passedCount++
		} else if allPassed {
			allPassed = false
			id := r.testCaseID.Hex()
			failedTC = &id
		}
	}

	status := models.ProblemFailed
	if allPassed {
		status = models.ProblemSuccess
	}
	now := time.Now().UTC()
	submission := models.ProblemSubmission{
		ID: primitive.NewObjectID(), Code: body.Code, Status: status, StudentID: studentID,
		ProblemID: problemID, FailedTestCase: failedTC, SubmittedAt: now, CreatedAt: now, UpdatedAt: now,
	}
	if _, err := h.coll(models.CollProblemSubmissions).InsertOne(ctx, submission); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to submit code", err.Error())
	}

	docs := make([]any, 0, len(results))
	for _, r := range results {
		var runtime *string
		if r.runtimeCode != nil {
			s := formatRuntime(r.runtimeCode)
			runtime = &s
		}
		actual := r.actualOutput
		docs = append(docs, models.ProblemSubmissionTestCase{
			ID: primitive.NewObjectID(), SubmissionID: submission.ID, TestCaseID: r.testCaseID,
			Passed: r.passed, ActualOutput: &actual, Runtime: runtime, CreatedAt: now,
		})
	}
	if len(docs) > 0 {
		if _, err := h.coll(models.CollProblemSubmissionTestCases).InsertMany(ctx, docs); err != nil {
			return response.Err(c, http.StatusInternalServerError, "Failed to save test case results", err.Error())
		}
	}

	return response.OK(c, "Code submitted", map[string]any{
		"submission_id": submission.ID.Hex(), "status": status, "all_passed": allPassed,
		"total_test_cases": len(testCases), "passed_count": passedCount,
	})
}

// formatRuntime mirrors str(tcr["runtime"]): Piston's run.code is a JSON
// number (process exit code), rendered without a trailing ".0" for whole
// values so it reads the same as Python's str(int).
func formatRuntime(v any) string {
	f, ok := v.(float64)
	if !ok {
		return fmt.Sprintf("%v", v)
	}
	if f == math.Trunc(f) {
		return strconv.FormatInt(int64(f), 10)
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}
