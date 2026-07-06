from typing import Optional
from shared.schemas.common import CamelModel


class CreateBatchRequest(CamelModel):
    batchname: str
    branch: str
    batch_end_year: str
    institution_id: str


class UpdateBatchRequest(CamelModel):
    batchname: Optional[str] = None
    branch: Optional[str] = None
    batch_end_year: Optional[str] = None
