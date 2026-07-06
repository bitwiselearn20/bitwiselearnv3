// Package assessment ports apps/python-server/routers/assessment.py to Go.
package assessment

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bitwiselearn/go-server/internal/db"
	"github.com/bitwiselearn/go-server/internal/dbutil"
	"github.com/bitwiselearn/go-server/internal/jobs"
	"github.com/bitwiselearn/go-server/internal/middleware"
	"github.com/bitwiselearn/go-server/internal/models"
	"github.com/bitwiselearn/go-server/internal/response"
	"github.com/bitwiselearn/go-server/internal/services/piston"
	"github.com/bitwiselearn/go-server/internal/services/queue"
)

const reqTimeout = 10 * time.Second

// Deps holds the dependencies the assessment handlers need.
type Deps struct {
	DB        *db.Client
	Auth      *middleware.Auth
	Piston    *piston.Client
	Publisher *queue.Publisher
}

// Register mounts every assessment route on g (expected prefix "/api/v1/assessments").
func Register(g *echo.Group, d Deps) {
	h := &handler{db: d.DB, piston: d.Piston, pub: d.Publisher}
	authed := d.Auth.Required
	notStudent := d.Auth.NotStudent()
	adminRoles := d.Auth.RequireRoles(models.RoleSuperadmin, models.RoleAdmin, models.RoleInstitution, models.RoleVendor)
	studentOnly := d.Auth.RequireRoles(models.RoleStudent)

	g.POST("/create-assessment", h.createAssessment, authed, notStudent)
	g.GET("/get-all-assessment", h.getAllAssessments, authed, adminRoles)
	g.GET("/get-assessment-by-id/:id", h.getAssessmentByID, authed)
	g.GET("/get-assessment-by-institution/:id", h.getAssessmentByInstitution, authed)
	g.GET("/get-assessment-by-batch/:id", h.getAssessmentByBatch, authed)
	g.PUT("/update-assessment-by-id/:id", h.updateAssessment, authed, notStudent)
	g.PUT("/update-assessment-status/:id", h.updateAssessmentStatus, authed, notStudent)
	g.DELETE("/delete-assessment-by-id/:id", h.deleteAssessment, authed, notStudent)

	g.GET("/get-sections-for-assessment/:id", h.getSections, authed)
	g.POST("/add-assessment-section", h.addSection, authed, adminRoles)
	g.PUT("/update-assessment-section/:id", h.updateSection, authed, adminRoles)
	g.DELETE("/delete-assessment-section/:id", h.deleteSection, authed, adminRoles)

	g.POST("/add-assessment-question/:id", h.addQuestion, authed, adminRoles)
	g.PUT("/update-assessment-question/:id", h.updateQuestion, authed, adminRoles)
	g.DELETE("/delete-assessment-question/:id", h.deleteQuestion, authed, adminRoles)
	g.GET("/get-questions-by-sectionId/:id", h.getQuestionsBySection, authed)

	g.POST("/submit-assessment-by-id/:id", h.submitAssessment, authed, studentOnly)
	g.POST("/submit-assessment-question-by-id/:id", h.submitAssessmentQuestion, authed, studentOnly)

	g.POST("/assignment-report/:id", h.triggerReport, authed, notStudent)
}

type handler struct {
	db     *db.Client
	piston *piston.Client
	pub    *queue.Publisher
}

func (h *handler) coll(name string) *mongo.Collection { return h.db.Coll(name) }
func iso(t time.Time) string                          { return t.Format(time.RFC3339) }

func (h *handler) getAssessment(ctx context.Context, id primitive.ObjectID) (*models.Assessment, error) {
	var a models.Assessment
	err := h.coll(models.CollAssessments).FindOne(ctx, bson.M{"_id": id}).Decode(&a)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &a, err
}

func (h *handler) getSection(ctx context.Context, id primitive.ObjectID) (*models.AssessmentSection, error) {
	var s models.AssessmentSection
	err := h.coll(models.CollAssessmentSections).FindOne(ctx, bson.M{"_id": id}).Decode(&s)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &s, err
}

