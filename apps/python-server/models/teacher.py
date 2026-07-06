from datetime import datetime, timezone
from typing import Optional
from beanie import Document, PydanticObjectId
from pydantic import EmailStr, Field


class Teacher(Document):
    name: str
    email: EmailStr
    phone_number: str
    login_password: str
    vendor_id: Optional[PydanticObjectId] = None
    institute_id: PydanticObjectId
    batch_id: PydanticObjectId
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    class Settings:
        name = "teachers"
