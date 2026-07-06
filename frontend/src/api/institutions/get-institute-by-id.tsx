import axiosInstance from "@/lib/axios";
import toast from "react-hot-toast";

export const getInstituteData = async (statefn: any, paramId: string) => {
  try {
    const getInstitution = await axiosInstance.get(
      "/api/v1/institutions/get-institution-by-id/" + paramId,
    );
    statefn(getInstitution.data.data);
  } catch (error) {
    toast.error("failed to get institution");
  }
};
