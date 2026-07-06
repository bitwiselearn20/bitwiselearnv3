"use client";

import { AnimatePresence, motion } from "framer-motion";
import { useRouter, useSearchParams } from "next/navigation";
import LeftSection from "@/component/courses/Course/v1/LeftSection";
import RightSection from "@/component/courses/Course/v1/RightSection";
import CourseStudyModeV1 from "../../CourseStudyMode/v1/CourseStudyModeV1";
import { useState, useEffect, useRef } from "react";

export default function CourseV1() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const lastTriggerRef = useRef(0);

  const studyMode = searchParams.get("mode") === "study";

  const setStudyMode = (value: boolean) => {
    router.replace(value ? "/courses/1?mode=study" : "/courses/1");
  };

  const isTouchDevice = () =>
    typeof window !== "undefined" &&
    ("ontouchstart" in window || navigator.maxTouchPoints > 0);

  useEffect(() => {
    if (isTouchDevice()) return;

    const handleKeyDown = (e: KeyboardEvent) => {
      const now = Date.now();
      if (now - lastTriggerRef.current < 400) return;
      lastTriggerRef.current = now;

      const tag = (e.target as HTMLElement).tagName;
      if (["INPUT", "TEXTAREA"].includes(tag)) return;

      if (e.key.toLowerCase() === "s") {
        e.preventDefault();
        setStudyMode(!studyMode);
      }
    };

    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [studyMode]);

  return (
    <div className="min-h-screen bg-[#121313] p-4">
      <AnimatePresence mode="wait">
        {!studyMode ? (
          <motion.div
            key="normal"
            initial={{ opacity: 0, scale: 0.98 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 1.02 }}
            transition={{ duration: 0.35, ease: "easeInOut" }}
            className="
              grid grid-cols-1 lg:grid-cols-3 gap-4
              min-h-[calc(100vh-2rem)]
            "
          >
            {/* LEFT */}
            <div className="lg:col-span-1 h-full">
              <LeftSection />
            </div>

            {/* RIGHT – THIS MUST SCROLL */}
            <div className="lg:col-span-2 h-full overflow-y-auto">
              <RightSection studyMode={studyMode} setStudyMode={setStudyMode} />
            </div>
          </motion.div>
        ) : (
          <motion.div
            key="study"
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 1.05 }}
            transition={{ duration: 0.4, ease: "easeInOut" }}
            className="min-h-[calc(100vh-2rem)]"
          >
            <CourseStudyModeV1
              studyMode={studyMode}
              setStudyMode={setStudyMode}
            />
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
