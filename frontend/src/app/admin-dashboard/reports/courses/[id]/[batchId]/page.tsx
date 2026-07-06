"use client";
import IndividualCourseReport from "@/component/report/individualCourseReport/IndividualCourseReport";
import { useParams } from "next/navigation";

function page() {
  const params = useParams();
  return (
    <div>
      <IndividualCourseReport
        courseId={params.id as string}
        batchId={params.batchId as string}
      />
    </div>
  );
}

export default page;
