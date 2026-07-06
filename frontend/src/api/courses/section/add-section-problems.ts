import axiosInstance from "@/lib/axios";

// Attach DSA coding problems to a section (max 5 enforced server-side).
export const addSectionProblems = async (
  sectionId: string,
  problemIds: string[],
) => {
  const res = await axiosInstance.post(
    `/api/v1/courses/add-coding-problems-to-section/${sectionId}`,
    { problemIds },
  );
  return res.data;
};
