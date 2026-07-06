import { NextRequest, NextResponse } from "next/server";

export async function PUT(
    req: NextRequest,
    context: { params: Promise<{ id: string }> }
) {
    try {
        const body = await req.json();
        const assessmentId = (await context.params).id;
        const token = req.cookies.get("token") || "";
        if (!token) throw new Error("Token not found");
        const cookieHeader = req.headers.get("cookie");


        const res = await fetch(
            `${process.env.BACKEND_URL}/api/v1/assessments/update-assessment-status/${assessmentId}`, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
                Cookie: cookieHeader || "",
            },
            credentials: "include",
            body: JSON.stringify(body),
        }
        );

        const data = await res.json();

        return NextResponse.json(data, { status: res.status });

    } catch (error: any) {
        return NextResponse.json(
            { message: error.message },
            { status: 500 }
        );
    }
}