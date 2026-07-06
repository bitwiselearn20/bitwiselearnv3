package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Collection names match the legacy Beanie Settings.name values.
const (
	CollUsers        = "users"
	CollInstitutions = "institutions"
	CollVendors      = "vendors"
	CollTeachers     = "teachers"
	CollStudents     = "students"
)

// User is a superadmin/admin account (collection "users").
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`
	Role      string             `bson:"role" json:"role"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}

// Institution account (collection "institutions").
type Institution struct {
	ID                   primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Name                 string              `bson:"name" json:"name"`
	Address              string              `bson:"address" json:"address"`
	PinCode              string              `bson:"pin_code" json:"pinCode"`
	Tagline              string              `bson:"tagline" json:"tagline"`
	WebsiteLink          string              `bson:"website_link" json:"websiteLink"`
	LoginPassword        string              `bson:"login_password" json:"-"`
	Email                string              `bson:"email" json:"email"`
	SecondaryEmail       *string             `bson:"secondary_email,omitempty" json:"secondaryEmail,omitempty"`
	PhoneNumber          string              `bson:"phone_number" json:"phoneNumber"`
	SecondaryPhoneNumber *string             `bson:"secondary_phone_number,omitempty" json:"secondaryPhoneNumber,omitempty"`
	CreatedBy            primitive.ObjectID  `bson:"created_by" json:"createdBy"`
	CreatedByVendorID    *primitive.ObjectID `bson:"created_by_vendor_id,omitempty" json:"createdByVendorId,omitempty"`
	CreatedAt            time.Time           `bson:"created_at" json:"createdAt"`
	UpdatedAt            time.Time           `bson:"updated_at" json:"updatedAt"`
}

// Vendor account (collection "vendors").
type Vendor struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name                 string             `bson:"name" json:"name"`
	Email                string             `bson:"email" json:"email"`
	SecondaryEmail       *string            `bson:"secondary_email,omitempty" json:"secondaryEmail,omitempty"`
	Tagline              string             `bson:"tagline" json:"tagline"`
	PhoneNumber          string             `bson:"phone_number" json:"phoneNumber"`
	SecondaryPhoneNumber *string            `bson:"secondary_phone_number,omitempty" json:"secondaryPhoneNumber,omitempty"`
	WebsiteLink          string             `bson:"website_link" json:"websiteLink"`
	LoginPassword        string             `bson:"login_password" json:"-"`
	CreatedAt            time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt            time.Time          `bson:"updated_at" json:"updatedAt"`
}

// Teacher account (collection "teachers").
type Teacher struct {
	ID            primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Name          string              `bson:"name" json:"name"`
	Email         string              `bson:"email" json:"email"`
	PhoneNumber   string              `bson:"phone_number" json:"phoneNumber"`
	LoginPassword string              `bson:"login_password" json:"-"`
	VendorID      *primitive.ObjectID `bson:"vendor_id,omitempty" json:"vendorId,omitempty"`
	InstituteID   primitive.ObjectID  `bson:"institute_id" json:"instituteId"`
	BatchID       primitive.ObjectID  `bson:"batch_id" json:"batchId"`
	CreatedAt     time.Time           `bson:"created_at" json:"createdAt"`
	UpdatedAt     time.Time           `bson:"updated_at" json:"updatedAt"`
}

// Student account (collection "students").
type Student struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name"`
	RollNumber    string             `bson:"roll_number" json:"rollNumber"`
	Email         string             `bson:"email" json:"email"`
	LoginPassword string             `bson:"login_password" json:"-"`
	CloudPlatform string             `bson:"cloud_platform" json:"cloudPlatform"`
	CloudName     *string            `bson:"cloudname,omitempty" json:"cloudname,omitempty"`
	CloudPass     *string            `bson:"cloudpass,omitempty" json:"cloudpass,omitempty"`
	CloudURL      *string            `bson:"cloudurl,omitempty" json:"cloudurl,omitempty"`
	BatchID       primitive.ObjectID `bson:"batch_id" json:"batchId"`
	InstituteID   primitive.ObjectID `bson:"institute_id" json:"instituteId"`
	CreatedAt     time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updatedAt"`
}
