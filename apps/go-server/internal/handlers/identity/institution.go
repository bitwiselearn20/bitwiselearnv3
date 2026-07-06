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

func registerInstitutionRoutes(g *echo.Group, am *middleware.Auth, h *handler) {
	authed := am.Required
	createRoles := am.RequireRoles(models.RoleSuperadmin, models.RoleAdmin, models.RoleVendor)
	deleteRoles := am.RequireRoles(models.RoleSuperadmin, models.RoleAdmin, models.RoleVendor)
	updateRoles := am.RequireRoles(models.RoleSuperadmin, models.RoleAdmin, models.RoleInstitution, models.RoleVendor)
	readRoles := updateRoles

	g.POST("/create-institution", h.createInstitution, authed, createRoles)
	g.GET("/get-all-institution", h.getAllInstitutions, authed, readRoles)
	g.GET("/get-institution-by-id/:id", h.getInstitutionByID, authed)
	g.GET("/get-institution-by-vendor/:id", h.getInstitutionsByVendor, authed, readRoles)
	g.PUT("/update-institution-by-id/:id", h.updateInstitution, authed, updateRoles)
	g.DELETE("/delete-institution-by-id/:id", h.deleteInstitution, authed, deleteRoles)
	g.GET("/dashboard/:institution_id", h.institutionDashboard, authed)
}

func (h *handler) createInstitution(c echo.Context) error {
	var body createInstitutionRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	user := middleware.UserFrom(c)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	existing := h.coll(models.CollInstitutions).FindOne(ctx, bson.M{"email": body.Email})
	if existing.Err() == nil {
		return response.Err(c, http.StatusBadRequest, "Email already exists", "Duplicate email")
	}

	rawPassword := genPassword()
	hashed, err := auth.HashPassword(rawPassword)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to create institution", err.Error())
	}

	createdBy, _ := dbutil.ParseID(user.ID)
	now := time.Now().UTC()
	inst := models.Institution{
		ID: primitive.NewObjectID(), Name: body.Name, Address: body.Address, PinCode: body.PinCode,
		Tagline: body.Tagline, WebsiteLink: body.WebsiteLink, LoginPassword: hashed, Email: body.Email,
		SecondaryEmail: body.SecondaryEmail, PhoneNumber: body.PhoneNumber,
		SecondaryPhoneNumber: body.SecondaryPhoneNumber, CreatedBy: createdBy,
		CreatedAt: now, UpdatedAt: now,
	}
	if user.Type == models.RoleVendor {
		vid := createdBy
		inst.CreatedByVendorID = &vid
	}
	if _, err := h.coll(models.CollInstitutions).InsertOne(ctx, inst); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to create institution", err.Error())
	}

	h.publishWelcomeEmail(ctx, body.Email, body.Name, rawPassword, "Institution")

	return response.JSON(c, http.StatusCreated, "Institution created", map[string]any{
		"id": inst.ID.Hex(), "name": inst.Name, "email": inst.Email,
	}, nil)
}

func institutionToMap(i models.Institution) map[string]any {
	var vendorID any
	if i.CreatedByVendorID != nil {
		vendorID = i.CreatedByVendorID.Hex()
	}
	var createdAt any
	if !i.CreatedAt.IsZero() {
		createdAt = isoMillis(i.CreatedAt)
	}
	return map[string]any{
		"id": i.ID.Hex(), "name": i.Name, "email": i.Email, "address": i.Address,
		"pin_code": i.PinCode, "tagline": i.Tagline, "website_link": i.WebsiteLink,
		"phone_number": i.PhoneNumber, "secondary_phone_number": i.SecondaryPhoneNumber,
		"secondary_email": i.SecondaryEmail, "created_by": i.CreatedBy.Hex(),
		"created_by_vendor_id": vendorID, "created_at": createdAt,
	}
}

