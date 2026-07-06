package assessment

import "time"

// Request bodies ported from schemas/assessment.py (CamelModel -> camelCase JSON).

type createAssessmentRequest struct {
	Name                       string    `json:"name"`
	Description                string    `json:"description"`
	Instruction                string    `json:"instruction"`
	StartTime                  time.Time `json:"startTime"`
	EndTime                    time.Time `json:"endTime"`
	IndividualSectionTimeLimit *int      `json:"individualSectionTimeLimit"`
	AutoSubmit                 bool      `json:"autoSubmit"`
	BatchID                    string    `json:"batchId"`
	TeacherID                  *string   `json:"teacherId"`
}

type updateAssessmentRequest struct {
	Name                       *string    `json:"name"`
	Description                *string    `json:"description"`
	Instruction                *string    `json:"instruction"`
	StartTime                  *time.Time `json:"startTime"`
	EndTime                    *time.Time `json:"endTime"`
	IndividualSectionTimeLimit *int       `json:"individualSectionTimeLimit"`
	AutoSubmit                 *bool      `json:"autoSubmit"`
}

type updateAssessmentStatusRequest struct {
	Status string `json:"status"`
}

type createAssessmentSectionRequest struct {
	Name             string `json:"name"`
	MarksPerQuestion int    `json:"marksPerQuestion"`
	AssessmentType   string `json:"assessmentType"`
	AssessmentID     string `json:"assessmentId"`
}

type updateAssessmentSectionRequest struct {
	Name             *string `json:"name"`
	MarksPerQuestion *int    `json:"marksPerQuestion"`
	AssessmentType   *string `json:"assessmentType"`
}

type addAssessmentQuestionRequest struct {
	Question      *string  `json:"question"`
	Options       []string `json:"options"`
	CorrectOption *string  `json:"correctOption"`
	ProblemID     *string  `json:"problemId"`
	MaxMarks      int      `json:"maxMarks"`
}

type updateAssessmentQuestionRequest struct {
	Question      *string   `json:"question"`
	Options       *[]string `json:"options"`
	CorrectOption *string   `json:"correctOption"`
	MaxMarks      *int      `json:"maxMarks"`
}

type submitAssessmentRequest struct {
	TabSwitchCount   int    `json:"tabSwitchCount"`
	ProctoringStatus string `json:"proctoringStatus"`
	StudentIP        string `json:"studentIp"`
}

type submitAssessmentQuestionRequest struct {
	QuestionID string  `json:"questionId"`
	Answer     *string `json:"answer"`
	Code       *string `json:"code"`
	Language   *string `json:"language"`
}
