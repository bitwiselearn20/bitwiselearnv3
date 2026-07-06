package course

// Request bodies ported from schemas/course.py (CamelModel -> camelCase JSON).

type createCourseRequest struct {
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	Level          string  `json:"level"`
	Duration       *string `json:"duration"`
	InstructorName string  `json:"instructorName"`
}

type updateCourseRequest struct {
	Name           *string `json:"name"`
	Description    *string `json:"description"`
	Level          *string `json:"level"`
	Duration       *string `json:"duration"`
	InstructorName *string `json:"instructorName"`
}

type createSectionRequest struct {
	Name string `json:"name"`
}

type updateSectionRequest struct {
	Name *string `json:"name"`
}

type addContentRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	SectionID   string  `json:"sectionId"`
	VideoURL    string  `json:"videoUrl"`
	Transcript  string  `json:"transcript"`
}

type updateContentRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	VideoURL    *string `json:"videoUrl"`
	Transcript  *string `json:"transcript"`
}

type createAssignmentRequest struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	Instruction      string `json:"instruction"`
	MarksPerQuestion int    `json:"marksPerQuestion"`
	SectionID        string `json:"sectionId"`
}

type updateAssignmentRequest struct {
	Name             *string `json:"name"`
	Description      *string `json:"description"`
	Instruction      *string `json:"instruction"`
	MarksPerQuestion *int    `json:"marksPerQuestion"`
}

type addAssignmentQuestionRequest struct {
	Question      string   `json:"question"`
	Options       []string `json:"options"`
	CorrectAnswer []string `json:"correctAnswer"`
	Type          string   `json:"type"`
}

type updateAssignmentQuestionRequest struct {
	Question      *string   `json:"question"`
	Options       *[]string `json:"options"`
	CorrectAnswer *[]string `json:"correctAnswer"`
	Type          *string   `json:"type"`
}

type submittedAnswer struct {
	QuestionID string   `json:"questionId"`
	Answer     []string `json:"answer"`
}

type submitAssignmentRequest struct {
	Answers []submittedAnswer `json:"answers"`
}

type addEnrollmentRequest struct {
	CourseID      string  `json:"courseId"`
	BatchID       string  `json:"batchId"`
	InstitutionID *string `json:"institutionId"`
}
