from datetime import datetime, timezone
from typing import Optional
from beanie import Document, PydanticObjectId
from pydantic import Field


class CourseLearningContent(Document):
    name: str
    description: Optional[str] = None
    creator_id: PydanticObjectId
    section_id: PydanticObjectId
    video_url: str = ""
    transcript: str = ""
    file: str = ""
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    class Settings:
        name = "course_learning_contents"
