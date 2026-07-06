import PublishCourseV1 from "./v1/PublishCourseV1";

export type Requirement = {
  label: string;
  satisfied: boolean;
};

export type PublishCourseProps = {
  open: boolean;
  onClose: () => void;
  onConfirm: () => void;
  requirements: Requirement[];
  isPublished: boolean;
};

const PublishCourse = ({
  open,
  onClose,
  onConfirm,
  requirements,
  isPublished,
}: PublishCourseProps) => {
  return (
    <PublishCourseV1
      open={open}
      onClose={onClose}
      onConfirm={onConfirm}
      requirements={requirements}
      isPublished={isPublished}
    />
  );
};

export default PublishCourse;
