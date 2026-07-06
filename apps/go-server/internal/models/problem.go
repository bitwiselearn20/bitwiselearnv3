package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Collection names for the DSA problem domain (ported from models/problem*.py).
const (
	CollProblems                   = "problems"
	CollProblemTopics              = "problem_topics"
	CollProblemTemplates           = "problem_templates"
	CollProblemTestCases           = "problem_test_cases"
	CollProblemSolutions           = "problem_solutions"
	CollProblemSubmissions         = "problem_submissions"
	CollProblemSubmissionTestCases = "problem_submission_test_cases"
)

// Problem (collection "problems").
type Problem struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Name        string              `bson:"name" json:"name"`
	Description string              `bson:"description" json:"description"`
	Hints       []string            `bson:"hints" json:"hints"`
	Difficulty  string              `bson:"difficulty" json:"difficulty"`
	CreatedBy   *string             `bson:"created_by,omitempty" json:"createdBy,omitempty"`
	CreatorType string              `bson:"creator_type" json:"creatorType"`
	SectionID   *primitive.ObjectID `bson:"section_id,omitempty" json:"sectionId,omitempty"`
	Published   string              `bson:"published" json:"published"`
	UserID      *primitive.ObjectID `bson:"user_id,omitempty" json:"userId,omitempty"`
	CreatedAt   time.Time           `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time           `bson:"updated_at" json:"updatedAt"`
}

// ProblemTopic (collection "problem_topics").
type ProblemTopic struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProblemID primitive.ObjectID `bson:"problem_id" json:"problemId"`
	TagName   []string           `bson:"tag_name" json:"tagName"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}

// ProblemTemplate (collection "problem_templates").
type ProblemTemplate struct {
	ID           primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	ProblemID    primitive.ObjectID  `bson:"problem_id" json:"problemId"`
	FunctionBody string              `bson:"function_body" json:"functionBody"`
	DefaultCode  string              `bson:"default_code" json:"defaultCode"`
	Language     string              `bson:"language" json:"language"`
	StudentsID   *primitive.ObjectID `bson:"students_id,omitempty" json:"studentsId,omitempty"`
	CreatedAt    time.Time           `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time           `bson:"updated_at" json:"updatedAt"`
}

// ProblemTestCase (collection "problem_test_cases").
type ProblemTestCase struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TestType  string             `bson:"test_type" json:"testType"`
	Input     string             `bson:"input" json:"input"`
	Output    string             `bson:"output" json:"output"`
	ProblemID primitive.ObjectID `bson:"problem_id" json:"problemId"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}

// ProblemSolution (collection "problem_solutions").
type ProblemSolution struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Solution      string             `bson:"solution" json:"solution"`
	VideoSolution *string            `bson:"video_solution,omitempty" json:"videoSolution,omitempty"`
	ProblemID     primitive.ObjectID `bson:"problem_id" json:"problemId"`
	CreatedAt     time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updatedAt"`
}

// ProblemSubmission (collection "problem_submissions").
type ProblemSubmission struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Code           string             `bson:"code" json:"code"`
	Runtime        *string            `bson:"runtime,omitempty" json:"runtime,omitempty"`
	Memory         *string            `bson:"memory,omitempty" json:"memory,omitempty"`
	Status         string             `bson:"status" json:"status"`
	StudentID      primitive.ObjectID `bson:"student_id" json:"studentId"`
	ProblemID      primitive.ObjectID `bson:"problem_id" json:"problemId"`
	FailedTestCase *string            `bson:"failed_test_case,omitempty" json:"failedTestCase,omitempty"`
	SubmittedAt    time.Time          `bson:"submitted_at" json:"submittedAt"`
	CreatedAt      time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updatedAt"`
}

// ProblemSubmissionTestCase (collection "problem_submission_test_cases").
type ProblemSubmissionTestCase struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SubmissionID primitive.ObjectID `bson:"submission_id" json:"submissionId"`
	TestCaseID   primitive.ObjectID `bson:"test_case_id" json:"testCaseId"`
	Passed       bool               `bson:"passed" json:"passed"`
	ActualOutput *string            `bson:"actual_output,omitempty" json:"actualOutput,omitempty"`
	Runtime      *string            `bson:"runtime,omitempty" json:"runtime,omitempty"`
	Memory       *string            `bson:"memory,omitempty" json:"memory,omitempty"`
	CreatedAt    time.Time          `bson:"created_at" json:"createdAt"`
}
