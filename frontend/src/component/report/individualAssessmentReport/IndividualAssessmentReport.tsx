import React from "react";
import IndividualAssessmentReportV1 from "./v1/IndividualAssessmentReportV1";

function IndividualAssessmentReport({
  assessmentId,
}: {
  assessmentId: string;
}) {
  return <IndividualAssessmentReportV1 assessmentId={assessmentId} />;
}

export default IndividualAssessmentReport;
