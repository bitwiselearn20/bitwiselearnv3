import CourseFormV1 from "./v1/CourseFormV1";

type CourseFormProps = {
  onClose: () => void;
  onSuccess: () => void;
};

const CourseForm: React.FC<CourseFormProps> = ({ onClose, onSuccess }) => {
  return <CourseFormV1 onClose={onClose} onSuccess={onSuccess} />;
};

export default CourseForm;
