// Package bulkupload ports apps/python-server/routers/bulk_upload.py to Go.
//
// Two endpoints (batches, cloud-info) reference model fields that don't
// exist on the current Batch/Student structs in the legacy Python code
// (Batch.name, Student.cloud_id/cloud_provider) — those routes would raise
// AttributeError in production today. This port maps them onto the fields
// that actually exist (Batch.batchname; Student.cloudname/cloud_platform)
// instead of reproducing the crash, noted inline at each spot.
package bulkupload

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bitwiselearn/go-server/internal/auth"
	"github.com/bitwiselearn/go-server/internal/db"
	"github.com/bitwiselearn/go-server/internal/dbutil"
	"github.com/bitwiselearn/go-server/internal/middleware"
	"github.com/bitwiselearn/go-server/internal/models"
	"github.com/bitwiselearn/go-server/internal/response"
)

const reqTimeout = 30 * time.Second

// Deps holds the dependencies the bulk-upload handlers need.
type Deps struct {
	DB   *db.Client
	Auth *middleware.Auth
}

// Register mounts every bulk-upload route on g (expected prefix "/api/v1/bulk-upload").
func Register(g *echo.Group, d Deps) {
	h := &handler{db: d.DB}
	authed := d.Auth.Required
	notStudent := d.Auth.NotStudent()
	adminOnly := d.Auth.AdminOnly()

	g.POST("/students/:id", h.bulkUploadStudents, authed, notStudent)
	g.POST("/batches/:id", h.bulkUploadBatches, authed, notStudent)
	g.POST("/testcases/:id", h.bulkUploadTestcases, authed, adminOnly)
	g.POST("/cloud-info/", h.bulkUploadCloudInfo, authed, adminOnly)
	g.POST("/assignment/:id", h.bulkUploadAssignment, authed, notStudent)
	g.POST("/assessment/:id", h.bulkUploadAssessment, authed, notStudent)
}

type handler struct {
	db *db.Client
}

func (h *handler) coll(name string) *mongo.Collection { return h.db.Coll(name) }

var emailRe = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

func looksLikeEmail(v string) bool { return emailRe.MatchString(v) }

// generatePassword mirrors _generate_password: 10 chars from
// ascii_letters+digits chosen with a CSPRNG.
func generatePassword() string {
	const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	buf := make([]byte, 10)
	_, _ = rand.Read(buf)
	out := make([]byte, 10)
	for i, b := range buf {
		out[i] = alphabet[int(b)%len(alphabet)]
	}
	return string(out)
}

func cell(row []string, i int) string {
	if i < len(row) {
		return strings.TrimSpace(row[i])
	}
	return ""
}

// readRows reads the uploaded "file" field as CSV (if the filename ends in
// .csv) or the first sheet of an Excel workbook otherwise, matching the
// legacy handler's format dispatch.
func readRows(c echo.Context) ([][]string, error) {
	fh, err := c.FormFile("file")
	if err != nil {
		return nil, err
	}
	f, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(strings.ToLower(fh.Filename), ".csv") {
		content = bytes.TrimPrefix(content, []byte{0xEF, 0xBB, 0xBF}) // utf-8-sig BOM
		return csv.NewReader(bytes.NewReader(content)).ReadAll()
	}
	return readXLSXRows(content)
}

func readXLSXRows(content []byte) ([][]string, error) {
	xl, err := excelize.OpenReader(bytes.NewReader(content))
	if err != nil {
		return nil, err
	}
	defer xl.Close()
	sheets := xl.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("workbook has no sheets")
	}
	return xl.GetRows(sheets[0])
}

// ---- students ----

