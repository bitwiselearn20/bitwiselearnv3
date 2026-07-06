import { create } from "zustand";
import { persist } from "zustand/middleware";

interface Teacher {
  data: {
      name: string;
      email: string;
      phoneNumber: string;
      institution: {
        id: string;
        name: string;
      };
      batch: {
        id: string;
        batchName: string;
        branch: string;
      };
  }
}
interface TeacherStore {
  info: Teacher | null;
  setData: (data: Teacher) => void;
  logout: () => void;
}
// vendor Information setup
export const useTeacher = create<TeacherStore>()(
  persist(
    (set) => ({
      info: null,
      setData: (data) => set({ info: data }),
      logout: () => set({ info: null }),
    }),
    {
      name: "teacher-storage",
      partialize: (state) => ({
        info: state.info, // persist ONLY data
      }),
    },
  ),
);
