// Package course ports apps/python-server/routers/course.py to Go.
package course

import (
	"context"
	"io"
	"math"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bitwiselearn/go-server/internal/db"
	"github.com/bitwiselearn/go-server/internal/dbutil"
	"github.com/bitwiselearn/go-server/internal/middleware"
	"github.com/bitwiselearn/go-server/internal/models"
	"github.com/bitwiselearn/go-server/internal/response"
	"github.com/bitwiselearn/go-server/internal/services/blob"
)

const reqTimeout = 10 * time.Second

// Deps holds the dependencies the course handlers need.
type Deps struct {
	DB    *db.Client
	Auth  *middleware.Auth
	Store *blob.Store
}

// Register mounts every course route on g (expected prefix "/api/v1/courses").
func Register(g *echo.Group, d Deps) {
	h := &handler{db: d.DB, store: d.Store}
	authed := d.Auth.Required
	adminOnly := d.Auth.AdminOnly()
	notStudent := d.Auth.NotStudent()
	studentOnly := d.Auth.RequireRoles(models.RoleStudent)
	enrollmentViewers := d.Auth.RequireRoles(models.RoleSuperadmin, models.RoleAdmin, models.RoleInstitution)
	adminVendor := d.Auth.RequireRoles(models.RoleSuperadmin, models.RoleAdmin, models.RoleVendor)

	// Course CRUD
	g.POST("/create-course", h.createCourse, authed, adminOnly)
	g.POST("/upload-thumbnail/:id", h.uploadThumbnail, authed, adminOnly)
	g.POST("/upload-completion-certificate/:id", h.uploadCertificate, authed, adminOnly)
	g.PUT("/change-publish-status/:id", h.changePublishStatus, authed, adminOnly)
	g.PUT("/update-course/:id", h.updateCourse, authed, adminOnly)
	g.GET("/get-all-courses-by-admin", h.getAllCoursesByAdmin, authed, adminVendor)
	g.GET("/get-course-by-id/:id", h.getCourseByID, authed)
	g.GET("/get-course-by-institution/:id", h.getCourseByInstitution, authed)
	g.GET("/get-all-sections-by-course/:id", h.getAllSectionsByCourse, authed)
	g.DELETE("/delete-course/:id", h.deleteCourse, authed, adminOnly)
	g.GET("/get-student-courses", h.getStudentCourses, authed, studentOnly)
	g.GET("/listed-courses", h.listedCourses)

	// Sections
	g.GET("/get-course-section/:id", h.getCourseSection, authed)
	g.POST("/add-course-section/:id", h.addCourseSection, authed)
	g.PUT("/update-course-section/:id", h.updateCourseSection, authed)
	g.DELETE("/delete-course-section/:id", h.deleteCourseSection, authed)

	// Content
	g.POST("/add-content-to-section", h.addContentToSection, authed, adminOnly)
	g.DELETE("/delete-content/:id", h.deleteContent, authed, adminOnly)
	g.PUT("/update-content-to-section/:id", h.updateContent, authed, adminOnly)
	g.POST("/upload-file-in-content/:id", h.uploadFileInContent, authed, adminOnly)
	g.DELETE("/remove-file-in-content/:id", h.removeFileInContent, authed, adminOnly)

	// Assignments
	g.POST("/add-assignment-to-section/", h.addAssignmentToSection, authed, adminOnly)
	g.PUT("/update-assignment-to-section/:id", h.updateAssignment, authed, adminOnly)
	g.DELETE("/remove-assignment-from-section/:id", h.removeAssignment, authed, adminOnly)
	g.GET("/get-assignment-by-id/:id", h.getAssignmentByID, authed)
	g.POST("/add-assignment-question/:id", h.addAssignmentQuestion, authed, adminOnly)
	g.PUT("/update-assignment-question/:id", h.updateAssignmentQuestion, authed, adminOnly)
	g.DELETE("/remove-assignment-question/:id", h.removeAssignmentQuestion, authed, adminOnly)
	g.GET("/get-all-section-assignments/:id", h.getAllSectionAssignments, authed, adminOnly)
	g.GET("/get-student-section-assignments/:id", h.getStudentSectionAssignments, authed, studentOnly)

	// Grades
	g.GET("/get-all-assignment-marks/", h.getAllAssignmentMarks, authed, studentOnly)
	g.GET("/get-all-assignment-marks-by-courseId/:id", h.getAssignmentMarksByCourse, authed, studentOnly)
	g.GET("/get-assignment-report/:id", h.getAssignmentReport, authed, studentOnly)
	g.POST("/submit-course-assignment/:id", h.submitCourseAssignment, authed, studentOnly)

	// Progress
	g.POST("/mark-content-as-done/:id", h.markContentDone, authed, studentOnly)
	g.POST("/unmark-content-as-done/:id", h.unmarkContentDone, authed, studentOnly)
	g.GET("/get-all-course-progress/", h.getAllCourseProgress, authed, studentOnly)
	g.GET("/get-individual-course-progress/:id", h.getIndividualCourseProgress, authed, studentOnly)

	// Enrollments
	g.GET("/get-course-enrollments/:id", h.getCourseEnrollments, authed, enrollmentViewers)
	g.GET("/get-course-enrollments-by-batch/:id", h.getCourseEnrollmentsByBatch, authed)
	g.POST("/add-course-enrollment/", h.addCourseEnrollment, authed, notStudent)
	g.DELETE("/remove-course-enrollment/:id", h.removeCourseEnrollment, authed, notStudent)
}

type handler struct {
	db    *db.Client
	store *blob.Store
}

func (h *handler) coll(name string) *mongo.Collection { return h.db.Coll(name) }

func iso(t time.Time) string { return t.Format(time.RFC3339) }

func (h *handler) getCourse(ctx context.Context, id primitive.ObjectID) (*models.Course, error) {
	var c models.Course
	err := h.coll(models.CollCourses).FindOne(ctx, bson.M{"_id": id}).Decode(&c)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &c, err
}

func (h *handler) getSection(ctx context.Context, id primitive.ObjectID) (*models.CourseSection, error) {
	var s models.CourseSection
	err := h.coll(models.CollCourseSections).FindOne(ctx, bson.M{"_id": id}).Decode(&s)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &s, err
}

func (h *handler) getContent(ctx context.Context, id primitive.ObjectID) (*models.CourseLearningContent, error) {
	var c models.CourseLearningContent
	err := h.coll(models.CollCourseLearningContents).FindOne(ctx, bson.M{"_id": id}).Decode(&c)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &c, err
}

func (h *handler) getAssignment(ctx context.Context, id primitive.ObjectID) (*models.CourseAssignment, error) {
	var a models.CourseAssignment
	err := h.coll(models.CollCourseAssignments).FindOne(ctx, bson.M{"_id": id}).Decode(&a)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &a, err
}

func (h *handler) getBatch(ctx context.Context, id primitive.ObjectID) (*models.Batch, error) {
	var b models.Batch
	err := h.coll(models.CollBatches).FindOne(ctx, bson.M{"_id": id}).Decode(&b)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &b, err
}

func (h *handler) getInstitution(ctx context.Context, id primitive.ObjectID) (*models.Institution, error) {
	var inst models.Institution
	err := h.coll(models.CollInstitutions).FindOne(ctx, bson.M{"_id": id}).Decode(&inst)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &inst, err
}

