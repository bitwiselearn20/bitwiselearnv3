import { NextRequest, NextResponse } from "next/server";
import axiosInstance from "@/lib/axios";

export async function POST(req: NextRequest) {
  try {
    const backendUrl = process.env.BACKEND_URL;

    if (!backendUrl) {
      return NextResponse.json(
        { error: "Backend URL not configured" },
        { status: 500 },
      );
    }

    const body = await req.json();
    const token = req.cookies.get("token") || "";
    if (!token) throw new Error("Token not found");
    const cookieHeader = req.headers.get("cookie");
    const response = await axiosInstance.post(
      `${backendUrl}/api/v1/teachers/create-teacher`,
      body,
      {
        headers: {
          Cookie: cookieHeader || "",
        },
        withCredentials: true,
      }
    );

    return NextResponse.json(response.data, { status: 201 });
  } catch (error: any) {
    console.error("Error creating teacher:", error.message);

    return NextResponse.json(
      {
        error: "Failed to create teacher",
        details: error?.response?.data || null,
      },
      { status: error?.response?.status || 500 },
    );
  }
}
