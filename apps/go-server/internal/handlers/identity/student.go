package identity

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/bitwiselearn/go-server/internal/auth"
	"github.com/bitwiselearn/go-server/internal/dbutil"
	"github.com/bitwiselearn/go-server/internal/middleware"
	"github.com/bitwiselearn/go-server/internal/models"
	"github.com/bitwiselearn/go-server/internal/response"
)

func registerStudentRoutes(g *echo.Group, am *middleware.Auth, h *handler) {
	authed := am.Required
	notStudent := am.NotStudent()
	manageRoles := am.RequireRoles(models.RoleSuperadmin, models.RoleAdmin, models.RoleInstitution, models.RoleVendor)

	g.POST("/create-student", h.createStudent, authed, notStudent)
	g.GET("/get-all-student", h.getAllStudents, authed)
	g.GET("/get-student-by-id/:id", h.getStudentByID, authed)
	g.GET("/get-student-by-batch/:id", h.getStudentsByBatch, authed)
	g.PUT("/update-student-by-id/:id", h.updateStudent, authed)
	g.DELETE("/delete-student-by-id/:id", h.deleteStudent, authed, manageRoles)
	g.GET("/dashboard", h.studentDashboard, authed)
}

func (h *handler) createStudent(c echo.Context) error {
	var body createStudentRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	instID, ok := dbutil.ParseID(body.InstitutionID)
	if !ok {
		return response.Err(c, http.StatusNotFound, "Institution not found", "Not found")
	}
	batchID, ok := dbutil.ParseID(body.BatchID)
	if !ok {
		return response.Err(c, http.StatusNotFound, "Batch not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	if inst, err := h.getInstitution(ctx, instID); err != nil || inst == nil {
		return response.Err(c, http.StatusNotFound, "Institution not found", "Not found")
	}
	if batch, err := h.getBatch(ctx, batchID); err != nil || batch == nil {
		return response.Err(c, http.StatusNotFound, "Batch not found", "Not found")
	}

	existing := h.coll(models.CollStudents).FindOne(ctx, bson.M{"email": body.Email})
	if existing.Err() == nil {
		return response.Err(c, http.StatusBadRequest, "Email already exists", "Duplicate email")
	}

	rawPassword := genPassword()
	if body.LoginPassword != nil && *body.LoginPassword != "" {
		rawPassword = *body.LoginPassword
	}
	hashed, err := auth.HashPassword(rawPassword)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to create student", err.Error())
	}
	cloudPlatform := body.CloudPlatform
	if cloudPlatform == "" {
		cloudPlatform = "AWS"
	}

	now := time.Now().UTC()
	student := models.Student{
		ID: primitive.NewObjectID(), Name: body.Name, RollNumber: body.RollNumber, Email: body.Email,
		LoginPassword: hashed, CloudPlatform: cloudPlatform, BatchID: batchID, InstituteID: instID,
		CreatedAt: now, UpdatedAt: now,
	}
	if _, err := h.coll(models.CollStudents).InsertOne(ctx, student); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to create student", err.Error())
	}

	h.publishWelcomeEmail(ctx, body.Email, body.Name, rawPassword, "Student")

	return response.JSON(c, http.StatusCreated, "Student created", map[string]any{
		"id": student.ID.Hex(), "name": student.Name, "email": student.Email,
	}, nil)
}

func (h *handler) getAllStudents(c echo.Context) error {
	user := middleware.UserFrom(c)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	filter := bson.M{}
	if user.Type == models.RoleInstitution {
		if oid, ok := dbutil.ParseID(user.ID); ok {
			filter["institute_id"] = oid
		}
	}
	cur, err := h.coll(models.CollStudents).Find(ctx, filter)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch students", err.Error())
	}
	var students []models.Student
	if err := cur.All(ctx, &students); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch students", err.Error())
	}
	data := make([]map[string]any, 0, len(students))
	for _, s := range students {
		data = append(data, map[string]any{
			"id": s.ID.Hex(), "name": s.Name, "email": s.Email, "roll_number": s.RollNumber,
			"batch_id": s.BatchID.Hex(), "institute_id": s.InstituteID.Hex(),
			"cloud_platform": s.CloudPlatform, "created_at": iso(s.CreatedAt),
		})
	}
	return response.OK(c, "Students fetched", data)
}