func (h *handler) getStudent(ctx context.Context, id primitive.ObjectID) (*models.Student, error) {
	var s models.Student
	err := h.coll(models.CollStudents).FindOne(ctx, bson.M{"_id": id}).Decode(&s)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &s, err
}

// ================= Course CRUD =================

func (h *handler) createCourse(c echo.Context) error {
	var body createCourseRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	user := middleware.UserFrom(c)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	existing := h.coll(models.CollCourses).FindOne(ctx, bson.M{"name": body.Name})
	if existing.Err() == nil {
		return response.Err(c, http.StatusBadRequest, "Course name already exists", "Duplicate name")
	}

	level := body.Level
	if level == "" {
		level = models.CourseLevelBasic
	}
	createdBy, _ := dbutil.ParseID(user.ID)
	now := time.Now().UTC()
	course := models.Course{
		ID: primitive.NewObjectID(), Name: body.Name, Description: body.Description, Level: level,
		Duration: body.Duration, InstructorName: body.InstructorName, IsPublished: models.CourseNotPublished,
		CreatedBy: createdBy, CreatedAt: now, UpdatedAt: now,
	}
	if _, err := h.coll(models.CollCourses).InsertOne(ctx, course); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to create course", err.Error())
	}
	return response.JSON(c, http.StatusCreated, "Course created", map[string]any{
		"id": course.ID.Hex(), "name": course.Name,
	}, nil)
}

func readUploadedFile(c echo.Context, field string) ([]byte, string, string, error) {
	fh, err := c.FormFile(field)
	if err != nil {
		return nil, "", "", err
	}
	f, err := fh.Open()
	if err != nil {
		return nil, "", "", err
	}
	defer f.Close()
	content, err := io.ReadAll(f)
	if err != nil {
		return nil, "", "", err
	}
	return content, fh.Filename, fh.Header.Get("Content-Type"), nil
}

func (h *handler) uploadThumbnail(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Course not found", "Not found")
	}
	content, filename, contentType, err := readUploadedFile(c, "file")
	if err != nil {
		return response.Err(c, http.StatusBadRequest, "Failed to upload thumbnail", err.Error())
	}
	if filename == "" {
		filename = "thumbnail"
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	course, err := h.getCourse(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch course", err.Error())
	}
	if course == nil {
		return response.Err(c, http.StatusNotFound, "Course not found", "Not found")
	}

	url, err := h.store.Upload(ctx, content, "course-thumbnails", filename, contentType)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to upload thumbnail", err.Error())
	}
	if course.Thumbnail != nil {
		_ = h.store.Delete(ctx, *course.Thumbnail)
	}
	_, err = h.coll(models.CollCourses).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{
		"thumbnail": url, "updated_at": time.Now().UTC(),
	}})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to upload thumbnail", err.Error())
	}
	return response.OK(c, "Thumbnail uploaded", map[string]any{"thumbnail": url})
}

func (h *handler) uploadCertificate(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Course not found", "Not found")
	}
	content, filename, contentType, err := readUploadedFile(c, "file")
	if err != nil {
		return response.Err(c, http.StatusBadRequest, "Failed to upload certificate", err.Error())
	}
	if filename == "" {
		filename = "certificate"
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	course, err := h.getCourse(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch course", err.Error())
	}
	if course == nil {
		return response.Err(c, http.StatusNotFound, "Course not found", "Not found")
	}

	url, err := h.store.Upload(ctx, content, "course-certificates", filename, contentType)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to upload certificate", err.Error())
	}
	if course.Certificate != nil {
		_ = h.store.Delete(ctx, *course.Certificate)
	}
	_, err = h.coll(models.CollCourses).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{
		"certificate": url, "updated_at": time.Now().UTC(),
	}})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to upload certificate", err.Error())
	}
	return response.OK(c, "Certificate uploaded", map[string]any{"certificate": url})
}

func (h *handler) changePublishStatus(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Course not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	course, err := h.getCourse(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch course", err.Error())
	}
	if course == nil {
		return response.Err(c, http.StatusNotFound, "Course not found", "Not found")
	}

	newStatus := models.CoursePublished
	if course.IsPublished == models.CoursePublished {
		newStatus = models.CourseNotPublished
		_, _ = h.coll(models.CollCourseEnrollments).DeleteMany(ctx, bson.M{"course_id": course.ID})
	}
	_, err = h.coll(models.CollCourses).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{
		"is_published": newStatus, "updated_at": time.Now().UTC(),
	}})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to update status", err.Error())
	}
	return response.OK(c, "Status changed", map[string]any{"is_published": newStatus})
}

func (h *handler) updateCourse(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Course not found", "Not found")
	}
	var body updateCourseRequest
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
	if body.Level != nil {
		set["level"] = *body.Level
	}
	if body.Duration != nil {
		set["duration"] = *body.Duration
	}
	if body.InstructorName != nil {
		set["instructor_name"] = *body.InstructorName
	}
	res, err := h.coll(models.CollCourses).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": set})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to update course", err.Error())
	}
	if res.MatchedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Course not found", "Not found")
	}
	name := ""
	if body.Name != nil {
		name = *body.Name
	}
	return response.OK(c, "Course updated", map[string]any{"id": oid.Hex(), "name": name})
}

func (h *handler) getAllCoursesByAdmin(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	cur, err := h.coll(models.CollCourses).Find(ctx, bson.M{})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch courses", err.Error())
	}
	var courses []models.Course
	if err := cur.All(ctx, &courses); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch courses", err.Error())
	}
	data := make([]map[string]any, 0, len(courses))
	for _, cs := range courses {
		data = append(data, map[string]any{
			"id": cs.ID.Hex(), "name": cs.Name, "description": cs.Description, "level": cs.Level,
			"duration": cs.Duration, "thumbnail": cs.Thumbnail, "instructor_name": cs.InstructorName,
			"is_published": cs.IsPublished, "created_at": iso(cs.CreatedAt),
		})
	}
	return response.OK(c, "Courses fetched", data)
}

