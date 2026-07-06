package identity

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bitwiselearn/go-server/internal/auth"
	"github.com/bitwiselearn/go-server/internal/dbutil"
	"github.com/bitwiselearn/go-server/internal/middleware"
	"github.com/bitwiselearn/go-server/internal/models"
	"github.com/bitwiselearn/go-server/internal/response"
)

func registerVendorRoutes(g *echo.Group, am *middleware.Auth, h *handler) {
	authed := am.Required
	createRoles := am.RequireRoles(models.RoleSuperadmin, models.RoleAdmin)
	readRoles := am.RequireRoles(models.RoleSuperadmin, models.RoleAdmin, models.RoleInstitution)
	deleteRoles := am.RequireRoles(models.RoleSuperadmin, models.RoleAdmin, models.RoleInstitution, models.RoleVendor)
	updateRoles := am.RequireRoles(models.RoleSuperadmin, models.RoleAdmin, models.RoleVendor)

	g.POST("/create-vendor", h.createVendor, authed, createRoles)
	g.GET("/get-all-vendor", h.getAllVendors, authed, readRoles)
	g.GET("/get-vendor-by-id/:id", h.getVendorByID, authed, updateRoles)
	g.PUT("/update-vendor-by-id/:id", h.updateVendor, authed, updateRoles)
	g.DELETE("/delete-vendor-by-id/:id", h.deleteVendor, authed, deleteRoles)
	g.GET("/dashboard", h.vendorDashboard, authed)
}

func (h *handler) getVendor(ctx context.Context, id primitive.ObjectID) (*models.Vendor, error) {
	var v models.Vendor
	err := h.coll(models.CollVendors).FindOne(ctx, bson.M{"_id": id}).Decode(&v)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &v, err
}

func (h *handler) createVendor(c echo.Context) error {
	var body createVendorRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	existing := h.coll(models.CollVendors).FindOne(ctx, bson.M{"email": body.Email})
	if existing.Err() == nil {
		return response.Err(c, http.StatusBadRequest, "Email already exists", "Duplicate email")
	}

	rawPassword := genPassword()
	hashed, err := auth.HashPassword(rawPassword)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to create vendor", err.Error())
	}

	now := time.Now().UTC()
	vendor := models.Vendor{
		ID: primitive.NewObjectID(), Name: body.Name, Email: body.Email, SecondaryEmail: body.SecondaryEmail,
		Tagline: body.Tagline, PhoneNumber: body.PhoneNumber, SecondaryPhoneNumber: body.SecondaryPhoneNumber,
		WebsiteLink: body.WebsiteLink, LoginPassword: hashed, CreatedAt: now, UpdatedAt: now,
	}
	if _, err := h.coll(models.CollVendors).InsertOne(ctx, vendor); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to create vendor", err.Error())
	}

	h.publishWelcomeEmail(ctx, body.Email, body.Name, rawPassword, "Vendor")

	return response.JSON(c, http.StatusCreated, "Vendor created", map[string]any{
		"id": vendor.ID.Hex(), "name": vendor.Name, "email": vendor.Email,
	}, nil)
}

func (h *handler) getAllVendors(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	cur, err := h.coll(models.CollVendors).Find(ctx, bson.M{})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch vendors", err.Error())
	}
	var vendors []models.Vendor
	if err := cur.All(ctx, &vendors); err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch vendors", err.Error())
	}
	data := make([]map[string]any, 0, len(vendors))
	for _, v := range vendors {
		var createdAt any
		if !v.CreatedAt.IsZero() {
			createdAt = isoMillis(v.CreatedAt)
		}
		data = append(data, map[string]any{
			"id": v.ID.Hex(), "name": v.Name, "email": v.Email, "phone_number": v.PhoneNumber,
			"tagline": v.Tagline, "website_link": v.WebsiteLink, "created_at": createdAt,
		})
	}
	return response.OK(c, "Vendors fetched", data)
}

