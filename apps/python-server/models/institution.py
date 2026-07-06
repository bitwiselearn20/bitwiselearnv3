from datetime import datetime, timezone
from typing import Optional
from beanie import Document
from pydantic import EmailStr, Field
from beanie import PydanticObjectId


class Institution(Document):
    name: str
    address: str
    pin_code: str
    tagline: str = ""
    website_link: str = ""
    login_password: str
    email: EmailStr
    secondary_email: Optional[str] = None
    phone_number: str
    secondary_phone_number: Optional[str] = None
    created_by: PydanticObjectId
    created_by_vendor_id: Optional[PydanticObjectId] = None
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    class Settings:
        name = "institutions"
        indexes = [
            "email",
            "phone_number",
        ]
