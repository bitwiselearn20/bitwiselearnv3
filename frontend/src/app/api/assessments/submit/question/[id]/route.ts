import { NextResponse, NextRequest } from "next/server";

export async function POST(
  req: NextRequest,
  context: { params: Promise<{ id: string }> },
) {
  try {
    const { id } = await context.params;
    const input = await req.json();
    const token = req.cookies.get("token") || "";
    if (!token) throw new Error("Token not found");
    const cookieHeader = req.headers.get("cookie");

    const res = await fetch(
      `${process.env.BACKEND_URL}/api/v1/assessments/submit-assessment-question-by-id/${id}`,
      {
        method: "POST",
        body: JSON.stringify(input),
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
    return NextResponse.json({ message: error.message }, { status: 500 });
  }
}
