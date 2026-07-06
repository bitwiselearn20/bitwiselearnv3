// Package db owns the MongoDB connection (with pooling) and startup index sync.
package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client wraps the Mongo connection and the resolved database handle.
type Client struct {
	mongo *mongo.Client
	DB    *mongo.Database
}

// Connect dials MongoDB with a bounded connection pool and verifies reachability
// with a ping (mirrors the Python connect_to_mongo behaviour).
func Connect(ctx context.Context, uri, dbName string, serverSelectionTimeoutMS int) (*Client, error) {
	if uri == "" {
		return nil, fmt.Errorf("DATABASE_URL is not configured")
	}

	opts := options.Client().
		ApplyURI(uri).
		SetServerSelectionTimeout(time.Duration(serverSelectionTimeoutMS) * time.Millisecond).
		// Pooling: keep replicas lean while supporting bursts under autoscale.
		SetMaxPoolSize(100).
		SetMinPoolSize(5).
		SetMaxConnIdleTime(60 * time.Second)

	m, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, time.Duration(serverSelectionTimeoutMS)*time.Millisecond)
	defer cancel()
	if err := m.Ping(pingCtx, nil); err != nil {
		_ = m.Disconnect(ctx)
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return &Client{mongo: m, DB: m.Database(dbName)}, nil
}

// Ping checks DB connectivity (used by the /ready probe).
func (c *Client) Ping(ctx context.Context) error {
	return c.mongo.Ping(ctx, nil)
}

// Disconnect closes the connection pool.
func (c *Client) Disconnect(ctx context.Context) error {
	return c.mongo.Disconnect(ctx)
}

// Coll returns a collection handle.
func (c *Client) Coll(name string) *mongo.Collection {
	return c.DB.Collection(name)
}

// indexSpec declares a single-field index to ensure at startup.
type indexSpec struct {
	collection string
	field      string
}

// requiredIndexes are the foreign-key / lookup fields that the legacy code
// queries on hot paths but never indexed — adding them removes collection scans
// and lowers the database tier (and Cosmos RU cost) needed to serve 10k users.
var requiredIndexes = []indexSpec{
	// identity lookups
	{"users", "email"},
	{"institutions", "email"},
	{"vendors", "email"},
	{"teachers", "email"},
	{"students", "email"},
	{"teachers", "institute_id"},
	{"students", "institute_id"},
	{"students", "batch_id"},
	{"teachers", "batch_id"},
	{"batches", "institution_id"},
	// problem domain
	{"problem_topics", "problem_id"},
	{"problem_templates", "problem_id"},
	{"problem_test_cases", "problem_id"},
	// course domain
	{"course_sections", "course_id"},
	{"course_learning_contents", "section_id"},
	{"course_assignments", "section_id"},
	{"course_enrollments", "batch_id"},
	// assessment domain
	{"assessment_sections", "assessment_id"},
	{"assessment_questions", "section_id"},
	{"assessment_submissions", "assessment_id"},
}

// EnsureIndexes creates the required indexes if missing. Idempotent and safe to
// run on every startup; creating an index on a non-existent collection creates
// the (empty) collection, which is harmless.
func (c *Client) EnsureIndexes(ctx context.Context) error {
	for _, spec := range requiredIndexes {
		model := mongo.IndexModel{
			Keys: bson.D{{Key: spec.field, Value: 1}},
		}
		if _, err := c.Coll(spec.collection).Indexes().CreateOne(ctx, model); err != nil {
			return fmt.Errorf("index %s.%s: %w", spec.collection, spec.field, err)
		}
	}
	return nil
}