func (h *handler) getCourseByID(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Course not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	course, err := h.getCourse(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch course", err.Error())
	}
	if course == nil {
		return response.Err(c, http.StatusNotFound, "Course not found", "Not found")
	}

	sectionsCur, err := h.coll(models.CollCourseSections).Find(ctx, bson.M{"course_id": course.ID})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch course", err.Error())
	}
	var sections []models.CourseSection
	_ = sectionsCur.All(ctx, &sections)

	sectionData := make([]map[string]any, 0, len(sections))
	for _, s := range sections {
		var contents []models.CourseLearningContent
		if cur, err := h.coll(models.CollCourseLearningContents).Find(ctx, bson.M{"section_id": s.ID}); err == nil {
			_ = cur.All(ctx, &contents)
		}
		var assignments []models.CourseAssignment
		if cur, err := h.coll(models.CollCourseAssignments).Find(ctx, bson.M{"section_id": s.ID}); err == nil {
			_ = cur.All(ctx, &assignments)
		}
		contentData := make([]map[string]any, 0, len(contents))
		for _, ct := range contents {
			contentData = append(contentData, map[string]any{
				"id": ct.ID.Hex(), "name": ct.Name, "description": ct.Description,
				"video_url": ct.VideoURL, "transcript": ct.Transcript, "file": ct.File,
			})
		}
		assignmentData := make([]map[string]any, 0, len(assignments))
		for _, a := range assignments {
			assignmentData = append(assignmentData, map[string]any{
				"id": a.ID.Hex(), "name": a.Name, "description": a.Description,
				"instruction": a.Instruction, "marks_per_question": a.MarksPerQuestion,
			})
		}
		sectionData = append(sectionData, map[string]any{
			"id": s.ID.Hex(), "name": s.Name, "contents": contentData, "assignments": assignmentData,
		})
	}

	return response.OK(c, "Course fetched", map[string]any{
		"id": course.ID.Hex(), "name": course.Name, "description": course.Description,
		"level": course.Level, "duration": course.Duration, "thumbnail": course.Thumbnail,
		"instructor_name": course.InstructorName, "certificate": course.Certificate,
		"is_published": course.IsPublished, "sections": sectionData, "created_at": iso(course.CreatedAt),
	})
}

func (h *handler) getCourseByInstitution(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.OK(c, "No courses", []map[string]any{})
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	cur, err := h.coll(models.CollBatches).Find(ctx, bson.M{"institution_id": oid})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch courses", err.Error())
	}
	var batches []models.Batch
	_ = cur.All(ctx, &batches)
	if len(batches) == 0 {
		return response.OK(c, "No courses", []map[string]any{})
	}
	batchIDs := make([]primitive.ObjectID, 0, len(batches))
	for _, b := range batches {
		batchIDs = append(batchIDs, b.ID)
	}

	ecur, err := h.coll(models.CollCourseEnrollments).Find(ctx, bson.M{"batch_id": bson.M{"$in": batchIDs}})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch courses", err.Error())
	}
	var enrollments []models.CourseEnrollment
	_ = ecur.All(ctx, &enrollments)

	seen := map[primitive.ObjectID]struct{}{}
	data := []map[string]any{}
	for _, e := range enrollments {
		if _, ok := seen[e.CourseID]; ok {
			continue
		}
		seen[e.CourseID] = struct{}{}
		cs, err := h.getCourse(ctx, e.CourseID)
		if err == nil && cs != nil {
			data = append(data, map[string]any{
				"id": cs.ID.Hex(), "name": cs.Name, "level": cs.Level,
				"thumbnail": cs.Thumbnail, "instructor_name": cs.InstructorName,
			})
		}
	}
	return response.OK(c, "Courses fetched", data)
}

func (h *handler) getAllSectionsByCourse(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.OK(c, "Sections fetched", []map[string]any{})
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	cur, err := h.coll(models.CollCourseSections).Find(ctx, bson.M{"course_id": oid})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch sections", err.Error())
	}
	var sections []models.CourseSection
	_ = cur.All(ctx, &sections)

	data := make([]map[string]any, 0, len(sections))
	for _, s := range sections {
		var contents []models.CourseLearningContent
		if ccur, err := h.coll(models.CollCourseLearningContents).Find(ctx, bson.M{"section_id": s.ID}); err == nil {
			_ = ccur.All(ctx, &contents)
		}
		contentData := make([]map[string]any, 0, len(contents))
		for _, ct := range contents {
			contentData = append(contentData, map[string]any{
				"id": ct.ID.Hex(), "name": ct.Name, "description": ct.Description,
				"video_url": ct.VideoURL, "transcript": ct.Transcript, "file": ct.File,
			})
		}
		data = append(data, map[string]any{
			"id": s.ID.Hex(), "name": s.Name, "course_id": s.CourseID.Hex(),
			"course_learning_contents": contentData,
		})
	}
	return response.OK(c, "Sections fetched", data)
}

func (h *handler) deleteCourse(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Course not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	course, err := h.getCourse(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch course", err.Error())
	}
	if course == nil {
		return response.Err(c, http.StatusNotFound, "Course not found", "Not found")
	}

	sectionsCur, err := h.coll(models.CollCourseSections).Find(ctx, bson.M{"course_id": course.ID})
	if err == nil {
		var sections []models.CourseSection
		_ = sectionsCur.All(ctx, &sections)
		for _, s := range sections {
			var contents []models.CourseLearningContent
			if cur, err := h.coll(models.CollCourseLearningContents).Find(ctx, bson.M{"section_id": s.ID}); err == nil {
				_ = cur.All(ctx, &contents)
			}
			for _, ct := range contents {
				if ct.File != "" {
					_ = h.store.Delete(ctx, ct.File)
				}
				_, _ = h.coll(models.CollCourseProgresses).DeleteMany(ctx, bson.M{"content_id": ct.ID})
				_, _ = h.coll(models.CollCourseLearningContents).DeleteOne(ctx, bson.M{"_id": ct.ID})
			}
			_, _ = h.coll(models.CollCourseAssignments).DeleteMany(ctx, bson.M{"section_id": s.ID})
			_, _ = h.coll(models.CollCourseSections).DeleteOne(ctx, bson.M{"_id": s.ID})
		}
	}
	_, _ = h.coll(models.CollCourseEnrollments).DeleteMany(ctx, bson.M{"course_id": course.ID})
	if course.Thumbnail != nil {
		_ = h.store.Delete(ctx, *course.Thumbnail)
	}
	if course.Certificate != nil {
		_ = h.store.Delete(ctx, *course.Certificate)
	}
	if _, err := h.coll(models.CollCourses).DeleteOne(ctx, bson.M{"_id": oid}); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to delete course", err.Error())
	}
	return response.OK(c, "Course deleted", nil)
}

func (h *handler) getStudentCourses(c echo.Context) error {
	user := middleware.UserFrom(c)
	uid, _ := dbutil.ParseID(user.ID)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	student, err := h.getStudent(ctx, uid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch student", err.Error())
	}
	if student == nil {
		return response.Err(c, http.StatusNotFound, "Student not found", "Not found")
	}

	cur, err := h.coll(models.CollCourseEnrollments).Find(ctx, bson.M{"batch_id": student.BatchID})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch courses", err.Error())
	}
	var enrollments []models.CourseEnrollment
	_ = cur.All(ctx, &enrollments)

	data := []map[string]any{}
	for _, e := range enrollments {
		cs, err := h.getCourse(ctx, e.CourseID)
		if err == nil && cs != nil && cs.IsPublished == models.CoursePublished {
			data = append(data, map[string]any{
				"id": cs.ID.Hex(), "name": cs.Name, "description": cs.Description,
				"level": cs.Level, "thumbnail": cs.Thumbnail, "instructor_name": cs.InstructorName,
			})
		}
	}
	return response.OK(c, "Student courses", data)
}

func (h *handler) listedCourses(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	cur, err := h.coll(models.CollCourses).Find(ctx, bson.M{"is_published": models.CoursePublished})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch courses", err.Error())
	}
	var courses []models.Course
	if err := cur.All(ctx, &courses); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch courses", err.Error())
	}
	data := make([]map[string]any, 0, len(courses))
	for _, cs := range courses {
		data = append(data, map[string]any{
			"id": cs.ID.Hex(), "name": cs.Name, "description": cs.Description,
			"level": cs.Level, "thumbnail": cs.Thumbnail, "instructor_name": cs.InstructorName,
			"duration": cs.Duration,
		})
	}
	return response.OK(c, "Listed courses", data)
}

