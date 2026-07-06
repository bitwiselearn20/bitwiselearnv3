"use client";

import { useEffect, useState } from "react";
import Study from "./Study";
import Assignments from "./Assignments";
import StudyModeToggle from "./StudyModeToggle";
import { motion } from "framer-motion";
import { AnimatePresence } from "framer-motion";

type ActiveTab = "study" | "assignments";

export default function RightSection({
  studyMode,
  setStudyMode,
}: {
  studyMode: boolean;
  setStudyMode: (v: boolean) => void;
}) {
  const [activeTab, setActiveTab] = useState<ActiveTab>("study");
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const timer = setTimeout(() => setLoading(false), 1200);
    return () => clearTimeout(timer);
  }, []);

  if (loading) {
    return (
      <div className="bg-[#1E1E1E] p-4 w-full h-full rounded-xl flex flex-col gap-6 animate-pulse">
        <div className="flex justify-center gap-6">
          <div className="h-10 w-32 bg-[#2a2a2a] rounded-md" />
          <div className="h-10 w-32 bg-[#2a2a2a] rounded-md" />
        </div>

        <div className="space-y-4">
          <div className="h-6 w-1/3 bg-[#2a2a2a] rounded" />
          <div className="h-40 w-full bg-[#2a2a2a] rounded-lg" />
          <div className="h-40 w-full bg-[#2a2a2a] rounded-lg" />
        </div>
      </div>
    );
  }

  return (
    <div className="bg-[#1E1E1E] p-4 w-full h-full rounded-xl flex flex-col gap-6 text-white font-mono">
      <div className="flex items-center">
        <div className="w-fit flex items-center mx-auto py-2 px-6 rounded-md gap-6">
          <motion.button
            onClick={() => setActiveTab("study")}
            animate={{
              backgroundColor: activeTab === "study" ? "#64ACFF" : "#4b5563",
              color: activeTab === "study" ? "#000000" : "#ffffff",
            }}
            transition={{ duration: 0.3, ease: "easeInOut" }}
            className={`text-xl px-4 py-2 rounded-md ${
              activeTab === "study" ? "cursor-not-allowed" : "cursor-pointer"
            }`}
          >
            Study
          </motion.button>

          <div className="border border-gray-600 h-10" />

          <motion.button
            onClick={() => setActiveTab("assignments")}
            animate={{
              backgroundColor:
                activeTab === "assignments" ? "#64ACFF" : "#4b5563",
              color: activeTab === "assignments" ? "#000000" : "#ffffff",
            }}
            transition={{ duration: 0.3, ease: "easeInOut" }}
            className={`text-xl rounded-md px-4 py-2 ${
              activeTab === "assignments"
                ? "cursor-not-allowed"
                : "cursor-pointer"
            }`}
          >
            Assignments
          </motion.button>
        </div>

        <StudyModeToggle
          enabled={studyMode}
          onToggle={() => setStudyMode(true)}
        />
      </div>

      <div className="flex-1 w-full overflow-hidden">
        <AnimatePresence mode="wait">
          {activeTab === "study" && (
            <motion.div
              key="study"
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -10 }}
              transition={{ duration: 0.25 }}
              className="h-full"
            >
              <Study />
            </motion.div>
          )}

          {activeTab === "assignments" && (
            <motion.div
              key="assignments"
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -10 }}
              transition={{ duration: 0.25 }}
              className="h-full"
            >
              <Assignments />
            </motion.div>
          )}
        </AnimatePresence>
      </div>
    </div>
  );
}
