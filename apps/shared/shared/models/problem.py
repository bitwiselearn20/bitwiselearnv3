from datetime import datetime, timezone
from typing import List, Optional
from beanie import Document, PydanticObjectId
from pydantic import Field
from shared.enums import ProblemLevel, ListingStatus


class Problem(Document):
    name: str
    description: str
    hints: List[str] = []
    difficulty: str = Field(default=ProblemLevel.EASY)
    created_by: Optional[str] = None
    creator_type: str = "ADMIN"
    section_id: Optional[PydanticObjectId] = None
    published: str = Field(default=ListingStatus.NOT_LISTED)
    user_id: Optional[PydanticObjectId] = None
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    class Settings:
        name = "problems"
