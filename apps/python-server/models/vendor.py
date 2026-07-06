from datetime import datetime, timezone
from typing import Optional
from beanie import Document
from pydantic import EmailStr, Field


class Vendor(Document):
    name: str
    email: EmailStr
    secondary_email: Optional[str] = None
    tagline: str = ""
    phone_number: str
    secondary_phone_number: Optional[str] = None
    website_link: str = ""
    login_password: str
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    class Settings:
        name = "vendors"
