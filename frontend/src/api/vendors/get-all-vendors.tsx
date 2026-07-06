import axiosInstance from "@/lib/axios";
import toast from "react-hot-toast";

export const getAllVendors = async (stateFn: any) => {
  try {
    const data = await axiosInstance.get("/api/v1/vendors/get-all-vendor");
    stateFn(data.data.data || []);
  } catch (error) {
    toast.error("failed to get vendors");
  }
};