func (h *handler) getQuestion(ctx context.Context, id primitive.ObjectID) (*models.AssessmentQuestion, error) {
	var q models.AssessmentQuestion
	err := h.coll(models.CollAssessmentQuestions).FindOne(ctx, bson.M{"_id": id}).Decode(&q)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &q, err
}

// ================= Assessment CRUD =================

func (h *handler) createAssessment(c echo.Context) error {
	var body createAssessmentRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	batchID, ok := dbutil.ParseID(body.BatchID)
	if !ok {
		return response.Err(c, http.StatusNotFound, "Batch not found", "Not found")
	}
	user := middleware.UserFrom(c)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	now := time.Now().UTC()
	assessment := models.Assessment{
		ID: primitive.NewObjectID(), Name: body.Name, Description: body.Description,
		Instruction: body.Instruction, StartTime: body.StartTime, EndTime: body.EndTime,
		IndividualSectionTimeLimit: body.IndividualSectionTimeLimit, Status: models.AssessmentUpcoming,
		ReportStatus: models.ReportNotRequested, AutoSubmit: body.AutoSubmit, BatchID: batchID,
		CreatorID: user.ID, CreatorType: user.Type, CreatedAt: now, UpdatedAt: now,
	}
	if body.TeacherID != nil {
		if tid, ok := dbutil.ParseID(*body.TeacherID); ok {
			assessment.TeacherID = &tid
		}
	}
	if _, err := h.coll(models.CollAssessments).InsertOne(ctx, assessment); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to create assessment", err.Error())
	}
	return response.JSON(c, http.StatusCreated, "Assessment created", map[string]any{
		"id": assessment.ID.Hex(), "name": assessment.Name,
	}, nil)
}

func (h *handler) getAllAssessments(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	cur, err := h.coll(models.CollAssessments).Find(ctx, bson.M{})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch assessments", err.Error())
	}
	var assessments []models.Assessment
	if err := cur.All(ctx, &assessments); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch assessments", err.Error())
	}
	data := make([]map[string]any, 0, len(assessments))
	for _, a := range assessments {
		data = append(data, map[string]any{
			"id": a.ID.Hex(), "name": a.Name, "description": a.Description, "status": a.Status,
			"start_time": iso(a.StartTime), "end_time": iso(a.EndTime), "batch_id": a.BatchID.Hex(),
			"report_status": a.ReportStatus, "created_at": iso(a.CreatedAt),
		})
	}
	return response.OK(c, "Assessments fetched", data)
}

func (h *handler) getAssessmentByID(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Assessment not found", "Not found")
	}
	user := middleware.UserFrom(c)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	assessment, err := h.getAssessment(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch assessment", err.Error())
	}
	if assessment == nil {
		return response.Err(c, http.StatusNotFound, "Assessment not found", "Not found")
	}

	var teacherID any
	if assessment.TeacherID != nil {
		teacherID = assessment.TeacherID.Hex()
	}
	data := map[string]any{
		"id": assessment.ID.Hex(), "name": assessment.Name, "description": assessment.Description,
		"instruction": assessment.Instruction, "status": assessment.Status,
		"start_time": iso(assessment.StartTime), "end_time": iso(assessment.EndTime),
		"individual_section_time_limit": assessment.IndividualSectionTimeLimit,
		"auto_submit":                   assessment.AutoSubmit, "batch_id": assessment.BatchID.Hex(),
		"report": assessment.Report, "report_status": assessment.ReportStatus, "teacher_id": teacherID,
	}

	if user.Type == models.RoleStudent {
		studentID, _ := dbutil.ParseID(user.ID)
		var sub models.AssessmentSubmission
		err := h.coll(models.CollAssessmentSubmissions).FindOne(ctx, bson.M{
			"assessment_id": assessment.ID, "student_id": studentID,
		}).Decode(&sub)
		if err == nil {
			data["has_submitted"] = sub.IsSubmitted
			data["submission_id"] = sub.ID.Hex()
		} else {
			data["has_submitted"] = false
			data["submission_id"] = nil
		}
	}

	return response.OK(c, "Assessment fetched", data)
}

