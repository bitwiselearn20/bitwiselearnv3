import AssignmentV2 from "./v2/AssignmentV2";

export default function Assignment({
  assignments = [],
  map = {},
}: {
  assignments?: any[];
  map: Object;
}) {
  return <AssignmentV2 assignments={assignments} map={map} />;
}
