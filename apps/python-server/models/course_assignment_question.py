from datetime import datetime, timezone
from typing import List
from beanie import Document, PydanticObjectId
from pydantic import Field
from enums import AssignmentType


class CourseAssignmentQuestion(Document):
    question: str
    options: List[str] = []
    correct_answer: List[str] = []
    assignment_id: PydanticObjectId
    type: str = Field(default=AssignmentType.SCQ)
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    class Settings:
        name = "course_assignment_questions"
