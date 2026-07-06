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

func registerTeacherRoutes(g *echo.Group, am *middleware.Auth, h *handler) {
	authed := am.Required
	notStudent := am.NotStudent()

	g.POST("/create-teacher", h.createTeacher, authed, notStudent)
	g.GET("/get-all-teacher", h.getAllTeachers, authed)
	g.GET("/get-teacher-by-id/:id", h.getTeacherByID, authed)
	g.GET("/get-teacher-by-institute/:id", h.getTeachersByInstitute, authed)
	g.GET("/get-teacher-by-batch/:id", h.getTeachersByBatch, authed)
	g.PUT("/update-teacher-by-id/:id", h.updateTeacher, authed, notStudent)
	g.DELETE("/delete-teacher-by-id/:id", h.deleteTeacher, authed, notStudent)
}

func (h *handler) createTeacher(c echo.Context) error {
	var body createTeacherRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	instID, ok := dbutil.ParseID(body.InstituteID)
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
	existing := h.coll(models.CollTeachers).FindOne(ctx, bson.M{"email": body.Email})
	if existing.Err() == nil {
		return response.Err(c, http.StatusBadRequest, "Email already exists", "Duplicate email")
	}

	rawPassword := genPassword()
	hashed, err := auth.HashPassword(rawPassword)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to create teacher", err.Error())
	}

	now := time.Now().UTC()
	teacher := models.Teacher{
		ID: primitive.NewObjectID(), Name: body.Name, Email: body.Email, PhoneNumber: body.PhoneNumber,
		LoginPassword: hashed, InstituteID: instID, BatchID: batchID, CreatedAt: now, UpdatedAt: now,
	}
	if body.VendorID != nil {
		if vid, ok := dbutil.ParseID(*body.VendorID); ok {
			teacher.VendorID = &vid
		}
	}
	if _, err := h.coll(models.CollTeachers).InsertOne(ctx, teacher); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to create teacher", err.Error())
	}

	h.publishWelcomeEmail(ctx, body.Email, body.Name, rawPassword, "Teacher")

	return response.JSON(c, http.StatusCreated, "Teacher created", map[string]any{
		"id": teacher.ID.Hex(), "name": teacher.Name, "email": teacher.Email,
	}, nil)
}

func (h *handler) getAllTeachers(c echo.Context) error {
	user := middleware.UserFrom(c)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	filter := bson.M{}
	if user.Type == models.RoleInstitution {
		if oid, ok := dbutil.ParseID(user.ID); ok {
			filter["institute_id"] = oid
		}
	} else if user.Type == models.RoleVendor {
		if oid, ok := dbutil.ParseID(user.ID); ok {
			filter["vendor_id"] = oid
		}
	}
	cur, err := h.coll(models.CollTeachers).Find(ctx, filter)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch teachers", err.Error())
	}
	var teachers []models.Teacher
	if err := cur.All(ctx, &teachers); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch teachers", err.Error())
	}
	data := make([]map[string]any, 0, len(teachers))
	for _, t := range teachers {
		data = append(data, map[string]any{
			"id": t.ID.Hex(), "name": t.Name, "email": t.Email, "phone_number": t.PhoneNumber,
			"institute_id": t.InstituteID.Hex(), "batch_id": t.BatchID.Hex(), "created_at": iso(t.CreatedAt),
		})
	}
	return response.OK(c, "Teachers fetched", data)
}

func (h *handler) getTeacherByID(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Teacher not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	teacher, err := h.getTeacher(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch teacher", err.Error())
	}
	if teacher == nil {
		return response.Err(c, http.StatusNotFound, "Teacher not found", "Not found")
	}
	var vendorID any
	if teacher.VendorID != nil {
		vendorID = teacher.VendorID.Hex()
	}
	return response.OK(c, "Teacher fetched", map[string]any{
		"id": teacher.ID.Hex(), "name": teacher.Name, "email": teacher.Email,
		"phone_number": teacher.PhoneNumber, "institute_id": teacher.InstituteID.Hex(),
		"batch_id": teacher.BatchID.Hex(), "vendor_id": vendorID, "created_at": iso(teacher.CreatedAt),
	})
}

func (h *handler) getTeachersByInstitute(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.OK(c, "Teachers fetched", []map[string]any{})
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	cur, err := h.coll(models.CollTeachers).Find(ctx, bson.M{"institute_id": oid})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch teachers", err.Error())
	}
	var teachers []models.Teacher
	_ = cur.All(ctx, &teachers)
	data := make([]map[string]any, 0, len(teachers))
	for _, t := range teachers {
		data = append(data, map[string]any{
			"id": t.ID.Hex(), "name": t.Name, "email": t.Email,
			"phone_number": t.PhoneNumber, "batch_id": t.BatchID.Hex(),
		})
	}
	return response.OK(c, "Teachers fetched", data)
}

func (h *handler) getTeachersByBatch(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.OK(c, "Teachers fetched", []map[string]any{})
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	cur, err := h.coll(models.CollTeachers).Find(ctx, bson.M{"batch_id": oid})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch teachers", err.Error())
	}
	var teachers []models.Teacher
	_ = cur.All(ctx, &teachers)
	data := make([]map[string]any, 0, len(teachers))
	for _, t := range teachers {
		data = append(data, map[string]any{"id": t.ID.Hex(), "name": t.Name, "email": t.Email, "phone_number": t.PhoneNumber})
	}
	return response.OK(c, "Teachers fetched", data)
}

func (h *handler) updateTeacher(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Teacher not found", "Not found")
	}
	var body updateTeacherRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	set := bson.M{"updated_at": time.Now().UTC()}
	if body.Name != nil {
		set["name"] = *body.Name
	}
	if body.Email != nil {
		set["email"] = *body.Email
	}
	if body.PhoneNumber != nil {
		set["phone_number"] = *body.PhoneNumber
	}
	res, err := h.coll(models.CollTeachers).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": set})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to update teacher", err.Error())
	}
	if res.MatchedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Teacher not found", "Not found")
	}
	teacher, _ := h.getTeacher(ctx, oid)
	name := ""
	if teacher != nil {
		name = teacher.Name
	}
	return response.OK(c, "Teacher updated", map[string]any{"id": oid.Hex(), "name": name})
}

func (h *handler) deleteTeacher(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Teacher not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	res, err := h.coll(models.CollTeachers).DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to delete teacher", err.Error())
	}
	if res.DeletedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Teacher not found", "Not found")
	}
	return response.OK(c, "Teacher deleted", nil)
}
