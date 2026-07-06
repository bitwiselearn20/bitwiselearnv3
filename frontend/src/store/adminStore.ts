import { create } from "zustand";
import { persist } from "zustand/middleware";

interface Admin {
  data: {
    email: string;
    id: string;
    name: string;
    ROLE: string;
  };
}
interface AdminStore {
  info: Admin | null;
  setData: (data: Admin) => void;
  logout: () => void;
}
// admin Information setup
export const useAdmin = create<AdminStore>()(
  persist(
    (set) => ({
      info: null,
      setData: (data) => set({ info: data }),
      logout: () => set({ info: null }),
    }),
    {
      name: "admin-storage",
      partialize: (state) => ({
        info: state.info, // persist ONLY data
      }),
    },
  ),
);