func (h *handler) bulkUploadStudents(c echo.Context) error {
	batchID, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Batch not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	var batch models.Batch
	if err := h.coll(models.CollBatches).FindOne(ctx, bson.M{"_id": batchID}).Decode(&batch); err != nil {
		return response.Err(c, http.StatusNotFound, "Batch not found", "Not found")
	}

	rows, err := readRows(c)
	if err != nil {
		return response.Err(c, http.StatusBadRequest, "Failed to parse file", err.Error())
	}

	created := 0
	errs := []string{}
	for idx, row := range rows {
		rowNum := idx + 1
		if idx == 0 {
			continue // header
		}
		if len(row) == 0 || cell(row, 0) == "" {
			continue
		}
		name := cell(row, 0)
		col2, col3, col4 := cell(row, 1), cell(row, 2), cell(row, 3)

		var email, rollNumber, batchName string
		if looksLikeEmail(col2) {
			email = col2
			rollNumber = fmt.Sprintf("AUTO-%d", rowNum)
		} else {
			rollNumber = col2
			email = col3
			batchName = col4
		}

		if rollNumber == "" {
			errs = append(errs, fmt.Sprintf("Row %d: Missing roll number", rowNum))
			continue
		}
		if email == "" {
			errs = append(errs, fmt.Sprintf("Row %d: Missing email", rowNum))
			continue
		}
		if !looksLikeEmail(email) {
			errs = append(errs, fmt.Sprintf("Row %d: Invalid email '%s'", rowNum, email))
			continue
		}
		if batchName != "" && !strings.EqualFold(strings.TrimSpace(batchName), strings.TrimSpace(batch.BatchName)) {
			errs = append(errs, fmt.Sprintf("Row %d: Batch name '%s' does not match current batch '%s'", rowNum, batchName, batch.BatchName))
			continue
		}

		if h.coll(models.CollStudents).FindOne(ctx, bson.M{"email": email}).Err() == nil {
			errs = append(errs, fmt.Sprintf("Row %d: Email %s already exists", rowNum, email))
			continue
		}
		if h.coll(models.CollStudents).FindOne(ctx, bson.M{"roll_number": rollNumber, "batch_id": batchID}).Err() == nil {
			errs = append(errs, fmt.Sprintf("Row %d: Roll number %s already exists in this batch", rowNum, rollNumber))
			continue
		}

		hashed, err := auth.HashPassword(generatePassword())
		if err != nil {
			errs = append(errs, fmt.Sprintf("Row %d: Failed to hash password", rowNum))
			continue
		}
		now := time.Now().UTC()
		student := models.Student{
			ID: primitive.NewObjectID(), Name: name, RollNumber: rollNumber, Email: email,
			LoginPassword: hashed, CloudPlatform: "AWS", InstituteID: batch.InstitutionID,
			BatchID: batchID, CreatedAt: now, UpdatedAt: now,
		}
		if _, err := h.coll(models.CollStudents).InsertOne(ctx, student); err != nil {
			errs = append(errs, fmt.Sprintf("Row %d: Failed to create student: %v", rowNum, err))
			continue
		}
		created++
	}

	return response.OK(c, "Students uploaded", map[string]any{"created": created, "errors": errs})
}

// ---- batches ----
//
// The legacy handler references Batch.name / Batch(name=...) which doesn't
// exist on the model (it's `batchname`) — this maps row[0] onto batchname
// and leaves branch/batch_end_year blank rather than reproducing the crash.

func (h *handler) bulkUploadBatches(c echo.Context) error {
	instID, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Institution not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	if h.coll(models.CollInstitutions).FindOne(ctx, bson.M{"_id": instID}).Err() != nil {
		return response.Err(c, http.StatusNotFound, "Institution not found", "Not found")
	}

	rows, err := readRows(c)
	if err != nil {
		return response.Err(c, http.StatusBadRequest, "Failed to parse file", err.Error())
	}

	created := 0
	errs := []string{}
	for idx, row := range rows {
		rowNum := idx + 1
		if idx == 0 {
			continue
		}
		if len(row) == 0 || cell(row, 0) == "" {
			continue
		}
		name := cell(row, 0)

		if h.coll(models.CollBatches).FindOne(ctx, bson.M{"batchname": name, "institution_id": instID}).Err() == nil {
			errs = append(errs, fmt.Sprintf("Row %d: Batch '%s' already exists", rowNum, name))
			continue
		}

		now := time.Now().UTC()
		batch := models.Batch{
			ID: primitive.NewObjectID(), BatchName: name, InstitutionID: instID,
			CreatedAt: now, UpdatedAt: now,
		}
		if _, err := h.coll(models.CollBatches).InsertOne(ctx, batch); err != nil {
			errs = append(errs, fmt.Sprintf("Row %d: Failed to create batch: %v", rowNum, err))
			continue
		}
		created++
	}

	return response.OK(c, "Batches uploaded", map[string]any{"created": created, "errors": errs})
}

// ---- test cases ----

