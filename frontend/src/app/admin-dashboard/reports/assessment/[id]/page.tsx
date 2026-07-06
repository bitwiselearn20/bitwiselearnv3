"use client";
import IndividualAssessmentReport from "@/component/report/individualAssessmentReport/IndividualAssessmentReport";
import { useParams } from "next/navigation";

function page() {
  const params = useParams();
  return <IndividualAssessmentReport assessmentId={params.id as string} />;
}

export default page;
