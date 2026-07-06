"use client";
import V1HomeNav from "./V1HomeNav";
import Image from "next/image";
import bgIMG from "./V1bgIMG.png";
import studyIMG from "./study.png";
import { Klee_One } from "next/font/google";
import Chart from "./Chart.png";
import { motion, Variants } from "framer-motion";
import Link from "next/link";

import {
  Lightbulb,
  Play,
  ClipboardCheck,
  CheckCircle2,
  HelpCircle,
  Brain,
  Code2,
  Users,
  FolderKanban,
  Star,
  GraduationCap,
} from "lucide-react";
import Footer from "@/component/general/Footer";
import Contributors from "@/component/Contributors/Contributors";

const kleeOne = Klee_One({
  subsets: ["latin"],
  weight: ["400", "600"],
});

function V1Home() {
  const containerVariants = {
    hidden: {},
    show: {
      transition: {
        staggerChildren: 0.18,
      },
    },
  };

  const cardVariants: Variants = {
    hidden: {
      opacity: 0,
      y: 40,
      rotateX: -8,
      scale: 0.96,
    },
    show: {
      opacity: 1,
      y: 0,
      rotateX: 0,
      scale: 1,
      transition: {
        duration: 0.7,
        ease: "easeOut",
      },
    },
  };

  return (
    <div className="relative text-white">
      {/* ================= FIXED BACKGROUND ================= */}
      <div className="fixed inset-0 -z-10">
        <Image
          src={bgIMG}
          alt="Background"
          fill
          priority
          className="object-cover"
        />
      </div>

      {/* ================= CONTENT ================= */}
      <div className="relative min-h-screen overflow-x-hidden">
        {/* Navbar */}
        <V1HomeNav />

        {/* Main wrapper */}
        <main className="pt-20 mt-20">
          {/* ================= HERO ================= */}
          <section
            className={`${kleeOne.className}
            flex flex-col items-center justify-center
            text-center
            -translate-y-6`}
          >
            <h1 className="text-6xl md:text-7xl tracking-widest font-normal">
              Learn. Code. Grow.
            </h1>

            <p className="mt-4 max-w-xl text-sm md:text-base text-white/80 leading-relaxed font-semibold">
              Designed to help you learn better, stay on track
              <br />
              and grow with purpose.
            </p>
          </section>

          {/* ================= CARDS ================= */}
          <section className={`${kleeOne.className} mt-16`}>
            <motion.div
              variants={containerVariants}
              initial="hidden"
              animate="show"
              className="max-w-7xl mx-auto px-6 grid grid-cols-3 gap-8 items-start"
            >
              {/* ================= CARD 1 ================= */}
              <motion.div
                variants={cardVariants}
                className="relative rounded-2xl px-6 py-10 bg-white/6 border 
                border-white/12 backdrop-blur-xl shadow-[0_20px_40px_rgba(0,0,0,0.45)]
                transition-all duration-500 ease-out
                hover:-translate-y-3
                hover:rotate-y-6
                hover:-rotate-x-3
                transform-3d
                perspective-[1000px]"
              >
                <div className="absolute inset-0 rounded-2xl bg-linear-to-b from-white/[0.14] to-transparent pointer-events-none" />

                <h3 className="relative text-2xl font-medium tracking-wide text-white/90">
                  Your Learning Progress
                </h3>

                <div className="relative mt-[-20] flex items-end justify-between">
                  <Image src={Chart} alt="Chart.png" height={130} />
                  <div className="text-right">
                    <p className="text-4xl font-semibold">75%</p>
                    <p className="text-md text-white/60">Complete</p>
                  </div>
                </div>

                <p className="relative mt-4 text-lg text-white/70 leading-relaxed text-center">
                  Stay consistent and
                  <br />
                  see your skills grow over time.
                </p>
              </motion.div>

              {/* ================= CARD 2 ================= */}
              <motion.div
                variants={cardVariants}
                className="relative flex justify-center transition-all duration-500 ease-out
                hover:-translate-y-4
                hover:-rotate-x-6
                transform-3d
                perspective-[1000px]"
              >
                <div
                  className="relative w-full max-w-90 h-90
                  rounded-2xl px-6 py-6
                  bg-white/6 border border-white/12
                  backdrop-blur-xl
                  shadow-[0_20px_40px_rgba(0,0,0,0.45)]
                  overflow-hidden"
                >
                  <div className="absolute inset-0 rounded-2xl bg-linear-to-b from-white/[0.14] to-transparent pointer-events-none" />

                  <h3 className="relative text-2xl font-medium tracking-wide text-white/90 text-center">
                    Your Learning Path
                  </h3>

                  <div className="relative mt-6 h-57.5">
                    <svg
                      className="absolute inset-0"
                      viewBox="0 0 360 230"
                      fill="none"
                    >
                      <path
                        d="M140 45 H180 V45 H220"
                        stroke="white"
                        strokeWidth="2"
                        strokeLinecap="round"
                      />
                      <path
                        d="
                          M240 92
                          V150
                          Q240 160 230 160
                          H110
                          Q100 160 100 170
                          V174
                        "
                        stroke="white"
                        strokeWidth="2"
                        strokeLinecap="round"
                        fill="none"
                      />
                    </svg>

                    <div
                      className="absolute left-0 top-0 w-30 h-20
                      rounded-xl border-2 border-white
                      flex flex-col items-center justify-center gap-1 bg-black/20"
                    >
                      <Play size={22} fill="white" />
                      <span className="text-base">Lesson</span>
                    </div>

                    <div
                      className="absolute right-0 top-0 w-30 h-20
                      rounded-xl border-2 border-white
                      flex flex-col items-center justify-center gap-1 bg-black/20"
                    >
                      <ClipboardCheck size={22} />
                      <span className="text-base">Task</span>
                    </div>

                    <div
                      className="absolute left-0 bottom-0 w-30 h-20
                      rounded-xl border-2 border-white
                      flex flex-col items-center justify-center gap-1 bg-black/20"
                    >
                      <Lightbulb size={22} />
                      <span className="text-base">Quiz</span>
                    </div>

                    <p
                      className="absolute left-35 bottom-2.5
                      text-sm text-white/80 leading-relaxed max-w-45"
                    >
                      Follow a structured path with lessons, tasks and
                      checkpoints
                    </p>
                  </div>
                </div>
              </motion.div>

              {/* ================= CARD 3 ================= */}
              <motion.div
                variants={cardVariants}
                className="relative rounded-2xl px-6 py-6 bg-white/6 border border-white/12 
                backdrop-blur-xl shadow-[0_20px_40px_rgba(0,0,0,0.45)]
                transition-all duration-500 ease-out
                hover:-translate-y-3
                hover:-rotate-y-6
                hover:-rotate-x-3
                transform-3d
                perspective-[1000px]"
              >
                <div className="absolute inset-0 rounded-2xl bg-linear-to-b from-white/[0.14] to-transparent pointer-events-none" />

                <h3 className="relative text-2xl font-medium tracking-wide text-white/90 flex items-center gap-2">
                  Learning Insights
                  <Lightbulb size={35} />
                </h3>

                <h1 className="text-lg">What should I do next?</h1>

                <ul className="relative mt-4 space-y-3 text-sm">
                  <li className="flex items-center gap-2 text-white/80">
                    <CheckCircle2 size={16} className="text-blue-400" />
                    Which topic needs more practice?
                  </li>

                  <li className="flex items-center gap-2 text-white/80">
                    <HelpCircle size={16} className="text-red-400" />
                    Where am I stuck?
                  </li>

                  <li className="flex items-center gap-2 text-white/80">
                    <CheckCircle2 size={16} className="text-green-400" />
                    How consistent am I?
                  </li>
                </ul>

                <p className="relative mt-4 text-xl text-white/70 leading-relaxed">
                  Get clarity from your learning
                  <br />
                  data and improve faster.
                </p>
              </motion.div>
            </motion.div>
          </section>
        </main>

        {/* ================= SMOOTH TRANSITION ================= */}
        <div className="h-68 bg-linear-to-b from-transparent via-black/60 to-black/90" />

        {/* ================= DEMO SCROLL SECTION ================= */}
        <section
          className={`${kleeOne.className} relative min-h-screen bg-black/90 flex items-center justify-center px-16`}
        >
          <div className="max-w-7xl w-full grid grid-cols-2 gap-16 items-center">
            {/* ================= LEFT SECTION ================= */}
            <motion.div
              initial={{ opacity: 0, y: 40 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.8, ease: "easeOut" }}
              viewport={{ once: true }}
            >
              <div className="text-4xl font-extrabold mb-6 leading-tight">
                <h1>Bridging the Gap between</h1>
                <span className="text-primaryBlue">Learning</span>{" "}
                <span>& Your Tech Career</span>
              </div>

              <div className="text-white/70 text-base leading-relaxed max-w-xl mb-8">
                <p>
                  We transform aspiring coders into job-ready developers.
                  Bitwise Learn guides you through the journey from beginner to
                  industry professional with well-planned courses, hands-on
                  practice, and personalized mentorship.
                </p>
              </div>

              {/* FEATURES */}
              <div className="space-y-4 mb-8">
                {[
                  {
                    icon: <Brain size={22} />,
                    title: "Smart Learning Paths",
                    desc: "AI-assisted, structured progress from basics to mastery",
                  },
                  {
                    icon: <Code2 size={22} />,
                    title: "Hands-On Coding",
                    desc: "Practice concepts with real-world coding tasks",
                  },
                  {
                    icon: <Users size={22} />,
                    title: "1:1 Mentorship",
                    desc: "Get guidance and feedback from experienced developers",
                  },
                  {
                    icon: <FolderKanban size={22} />,
                    title: "Project-Based Learning",
                    desc: "Build portfolio-ready projects as you learn",
                  },
                ].map((item, i) => (
                  <motion.div
                    key={i}
                    initial={{ opacity: 0, x: -30 }}
                    whileInView={{ opacity: 1, x: 0 }}
                    transition={{ delay: i * 0.1, duration: 0.6 }}
                    viewport={{ once: true }}
                    className="flex gap-4 items-start p-4 rounded-xl bg-white/5 border border-white/10 backdrop-blur-md"
                  >
                    <div className="text-primaryBlue">{item.icon}</div>
                    <div>
                      <h4 className="font-semibold text-white">{item.title}</h4>
                      <p className="text-sm text-white/60">{item.desc}</p>
                    </div>
                  </motion.div>
                ))}
              </div>

              {/* CTA */}
              <Link href="/about">
                <motion.button
                  whileHover={{ scale: 1.05 }}
                  whileTap={{ scale: 0.95 }}
                  className="px-6 py-3 rounded-full bg-primaryBlue text-black font-semibold"
                >
                  Discover Our Mission →
                </motion.button>
              </Link>
            </motion.div>

            {/* ================= RIGHT SECTION ================= */}
            <motion.div
              initial={{ opacity: 0, y: 60 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.9, ease: "easeOut" }}
              viewport={{ once: true }}
              className="relative flex items-center justify-center"
            >
              <div className="relative w-105 h-70 rounded-2xl bg-linear-to-br from-white/10 to-white/5 border border-white/10 backdrop-blur-xl shadow-[0_30px_80px_rgba(0,0,0,0.6)] flex items-center justify-center">
                <Image src={studyIMG} alt="study.png" className="rounded-2xl" />
              </div>

              {/* FLOATING STATS */}
              <motion.div
                initial={{ opacity: 0, scale: 0.8 }}
                whileInView={{ opacity: 1, scale: 1 }}
                transition={{ delay: 0.4, duration: 0.6 }}
                className="absolute -right-10 -top-6 px-4 py-3 rounded-xl bg-black/70 border border-white/15 backdrop-blur-lg"
              >
                <div className="flex items-center gap-2">
                  <Star className="text-yellow-400" size={18} />
                  <span className="font-semibold">4.8</span>
                  <span className="text-sm text-white/60">(300+ reviews)</span>
                </div>
              </motion.div>

              <motion.div
                initial={{ opacity: 0, scale: 0.8 }}
                whileInView={{ opacity: 1, scale: 1 }}
                transition={{ delay: 0.6, duration: 0.6 }}
                className="absolute -left-12 bottom-6 px-5 py-4 rounded-xl bg-black/70 border border-white/15 backdrop-blur-lg space-y-1"
              >
                <p className="text-lg font-semibold">12,000+</p>
                <p className="text-sm text-white/60">Learners Growing Daily</p>
              </motion.div>
            </motion.div>
          </div>
        </section>

        <section className="bg-black/90">
          <Contributors />
        </section>

        {/* ================= BLACK → BLUE FADE ================= */}
        <div className="h-64 bg-linear-to-b from-black via-zinc/900 to-zinc-800" />

        {/* ================= CTA SECTION ================= */}
        <section className="relative min-h-[70vh] bg-linear-to-b from-zinc-800 to-black flex items-center justify-center overflow-hidden">
          {/* subtle background glow */}
          <div className="absolute inset-0 bg-[radial-gradient(circle_at_center,rgba(255,255,255,0.08),transparent_60%)] pointer-events-none" />

          <motion.div
            initial={{ opacity: 0, y: 60 }}
            whileInView={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.9, ease: "easeOut" }}
            viewport={{ once: true }}
            className="relative z-10 text-center max-w-3xl px-6"
          >
            {/* Heading */}
            <motion.h1
              initial={{ opacity: 0, y: 30 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.2, duration: 0.8 }}
              viewport={{ once: true }}
              className="text-4xl md:text-5xl font-extrabold text-white mb-4"
            >
              Ready to Launch Your
              <br />
              Tech Career?
            </motion.h1>

            {/* Subtext */}
            <motion.p
              initial={{ opacity: 0 }}
              whileInView={{ opacity: 1 }}
              transition={{ delay: 0.4, duration: 0.8 }}
              viewport={{ once: true }}
              className="text-white/80 text-sm md:text-base mb-10"
            >
              Join thousands of students who have transformed their lives with
              Bitwise Learn.
            </motion.p>

            {/* Buttons */}
            <motion.div
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.6, duration: 0.7 }}
              viewport={{ once: true }}
              className="flex items-center justify-center gap-4"
            >
              <Link href="/student-login">
                <motion.button
                  whileHover={{ scale: 1.07 }}
                  whileTap={{ scale: 0.95 }}
                  className="px-6 py-3 rounded-full bg-white text-blue-700 font-semibold shadow-lg"
                >
                  Get Started Now →
                </motion.button>
              </Link>

              <Link href="/contact">
                <motion.button
                  whileHover={{ scale: 1.05 }}
                  whileTap={{ scale: 0.95 }}
                  className="px-6 py-3 rounded-full border border-white/50 text-white font-medium backdrop-blur-sm"
                >
                  Contact Us
                </motion.button>
              </Link>
            </motion.div>
          </motion.div>
        </section>
      </div>

      <Footer />
    </div>
  );
}

export default V1Home;
