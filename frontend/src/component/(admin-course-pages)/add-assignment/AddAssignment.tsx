import AddAssignmentV1 from "./v1/AddAssignmentV1";

type Props = {
  sectionId: string;
  onClose: () => void;
};

export default function AddAssignment({ sectionId, onClose }: Props) {
  return <AddAssignmentV1 sectionId={sectionId} onClose={onClose} />;
}
