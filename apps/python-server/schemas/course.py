from typing import Optional, List
from schemas.common import CamelModel


class CreateCourseRequest(CamelModel):
    name: str
    description: str
    level: str = "BASIC"
    duration: Optional[str] = None
    instructor_name: str


class UpdateCourseRequest(CamelModel):
    name: Optional[str] = None
    description: Optional[str] = None
    level: Optional[str] = None
    duration: Optional[str] = None
    instructor_name: Optional[str] = None


class CreateSectionRequest(CamelModel):
    name: str


class UpdateSectionRequest(CamelModel):
    name: Optional[str] = None


class AddContentRequest(CamelModel):
    name: str
    description: Optional[str] = None
    section_id: str
    video_url: str = ""
    transcript: str = ""


class UpdateContentRequest(CamelModel):
    name: Optional[str] = None
    description: Optional[str] = None
    video_url: Optional[str] = None
    transcript: Optional[str] = None


class CreateAssignmentRequest(CamelModel):
    name: str
    description: str
    instruction: str
    marks_per_question: int
    section_id: str


class UpdateAssignmentRequest(CamelModel):
    name: Optional[str] = None
    description: Optional[str] = None
    instruction: Optional[str] = None
    marks_per_question: Optional[int] = None


class AddAssignmentQuestionRequest(CamelModel):
    question: str
    options: List[str]
    correct_answer: List[str]
    type: str = "SCQ"


class UpdateAssignmentQuestionRequest(CamelModel):
    question: Optional[str] = None
    options: Optional[List[str]] = None
    correct_answer: Optional[List[str]] = None
    type: Optional[str] = None


class SubmitAssignmentRequest(CamelModel):
    answers: List[dict]  # [{question_id, answer: [str]}]


class AddEnrollmentRequest(CamelModel):
    course_id: str
    batch_id: str
    institution_id: Optional[str] = None