// ================= Sections =================

func (h *handler) getCourseSection(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Section not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	section, err := h.getSection(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch section", err.Error())
	}
	if section == nil {
		return response.Err(c, http.StatusNotFound, "Section not found", "Not found")
	}
	cur, err := h.coll(models.CollCourseLearningContents).Find(ctx, bson.M{"section_id": section.ID})
	var contents []models.CourseLearningContent
	if err == nil {
		_ = cur.All(ctx, &contents)
	}
	contentData := make([]map[string]any, 0, len(contents))
	for _, ct := range contents {
		contentData = append(contentData, map[string]any{
			"id": ct.ID.Hex(), "name": ct.Name, "description": ct.Description,
			"video_url": ct.VideoURL, "transcript": ct.Transcript, "file": ct.File,
		})
	}
	return response.OK(c, "Section fetched", map[string]any{
		"id": section.ID.Hex(), "name": section.Name, "course_id": section.CourseID.Hex(),
		"contents": contentData,
	})
}

func (h *handler) addCourseSection(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Course not found", "Not found")
	}
	var body createSectionRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	user := middleware.UserFrom(c)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	course, err := h.getCourse(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch course", err.Error())
	}
	if course == nil {
		return response.Err(c, http.StatusNotFound, "Course not found", "Not found")
	}

	creatorID, _ := dbutil.ParseID(user.ID)
	now := time.Now().UTC()
	section := models.CourseSection{
		ID: primitive.NewObjectID(), Name: body.Name, CreatorID: creatorID, CourseID: course.ID,
		CreatedAt: now, UpdatedAt: now,
	}
	if _, err := h.coll(models.CollCourseSections).InsertOne(ctx, section); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to create section", err.Error())
	}
	return response.JSON(c, http.StatusCreated, "Section created", map[string]any{
		"id": section.ID.Hex(), "name": section.Name,
	}, nil)
}

func (h *handler) updateCourseSection(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Section not found", "Not found")
	}
	var body updateSectionRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	set := bson.M{"updated_at": time.Now().UTC()}
	if body.Name != nil && *body.Name != "" {
		set["name"] = *body.Name
	}
	res, err := h.coll(models.CollCourseSections).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": set})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to update section", err.Error())
	}
	if res.MatchedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Section not found", "Not found")
	}
	section, _ := h.getSection(ctx, oid)
	name := ""
	if section != nil {
		name = section.Name
	}
	return response.OK(c, "Section updated", map[string]any{"id": oid.Hex(), "name": name})
}

func (h *handler) deleteCourseSection(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Section not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	section, err := h.getSection(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch section", err.Error())
	}
	if section == nil {
		return response.Err(c, http.StatusNotFound, "Section not found", "Not found")
	}

	cur, err := h.coll(models.CollCourseLearningContents).Find(ctx, bson.M{"section_id": section.ID})
	if err == nil {
		var contents []models.CourseLearningContent
		_ = cur.All(ctx, &contents)
		for _, ct := range contents {
			if ct.File != "" {
				_ = h.store.Delete(ctx, ct.File)
			}
			_, _ = h.coll(models.CollCourseProgresses).DeleteMany(ctx, bson.M{"content_id": ct.ID})
			_, _ = h.coll(models.CollCourseLearningContents).DeleteOne(ctx, bson.M{"_id": ct.ID})
		}
	}
	_, _ = h.coll(models.CollCourseAssignments).DeleteMany(ctx, bson.M{"section_id": section.ID})
	if _, err := h.coll(models.CollCourseSections).DeleteOne(ctx, bson.M{"_id": oid}); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to delete section", err.Error())
	}
	return response.OK(c, "Section deleted", nil)
}

// ================= Content =================

func (h *handler) addContentToSection(c echo.Context) error {
	var body addContentRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	sectionID, ok := dbutil.ParseID(body.SectionID)
	if !ok {
		return response.Err(c, http.StatusNotFound, "Section not found", "Not found")
	}
	user := middleware.UserFrom(c)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	section, err := h.getSection(ctx, sectionID)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch section", err.Error())
	}
	if section == nil {
		return response.Err(c, http.StatusNotFound, "Section not found", "Not found")
	}

	creatorID, _ := dbutil.ParseID(user.ID)
	now := time.Now().UTC()
	content := models.CourseLearningContent{
		ID: primitive.NewObjectID(), Name: body.Name, Description: body.Description,
		CreatorID: creatorID, SectionID: section.ID, VideoURL: body.VideoURL, Transcript: body.Transcript,
		CreatedAt: now, UpdatedAt: now,
	}
	if _, err := h.coll(models.CollCourseLearningContents).InsertOne(ctx, content); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to add content", err.Error())
	}
	return response.JSON(c, http.StatusCreated, "Content added", map[string]any{
		"id": content.ID.Hex(), "name": content.Name,
	}, nil)
}

func (h *handler) deleteContent(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Content not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	content, err := h.getContent(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch content", err.Error())
	}
	if content == nil {
		return response.Err(c, http.StatusNotFound, "Content not found", "Not found")
	}
	if content.File != "" {
		_ = h.store.Delete(ctx, content.File)
	}
	_, _ = h.coll(models.CollCourseProgresses).DeleteMany(ctx, bson.M{"content_id": content.ID})
	if _, err := h.coll(models.CollCourseLearningContents).DeleteOne(ctx, bson.M{"_id": oid}); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to delete content", err.Error())
	}
	return response.OK(c, "Content deleted", nil)
}

func (h *handler) updateContent(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Content not found", "Not found")
	}
	var body updateContentRequest
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
	if body.VideoURL != nil {
		set["video_url"] = *body.VideoURL
	}
	if body.Transcript != nil {
		set["transcript"] = *body.Transcript
	}
	res, err := h.coll(models.CollCourseLearningContents).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": set})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to update content", err.Error())
	}
	if res.MatchedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Content not found", "Not found")
	}
	return response.OK(c, "Content updated", map[string]any{"id": oid.Hex()})
}

func (h *handler) uploadFileInContent(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Content not found", "Not found")
	}
	content, filename, contentType, err := readUploadedFile(c, "file")
	if err != nil {
		return response.Err(c, http.StatusBadRequest, "Failed to upload file", err.Error())
	}
	if filename == "" {
		filename = "file"
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	ct, err := h.getContent(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch content", err.Error())
	}
	if ct == nil {
		return response.Err(c, http.StatusNotFound, "Content not found", "Not found")
	}

	url, err := h.store.Upload(ctx, content, "course-content", filename, contentType)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to upload file", err.Error())
	}
	if ct.File != "" {
		_ = h.store.Delete(ctx, ct.File)
	}
	_, err = h.coll(models.CollCourseLearningContents).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{
		"file": url, "updated_at": time.Now().UTC(),
	}})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to upload file", err.Error())
	}
	return response.OK(c, "File uploaded", map[string]any{"file": url})
}

