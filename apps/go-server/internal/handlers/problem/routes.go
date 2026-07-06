// Package problem ports apps/python-server/routers/dsa_problem.py to Go.
package problem

import (
	"context"
	"net/http"
	"regexp"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/bitwiselearn/go-server/internal/db"
	"github.com/bitwiselearn/go-server/internal/dbutil"
	"github.com/bitwiselearn/go-server/internal/middleware"
	"github.com/bitwiselearn/go-server/internal/models"
	"github.com/bitwiselearn/go-server/internal/response"
)

// Deps holds the dependencies the problem handlers need.
type Deps struct {
	DB   *db.Client
	Auth *middleware.Auth
}

const reqTimeout = 8 * time.Second

// Register mounts every DSA problem route on g (expected prefix "/api/v1/problems").
func Register(g *echo.Group, d Deps) {
	h := &handler{db: d.DB}
	notStudent := d.Auth.NotStudent()
	authed := d.Auth.Required

	// Public
	g.GET("/get-all-dsa-problem/", h.getAllProblems)
	g.GET("/get-all-listed-problem/", h.getAllListedProblems)
	g.POST("/search-question", h.searchQuestion)
	g.GET("/get-dsa-problem/:id/", h.getProblem)
	g.GET("/get-dsa-problems-by-tag", h.getProblemsByTag)

	// Admin (not_student)
	g.PUT("/change-status/:id", h.changeStatus, authed, notStudent)
	g.GET("/admin/get-dsa-problem/:id", h.adminGetProblem, authed, notStudent)
	g.GET("/admin/get-dsa-problem/testcases/:id", h.adminGetTestCases, authed, notStudent)
	g.GET("/admin/get-dsa-problem/solution/:id", h.adminGetSolutions, authed, notStudent)
	g.GET("/admin/get-dsa-problem/submission/:id", h.adminGetSubmissions, authed, notStudent)
	g.GET("/admin/get-dsa-problem/templates/:id", h.adminGetTemplates, authed, notStudent)
	g.POST("/add-problem/", h.addProblem, authed, notStudent)
	g.PATCH("/update-problem/:id", h.updateProblem, authed, notStudent)
	g.DELETE("/delete-problem/:id", h.deleteProblem, authed, notStudent)

	g.POST("/add-topic-to-problem/:id", h.addTopic, authed, notStudent)
	g.PATCH("/update-topic-to-problem/:id", h.updateTopic, authed, notStudent)
	g.DELETE("/delete-topic-from-problem/:id", h.deleteTopic, authed, notStudent)

	g.POST("/add-template-to-problem/:id", h.addTemplate, authed, notStudent)
	g.PATCH("/update-template-to-problem/:id", h.updateTemplate, authed, notStudent)
	g.DELETE("/delete-template-from-problem/:id", h.deleteTemplate, authed, notStudent)

	g.POST("/add-testcase-to-problem/:id", h.addTestCase, authed, notStudent)
	g.PATCH("/update-testcase-to-problem/:id", h.updateTestCase, authed, notStudent)
	g.DELETE("/delete-testcase-to-problem/:id", h.deleteTestCase, authed, notStudent)

	g.POST("/add-solution-to-problem/:id", h.addSolution, authed, notStudent)
	g.PATCH("/update-solution-to-problem/:id", h.updateSolution, authed, notStudent)
	g.DELETE("/delete-solution-to-problem/:id", h.deleteSolution, authed, notStudent)

	g.GET("/admin/get-user-solved-questions", h.adminGetUserSolved, authed, notStudent)

	// Authenticated (any role)
	g.GET("/get-user-solved-questions/", h.getUserSolvedQuestions, authed)
	g.GET("/get-submission/:id", h.getSubmission, authed)
}

type handler struct {
	db *db.Client
}

func (h *handler) coll(name string) *mongo.Collection { return h.db.Coll(name) }

// ---- shared fetch helpers (mirror the repeated Beanie queries in Python) ----

func (h *handler) tagsForProblem(ctx context.Context, problemID primitive.ObjectID) []string {
	cur, err := h.coll(models.CollProblemTopics).Find(ctx, bson.M{"problem_id": problemID})
	if err != nil {
		return []string{}
	}
	var topics []models.ProblemTopic
	_ = cur.All(ctx, &topics)
	tags := []string{}
	for _, t := range topics {
		tags = append(tags, t.TagName...)
	}
	return tags
}

