package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Collection names for the course domain (ported from models/course*.py).
const (
	CollCourses                     = "courses"
	CollCourseSections              = "course_sections"
	CollCourseLearningContents      = "course_learning_contents"
	CollCourseAssignments           = "course_assignments"
	CollCourseAssignmentQuestions   = "course_assignment_questions"
	CollCourseAssignmentSubmissions = "course_assignment_submissions"
	CollCourseEnrollments           = "course_enrollments"
	CollCourseProgresses            = "course_progresses"
	CollBatches                     = "batches"
)

// Course (collection "courses").
type Course struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name           string             `bson:"name" json:"name"`
	Description    string             `bson:"description" json:"description"`
	Level          string             `bson:"level" json:"level"`
	Duration       *string            `bson:"duration,omitempty" json:"duration,omitempty"`
	Thumbnail      *string            `bson:"thumbnail,omitempty" json:"thumbnail,omitempty"`
	InstructorName string             `bson:"instructor_name" json:"instructorName"`
	Certificate    *string            `bson:"certificate,omitempty" json:"certificate,omitempty"`
	IsPublished    string             `bson:"is_published" json:"isPublished"`
	CreatedBy      primitive.ObjectID `bson:"created_by" json:"createdBy"`
	CreatedAt      time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updatedAt"`
}

// CourseSection (collection "course_sections").
type CourseSection struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	CreatorID primitive.ObjectID `bson:"creator_id" json:"creatorId"`
	CourseID  primitive.ObjectID `bson:"course_id" json:"courseId"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}

// CourseLearningContent (collection "course_learning_contents").
type CourseLearningContent struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description *string            `bson:"description,omitempty" json:"description,omitempty"`
	CreatorID   primitive.ObjectID `bson:"creator_id" json:"creatorId"`
	SectionID   primitive.ObjectID `bson:"section_id" json:"sectionId"`
	VideoURL    string             `bson:"video_url" json:"videoUrl"`
	Transcript  string             `bson:"transcript" json:"transcript"`
	File        string             `bson:"file" json:"file"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
}

// CourseAssignment (collection "course_assignments").
type CourseAssignment struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name             string             `bson:"name" json:"name"`
	Description      string             `bson:"description" json:"description"`
	Instruction      string             `bson:"instruction" json:"instruction"`
	MarksPerQuestion int                `bson:"marks_per_question" json:"marksPerQuestion"`
	SectionID        primitive.ObjectID `bson:"section_id" json:"sectionId"`
	CreatedAt        time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updatedAt"`
}

// CourseAssignmentQuestion (collection "course_assignment_questions").
type CourseAssignmentQuestion struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Question      string             `bson:"question" json:"question"`
	Options       []string           `bson:"options" json:"options"`
	CorrectAnswer []string           `bson:"correct_answer" json:"correctAnswer"`
	AssignmentID  primitive.ObjectID `bson:"assignment_id" json:"assignmentId"`
	Type          string             `bson:"type" json:"type"`
	CreatedAt     time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updatedAt"`
}

// CourseAssignmentSubmission (collection "course_assignment_submissions").
type CourseAssignmentSubmission struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	QuestionID    primitive.ObjectID `bson:"question_id" json:"questionId"`
	StudentID     primitive.ObjectID `bson:"student_id" json:"studentId"`
	Answer        []string           `bson:"answer" json:"answer"`
	AssignmentID  primitive.ObjectID `bson:"assignment_id" json:"assignmentId"`
	MarksObtained *float64           `bson:"marks_obtained,omitempty" json:"marksObtained,omitempty"`
	IsCorrect     *bool              `bson:"is_correct,omitempty" json:"isCorrect,omitempty"`
	SubmittedAt   time.Time          `bson:"submitted_at" json:"submittedAt"`
	CreatedAt     time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updatedAt"`
}

// CourseEnrollment (collection "course_enrollments").
type CourseEnrollment struct {
	ID            primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	CourseID      primitive.ObjectID  `bson:"course_id" json:"courseId"`
	BatchID       primitive.ObjectID  `bson:"batch_id" json:"batchId"`
	InstitutionID *primitive.ObjectID `bson:"institution_id,omitempty" json:"institutionId,omitempty"`
	EnrolledAt    time.Time           `bson:"enrolled_at" json:"enrolledAt"`
	CreatedAt     time.Time           `bson:"created_at" json:"createdAt"`
	UpdatedAt     time.Time           `bson:"updated_at" json:"updatedAt"`
}

// CourseProgress (collection "course_progresses").
type CourseProgress struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	StudentID primitive.ObjectID `bson:"student_id" json:"studentId"`
	ContentID primitive.ObjectID `bson:"content_id" json:"contentId"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}

// Batch (collection "batches").
type Batch struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BatchName     string             `bson:"batchname" json:"batchname"`
	Branch        string             `bson:"branch" json:"branch"`
	BatchEndYear  string             `bson:"batch_end_year" json:"batchEndYear"`
	InstitutionID primitive.ObjectID `bson:"institution_id" json:"institutionId"`
	CreatedAt     time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updatedAt"`
}
