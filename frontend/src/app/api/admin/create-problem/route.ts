import axiosInstance from "@/lib/axios";
import { NextRequest, NextResponse } from "next/server";

export async function POST(req: NextRequest) {
  try {
    // Parse JSON body
    const data = await req.json();

    const backendUrl = `${process.env.BACKEND_URL}/api/v1/problems/add-problem/`;
    console.log(backendUrl);
    // Send the data directly without JSON.parse
    const token = req.cookies.get("token") || "";
    if (!token) throw new Error("Token not found");
    const cookieHeader = req.headers.get("cookie");

    const response = await axiosInstance.post(backendUrl, data, {
      headers: {
        Cookie: cookieHeader || "",
      },
      withCredentials: true,
    });
    // console.log(response);
    return NextResponse.json({
      success: true,
      data: response.data,
    });
  } catch (error: any) {
    console.error("Error handling request:", error);
    return NextResponse.json(
      { success: false, message: error.message || "Something went wrong" },
      { status: 400 },
    );
  }
}
