from models.user import User
from models.institution import Institution
from models.vendor import Vendor
from models.batch import Batch
from models.teacher import Teacher
from models.student import Student
from models.course import Course
from models.course_section import CourseSection
from models.course_content import CourseLearningContent
from models.course_assignment import CourseAssignment
from models.course_assignment_question import CourseAssignmentQuestion
from models.course_assignment_submission import CourseAssignmentSubmission
from models.course_enrollment import CourseEnrollment
from models.course_progress import CourseProgress
from models.problem import Problem
from models.problem_test_case import ProblemTestCase
from models.problem_solution import ProblemSolution
from models.problem_submission import ProblemSubmission
from models.problem_submission_test_case import ProblemSubmissionTestCase
from models.problem_template import ProblemTemplate
from models.problem_topic import ProblemTopic
from models.assessment import Assessment
from models.assessment_section import AssessmentSection
from models.assessment_question import AssessmentQuestion
from models.assessment_submission import AssessmentSubmission
from models.assessment_question_submission import AssessmentQuestionSubmission

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
