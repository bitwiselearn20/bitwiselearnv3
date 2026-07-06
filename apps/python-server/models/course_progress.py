from datetime import datetime, timezone
from beanie import Document, PydanticObjectId
from pydantic import Field
import pymongo


class CourseProgress(Document):
    student_id: PydanticObjectId
    content_id: PydanticObjectId
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    class Settings:
        name = "course_progresses"
        indexes = [
            pymongo.IndexModel(
                [("student_id", 1), ("content_id", 1)],
                unique=True,
            ),
        ]
