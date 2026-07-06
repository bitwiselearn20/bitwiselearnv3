from datetime import datetime, timezone
from typing import List, Optional
from beanie import Document, PydanticObjectId
from pydantic import Field


class CourseAssignmentSubmission(Document):
    question_id: PydanticObjectId
    student_id: PydanticObjectId
    answer: List[str] = []
    assignment_id: PydanticObjectId
    marks_obtained: Optional[float] = None
    is_correct: Optional[bool] = None
    submitted_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    class Settings:
        name = "course_assignment_submissions"
