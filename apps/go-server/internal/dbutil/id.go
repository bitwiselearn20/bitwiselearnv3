// Package dbutil holds small MongoDB helpers shared across handler packages.
package dbutil

import "go.mongodb.org/mongo-driver/bson/primitive"

// ParseID parses a hex string into an ObjectID, matching the legacy
// PydanticObjectId(id) coercion used throughout the FastAPI routers.
func ParseID(hex string) (primitive.ObjectID, bool) {
	oid, err := primitive.ObjectIDFromHex(hex)
	return oid, err == nil
}