func (h *handler) removeFileInContent(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Content not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	ct, err := h.getContent(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch content", err.Error())
	}
	if ct == nil {
		return response.Err(c, http.StatusNotFound, "Content not found", "Not found")
	}
	if ct.File != "" {
		_ = h.store.Delete(ctx, ct.File)
		_, err = h.coll(models.CollCourseLearningContents).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{
			"file": "", "updated_at": time.Now().UTC(),
		}})
		if err != nil {
			return response.Err(c, http.StatusInternalServerError, "Failed to remove file", err.Error())
		}
	}
	return response.OK(c, "File removed", nil)
}

// ================= Assignments =================

func (h *handler) addAssignmentToSection(c echo.Context) error {
	var body createAssignmentRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	sectionID, ok := dbutil.ParseID(body.SectionID)
	if !ok {
		return response.Err(c, http.StatusNotFound, "Section not found", "Not found")
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
	assignment := models.CourseAssignment{
		ID: primitive.NewObjectID(), Name: body.Name, Description: body.Description,
		Instruction: body.Instruction, MarksPerQuestion: body.MarksPerQuestion, SectionID: section.ID,
		CreatedAt: now, UpdatedAt: now,
	}
	if _, err := h.coll(models.CollCourseAssignments).InsertOne(ctx, assignment); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to create assignment", err.Error())
	}
	return response.JSON(c, http.StatusCreated, "Assignment created", map[string]any{
		"id": assignment.ID.Hex(), "name": assignment.Name,
	}, nil)
}

func (h *handler) updateAssignment(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Assignment not found", "Not found")
	}
	var body updateAssignmentRequest
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
	if body.MarksPerQuestion != nil {
		set["marks_per_question"] = *body.MarksPerQuestion
	}
	res, err := h.coll(models.CollCourseAssignments).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": set})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to update assignment", err.Error())
	}
	if res.MatchedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Assignment not found", "Not found")
	}
	return response.OK(c, "Assignment updated", map[string]any{"id": oid.Hex()})
}

func (h *handler) removeAssignment(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Assignment not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	assignment, err := h.getAssignment(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch assignment", err.Error())
	}
	if assignment == nil {
		return response.Err(c, http.StatusNotFound, "Assignment not found", "Not found")
	}
	_, _ = h.coll(models.CollCourseAssignmentQuestions).DeleteMany(ctx, bson.M{"assignment_id": assignment.ID})
	_, _ = h.coll(models.CollCourseAssignmentSubmissions).DeleteMany(ctx, bson.M{"assignment_id": assignment.ID})
	if _, err := h.coll(models.CollCourseAssignments).DeleteOne(ctx, bson.M{"_id": oid}); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to delete assignment", err.Error())
	}
	return response.OK(c, "Assignment deleted", nil)
}

func (h *handler) getAssignmentByID(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Assignment not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	assignment, err := h.getAssignment(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch assignment", err.Error())
	}
	if assignment == nil {
		return response.Err(c, http.StatusNotFound, "Assignment not found", "Not found")
	}
	questions := h.questionsForAssignment(ctx, assignment.ID)
	qData := make([]map[string]any, 0, len(questions))
	for _, q := range questions {
		qData = append(qData, map[string]any{
			"id": q.ID.Hex(), "question": q.Question, "options": q.Options,
			"correct_answer": q.CorrectAnswer, "type": q.Type,
		})
	}
	return response.OK(c, "Assignment fetched", map[string]any{
		"id": assignment.ID.Hex(), "name": assignment.Name, "description": assignment.Description,
		"instruction": assignment.Instruction, "marks_per_question": assignment.MarksPerQuestion,
		"section_id": assignment.SectionID.Hex(), "questions": qData,
	})
}

func (h *handler) questionsForAssignment(ctx context.Context, assignmentID primitive.ObjectID) []models.CourseAssignmentQuestion {
	cur, err := h.coll(models.CollCourseAssignmentQuestions).Find(ctx, bson.M{"assignment_id": assignmentID})
	if err != nil {
		return nil
	}
	var qs []models.CourseAssignmentQuestion
	_ = cur.All(ctx, &qs)
	return qs
}

func (h *handler) submissionsFor(ctx context.Context, assignmentID, studentID primitive.ObjectID) []models.CourseAssignmentSubmission {
	cur, err := h.coll(models.CollCourseAssignmentSubmissions).Find(ctx, bson.M{
		"assignment_id": assignmentID, "student_id": studentID,
	})
	if err != nil {
		return nil
	}
	var subs []models.CourseAssignmentSubmission
	_ = cur.All(ctx, &subs)
	return subs
}

func (h *handler) addAssignmentQuestion(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Assignment not found", "Not found")
	}
	body := addAssignmentQuestionRequest{Type: models.AssignmentSCQ}
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	assignment, err := h.getAssignment(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch assignment", err.Error())
	}
	if assignment == nil {
		return response.Err(c, http.StatusNotFound, "Assignment not found", "Not found")
	}

	now := time.Now().UTC()
	q := models.CourseAssignmentQuestion{
		ID: primitive.NewObjectID(), Question: body.Question, Options: body.Options,
		CorrectAnswer: body.CorrectAnswer, AssignmentID: assignment.ID, Type: body.Type,
		CreatedAt: now, UpdatedAt: now,
	}
	if _, err := h.coll(models.CollCourseAssignmentQuestions).InsertOne(ctx, q); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to add question", err.Error())
	}
	return response.JSON(c, http.StatusCreated, "Question added", map[string]any{"id": q.ID.Hex()}, nil)
}

func (h *handler) updateAssignmentQuestion(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Question not found", "Not found")
	}
	var body updateAssignmentQuestionRequest
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
	if body.CorrectAnswer != nil {
		set["correct_answer"] = *body.CorrectAnswer
	}
	if body.Type != nil {
		set["type"] = *body.Type
	}
	res, err := h.coll(models.CollCourseAssignmentQuestions).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": set})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to update question", err.Error())
	}
	if res.MatchedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Question not found", "Not found")
	}
	return response.OK(c, "Question updated", map[string]any{"id": oid.Hex()})
}

func (h *handler) removeAssignmentQuestion(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Question not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	_, _ = h.coll(models.CollCourseAssignmentSubmissions).DeleteMany(ctx, bson.M{"question_id": oid})
	res, err := h.coll(models.CollCourseAssignmentQuestions).DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to delete question", err.Error())
	}
	if res.DeletedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Question not found", "Not found")
	}
	return response.OK(c, "Question deleted", nil)
}

func (h *handler) getAllSectionAssignments(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.OK(c, "Assignments fetched", []map[string]any{})
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	cur, err := h.coll(models.CollCourseAssignments).Find(ctx, bson.M{"section_id": oid})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch assignments", err.Error())
	}
	var assignments []models.CourseAssignment
	_ = cur.All(ctx, &assignments)
	data := make([]map[string]any, 0, len(assignments))
	for _, a := range assignments {
		data = append(data, map[string]any{
			"id": a.ID.Hex(), "name": a.Name, "description": a.Description,
			"marks_per_question": a.MarksPerQuestion,
		})
	}
	return response.OK(c, "Assignments fetched", data)
}

