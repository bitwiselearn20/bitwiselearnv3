// Package report ports apps/python-server/routers/report.py to Go.
package report

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/bitwiselearn/go-server/internal/db"
	"github.com/bitwiselearn/go-server/internal/dbutil"
	"github.com/bitwiselearn/go-server/internal/jobs"
	"github.com/bitwiselearn/go-server/internal/middleware"
	"github.com/bitwiselearn/go-server/internal/models"
	"github.com/bitwiselearn/go-server/internal/response"
	"github.com/bitwiselearn/go-server/internal/services/queue"
)

const reqTimeout = 10 * time.Second

// Deps holds the dependencies the report handlers need.
type Deps struct {
	DB        *db.Client
	Auth      *middleware.Auth
	Publisher *queue.Publisher
}

// Register mounts every report route on g (expected prefix "/api/v1/reports").
func Register(g *echo.Group, d Deps) {
	h := &handler{db: d.DB, pub: d.Publisher}
	authed := d.Auth.Required
	adminOnly := d.Auth.AdminOnly()
	notStudent := d.Auth.NotStudent()

	g.GET("/get-stats-count", h.getStatsCount, authed, adminOnly)
	g.GET("/assessment-report/:id", h.getAssessmentReport, authed, notStudent)
	g.GET("/course-report/:batch_id/:course_id", h.getCourseReport, authed, notStudent)
	g.GET("/full-assessment-report/:id", h.triggerFullAssessmentReport, authed, notStudent)
}

type handler struct {
	db  *db.Client
	pub *queue.Publisher
}

func (h *handler) coll(name string) *mongo.Collection { return h.db.Coll(name) }
func iso(t time.Time) string                          { return t.Format(time.RFC3339) }

func pageLimit(c echo.Context, defaultLimit, maxLimit int) (page, limit int) {
	page = 1
	if v, err := strconv.Atoi(c.QueryParam("page")); err == nil && v >= 1 {
		page = v
	}
	limit = defaultLimit
	if v, err := strconv.Atoi(c.QueryParam("limit")); err == nil && v >= 1 && v <= maxLimit {
		limit = v
	}
	return page, limit
}

func totalPages(total, limit int64) int64 {
	if total <= 0 {
		return 1
	}
	return int64(math.Ceil(float64(total) / float64(limit)))
}

func (h *handler) getStatsCount(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	count := func(coll string, filter bson.M) int64 {
		if filter == nil {
			filter = bson.M{}
		}
		n, _ := h.coll(coll).CountDocuments(ctx, filter)
		return n
	}

	return response.OK(c, "Stats fetched", map[string]any{
		"admins":       count(models.CollUsers, bson.M{"role": models.RoleAdmin}),
		"institutions": count(models.CollInstitutions, nil),
		"vendors":      count(models.CollVendors, nil),
		"batches":      count(models.CollBatches, nil),
		"teachers":     count(models.CollTeachers, nil),
		"students":     count(models.CollStudents, nil),
		"courses":      count(models.CollCourses, nil),
		"assessments":  count(models.CollAssessments, nil),
	})
}

