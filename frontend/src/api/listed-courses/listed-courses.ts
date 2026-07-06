import axiosInstance from "@/lib/axios";

export const getAllListedCourses = async () => {
  const response = await axiosInstance.get("/api/listed-courses");
  return response.data.data; 
};