func (h *handler) getStudentByID(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Student not found", "Not found")
	}
	user := middleware.UserFrom(c)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	student, err := h.getStudent(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch student", err.Error())
	}
	if student == nil {
		return response.Err(c, http.StatusNotFound, "Student not found", "Not found")
	}

	switch user.Type {
	case models.RoleStudent:
		if student.ID.Hex() != user.ID {
			return response.Err(c, http.StatusForbidden, "Not authorized", "Forbidden")
		}
	case models.RoleInstitution:
		if student.InstituteID.Hex() != user.ID {
			return response.Err(c, http.StatusForbidden, "Not authorized", "Forbidden")
		}
	case models.RoleVendor:
		inst, err := h.getInstitution(ctx, student.InstituteID)
		if err != nil || inst == nil || inst.CreatedByVendorID == nil || inst.CreatedByVendorID.Hex() != user.ID {
			return response.Err(c, http.StatusForbidden, "Not authorized", "Forbidden")
		}
	case models.RoleTeacher:
		tOID, ok := dbutil.ParseID(user.ID)
		if !ok {
			return response.Err(c, http.StatusForbidden, "Not authorized", "Forbidden")
		}
		teacher, err := h.getTeacher(ctx, tOID)
		if err != nil || teacher == nil || teacher.BatchID.Hex() != student.BatchID.Hex() {
			return response.Err(c, http.StatusForbidden, "Not authorized", "Forbidden")
		}
	}

	return response.OK(c, "Student fetched", map[string]any{
		"id": student.ID.Hex(), "name": student.Name, "email": student.Email,
		"roll_number": student.RollNumber, "batch_id": student.BatchID.Hex(),
		"institute_id": student.InstituteID.Hex(), "cloud_platform": student.CloudPlatform,
		"cloudname": student.CloudName, "cloudpass": student.CloudPass, "cloudurl": student.CloudURL,
		"created_at": iso(student.CreatedAt),
	})
}

func (h *handler) getStudentsByBatch(c echo.Context) error {
	batchID, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Batch not found", "Not found")
	}
	user := middleware.UserFrom(c)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	batch, err := h.getBatch(ctx, batchID)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch batch", err.Error())
	}
	if batch == nil {
		return response.Err(c, http.StatusNotFound, "Batch not found", "Not found")
	}

	switch user.Type {
	case models.RoleStudent:
		uOID, ok := dbutil.ParseID(user.ID)
		if !ok {
			return response.Err(c, http.StatusForbidden, "Not authorized", "Forbidden")
		}
		me, err := h.getStudent(ctx, uOID)
		if err != nil || me == nil || me.BatchID.Hex() != c.Param("id") {
			return response.Err(c, http.StatusForbidden, "Not authorized", "Forbidden")
		}
	case models.RoleTeacher:
		tOID, ok := dbutil.ParseID(user.ID)
		if !ok {
			return response.Err(c, http.StatusForbidden, "Not authorized", "Forbidden")
		}
		teacher, err := h.getTeacher(ctx, tOID)
		if err != nil || teacher == nil || teacher.BatchID.Hex() != c.Param("id") {
			return response.Err(c, http.StatusForbidden, "Not authorized", "Forbidden")
		}
	case models.RoleInstitution:
		if batch.InstitutionID.Hex() != user.ID {
			return response.Err(c, http.StatusForbidden, "Not authorized", "Forbidden")
		}
	case models.RoleVendor:
		inst, err := h.getInstitution(ctx, batch.InstitutionID)
		if err != nil || inst == nil || inst.CreatedByVendorID == nil || inst.CreatedByVendorID.Hex() != user.ID {
			return response.Err(c, http.StatusForbidden, "Not authorized", "Forbidden")
		}
	}

	cur, err := h.coll(models.CollStudents).Find(ctx, bson.M{"batch_id": batchID})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch students", err.Error())
	}
	var students []models.Student
	if err := cur.All(ctx, &students); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch students", err.Error())
	}
	data := make([]map[string]any, 0, len(students))
	for _, s := range students {
		data = append(data, map[string]any{"id": s.ID.Hex(), "name": s.Name, "email": s.Email, "roll_number": s.RollNumber})
	}
	return response.OK(c, "Students fetched", data)
}