func (h *handler) templatesForProblem(ctx context.Context, problemID primitive.ObjectID) []map[string]any {
	cur, err := h.coll(models.CollProblemTemplates).Find(ctx, bson.M{"problem_id": problemID})
	if err != nil {
		return []map[string]any{}
	}
	var templates []models.ProblemTemplate
	_ = cur.All(ctx, &templates)
	out := make([]map[string]any, 0, len(templates))
	for _, t := range templates {
		out = append(out, map[string]any{
			"id": t.ID.Hex(), "language": t.Language, "default_code": t.DefaultCode,
			"function_body": t.FunctionBody,
		})
	}
	return out
}

func (h *handler) testCasesForProblem(ctx context.Context, problemID primitive.ObjectID, testType string) []map[string]any {
	filter := bson.M{"problem_id": problemID}
	if testType != "" {
		filter["test_type"] = testType
	}
	cur, err := h.coll(models.CollProblemTestCases).Find(ctx, filter)
	if err != nil {
		return []map[string]any{}
	}
	var cases []models.ProblemTestCase
	_ = cur.All(ctx, &cases)
	out := make([]map[string]any, 0, len(cases))
	for _, tc := range cases {
		out = append(out, map[string]any{
			"id": tc.ID.Hex(), "input": tc.Input, "output": tc.Output, "test_type": tc.TestType,
		})
	}
	return out
}

func (h *handler) solutionsForProblem(ctx context.Context, problemID primitive.ObjectID) []map[string]any {
	cur, err := h.coll(models.CollProblemSolutions).Find(ctx, bson.M{"problem_id": problemID})
	if err != nil {
		return []map[string]any{}
	}
	var sols []models.ProblemSolution
	_ = cur.All(ctx, &sols)
	out := make([]map[string]any, 0, len(sols))
	for _, s := range sols {
		out = append(out, map[string]any{
			"id": s.ID.Hex(), "solution": s.Solution, "video_solution": s.VideoSolution,
		})
	}
	return out
}

func (h *handler) getProblemByID(ctx context.Context, id primitive.ObjectID) (*models.Problem, error) {
	var p models.Problem
	err := h.coll(models.CollProblems).FindOne(ctx, bson.M{"_id": id}).Decode(&p)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &p, err
}

// ---- public endpoints ----

func (h *handler) getAllProblems(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	cur, err := h.coll(models.CollProblems).Find(ctx, bson.M{})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch problems", err.Error())
	}
	var problems []models.Problem
	if err := cur.All(ctx, &problems); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch problems", err.Error())
	}

	data := make([]map[string]any, 0, len(problems))
	for _, p := range problems {
		data = append(data, map[string]any{
			"id": p.ID.Hex(), "name": p.Name, "difficulty": p.Difficulty,
			"published": p.Published, "tags": h.tagsForProblem(ctx, p.ID),
		})
	}
	return response.OK(c, "Problems fetched", data)
}

func (h *handler) getAllListedProblems(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	cur, err := h.coll(models.CollProblems).Find(ctx, bson.M{"published": models.StatusListed})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch problems", err.Error())
	}
	var problems []models.Problem
	if err := cur.All(ctx, &problems); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch problems", err.Error())
	}

	data := make([]map[string]any, 0, len(problems))
	for _, p := range problems {
		data = append(data, map[string]any{
			"id": p.ID.Hex(), "name": p.Name, "difficulty": p.Difficulty,
			"tags": h.tagsForProblem(ctx, p.ID),
		})
	}
	return response.OK(c, "Listed problems fetched", data)
}

func (h *handler) searchQuestion(c echo.Context) error {
	var body searchProblemRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	pattern := regexp.QuoteMeta(body.Query)
	filter := bson.M{"name": bson.M{"$regex": pattern, "$options": "i"}}
	cur, err := h.coll(models.CollProblems).Find(ctx, filter)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Search failed", err.Error())
	}
	var problems []models.Problem
	if err := cur.All(ctx, &problems); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Search failed", err.Error())
	}

	data := make([]map[string]any, 0, len(problems))
	for _, p := range problems {
		data = append(data, map[string]any{
			"id": p.ID.Hex(), "name": p.Name, "difficulty": p.Difficulty, "published": p.Published,
		})
	}
	return response.OK(c, "Search results", data)
}

