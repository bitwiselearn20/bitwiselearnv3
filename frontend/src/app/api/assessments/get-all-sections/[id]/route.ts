import { NextRequest, NextResponse } from "next/server";

export async function GET(
  req: NextRequest,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const assessmentId = (await params).id;
    const token = req.cookies.get("token") || "";
    if (!token) throw new Error("Token not found");
    const cookieHeader = req.headers.get("cookie");


    if (!assessmentId) {
      return NextResponse.json(
        { message: "Assessment ID is required" },
        { status: 400 },
      );
    }

    const backendRes = await fetch(
      `${process.env.BACKEND_URL}/api/v1/assessments/get-sections-for-assessment/${assessmentId}`,
      {
        method: "GET",
        headers: {
          Cookie: cookieHeader || "",
        },
        credentials: "include",
      },
    );

    const data = await backendRes.json();

    return NextResponse.json(data, {
      status: backendRes.status,
    });
  } catch (error: any) {
    return NextResponse.json({ message: error.message }, { status: 500 });
  }
}