func (h *handler) getAssessmentByInstitution(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.OK(c, "No assessments", []map[string]any{})
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	cur, err := h.coll(models.CollBatches).Find(ctx, bson.M{"institution_id": oid})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch assessments", err.Error())
	}
	var batches []models.Batch
	_ = cur.All(ctx, &batches)
	if len(batches) == 0 {
		return response.OK(c, "No assessments", []map[string]any{})
	}
	batchIDs := make([]primitive.ObjectID, 0, len(batches))
	for _, b := range batches {
		batchIDs = append(batchIDs, b.ID)
	}

	acur, err := h.coll(models.CollAssessments).Find(ctx, bson.M{"batch_id": bson.M{"$in": batchIDs}})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch assessments", err.Error())
	}
	var assessments []models.Assessment
	_ = acur.All(ctx, &assessments)
	data := make([]map[string]any, 0, len(assessments))
	for _, a := range assessments {
		data = append(data, map[string]any{
			"id": a.ID.Hex(), "name": a.Name, "status": a.Status, "batch_id": a.BatchID.Hex(),
			"start_time": iso(a.StartTime), "end_time": iso(a.EndTime),
		})
	}
	return response.OK(c, "Assessments fetched", data)
}

func (h *handler) getAssessmentByBatch(c echo.Context) error {
	batchID, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.OK(c, "Assessments fetched", []map[string]any{})
	}
	user := middleware.UserFrom(c)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	cur, err := h.coll(models.CollAssessments).Find(ctx, bson.M{"batch_id": batchID})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch assessments", err.Error())
	}
	var assessments []models.Assessment
	_ = cur.All(ctx, &assessments)

	data := []map[string]any{}
	if user.Type == models.RoleStudent {
		studentID, _ := dbutil.ParseID(user.ID)
		for _, a := range assessments {
			if a.Status != models.AssessmentLive {
				continue
			}
			var sub models.AssessmentSubmission
			err := h.coll(models.CollAssessmentSubmissions).FindOne(ctx, bson.M{
				"assessment_id": a.ID, "student_id": studentID,
			}).Decode(&sub)
			canAccess := err != nil || !sub.IsSubmitted
			data = append(data, map[string]any{
				"id": a.ID.Hex(), "name": a.Name, "status": a.Status,
				"start_time": iso(a.StartTime), "end_time": iso(a.EndTime), "canAccessTest": canAccess,
			})
		}
	} else {
		for _, a := range assessments {
			data = append(data, map[string]any{
				"id": a.ID.Hex(), "name": a.Name, "status": a.Status,
				"start_time": iso(a.StartTime), "end_time": iso(a.EndTime),
				"batch_id": a.BatchID.Hex(), "report_status": a.ReportStatus,
			})
		}
	}
	return response.OK(c, "Assessments fetched", data)
}

func (h *handler) updateAssessment(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Assessment not found", "Not found")
	}
	var body updateAssessmentRequest
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
	if body.Instruction != nil {
		set["instruction"] = *body.Instruction
	}
	if body.StartTime != nil {
		set["start_time"] = *body.StartTime
	}
	if body.EndTime != nil {
		set["end_time"] = *body.EndTime
	}
	if body.IndividualSectionTimeLimit != nil {
		set["individual_section_time_limit"] = *body.IndividualSectionTimeLimit
	}
	if body.AutoSubmit != nil {
		set["auto_submit"] = *body.AutoSubmit
	}
	res, err := h.coll(models.CollAssessments).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": set})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to update assessment", err.Error())
	}
	if res.MatchedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Assessment not found", "Not found")
	}
	return response.OK(c, "Assessment updated", map[string]any{"id": oid.Hex()})
}

func (h *handler) updateAssessmentStatus(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Assessment not found", "Not found")
	}
	var body updateAssessmentStatusRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	res, err := h.coll(models.CollAssessments).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{
		"status": body.Status, "updated_at": time.Now().UTC(),
	}})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to update status", err.Error())
	}
	if res.MatchedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Assessment not found", "Not found")
	}
	return response.OK(c, "Status updated", map[string]any{"status": body.Status})
}

