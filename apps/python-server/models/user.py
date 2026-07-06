from datetime import datetime, timezone
from beanie import Document
from pydantic import EmailStr, Field
from enums import UserType


class User(Document):
    name: str
    email: EmailStr
    password: str
    role: str = Field(default=UserType.SUPERADMIN)
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    updated_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    class Settings:
        name = "users"

    def save_updated(self):
        self.updated_at = datetime.now(timezone.utc)
