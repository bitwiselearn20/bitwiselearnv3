from pydantic import EmailStr
from typing import Optional
from shared.schemas.common import CamelModel


class CreateTeacherRequest(CamelModel):
    name: str
    email: EmailStr
    phone_number: str
    vendor_id: Optional[str] = None
    institute_id: str
    batch_id: str


class UpdateTeacherRequest(CamelModel):
    name: Optional[str] = None
    email: Optional[EmailStr] = None
    phone_number: Optional[str] = None
