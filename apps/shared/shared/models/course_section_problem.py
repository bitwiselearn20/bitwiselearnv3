from datetime import datetime, timezone
from beanie import Document, PydanticObjectId
from pydantic import Field
import pymongo


class CourseSectionProblem(Document):
    """Links a DSA coding problem to a course section (submodule).

    A section may have at most 5 coding problems (enforced in the service).
    Student completion is derived from a successful ProblemSubmission.
    """

    section_id: PydanticObjectId
    problem_id: PydanticObjectId
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    class Settings:
        name = "course_section_problems"
        indexes = [
            "section_id",
            pymongo.IndexModel(
                [("section_id", 1), ("problem_id", 1)],
                unique=True,
            ),
        ]
