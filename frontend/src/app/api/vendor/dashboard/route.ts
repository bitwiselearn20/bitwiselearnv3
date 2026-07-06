import { NextRequest, NextResponse } from "next/server";
import axios from "axios";

export async function GET(req: NextRequest) {
  try {
    const backendUrl = process.env.BACKEND_URL!;
    const cookieHeader = req.headers.get("cookie");

    const response = await axios.get(
      `${backendUrl}/api/v1/vendors/dashboard`,
      {
        headers: { Cookie: cookieHeader || "" },
        withCredentials: true,
      }
    );

    return NextResponse.json(response.data, { status: 200 });
  } catch (error) {
    return NextResponse.json(
      { error: "Failed to fetch vendor dashboard" },
      { status: 500 }
    );
  }
}