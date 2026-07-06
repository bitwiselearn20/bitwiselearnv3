// /:id
// /unmark-content-as-done/:id

import axiosInstance from "@/lib/axios";
import { NextRequest, NextResponse } from "next/server";

export async function POST(
  req: NextRequest,
  context: { params: Promise<{ id: string }> },
) {
  try {
    const { id } = await context.params;
    const body = await req.json();
    const token = req.cookies.get("token") || "";
    if (!token) throw new Error("Token not found");
    const cookieHeader = req.headers.get("cookie");
    let res;
    if (body.currentStatus === "DONE") {
      res = await axiosInstance.post(
        `${process.env.BACKEND_URL}/api/v1/courses/mark-content-as-done/${id}`,
        {},
        {
          headers: {
            Cookie: cookieHeader || "",
          },
          withCredentials: true,
        }
      );
    } else {
      res = await axiosInstance.post(
        `${process.env.BACKEND_URL}/api/v1/courses/unmark-content-as-done/${id}`,
        {},
        {
          headers: {
            Cookie: cookieHeader || "",
          },
          withCredentials: true,
        }
      );
    }

    const data = res.data;

    return NextResponse.json(data, { status: res.status });
  } catch (error: any) {
    return NextResponse.json({ message: error.message }, { status: 500 });
  }
}
