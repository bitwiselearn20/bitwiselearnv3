import axiosInstance from "@/lib/axios";
import { NextRequest, NextResponse } from "next/server";

export async function GET(
  req: NextRequest,
  context: { params: Promise<{ id: string }> }
) {
  try {
    const { id } = await context.params;

    const cookieHeader = req.headers.get("cookie");

    const res = await axiosInstance.get(
      `${process.env.BACKEND_URL}/api/v1/courses/get-individual-course-progress/${id}`,
      {
        headers: {
          Cookie: cookieHeader || "",
        },
        withCredentials: true,
      }
    );

    return NextResponse.json(res.data, { status: res.status });
  } catch (error: any) {
    return NextResponse.json(
      { message: error.message },
      { status: 500 }
    );
  }
}