func (h *handler) deleteAssessment(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Assessment not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	assessment, err := h.getAssessment(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch assessment", err.Error())
	}
	if assessment == nil {
		return response.Err(c, http.StatusNotFound, "Assessment not found", "Not found")
	}

	var sectionIDs []primitive.ObjectID
	if cur, err := h.coll(models.CollAssessmentSections).Find(ctx, bson.M{"assessment_id": oid}); err == nil {
		var sections []models.AssessmentSection
		_ = cur.All(ctx, &sections)
		for _, s := range sections {
			sectionIDs = append(sectionIDs, s.ID)
		}
	}
	_, _ = h.coll(models.CollAssessmentSections).DeleteMany(ctx, bson.M{"assessment_id": oid})
	if len(sectionIDs) > 0 {
		_, _ = h.coll(models.CollAssessmentQuestions).DeleteMany(ctx, bson.M{"section_id": bson.M{"$in": sectionIDs}})
	}
	_, _ = h.coll(models.CollAssessmentSubmissions).DeleteMany(ctx, bson.M{"assessment_id": oid})
	_, _ = h.coll(models.CollAssessmentQuestionSubmissions).DeleteMany(ctx, bson.M{"assessment_id": oid})
	if _, err := h.coll(models.CollAssessments).DeleteOne(ctx, bson.M{"_id": oid}); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to delete assessment", err.Error())
	}
	return response.OK(c, "Assessment deleted", nil)
}

// ================= Sections =================

func (h *handler) getSections(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.OK(c, "Sections fetched", []map[string]any{})
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	cur, err := h.coll(models.CollAssessmentSections).Find(ctx, bson.M{"assessment_id": oid})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch sections", err.Error())
	}
	var sections []models.AssessmentSection
	if err := cur.All(ctx, &sections); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch sections", err.Error())
	}
	data := make([]map[string]any, 0, len(sections))
	for _, s := range sections {
		data = append(data, map[string]any{
			"id": s.ID.Hex(), "name": s.Name, "marks_per_question": s.MarksPerQuestion,
			"assessment_type": s.AssessmentType, "assessment_id": s.AssessmentID.Hex(),
		})
	}
	return response.OK(c, "Sections fetched", data)
}

func (h *handler) addSection(c echo.Context) error {
	var body createAssessmentSectionRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	assessmentID, ok := dbutil.ParseID(body.AssessmentID)
	if !ok {
		return response.Err(c, http.StatusNotFound, "Assessment not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	now := time.Now().UTC()
	section := models.AssessmentSection{
		ID: primitive.NewObjectID(), Name: body.Name, MarksPerQuestion: body.MarksPerQuestion,
		AssessmentType: body.AssessmentType, AssessmentID: assessmentID, CreatedAt: now, UpdatedAt: now,
	}
	if _, err := h.coll(models.CollAssessmentSections).InsertOne(ctx, section); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to create section", err.Error())
	}
	return response.JSON(c, http.StatusCreated, "Section created", map[string]any{"id": section.ID.Hex()}, nil)
}

func (h *handler) updateSection(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Section not found", "Not found")
	}
	var body updateAssessmentSectionRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	set := bson.M{"updated_at": time.Now().UTC()}
	if body.Name != nil {
		set["name"] = *body.Name
	}
	if body.MarksPerQuestion != nil {
		set["marks_per_question"] = *body.MarksPerQuestion
	}
	if body.AssessmentType != nil {
		set["assessment_type"] = *body.AssessmentType
	}
	res, err := h.coll(models.CollAssessmentSections).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": set})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to update section", err.Error())
	}
	if res.MatchedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Section not found", "Not found")
	}
	return response.OK(c, "Section updated", nil)
}

func (h *handler) deleteSection(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Section not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	_, _ = h.coll(models.CollAssessmentQuestions).DeleteMany(ctx, bson.M{"section_id": oid})
	res, err := h.coll(models.CollAssessmentSections).DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to delete section", err.Error())
	}
	if res.DeletedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Section not found", "Not found")
	}
	return response.OK(c, "Section deleted", nil)
}

// ================= Questions =================

