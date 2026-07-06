import axiosInstance from "@/lib/axios";
import axios from "axios";
import toast from "react-hot-toast";

export const sebdAssessmentReport = async (id: string) => {
  try {
    await axiosInstance.get("/api/assessment-report/" + id);
  } catch (error) {
    toast.error("error sending assessment report");
  }
};
