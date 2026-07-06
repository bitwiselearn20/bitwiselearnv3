"use client";

import { useState } from "react";
import { motion, type Variants } from "framer-motion";
import { KeyRound } from "lucide-react";
import { verifyForgotPasswordOTP } from "@/api/auth/verifyForgotPassword";
/* ================= ANIMATIONS ================= */

const slideUp: Variants = {
  hidden: { opacity: 0, y: 30 },
  show: {
    opacity: 1,
    y: 0,
    transition: { duration: 0.6, ease: "easeOut" as const },
  },
};

const stagger: Variants = {
  hidden: {},
  show: { transition: { staggerChildren: 0.15 } },
};

/* ================= TYPES ================= */

type Role = "STUDENT" | "TEACHER" | "VENDOR" | "ADMIN" | "INSTITUTION";

interface OtpFormProps {
  email: string;
  role: Role;
  onVerified: (resetToken: string) => void;
  onBack?: () => void;
}

/* ================= COMPONENT ================= */

export default function OtpForm({
  email,
  role,
  onVerified,
  onBack,
}: OtpFormProps) {
  const [otp, setOtp] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function handleVerify() {
    if (!otp) return;

    setLoading(true);
    setError("");

    try {
      const data = await verifyForgotPasswordOTP({
        email,
        otp,
      });

      onVerified(data.data.resetToken);
    } catch (err: any) {
      setError(err.message || "Invalid OTP");
    } finally {
      setLoading(false);
    }
  }

  return (
    <motion.div
      variants={stagger}
      initial="hidden"
      animate="show"
      className="space-y-6"
    >
      <motion.h1 variants={slideUp} className="text-2xl font-bold">
        <span className="text-white">Verify</span>{" "}
        <span className="text-primaryBlue">OTP</span>
      </motion.h1>

      <motion.p variants={slideUp} className="text-sm text-neutral-400">
        Enter the OTP sent to <span className="text-white">{email}</span>
      </motion.p>

      <motion.div variants={slideUp}>
        <label className="text-lg text-white">OTP</label>
        <div className="relative">
          <div className="absolute top-1/2 left-3 -translate-y-1/2">
            <KeyRound size={22} color="white" />
          </div>
          <input
            type="text"
            value={otp}
            onChange={(e) => setOtp(e.target.value)}
            className="w-full pl-12 py-2 pr-4 rounded-lg bg-bg text-white
            focus:ring-2 focus:ring-primaryBlue focus:ring-offset-2 outline-none
            focus:ring-offset-bg"
            placeholder="Enter 6-digit OTP"
            maxLength={6}
          />
        </div>
      </motion.div>

      {error && (
        <motion.p variants={slideUp} className="text-sm text-red-400">
          {error}
        </motion.p>
      )}

      <motion.button
        type="button"
        variants={slideUp}
        onClick={handleVerify}
        disabled={loading}
        className="w-full py-3 rounded-lg bg-primaryBlue
        text-white text-lg font-semibold hover:opacity-90 transition disabled:opacity-60"
      >
        {loading ? "Verifying..." : "Verify OTP"}
      </motion.button>

      {onBack && (
        <motion.button
          variants={slideUp}
          type="button"
          onClick={onBack}
          className="text-sm text-neutral-300 hover:text-primaryBlue transition"
        >
          ← Back
        </motion.button>
      )}
    </motion.div>
  );
}
