// Package server wires the Echo HTTP server: middleware, CORS, health/ready
// probes, and (in later phases) the migrated route groups.
package server

import (
	"context"
	"net/http"
	"regexp"
	"time"

	"github.com/labstack/echo/v4"
	emw "github.com/labstack/echo/v4/middleware"

	"github.com/bitwiselearn/go-server/internal/auth"
	"github.com/bitwiselearn/go-server/internal/config"
	"github.com/bitwiselearn/go-server/internal/db"
	"github.com/bitwiselearn/go-server/internal/handlers/assessment"
	authh "github.com/bitwiselearn/go-server/internal/handlers/auth"
	"github.com/bitwiselearn/go-server/internal/handlers/bulkupload"
	"github.com/bitwiselearn/go-server/internal/handlers/code"
	"github.com/bitwiselearn/go-server/internal/handlers/contact"
	"github.com/bitwiselearn/go-server/internal/handlers/course"
	"github.com/bitwiselearn/go-server/internal/handlers/identity"
	"github.com/bitwiselearn/go-server/internal/handlers/problem"
	"github.com/bitwiselearn/go-server/internal/handlers/report"
	appmw "github.com/bitwiselearn/go-server/internal/middleware"
	"github.com/bitwiselearn/go-server/internal/services/blob"
	"github.com/bitwiselearn/go-server/internal/services/piston"
	"github.com/bitwiselearn/go-server/internal/services/queue"
)

// localOriginRe mirrors the CORS allow_origin_regex from the FastAPI app:
// localhost / 127.0.0.1 / private LAN ranges on any port.
var localOriginRe = regexp.MustCompile(
	`^https?://(localhost|127\.0\.0\.1|10(?:\.\d{1,3}){3}|192\.168(?:\.\d{1,3}){2}|172\.(1[6-9]|2\d|3[0-1])(?:\.\d{1,3}){2})(:\d+)?$`,
)

// Deps holds the constructed dependencies the server needs.
type Deps struct {
	Config    *config.Config
	DB        *db.Client
	JWT       *auth.Manager
	Auth      *appmw.Auth
	Store     *blob.Store
	Publisher *queue.Publisher
	OTP       auth.OTPStore
	Blocklist auth.TokenBlocklist
	Piston    *piston.Client
}

// New builds and configures the Echo instance.
func New(d Deps) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Core middleware. Recover keeps a panicking handler from taking down the
	// replica; RequestID + Logger feed Application Insights.
	e.Use(emw.Recover())
	e.Use(emw.RequestID())
	e.Use(emw.Logger())
	e.Use(emw.BodyLimit("25M")) // bound upload memory

	// CORS: exact frontend origin plus the local/LAN regex.
	e.Use(emw.CORSWithConfig(emw.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			if origin == d.Config.FrontendURL {
				return true, nil
			}
			return localOriginRe.MatchString(origin), nil
		},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{"*"},
	}))

	// Liveness: process is up.
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// Readiness: dependencies reachable (used by ACA readiness probe so traffic
	// only routes to replicas that can actually serve).
	e.GET("/ready", func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), 2*time.Second)
		defer cancel()
		if err := d.DB.Ping(ctx); err != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "db_unavailable"})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "ready"})
	})

	// Root parity with the FastAPI app.
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "BitwiseLearn API v2.0"})
	})

	// Route groups for migrated modules.
	problem.Register(e.Group("/api/v1/problems"), problem.Deps{DB: d.DB, Auth: d.Auth})
	course.Register(e.Group("/api/v1/courses"), course.Deps{DB: d.DB, Auth: d.Auth, Store: d.Store})
	authh.Register(e.Group("/api/v1/auth"), authh.Deps{
		DB: d.DB, JWT: d.JWT, OTP: d.OTP, Blocklist: d.Blocklist, Publisher: d.Publisher,
	})
	identity.Register(e, identity.Deps{DB: d.DB, Auth: d.Auth, Publisher: d.Publisher})
	assessment.Register(e.Group("/api/v1/assessments"), assessment.Deps{
		DB: d.DB, Auth: d.Auth, Piston: d.Piston, Publisher: d.Publisher,
	})
	report.Register(e.Group("/api/v1/reports"), report.Deps{DB: d.DB, Auth: d.Auth, Publisher: d.Publisher})
	bulkupload.Register(e.Group("/api/v1/bulk-upload"), bulkupload.Deps{DB: d.DB, Auth: d.Auth})
	contact.Register(e.Group("/api/v1/contact"), contact.Deps{Publisher: d.Publisher})
	code.Register(e.Group("/api/v1/code"), code.Deps{DB: d.DB, Auth: d.Auth, Piston: d.Piston})

	return e
}
