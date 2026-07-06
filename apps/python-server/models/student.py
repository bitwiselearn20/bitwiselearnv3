from datetime import datetime, timezone
from typing import Optional
from beanie import Document, PydanticObjectId
from pydantic import EmailStr, Field
from enums import CloudProvider


class Student(Document):
    name: str
    roll_number: str
    email: EmailStr
    login_password: str
    cloud_platform: str = Field(default=CloudProvider.AWS)
    cloudname: Optional[str] = None
    cloudpass: Optional[str] = None
    cloudurl: Optional[str] = None
    batch_id: PydanticObjectId
    institute_id: PydanticObjectId
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    class Settings:
        name = "students"
        indexes = [
            "email",
        ]
