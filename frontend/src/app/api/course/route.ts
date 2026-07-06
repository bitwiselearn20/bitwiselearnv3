import { NextRequest, NextResponse } from "next/server";

export async function GET(req: NextRequest) {
  try {
    const backendUrl = process.env.BACKEND_URL;
    const searchParam = req.nextUrl.searchParams;
    const status = searchParam.get("status");

    if (!backendUrl) {
      return NextResponse.json(
        { error: "Backend URL not configured" },
        { status: 500 },
      );
    }

    const token = req.cookies.get("token") || "";
    if (!token) throw new Error("Token not found");
    const cookieHeader = req.headers.get("cookie");

    const res = await fetch(
      `${backendUrl}/api/v1/courses/get-all-courses-by-admin`,
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Cookie: cookieHeader || "",
        },
        credentials: "include",
        cache: "no-store",
      },
    );

    const data = await res.json();
    console.log(data);
    let publishedCourse: any[] = data.data;
    if (status === "published") {
      publishedCourse = publishedCourse.filter(
        (course: any) => course.isPublished === "PUBLISHED",
      );
    }
    return NextResponse.json(publishedCourse, { status: res.status });
  } catch (error: any) {
    console.log("Error Fetching Courses : ", error);
    return NextResponse.json(
      { message: "Failed to load Courses" },
      { status: 500 },
    );
  }
}
