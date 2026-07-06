from pydantic import EmailStr
from typing import Optional
from shared.schemas.common import CamelModel


class CreateVendorRequest(CamelModel):
    name: str
    email: EmailStr
    secondary_email: Optional[str] = None
    tagline: str = ""
    phone_number: str
    secondary_phone_number: Optional[str] = None
    website_link: str = ""


class UpdateVendorRequest(CamelModel):
    name: Optional[str] = None
    email: Optional[EmailStr] = None
    secondary_email: Optional[str] = None
    tagline: Optional[str] = None
    phone_number: Optional[str] = None
    secondary_phone_number: Optional[str] = None
    website_link: Optional[str] = None
