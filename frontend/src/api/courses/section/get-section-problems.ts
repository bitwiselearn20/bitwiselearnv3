import axiosInstance from "@/lib/axios";

// Coding problems attached to a section (admin/teacher view).
export const getSectionProblems = async (sectionId: string) => {
  const res = await axiosInstance.get(
    `/api/v1/courses/get-section-coding-problems/${sectionId}`,
  );
  return res.data;
};
