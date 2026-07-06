import axiosInstance from "@/lib/axios";
import toast from "react-hot-toast";

export const getVendorInstitutions = async (statefn: any, paramId: string) => {
  try {
    const getInstitution = await axiosInstance.get(
      "/api/v1/institutions/get-institution-by-vendor/" + paramId,
    );
    statefn(getInstitution.data.data || []);
  } catch (error) {
    toast.error("failed to fetch vendor institutions");
  }
};
