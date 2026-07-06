import CourseBuilderV1 from "./v1/CourseBuilderV1";

type PageProps = {
  params: {
    id: string;
  };
};

const CourseBuilder = ({ params }: PageProps) => {
  return (
    <div>
      <CourseBuilderV1 courseId={params.id} />
    </div>
  );
};

export default CourseBuilder;