func (h *handler) getStudentSectionAssignments(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.OK(c, "Assignments fetched", []map[string]any{})
	}
	user := middleware.UserFrom(c)
	studentID, _ := dbutil.ParseID(user.ID)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	cur, err := h.coll(models.CollCourseAssignments).Find(ctx, bson.M{"section_id": oid})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch assignments", err.Error())
	}
	var assignments []models.CourseAssignment
	_ = cur.All(ctx, &assignments)

	data := make([]map[string]any, 0, len(assignments))
	for _, a := range assignments {
		subs := h.submissionsFor(ctx, a.ID, studentID)
		data = append(data, map[string]any{
			"id": a.ID.Hex(), "name": a.Name, "description": a.Description,
			"marks_per_question": a.MarksPerQuestion, "attempted": len(subs) > 0,
			"submission_count": len(subs),
		})
	}
	return response.OK(c, "Assignments fetched", data)
}

// ================= Grades =================

func round2(v float64) float64 { return math.Round(v*100) / 100 }

func (h *handler) getAllAssignmentMarks(c echo.Context) error {
	user := middleware.UserFrom(c)
	studentID, _ := dbutil.ParseID(user.ID)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	student, err := h.getStudent(ctx, studentID)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch student", err.Error())
	}
	if student == nil {
		return response.Err(c, http.StatusNotFound, "Student not found", "Not found")
	}

	cur, err := h.coll(models.CollCourseEnrollments).Find(ctx, bson.M{"batch_id": student.BatchID})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch enrollments", err.Error())
	}
	var enrollments []models.CourseEnrollment
	_ = cur.All(ctx, &enrollments)

	results := []map[string]any{}
	for _, e := range enrollments {
		course, err := h.getCourse(ctx, e.CourseID)
		if err != nil || course == nil {
			continue
		}
		sectionsCur, _ := h.coll(models.CollCourseSections).Find(ctx, bson.M{"course_id": course.ID})
		var sections []models.CourseSection
		if sectionsCur != nil {
			_ = sectionsCur.All(ctx, &sections)
		}
		totalMarks, obtainedMarks := 0, 0.0
		for _, sec := range sections {
			assignCur, _ := h.coll(models.CollCourseAssignments).Find(ctx, bson.M{"section_id": sec.ID})
			var assignments []models.CourseAssignment
			if assignCur != nil {
				_ = assignCur.All(ctx, &assignments)
			}
			for _, a := range assignments {
				questions := h.questionsForAssignment(ctx, a.ID)
				totalMarks += len(questions) * a.MarksPerQuestion
				subs := h.submissionsFor(ctx, a.ID, studentID)
				for _, s := range subs {
					if s.MarksObtained != nil {
						obtainedMarks += *s.MarksObtained
					}
				}
			}
		}
		results = append(results, map[string]any{
			"course_id": course.ID.Hex(), "course_name": course.Name,
			"total_marks": totalMarks, "obtained_marks": obtainedMarks,
		})
	}
	return response.OK(c, "Assignment marks", results)
}

func (h *handler) getAssignmentMarksByCourse(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Course not found", "Not found")
	}
	user := middleware.UserFrom(c)
	studentID, _ := dbutil.ParseID(user.ID)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	course, err := h.getCourse(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch course", err.Error())
	}
	if course == nil {
		return response.Err(c, http.StatusNotFound, "Course not found", "Not found")
	}

	sectionsCur, err := h.coll(models.CollCourseSections).Find(ctx, bson.M{"course_id": course.ID})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch sections", err.Error())
	}
	var sections []models.CourseSection
	_ = sectionsCur.All(ctx, &sections)

	results := []map[string]any{}
	for _, sec := range sections {
		assignCur, _ := h.coll(models.CollCourseAssignments).Find(ctx, bson.M{"section_id": sec.ID})
		var assignments []models.CourseAssignment
		if assignCur != nil {
			_ = assignCur.All(ctx, &assignments)
		}
		for _, a := range assignments {
			questions := h.questionsForAssignment(ctx, a.ID)
			total := len(questions) * a.MarksPerQuestion
			subs := h.submissionsFor(ctx, a.ID, studentID)
			obtained := 0.0
			for _, s := range subs {
				if s.MarksObtained != nil {
					obtained += *s.MarksObtained
				}
			}
			results = append(results, map[string]any{
				"assignment_id": a.ID.Hex(), "assignment_name": a.Name,
				"section_name": sec.Name, "total_marks": total, "obtained_marks": obtained,
			})
		}
	}
	return response.OK(c, "Assignment marks", results)
}

func buildAssignmentReport(assignment *models.CourseAssignment, questions []models.CourseAssignmentQuestion, submissions []models.CourseAssignmentSubmission) map[string]any {
	obtained := 0.0
	correct := 0
	var latest *time.Time
	for _, s := range submissions {
		if s.MarksObtained != nil {
			obtained += *s.MarksObtained
		}
		if s.IsCorrect != nil && *s.IsCorrect {
			correct++
		}
		if latest == nil || s.SubmittedAt.After(*latest) {
			t := s.SubmittedAt
			latest = &t
		}
	}
	total := float64(len(questions)) * float64(assignment.MarksPerQuestion)
	percentage := 0.0
	if total > 0 {
		percentage = round2(obtained / total * 100)
	}
	var attemptedAt any
	if latest != nil {
		attemptedAt = iso(*latest)
	}
	return map[string]any{
		"assignment_id":      assignment.ID.Hex(),
		"assignment_name":    assignment.Name,
		"total_questions":    len(questions),
		"answered_questions": len(submissions),
		"correct_answers":    correct,
		"obtained_marks":     obtained,
		"total_marks":        total,
		"percentage":         percentage,
		"attempted_at":       attemptedAt,
	}
}

func (h *handler) getAssignmentReport(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Assignment not found", "Not found")
	}
	user := middleware.UserFrom(c)
	studentID, _ := dbutil.ParseID(user.ID)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	assignment, err := h.getAssignment(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch assignment", err.Error())
	}
	if assignment == nil {
		return response.Err(c, http.StatusNotFound, "Assignment not found", "Not found")
	}
	questions := h.questionsForAssignment(ctx, assignment.ID)
	submissions := h.submissionsFor(ctx, assignment.ID, studentID)
	if len(submissions) == 0 {
		return response.Err(c, http.StatusNotFound, "Report not found", "Not found")
	}
	return response.OK(c, "Assignment report fetched", buildAssignmentReport(assignment, questions, submissions))
}

