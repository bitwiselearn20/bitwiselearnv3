import { NextRequest,NextResponse } from "next/server";
import axiosInstance from "@/lib/axios";

export async function POST(req:NextRequest){
    try {
        const body = await req.json();
        const response = await axiosInstance.post(
            `${process.env.BACKEND_URL}/api/v1/contact`,
            body,
            {
                headers:{
                    "Content-Type":"application/json",
                },
            },
        );

    return NextResponse.json(response.data, {
      status: response.status,
    });
  } catch (error: any) {
    return NextResponse.json(
      {
        message: error.response?.data?.message || error.message || "",
      },
      {
        status: error.response?.status || 500,
      },
    );
  }
}