func (h *handler) addQuestion(c echo.Context) error {
	sectionID, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Section not found", "Not found")
	}
	var body addAssessmentQuestionRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	section, err := h.getSection(ctx, sectionID)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch section", err.Error())
	}
	if section == nil {
		return response.Err(c, http.StatusNotFound, "Section not found", "Not found")
	}

	now := time.Now().UTC()
	q := models.AssessmentQuestion{
		ID: primitive.NewObjectID(), Question: body.Question, Options: body.Options,
		CorrectOption: body.CorrectOption, SectionID: section.ID, MaxMarks: body.MaxMarks,
		CreatedAt: now, UpdatedAt: now,
	}
	if body.ProblemID != nil {
		if pid, ok := dbutil.ParseID(*body.ProblemID); ok {
			q.ProblemID = &pid
		}
	}
	if _, err := h.coll(models.CollAssessmentQuestions).InsertOne(ctx, q); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to add question", err.Error())
	}
	return response.JSON(c, http.StatusCreated, "Question added", map[string]any{"id": q.ID.Hex()}, nil)
}

func (h *handler) updateQuestion(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Question not found", "Not found")
	}
	var body updateAssessmentQuestionRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	set := bson.M{"updated_at": time.Now().UTC()}
	if body.Question != nil {
		set["question"] = *body.Question
	}
	if body.Options != nil {
		set["options"] = *body.Options
	}
	if body.CorrectOption != nil {
		set["correct_option"] = *body.CorrectOption
	}
	if body.MaxMarks != nil {
		set["max_marks"] = *body.MaxMarks
	}
	res, err := h.coll(models.CollAssessmentQuestions).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": set})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to update question", err.Error())
	}
	if res.MatchedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Question not found", "Not found")
	}
	return response.OK(c, "Question updated", nil)
}

func (h *handler) deleteQuestion(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Question not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	res, err := h.coll(models.CollAssessmentQuestions).DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to delete question", err.Error())
	}
	if res.DeletedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Question not found", "Not found")
	}
	return response.OK(c, "Question deleted", nil)
}

func (h *handler) getQuestionsBySection(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.OK(c, "Questions fetched", []map[string]any{})
	}
	user := middleware.UserFrom(c)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	cur, err := h.coll(models.CollAssessmentQuestions).Find(ctx, bson.M{"section_id": oid})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch questions", err.Error())
	}
	var questions []models.AssessmentQuestion
	if err := cur.All(ctx, &questions); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch questions", err.Error())
	}
	data := make([]map[string]any, 0, len(questions))
	for _, q := range questions {
		var correctOption any
		if user.Type != models.RoleStudent {
			correctOption = q.CorrectOption
		}
		var problemID any
		if q.ProblemID != nil {
			problemID = q.ProblemID.Hex()
		}
		data = append(data, map[string]any{
			"id": q.ID.Hex(), "question": q.Question, "options": q.Options,
			"correct_option": correctOption, "problem_id": problemID,
			"max_marks": q.MaxMarks, "section_id": q.SectionID.Hex(),
		})
	}
	return response.OK(c, "Questions fetched", data)
}

// ================= Submissions =================

func (h *handler) submitAssessment(c echo.Context) error {
	assessmentID, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Assessment not found", "Not found")
	}
	var body submitAssessmentRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	user := middleware.UserFrom(c)
	studentID, _ := dbutil.ParseID(user.ID)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	assessment, err := h.getAssessment(ctx, assessmentID)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch assessment", err.Error())
	}
	if assessment == nil {
		return response.Err(c, http.StatusNotFound, "Assessment not found", "Not found")
	}

	qcur, err := h.coll(models.CollAssessmentQuestionSubmissions).Find(ctx, bson.M{
		"assessment_id": assessment.ID, "student_id": studentID,
	})
	var qSubs []models.AssessmentQuestionSubmission
	if err == nil {
		_ = qcur.All(ctx, &qSubs)
	}
	totalMarks := 0.0
	for _, qs := range qSubs {
		totalMarks += qs.MarksObtained
	}

	now := time.Now().UTC()
	var submission models.AssessmentSubmission
	err = h.coll(models.CollAssessmentSubmissions).FindOne(ctx, bson.M{
		"assessment_id": assessment.ID, "student_id": studentID,
	}).Decode(&submission)

	if err == nil {
		_, uerr := h.coll(models.CollAssessmentSubmissions).UpdateOne(ctx, bson.M{"_id": submission.ID}, bson.M{"$set": bson.M{
			"is_submitted": true, "submitted_at": now, "total_marks": totalMarks,
			"tab_switch_count": body.TabSwitchCount, "proctoring_status": body.ProctoringStatus,
			"updated_at": now,
		}})
		if uerr != nil {
			return response.Err(c, http.StatusInternalServerError, "Failed to submit assessment", uerr.Error())
		}
	} else {
		submission = models.AssessmentSubmission{
			ID: primitive.NewObjectID(), AssessmentID: assessment.ID, StudentID: studentID,
			StudentIP: body.StudentIP, IsSubmitted: true, SubmittedAt: &now, TotalMarks: &totalMarks,
			TabSwitchCount: body.TabSwitchCount, ProctoringStatus: body.ProctoringStatus,
			StartedAt: now, CreatedAt: now, UpdatedAt: now,
		}
		if _, ierr := h.coll(models.CollAssessmentSubmissions).InsertOne(ctx, submission); ierr != nil {
			return response.Err(c, http.StatusInternalServerError, "Failed to submit assessment", ierr.Error())
		}
	}

	return response.OK(c, "Assessment submitted", map[string]any{
		"submission_id": submission.ID.Hex(), "total_marks": totalMarks,
	})
}

