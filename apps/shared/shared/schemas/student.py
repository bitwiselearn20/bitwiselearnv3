from pydantic import EmailStr
from typing import Optional
from shared.schemas.common import CamelModel


class CreateStudentRequest(CamelModel):
    name: str
    roll_number: str
    email: EmailStr
    login_password: Optional[str] = None
    batch_id: str
    institution_id: str
    cloud_platform: str = "AWS"


class UpdateStudentRequest(CamelModel):
    name: Optional[str] = None
    roll_number: Optional[str] = None
    email: Optional[EmailStr] = None
    cloud_platform: Optional[str] = None
    cloudname: Optional[str] = None
    cloudpass: Optional[str] = None
    cloudurl: Optional[str] = None
