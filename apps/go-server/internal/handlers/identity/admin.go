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

func registerAdminRoutes(g *echo.Group, am *middleware.Auth, h *handler) {
	authed := am.Required
	superadminOnly := am.SuperadminOnly()
	adminOnly := am.AdminOnly()

	g.POST("/create-admin", h.createAdmin, authed, superadminOnly)
	g.GET("/get-all-admin", h.getAllAdmins, authed, superadminOnly)
	g.GET("/get-admin-by-id/:id", h.getAdminByID, authed, superadminOnly)
	g.PUT("/update-admin-by-id/:id", h.updateAdmin, authed, superadminOnly)
	g.DELETE("/delete-admin-by-id/:id", h.deleteAdmin, authed, superadminOnly)
	g.GET("/db-info", h.dbInfo, authed, adminOnly)
	g.GET("/dashboard", h.adminDashboard, authed, adminOnly)
}

func (h *handler) getUser(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var u models.User
	err := h.coll(models.CollUsers).FindOne(ctx, bson.M{"_id": id}).Decode(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (h *handler) createAdmin(c echo.Context) error {
	var body createAdminRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	existing := h.coll(models.CollUsers).FindOne(ctx, bson.M{"email": body.Email})
	if existing.Err() == nil {
		return response.Err(c, http.StatusBadRequest, "Email already exists", "Duplicate email")
	}

	rawPassword := genPassword()
	hashed, err := auth.HashPassword(rawPassword)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to create admin", err.Error())
	}
	role := body.Role
	if role == "" {
		role = models.RoleAdmin
	}

	now := time.Now().UTC()
	admin := models.User{
		ID: primitive.NewObjectID(), Name: body.Name, Email: body.Email, Password: hashed,
		Role: role, CreatedAt: now, UpdatedAt: now,
	}
	if _, err := h.coll(models.CollUsers).InsertOne(ctx, admin); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to create admin", err.Error())
	}

	h.publishWelcomeEmail(ctx, body.Email, body.Name, rawPassword, "Admin")

	return response.JSON(c, http.StatusCreated, "Admin created successfully", map[string]any{
		"id": admin.ID.Hex(), "name": admin.Name, "email": admin.Email, "role": admin.Role,
	}, nil)
}

func (h *handler) getAllAdmins(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	cur, err := h.coll(models.CollUsers).Find(ctx, bson.M{"role": models.RoleAdmin})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch admins", err.Error())
	}
	var admins []models.User
	if err := cur.All(ctx, &admins); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch admins", err.Error())
	}
	data := make([]map[string]any, 0, len(admins))
	for _, a := range admins {
		var createdAt any
		if !a.CreatedAt.IsZero() {
			createdAt = isoMillis(a.CreatedAt)
		}
		data = append(data, map[string]any{
			"id": a.ID.Hex(), "name": a.Name, "email": a.Email, "role": a.Role, "created_at": createdAt,
		})
	}
	return response.OK(c, "Admins fetched", data)
}

func (h *handler) getAdminByID(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Admin not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	admin, err := h.getUser(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusNotFound, "Admin not found", "Not found")
	}
	return response.OK(c, "Admin fetched", map[string]any{
		"id": admin.ID.Hex(), "name": admin.Name, "email": admin.Email, "role": admin.Role,
		"created_at": iso(admin.CreatedAt),
	})
}

func (h *handler) updateAdmin(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Admin not found", "Not found")
	}
	var body updateAdminRequest
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
	res, err := h.coll(models.CollUsers).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": set})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to update admin", err.Error())
	}
	if res.MatchedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Admin not found", "Not found")
	}
	admin, err := h.getUser(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusNotFound, "Admin not found", "Not found")
	}
	return response.OK(c, "Admin updated", map[string]any{
		"id": admin.ID.Hex(), "name": admin.Name, "email": admin.Email, "role": admin.Role,
	})
}

func (h *handler) deleteAdmin(c echo.Context) error {
	user := middleware.UserFrom(c)
	id := c.Param("id")
	if user.ID == id {
		return response.Err(c, http.StatusBadRequest, "Cannot delete yourself", "Self-deletion not allowed")
	}
	oid, ok := dbutil.ParseID(id)
	if !ok {
		return response.Err(c, http.StatusNotFound, "Admin not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	res, err := h.coll(models.CollUsers).DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to delete admin", err.Error())
	}
	if res.DeletedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Admin not found", "Not found")
	}
	return response.OK(c, "Admin deleted", nil)
}

func (h *handler) dbInfo(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	cur, err := h.coll(models.CollInstitutions).Find(ctx, bson.M{})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch db info", err.Error())
	}
	var institutions []models.Institution
	_ = cur.All(ctx, &institutions)

	instData := make([]map[string]any, 0, len(institutions))
	for _, inst := range institutions {
		bcur, err := h.coll(models.CollBatches).Find(ctx, bson.M{"institution_id": inst.ID})
		var batches []models.Batch
		if err == nil {
			_ = bcur.All(ctx, &batches)
		}
		batchData := make([]map[string]any, 0, len(batches))
		for _, b := range batches {
			batchData = append(batchData, map[string]any{
				"id": b.ID.Hex(), "batchname": b.BatchName, "branch": b.Branch,
				"student_count": h.countColl(ctx, models.CollStudents, bson.M{"batch_id": b.ID}),
				"teacher_count": h.countColl(ctx, models.CollTeachers, bson.M{"batch_id": b.ID}),
			})
		}
		instData = append(instData, map[string]any{
			"id": inst.ID.Hex(), "name": inst.Name, "email": inst.Email, "batches": batchData,
		})
	}

	return response.OK(c, "DB info fetched", map[string]any{
		"institutions":          h.countColl(ctx, models.CollInstitutions, nil),
		"vendors":               h.countColl(ctx, models.CollVendors, nil),
		"admins":                h.countColl(ctx, models.CollUsers, bson.M{"role": models.RoleAdmin}),
		"students":              h.countColl(ctx, models.CollStudents, nil),
		"teachers":              h.countColl(ctx, models.CollTeachers, nil),
		"batches":               h.countColl(ctx, models.CollBatches, nil),
		"courses":               h.countColl(ctx, models.CollCourses, nil),
		"assessments":           h.countColl(ctx, models.CollAssessments, nil),
		"institution_hierarchy": instData,
	})
}

func (h *handler) adminDashboard(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	return response.OK(c, "Dashboard data", map[string]any{
		"institutions": h.countColl(ctx, models.CollInstitutions, nil),
		"vendors":      h.countColl(ctx, models.CollVendors, nil),
		"students":     h.countColl(ctx, models.CollStudents, nil),
		"teachers":     h.countColl(ctx, models.CollTeachers, nil),
		"batches":      h.countColl(ctx, models.CollBatches, nil),
		"courses":      h.countColl(ctx, models.CollCourses, nil),
		"assessments":  h.countColl(ctx, models.CollAssessments, nil),
	})
}
