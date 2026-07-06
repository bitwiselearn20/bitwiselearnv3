import axiosInstance from "@/lib/axios";

let failureCount = 0;

export const getHealth = async (stateFn: any) => {
  const startTime = Date.now();

  try {
    await axiosInstance.get("/health");
    const latency = Date.now() - startTime;

    failureCount = 0; 

    if (latency < 500) {
      stateFn("good");
    } else if (latency < 2000) {
      stateFn("slow");
    } else {
      stateFn("bad");
    }

  } catch {
    failureCount += 1;

    if (failureCount >= 2) {
      stateFn("bad");
    }
  }
};
