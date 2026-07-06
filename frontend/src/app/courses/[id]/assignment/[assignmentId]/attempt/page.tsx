import AttemptAssignment from "@/component/courses/AttempAssignment/AttemptAssignment";

export default async function AttemptPage({
  params,
}: {
  params: { assignmentId: string; id: string };
}) {
  const { assignmentId, id } = await params;
  return <AttemptAssignment assignmentId={assignmentId} />;
}
