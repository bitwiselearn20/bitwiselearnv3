package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Collection names for the assessment domain.
const (
	CollAssessments                   = "assessments"
	CollAssessmentSections            = "assessment_sections"
	CollAssessmentQuestions           = "assessment_questions"
	CollAssessmentSubmissions         = "assessment_submissions"
	CollAssessmentQuestionSubmissions = "assessment_question_submissions"
)

// Assessment (collection "assessments"). Full CRUD lands in Phase 3
// (routers/assessment.py); this struct is needed now for dashboard/count
// endpoints in the identity routers.
type Assessment struct {
	ID                         primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Name                       string              `bson:"name" json:"name"`
	Description                string              `bson:"description" json:"description"`
	Instruction                string              `bson:"instruction" json:"instruction"`
	StartTime                  time.Time           `bson:"start_time" json:"startTime"`
	EndTime                    time.Time           `bson:"end_time" json:"endTime"`
	IndividualSectionTimeLimit *int                `bson:"individual_section_time_limit,omitempty" json:"individualSectionTimeLimit,omitempty"`
	Status                     string              `bson:"status" json:"status"`
	Report                     *string             `bson:"report,omitempty" json:"report,omitempty"`
	ReportStatus               string              `bson:"report_status" json:"reportStatus"`
	CreatorID                  string              `bson:"creator_id" json:"creatorId"`
	CreatorType                string              `bson:"creator_type" json:"creatorType"`
	AutoSubmit                 bool                `bson:"auto_submit" json:"autoSubmit"`
	BatchID                    primitive.ObjectID  `bson:"batch_id" json:"batchId"`
	TeacherID                  *primitive.ObjectID `bson:"teacher_id,omitempty" json:"teacherId,omitempty"`
	CreatedAt                  time.Time           `bson:"created_at" json:"createdAt"`
	UpdatedAt                  time.Time           `bson:"updated_at" json:"updatedAt"`
}

// AssessmentSection (collection "assessment_sections").
type AssessmentSection struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name             string             `bson:"name" json:"name"`
	MarksPerQuestion int                `bson:"marks_per_question" json:"marksPerQuestion"`
	AssessmentType   string             `bson:"assessment_type" json:"assessmentType"`
	AssessmentID     primitive.ObjectID `bson:"assessment_id" json:"assessmentId"`
	CreatedAt        time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updatedAt"`
}

// AssessmentQuestion (collection "assessment_questions").
type AssessmentQuestion struct {
	ID            primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Question      *string             `bson:"question,omitempty" json:"question,omitempty"`
	Options       []string            `bson:"options" json:"options"`
	CorrectOption *string             `bson:"correct_option,omitempty" json:"correctOption,omitempty"`
	SectionID     primitive.ObjectID  `bson:"section_id" json:"sectionId"`
	ProblemID     *primitive.ObjectID `bson:"problem_id,omitempty" json:"problemId,omitempty"`
	MaxMarks      int                 `bson:"max_marks" json:"maxMarks"`
	CreatedAt     time.Time           `bson:"created_at" json:"createdAt"`
	UpdatedAt     time.Time           `bson:"updated_at" json:"updatedAt"`
}

// AssessmentSubmission (collection "assessment_submissions").
type AssessmentSubmission struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AssessmentID     primitive.ObjectID `bson:"assessment_id" json:"assessmentId"`
	TabSwitchCount   int                `bson:"tab_switch_count" json:"tabSwitchCount"`
	StudentID        primitive.ObjectID `bson:"student_id" json:"studentId"`
	StudentIP        string             `bson:"student_ip" json:"studentIp"`
	ProctoringStatus string             `bson:"proctoring_status" json:"proctoringStatus"`
	StartedAt        time.Time          `bson:"started_at" json:"startedAt"`
	SubmittedAt      *time.Time         `bson:"submitted_at,omitempty" json:"submittedAt,omitempty"`
	TotalMarks       *float64           `bson:"total_marks,omitempty" json:"totalMarks,omitempty"`
	IsSubmitted      bool               `bson:"is_submitted" json:"isSubmitted"`
	CreatedAt        time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updatedAt"`
}

// AssessmentQuestionSubmission (collection "assessment_question_submissions").
type AssessmentQuestionSubmission struct {
	ID                     primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	QuestionID             primitive.ObjectID  `bson:"question_id" json:"questionId"`
	AssessmentID           primitive.ObjectID  `bson:"assessment_id" json:"assessmentId"`
	StudentID              primitive.ObjectID  `bson:"student_id" json:"studentId"`
	Answer                 *string             `bson:"answer,omitempty" json:"answer,omitempty"`
	MarksObtained          float64             `bson:"marks_obtained" json:"marksObtained"`
	AssessmentSubmissionID *primitive.ObjectID `bson:"assessment_submission_id,omitempty" json:"assessmentSubmissionId,omitempty"`
	CreatedAt              time.Time           `bson:"created_at" json:"createdAt"`
	UpdatedAt              time.Time           `bson:"updated_at" json:"updatedAt"`
}
