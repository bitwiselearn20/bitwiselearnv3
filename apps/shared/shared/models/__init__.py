from shared.models.user import User
from shared.models.institution import Institution
from shared.models.vendor import Vendor
from shared.models.batch import Batch
from shared.models.teacher import Teacher
from shared.models.student import Student
from shared.models.course import Course
from shared.models.course_section import CourseSection
from shared.models.course_content import CourseLearningContent
from shared.models.course_assignment import CourseAssignment
from shared.models.course_assignment_question import CourseAssignmentQuestion
from shared.models.course_assignment_submission import CourseAssignmentSubmission
from shared.models.course_enrollment import CourseEnrollment
from shared.models.course_progress import CourseProgress
from shared.models.problem import Problem
from shared.models.problem_test_case import ProblemTestCase
from shared.models.problem_solution import ProblemSolution
from shared.models.problem_submission import ProblemSubmission
from shared.models.problem_submission_test_case import ProblemSubmissionTestCase
from shared.models.problem_template import ProblemTemplate
from shared.models.problem_topic import ProblemTopic
from shared.models.assessment import Assessment
from shared.models.assessment_section import AssessmentSection
from shared.models.assessment_question import AssessmentQuestion
from shared.models.assessment_submission import AssessmentSubmission
from shared.models.assessment_question_submission import AssessmentQuestionSubmission

ALL_MODELS = [
    User,
    Institution,
    Vendor,
    Batch,
    Teacher,
    Student,
    Course,
    CourseSection,
    CourseLearningContent,
    CourseAssignment,
    CourseAssignmentQuestion,
    CourseAssignmentSubmission,
    CourseEnrollment,
    CourseProgress,
    Problem,
    ProblemTestCase,
    ProblemSolution,
    ProblemSubmission,
    ProblemSubmissionTestCase,
    ProblemTemplate,
    ProblemTopic,
    Assessment,
    AssessmentSection,
    AssessmentQuestion,
    AssessmentSubmission,
    AssessmentQuestionSubmission,
]
