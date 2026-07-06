// Package auth ports apps/python-server/routers/auth.py to Go.
package auth

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	appauth "github.com/bitwiselearn/go-server/internal/auth"
	"github.com/bitwiselearn/go-server/internal/db"
	"github.com/bitwiselearn/go-server/internal/dbutil"
	"github.com/bitwiselearn/go-server/internal/jobs"
	"github.com/bitwiselearn/go-server/internal/models"
	"github.com/bitwiselearn/go-server/internal/response"
	"github.com/bitwiselearn/go-server/internal/services/queue"
)

const reqTimeout = 8 * time.Second

// Deps holds the dependencies the auth handlers need.
type Deps struct {
	DB        *db.Client
	JWT       *appauth.Manager
	OTP       appauth.OTPStore
	Blocklist appauth.TokenBlocklist
	Publisher *queue.Publisher
}

// Register mounts every auth route on g (expected prefix "/api/v1/auth").
func Register(g *echo.Group, d Deps) {
	h := &handler{db: d.DB, jwt: d.JWT, otp: d.OTP, blocklist: d.Blocklist, pub: d.Publisher}

	g.POST("/admin/login", h.adminLogin)
	g.POST("/institution/login", h.institutionLogin)
	g.POST("/vendor/login", h.vendorLogin)
	g.POST("/teacher/login", h.teacherLogin)
	g.POST("/student/login", h.studentLogin)
	g.POST("/refresh", h.refresh)
	g.POST("/forgot-password", h.forgotPassword)
	g.POST("/verify-forgot-password", h.verifyForgotPassword)
	g.POST("/reset-password", h.resetPassword)
}

type handler struct {
	db        *db.Client
	jwt       *appauth.Manager
	otp       appauth.OTPStore
	blocklist appauth.TokenBlocklist
	pub       *queue.Publisher
}

func (h *handler) coll(name string) *mongo.Collection { return h.db.Coll(name) }

func setAuthCookies(c echo.Context, access, refresh string) {
	c.SetCookie(&http.Cookie{
		Name: "token", Value: access, HttpOnly: true, Secure: true,
		SameSite: http.SameSiteNoneMode, Path: "/", MaxAge: 86400,
	})
	c.SetCookie(&http.Cookie{
		Name: "refreshToken", Value: refresh, HttpOnly: true, Secure: true,
		SameSite: http.SameSiteNoneMode, Path: "/", MaxAge: 86400 * 20,
	})
}

func invalidCreds(c echo.Context) error {
	return response.Err(c, http.StatusUnauthorized, "Invalid email or password", "Invalid credentials")
}

// ---- role logins ----

func (h *handler) adminLogin(c echo.Context) error {
	var body loginRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	var user models.User
	err := h.coll(models.CollUsers).FindOne(ctx, bson.M{"email": body.Email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return invalidCreds(c)
	}
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Login failed", err.Error())
	}
	if !appauth.VerifyPassword(body.Password, user.Password) {
		return invalidCreds(c)
	}

	access, _ := h.jwt.GenerateAccessToken(user.ID.Hex(), user.Role)
	refresh, _ := h.jwt.GenerateRefreshToken(user.ID.Hex(), user.Role)
	if err := response.OK(c, "Login successful", map[string]any{
		"id": user.ID.Hex(), "name": user.Name, "email": user.Email, "role": user.Role,
		"tokens": map[string]any{"access_token": access, "refresh_token": refresh},
	}); err != nil {
		return err
	}
	setAuthCookies(c, access, refresh)
	return nil
}

func (h *handler) institutionLogin(c echo.Context) error {
	var body loginRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	var inst models.Institution
	err := h.coll(models.CollInstitutions).FindOne(ctx, bson.M{"email": body.Email}).Decode(&inst)
	if err == mongo.ErrNoDocuments {
		return invalidCreds(c)
	}
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Login failed", err.Error())
	}
	if !appauth.VerifyPassword(body.Password, inst.LoginPassword) {
		return invalidCreds(c)
	}

	access, _ := h.jwt.GenerateAccessToken(inst.ID.Hex(), models.RoleInstitution)
	refresh, _ := h.jwt.GenerateRefreshToken(inst.ID.Hex(), models.RoleInstitution)
	if err := response.OK(c, "Login successful", map[string]any{
		"id": inst.ID.Hex(), "name": inst.Name, "email": inst.Email, "type": models.RoleInstitution,
		"tokens": map[string]any{"access_token": access, "refresh_token": refresh},
	}); err != nil {
		return err
	}
	setAuthCookies(c, access, refresh)
	return nil
}

func (h *handler) vendorLogin(c echo.Context) error {
	var body loginRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	var vendor models.Vendor
	err := h.coll(models.CollVendors).FindOne(ctx, bson.M{"email": body.Email}).Decode(&vendor)
	if err == mongo.ErrNoDocuments {
		return invalidCreds(c)
	}
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Login failed", err.Error())
	}
	if !appauth.VerifyPassword(body.Password, vendor.LoginPassword) {
		return invalidCreds(c)
	}

	access, _ := h.jwt.GenerateAccessToken(vendor.ID.Hex(), models.RoleVendor)
	refresh, _ := h.jwt.GenerateRefreshToken(vendor.ID.Hex(), models.RoleVendor)
	if err := response.OK(c, "Login successful", map[string]any{
		"id": vendor.ID.Hex(), "name": vendor.Name, "email": vendor.Email, "type": models.RoleVendor,
		"tokens": map[string]any{"access_token": access, "refresh_token": refresh},
	}); err != nil {
		return err
	}
	setAuthCookies(c, access, refresh)
	return nil
}

