import axiosInstance from "@/lib/axios";
import toast from "react-hot-toast";

export const getAllInstitutions = async (stateFn: any) => {
  try {
    const data = await axiosInstance.get("/api/v1/institutions/get-all-institution");
    stateFn(data.data.data || []);
  } catch (error) {
    toast.error("failed to get all institution");
  }
};
