"use client";

import AdminLoginIMG from "../v1/AdminLoginIMG.png";
import Image from "next/image";
import { Mail, Lock, Eye, EyeOff } from "lucide-react";
import { useState, useEffect } from "react";
import { motion } from "framer-motion";
import handleLogin from "@/api/handleLogin";
import toast from "react-hot-toast";
import { useRouter } from "next/dist/client/components/navigation";
import ResetPasswordForm from "@/component/auth/ResetPasswordForm";
import ForgotPasswordForm from "@/component/auth/ForgotPasswordForm";
import OtpForm from "@/component/auth/OtpForm";
/* ================= ANIMATION VARIANTS ================= */

const pageFade = {
  hidden: { opacity: 0 },
  show: {
    opacity: 1,
    transition: { duration: 0.6, ease: "easeOut" },
  },
};

const slideUp = {
  hidden: { opacity: 0, y: 30 },
  show: {
    opacity: 1,
    y: 0,
    transition: { duration: 0.7, ease: "easeOut" },
  },
};

const stagger = {
  hidden: {},
  show: {
    transition: {
      staggerChildren: 0.15,
    },
  },
};

const imageReveal = {
  hidden: { opacity: 0, scale: 1.05 },
  show: {
    opacity: 1,
    scale: 1,
    transition: { duration: 0.9, ease: "easeOut" },
  },
};

/* ================= TYPEWRITER ================= */

function WelcomeTypewriter() {
  const fullText = "Welcome back, Admin!";
  const adminIndex = fullText.indexOf("Admin");

  const [text, setText] = useState("");
  const [isDeleting, setIsDeleting] = useState(false);
  const [index, setIndex] = useState(0);

  useEffect(() => {
    const speed = isDeleting ? 40 : 80;

    const timeout = setTimeout(() => {
      if (!isDeleting && index < fullText.length) {
        setText(fullText.slice(0, index + 1));
        setIndex(index + 1);
      } else if (isDeleting && index > 0) {
        setText(fullText.slice(0, index - 1));
        setIndex(index - 1);
      } else if (!isDeleting && index === fullText.length) {
        setTimeout(() => setIsDeleting(true), 1000);
      } else if (isDeleting && index === 0) {
        setIsDeleting(false);
      }
    }, speed);

    return () => clearTimeout(timeout);
  }, [index, isDeleting]);

  return (
    <motion.p
      variants={slideUp}
      className="mt-8 text-xl text-neutral-400 h-5 text-center md:text-left"
    >
      {text.slice(0, adminIndex)}
      <span className="text-primaryBlue">{text.slice(adminIndex)}</span>
    </motion.p>
  );
}

/* ================= MAIN COMPONENT ================= */
const ROLE = "ADMIN" as const;

export default function AdminLoginV1() {
  const [step, setStep] = useState<"LOGIN" | "EMAIL" | "OTP" | "RESET">(
    "LOGIN",
  );
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [remember, setRemember] = useState(false);
  const [showPassword, setShowPassword] = useState(false);

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();

  const isDisabled = !email || !password || loading;

  async function fetchLoginData(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();

    try {
      setLoading(true);
      await handleLogin({
        data: { email, password, role: "ADMIN" },
      });
      router.push("/admin-dashboard");
    } catch (err: any) {
      toast.error(err?.response?.data?.error || err?.message || "Login failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <motion.div
      variants={pageFade}
      initial="hidden"
      animate="show"
      className="bg-bg min-h-screen w-screen flex flex-col md:flex-row"
    >
      {/* LEFT SECTION */}
      <motion.div
        variants={stagger}
        initial="hidden"
        animate="show"
        className="flex-1 flex flex-col px-6 py-10 md:p-16"
      >
        <motion.div
          variants={slideUp}
          className="flex justify-center md:justify-start"
        >
          <h1 className="text-3xl">
            <span className="text-primaryBlue font-bold">B</span>
            <span className="font-bold text-white">itwise</span>{" "}
            <span className="text-white">Learn</span>
          </h1>
        </motion.div>

        <motion.div
          key={step} // 🔥 THIS IS THE FIX
          variants={slideUp}
          initial="hidden"
          animate="show"
          className="relative w-full md:w-[60%] bg-divBg mt-10 md:mt-16 md:ml-16 rounded-3xl p-8"
        >
          {step === "LOGIN" && (
            <>
              <h1 className="text-2xl font-bold mb-6">
                <span className="text-white">Log</span>{" "}
                <span className="text-primaryBlue">in</span>
              </h1>

              <motion.form
                variants={stagger}
                onSubmit={fetchLoginData}
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
                      className="w-full pl-12 py-2 pr-4 rounded-lg bg-bg text-white"
                      required
                    />
                  </div>
                </motion.div>

                <motion.div variants={slideUp}>
                  <label className="text-lg text-white">Password</label>
                  <div className="relative">
                    <div className="absolute top-1/2 left-3 -translate-y-1/2">
                      <Lock size={24} color="white" />
                    </div>
                    <input
                      type={showPassword ? "text" : "password"}
                      value={password}
                      onChange={(e) => setPassword(e.target.value)}
                      className="w-full pl-12 py-2 pr-4 rounded-lg bg-bg text-white"
                      required
                    />
                    <button
                      type="button"
                      onClick={() => setShowPassword(!showPassword)}
                      className="absolute top-1/2 right-3 -translate-y-1/2 text-white"
                    >
                      {showPassword ? <Eye size={17} /> : <EyeOff size={17} />}
                    </button>
                  </div>
                </motion.div>

                <motion.button
                  variants={slideUp}
                  type="submit"
                  disabled={loading}
                  className="w-full py-3 rounded-lg bg-primaryBlue text-white font-semibold"
                >
                  {loading ? "Logging in..." : "Log in"}
                </motion.button>

                <motion.button
                  variants={slideUp}
                  type="button"
                  onClick={() => setStep("EMAIL")}
                  className="text-sm text-neutral-300 hover:text-primaryBlue"
                >
                  Forgot Password?
                </motion.button>
              </motion.form>
            </>
          )}
          {step === "EMAIL" && (
            <ForgotPasswordForm
              role={ROLE}
              email={email}
              setEmail={setEmail}
              onSuccess={() => setStep("OTP")}
              onBack={() => setStep("LOGIN")}
            />
          )}

          {step === "OTP" && (
            <OtpForm
              role={ROLE}
              email={email}
              onVerified={() => setStep("RESET")}
            />
          )}

          {step === "RESET" && (
            <ResetPasswordForm
              role={ROLE}
              onSuccess={() => router.push("/admin-login")}
            />
          )}
        </motion.div>

        <WelcomeTypewriter />
      </motion.div>

      {/* RIGHT IMAGE */}
      <motion.div
        variants={imageReveal}
        initial="hidden"
        animate="show"
        className="relative hidden lg:block lg:w-[38%] lg:h-screen"
      >
        <Image
          src={AdminLoginIMG}
          alt="Admin login illustration"
          fill
          className="object-cover"
          priority
        />
      </motion.div>
    </motion.div>
  );
}
