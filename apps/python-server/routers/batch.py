from datetime import datetime, timezone
from fastapi import APIRouter, Depends
from beanie import PydanticObjectId
from schemas.batch import CreateBatchRequest, UpdateBatchRequest
from utils.api_response import api_response
from middleware.auth import get_current_user, require_roles, not_student
from models.batch import Batch
from models.institution import Institution
from models.student import Student
from models.teacher import Teacher
from models.course_enrollment import CourseEnrollment
from enums import UserType

router = APIRouter(prefix="/api/v1/batches", tags=["Batches"])

_delete_roles = require_roles(UserType.SUPERADMIN, UserType.ADMIN, UserType.INSTITUTION, UserType.VENDOR)


@router.post("/create-batch")
async def create_batch(body: CreateBatchRequest, current_user: dict = Depends(not_student)):
    inst = await Institution.get(PydanticObjectId(body.institution_id))
    if not inst:
        return api_response(404, "Institution not found", error="Not found")

    batch = Batch(
        batchname=body.batchname,
        branch=body.branch,
        batch_end_year=body.batch_end_year,
        institution_id=PydanticObjectId(body.institution_id),
    )
    await batch.insert()
    return api_response(201, "Batch created", data={
        "id": str(batch.id), "batchname": batch.batchname,
        "branch": batch.branch, "batch_end_year": batch.batch_end_year,
        "institution_id": str(batch.institution_id)
    })


@router.get("/get-all-batch")
async def get_all_batches_global(current_user: dict = Depends(not_student)):
    batches = await Batch.find_all().to_list()
    data = []
    for b in batches:
        student_count = await Student.find(Student.batch_id == b.id).count()
        data.append({
            "id": str(b.id), "batchname": b.batchname, "branch": b.branch,
            "batch_end_year": b.batch_end_year, "institution_id": str(b.institution_id),
            "student_count": student_count, "created_at": b.created_at.isoformat()
        })
    return api_response(200, "Batches fetched", data=data)


@router.get("/get-all-batch/{id}")
async def get_all_batches(id: str, current_user: dict = Depends(not_student)):
    batches = await Batch.find(Batch.institution_id == PydanticObjectId(id)).to_list()
    data = []
    for b in batches:
        student_count = await Student.find(Student.batch_id == b.id).count()
        data.append({
            "id": str(b.id), "batchname": b.batchname, "branch": b.branch,
            "batch_end_year": b.batch_end_year, "institution_id": str(b.institution_id),
            "student_count": student_count, "created_at": b.created_at.isoformat()
        })
    return api_response(200, "Batches fetched", data=data)


@router.get("/get-batch-by-id/{id}")
async def get_batch_by_id(id: str, current_user: dict = Depends(get_current_user)):
    batch = await Batch.get(PydanticObjectId(id))
    if not batch:
        return api_response(404, "Batch not found", error="Not found")

    students = await Student.find(Student.batch_id == batch.id).to_list()
    student_data = [{
        "id": str(s.id), "name": s.name, "email": s.email,
        "roll_number": s.roll_number
    } for s in students]

    return api_response(200, "Batch fetched", data={
        "id": str(batch.id), "batchname": batch.batchname, "branch": batch.branch,
        "batch_end_year": batch.batch_end_year, "institution_id": str(batch.institution_id),
        "students": student_data, "created_at": batch.created_at.isoformat()
    })


@router.put("/update-batch-by-id/{id}")
async def update_batch(id: str, body: UpdateBatchRequest, current_user: dict = Depends(not_student)):
    batch = await Batch.get(PydanticObjectId(id))
    if not batch:
        return api_response(404, "Batch not found", error="Not found")

    update_data = body.model_dump(exclude_none=True)
    for key, val in update_data.items():
        setattr(batch, key, val)
    batch.updated_at = datetime.now(timezone.utc)
    await batch.save()
    return api_response(200, "Batch updated", data={"id": str(batch.id), "batchname": batch.batchname})


@router.delete("/delete-batch-by-id/{id}")
async def delete_batch(id: str, current_user: dict = Depends(_delete_roles)):
    batch = await Batch.get(PydanticObjectId(id))
    if not batch:
        return api_response(404, "Batch not found", error="Not found")

    bid = batch.id
    await Student.find(Student.batch_id == bid).delete()
    await Teacher.find(Teacher.batch_id == bid).delete()
    await CourseEnrollment.find(CourseEnrollment.batch_id == bid).delete()
    await batch.delete()
    return api_response(200, "Batch deleted")