func (h *handler) getVendorByID(c echo.Context) error {
	id := c.Param("id")
	user := middleware.UserFrom(c)
	if user.Type == models.RoleVendor && user.ID != id {
		return response.Err(c, http.StatusForbidden, "Not authorized", "Forbidden")
	}
	oid, ok := dbutil.ParseID(id)
	if !ok {
		return response.Err(c, http.StatusNotFound, "Vendor not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	vendor, err := h.getVendor(ctx, oid)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch vendor", err.Error())
	}
	if vendor == nil {
		return response.Err(c, http.StatusNotFound, "Vendor not found", "Not found")
	}
	var createdAt any
	if !vendor.CreatedAt.IsZero() {
		createdAt = isoMillis(vendor.CreatedAt)
	}
	return response.OK(c, "Vendor fetched", map[string]any{
		"id": vendor.ID.Hex(), "name": vendor.Name, "email": vendor.Email,
		"secondary_email": vendor.SecondaryEmail, "phone_number": vendor.PhoneNumber,
		"secondary_phone_number": vendor.SecondaryPhoneNumber, "tagline": vendor.Tagline,
		"website_link": vendor.WebsiteLink, "created_at": createdAt,
	})
}

func (h *handler) updateVendor(c echo.Context) error {
	id := c.Param("id")
	user := middleware.UserFrom(c)
	if user.Type == models.RoleVendor && user.ID != id {
		return response.Err(c, http.StatusForbidden, "Not authorized", "Forbidden")
	}
	oid, ok := dbutil.ParseID(id)
	if !ok {
		return response.Err(c, http.StatusNotFound, "Vendor not found", "Not found")
	}
	var body updateVendorRequest
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
	if body.SecondaryEmail != nil {
		set["secondary_email"] = *body.SecondaryEmail
	}
	if body.Tagline != nil {
		set["tagline"] = *body.Tagline
	}
	if body.PhoneNumber != nil {
		set["phone_number"] = *body.PhoneNumber
	}
	if body.SecondaryPhoneNumber != nil {
		set["secondary_phone_number"] = *body.SecondaryPhoneNumber
	}
	if body.WebsiteLink != nil {
		set["website_link"] = *body.WebsiteLink
	}
	res, err := h.coll(models.CollVendors).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": set})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to update vendor", err.Error())
	}
	if res.MatchedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Vendor not found", "Not found")
	}
	vendor, _ := h.getVendor(ctx, oid)
	name := ""
	if vendor != nil {
		name = vendor.Name
	}
	return response.OK(c, "Vendor updated", map[string]any{"id": oid.Hex(), "name": name})
}

func (h *handler) deleteVendor(c echo.Context) error {
	user := middleware.UserFrom(c)
	if user.Type == models.RoleTeacher || user.Type == models.RoleStudent {
		return response.Err(c, http.StatusForbidden, "Not authorized", "Forbidden")
	}
	oid, ok := dbutil.ParseID(c.Param("id"))
	if !ok {
		return response.Err(c, http.StatusNotFound, "Vendor not found", "Not found")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	res, err := h.coll(models.CollVendors).DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to delete vendor", err.Error())
	}
	if res.DeletedCount == 0 {
		return response.Err(c, http.StatusNotFound, "Vendor not found", "Not found")
	}
	return response.OK(c, "Vendor deleted", nil)
}

func (h *handler) vendorDashboard(c echo.Context) error {
	user := middleware.UserFrom(c)
	if user.Type != models.RoleVendor {
		return response.Err(c, http.StatusForbidden, "Not authorized", "Forbidden")
	}
	oid, ok := dbutil.ParseID(user.ID)
	if !ok {
		return response.Err(c, http.StatusForbidden, "Not authorized", "Forbidden")
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	cur, err := h.coll(models.CollInstitutions).Find(ctx, bson.M{"created_by_vendor_id": oid})
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to fetch dashboard", err.Error())
	}
	var institutions []models.Institution
	_ = cur.All(ctx, &institutions)

	var totalStudents, totalTeachers, totalBatches int64
	for _, inst := range institutions {
		totalBatches += h.countColl(ctx, models.CollBatches, bson.M{"institution_id": inst.ID})
		totalStudents += h.countColl(ctx, models.CollStudents, bson.M{"institute_id": inst.ID})
		totalTeachers += h.countColl(ctx, models.CollTeachers, bson.M{"institute_id": inst.ID})
	}

	return response.OK(c, "Vendor dashboard", map[string]any{
		"institutions": len(institutions), "batches": totalBatches,
		"students": totalStudents, "teachers": totalTeachers,
	})
}
