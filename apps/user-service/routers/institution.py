import secrets
from datetime import datetime, timezone
from fastapi import APIRouter, Depends
from beanie import PydanticObjectId
from shared.schemas.institution import CreateInstitutionRequest, UpdateInstitutionRequest
from shared.utils.api_response import api_response
from shared.utils.password import hash_password
from shared.middleware.auth import get_current_user, require_roles
from shared.models.institution import Institution
from shared.models.batch import Batch
from shared.models.teacher import Teacher
from shared.models.student import Student
from shared.models.course_enrollment import CourseEnrollment
from shared.models.assessment import Assessment
from shared.enums import UserType
from shared.services.email import send_welcome_email

router = APIRouter(prefix="/api/v1/institutions", tags=["Institutions"])

_create_roles = require_roles(UserType.SUPERADMIN, UserType.ADMIN, UserType.VENDOR)
_delete_roles = require_roles(UserType.SUPERADMIN, UserType.ADMIN, UserType.VENDOR)


@router.post("/create-institution")
async def create_institution(body: CreateInstitutionRequest, current_user: dict = Depends(_create_roles)):
    existing = await Institution.find_one(Institution.email == body.email)
    if existing:
        return api_response(400, "Email already exists", error="Duplicate email")

    raw_password = secrets.token_urlsafe(10)
    hashed = hash_password(raw_password)

    vendor_id = None
    if current_user["type"] == UserType.VENDOR:
        vendor_id = PydanticObjectId(current_user["id"])

    inst = Institution(
        name=body.name,
        address=body.address,
        pin_code=body.pin_code,
        tagline=body.tagline,
        website_link=body.website_link,
        login_password=hashed,
        display_password=raw_password,
        email=body.email,
        secondary_email=body.secondary_email,
        phone_number=body.phone_number,
        secondary_phone_number=body.secondary_phone_number,
        created_by=PydanticObjectId(current_user["id"]),
        created_by_vendor_id=vendor_id,
    )
    await inst.insert()

    try:
        send_welcome_email(body.email, body.name, raw_password, "Institution")
    except Exception:
        pass

    return api_response(201, "Institution created", data={
        "id": str(inst.id), "name": inst.name, "email": inst.email
    })


@router.get("/get-all-institution")
async def get_all_institutions(current_user: dict = Depends(get_current_user)):
    if current_user["type"] == UserType.VENDOR:
        institutions = await Institution.find(
            Institution.created_by_vendor_id == PydanticObjectId(current_user["id"])
        ).to_list()
    else:
        institutions = await Institution.find_all().to_list()

    data = [{
        "id": str(i.id), "name": i.name, "email": i.email, "address": i.address,
        "pin_code": i.pin_code, "tagline": i.tagline, "website_link": i.website_link,
        "phone_number": i.phone_number, "secondary_phone_number": i.secondary_phone_number,
        "secondary_email": i.secondary_email, "display_password": i.display_password,
        "created_by": str(i.created_by), "created_by_vendor_id": str(i.created_by_vendor_id) if i.created_by_vendor_id else None,
        "created_at": i.created_at.isoformat()
    } for i in institutions]
    return api_response(200, "Institutions fetched", data=data)


@router.get("/get-institution-by-id/{id}")
async def get_institution_by_id(id: str, current_user: dict = Depends(get_current_user)):
    inst = await Institution.get(PydanticObjectId(id))
    if not inst:
        return api_response(404, "Institution not found", error="Not found")
    return api_response(200, "Institution fetched", data={
        "id": str(inst.id), "name": inst.name, "email": inst.email,
        "address": inst.address, "pin_code": inst.pin_code, "tagline": inst.tagline,
        "website_link": inst.website_link, "phone_number": inst.phone_number,
        "secondary_phone_number": inst.secondary_phone_number, "secondary_email": inst.secondary_email,
        "display_password": inst.display_password,
        "created_by": str(inst.created_by),
        "created_by_vendor_id": str(inst.created_by_vendor_id) if inst.created_by_vendor_id else None,
        "created_at": inst.created_at.isoformat()
    })


@router.get("/get-institution-by-vendor/{id}")
async def get_institutions_by_vendor(id: str, current_user: dict = Depends(get_current_user)):
    institutions = await Institution.find(
        Institution.created_by_vendor_id == PydanticObjectId(id)
    ).to_list()
    data = [{
        "id": str(i.id), "name": i.name, "email": i.email,
        "address": i.address, "phone_number": i.phone_number,
        "created_at": i.created_at.isoformat()
    } for i in institutions]
    return api_response(200, "Institutions fetched", data=data)


@router.put("/update-insititution-by-id/{id}")
async def update_institution(id: str, body: UpdateInstitutionRequest, current_user: dict = Depends(get_current_user)):
    inst = await Institution.get(PydanticObjectId(id))
    if not inst:
        return api_response(404, "Institution not found", error="Not found")

    update_data = body.model_dump(exclude_none=True)
    for key, val in update_data.items():
        setattr(inst, key, val)
    inst.updated_at = datetime.now(timezone.utc)
    await inst.save()
    return api_response(200, "Institution updated", data={"id": str(inst.id), "name": inst.name})


@router.post("/regenerate-password/{id}")
async def regenerate_institution_password(id: str, current_user: dict = Depends(_delete_roles)):
    inst = await Institution.get(PydanticObjectId(id))
    if not inst:
        return api_response(404, "Institution not found", error="Not found")
    raw_password = secrets.token_urlsafe(10)
    inst.login_password = hash_password(raw_password)
    inst.display_password = raw_password
    inst.updated_at = datetime.now(timezone.utc)
    await inst.save()
    return api_response(200, "Password regenerated", data={
        "id": str(inst.id), "email": inst.email, "password": raw_password,
    })


@router.delete("/delete-institution-by-id/{id}")
async def delete_institution(id: str, current_user: dict = Depends(_delete_roles)):
    inst = await Institution.get(PydanticObjectId(id))
    if not inst:
        return api_response(404, "Institution not found", error="Not found")

    oid = PydanticObjectId(id)
    # Cascade: delete batches, teachers, students
    batches = await Batch.find(Batch.institution_id == oid).to_list()
    for b in batches:
        await Student.find(Student.batch_id == b.id).delete()
        await Teacher.find(Teacher.batch_id == b.id).delete()
        await CourseEnrollment.find(CourseEnrollment.batch_id == b.id).delete()
        await b.delete()
    await Teacher.find(Teacher.institute_id == oid).delete()
    await Student.find(Student.institute_id == oid).delete()
    await inst.delete()
    return api_response(200, "Institution deleted")


@router.get("/dashboard/{institution_id}")
async def institution_dashboard(institution_id: str, current_user: dict = Depends(get_current_user)):
    oid = PydanticObjectId(institution_id)
    inst = await Institution.get(oid)
    if not inst:
        return api_response(404, "Institution not found", error="Not found")

    batches = await Batch.find(Batch.institution_id == oid).count()
    students = await Student.find(Student.institute_id == oid).count()
    teachers = await Teacher.find(Teacher.institute_id == oid).count()

    batch_ids = [b.id for b in await Batch.find(Batch.institution_id == oid).to_list()]
    enrollments = await CourseEnrollment.find({"batch_id": {"$in": batch_ids}}).count() if batch_ids else 0
    assessments = await Assessment.find({"batch_id": {"$in": batch_ids}}).count() if batch_ids else 0

    return api_response(200, "Dashboard data", data={
        "institution": {"id": str(inst.id), "name": inst.name},
        "batches": batches,
        "students": students,
        "teachers": teachers,
        "enrollments": enrollments,
        "assessments": assessments,
    })
