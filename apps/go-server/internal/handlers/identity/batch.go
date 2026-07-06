package identity

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/bitwiselearn/go-server/internal/dbutil"
	"github.com/bitwiselearn/go-server/internal/middleware"
	"github.com/bitwiselearn/go-server/internal/models"
	"github.com/bitwiselearn/go-server/internal/response"
)

func registerBatchRoutes(g *echo.Group, am *middleware.Auth, h *handler) {
	authed := am.Required
	notStudent := am.NotStudent()
	deleteRoles := am.RequireRoles(models.RoleSuperadmin, models.RoleAdmin, models.RoleInstitution, models.RoleVendor)

	g.POST("/create-batch", h.createBatch, authed, notStudent)
	g.GET("/get-all-batch", h.getAllBatchesGlobal, authed, notStudent)
	g.GET("/get-all-batch/:id", h.getAllBatchesByInstitution, authed, notStudent)
	g.GET("/get-batch-by-id/:id", h.getBatchByID, authed)
	g.PUT("/update-batch-by-id/:id", h.updateBatch, authed, notStudent)
	g.DELETE("/delete-batch-by-id/:id", h.deleteBatch, authed, deleteRoles)
}

func (h *handler) createBatch(c echo.Context) error {
	var body createBatchRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	instID, ok := dbutil.ParseID(body.InstitutionID)
	if !ok {
		return response.Err(c, http.StatusNotFound, "Institution not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	if inst, err := h.getInstitution(ctx, instID); err != nil || inst == nil {
		return response.Err(c, http.StatusNotFound, "Institution not found", "Not found")
	}

	now := time.Now().UTC()
	batch := models.Batch{
		ID: primitive.NewObjectID(), BatchName: body.BatchName, Branch: body.Branch,
		BatchEndYear: body.BatchEndYear, InstitutionID: instID, CreatedAt: now, UpdatedAt: now,
	}
	if _, err := h.coll(models.CollBatches).InsertOne(ctx, batch); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to create batch", err.Error())
	}
	return response.JSON(c, http.StatusCreated, "Batch created", map[string]any{
		"id": batch.ID.Hex(), "batchname": batch.BatchName, "branch": batch.Branch,
		"batch_end_year": batch.BatchEndYear, "institution_id": batch.InstitutionID.Hex(),
	}, nil)
}

func (h *handler) batchToMap(ctx context.Context, b models.Batch) map[string]any {
	studentCount := h.countColl(ctx, models.CollStudents, bson.M{"batch_id": b.ID})
	return map[string]any{
		"id": b.ID.Hex(), "batchname": b.BatchName, "branch": b.Branch,
		"batch_end_year": b.BatchEndYear, "institution_id": b.InstitutionID.Hex(),
		"student_count": studentCount, "created_at": iso(b.CreatedAt),
	}
}

func (h *handler) getAllBatchesGlobal(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	cur, err := h.coll(models.CollBatches).Find(ctx, bson.M{})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch batches", err.Error())
	}
	var batches []models.Batch
	if err := cur.All(ctx, &batches); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch batches", err.Error())
	}
	data := make([]map[string]any, 0, len(batches))
	for _, b := range batches {
		data = append(data, h.batchToMap(ctx, b))
	}
	return response.OK(c, "Batches fetched", data)
}

func (h *handler) getAllBatchesByInstitution(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.OK(c, "Batches fetched", []map[string]any{})
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	cur, err := h.coll(models.CollBatches).Find(ctx, bson.M{"institution_id": oid})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch batches", err.Error())
	}
	var batches []models.Batch
	if err := cur.All(ctx, &batches); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch batches", err.Error())
	}
	data := make([]map[string]any, 0, len(batches))
	for _, b := range batches {
		data = append(data, h.batchToMap(ctx, b))
	}
	return response.OK(c, "Batches fetched", data)
}

func (h *handler) getBatchByID(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Batch not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	batch, err := h.getBatch(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch batch", err.Error())
	}
	if batch == nil {
		return response.Err(c, http.StatusNotFound, "Batch not found", "Not found")
	}

	cur, err := h.coll(models.CollStudents).Find(ctx, bson.M{"batch_id": batch.ID})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch batch", err.Error())
	}
	var students []models.Student
	_ = cur.All(ctx, &students)
	studentData := make([]map[string]any, 0, len(students))
	for _, s := range students {
		studentData = append(studentData, map[string]any{
			"id": s.ID.Hex(), "name": s.Name, "email": s.Email, "roll_number": s.RollNumber,
		})
	}

	return response.OK(c, "Batch fetched", map[string]any{
		"id": batch.ID.Hex(), "batchname": batch.BatchName, "branch": batch.Branch,
		"batch_end_year": batch.BatchEndYear, "institution_id": batch.InstitutionID.Hex(),
		"students": studentData, "created_at": iso(batch.CreatedAt),
	})
}

func (h *handler) updateBatch(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Batch not found", "Not found")
	}
	var body updateBatchRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	set := bson.M{"updated_at": time.Now().UTC()}
	if body.BatchName != nil {
		set["batchname"] = *body.BatchName
	}
	if body.Branch != nil {
		set["branch"] = *body.Branch
	}
	if body.BatchEndYear != nil {
		set["batch_end_year"] = *body.BatchEndYear
	}
	res, err := h.coll(models.CollBatches).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": set})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to update batch", err.Error())
	}
	if res.MatchedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Batch not found", "Not found")
	}
	batch, _ := h.getBatch(ctx, oid)
	name := ""
	if batch != nil {
		name = batch.BatchName
	}
	return response.OK(c, "Batch updated", map[string]any{"id": oid.Hex(), "batchname": name})
}

func (h *handler) deleteBatch(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Batch not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	batch, err := h.getBatch(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch batch", err.Error())
	}
	if batch == nil {
		return response.Err(c, http.StatusNotFound, "Batch not found", "Not found")
	}
	_, _ = h.coll(models.CollStudents).DeleteMany(ctx, bson.M{"batch_id": oid})
	_, _ = h.coll(models.CollTeachers).DeleteMany(ctx, bson.M{"batch_id": oid})
	_, _ = h.coll(models.CollCourseEnrollments).DeleteMany(ctx, bson.M{"batch_id": oid})
	if _, err := h.coll(models.CollBatches).DeleteOne(ctx, bson.M{"_id": oid}); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to delete batch", err.Error())
	}
	return response.OK(c, "Batch deleted", nil)
}