func (h *handler) getProblem(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Problem not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	problem, err := h.getProblemByID(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch problem", err.Error())
	}
	if problem == nil {
		return response.Err(c, http.StatusNotFound, "Problem not found", "Not found")
	}

	tags := h.tagsForProblem(ctx, problem.ID)
	templates := h.templatesForProblem(ctx, problem.ID)
	testCases := h.testCasesForProblem(ctx, problem.ID, models.TestcaseExample)
	solutions := h.solutionsForProblem(ctx, problem.ID)

	return response.OK(c, "Problem fetched", map[string]any{
		"id": problem.ID.Hex(), "name": problem.Name, "description": problem.Description,
		"hints": problem.Hints, "difficulty": problem.Difficulty, "published": problem.Published,
		"tags": tags, "templates": templates, "test_cases": testCases, "solutions": solutions,
	})
}

func (h *handler) getProblemsByTag(c echo.Context) error {
	tag := c.QueryParam("tag")
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	cur, err := h.coll(models.CollProblemTopics).Find(ctx, bson.M{"tag_name": tag})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch problems", err.Error())
	}
	var topics []models.ProblemTopic
	if err := cur.All(ctx, &topics); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch problems", err.Error())
	}
	if len(topics) == 0 {
		return response.OK(c, "No problems found", []map[string]any{})
	}
	ids := make([]primitive.ObjectID, 0, len(topics))
	for _, t := range topics {
		ids = append(ids, t.ProblemID)
	}

	pcur, err := h.coll(models.CollProblems).Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch problems", err.Error())
	}
	var problems []models.Problem
	if err := pcur.All(ctx, &problems); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch problems", err.Error())
	}
	data := make([]map[string]any, 0, len(problems))
	for _, p := range problems {
		data = append(data, map[string]any{
			"id": p.ID.Hex(), "name": p.Name, "difficulty": p.Difficulty, "published": p.Published,
		})
	}
	return response.OK(c, "Problems fetched", data)
}

// ---- admin endpoints ----

func (h *handler) changeStatus(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Problem not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	problem, err := h.getProblemByID(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch problem", err.Error())
	}
	if problem == nil {
		return response.Err(c, http.StatusNotFound, "Problem not found", "Not found")
	}
	newStatus := models.StatusListed
	if problem.Published == models.StatusListed {
		newStatus = models.StatusNotListed
	}
	_, err = h.coll(models.CollProblems).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{
		"published": newStatus, "updated_at": time.Now().UTC(),
	}})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to update status", err.Error())
	}
	return response.OK(c, "Status changed", map[string]any{"published": newStatus})
}

func (h *handler) adminGetProblem(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Problem not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	problem, err := h.getProblemByID(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch problem", err.Error())
	}
	if problem == nil {
		return response.Err(c, http.StatusNotFound, "Problem not found", "Not found")
	}

	var topics []models.ProblemTopic
	if cur, err := h.coll(models.CollProblemTopics).Find(ctx, bson.M{"problem_id": problem.ID}); err == nil {
		_ = cur.All(ctx, &topics)
	}
	topicsOut := make([]map[string]any, 0, len(topics))
	for _, t := range topics {
		topicsOut = append(topicsOut, map[string]any{"id": t.ID.Hex(), "tag_name": t.TagName})
	}

	var submissions []models.ProblemSubmission
	if cur, err := h.coll(models.CollProblemSubmissions).Find(ctx, bson.M{"problem_id": problem.ID}); err == nil {
		_ = cur.All(ctx, &submissions)
	}
	subsOut := make([]map[string]any, 0, len(submissions))
	for _, s := range submissions {
		subsOut = append(subsOut, map[string]any{
			"id": s.ID.Hex(), "status": s.Status, "student_id": s.StudentID.Hex(),
		})
	}

	return response.OK(c, "Problem fetched", map[string]any{
		"id": problem.ID.Hex(), "name": problem.Name, "description": problem.Description,
		"hints": problem.Hints, "difficulty": problem.Difficulty, "published": problem.Published,
		"topics":      topicsOut,
		"templates":   h.templatesForProblem(ctx, problem.ID),
		"test_cases":  h.testCasesForProblem(ctx, problem.ID, ""),
		"solutions":   h.solutionsForProblem(ctx, problem.ID),
		"submissions": subsOut,
	})
}

func (h *handler) adminGetTestCases(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Problem not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	return response.OK(c, "Test cases fetched", h.testCasesForProblem(ctx, oid, ""))
}

func (h *handler) adminGetSolutions(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Problem not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	return response.OK(c, "Solutions fetched", h.solutionsForProblem(ctx, oid))
}

