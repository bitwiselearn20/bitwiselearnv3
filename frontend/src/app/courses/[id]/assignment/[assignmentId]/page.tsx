"use client";
import AttemptAssignment from "@/component/courses/AttempAssignment/AttemptAssignment";
import { useParams } from "next/navigation";

export default function Page() {
  const params = useParams();
  const { assignmentId } = params;
  return <AttemptAssignment assignmentId={assignmentId! as string} />;
}
