import axiosInstance from "@/lib/axios";
import toast from "react-hot-toast";

export const getAllProblemSubmission = async (statefn: any, id: string) => {
  try {
    const getProblem = await axiosInstance.get(
      "/api/v1/problems/get-submission/" + id,
    );
    // unwrap the API envelope { statusCode, message, data, error } -> data array
    const list = getProblem.data?.data;
    statefn(Array.isArray(list) ? list : []);
  } catch (error) {
    toast.error("failed to get problem submission");
    statefn([]);
  }
};
