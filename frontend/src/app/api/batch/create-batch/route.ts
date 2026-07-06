import { NextRequest, NextResponse } from "next/server";

export async function POST(req: NextRequest) {
    try {
        const body = await req.json();
        const token = req.cookies.get("token") || "";
        if (!token) throw new Error("Token not found");
        const cookieHeader = req.headers.get("cookie");


        const res = await fetch(
            `${process.env.BACKEND_URL}/api/v1/batches/create-batch`,
            {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    Cookie: cookieHeader || "",
                },
                credentials: "include",
                body: JSON.stringify(body),
            },
        );

        const data = await res.json();
        return NextResponse.json(data, { status: res.status });
    } catch (error: any) {
        return NextResponse.json({ message: error.message }, { status: 500 });
    }
}
