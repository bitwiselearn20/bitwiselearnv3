from enum import Enum


class ReportStatus(str, Enum):
    NOT_REQUESTED = "NOT_REQUESTED"
    PROCESSING = "PROCESSING"
    PROCESSED = "PROCESSED"


class CourseStatus(str, Enum):
    PUBLISHED = "PUBLISHED"
    NOT_PUBLISHED = "NOT_PUBLISHED"


class CloudProvider(str, Enum):
    AWS = "AWS"
    GCP = "GCP"
    AZURE = "AZURE"


class ListingStatus(str, Enum):
    LISTED = "LISTED"
    NOT_LISTED = "NOT_LISTED"


class Proctor(str, Enum):
    CHEATED = "CHEATED"
    NOT_CHEATED = "NOT_CHEATED"


class CourseLevel(str, Enum):
    BASIC = "BASIC"
    INTERMEDIATE = "INTERMEDIATE"
    ADVANCE = "ADVANCE"


class ProblemLevel(str, Enum):
    EASY = "EASY"
    MEDIUM = "MEDIUM"
    HARD = "HARD"


class AssignmentType(str, Enum):
    MCQ = "MCQ"
    SCQ = "SCQ"


class ProblemStatus(str, Enum):
    SUCCESS = "SUCCESS"
    FAILED = "FAILED"


class TestcaseType(str, Enum):
    EXAMPLE = "EXAMPLE"
    HIDDEN = "HIDDEN"


class AssessmentType(str, Enum):
    CODE = "CODE"
    NO_CODE = "NO_CODE"


class AssessmentStatus(str, Enum):
    UPCOMING = "UPCOMING"
    LIVE = "LIVE"
    ENDED = "ENDED"


class Languages(str, Enum):
    JAVASCRIPT = "JAVASCRIPT"
    JAVA = "JAVA"
    PYTHON = "PYTHON"
    C = "C"
    CPP = "CPP"


class UserType(str, Enum):
    SUPERADMIN = "SUPERADMIN"
    ADMIN = "ADMIN"
    TEACHER = "TEACHER"
    INSTITUTION = "INSTITUTION"
    VENDOR = "VENDOR"
    STUDENT = "STUDENT"