func (h *handler) submitCourseAssignment(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Assignment not found", "Not found")
	}
	var body submitAssignmentRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	user := middleware.UserFrom(c)
	studentID, _ := dbutil.ParseID(user.ID)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	assignment, err := h.getAssignment(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch assignment", err.Error())
	}
	if assignment == nil {
		return response.Err(c, http.StatusNotFound, "Assignment not found", "Not found")
	}

	existing := h.submissionsFor(ctx, assignment.ID, studentID)
	if len(existing) > 0 {
		questions := h.questionsForAssignment(ctx, assignment.ID)
		return response.JSON(c, http.StatusConflict, "Assignment already attempted", map[string]any{
			"report": buildAssignmentReport(assignment, questions, existing),
		}, "Only one attempt is allowed")
	}

	questions := h.questionsForAssignment(ctx, assignment.ID)
	questionMap := make(map[string]models.CourseAssignmentQuestion, len(questions))
	for _, q := range questions {
		questionMap[q.ID.Hex()] = q
	}

	type result struct {
		QuestionID string `json:"question_id"`
		IsCorrect  bool   `json:"is_correct"`
		Marks      int    `json:"marks"`
	}
	results := make([]result, 0, len(body.Answers))
	for _, ans := range body.Answers {
		if ans.QuestionID == "" {
			continue
		}
		question, ok := questionMap[ans.QuestionID]
		if !ok {
			continue
		}
		isCorrect := sortedEqual(ans.Answer, question.CorrectAnswer)
		marks := 0
		if isCorrect {
			marks = assignment.MarksPerQuestion
		}
		now := time.Now().UTC()
		sub := models.CourseAssignmentSubmission{
			ID: primitive.NewObjectID(), QuestionID: question.ID, StudentID: studentID,
			Answer: ans.Answer, AssignmentID: assignment.ID,
			MarksObtained: floatPtr(float64(marks)), IsCorrect: boolPtr(isCorrect),
			SubmittedAt: now, CreatedAt: now, UpdatedAt: now,
		}
		if _, err := h.coll(models.CollCourseAssignmentSubmissions).InsertOne(ctx, sub); err != nil {
			return response.Err(c, http.StatusInternalServerError, "Failed to submit assignment", err.Error())
		}
		results = append(results, result{QuestionID: question.ID.Hex(), IsCorrect: isCorrect, Marks: marks})
	}

	resultsData := make([]map[string]any, 0, len(results))
	obtainedMarks := 0
	correctCount := 0
	for _, r := range results {
		resultsData = append(resultsData, map[string]any{
			"question_id": r.QuestionID, "is_correct": r.IsCorrect, "marks": r.Marks,
		})
		obtainedMarks += r.Marks
		if r.IsCorrect {
			correctCount++
		}
	}
	totalQuestions := len(questionMap)
	totalMarks := totalQuestions * assignment.MarksPerQuestion
	percentage := 0.0
	if totalMarks > 0 {
		percentage = round2(float64(obtainedMarks) / float64(totalMarks) * 100)
	}

	return response.OK(c, "Assignment submitted", map[string]any{
		"results": resultsData,
		"report": map[string]any{
			"assignment_id":      assignment.ID.Hex(),
			"assignment_name":    assignment.Name,
			"total_questions":    totalQuestions,
			"answered_questions": len(results),
			"correct_answers":    correctCount,
			"obtained_marks":     obtainedMarks,
			"total_marks":        totalMarks,
			"percentage":         percentage,
			"attempted_at":       iso(time.Now().UTC()),
		},
	})
}

func floatPtr(v float64) *float64 { return &v }
func boolPtr(v bool) *bool        { return &v }

func sortedEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	ac := append([]string(nil), a...)
	bc := append([]string(nil), b...)
	sortStrings(ac)
	sortStrings(bc)
	for i := range ac {
		if ac[i] != bc[i] {
			return false
		}
	}
	return true
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j-1] > s[j]; j-- {
			s[j-1], s[j] = s[j], s[j-1]
		}
	}
}

// ================= Progress =================

func (h *handler) markContentDone(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Content not found", "Not found")
	}
	user := middleware.UserFrom(c)
	studentID, _ := dbutil.ParseID(user.ID)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	content, err := h.getContent(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch content", err.Error())
	}
	if content == nil {
		return response.Err(c, http.StatusNotFound, "Content not found", "Not found")
	}

	existing := h.coll(models.CollCourseProgresses).FindOne(ctx, bson.M{"student_id": studentID, "content_id": content.ID})
	if existing.Err() == nil {
		return response.OK(c, "Already marked as done", nil)
	}
	now := time.Now().UTC()
	progress := models.CourseProgress{ID: primitive.NewObjectID(), StudentID: studentID, ContentID: content.ID, CreatedAt: now, UpdatedAt: now}
	if _, err := h.coll(models.CollCourseProgresses).InsertOne(ctx, progress); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to mark content done", err.Error())
	}
	return response.JSON(c, http.StatusCreated, "Content marked as done", nil, nil)
}

func (h *handler) unmarkContentDone(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Content not found", "Not found")
	}
	user := middleware.UserFrom(c)
	studentID, _ := dbutil.ParseID(user.ID)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	content, err := h.getContent(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch content", err.Error())
	}
	if content == nil {
		return response.Err(c, http.StatusNotFound, "Content not found", "Not found")
	}
	_, _ = h.coll(models.CollCourseProgresses).DeleteMany(ctx, bson.M{"student_id": studentID, "content_id": content.ID})
	return response.OK(c, "Content unmarked", nil)
}

func (h *handler) getAllCourseProgress(c echo.Context) error {
	user := middleware.UserFrom(c)
	studentID, _ := dbutil.ParseID(user.ID)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	student, err := h.getStudent(ctx, studentID)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch student", err.Error())
	}
	if student == nil {
		return response.Err(c, http.StatusNotFound, "Student not found", "Not found")
	}

	cur, err := h.coll(models.CollCourseEnrollments).Find(ctx, bson.M{"batch_id": student.BatchID})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch enrollments", err.Error())
	}
	var enrollments []models.CourseEnrollment
	_ = cur.All(ctx, &enrollments)

	results := []map[string]any{}
	for _, e := range enrollments {
		course, err := h.getCourse(ctx, e.CourseID)
		if err != nil || course == nil {
			continue
		}
		sectionsCur, _ := h.coll(models.CollCourseSections).Find(ctx, bson.M{"course_id": course.ID})
		var sections []models.CourseSection
		if sectionsCur != nil {
			_ = sectionsCur.All(ctx, &sections)
		}
		totalContent, doneContent := 0, 0
		for _, sec := range sections {
			contentsCur, _ := h.coll(models.CollCourseLearningContents).Find(ctx, bson.M{"section_id": sec.ID})
			var contents []models.CourseLearningContent
			if contentsCur != nil {
				_ = contentsCur.All(ctx, &contents)
			}
			totalContent += len(contents)
			for _, ct := range contents {
				res := h.coll(models.CollCourseProgresses).FindOne(ctx, bson.M{"student_id": student.ID, "content_id": ct.ID})
				if res.Err() == nil {
					doneContent++
				}
			}
		}
		percentage := 0.0
		if totalContent > 0 {
			percentage = round2(float64(doneContent) / float64(totalContent) * 100)
		}
		results = append(results, map[string]any{
			"course_id": course.ID.Hex(), "course_name": course.Name,
			"total_content": totalContent, "completed_content": doneContent, "percentage": percentage,
		})
	}
	return response.OK(c, "Course progress", results)
}