func (h *handler) getAllInstitutions(c echo.Context) error {
	user := middleware.UserFrom(c)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	var institutions []models.Institution
	switch user.Type {
	case models.RoleVendor:
		oid, _ := dbutil.ParseID(user.ID)
		cur, err := h.coll(models.CollInstitutions).Find(ctx, bson.M{"created_by_vendor_id": oid})
		if err != nil {
			return response.Err(c, http.StatusInternalServerError, "Failed to fetch institutions", err.Error())
		}
		_ = cur.All(ctx, &institutions)
	case models.RoleInstitution:
		oid, ok := dbutil.ParseID(user.ID)
		if ok {
			inst, err := h.getInstitution(ctx, oid)
			if err != nil {
				return response.Err(c, http.StatusInternalServerError, "Failed to fetch institutions", err.Error())
			}
			if inst != nil {
				institutions = []models.Institution{*inst}
			}
		}
	default:
		cur, err := h.coll(models.CollInstitutions).Find(ctx, bson.M{})
		if err != nil {
			return response.Err(c, http.StatusInternalServerError, "Failed to fetch institutions", err.Error())
		}
		_ = cur.All(ctx, &institutions)
	}

	data := make([]map[string]any, 0, len(institutions))
	for _, i := range institutions {
		data = append(data, institutionToMap(i))
	}
	return response.OK(c, "Institutions fetched", data)
}

func (h *handler) getInstitutionByID(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Institution not found", "Not found")
	}
	user := middleware.UserFrom(c)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	inst, err := h.getInstitution(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch institution", err.Error())
	}
	if inst == nil {
		return response.Err(c, http.StatusNotFound, "Institution not found", "Not found")
	}

	switch user.Type {
	case models.RoleInstitution:
		if inst.ID.Hex() != user.ID {
			return response.Err(c, http.StatusForbidden, "Not authorized", "Forbidden")
		}
	case models.RoleVendor:
		if inst.CreatedByVendorID == nil || inst.CreatedByVendorID.Hex() != user.ID {
			return response.Err(c, http.StatusForbidden, "Not authorized", "Forbidden")
		}
	case models.RoleStudent:
		uOID, ok := dbutil.ParseID(user.ID)
		if !ok {
			return response.Err(c, http.StatusForbidden, "Not authorized", "Forbidden")
		}
		student, err := h.getStudent(ctx, uOID)
		if err != nil || student == nil || student.InstituteID.Hex() != inst.ID.Hex() {
			return response.Err(c, http.StatusForbidden, "Not authorized", "Forbidden")
		}
	case models.RoleTeacher:
		uOID, ok := dbutil.ParseID(user.ID)
		if !ok {
			return response.Err(c, http.StatusForbidden, "Not authorized", "Forbidden")
		}
		teacher, err := h.getTeacher(ctx, uOID)
		if err != nil || teacher == nil || teacher.InstituteID.Hex() != inst.ID.Hex() {
			return response.Err(c, http.StatusForbidden, "Not authorized", "Forbidden")
		}
	}

	return response.OK(c, "Institution fetched", institutionToMap(*inst))
}

func (h *handler) getInstitutionsByVendor(c echo.Context) error {
	id := c.Param("id")
	user := middleware.UserFrom(c)
	if user.Type == models.RoleVendor && user.ID != id {
		return response.Err(c, http.StatusForbidden, "Not authorized", "Forbidden")
	}
	oid, ok := dbutil.ParseID(id)
	if !ok {
		return response.OK(c, "Institutions fetched", []map[string]any{})
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	cur, err := h.coll(models.CollInstitutions).Find(ctx, bson.M{"created_by_vendor_id": oid})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch institutions", err.Error())
	}
	var institutions []models.Institution
	_ = cur.All(ctx, &institutions)
	data := make([]map[string]any, 0, len(institutions))
	for _, i := range institutions {
		var createdAt any
		if !i.CreatedAt.IsZero() {
			createdAt = isoMillis(i.CreatedAt)
		}
		data = append(data, map[string]any{
			"id": i.ID.Hex(), "name": i.Name, "email": i.Email,
			"address": i.Address, "phone_number": i.PhoneNumber, "created_at": createdAt,
		})
	}
	return response.OK(c, "Institutions fetched", data)
}