func (h *handler) adminGetSubmissions(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Problem not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	cur, err := h.coll(models.CollProblemSubmissions).Find(ctx, bson.M{"problem_id": oid})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch submissions", err.Error())
	}
	var subs []models.ProblemSubmission
	if err := cur.All(ctx, &subs); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch submissions", err.Error())
	}
	data := make([]map[string]any, 0, len(subs))
	for _, s := range subs {
		data = append(data, map[string]any{
			"id": s.ID.Hex(), "status": s.Status, "code": s.Code, "runtime": s.Runtime,
			"memory": s.Memory, "student_id": s.StudentID.Hex(),
			"submitted_at": s.SubmittedAt.Format(time.RFC3339),
		})
	}
	return response.OK(c, "Submissions fetched", data)
}

func (h *handler) adminGetTemplates(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Problem not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	return response.OK(c, "Templates fetched", h.templatesForProblem(ctx, oid))
}

func (h *handler) addProblem(c echo.Context) error {
	var body createProblemRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	user := middleware.UserFrom(c)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	if body.Difficulty == "" {
		body.Difficulty = models.ProblemEasy
	}
	if body.Hints == nil {
		body.Hints = []string{}
	}
	now := time.Now().UTC()
	problem := models.Problem{
		ID:          primitive.NewObjectID(),
		Name:        body.Name,
		Description: body.Description,
		Hints:       body.Hints,
		Difficulty:  body.Difficulty,
		CreatedBy:   &user.Type,
		CreatorType: user.Type,
		Published:   models.StatusNotListed,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if body.SectionID != nil {
		if sid, ok := dbutil.ParseID(*body.SectionID); ok {
			problem.SectionID = &sid
		}
	}
	if uid, ok := dbutil.ParseID(user.ID); ok {
		problem.UserID = &uid
	}

	if _, err := h.coll(models.CollProblems).InsertOne(ctx, problem); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to create problem", err.Error())
	}
	return response.JSON(c, http.StatusCreated, "Problem created", map[string]any{
		"id": problem.ID.Hex(), "name": problem.Name,
	}, nil)
}

func (h *handler) updateProblem(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Problem not found", "Not found")
	}
	var body updateProblemRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	set := bson.M{"updated_at": time.Now().UTC()}
	if body.Name != nil {
		set["name"] = *body.Name
	}
	if body.Description != nil {
		set["description"] = *body.Description
	}
	if body.Hints != nil {
		set["hints"] = *body.Hints
	}
	if body.Difficulty != nil {
		set["difficulty"] = *body.Difficulty
	}

	res, err := h.coll(models.CollProblems).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": set})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to update problem", err.Error())
	}
	if res.MatchedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Problem not found", "Not found")
	}
	return response.OK(c, "Problem updated", map[string]any{"id": oid.Hex()})
}

func (h *handler) deleteProblem(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Problem not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	problem, err := h.getProblemByID(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch problem", err.Error())
	}
	if problem == nil {
		return response.Err(c, http.StatusNotFound, "Problem not found", "Not found")
	}

	filter := bson.M{"problem_id": problem.ID}
	_, _ = h.coll(models.CollProblemTopics).DeleteMany(ctx, filter)
	_, _ = h.coll(models.CollProblemTemplates).DeleteMany(ctx, filter)
	_, _ = h.coll(models.CollProblemTestCases).DeleteMany(ctx, filter)
	_, _ = h.coll(models.CollProblemSolutions).DeleteMany(ctx, filter)
	_, _ = h.coll(models.CollProblemSubmissions).DeleteMany(ctx, filter)
	if _, err := h.coll(models.CollProblems).DeleteOne(ctx, bson.M{"_id": oid}); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to delete problem", err.Error())
	}
	return response.OK(c, "Problem deleted", nil)
}

// ---- topics ----

func (h *handler) addTopic(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Problem not found", "Not found")
	}
	var body addTopicRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	problem, err := h.getProblemByID(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch problem", err.Error())
	}
	if problem == nil {
		return response.Err(c, http.StatusNotFound, "Problem not found", "Not found")
	}

	now := time.Now().UTC()
	topic := models.ProblemTopic{ID: primitive.NewObjectID(), ProblemID: problem.ID, TagName: body.TagName, CreatedAt: now, UpdatedAt: now}
	if _, err := h.coll(models.CollProblemTopics).InsertOne(ctx, topic); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to add topic", err.Error())
	}
	return response.JSON(c, http.StatusCreated, "Topic added", map[string]any{"id": topic.ID.Hex()}, nil)
}

