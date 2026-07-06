import axiosInstance from "@/lib/axios";

// Student view: coding problems attached to a section, each with completed status.
export const getStudentSectionProblems = async (sectionId: string) => {
  const res = await axiosInstance.get(
    `/api/v1/courses/get-student-section-coding-problems/${sectionId}`,
  );
  return res.data?.data || { problems: [], total: 0, completed: 0 };
};