func (h *handler) bulkUploadTestcases(c echo.Context) error {
	problemID, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Problem not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	rows, err := readRows(c)
	if err != nil {
		return response.Err(c, http.StatusBadRequest, "Failed to parse file", err.Error())
	}

	created := 0
	errs := []string{}
	for idx, row := range rows {
		if idx == 0 {
			continue
		}
		if len(row) == 0 || cell(row, 0) == "" {
			continue
		}
		input := cell(row, 0)
		output := cell(row, 1)
		tcType := strings.ToUpper(cell(row, 2))
		if tcType == "" {
			tcType = models.TestcaseHidden
		}
		if tcType != models.TestcaseExample && tcType != models.TestcaseHidden {
			tcType = models.TestcaseHidden
		}

		now := time.Now().UTC()
		tc := models.ProblemTestCase{
			ID: primitive.NewObjectID(), Input: input, Output: output, TestType: tcType,
			ProblemID: problemID, CreatedAt: now, UpdatedAt: now,
		}
		if _, err := h.coll(models.CollProblemTestCases).InsertOne(ctx, tc); err != nil {
			errs = append(errs, fmt.Sprintf("Row %d: Failed to create test case: %v", idx+1, err))
			continue
		}
		created++
	}

	return response.OK(c, "Testcases uploaded", map[string]any{"created": created, "errors": errs})
}

// ---- cloud info ----
//
// The legacy handler sets student.cloud_id / student.cloud_provider, which
// don't exist on the Student model — mapped onto cloudname / cloud_platform
// (the closest existing fields) instead of reproducing the crash.

func (h *handler) bulkUploadCloudInfo(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	rows, err := readRows(c)
	if err != nil {
		return response.Err(c, http.StatusBadRequest, "Failed to parse file", err.Error())
	}

	updated := 0
	errs := []string{}
	for idx, row := range rows {
		rowNum := idx + 1
		if idx == 0 {
			continue
		}
		if len(row) == 0 || cell(row, 0) == "" {
			continue
		}
		email := cell(row, 0)
		cloudID := cell(row, 1)
		cloudProvider := strings.ToUpper(cell(row, 2))

		res, err := h.coll(models.CollStudents).UpdateOne(ctx, bson.M{"email": email}, bson.M{"$set": bson.M{
			"cloudname": cloudID, "cloud_platform": cloudProvider, "updated_at": time.Now().UTC(),
		}})
		if err != nil || res.MatchedCount == 0 {
			errs = append(errs, fmt.Sprintf("Row %d: Student %s not found", rowNum, email))
			continue
		}
		updated++
	}

	return response.OK(c, "Cloud info updated", map[string]any{"updated": updated, "errors": errs})
}

// ---- assignment preview (parse-only, no ID lookup needed) ----

func (h *handler) bulkUploadAssignment(c echo.Context) error {
	rows, err := readRows(c)
	if err != nil {
		return response.Err(c, http.StatusBadRequest, "Failed to parse file", err.Error())
	}

	rowsData := []map[string]any{}
	for idx, row := range rows {
		if idx == 0 {
			continue
		}
		if len(row) == 0 || cell(row, 0) == "" {
			continue
		}
		options := []string{}
		for i := 1; i <= 4; i++ {
			if v := cell(row, i); v != "" {
				options = append(options, v)
			}
		}
		rowsData = append(rowsData, map[string]any{
			"question": cell(row, 0), "options": options, "correct_option": cell(row, 5),
		})
	}

	return response.OK(c, "Assignment data parsed", map[string]any{"rows": rowsData, "count": len(rowsData)})
}

// ---- assessment questions ----

func (h *handler) bulkUploadAssessment(c echo.Context) error {
	sectionID, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Section not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	if h.coll(models.CollAssessmentSections).FindOne(ctx, bson.M{"_id": sectionID}).Err() != nil {
		return response.Err(c, http.StatusNotFound, "Section not found", "Not found")
	}

	rows, err := readRows(c)
	if err != nil {
		return response.Err(c, http.StatusBadRequest, "Failed to parse file", err.Error())
	}

	created := 0
	for idx, row := range rows {
		if idx == 0 {
			continue
		}
		if len(row) == 0 || cell(row, 0) == "" {
			continue
		}
		options := []string{}
		for i := 1; i <= 4; i++ {
			if v := cell(row, i); v != "" {
				options = append(options, v)
			}
		}
		var correctOption *string
		if v := cell(row, 5); v != "" {
			correctOption = &v
		}
		maxMarks := 1
		if v := cell(row, 6); v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				maxMarks = n
			}
		}
		question := cell(row, 0)

		now := time.Now().UTC()
		q := models.AssessmentQuestion{
			ID: primitive.NewObjectID(), Question: &question, Options: options,
			CorrectOption: correctOption, SectionID: sectionID, MaxMarks: maxMarks,
			CreatedAt: now, UpdatedAt: now,
		}
		if _, err := h.coll(models.CollAssessmentQuestions).InsertOne(ctx, q); err != nil {
			continue
		}
		created++
	}

	return response.OK(c, "Assessment questions uploaded", map[string]any{"created": created})
}
