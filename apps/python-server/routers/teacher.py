import secrets
from datetime import datetime, timezone
from fastapi import APIRouter, Depends
from beanie import PydanticObjectId
from schemas.teacher import CreateTeacherRequest, UpdateTeacherRequest
from utils.api_response import api_response
from utils.password import hash_password
from middleware.auth import get_current_user, not_student
from models.teacher import Teacher
from models.institution import Institution
from models.batch import Batch
from enums import UserType
from services.email import send_welcome_email

router = APIRouter(prefix="/api/v1/teachers", tags=["Teachers"])


@router.post("/create-teacher")
async def create_teacher(body: CreateTeacherRequest, current_user: dict = Depends(not_student)):
    inst = await Institution.get(PydanticObjectId(body.institute_id))
    if not inst:
        return api_response(404, "Institution not found", error="Not found")
    batch = await Batch.get(PydanticObjectId(body.batch_id))
    if not batch:
        return api_response(404, "Batch not found", error="Not found")

    existing = await Teacher.find_one(Teacher.email == body.email)
    if existing:
        return api_response(400, "Email already exists", error="Duplicate email")

    raw_password = secrets.token_urlsafe(10)
    hashed = hash_password(raw_password)

    teacher = Teacher(
        name=body.name,
        email=body.email,
        phone_number=body.phone_number,
        login_password=hashed,
        vendor_id=PydanticObjectId(body.vendor_id) if body.vendor_id else None,
        institute_id=PydanticObjectId(body.institute_id),
        batch_id=PydanticObjectId(body.batch_id),
    )
    await teacher.insert()

    try:
        send_welcome_email(body.email, body.name, raw_password, "Teacher")
    except Exception:
        pass

    return api_response(201, "Teacher created", data={
        "id": str(teacher.id), "name": teacher.name, "email": teacher.email
    })


@router.get("/get-all-teacher")
async def get_all_teachers(current_user: dict = Depends(get_current_user)):
    if current_user["type"] == UserType.INSTITUTION:
        teachers = await Teacher.find(
            Teacher.institute_id == PydanticObjectId(current_user["id"])
        ).to_list()
    elif current_user["type"] == UserType.VENDOR:
        teachers = await Teacher.find(
            Teacher.vendor_id == PydanticObjectId(current_user["id"])
        ).to_list()
    else:
        teachers = await Teacher.find_all().to_list()

    data = [{
        "id": str(t.id), "name": t.name, "email": t.email,
        "phone_number": t.phone_number, "institute_id": str(t.institute_id),
        "batch_id": str(t.batch_id), "created_at": t.created_at.isoformat()
    } for t in teachers]
    return api_response(200, "Teachers fetched", data=data)


@router.get("/get-teacher-by-id/{id}")
async def get_teacher_by_id(id: str, current_user: dict = Depends(get_current_user)):
    teacher = await Teacher.get(PydanticObjectId(id))
    if not teacher:
        return api_response(404, "Teacher not found", error="Not found")
    return api_response(200, "Teacher fetched", data={
        "id": str(teacher.id), "name": teacher.name, "email": teacher.email,
        "phone_number": teacher.phone_number, "institute_id": str(teacher.institute_id),
        "batch_id": str(teacher.batch_id), "vendor_id": str(teacher.vendor_id) if teacher.vendor_id else None,
        "created_at": teacher.created_at.isoformat()
    })


@router.get("/get-teacher-by-institute/{id}")
async def get_teachers_by_institute(id: str, current_user: dict = Depends(get_current_user)):
    teachers = await Teacher.find(Teacher.institute_id == PydanticObjectId(id)).to_list()
    data = [{"id": str(t.id), "name": t.name, "email": t.email,
             "phone_number": t.phone_number, "batch_id": str(t.batch_id)} for t in teachers]
    return api_response(200, "Teachers fetched", data=data)


@router.get("/get-teacher-by-batch/{id}")
async def get_teachers_by_batch(id: str, current_user: dict = Depends(get_current_user)):
    teachers = await Teacher.find(Teacher.batch_id == PydanticObjectId(id)).to_list()
    data = [{"id": str(t.id), "name": t.name, "email": t.email,
             "phone_number": t.phone_number} for t in teachers]
    return api_response(200, "Teachers fetched", data=data)


@router.put("/update-teacher-by-id/{id}")
async def update_teacher(id: str, body: UpdateTeacherRequest, current_user: dict = Depends(not_student)):
    teacher = await Teacher.get(PydanticObjectId(id))
    if not teacher:
        return api_response(404, "Teacher not found", error="Not found")

    update_data = body.model_dump(exclude_none=True)
    for key, val in update_data.items():
        setattr(teacher, key, val)
    teacher.updated_at = datetime.now(timezone.utc)
    await teacher.save()
    return api_response(200, "Teacher updated", data={"id": str(teacher.id), "name": teacher.name})


@router.delete("/delete-teacher-by-id/{id}")
async def delete_teacher(id: str, current_user: dict = Depends(not_student)):
    teacher = await Teacher.get(PydanticObjectId(id))
    if not teacher:
        return api_response(404, "Teacher not found", error="Not found")
    await teacher.delete()
    return api_response(200, "Teacher deleted")