func (h *handler) updateTopic(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Topic not found", "Not found")
	}
	var body updateTopicRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	res, err := h.coll(models.CollProblemTopics).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{
		"tag_name": body.TagName, "updated_at": time.Now().UTC(),
	}})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to update topic", err.Error())
	}
	if res.MatchedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Topic not found", "Not found")
	}
	return response.OK(c, "Topic updated", nil)
}

func (h *handler) deleteTopic(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Topic not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	res, err := h.coll(models.CollProblemTopics).DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to delete topic", err.Error())
	}
	if res.DeletedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Topic not found", "Not found")
	}
	return response.OK(c, "Topic deleted", nil)
}

// ---- templates ----

func (h *handler) addTemplate(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Problem not found", "Not found")
	}
	var body addTemplateRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	problem, err := h.getProblemByID(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch problem", err.Error())
	}
	if problem == nil {
		return response.Err(c, http.StatusNotFound, "Problem not found", "Not found")
	}

	now := time.Now().UTC()
	tmpl := models.ProblemTemplate{
		ID: primitive.NewObjectID(), ProblemID: problem.ID, FunctionBody: body.FunctionBody,
		DefaultCode: body.DefaultCode, Language: body.Language, CreatedAt: now, UpdatedAt: now,
	}
	if _, err := h.coll(models.CollProblemTemplates).InsertOne(ctx, tmpl); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to add template", err.Error())
	}
	return response.JSON(c, http.StatusCreated, "Template added", map[string]any{"id": tmpl.ID.Hex()}, nil)
}

func (h *handler) updateTemplate(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Template not found", "Not found")
	}
	var body updateTemplateRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	set := bson.M{"updated_at": time.Now().UTC()}
	if body.FunctionBody != nil {
		set["function_body"] = *body.FunctionBody
	}
	if body.DefaultCode != nil {
		set["default_code"] = *body.DefaultCode
	}
	if body.Language != nil {
		set["language"] = *body.Language
	}
	res, err := h.coll(models.CollProblemTemplates).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": set})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to update template", err.Error())
	}
	if res.MatchedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Template not found", "Not found")
	}
	return response.OK(c, "Template updated", nil)
}

func (h *handler) deleteTemplate(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Template not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	res, err := h.coll(models.CollProblemTemplates).DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to delete template", err.Error())
	}
	if res.DeletedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Template not found", "Not found")
	}
	return response.OK(c, "Template deleted", nil)
}

// ---- test cases ----

func (h *handler) addTestCase(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Problem not found", "Not found")
	}
	body := addTestCaseRequest{TestType: models.TestcaseExample}
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	problem, err := h.getProblemByID(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch problem", err.Error())
	}
	if problem == nil {
		return response.Err(c, http.StatusNotFound, "Problem not found", "Not found")
	}

	now := time.Now().UTC()
	tc := models.ProblemTestCase{
		ID: primitive.NewObjectID(), TestType: body.TestType, Input: body.Input, Output: body.Output,
		ProblemID: problem.ID, CreatedAt: now, UpdatedAt: now,
	}
	if _, err := h.coll(models.CollProblemTestCases).InsertOne(ctx, tc); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to add test case", err.Error())
	}
	return response.JSON(c, http.StatusCreated, "Test case added", map[string]any{"id": tc.ID.Hex()}, nil)
}

func (h *handler) updateTestCase(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Test case not found", "Not found")
	}
	var body updateTestCaseRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	set := bson.M{"updated_at": time.Now().UTC()}
	if body.TestType != nil {
		set["test_type"] = *body.TestType
	}
	if body.Input != nil {
		set["input"] = *body.Input
	}
	if body.Output != nil {
		set["output"] = *body.Output
	}
	res, err := h.coll(models.CollProblemTestCases).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": set})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to update test case", err.Error())
	}
	if res.MatchedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Test case not found", "Not found")
	}
	return response.OK(c, "Test case updated", nil)
}

func (h *handler) deleteTestCase(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Test case not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	res, err := h.coll(models.CollProblemTestCases).DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to delete test case", err.Error())
	}
	if res.DeletedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Test case not found", "Not found")
	}
	return response.OK(c, "Test case deleted", nil)
}

// ---- solutions ----

