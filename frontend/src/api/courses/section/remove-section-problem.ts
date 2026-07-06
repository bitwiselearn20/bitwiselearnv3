import axiosInstance from "@/lib/axios";

// Detach a coding problem from a section.
export const removeSectionProblem = async (
  sectionId: string,
  problemId: string,
) => {
  const res = await axiosInstance.delete(
    `/api/v1/courses/remove-coding-problem-from-section/${sectionId}/${problemId}`,
  );
  return res.data;
};
