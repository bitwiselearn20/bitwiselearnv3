// Package identity ports the student/teacher/institution/vendor/batch/admin
// CRUD routers (apps/python-server/routers/{student,teacher,institution,
// vendor,batch,admin}.py) to Go.
package identity

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bitwiselearn/go-server/internal/db"
	"github.com/bitwiselearn/go-server/internal/jobs"
	"github.com/bitwiselearn/go-server/internal/middleware"
	"github.com/bitwiselearn/go-server/internal/models"
	"github.com/bitwiselearn/go-server/internal/services/queue"
)

const reqTimeout = 8 * time.Second

// Deps holds the dependencies every identity handler needs.
type Deps struct {
	DB        *db.Client
	Auth      *middleware.Auth
	Publisher *queue.Publisher
}

// Register mounts all six identity route groups on e.
func Register(e *echo.Echo, d Deps) {
	h := &handler{db: d.DB, pub: d.Publisher}
	registerStudentRoutes(e.Group("/api/v1/students"), d.Auth, h)
	registerTeacherRoutes(e.Group("/api/v1/teachers"), d.Auth, h)
	registerInstitutionRoutes(e.Group("/api/v1/institutions"), d.Auth, h)
	registerVendorRoutes(e.Group("/api/v1/vendors"), d.Auth, h)
	registerBatchRoutes(e.Group("/api/v1/batches"), d.Auth, h)
	registerAdminRoutes(e.Group("/api/v1/admins"), d.Auth, h)
}

type handler struct {
	db  *db.Client
	pub *queue.Publisher
}

func (h *handler) coll(name string) *mongo.Collection { return h.db.Coll(name) }

// genPassword mirrors Python's secrets.token_urlsafe(10).
func genPassword() string {
	buf := make([]byte, 10)
	_, _ = rand.Read(buf)
	return base64.RawURLEncoding.EncodeToString(buf)
}

func (h *handler) publishWelcomeEmail(ctx context.Context, to, name, password, role string) {
	if h.pub == nil {
		return
	}
	_ = h.pub.Publish(ctx, jobs.EmailQueue, jobs.EmailJob{
		Kind: jobs.EmailKindWelcome, To: to, Name: name, Password: password, Role: role,
	})
}

func iso(t time.Time) string { return t.Format(time.RFC3339) }

func isoMillis(t time.Time) string { return t.Format("2006-01-02T15:04:05.000Z07:00") }

// ---- shared lookups ----

func (h *handler) getInstitution(ctx context.Context, id primitive.ObjectID) (*models.Institution, error) {
	var inst models.Institution
	err := h.coll(models.CollInstitutions).FindOne(ctx, bson.M{"_id": id}).Decode(&inst)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &inst, err
}

func (h *handler) getBatch(ctx context.Context, id primitive.ObjectID) (*models.Batch, error) {
	var b models.Batch
	err := h.coll(models.CollBatches).FindOne(ctx, bson.M{"_id": id}).Decode(&b)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &b, err
}

func (h *handler) getTeacher(ctx context.Context, id primitive.ObjectID) (*models.Teacher, error) {
	var t models.Teacher
	err := h.coll(models.CollTeachers).FindOne(ctx, bson.M{"_id": id}).Decode(&t)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &t, err
}

func (h *handler) getStudent(ctx context.Context, id primitive.ObjectID) (*models.Student, error) {
	var s models.Student
	err := h.coll(models.CollStudents).FindOne(ctx, bson.M{"_id": id}).Decode(&s)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &s, err
}

func (h *handler) countColl(ctx context.Context, coll string, filter bson.M) int64 {
	if filter == nil {
		filter = bson.M{}
	}
	n, _ := h.coll(coll).CountDocuments(ctx, filter)
	return n
}
