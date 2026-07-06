import axiosInstance from "@/lib/axios";
import { NextRequest, NextResponse } from "next/server";

export async function DELETE(
  req: NextRequest,
  context: { params: Promise<{ id: string }> },
) {
  try {
    const { id } = await context.params;

    const backendUrl = process.env.BACKEND_URL;
    if (!backendUrl) {
      return NextResponse.json(
        { error: "Backend URL not configured" },
        { status: 500 },
      );
    }
    const token = req.cookies.get("token") || "";
    if (!token) throw new Error("Token not found");
    const cookieHeader = req.headers.get("cookie");

    const response = await axiosInstance.delete(
      backendUrl + "/api/v1/problems/delete-problem/" + id,
      {
        headers: {
          Cookie: cookieHeader || "",
        },
        withCredentials: true,
      },
    );

    return NextResponse.json(response.data.data, { status: 200 });
  } catch (error: any) {
    console.error("Error fetching problem:", error.message);

    return NextResponse.json(
      { error: "Failed to fetch problem" },
      { status: 500 },
    );
  }
}
