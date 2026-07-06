import { create } from "zustand";
import { persist } from "zustand/middleware";

interface Student {
  data: {
    email: string;
    name: string;
    rollNumber: string;
    batch: {
      id: string;
      batchname: string;
      branch: string;
      batchEndYear: string;
    };
    insitution: {
      id: string;
      name: string;
      tagline: string;
      websiteLink: string;
    };
  };
}

interface StudentStore {
  info: Student | null;
  setData: (data: Student) => void;
  logout: () => void;
}

export const useStudent = create<StudentStore>()(
  persist(
    (set) => ({
      info: null,
      setData: (data) => set({ info: data }),
      logout: () => set({ info: null }),
    }),
    {
      name: "student-storage",
      partialize: (state) => ({
        info: state.info, // persist ONLY data
      }),
    },
  ),
);
