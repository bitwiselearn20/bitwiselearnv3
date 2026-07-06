from datetime import datetime, timedelta, timezone
from typing import Dict
import jwt
from config import get_settings

settings = get_settings()

# Track used/expired tokens
_used_tokens: Dict[str, bool] = {}


def generate_reset_token(email: str, user_type: str, user_id: str) -> str:
    payload = {
        "email": email,
        "type": user_type,
        "id": user_id,
        "exp": datetime.now(timezone.utc) + timedelta(minutes=10),
        "iat": datetime.now(timezone.utc),
    }
    return jwt.encode(payload, settings.RESET_TOKEN_SECRET, algorithm="HS256")


def verify_reset_token(token: str) -> dict | None:
    if _used_tokens.get(token):
        return None
    try:
        payload = jwt.decode(token, settings.RESET_TOKEN_SECRET, algorithms=["HS256"])
        return payload
    except jwt.PyJWTError:
        return None


def invalidate_reset_token(token: str):
    _used_tokens[token] = True
