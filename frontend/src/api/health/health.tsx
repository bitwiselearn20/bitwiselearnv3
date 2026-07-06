import axiosInstance from "@/lib/axios";
import toast from "react-hot-toast";
export const getHealth = async (stateFn: any) => {
  try {
    const data = await axiosInstance.get("/health");
    stateFn(data.data);
  } catch (error) {
    toast.error("unknown error occured");
  }
};
