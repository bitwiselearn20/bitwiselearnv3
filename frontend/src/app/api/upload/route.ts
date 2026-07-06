import { NextRequest, NextResponse } from "next/server";

export async function POST(req: NextRequest) {
  try {
    const token = req.cookies.get("token") || "";
    if (!token) throw new Error("Token not found");
    const cookieHeader = req.headers.get("cookie");
    const formData = await req.formData();
    console.log("route hit");
    const res = await fetch(
      `${process.env.BACKEND_URL}/api/v1/bulk-upload/cloud-info/`,
      {
        method: "POST",
        headers: {
          Cookie: cookieHeader || "",
        },
        credentials: "include",
        body: formData,
      },
    );
    console.log("output recived");
    const data = await res.json();
    return NextResponse.json(data);
  } catch (error: any) {
    return NextResponse.json({ message: error.message }, { status: 500 });
  }
}
