import axiosInstance from "@/lib/axios";

export async function handleReport(id: string) {
  await axiosInstance.get("/api/full-assessment-report/" + id);
}
