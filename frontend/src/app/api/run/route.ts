import { NextRequest, NextResponse } from "next/server";
import axiosInstance from "@/lib/axios";

export async function POST(req: NextRequest) {
  try {
    const backendUrl = process.env.BACKEND_URL;
    const data = await req.json();

    if (!backendUrl) {
      return NextResponse.json(
        { error: "Backend URL not configured" },
        { status: 500 },
      );
    }
    const token = req.cookies.get("token") || "";
    if (!token) throw new Error("Token not found");
    const cookieHeader = req.headers.get("cookie");

    const response = await axiosInstance.post(
      `${backendUrl}/api/v1/code/run`,
      data,
      {
        headers: {
          Cookie: cookieHeader || "",
        },
        withCredentials: true,
      }
    );

    return NextResponse.json(response.data.data, { status: 200 });
  } catch (error: any) {
    console.error("Error fetching institutions:", error.message);

    return NextResponse.json(
      { error: "Failed to fetch institutions" },
      { status: 500 },
    );
  }
}
