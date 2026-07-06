from datetime import datetime, timezone
from typing import Optional
from beanie import Document, PydanticObjectId
from pydantic import Field
from shared.enums import Proctor
import pymongo


class AssessmentSubmission(Document):
    assessment_id: PydanticObjectId
    tab_switch_count: int = 0
    student_id: PydanticObjectId
    student_ip: str
    proctoring_status: str = Field(default=Proctor.NOT_CHEATED)
    started_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    submitted_at: Optional[datetime] = None
    total_marks: Optional[float] = None
    is_submitted: bool = False
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    class Settings:
        name = "assessment_submissions"
        indexes = [
            pymongo.IndexModel(
                [("assessment_id", 1), ("student_id", 1)],
                unique=True,
            ),
        ]
