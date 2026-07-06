import { create } from "zustand";
import { persist } from "zustand/middleware";
interface Institute {
  data: {
    email: string;
    id: string;
    name: string;
    address: string;
    pinCode: string;
    tagline: string;
    websiteLink: string;
    phoneNumber: string;
  };
}
interface InstituteStore {
  info: Institute | null;
  setData: (data: Institute) => void;
  logout: () => void;
}
// institute Information setup
export const useInstitution = create<InstituteStore>()(
  persist(
    (set) => ({
      info: null,
      setData: (data) => set({ info: data }),
      logout: () => set({ info: null }),
    }),
    {
      name: "institution-storage-v2",
      partialize: (state) => ({
        info: state.info, // persist ONLY data
      }),
    },
  ),
);
