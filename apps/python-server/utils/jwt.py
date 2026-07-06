from datetime import datetime, timedelta, timezone
import jwt
from config import get_settings

settings = get_settings()


def generate_access_token(user_id: str, user_type: str) -> str:
    payload = {
        "id": user_id,
        "type": user_type,
        "exp": datetime.now(timezone.utc) + timedelta(days=1),
        "iat": datetime.now(timezone.utc),
    }
    return jwt.encode(payload, settings.JWT_ACCESS_SECRET, algorithm="HS256")


def generate_refresh_token(user_id: str, user_type: str) -> str:
    payload = {
        "id": user_id,
        "type": user_type,
        "exp": datetime.now(timezone.utc) + timedelta(days=20),
        "iat": datetime.now(timezone.utc),
    }
    return jwt.encode(payload, settings.JWT_REFRESH_SECRET, algorithm="HS256")


def verify_access_token(token: str) -> dict | None:
    try:
        return jwt.decode(token, settings.JWT_ACCESS_SECRET, algorithms=["HS256"])
    except jwt.PyJWTError:
        return None


def verify_refresh_token(token: str) -> dict | None:
    try:
        return jwt.decode(token, settings.JWT_REFRESH_SECRET, algorithms=["HS256"])
    except jwt.PyJWTError:
        return None