func (h *handler) getAssessmentReport(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Assessment not found", "Not found")
	}
	page, limit := pageLimit(c, 10, 1000)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	var assessment models.Assessment
	err := h.coll(models.CollAssessments).FindOne(ctx, bson.M{"_id": oid}).Decode(&assessment)
	if err == mongo.ErrNoDocuments {
		return response.Err(c, http.StatusNotFound, "Assessment not found", "Not found")
	}
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch assessment", err.Error())
	}

	total, _ := h.coll(models.CollAssessmentSubmissions).CountDocuments(ctx, bson.M{"assessment_id": oid})

	findOpts := options.Find().SetSkip(int64((page - 1) * limit)).SetLimit(int64(limit))
	cur, err := h.coll(models.CollAssessmentSubmissions).Find(ctx, bson.M{"assessment_id": oid}, findOpts)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch submissions", err.Error())
	}
	var submissions []models.AssessmentSubmission
	if err := cur.All(ctx, &submissions); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch submissions", err.Error())
	}

	studentIDs := make([]primitive.ObjectID, 0, len(submissions))
	for _, s := range submissions {
		studentIDs = append(studentIDs, s.StudentID)
	}
	studentsByID := h.loadStudentsMap(ctx, studentIDs)

	data := make([]map[string]any, 0, len(submissions))
	for _, sub := range submissions {
		student, ok := studentsByID[sub.StudentID.Hex()]
		name, email, roll := "Unknown", "", ""
		if ok {
			name, email, roll = student.Name, student.Email, student.RollNumber
		}
		var startedAt, submittedAt any
		if !sub.StartedAt.IsZero() {
			startedAt = iso(sub.StartedAt)
		}
		if sub.SubmittedAt != nil {
			submittedAt = iso(*sub.SubmittedAt)
		}
		data = append(data, map[string]any{
			"id": sub.ID.Hex(), "student_id": sub.StudentID.Hex(), "student_name": name,
			"student_email": email, "student_roll_number": roll, "is_submitted": sub.IsSubmitted,
			"total_marks": sub.TotalMarks, "tab_switch_count": sub.TabSwitchCount,
			"proctoring_status": sub.ProctoringStatus, "student_ip": sub.StudentIP,
			"started_at": startedAt, "submitted_at": submittedAt,
		})
	}

	return response.OK(c, "Assessment report fetched", map[string]any{
		"assessment": map[string]any{
			"id": assessment.ID.Hex(), "name": assessment.Name, "description": assessment.Description,
			"start_time": iso(assessment.StartTime), "end_time": iso(assessment.EndTime),
			"status": assessment.Status, "report_status": assessment.ReportStatus,
			"individual_section_time_limit": assessment.IndividualSectionTimeLimit,
		},
		"submissions": data, "total": total, "page": page, "total_pages": totalPages(total, int64(limit)),
	})
}

func (h *handler) loadStudentsMap(ctx context.Context, ids []primitive.ObjectID) map[string]models.Student {
	out := map[string]models.Student{}
	if len(ids) == 0 {
		return out
	}
	cur, err := h.coll(models.CollStudents).Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return out
	}
	var students []models.Student
	_ = cur.All(ctx, &students)
	for _, s := range students {
		out[s.ID.Hex()] = s
	}
	return out
}

func groupByStudent[T any](records []T, studentID func(T) primitive.ObjectID) map[string][]T {
	grouped := map[string][]T{}
	for _, r := range records {
		key := studentID(r).Hex()
		grouped[key] = append(grouped[key], r)
	}
	return grouped
}

