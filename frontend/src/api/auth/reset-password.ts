import toast from "react-hot-toast";

export async function resetPassword(data: {
  newPassword: string;
  role: "STUDENT" | "INSTITUTION" | "ADMIN" | "VENDOR" | "TEACHER";
}) {
  const res = await fetch("/api/auth/reset-password", {
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
