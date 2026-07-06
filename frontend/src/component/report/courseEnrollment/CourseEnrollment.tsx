import React from "react";
import CourseEnrollmentV1 from "./v1/CourseEnrollmentV1";

function CourseEnrollment({ courseId }: { courseId: string }) {
  return <CourseEnrollmentV1 courseId={courseId} />;
}

export default CourseEnrollment;
