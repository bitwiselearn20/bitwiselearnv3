"use client";

import V1HomeNav from "@/component/Home/V1/V1HomeNav";
import { Klee_One } from "next/font/google";
import { CheckCircle2, Mail, Phone } from "lucide-react";
import { motion, Variants } from "framer-motion";
import Footer from "@/component/general/Footer";

const kleeOne = Klee_One({
  subsets: ["latin"],
  weight: ["400", "600"],
});

const container = {
  hidden: {},
  show: {
    transition: {
      staggerChildren: 0.15,
    },
  },
};

const fadeUp: Variants = {
  hidden: { opacity: 0, y: 30 },
  show: {
    opacity: 1,
    y: 0,
    transition: { duration: 0.7, ease: "easeOut" },
  },
};

export default function ContactV1() {
  async function getFormData(formData: FormData) {
    const name = formData.get("name");
    const email = formData.get("email");
    const message = formData.get("message");
  }

  return (
    <div className={kleeOne.className}>
      <V1HomeNav />

      <motion.div
        variants={container}
        initial="hidden"
        animate="show"
        className="text-white max-w-7xl mx-auto px-6 pt-24"
      >
        <motion.h1
          variants={fadeUp}
          className="text-4xl md:text-6xl text-center"
        >
          Get in touch with us
        </motion.h1>

        <motion.p
          variants={fadeUp}
          className="text-center text-white/70 mt-4 text-sm md:text-base"
        >
          Fill out form below or schedule a meeting with us at your{" "}
          <br className="hidden sm:block" />
          convenience
        </motion.p>

        <motion.div
          variants={container}
          className="mt-10 grid grid-cols-1 lg:grid-cols-2 gap-16 items-start"
        >
          {/* FORM */}
          <motion.div variants={fadeUp} className="w-full max-w-md">
            <form action={getFormData} className="space-y-6">
              {[
                {
                  label: "Name",
                  name: "name",
                  type: "text",
                  placeholder: "Enter Your Name",
                },
                {
                  label: "Email",
                  name: "email",
                  type: "email",
                  placeholder: "Enter Your Email",
                },
              ].map((field) => (
                <div key={field.name} className="flex flex-col gap-2">
                  <label className="text-sm text-white/80">{field.label}</label>
                  <input
                    type={field.type}
                    name={field.name}
                    placeholder={field.placeholder}
                    className="px-5 py-3 rounded-xl bg-white/5 text-white placeholder:text-white/40
                      border border-white/15 backdrop-blur-md shadow-inner
                      focus:outline-none focus:border-blue-400
                      focus:ring-2 focus:ring-blue-400/30 transition"
                  />
                </div>
              ))}

              <div className="flex flex-col gap-2">
                <label className="text-sm text-white/80">Message</label>
                <textarea
                  name="message"
                  rows={4}
                  placeholder="Enter Your Message"
                  className="px-5 py-3 rounded-xl bg-white/5 text-white placeholder:text-white/40
                    border border-white/15 backdrop-blur-md shadow-inner resize-none
                    focus:outline-none focus:border-blue-400
                    focus:ring-2 focus:ring-blue-400/30 transition"
                />
              </div>

              <div className="flex justify-center pt-4">
                <button
                  type="submit"
                  className="px-10 mb-10 py-3 rounded-xl bg-primaryBlue text-black font-semibold
                    hover:bg-blue-500 transition shadow-[0_8px_25px_rgba(59,130,246,0.4)]"
                >
                  Connect with us
                </button>
              </div>
            </form>
          </motion.div>

          {/* SERVICES */}
          <motion.div variants={fadeUp} className="space-y-6">
            <h2 className="text-2xl">With our services you can</h2>

            <div className="space-y-3">
              {[
                {
                  text: "Learn with structured, goal-driven paths",
                  color: "text-red-300",
                },
                {
                  text: "Track progress and stay consistent",
                  color: "text-blue-300",
                },
                {
                  text: "Practice with tasks, quizzes, and checkpoints",
                  color: "text-green-300",
                },
                {
                  text: "Get insights on what to learn next",
                  color: "text-purple-300",
                },
              ].map((item, i) => (
                <div key={i} className="flex items-center gap-3 text-white/80">
                  <CheckCircle2 className={item.color} size={18} />
                  <span>{item.text}</span>
                </div>
              ))}
            </div>

            <div className="h-px bg-white/20 my-8" />

            <p className="text-center text-white/70">
              You can also contact us via
            </p>

            <div className="flex flex-col sm:flex-row gap-6 justify-center text-white/80">
              <div className="flex items-center gap-2">
                <Mail size={18} />
                <span>johndoe@example.com</span>
              </div>
              <div className="flex items-center gap-2">
                <Phone size={18} />
                <span>+91 987 8762 876</span>
              </div>
            </div>
          </motion.div>
        </motion.div>
      </motion.div>
      <Footer />
    </div>
  );
}
