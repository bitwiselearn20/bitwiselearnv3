import axiosInstance from "@/lib/axios";
import toast from "react-hot-toast";

export const getVendorData = async (statefn: any, paramId: string) => {
  try {
    const getVendor = await axiosInstance.get("/api/v1/vendors/get-vendor-by-id/" + paramId);
    statefn(getVendor.data.data);
  } catch (error) {
    toast.error("failed to get all vendors");
  }
};
