import axiosInstance from "@/lib/axios";
import toast from "react-hot-toast";

export const createInstitution = async (data: any) => {
  try {
    const createInstitution = await axiosInstance.post(
      "/api/v1/institutions/create-institution",
      data,
    );

    return createInstitution.data;
  } catch (error: any) {
    toast.error(error?.response?.data?.error || "failed to create institution");
    throw error;
  }
};