func (h *handler) teacherLogin(c echo.Context) error {
	var body loginRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	var teacher models.Teacher
	err := h.coll(models.CollTeachers).FindOne(ctx, bson.M{"email": body.Email}).Decode(&teacher)
	if err == mongo.ErrNoDocuments {
		return invalidCreds(c)
	}
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Login failed", err.Error())
	}
	if !appauth.VerifyPassword(body.Password, teacher.LoginPassword) {
		return invalidCreds(c)
	}

	access, _ := h.jwt.GenerateAccessToken(teacher.ID.Hex(), models.RoleTeacher)
	refresh, _ := h.jwt.GenerateRefreshToken(teacher.ID.Hex(), models.RoleTeacher)
	if err := response.OK(c, "Login successful", map[string]any{
		"id": teacher.ID.Hex(), "name": teacher.Name, "email": teacher.Email, "type": models.RoleTeacher,
		"tokens": map[string]any{"access_token": access, "refresh_token": refresh},
	}); err != nil {
		return err
	}
	setAuthCookies(c, access, refresh)
	return nil
}

func (h *handler) studentLogin(c echo.Context) error {
	var body loginRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	var student models.Student
	err := h.coll(models.CollStudents).FindOne(ctx, bson.M{"email": body.Email}).Decode(&student)
	if err == mongo.ErrNoDocuments {
		return invalidCreds(c)
	}
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Login failed", err.Error())
	}
	if !appauth.VerifyPassword(body.Password, student.LoginPassword) {
		return invalidCreds(c)
	}

	var batch *models.Batch
	if b, err := h.getBatch(ctx, student.BatchID); err == nil {
		batch = b
	}
	var institute *models.Institution
	if i, err := h.getInstitutionByID(ctx, student.InstituteID); err == nil {
		institute = i
	}

	access, _ := h.jwt.GenerateAccessToken(student.ID.Hex(), models.RoleStudent)
	refresh, _ := h.jwt.GenerateRefreshToken(student.ID.Hex(), models.RoleStudent)

	var batchData, instData any
	if batch != nil {
		batchData = map[string]any{
			"id": batch.ID.Hex(), "batchname": batch.BatchName, "branch": batch.Branch,
			"batchEndYear": batch.BatchEndYear,
		}
	}
	if institute != nil {
		instData = map[string]any{"id": institute.ID.Hex(), "name": institute.Name}
	}

	if err := response.OK(c, "Login successful", map[string]any{
		"id": student.ID.Hex(), "name": student.Name, "email": student.Email, "type": models.RoleStudent,
		"roll_number": student.RollNumber, "batch": batchData, "institution": instData,
		"tokens": map[string]any{"access_token": access, "refresh_token": refresh},
	}); err != nil {
		return err
	}
	setAuthCookies(c, access, refresh)
	return nil
}

func (h *handler) getBatch(ctx context.Context, id primitive.ObjectID) (*models.Batch, error) {
	var b models.Batch
	err := h.coll(models.CollBatches).FindOne(ctx, bson.M{"_id": id}).Decode(&b)
	return &b, err
}

func (h *handler) getInstitutionByID(ctx context.Context, id primitive.ObjectID) (*models.Institution, error) {
	var inst models.Institution
	err := h.coll(models.CollInstitutions).FindOne(ctx, bson.M{"_id": id}).Decode(&inst)
	return &inst, err
}

// ---- refresh ----

func (h *handler) refresh(c echo.Context) error {
	cookie, err := c.Cookie("refreshToken")
	if err != nil || cookie.Value == "" {
		return response.Err(c, http.StatusUnauthorized, "No refresh token provided", "No refresh token")
	}
	claims := h.jwt.VerifyRefreshToken(cookie.Value)
	if claims == nil {
		return response.Err(c, http.StatusUnauthorized, "Invalid refresh token", "Invalid refresh token")
	}
	userID, _ := claims["id"].(string)
	userType, _ := claims["type"].(string)

	access, _ := h.jwt.GenerateAccessToken(userID, userType)
	refresh, _ := h.jwt.GenerateRefreshToken(userID, userType)
	if err := response.OK(c, "Token refreshed", map[string]any{"id": userID, "type": userType}); err != nil {
		return err
	}
	setAuthCookies(c, access, refresh)
	return nil
}

// ---- forgot password / OTP / reset ----

// foundAccount is the result of searching every identity collection by email,
// mirroring the linear User -> Institution -> Vendor -> Teacher -> Student
// search in the legacy forgot_password/verify_forgot_password handlers.
type foundAccount struct {
	Type string
	ID   string
}

