from datetime import datetime, timezone
from typing import Optional
from beanie import Document, PydanticObjectId
from pydantic import Field


class ProblemSolution(Document):
    solution: str
    video_solution: Optional[str] = None
    problem_id: PydanticObjectId
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    class Settings:
        name = "problem_solutions"
        indexes = [
            "problem_id",
        ]
