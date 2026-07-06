import toast from "react-hot-toast";

export async function sendForgotPasswordOTP(data: {
  email: string;
  role: "STUDENT" | "INSTITUTION" | "ADMIN" | "VENDOR" | "TEACHER";
}) {
  const res = await fetch("/api/auth/forgot-password", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });

  const result = await res.json();

  if (!res.ok) {
    // throw new Error(result.message || "Failed to send OTP");
    toast.error(result.message || "failed to send OTP");
  }

  return result;
}
