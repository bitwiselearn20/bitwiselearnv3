import { NextResponse } from "next/server";
import axios from "axios";

export async function GET() {
  try {
    const backendUrl = process.env.BACKEND_URL;

    const response = await axios.get(
      `${backendUrl}/api/v1/courses/listed-courses`
    );

    return NextResponse.json(response.data, {
      status: response.status,
    });

  } catch (error: any) {
    console.error("Listed Courses Error:", error.message);

    return NextResponse.json(
      { message: "Failed to fetch listed courses" },
      { status: 500 }
    );
  }
}