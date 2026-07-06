import cloudinary
import cloudinary.uploader
import base64
from config import get_settings
from urllib.parse import urlparse

settings = get_settings()

cloudinary.config(
    cloud_name=settings.CLOUDINARY_CLOUD_NAME,
    api_key=settings.CLOUDINARY_API_KEY,
    api_secret=settings.CLOUDINARY_API_SECRET,
)


def upload_to_cloudinary(file_bytes: bytes, folder: str, filename: str) -> str:
    b64 = base64.b64encode(file_bytes).decode("utf-8")
    data_uri = f"data:application/octet-stream;base64,{b64}"
    result = cloudinary.uploader.upload(
        data_uri,
        folder=folder,
        public_id=filename,
        resource_type="auto",
    )
    return result.get("secure_url", "")


def delete_from_cloudinary(file_url: str) -> bool:
    try:
        parsed = urlparse(file_url)
        path = parsed.path
        # Extract public_id: remove /vXXXX/ version and file extension
        parts = path.split("/upload/")
        if len(parts) < 2:
            return False
        after_upload = parts[1]
        # Remove version prefix if present (e.g., v1234567890/)
        segments = after_upload.split("/", 1)
        if segments[0].startswith("v") and segments[0][1:].isdigit():
            after_upload = segments[1] if len(segments) > 1 else segments[0]
        # Remove file extension
        public_id = after_upload.rsplit(".", 1)[0]
        cloudinary.uploader.destroy(public_id)
        return True
    except Exception:
        return False