func (h *handler) updateStudent(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Student not found", "Not found")
	}
	var body updateStudentRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	set := bson.M{"updated_at": time.Now().UTC()}
	if body.Name != nil {
		set["name"] = *body.Name
	}
	if body.RollNumber != nil {
		set["roll_number"] = *body.RollNumber
	}
	if body.Email != nil {
		set["email"] = *body.Email
	}
	if body.CloudPlatform != nil {
		set["cloud_platform"] = *body.CloudPlatform
	}
	if body.CloudName != nil {
		set["cloudname"] = *body.CloudName
	}
	if body.CloudPass != nil {
		set["cloudpass"] = *body.CloudPass
	}
	if body.CloudURL != nil {
		set["cloudurl"] = *body.CloudURL
	}
	res, err := h.coll(models.CollStudents).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": set})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to update student", err.Error())
	}
	if res.MatchedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Student not found", "Not found")
	}
	student, _ := h.getStudent(ctx, oid)
	name := ""
	if student != nil {
		name = student.Name
	}
	return response.OK(c, "Student updated", map[string]any{"id": oid.Hex(), "name": name})
}

func (h *handler) deleteStudent(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Student not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	res, err := h.coll(models.CollStudents).DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to delete student", err.Error())
	}
	if res.DeletedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Student not found", "Not found")
	}
	return response.OK(c, "Student deleted", nil)
}

func (h *handler) studentDashboard(c echo.Context) error {
	user := middleware.UserFrom(c)
	if user.Type != models.RoleStudent {
		return response.Err(c, http.StatusForbidden, "Not authorized", "Forbidden")
	}
	oid, ok := dbutil.ParseID(user.ID)
	if !ok {
		return response.Err(c, http.StatusForbidden, "Not authorized", "Forbidden")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	student, err := h.getStudent(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch student", err.Error())
	}
	if student == nil {
		return response.Err(c, http.StatusNotFound, "Student not found", "Not found")
	}
	batch, _ := h.getBatch(ctx, student.BatchID)
	inst, _ := h.getInstitution(ctx, student.InstituteID)

	var courses []map[string]any
	if cur, err := h.coll(models.CollCourseEnrollments).Find(ctx, bson.M{"batch_id": student.BatchID}); err == nil {
		var enrollments []models.CourseEnrollment
		_ = cur.All(ctx, &enrollments)
		for _, e := range enrollments {
			var course models.Course
			if err := h.coll(models.CollCourses).FindOne(ctx, bson.M{"_id": e.CourseID}).Decode(&course); err == nil {
				courses = append(courses, map[string]any{"id": course.ID.Hex(), "name": course.Name, "level": course.Level})
			}
		}
	}

	var assessmentData []map[string]any
	if cur, err := h.coll(models.CollAssessments).Find(ctx, bson.M{"batch_id": student.BatchID}); err == nil {
		var assessments []models.Assessment
		_ = cur.All(ctx, &assessments)
		for _, a := range assessments {
			assessmentData = append(assessmentData, map[string]any{
				"id": a.ID.Hex(), "name": a.Name, "status": a.Status,
				"start_time": iso(a.StartTime), "end_time": iso(a.EndTime),
			})
		}
	}

	var batchData, instData any
	if batch != nil {
		batchData = map[string]any{"id": batch.ID.Hex(), "batchname": batch.BatchName}
	}
	if inst != nil {
		instData = map[string]any{"id": inst.ID.Hex(), "name": inst.Name}
	}

	return response.OK(c, "Student dashboard", map[string]any{
		"student":     map[string]any{"id": student.ID.Hex(), "name": student.Name, "email": student.Email},
		"batch":       batchData,
		"institution": instData,
		"courses":     courses,
		"assessments": assessmentData,
	})
}
