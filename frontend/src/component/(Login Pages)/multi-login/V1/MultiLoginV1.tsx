"use client";

type LoginRole = "INSTITUTION" | "VENDOR" | "TEACHER";

import MutliLoginIMG from "../V1/MultiLoginIMG.png";
import Image from "next/image";
import {
  Mail,
  Lock,
  Eye,
  EyeOff,
  GraduationCap,
  School,
  Handshake,
} from "lucide-react";
import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";

import { motion } from "framer-motion";
import handleLogin from "@/api/handleLogin";
import toast from "react-hot-toast";
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
  const fullText = "Teach with purpose. Lead with impact.";

  const [text, setText] = useState("");
  const [isDeleting, setIsDeleting] = useState(false);
  const [index, setIndex] = useState(0);

  useEffect(() => {
    const speed = isDeleting ? 30 : 40;

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

  const renderText = () => {
    const parts = text.split(/(Teach|Lead)/g);

    return parts.map((part, i) =>
      part === "Teach" || part === "Lead" ? (
        <span key={i} className="text-primaryBlue">
          {part}
        </span>
      ) : (
        <span key={i}>{part}</span>
      ),
    );
  };

  return (
    <motion.p
      variants={slideUp}
      className="mt-8 text-xl text-neutral-400 h-5 text-center md:text-left"
    >
      {renderText()}
    </motion.p>
  );
}

/* ================= MAIN COMPONENT ================= */

export default function AdminLoginV1() {
  const [step, setStep] = useState<"LOGIN" | "EMAIL" | "OTP" | "RESET">(
    "LOGIN",
  );
  const [role, setRole] = useState<LoginRole>("TEACHER");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [remember, setRemember] = useState(false);
  const [showPassword, setShowPassword] = useState(false);

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();

  const isDisabled = !email || !password || loading;

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setLoading(true);

    try {
      const payload = {
        email,
        password,
        role,
      };

      await handleLogin({
        data: { email, password, role: role },
      });
      toast.success("login successful");
      router.push(`/${role.toLowerCase()}-dashboard`);
    } catch (err) {
      setError("Invalid email or password");
      toast.error("login failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <motion.div className="bg-bg min-h-screen w-screen flex flex-col md:flex-row">
      {/* LEFT */}
      <div className="flex-1 flex flex-col px-6 py-10 md:p-16">
        <h1 className="text-3xl mb-10">
          <span className="text-primaryBlue font-bold">B</span>
          <span className="font-bold text-white">itwise</span>{" "}
          <span className="text-white">Learn</span>
        </h1>

        <motion.div
          key={step} // 🔥 forces remount when step changes
          variants={slideUp}
          initial="hidden"
          animate="show"
          className="relative w-full md:w-[60%] bg-divBg rounded-3xl p-8"
        >
          {step === "LOGIN" && (
            <>
              <h2 className="text-2xl font-bold mb-4">
                <span className="text-white">Log</span>{" "}
                <span className="text-primaryBlue">in</span>
              </h2>

              {/* ROLE SELECT */}
              <div className="flex gap-2 mb-6 bg-bg p-2 rounded-xl">
                {[
                  { label: "TEACHER", icon: GraduationCap },
                  { label: "INSTITUTION", icon: School },
                  { label: "VENDOR", icon: Handshake },
                ].map(({ label, icon: Icon }) => (
                  <button
                    key={label}
                    type="button"
                    onClick={() => setRole(label as LoginRole)}
                    className={`flex-1 py-2 rounded-lg flex items-center justify-center gap-2 text-sm font-medium transition
                  ${
                    role === label
                      ? "bg-primaryBlue text-white"
                      : "bg-[#3B82F6]/60 text-white hover:bg-[#3B82F6]"
                  }`}
                  >
                    <Icon size={18} />
                    {label.toLowerCase()}
                  </button>
                ))}
              </div>

              {/* FORM */}
              <form onSubmit={handleSubmit} className="space-y-5">
                {/* EMAIL */}
                <div>
                  <label className="text-white text-sm">Email Address</label>
                  <div className="relative mt-1">
                    <Mail className="absolute left-3 top-1/2 -translate-y-1/2 text-white" />
                    <input
                      type="email"
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                      className="w-full pl-11 py-2 rounded-lg bg-bg text-white focus:ring-2 focus:ring-primaryBlue outline-none"
                      placeholder="johndoe@example.com"
                    />
                  </div>
                </div>

                {/* PASSWORD */}
                <div>
                  <label className="text-white text-sm">Password</label>
                  <div className="relative mt-1">
                    <Lock className="absolute left-3 top-1/2 -translate-y-1/2 text-white" />
                    <input
                      type={showPassword ? "text" : "password"}
                      value={password}
                      onChange={(e) => setPassword(e.target.value)}
                      className="w-full pl-11 pr-10 py-2 rounded-lg bg-bg text-white focus:ring-2 focus:ring-primaryBlue outline-none"
                      placeholder="••••••••"
                    />
                    <button
                      type="button"
                      onClick={() => setShowPassword((p) => !p)}
                      className="absolute right-3 top-1/2 -translate-y-1/2 text-white"
                    >
                      {showPassword ? <Eye size={18} /> : <EyeOff size={18} />}
                    </button>
                  </div>
                </div>

                {/* REMEMBER + ERROR */}
                <div className="flex items-center justify-between text-sm">
                  <label className="flex items-center gap-2 text-white">
                    <input
                      type="checkbox"
                      checked={remember}
                      onChange={(e) => setRemember(e.target.checked)}
                      className="accent-primaryBlue"
                    />
                    Remember me
                  </label>

                  <button
                    type="button"
                    onClick={() => setStep("EMAIL")}
                    className="text-neutral-300 hover:text-primaryBlue transition"
                  >
                    Forgot password?
                  </button>
                </div>

                {error && (
                  <p className="text-red-400 text-sm bg-red-400/10 p-2 rounded-lg">
                    {error}
                  </p>
                )}

                {/* SUBMIT */}
                <button
                  type="submit"
                  disabled={isDisabled}
                  className="w-full py-3 rounded-lg bg-primaryBlue text-white font-semibold
              disabled:opacity-50 disabled:cursor-not-allowed transition"
                >
                  {loading ? "Signing in..." : "Log in"}
                </button>
              </form>
            </>
          )}
          {step === "EMAIL" && (
            <ForgotPasswordForm
              role={role}
              email={email}
              setEmail={setEmail}
              onSuccess={() => {
                setStep("OTP");
              }}
              onBack={() => setStep("LOGIN")}
            />
          )}
          {step === "OTP" && (
            <OtpForm
              role={role}
              email={email}
              onVerified={() => setStep("RESET")}
            />
          )}
          {step === "RESET" && (
            <ResetPasswordForm
              role={role}
              onSuccess={() => window.location.href=("/multi-login")}
            />
          )}
        </motion.div>
      </div>

      {/* RIGHT IMAGE */}
      <div className="hidden lg:block lg:w-[38%] relative">
        <Image
          src={MutliLoginIMG}
          alt="Admin login illustration"
          fill
          className="object-cover"
          priority
        />
      </div>
    </motion.div>
  );
}
