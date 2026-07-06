"use client";

import { useState } from "react";
import { motion, type Variants } from "framer-motion";
import { Lock, Eye, EyeOff } from "lucide-react";
import { resetPassword } from "@/api/auth/reset-password";

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

type Role =
    | "STUDENT"
    | "TEACHER"
    | "VENDOR"
    | "ADMIN"
    | "INSTITUTION";

interface ResetPasswordFormProps {
    resetToken?: string;
    role: Role;
    onSuccess: () => void;
}

/* ================= COMPONENT ================= */

export default function ResetPasswordForm({
    resetToken,
    role,
    onSuccess,
}: ResetPasswordFormProps) {
    const [newPassword, setNewPassword] = useState("");
    const [showPassword, setShowPassword] = useState(false);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");

    async function handleReset() {
        if (!newPassword) return;

        setLoading(true);
        setError("");

        try {
            await resetPassword({
                newPassword,
                role,
            });
            ;

            onSuccess();
        } catch (err: any) {
            setError(err.message || "Password reset failed");
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
                <span className="text-white">Reset</span>{" "}
                <span className="text-primaryBlue">Password</span>
            </motion.h1>

            <motion.div variants={slideUp}>
                <label className="text-lg text-white">New Password</label>
                <div className="relative">
                    <div className="absolute top-1/2 left-3 -translate-y-1/2">
                        <Lock size={22} color="white" />
                    </div>
                    <input
                        type={showPassword ? "text" : "password"}
                        value={newPassword}
                        onChange={(e) => setNewPassword(e.target.value)}
                        className="w-full pl-12 py-2 pr-4 rounded-lg bg-bg text-white
            focus:ring-2 focus:ring-primaryBlue focus:ring-offset-2 outline-none
            focus:ring-offset-bg"
                        placeholder="Enter new password"
                    />
                    <button
                        type="button"
                        onClick={() => setShowPassword(!showPassword)}
                        className="absolute top-1/2 right-3 -translate-y-1/2 text-white"
                    >
                        {showPassword ? <Eye size={18} /> : <EyeOff size={18} />}
                    </button>
                </div>
            </motion.div>

            {error && (
                <motion.p variants={slideUp} className="text-sm text-red-400">
                    {error}
                </motion.p>
            )}

            <motion.button
                variants={slideUp}
                onClick={handleReset}
                disabled={loading}
                className="w-full py-3 rounded-lg bg-primaryBlue
        text-white text-lg font-semibold hover:opacity-90 transition disabled:opacity-60"
            >
                {loading ? "Resetting..." : "Reset Password"}
            </motion.button>
        </motion.div>
    );
}
