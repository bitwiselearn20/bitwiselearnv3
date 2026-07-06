from datetime import datetime, timezone
from typing import List
from beanie import Document, PydanticObjectId
from pydantic import Field


class ProblemTopic(Document):
    problem_id: PydanticObjectId
    tag_name: List[str] = []
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    class Settings:
        name = "problem_topics"
