from pydantic import EmailStr
from typing import Optional
from schemas.common import CamelModel


class CreateInstitutionRequest(CamelModel):
    name: str
    address: str
    pin_code: str
    tagline: str = ""
    website_link: str = ""
    email: EmailStr
    secondary_email: Optional[str] = None
    phone_number: str
    secondary_phone_number: Optional[str] = None


class UpdateInstitutionRequest(CamelModel):
    name: Optional[str] = None
    address: Optional[str] = None
    pin_code: Optional[str] = None
    tagline: Optional[str] = None
    website_link: Optional[str] = None
    email: Optional[EmailStr] = None
    secondary_email: Optional[str] = None
    phone_number: Optional[str] = None
    secondary_phone_number: Optional[str] = None
