"use client";

import { Mail } from "lucide-react";
import { useState } from "react";
import { motion, type Variants } from "framer-motion";
import { sendForgotPasswordOTP } from "@/api/auth/forgot-password";


/* ================= ANIMATION VARIANTS ================= */

const slideUp: Variants = {
    hidden: { opacity: 0, y: 30 },
    show: { opacity: 1, y: 0, transition: { duration: 0.7, ease: "easeOut" as const } },
};

const stagger: Variants = {
    hidden: {},
    show: { transition: { staggerChildren: 0.15 } },
};

/* ================= RESET PASSWORD COMPONENT ================= */

interface ForgotPasswordFormProps {
    role: "STUDENT" | "TEACHER" | "VENDOR" | "ADMIN" | "INSTITUTION";
    email: string;
    setEmail: (email: string) => void;
    onSuccess: () => void;
    onBack: () => void;
}

export default function ForgotPasswordForm({
    role,
    email,
    setEmail,
    onSuccess,
    onBack,
}: ForgotPasswordFormProps) {
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState("");

    async function handleForgot(e: React.FormEvent<HTMLFormElement>) {
        e.preventDefault();
        setLoading(true);
        setMessage("");

        try {
            await sendForgotPasswordOTP({
                email,
                role,
            });

            // move to OTP screen
            onSuccess();
        } catch (err: any) {
            setMessage(err.message || "Failed to send OTP. Please try again.");
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
                <span className="text-white">Forgot</span>{" "}
                <span className="text-primaryBlue">Password</span>
            </motion.h1>

            <motion.p variants={slideUp} className="text-neutral-400 text-sm">
                Enter your registered email. We'll send you a reset link 🔐
            </motion.p>

            <motion.form
                variants={stagger}
                onSubmit={handleForgot}
                className="space-y-6"
            >
                <motion.div variants={slideUp}>
                    <label className="text-lg text-white">Email Address</label>
                    <div className="relative">
                        <div className="absolute top-1/2 left-3 -translate-y-1/2">
                            <Mail size={24} color="white" />
                        </div>
                        <input
                            type="email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            className="w-full pl-12 py-2 pr-4 rounded-lg bg-bg text-white
              focus:ring-2 focus:ring-primaryBlue focus:ring-offset-2 outline-none
              focus:ring-offset-bg"
                            placeholder="johndoe@example.com"
                            required
                        />
                    </div>
                </motion.div>

                {message && (
                    <motion.p
                        variants={slideUp}
                        className={`text-sm ${message.includes("sent")
                            ? "text-green-400"
                            : "text-red-400"
                            }`}
                    >
                        {message}
                    </motion.p>
                )}

                <motion.button
                    variants={slideUp}
                    type="submit"
                    disabled={loading}
                    className="w-full py-3 rounded-lg bg-primaryBlue
          text-white text-lg font-semibold hover:opacity-90 transition disabled:opacity-60"
                >
                    {loading ? "Sending OTP..." : "Send OTP"}

                </motion.button>

                <motion.button
                    variants={slideUp}
                    type="button"
                    onClick={onBack}
                    className="text-sm text-neutral-300 hover:text-primaryBlue transition"
                >
                    ← Back to Login
                </motion.button>
            </motion.form>
        </motion.div>
    );
}