func (h *handler) addSolution(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Problem not found", "Not found")
	}
	var body addSolutionRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	problem, err := h.getProblemByID(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch problem", err.Error())
	}
	if problem == nil {
		return response.Err(c, http.StatusNotFound, "Problem not found", "Not found")
	}

	now := time.Now().UTC()
	sol := models.ProblemSolution{
		ID: primitive.NewObjectID(), Solution: body.Solution, VideoSolution: body.VideoSolution,
		ProblemID: problem.ID, CreatedAt: now, UpdatedAt: now,
	}
	if _, err := h.coll(models.CollProblemSolutions).InsertOne(ctx, sol); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to add solution", err.Error())
	}
	return response.JSON(c, http.StatusCreated, "Solution added", map[string]any{"id": sol.ID.Hex()}, nil)
}

func (h *handler) updateSolution(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Solution not found", "Not found")
	}
	var body updateSolutionRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	set := bson.M{"updated_at": time.Now().UTC()}
	if body.Solution != nil {
		set["solution"] = *body.Solution
	}
	if body.VideoSolution != nil {
		set["video_solution"] = *body.VideoSolution
	}
	res, err := h.coll(models.CollProblemSolutions).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": set})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to update solution", err.Error())
	}
	if res.MatchedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Solution not found", "Not found")
	}
	return response.OK(c, "Solution updated", nil)
}

func (h *handler) deleteSolution(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Solution not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	res, err := h.coll(models.CollProblemSolutions).DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to delete solution", err.Error())
	}
	if res.DeletedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Solution not found", "Not found")
	}
	return response.OK(c, "Solution deleted", nil)
}

// ---- authenticated endpoints ----

func (h *handler) getUserSolvedQuestions(c echo.Context) error {
	user := middleware.UserFrom(c)
	uid, ok := dbutil.ParseID(user.ID)
	if !ok {
		return response.Err(c, http.StatusUnauthorized, "Invalid user ID", "Invalid user ID")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	cur, err := h.coll(models.CollProblemSubmissions).Find(ctx, bson.M{"student_id": uid})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch submissions", err.Error())
	}
	var subs []models.ProblemSubmission
	if err := cur.All(ctx, &subs); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch submissions", err.Error())
	}

	solved := map[string]struct{}{}
	for _, s := range subs {
		if s.Status == models.ProblemSuccess {
			solved[s.ProblemID.Hex()] = struct{}{}
		}
	}

	data := make([]map[string]any, 0, len(solved))
	for pidHex := range solved {
		oid, _ := dbutil.ParseID(pidHex)
		p, err := h.getProblemByID(ctx, oid)
		if err == nil && p != nil {
			data = append(data, map[string]any{"id": p.ID.Hex(), "name": p.Name, "difficulty": p.Difficulty})
		}
	}
	return response.OK(c, "Solved questions", data)
}

func (h *handler) adminGetUserSolved(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	cur, err := h.coll(models.CollProblemSubmissions).Find(ctx, bson.M{})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch submissions", err.Error())
	}
	var subs []models.ProblemSubmission
	if err := cur.All(ctx, &subs); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch submissions", err.Error())
	}

	userSolved := map[string]map[string]struct{}{}
	for _, s := range subs {
		key := s.StudentID.Hex()
		if userSolved[key] == nil {
			userSolved[key] = map[string]struct{}{}
		}
		if s.Status == models.ProblemSuccess {
			userSolved[key][s.ProblemID.Hex()] = struct{}{}
		}
	}

	data := make([]map[string]any, 0, len(userSolved))
	for k, v := range userSolved {
		data = append(data, map[string]any{"student_id": k, "solved_count": len(v)})
	}
	return response.OK(c, "User solved data", data)
}

func (h *handler) getSubmission(c echo.Context) error {
	user := middleware.UserFrom(c)
	uid, ok := dbutil.ParseID(user.ID)
	if !ok {
		return response.Err(c, http.StatusUnauthorized, "Invalid user ID", "Invalid user ID")
	}
	pid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.OK(c, "Submissions fetched", []map[string]any{})
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "submitted_at", Value: -1}})
	cur, err := h.coll(models.CollProblemSubmissions).Find(ctx, bson.M{"student_id": uid, "problem_id": pid}, opts)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch submissions", err.Error())
	}
	var subs []models.ProblemSubmission
	if err := cur.All(ctx, &subs); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch submissions", err.Error())
	}
	data := make([]map[string]any, 0, len(subs))
	for _, s := range subs {
		data = append(data, map[string]any{
			"id": s.ID.Hex(), "code": s.Code, "status": s.Status, "runtime": s.Runtime,
			"memory": s.Memory, "failed_test_case": s.FailedTestCase,
			"submitted_at": s.SubmittedAt.Format(time.RFC3339),
		})
	}
	return response.OK(c, "Submissions fetched", data)
}