func (h *handler) getIndividualCourseProgress(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Course not found", "Not found")
	}
	user := middleware.UserFrom(c)
	studentID, _ := dbutil.ParseID(user.ID)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	course, err := h.getCourse(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch course", err.Error())
	}
	if course == nil {
		return response.Err(c, http.StatusNotFound, "Course not found", "Not found")
	}

	sectionsCur, err := h.coll(models.CollCourseSections).Find(ctx, bson.M{"course_id": course.ID})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch sections", err.Error())
	}
	var sections []models.CourseSection
	_ = sectionsCur.All(ctx, &sections)

	results := []map[string]any{}
	for _, sec := range sections {
		contentsCur, _ := h.coll(models.CollCourseLearningContents).Find(ctx, bson.M{"section_id": sec.ID})
		var contents []models.CourseLearningContent
		if contentsCur != nil {
			_ = contentsCur.All(ctx, &contents)
		}
		contentProgress := make([]map[string]any, 0, len(contents))
		done := 0
		for _, ct := range contents {
			res := h.coll(models.CollCourseProgresses).FindOne(ctx, bson.M{"student_id": studentID, "content_id": ct.ID})
			completed := res.Err() == nil
			if completed {
				done++
			}
			contentProgress = append(contentProgress, map[string]any{
				"content_id": ct.ID.Hex(), "name": ct.Name, "completed": completed,
			})
		}
		total := len(contents)
		percentage := 0.0
		if total > 0 {
			percentage = round2(float64(done) / float64(total) * 100)
		}
		results = append(results, map[string]any{
			"section_id": sec.ID.Hex(), "section_name": sec.Name,
			"total": total, "completed": done, "percentage": percentage, "contents": contentProgress,
		})
	}
	return response.OK(c, "Individual course progress", results)
}

// ================= Enrollments =================

func (h *handler) getCourseEnrollments(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Course not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	course, err := h.getCourse(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch course", err.Error())
	}
	if course == nil {
		return response.Err(c, http.StatusNotFound, "Course not found", "Not found")
	}

	cur, err := h.coll(models.CollCourseEnrollments).Find(ctx, bson.M{"course_id": course.ID})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch enrollments", err.Error())
	}
	var enrollments []models.CourseEnrollment
	_ = cur.All(ctx, &enrollments)

	data := make([]map[string]any, 0, len(enrollments))
	for _, e := range enrollments {
		batch, _ := h.getBatch(ctx, e.BatchID)
		var institution *models.Institution
		if e.InstitutionID != nil {
			institution, _ = h.getInstitution(ctx, *e.InstitutionID)
		} else if batch != nil {
			institution, _ = h.getInstitution(ctx, batch.InstitutionID)
		}

		instID, instName := "", "Unknown Institution"
		if institution != nil {
			instID, instName = institution.ID.Hex(), institution.Name
		} else if e.InstitutionID != nil {
			instID = e.InstitutionID.Hex()
		}
		batchID, batchName, batchBranch := e.BatchID.Hex(), "", ""
		if batch != nil {
			batchID, batchName, batchBranch = batch.ID.Hex(), batch.BatchName, batch.Branch
		}

		var batchNameField any
		if batch != nil {
			batchNameField = batch.BatchName
		}
		var institutionIDField any
		if e.InstitutionID != nil {
			institutionIDField = e.InstitutionID.Hex()
		}

		data = append(data, map[string]any{
			"institution": map[string]any{"id": instID, "name": instName},
			"batch":       map[string]any{"id": batchID, "batchname": batchName, "branch": batchBranch},

			"id": e.ID.Hex(), "course_id": e.CourseID.Hex(), "batch_id": e.BatchID.Hex(),
			"batch_name": batchNameField, "institution_id": institutionIDField,
			"enrolled_at": iso(e.EnrolledAt),
		})
	}

	var createdAt any
	if !course.CreatedAt.IsZero() {
		createdAt = iso(course.CreatedAt)
	}
	return response.OK(c, "Enrollments fetched", map[string]any{
		"course": map[string]any{
			"id": course.ID.Hex(), "name": course.Name, "description": course.Description,
			"level": course.Level, "duration": course.Duration, "thumbnail": course.Thumbnail,
			"instructor_name": course.InstructorName, "certificate": course.Certificate,
			"is_published": course.IsPublished, "created_at": createdAt,
		},
		"data": data,
	})
}

func (h *handler) getCourseEnrollmentsByBatch(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.OK(c, "Enrollments fetched", []map[string]any{})
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	cur, err := h.coll(models.CollCourseEnrollments).Find(ctx, bson.M{"batch_id": oid})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch enrollments", err.Error())
	}
	var enrollments []models.CourseEnrollment
	_ = cur.All(ctx, &enrollments)

	data := make([]map[string]any, 0, len(enrollments))
	for _, e := range enrollments {
		course, _ := h.getCourse(ctx, e.CourseID)
		var courseName, instructorName, level, createdAt any
		if course != nil {
			courseName, instructorName, level = course.Name, course.InstructorName, course.Level
			createdAt = iso(course.CreatedAt)
		}
		data = append(data, map[string]any{
			"id": e.ID.Hex(), "course_id": e.CourseID.Hex(), "course_name": courseName,
			"instructor_name": instructorName, "level": level, "created_at": createdAt,
			"batch_id": e.BatchID.Hex(), "enrolled_at": iso(e.EnrolledAt),
		})
	}
	return response.OK(c, "Enrollments fetched", data)
}

func (h *handler) addCourseEnrollment(c echo.Context) error {
	var body addEnrollmentRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	courseID, ok := dbutil.ParseID(body.CourseID)
	if !ok {
		return response.Err(c, http.StatusNotFound, "Course not found", "Not found")
	}
	batchID, ok := dbutil.ParseID(body.BatchID)
	if !ok {
		return response.Err(c, http.StatusNotFound, "Batch not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	course, err := h.getCourse(ctx, courseID)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch course", err.Error())
	}
	if course == nil {
		return response.Err(c, http.StatusNotFound, "Course not found", "Not found")
	}
	if course.IsPublished != models.CoursePublished {
		return response.Err(c, http.StatusBadRequest, "Publish course before assigning it to a batch", "Course is not published")
	}

	batch, err := h.getBatch(ctx, batchID)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch batch", err.Error())
	}
	if batch == nil {
		return response.Err(c, http.StatusNotFound, "Batch not found", "Not found")
	}

	existing := h.coll(models.CollCourseEnrollments).FindOne(ctx, bson.M{"course_id": courseID, "batch_id": batchID})
	if existing.Err() == nil {
		return response.Err(c, http.StatusBadRequest, "Already enrolled", "Duplicate enrollment")
	}

	now := time.Now().UTC()
	enrollment := models.CourseEnrollment{
		ID: primitive.NewObjectID(), CourseID: courseID, BatchID: batchID,
		EnrolledAt: now, CreatedAt: now, UpdatedAt: now,
	}
	if body.InstitutionID != nil {
		if instID, ok := dbutil.ParseID(*body.InstitutionID); ok {
			enrollment.InstitutionID = &instID
		}
	}
	if _, err := h.coll(models.CollCourseEnrollments).InsertOne(ctx, enrollment); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to create enrollment", err.Error())
	}
	return response.JSON(c, http.StatusCreated, "Enrollment created", map[string]any{"id": enrollment.ID.Hex()}, nil)
}

func (h *handler) removeCourseEnrollment(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Enrollment not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	res, err := h.coll(models.CollCourseEnrollments).DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to remove enrollment", err.Error())
	}
	if res.DeletedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Enrollment not found", "Not found")
	}
	return response.OK(c, "Enrollment removed", nil)
}