func (h *handler) submitAssessmentQuestion(c echo.Context) error {
	assessmentID, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Assessment not found", "Not found")
	}
	var body submitAssessmentQuestionRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	user := middleware.UserFrom(c)
	studentID, _ := dbutil.ParseID(user.ID)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	assessment, err := h.getAssessment(ctx, assessmentID)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch assessment", err.Error())
	}
	if assessment == nil {
		return response.Err(c, http.StatusNotFound, "Assessment not found", "Not found")
	}
	qID, ok := dbutil.ParseID(body.QuestionID)
	if !ok {
		return response.Err(c, http.StatusNotFound, "Question not found", "Not found")
	}
	question, err := h.getQuestion(ctx, qID)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch question", err.Error())
	}
	if question == nil {
		return response.Err(c, http.StatusNotFound, "Question not found", "Not found")
	}
	section, err := h.getSection(ctx, question.SectionID)
	if err != nil || section == nil || section.AssessmentID.Hex() != assessment.ID.Hex() {
		return response.Err(c, http.StatusBadRequest, "Question does not belong to this assessment", "Invalid assessment question")
	}

	now := time.Now().UTC()
	var submission models.AssessmentSubmission
	err = h.coll(models.CollAssessmentSubmissions).FindOne(ctx, bson.M{
		"assessment_id": assessment.ID, "student_id": studentID,
	}).Decode(&submission)
	if err != nil {
		submission = models.AssessmentSubmission{
			ID: primitive.NewObjectID(), AssessmentID: assessment.ID, StudentID: studentID,
			StudentIP: "", ProctoringStatus: "NOT_CHEATED", StartedAt: now, CreatedAt: now, UpdatedAt: now,
		}
		if _, ierr := h.coll(models.CollAssessmentSubmissions).InsertOne(ctx, submission); ierr != nil {
			return response.Err(c, http.StatusInternalServerError, "Failed to submit question", ierr.Error())
		}
	}

	marks := 0.0
	var answerStr *string = body.Answer

	switch section.AssessmentType {
	case models.AssessmentNoCode:
		if body.Answer != nil && question.CorrectOption != nil && *body.Answer == *question.CorrectOption {
			marks = float64(question.MaxMarks)
		}
	case models.AssessmentCode:
		if body.Code != nil && body.Language != nil && question.ProblemID != nil {
			result, err := h.gradeCode(ctx, *question.ProblemID, *body.Code, *body.Language)
			if err != nil {
				return response.Err(c, http.StatusServiceUnavailable, "Code execution failed", err.Error())
			}
			if result {
				marks = float64(question.MaxMarks)
			}
		}
		answerStr = body.Code
	}

	existing := h.coll(models.CollAssessmentQuestionSubmissions).FindOne(ctx, bson.M{
		"question_id": question.ID, "assessment_id": assessment.ID, "student_id": studentID,
	})
	if existing.Err() == nil {
		var ex models.AssessmentQuestionSubmission
		_ = existing.Decode(&ex)
		_, uerr := h.coll(models.CollAssessmentQuestionSubmissions).UpdateOne(ctx, bson.M{"_id": ex.ID}, bson.M{"$set": bson.M{
			"answer": answerStr, "marks_obtained": marks, "assessment_submission_id": submission.ID,
			"updated_at": now,
		}})
		if uerr != nil {
			return response.Err(c, http.StatusInternalServerError, "Failed to submit question", uerr.Error())
		}
	} else {
		qs := models.AssessmentQuestionSubmission{
			ID: primitive.NewObjectID(), QuestionID: question.ID, AssessmentID: assessment.ID,
			StudentID: studentID, Answer: answerStr, MarksObtained: marks,
			AssessmentSubmissionID: &submission.ID, CreatedAt: now, UpdatedAt: now,
		}
		if _, ierr := h.coll(models.CollAssessmentQuestionSubmissions).InsertOne(ctx, qs); ierr != nil {
			return response.Err(c, http.StatusInternalServerError, "Failed to submit question", ierr.Error())
		}
	}

	return response.OK(c, "Question submitted", map[string]any{"marks_obtained": marks})
}

