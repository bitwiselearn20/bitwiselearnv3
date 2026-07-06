from pydantic import EmailStr
from typing import Optional
from schemas.common import CamelModel


class CreateAdminRequest(CamelModel):
    name: str
    email: EmailStr
    role: str = "ADMIN"


class UpdateAdminRequest(CamelModel):
    name: Optional[str] = None
    email: Optional[EmailStr] = None
