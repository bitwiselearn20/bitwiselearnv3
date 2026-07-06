import { NextRequest, NextResponse } from "next/server";

export async function GET(
    req: NextRequest,
    context: { params: Promise<{ id: string }> }
) {
    try {
        const sectionId = (await context.params).id;
        const token = req.cookies.get("token") || "";
        if (!token) throw new Error("Token not found");
        const cookieHeader = req.headers.get("cookie");


        const res = await fetch(
            `${process.env.BACKEND_URL}/api/v1/assessments/get-questions-by-sectionId/${sectionId}`,
            {
                method: "GET",
                headers: {
                    Cookie: cookieHeader || "",
                },
                credentials: "include",
            }
        );

        const data = await res.json();
        return NextResponse.json(data, { status: res.status });

    } catch (error: any) {
        return NextResponse.json(
            { message: error.message },
            { status: 500 }
        )
    }
}