func (h *handler) findAccountByEmail(ctx context.Context, email string) *foundAccount {
	var user models.User
	if err := h.coll(models.CollUsers).FindOne(ctx, bson.M{"email": email}).Decode(&user); err == nil {
		return &foundAccount{Type: user.Role, ID: user.ID.Hex()}
	}
	var inst models.Institution
	if err := h.coll(models.CollInstitutions).FindOne(ctx, bson.M{"email": email}).Decode(&inst); err == nil {
		return &foundAccount{Type: models.RoleInstitution, ID: inst.ID.Hex()}
	}
	var vendor models.Vendor
	if err := h.coll(models.CollVendors).FindOne(ctx, bson.M{"email": email}).Decode(&vendor); err == nil {
		return &foundAccount{Type: models.RoleVendor, ID: vendor.ID.Hex()}
	}
	var teacher models.Teacher
	if err := h.coll(models.CollTeachers).FindOne(ctx, bson.M{"email": email}).Decode(&teacher); err == nil {
		return &foundAccount{Type: models.RoleTeacher, ID: teacher.ID.Hex()}
	}
	var student models.Student
	if err := h.coll(models.CollStudents).FindOne(ctx, bson.M{"email": email}).Decode(&student); err == nil {
		return &foundAccount{Type: models.RoleStudent, ID: student.ID.Hex()}
	}
	return nil
}

func (h *handler) forgotPassword(c echo.Context) error {
	var body forgotPasswordRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	email := strings.ToLower(body.Email)
	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()

	if h.findAccountByEmail(ctx, email) == nil {
		return response.Err(c, http.StatusNotFound, "Email not found", "Email not registered")
	}

	otp := h.otp.Generate(email)
	if h.pub != nil {
		_ = h.pub.Publish(ctx, jobs.EmailQueue, jobs.EmailJob{Kind: jobs.EmailKindOTP, To: email, OTP: otp})
	}
	return response.OK(c, "OTP sent to email", map[string]any{"email": email})
}

func (h *handler) verifyForgotPassword(c echo.Context) error {
	var body verifyOtpRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}
	email := strings.ToLower(body.Email)
	if !h.otp.Verify(email, body.OTP) {
		return response.Err(c, http.StatusBadRequest, "Invalid or expired OTP", "Invalid OTP")
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	account := h.findAccountByEmail(ctx, email)
	if account == nil {
		return response.Err(c, http.StatusNotFound, "User not found", "User not found")
	}

	resetTok, _ := h.jwt.GenerateResetToken(email, account.Type, account.ID)
	if err := response.OK(c, "OTP verified", map[string]any{"verified": true, "reset_token": resetTok}); err != nil {
		return err
	}
	c.SetCookie(&http.Cookie{
		Name: "reset_token", Value: resetTok, HttpOnly: true, Secure: true,
		SameSite: http.SameSiteNoneMode, Path: "/", MaxAge: 600,
	})
	return nil
}

func (h *handler) resetPassword(c echo.Context) error {
	cookie, err := c.Cookie("reset_token")
	if err != nil || cookie.Value == "" {
		return response.Err(c, http.StatusUnauthorized, "No reset token", "No reset token")
	}
	var body resetPasswordRequest
	if err := c.Bind(&body); err != nil {
		return response.Err(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}

	claims := h.jwt.VerifyResetToken(cookie.Value, h.blocklist)
	if claims == nil {
		return response.Err(c, http.StatusUnauthorized, "Invalid or expired reset token", "Invalid reset token")
	}
	userType, _ := claims["type"].(string)
	userID, _ := claims["id"].(string)

	hashed, err := appauth.HashPassword(body.NewPassword)
	if err != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to reset password", err.Error())
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), reqTimeout)
	defer cancel()
	oid, ok := dbutil.ParseID(userID)
	if !ok {
		return response.Err(c, http.StatusUnauthorized, "Invalid reset token", "Invalid reset token")
	}

	var updateErr error
	switch userType {
	case models.RoleSuperadmin, models.RoleAdmin:
		_, updateErr = h.coll(models.CollUsers).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"password": hashed}})
	case models.RoleInstitution:
		_, updateErr = h.coll(models.CollInstitutions).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"login_password": hashed}})
	case models.RoleVendor:
		_, updateErr = h.coll(models.CollVendors).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"login_password": hashed}})
	case models.RoleTeacher:
		_, updateErr = h.coll(models.CollTeachers).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"login_password": hashed}})
	case models.RoleStudent:
		_, updateErr = h.coll(models.CollStudents).UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"login_password": hashed}})
	}
	if updateErr != nil {
		return response.Err(c, http.StatusInternalServerError, "Failed to reset password", updateErr.Error())
	}

	h.blocklist.Invalidate(cookie.Value)

	access, _ := h.jwt.GenerateAccessToken(userID, userType)
	refresh, _ := h.jwt.GenerateRefreshToken(userID, userType)
	if err := response.OK(c, "Password reset successful", map[string]any{"id": userID, "type": userType}); err != nil {
		return err
	}
	setAuthCookies(c, access, refresh)
	c.SetCookie(&http.Cookie{Name: "reset_token", Value: "", Path: "/", MaxAge: -1})
	return nil
}
