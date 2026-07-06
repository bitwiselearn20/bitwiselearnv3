"use client";

import AddSectionV2 from "./v2/AddSectionV2";

type Props = {
  sectionNumber: number;
  sectionId: string;
  sectionData: {
    id: string;
    name: string;
    courseLearningContents: {
      id: string;
      name: string;
      description: string;
      transcript: string;
    }[];
  };
  onAddAssignment: (sectionId: string) => void;
  onSectionDeleted: () => void;
};

const AddSection = ({
  sectionId,
  sectionNumber,
  sectionData,
  onAddAssignment,
  onSectionDeleted,
}: Props) => {
  return (
    <AddSectionV2
      sectionId={sectionId}
      sectionNumber={sectionNumber}
      sectionData={sectionData}
      onAddAssignment={() => onAddAssignment}
      onSectionDeleted={() => onSectionDeleted}
    />
  );
};

export default AddSection;
