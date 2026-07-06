import { NextRequest, NextResponse } from "next/server";

export async function DELETE(
  req: NextRequest,
  context: { params: Promise<{ id: string }> },
) {
  try {
    const { id } = await context.params;
    const token = req.cookies.get("token") || "";
    if (!token) throw new Error("Token not found");
    const cookieHeader = req.headers.get("cookie");

    const res = await fetch(
      `${process.env.BACKEND_URL}/api/v1/courses/delete-content/${id}`,
      {
        method: "DELETE",
        headers: {
          Cookie: cookieHeader || "",
        },
        credentials: "include",
      }
    );

    const data = await res.json();

    return NextResponse.json(data, { status: res.status });
  } catch (error: any) {
    console.error(error);
    return NextResponse.json({ message: error.message }, { status: 500 });
  }
}
