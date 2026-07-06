// Package middleware ports middleware/auth.py: JWT extraction, user existence
// verification against the correct collection, and role guards.
package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/bitwiselearn/go-server/internal/auth"
	"github.com/bitwiselearn/go-server/internal/db"
	"github.com/bitwiselearn/go-server/internal/models"
	"github.com/bitwiselearn/go-server/internal/response"
)

// CurrentUser is the authenticated principal attached to the request context.
type CurrentUser struct {
	ID   string
	Type string
}

const currentUserKey = "currentUser"

// Auth builds the authentication middleware bound to the JWT manager and DB.
type Auth struct {
	jwt *auth.Manager
	db  *db.Client
}

// NewAuth constructs the auth middleware.
func NewAuth(jwtMgr *auth.Manager, dbc *db.Client) *Auth {
	return &Auth{jwt: jwtMgr, db: dbc}
}

// extractToken reads the JWT from the "token" cookie, then the Bearer header.
func extractToken(c echo.Context) string {
	if cookie, err := c.Cookie("token"); err == nil && cookie.Value != "" {
		return cookie.Value
	}
	header := c.Request().Header.Get("Authorization")
	if strings.HasPrefix(header, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
	}
	return ""
}

// collectionForRole maps a role to the collection where that principal lives.
func collectionForRole(role string) string {
	switch role {
	case models.RoleSuperadmin, models.RoleAdmin:
		return models.CollUsers
	case models.RoleInstitution:
		return models.CollInstitutions
	case models.RoleVendor:
		return models.CollVendors
	case models.RoleTeacher:
		return models.CollTeachers
	case models.RoleStudent:
		return models.CollStudents
	default:
		return ""
	}
}

// Required is middleware that authenticates the request (get_current_user).
func (a *Auth) Required(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := extractToken(c)
		if token == "" {
			return response.Err(c, http.StatusUnauthorized, "Not authenticated", "Not authenticated")
		}

		claims := a.jwt.VerifyAccessToken(token)
		if claims == nil {
			return response.Err(c, http.StatusUnauthorized, "Invalid or expired token", "Invalid or expired token")
		}

		userID, _ := claims["id"].(string)
		userType := strings.ToUpper(toString(claims["type"]))
		if userID == "" || userType == "" {
			return response.Err(c, http.StatusUnauthorized, "Invalid token payload", "Invalid token payload")
		}

		oid, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return response.Err(c, http.StatusUnauthorized, "Invalid user ID", "Invalid user ID")
		}

		coll := collectionForRole(userType)
		if coll == "" {
			return response.Err(c, http.StatusUnauthorized, "Invalid token payload", "Invalid token payload")
		}

		ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
		defer cancel()
		count, err := a.db.Coll(coll).CountDocuments(ctx, bson.M{"_id": oid})
		if err != nil || count == 0 {
			return response.Err(c, http.StatusUnauthorized, "User not found", "User not found")
		}

		c.Set(currentUserKey, &CurrentUser{ID: userID, Type: userType})
		return next(c)
	}
}

// RequireRoles guards a route to the given roles (require_roles).
// Apply after Required.
func (a *Auth) RequireRoles(roles ...string) echo.MiddlewareFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := UserFrom(c)
			if user == nil {
				return response.Err(c, http.StatusUnauthorized, "Not authenticated", "Not authenticated")
			}
			if _, ok := allowed[user.Type]; !ok {
				return response.Err(c, http.StatusForbidden, "Insufficient permissions", "Insufficient permissions")
			}
			return next(c)
		}
	}
}

// NotStudent guards a route to any authenticated role except STUDENT
// (ports the legacy `not_student = require_roles(SUPERADMIN, ADMIN, INSTITUTION, VENDOR, TEACHER)`).
// Apply after Required.
func (a *Auth) NotStudent() echo.MiddlewareFunc {
	return a.RequireRoles(models.RoleSuperadmin, models.RoleAdmin, models.RoleInstitution, models.RoleVendor, models.RoleTeacher)
}

// AdminOnly guards a route to SUPERADMIN/ADMIN only. Apply after Required.
func (a *Auth) AdminOnly() echo.MiddlewareFunc {
	return a.RequireRoles(models.RoleSuperadmin, models.RoleAdmin)
}

// SuperadminOnly guards a route to SUPERADMIN only. Apply after Required.
func (a *Auth) SuperadminOnly() echo.MiddlewareFunc {
	return a.RequireRoles(models.RoleSuperadmin)
}

// UserFrom returns the authenticated user previously set by Required.
func UserFrom(c echo.Context) *CurrentUser {
	if v, ok := c.Get(currentUserKey).(*CurrentUser); ok {
		return v
	}
	return nil
}

func toString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
