import { NextRequest, NextResponse } from "next/server";

export async function DELETE(
  req: NextRequest,
  context: { params: Promise<{ id: string }> },
) {
  try {
    const { id: assignemntId } = await context.params;
    const token = req.cookies.get("token") || "";
    if (!token) throw new Error("Token not found");
    const cookieHeader = req.headers.get("cookie");

    if (!assignemntId) {
      return NextResponse.json(
        { message: "Assignment ID is Required" },
        { status: 400 },
      );
    }

    const res = await fetch(
      `${process.env.BACKEND_URL}/api/v1/courses/remove-assignment-from-section/${assignemntId}`,
      {
        method: "DELETE",
        headers: {
          Cookie: cookieHeader || "",
        },
        credentials: "include",
      },
    );

    const data = await res.json();
    return NextResponse.json(data, { status: res.status });
  } catch (error: any) {
    console.error("Error in deleting Assignment : ", error);
    return NextResponse.json({ message: error.message }, { status: 500 });
  }
}
