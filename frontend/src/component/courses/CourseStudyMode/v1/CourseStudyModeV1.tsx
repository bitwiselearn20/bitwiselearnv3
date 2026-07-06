"use client";

import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { X } from "lucide-react";
import StudyModeToggle from "../../Course/v1/StudyModeToggle";
import VideoSection from "./VideoSection";
import NotesPanel from "./NotesPanel";
import CourseMeta from "./CourseMeta";
import Link from "next/link";

export default function CourseStudyModeV1({
  studyMode,
  setStudyMode,
}: {
  studyMode: boolean;
  setStudyMode: (v: boolean) => void;
}) {
  const [loading, setLoading] = useState(true);
  const [isOpen, setIsOpen] = useState(false);

  useEffect(() => {
    const t = setTimeout(() => setLoading(false), 1200);
    return () => clearTimeout(t);
  }, []);

  return (
    <div className="relative min-h-screen bg-[#121313] p-4 flex flex-col gap-4 overflow-hidden">
      {/* OVERLAY */}
      {isOpen && (
        <div
          onClick={() => setIsOpen(false)}
          className="absolute inset-0 bg-black/50 z-40"
        />
      )}

      {/* SIDEBAR */}
      <div
        className={`
          absolute top-0 left-0 h-full w-65
          bg-[#121313]
          border-r border-white/10
          z-50
          transform transition-transform duration-300
          ${isOpen ? "translate-x-0" : "-translate-x-full"}
        `}
      >
        <div className="flex items-center justify-between p-4 border-b border-white/10">
          <span className="text-sm text-gray-300 font-mono">
            Course Sidebar
          </span>

          <button
            onClick={() => setIsOpen(false)}
            className="p-1 rounded hover:bg-white/10 cursor-pointer"
          >
            <X size={18} color="red" />
          </button>
        </div>

        <div className="p-4 text-sm text-gray-400">Sidebar content here</div>
      </div>

      {/* TOP BAR */}
      <div className="flex items-center justify-between gap-5 rounded-xl p-1 flex-wrap relative z-30">
        <div className="flex items-center gap-2">
          <button
            onClick={() => setIsOpen(true)}
            className="bg-[#121313] text-white px-3 rounded-md cursor-pointer"
          >
            ☰
          </button>

          {loading ? (
            <div className="h-6 w-48 bg-[#1E1E1E] rounded animate-pulse" />
          ) : (
            <div className="flex items-center">
              <Link
                href="/courses"
                className="text-gray-400 font-mono text-xl mr-2 hover:text-white"
              >
                Course
              </Link>
              <span className="text-blue-400 text-xl">&gt;</span>
              <Link
                href="/courses/1"
                className="text-gray-400 font-mono text-xl ml-2 hover:text-white"
              >
                [Course Name]
              </Link>
            </div>
          )}
        </div>

        <StudyModeToggle
          enabled={studyMode}
          onToggle={() => setStudyMode(!studyMode)}
        />
      </div>

      {/* CONTENT */}
      <div className="flex flex-col lg:flex-row gap-4 flex-1 relative z-10">
        {/* VIDEO */}
        <motion.div
          layout
          className="flex-3 bg-[#1E1E1E] rounded-xl overflow-hidden"
          initial={{ scale: 0.92 }}
          animate={{ scale: 1 }}
          transition={{ duration: 0.4, ease: "easeInOut" }}
        >
          <VideoSection loading={loading} />
        </motion.div>

        {/* NOTES */}
        <motion.div
          layout
          className="flex-[1.3] bg-[#1E1E1E] rounded-xl overflow-hidden"
          initial={{ x: 80, opacity: 0 }}
          animate={{ x: 0, opacity: 1 }}
          transition={{ duration: 0.45, ease: "easeInOut" }}
        >
          <NotesPanel loading={loading} />
        </motion.div>
      </div>

      {/* META */}
      <div className="w-full lg:w-2/3 relative z-10">
        <CourseMeta loading={loading} />
      </div>
    </div>
  );
}
