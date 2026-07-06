import secrets
import time
from typing import Dict, Tuple

# In-memory OTP store: email -> (otp, expiry_timestamp)
_otp_store: Dict[str, Tuple[str, float]] = {}

OTP_EXPIRY_SECONDS = 600  # 10 minutes


def _cleanup_expired():
    now = time.time()
    expired = [k for k, v in _otp_store.items() if v[1] < now]
    for k in expired:
        del _otp_store[k]


def generate_otp(email: str) -> str:
    _cleanup_expired()
    otp = str(secrets.randbelow(900000) + 100000)  # 6-digit
    _otp_store[email.lower()] = (otp, time.time() + OTP_EXPIRY_SECONDS)
    return otp


def verify_otp(email: str, otp: str) -> bool:
    _cleanup_expired()
    key = email.lower()
    stored = _otp_store.get(key)
    if not stored:
        return False
    if stored[0] == otp and stored[1] > time.time():
        del _otp_store[key]
        return True
    return False
