from datetime import datetime, timezone
from typing import Optional
from beanie import Document, PydanticObjectId
from pydantic import Field
import pymongo


class ProblemSubmissionTestCase(Document):
    submission_id: PydanticObjectId
    test_case_id: PydanticObjectId
    passed: bool
    actual_output: Optional[str] = None
    runtime: Optional[str] = None
    memory: Optional[str] = None
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    class Settings:
        name = "problem_submission_test_cases"
        indexes = [
            pymongo.IndexModel(
                [("submission_id", 1), ("test_case_id", 1)],
                unique=True,
            ),
        ]
