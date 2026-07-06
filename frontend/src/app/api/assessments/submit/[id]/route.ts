import axiosInstance from "@/lib/axios";
import { NextResponse, NextRequest } from "next/server";

export async function POST(
  req: NextRequest,
  context: { params: Promise<{ id: string }> },
) {
  try {
    const { id } = await context.params;
    const token = req.cookies.get("token") || "";
    if (!token) throw new Error("Token not found");
    const cookieHeader = req.headers.get("cookie");

    const body = await req.json();
    const res = await fetch(
      `${process.env.BACKEND_URL}/api/v1/assessments/submit-assessment-by-id/${id}`,
      {
        method: "POST",
        body: JSON.stringify(body),
        headers: {
          "Content-Type": "application/json",
          Cookie: cookieHeader || "",
        },
        credentials: "include",
      },
    );

    const data = await res.json();

    return NextResponse.json(data, { status: res.status });
  } catch (error: any) {
    console.log(JSON.stringify(error));
    return NextResponse.json({ message: error.message }, { status: 500 });
  }
}
