package identity

// Request bodies ported from schemas/{student,teacher,institution,vendor,batch,user}.py.

type createStudentRequest struct {
	Name          string  `json:"name"`
	RollNumber    string  `json:"rollNumber"`
	Email         string  `json:"email"`
	LoginPassword *string `json:"loginPassword"`
	BatchID       string  `json:"batchId"`
	InstitutionID string  `json:"institutionId"`
	CloudPlatform string  `json:"cloudPlatform"`
}

type updateStudentRequest struct {
	Name          *string `json:"name"`
	RollNumber    *string `json:"rollNumber"`
	Email         *string `json:"email"`
	CloudPlatform *string `json:"cloudPlatform"`
	CloudName     *string `json:"cloudname"`
	CloudPass     *string `json:"cloudpass"`
	CloudURL      *string `json:"cloudurl"`
}

type createTeacherRequest struct {
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	PhoneNumber string  `json:"phoneNumber"`
	VendorID    *string `json:"vendorId"`
	InstituteID string  `json:"instituteId"`
	BatchID     string  `json:"batchId"`
}

type updateTeacherRequest struct {
	Name        *string `json:"name"`
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phoneNumber"`
}

type createInstitutionRequest struct {
	Name                 string  `json:"name"`
	Address              string  `json:"address"`
	PinCode              string  `json:"pinCode"`
	Tagline              string  `json:"tagline"`
	WebsiteLink          string  `json:"websiteLink"`
	Email                string  `json:"email"`
	SecondaryEmail       *string `json:"secondaryEmail"`
	PhoneNumber          string  `json:"phoneNumber"`
	SecondaryPhoneNumber *string `json:"secondaryPhoneNumber"`
}

type updateInstitutionRequest struct {
	Name                 *string `json:"name"`
	Address              *string `json:"address"`
	PinCode              *string `json:"pinCode"`
	Tagline              *string `json:"tagline"`
	WebsiteLink          *string `json:"websiteLink"`
	Email                *string `json:"email"`
	SecondaryEmail       *string `json:"secondaryEmail"`
	PhoneNumber          *string `json:"phoneNumber"`
	SecondaryPhoneNumber *string `json:"secondaryPhoneNumber"`
}

type createVendorRequest struct {
	Name                 string  `json:"name"`
	Email                string  `json:"email"`
	SecondaryEmail       *string `json:"secondaryEmail"`
	Tagline              string  `json:"tagline"`
	PhoneNumber          string  `json:"phoneNumber"`
	SecondaryPhoneNumber *string `json:"secondaryPhoneNumber"`
	WebsiteLink          string  `json:"websiteLink"`
}

type updateVendorRequest struct {
	Name                 *string `json:"name"`
	Email                *string `json:"email"`
	SecondaryEmail       *string `json:"secondaryEmail"`
	Tagline              *string `json:"tagline"`
	PhoneNumber          *string `json:"phoneNumber"`
	SecondaryPhoneNumber *string `json:"secondaryPhoneNumber"`
	WebsiteLink          *string `json:"websiteLink"`
}

type createBatchRequest struct {
	BatchName     string `json:"batchname"`
	Branch        string `json:"branch"`
	BatchEndYear  string `json:"batchEndYear"`
	InstitutionID string `json:"institutionId"`
}

type updateBatchRequest struct {
	BatchName    *string `json:"batchname"`
	Branch       *string `json:"branch"`
	BatchEndYear *string `json:"batchEndYear"`
}

type createAdminRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type updateAdminRequest struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}