// gradeCode composes the user's code with the problem's template, runs it
// against every test case, and reports whether all passed. Ports the
// _find_problem_template / _compose_full_code / execute_code loop.
func (h *handler) gradeCode(ctx context.Context, problemID primitive.ObjectID, code, language string) (bool, error) {
	var problem models.Problem
	if err := h.coll(models.CollProblems).FindOne(ctx, bson.M{"_id": problemID}).Decode(&problem); err != nil {
		return false, nil
	}

	var template models.ProblemTemplate
	err := h.coll(models.CollProblemTemplates).FindOne(ctx, bson.M{
		"problem_id": problem.ID, "language": language,
	}).Decode(&template)
	if err != nil {
		normalized := piston.NormalizeLanguage(language)
		if normalized != language {
			err = h.coll(models.CollProblemTemplates).FindOne(ctx, bson.M{
				"problem_id": problem.ID, "language": normalized,
			}).Decode(&template)
		}
	}
	fullCode := code
	if err == nil && template.FunctionBody != "" {
		if strings.Contains(template.FunctionBody, "_solution_") {
			fullCode = strings.ReplaceAll(template.FunctionBody, "_solution_", code)
		} else {
			fullCode = code + "\n" + template.FunctionBody
		}
	}

	cur, err := h.coll(models.CollProblemTestCases).Find(ctx, bson.M{"problem_id": problem.ID})
	if err != nil {
		return false, err
	}
	var testCases []models.ProblemTestCase
	if err := cur.All(ctx, &testCases); err != nil {
		return false, err
	}

	for _, tc := range testCases {
		result := h.piston.Execute(ctx, language, fullCode, tc.Input)
		if errMsg := result.Error(); errMsg != "" {
			details := result.Details()
			if details != "" {
				errMsg = errMsg + ": " + details
			}
			return false, &gradeError{msg: errMsg}
		}
		actual := strings.TrimSpace(strings.ReplaceAll(result.Stdout(), "\r\n", "\n"))
		expected := strings.TrimSpace(strings.ReplaceAll(tc.Output, "\r\n", "\n"))
		if actual != expected {
			return false, nil
		}
	}
	return true, nil
}

type gradeError struct{ msg string }

func (e *gradeError) Error() string { return e.msg }

// ================= Report =================

func (h *handler) triggerReport(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Assessment not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	assessment, err := h.getAssessment(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch assessment", err.Error())
	}
	if assessment == nil {
		return response.Err(c, http.StatusNotFound, "Assessment not found", "Not found")
	}

	_, err = h.coll(models.CollAssessments).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{
		"report_status": models.ReportProcessing, "updated_at": time.Now().UTC(),
	}})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to trigger report", err.Error())
	}

	if h.pub != nil {
		_ = h.pub.Publish(ctx, jobs.AssessmentReportQueue, jobs.AssessmentReportJob{AssessmentID: oid.Hex()})
	}

	return response.OK(c, "Report generation triggered", map[string]any{"report_status": models.ReportProcessing})
}
