import React from "react";
import IndividualCourseReportV1 from "./v1/IndividualCourseReportV1";

function IndividualCourseReport({
  courseId,
  batchId,
}: {
  courseId: string;
  batchId: string;
}) {
  return <IndividualCourseReportV1 courseId={courseId} batchId={batchId} />;
}

export default IndividualCourseReport;