func (h *handler) updateInstitution(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Institution not found", "Not found")
	}
	var body updateInstitutionRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	user := middleware.UserFrom(c)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	inst, err := h.getInstitution(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch institution", err.Error())
	}
	if inst == nil {
		return response.Err(c, http.StatusNotFound, "Institution not found", "Not found")
	}
	if user.Type == models.RoleVendor && (inst.CreatedByVendorID == nil || inst.CreatedByVendorID.Hex() != user.ID) {
		return response.Err(c, http.StatusForbidden, "Not authorized", "Forbidden")
	}
	if user.Type == models.RoleInstitution && inst.ID.Hex() != user.ID {
		return response.Err(c, http.StatusForbidden, "Not authorized", "Forbidden")
	}

	set := bson.M{"updated_at": time.Now().UTC()}
	if body.Name != nil {
		set["name"] = *body.Name
	}
	if body.Address != nil {
		set["address"] = *body.Address
	}
	if body.PinCode != nil {
		set["pin_code"] = *body.PinCode
	}
	if body.Tagline != nil {
		set["tagline"] = *body.Tagline
	}
	if body.WebsiteLink != nil {
		set["website_link"] = *body.WebsiteLink
	}
	if body.Email != nil {
		set["email"] = *body.Email
	}
	if body.SecondaryEmail != nil {
		set["secondary_email"] = *body.SecondaryEmail
	}
	if body.PhoneNumber != nil {
		set["phone_number"] = *body.PhoneNumber
	}
	if body.SecondaryPhoneNumber != nil {
		set["secondary_phone_number"] = *body.SecondaryPhoneNumber
	}
	if _, err := h.coll(models.CollInstitutions).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": set}); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to update institution", err.Error())
	}
	updated, _ := h.getInstitution(ctx, oid)
	name := inst.Name
	if updated != nil {
		name = updated.Name
	}
	return response.OK(c, "Institution updated", map[string]any{"id": oid.Hex(), "name": name})
}

func (h *handler) deleteInstitution(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Institution not found", "Not found")
	}
	user := middleware.UserFrom(c)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	inst, err := h.getInstitution(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch institution", err.Error())
	}
	if inst == nil {
		return response.Err(c, http.StatusNotFound, "Institution not found", "Not found")
	}
	if user.Type == models.RoleVendor && (inst.CreatedByVendorID == nil || inst.CreatedByVendorID.Hex() != user.ID) {
		return response.Err(c, http.StatusForbidden, "Not authorized", "Forbidden")
	}

	cur, err := h.coll(models.CollBatches).Find(ctx, bson.M{"institution_id": oid})
	if err == nil {
		var batches []models.Batch
		_ = cur.All(ctx, &batches)
		for _, b := range batches {
			_, _ = h.coll(models.CollStudents).DeleteMany(ctx, bson.M{"batch_id": b.ID})
			_, _ = h.coll(models.CollTeachers).DeleteMany(ctx, bson.M{"batch_id": b.ID})
			_, _ = h.coll(models.CollCourseEnrollments).DeleteMany(ctx, bson.M{"batch_id": b.ID})
			_, _ = h.coll(models.CollBatches).DeleteOne(ctx, bson.M{"_id": b.ID})
		}
	}
	_, _ = h.coll(models.CollTeachers).DeleteMany(ctx, bson.M{"institute_id": oid})
	_, _ = h.coll(models.CollStudents).DeleteMany(ctx, bson.M{"institute_id": oid})
	if _, err := h.coll(models.CollInstitutions).DeleteOne(ctx, bson.M{"_id": oid}); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to delete institution", err.Error())
	}
	return response.OK(c, "Institution deleted", nil)
}

func (h *handler) institutionDashboard(c echo.Context) error {
	oid, ok := dbutil.ParseID(c.Param("institution_id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Institution not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	inst, err := h.getInstitution(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch institution", err.Error())
	}
	if inst == nil {
		return response.Err(c, http.StatusNotFound, "Institution not found", "Not found")
	}

	batches := h.countColl(ctx, models.CollBatches, bson.M{"institution_id": oid})
	students := h.countColl(ctx, models.CollStudents, bson.M{"institute_id": oid})
	teachers := h.countColl(ctx, models.CollTeachers, bson.M{"institute_id": oid})

	var batchIDs []primitive.ObjectID
	if cur, err := h.coll(models.CollBatches).Find(ctx, bson.M{"institution_id": oid}); err == nil {
		var bs []models.Batch
		_ = cur.All(ctx, &bs)
		for _, b := range bs {
			batchIDs = append(batchIDs, b.ID)
		}
	}
	var enrollments, assessments int64
	if len(batchIDs) > 0 {
		enrollments = h.countColl(ctx, models.CollCourseEnrollments, bson.M{"batch_id": bson.M{"$in": batchIDs}})
		assessments = h.countColl(ctx, models.CollAssessments, bson.M{"batch_id": bson.M{"$in": batchIDs}})
	}

	return response.OK(c, "Dashboard data", map[string]any{
		"institution": map[string]any{"id": inst.ID.Hex(), "name": inst.Name},
		"batches":     batches, "students": students, "teachers": teachers,
		"enrollments": enrollments, "assessments": assessments,
	})
}
