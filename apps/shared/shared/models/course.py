from datetime import datetime, timezone
from typing import Optional
from beanie import Document, PydanticObjectId
from pydantic import Field
from shared.enums import CourseStatus, CourseLevel


class Course(Document):
    name: str
    description: str
    level: str = Field(default=CourseLevel.BASIC)
    duration: Optional[str] = None
    thumbnail: Optional[str] = None
    instructor_name: str
    certificate: Optional[str] = None
    is_published: str = Field(default=CourseStatus.NOT_PUBLISHED)
    created_by: PydanticObjectId
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    class Settings:
        name = "courses"
        indexes = [
            "name",
        ]
