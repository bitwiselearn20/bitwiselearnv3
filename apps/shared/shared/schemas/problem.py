from typing import Optional, List
from shared.schemas.common import CamelModel


class CreateProblemRequest(CamelModel):
    name: str
    description: str
    hints: List[str] = []
    difficulty: str = "EASY"
    section_id: Optional[str] = None


class UpdateProblemRequest(CamelModel):
    name: Optional[str] = None
    description: Optional[str] = None
    hints: Optional[List[str]] = None
    difficulty: Optional[str] = None


class AddTopicRequest(CamelModel):
    tag_name: List[str]


class UpdateTopicRequest(CamelModel):
    tag_name: List[str]


class AddTemplateRequest(CamelModel):
    function_body: str
    default_code: str
    language: str


class UpdateTemplateRequest(CamelModel):
    function_body: Optional[str] = None
    default_code: Optional[str] = None
    language: Optional[str] = None


class AddTestCaseRequest(CamelModel):
    test_type: str = "EXAMPLE"
    input: str
    output: str


class UpdateTestCaseRequest(CamelModel):
    test_type: Optional[str] = None
    input: Optional[str] = None
    output: Optional[str] = None


class AddSolutionRequest(CamelModel):
    solution: str
    video_solution: Optional[str] = None


class UpdateSolutionRequest(CamelModel):
    solution: Optional[str] = None
    video_solution: Optional[str] = None


class RunCodeRequest(CamelModel):
    code: str
    language: str
    problem_id: str


class CompileCodeRequest(CamelModel):
    code: str
    language: str
    stdin: str = ""


class SubmitCodeRequest(CamelModel):
    code: str
    language: str
    problem_id: str


class SearchProblemRequest(CamelModel):
    query: str