func (h *handler) getCourseReport(c echo.Context) error {
	batchOID, ok := dbutil.ParseID(c.Param("batch_id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Batch not found", "Not found")
	}
	courseOID, ok := dbutil.ParseID(c.Param("course_id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Course not found", "Not found")
	}
	page, limit := pageLimit(c, 10, 100)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	var course models.Course
	if err := h.coll(models.CollCourses).FindOne(ctx, bson.M{"_id": courseOID}).Decode(&course); err != nil {
		return response.Err(c, http.StatusNotFound, "Course not found", "Not found")
	}
	var batch models.Batch
	if err := h.coll(models.CollBatches).FindOne(ctx, bson.M{"_id": batchOID}).Decode(&batch); err != nil {
		return response.Err(c, http.StatusNotFound, "Batch not found", "Not found")
	}

	total, _ := h.coll(models.CollStudents).CountDocuments(ctx, bson.M{"batch_id": batchOID})
	findOpts := options.Find().SetSkip(int64((page - 1) * limit)).SetLimit(int64(limit))
	scur, err := h.coll(models.CollStudents).Find(ctx, bson.M{"batch_id": batchOID}, findOpts)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch students", err.Error())
	}
	var students []models.Student
	if err := scur.All(ctx, &students); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch students", err.Error())
	}

	studentIDs := make([]primitive.ObjectID, 0, len(students))
	for _, s := range students {
		studentIDs = append(studentIDs, s.ID)
	}

	var sections []models.CourseSection
	if cur, err := h.coll(models.CollCourseSections).Find(ctx, bson.M{"course_id": course.ID}); err == nil {
		_ = cur.All(ctx, &sections)
	}
	sectionIDs := make([]primitive.ObjectID, 0, len(sections))
	for _, s := range sections {
		sectionIDs = append(sectionIDs, s.ID)
	}

	var contents []models.CourseLearningContent
	var assignments []models.CourseAssignment
	if len(sectionIDs) > 0 {
		if cur, err := h.coll(models.CollCourseLearningContents).Find(ctx, bson.M{"section_id": bson.M{"$in": sectionIDs}}); err == nil {
			_ = cur.All(ctx, &contents)
		}
		if cur, err := h.coll(models.CollCourseAssignments).Find(ctx, bson.M{"section_id": bson.M{"$in": sectionIDs}}); err == nil {
			_ = cur.All(ctx, &assignments)
		}
	}
	contentIDs := make([]primitive.ObjectID, 0, len(contents))
	for _, ct := range contents {
		contentIDs = append(contentIDs, ct.ID)
	}
	assignmentIDs := make([]primitive.ObjectID, 0, len(assignments))
	for _, a := range assignments {
		assignmentIDs = append(assignmentIDs, a.ID)
	}

	var progresses []models.CourseProgress
	var submissions []models.CourseAssignmentSubmission
	if len(studentIDs) > 0 && len(contentIDs) > 0 {
		if cur, err := h.coll(models.CollCourseProgresses).Find(ctx, bson.M{
			"student_id": bson.M{"$in": studentIDs}, "content_id": bson.M{"$in": contentIDs},
		}); err == nil {
			_ = cur.All(ctx, &progresses)
		}
	}
	if len(studentIDs) > 0 && len(assignmentIDs) > 0 {
		if cur, err := h.coll(models.CollCourseAssignmentSubmissions).Find(ctx, bson.M{
			"student_id": bson.M{"$in": studentIDs}, "assignment_id": bson.M{"$in": assignmentIDs},
		}); err == nil {
			_ = cur.All(ctx, &submissions)
		}
	}

	progressesByStudent := groupByStudent(progresses, func(p models.CourseProgress) primitive.ObjectID { return p.StudentID })
	submissionsByStudent := groupByStudent(submissions, func(s models.CourseAssignmentSubmission) primitive.ObjectID { return s.StudentID })

	data := make([]map[string]any, 0, len(students))
	for _, student := range students {
		key := student.ID.Hex()
		studentProgresses := progressesByStudent[key]
		studentSubmissions := submissionsByStudent[key]

		progressData := make([]map[string]any, 0, len(studentProgresses))
		for _, p := range studentProgresses {
			progressData = append(progressData, map[string]any{"id": p.ID.Hex(), "contentId": p.ContentID.Hex()})
		}
		submissionData := make([]map[string]any, 0, len(studentSubmissions))
		for _, s := range studentSubmissions {
			submissionData = append(submissionData, map[string]any{"id": s.ID.Hex(), "assignmentId": s.AssignmentID.Hex()})
		}

		data = append(data, map[string]any{
			"id": key, "name": student.Name, "rollNumber": student.RollNumber,
			"courseProgresses": progressData,
			// Keep the legacy typo for existing frontend consumers, alongside the
			// corrected key — matches the Python response's deliberate duplication.
			"courseAssignemntSubmissions": submissionData,
			"courseAssignmentSubmissions": submissionData,
		})
	}

	return response.OK(c, "Course report fetched", map[string]any{
		"students":       data,
		"batch":          map[string]any{"id": batch.ID.Hex(), "batchname": batch.BatchName},
		"course":         map[string]any{"id": course.ID.Hex(), "name": course.Name},
		"total_students": total, "page": page, "total_pages": totalPages(total, int64(limit)),
		"totalCourseTopics": len(contentIDs), "totalAssignments": len(assignmentIDs),
	})
}

func (h *handler) triggerFullAssessmentReport(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Assessment not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	res, err := h.coll(models.CollAssessments).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{
		"report_status": models.ReportProcessing,
	}})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to trigger report", err.Error())
	}
	if res.MatchedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Assessment not found", "Not found")
	}

	if h.pub != nil {
		_ = h.pub.Publish(ctx, jobs.AssessmentReportQueue, jobs.AssessmentReportJob{AssessmentID: oid.Hex()})
	}

	return response.OK(c, "Full report generation triggered", map[string]any{"report_status": models.ReportProcessing})
}
