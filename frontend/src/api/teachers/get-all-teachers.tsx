import axiosInstance from "@/lib/axios";
import toast from "react-hot-toast";

export const getAllTeachers = async (statefn: any) => {
  try {
    const getAllTeachers = await axiosInstance.get("/api/teacher/");
    statefn(getAllTeachers.data);
  } catch (error) {
    toast.error("failed to get all teachers");
  }
};
