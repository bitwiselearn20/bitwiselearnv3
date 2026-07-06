"use client";

import { useState } from "react";
import { Check } from "lucide-react";

export default function CourseMeta({ loading }: { loading: boolean }) {
  const [completed, setCompleted] = useState(false);

  return (
    <div className="flex items-start justify-between flex-wrap gap-4">
      {/* LEFT */}
      <div className="flex flex-col gap-2">
        {loading ? (
          <>
            <div className="h-5 w-48 bg-divBg rounded animate-pulse" />
            <div className="h-4 w-32 bg-divBg rounded animate-pulse" />
          </>
        ) : (
          <>
            <h2 className="text-lg font-semibold text-white">
              Getting Started with HTML
            </h2>

            <div className="flex items-center gap-2">
              <div className="h-7 w-7 rounded-full bg-[#64ACFF] flex items-center justify-center text-black text-xs font-semibold">
                JD
              </div>
              <span className="text-sm text-gray-300">John Doe</span>
            </div>
          </>
        )}
      </div>

      {/* RIGHT */}
      {loading ? (
        <div className="h-5 w-32 bg-divBg rounded animate-pulse hidden sm:block" />
      ) : (
        <div
          onClick={() => setCompleted((p) => !p)}
          className="hidden sm:flex items-center gap-2 cursor-pointer select-none"
        >
          <div
            className={`h-5 w-5 rounded-md border flex items-center justify-center transition-all ${
              completed ? "bg-[#64ACFF] border-[#64ACFF]" : "border-gray-500"
            }`}
          >
            {completed && <Check size={14} className="text-black" />}
          </div>
          <span className="text-sm text-gray-300">Mark as Complete</span>
        </div>
      )}
    </div>
  );
}
