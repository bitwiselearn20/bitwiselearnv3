import axiosInstance from "@/lib/axios";

export type payload = {
    name:string;
    email:string;
    phone?:string;
    message:string;
}

export const sendComplaint = async (payload:payload)=>{
    try {
        const response = await axiosInstance.post(
            "/api/contact",
            payload
        );

        return response.data;
    } catch (error:any) {
        throw new Error("Error in Sending Mail");
    }
}