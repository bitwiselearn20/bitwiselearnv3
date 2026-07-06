import { checkJWT } from "@/lib/authJwt";
import { NextRequest, NextResponse } from "next/server"


const ROLE_MAP: Record<string, number> = {
    "SUPERADMIN": 0,
    "ADMIN": 1,
    "VENDOR": 2,
    "INSTITUTION": 3,
    "TEACHER": 4,
    "STUDENT": 5,
}
export async function GET(req: NextRequest) {
    try {
        const token = req.cookies.get("token") || "";
        if (!token) throw new Error("Token not found");
        const decodedToken = checkJWT(token.value as any);
        return NextResponse.json({ data: ROLE_MAP[decodedToken.type] }, { status: 200 });
    } catch (error) {
        return NextResponse.json({ error: "cant fetch role" }, { status: 500 });
    }
}