import { NextRequest, NextResponse } from "next/server";
import axiosInstance from "@/lib/axios";

export async function GET(req: NextRequest) {
  try {
    const backendUrl = process.env.BACKEND_URL;
    if (!backendUrl) {
      return NextResponse.json(
        { error: "Backend URL not configured" },
        { status: 500 },
      );
    }
    const params = req.nextUrl.searchParams;
    const isAdminRequest = params.get("q") === "valid";
    const token = req.cookies.get("token") || "";
    if (!token) throw new Error("Token not found");
    const cookieHeader = req.headers.get("cookie");
    let response;
    if (isAdminRequest) {
      response = await axiosInstance.get(
        `${backendUrl}/api/v1/problems/get-all-dsa-problem`,
        {
          headers: {
            Cookie: cookieHeader || "",
          },
          withCredentials: true,
        }
      );
    } else {
      response = await axiosInstance.get(
        `${backendUrl}/api/v1/problems/get-all-listed-problem`,
        {
          headers: {
            Cookie: cookieHeader || "",
          },
          withCredentials: true,
        }
      );
    }
    return NextResponse.json(response.data.data, { status: 200 });
  } catch (error: any) {
    console.error("Error fetching problem:", error.message);
    console.log(error);
    return NextResponse.json(
      { error: "Failed to fetch problem" },
      { status: 500 },
    );
  }
}
