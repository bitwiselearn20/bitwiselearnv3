from datetime import datetime, timezone
from beanie import Document, PydanticObjectId
from pydantic import Field


class CourseSection(Document):
    name: str
    creator_id: PydanticObjectId
    course_id: PydanticObjectId
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    class Settings:
        name = "course_sections"
