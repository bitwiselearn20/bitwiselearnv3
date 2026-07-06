import axiosInstance from "@/lib/axios";
import { useAdmin } from "@/store/adminStore";
import { useInstitution } from "@/store/institutionStore";
import { useStudent } from "@/store/studentStore";
import { useTeacher } from "@/store/teacherStore";
import useVendor from "@/store/vendorStore";

export const logoutUser = async () => {
  await axiosInstance.post("/api/auth/logout");

  useAdmin.getState().logout();
  useInstitution.getState().logout();
  useStudent.getState().logout();
  useTeacher.getState().logout();
  useVendor.getState().logout();

  localStorage.clear();
};

