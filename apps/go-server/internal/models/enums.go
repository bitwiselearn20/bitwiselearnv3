package models

// Ported verbatim from the legacy enums.py. String constants keep the exact
// stored values so existing documents validate unchanged.

// UserType roles.
const (
	RoleSuperadmin  = "SUPERADMIN"
	RoleAdmin       = "ADMIN"
	RoleTeacher     = "TEACHER"
	RoleInstitution = "INSTITUTION"
	RoleVendor      = "VENDOR"
	RoleStudent     = "STUDENT"
)

// ListingStatus.
const (
	StatusListed    = "LISTED"
	StatusNotListed = "NOT_LISTED"
)

// CourseStatus.
const (
	CoursePublished    = "PUBLISHED"
	CourseNotPublished = "NOT_PUBLISHED"
)

// CourseLevel.
const (
	CourseLevelBasic        = "BASIC"
	CourseLevelIntermediate = "INTERMEDIATE"
	CourseLevelAdvance      = "ADVANCE"
)

// ProblemLevel.
const (
	ProblemEasy   = "EASY"
	ProblemMedium = "MEDIUM"
	ProblemHard   = "HARD"
)

// TestcaseType.
const (
	TestcaseExample = "EXAMPLE"
	TestcaseHidden  = "HIDDEN"
)

// Languages.
const (
	LangJavaScript = "JAVASCRIPT"
	LangJava       = "JAVA"
	LangPython     = "PYTHON"
	LangC          = "C"
	LangCPP        = "CPP"
)

// ProblemStatus.
const (
	ProblemSuccess = "SUCCESS"
	ProblemFailed  = "FAILED"
)

// AssignmentType.
const (
	AssignmentMCQ = "MCQ"
	AssignmentSCQ = "SCQ"
)

// AssessmentType.
const (
	AssessmentCode   = "CODE"
	AssessmentNoCode = "NO_CODE"
)

// AssessmentStatus.
const (
	AssessmentUpcoming = "UPCOMING"
	AssessmentLive     = "LIVE"
	AssessmentEnded    = "ENDED"
)

// ReportStatus.
const (
	ReportNotRequested = "NOT_REQUESTED"
	ReportProcessing   = "PROCESSING"
	ReportProcessed    = "PROCESSED"
)
