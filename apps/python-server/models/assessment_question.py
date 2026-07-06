from datetime import datetime, timezone
from typing import List, Optional
from beanie import Document, PydanticObjectId
from pydantic import Field


class AssessmentQuestion(Document):
    question: Optional[str] = None
    options: List[str] = []
    correct_option: Optional[str] = None
    section_id: PydanticObjectId
    problem_id: Optional[PydanticObjectId] = None
    max_marks: int
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    class Settings:
        name = "assessment_questions"
