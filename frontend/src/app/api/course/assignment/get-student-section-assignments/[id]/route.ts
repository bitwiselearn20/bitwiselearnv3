import axios from "axios";
import { NextRequest, NextResponse } from "next/server";

export async function GET(
  req: NextRequest,
  context: { params: Promise<{ id: string }> },
) {
  try {
    const { id } = await context.params;
    const token = req.cookies.get("token") || "";
    if (!token) throw new Error("Token not found");
    const cookieHeader = req.headers.get("cookie");

    const res = await axios.get(
      `${process.env.BACKEND_URL}/api/v1/courses/get-student-section-assignments/${id}`,
      {
        withCredentials: true,
        headers: {
          Cookie: cookieHeader || "",
        },
      },
    );

    const data = res.data.data;
    console.log(data);
    return NextResponse.json(data, { status: res.status });
  } catch (error: any) {
    console.error(error);
    return NextResponse.json({ message: error.message }, { status: 500 });
  }
}
