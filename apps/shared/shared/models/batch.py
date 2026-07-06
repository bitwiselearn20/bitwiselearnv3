from datetime import datetime, timezone
from beanie import Document, PydanticObjectId
from pydantic import Field


class Batch(Document):
    batchname: str
    branch: str
    batch_end_year: str
    institution_id: PydanticObjectId
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    class Settings:
        name = "batches"
