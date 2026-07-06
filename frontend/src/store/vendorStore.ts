import { create } from "zustand";
import { persist } from "zustand/middleware";

interface Vendor {
  data:{
  email: string;
  id: string;
  name: string;
  tagline: string;
  websiteLink: string;
  phoneNumber: string;
  }
}
interface VendorStore {
  info: Vendor | null;
  setData: (data: Vendor) => void;
  logout: () => void;
}
// vendor Information setup
const useVendor = create<VendorStore>()(
  persist(
    (set) => ({
      info: null,
      setData: (data) => set({ info: data }),
      logout: () => set({ info: null }),
    }),
    {
      name: "vendor-storage",
      partialize: (state) => ({
        info: state.info, // persist ONLY data
      }),
    },
  ),
);

export default useVendor;