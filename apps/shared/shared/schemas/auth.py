from pydantic import EmailStr
from typing import Optional
from shared.schemas.common import CamelModel


class LoginRequest(CamelModel):
    email: EmailStr
    password: str


class ForgotPasswordRequest(CamelModel):
    email: EmailStr
    role: Optional[str] = None


class VerifyOtpRequest(CamelModel):
    email: EmailStr
    otp: str


class ResetPasswordRequest(CamelModel):
    new_password: str
    role: Optional[str] = None
