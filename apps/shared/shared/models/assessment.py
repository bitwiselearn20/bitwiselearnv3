from datetime import datetime, timezone
from typing import Optional
from beanie import Document, PydanticObjectId
from pydantic import Field
from shared.enums import AssessmentStatus, ReportStatus


class Assessment(Document):
    name: str
    description: str
    instruction: str
    start_time: datetime
    end_time: datetime
    individual_section_time_limit: Optional[int] = None
    status: str = Field(default=AssessmentStatus.UPCOMING)
    report: Optional[str] = None
    report_status: str = Field(default=ReportStatus.NOT_REQUESTED)
    creator_id: str
    creator_type: str = "ADMIN"
    auto_submit: bool = True
    batch_id: PydanticObjectId
    teacher_id: Optional[PydanticObjectId] = None
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    class Settings:
        name = "assessments"
