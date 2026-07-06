import { checkJWT } from "@/lib/authJwt";
import axiosInstance from "@/lib/axios";
import { NextRequest, NextResponse } from "next/server";

export async function GET(req: NextRequest) {
  try {
    const backendUrl = process.env.BACKEND_URL;
    console.log("NEXT COOKEIESS");
    console.log(req.cookies);
    if (!backendUrl) {
      return NextResponse.json(
        { error: "Backend URL not configured" },
        { status: 500 },
      );
    }
    // 
    const token = req.cookies.get("token") || "";
    if (!token) throw new Error("Token not found");
    const cookieHeader = req.headers.get("cookie");
    // 
    const response = await axiosInstance.get(

      `${backendUrl}/api/v1/problems/get-user-solved-questions`,
      // 
      {
        headers: {
          Cookie: cookieHeader || "",
        },
        withCredentials: true,
      }
      // 
    );

    return NextResponse.json(response.data.data, { status: 200 });
  } catch (error: any) {
    console.log(error);
    return NextResponse.json({ error: error.message });
  }
}
