"use client";

import StudentLoginIMG from "../v1/StudentLoginIMG.png";
import Image from "next/image";
import { Mail, Lock, Eye, EyeOff } from "lucide-react";
import { useState, useEffect } from "react";
import handleLogin from "@/api/handleLogin";
import { useRouter } from "next/navigation";
import { motion, type Variants } from "framer-motion";

import ResetPasswordForm from "@/component/auth/ResetPasswordForm";
import ForgotPasswordForm from "@/component/auth/ForgotPasswordForm";
import OtpForm from "@/component/auth/OtpForm";
import { toast } from "react-hot-toast";

/* ================= ANIMATION VARIANTS ================= */

const pageFade: Variants = {
  hidden: { opacity: 0 },
  show: { opacity: 1, transition: { duration: 0.6, ease: "easeOut" as const } },
};

const slideUp: Variants = {
  hidden: { opacity: 0, y: 30 },
  show: {
    opacity: 1,
    y: 0,
    transition: { duration: 0.7, ease: "easeOut" as const },
  },
};

const stagger: Variants = {
  hidden: {},
  show: { transition: { staggerChildren: 0.15 } },
};

const imageReveal: Variants = {
  hidden: { opacity: 0, scale: 1.05 },
  show: {
    opacity: 1,
    scale: 1,
    transition: { duration: 0.9, ease: "easeOut" as const },
  },
};

/* ================= TYPEWRITER ================= */
const ROLE = "STUDENT" as const;

function WelcomeTypewriter() {
  const fullText = "Learn a little today. Grow a lot tomorrow.";
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

  const renderText = () =>
    text.split(/(Learn|Grow)/g).map((part, i) =>
      part === "Learn" || part === "Grow" ? (
        <span key={i} className="text-primaryBlue">
          {part}
        </span>
      ) : (
        <span key={i}>{part}</span>
      ),
    );

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

export default function StudentLoginV1() {
  const [step, setStep] = useState<"LOGIN" | "EMAIL" | "OTP" | "RESET">(
    "LOGIN",
  );
  const [showPassword, setShowPassword] = useState(false);

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const router = useRouter();

  async function fetchLoginData(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();

    try {
      setLoading(true);
      await handleLogin({
        data: { email, password, role: "STUDENT" },
      });
      toast.success("Login successful");
      router.push("/dashboard");
    } catch (err) {
      // console.error("Login failed", err);
      toast.error("login failed");
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
      {/* LEFT */}
      <motion.div
        variants={stagger}
        initial="hidden"
        animate="show"
        className="flex-1 flex flex-col px-6 py-10 md:p-16"
      >
        <motion.h1
          variants={slideUp}
          className="text-3xl text-center md:text-left"
        >
          <span className="text-primaryBlue font-bold">B</span>
          <span className="font-bold text-white">itwise</span>{" "}
          <span className="text-white">Learn</span>
        </motion.h1>

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
              onSuccess={() => window.location.reload()}
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
          src={StudentLoginIMG}
          alt="Login Illustration"
          fill
          className="object-cover"
        />
      </motion.div>
    </motion.div>
  );
}
