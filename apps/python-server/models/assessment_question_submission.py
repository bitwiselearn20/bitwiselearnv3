from datetime import datetime, timezone
from typing import Optional
from beanie import Document, PydanticObjectId
from pydantic import Field


class AssessmentQuestionSubmission(Document):
    question_id: PydanticObjectId
    assessment_id: PydanticObjectId
    student_id: PydanticObjectId
    answer: Optional[str] = None
    marks_obtained: float = 0
    assessment_submission_id: Optional[PydanticObjectId] = None
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    class Settings:
        name = "assessment_question_submissions"
