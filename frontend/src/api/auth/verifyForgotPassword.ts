import toast from "react-hot-toast";

export async function verifyForgotPasswordOTP(data: {
  email: string;
  otp: string;
}) {
  const res = await fetch("/api/auth/verify-forgot-password", {
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
