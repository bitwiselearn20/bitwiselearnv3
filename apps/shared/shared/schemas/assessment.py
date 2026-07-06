from typing import Optional, List
from datetime import datetime
from shared.schemas.common import CamelModel


class CreateAssessmentRequest(CamelModel):
    name: str
    description: str
    instruction: str
    start_time: datetime
    end_time: datetime
    individual_section_time_limit: Optional[int] = None
    auto_submit: bool = True
    batch_id: str
    teacher_id: Optional[str] = None


class UpdateAssessmentRequest(CamelModel):
    name: Optional[str] = None
    description: Optional[str] = None
    instruction: Optional[str] = None
    start_time: Optional[datetime] = None
    end_time: Optional[datetime] = None
    individual_section_time_limit: Optional[int] = None
    auto_submit: Optional[bool] = None


class UpdateAssessmentStatusRequest(CamelModel):
    status: str


class CreateAssessmentSectionRequest(CamelModel):
    name: str
    marks_per_question: int
    assessment_type: str
    assessment_id: str


class UpdateAssessmentSectionRequest(CamelModel):
    name: Optional[str] = None
    marks_per_question: Optional[int] = None
    assessment_type: Optional[str] = None


class AddAssessmentQuestionRequest(CamelModel):
    question: Optional[str] = None
    options: List[str] = []
    correct_option: Optional[str] = None
    problem_id: Optional[str] = None
    max_marks: int


class UpdateAssessmentQuestionRequest(CamelModel):
    question: Optional[str] = None
    options: Optional[List[str]] = None
    correct_option: Optional[str] = None
    max_marks: Optional[int] = None


class SubmitAssessmentRequest(CamelModel):
    tab_switch_count: int = 0
    proctoring_status: str = "NOT_CHEATED"
    student_ip: str = ""


class SubmitAssessmentQuestionRequest(CamelModel):
    question_id: str
    answer: Optional[str] = None
    code: Optional[str] = None
    language: Optional[str] = None
