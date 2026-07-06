import axiosInstance from "@/lib/axios";
import { NextRequest, NextResponse } from "next/server";

export async function GET(req: NextRequest) {
  try {
    const token = req.cookies.get("token") || "";
    if (!token) throw new Error("Token not found");
    const cookieHeader = req.headers.get("cookie");

    const res = await axiosInstance.get(
      `${process.env.BACKEND_URL}/api/v1/courses/get-all-course-progress`,
      {
        headers: {
          Cookie: cookieHeader || "",
        },
        withCredentials: true,
      },
    );

    return NextResponse.json(res.data, { status: res.status });
  } catch (error: any) {
    return NextResponse.json(
      { message: error.message },
      { status: 500 },
    );
  }
}