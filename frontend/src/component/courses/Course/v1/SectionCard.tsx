"use client";

import { ChevronDown, Lock } from "lucide-react";
import { useState } from "react";

type Topic = {
  title: string;
  duration: string;
  completed?: boolean;
};

type Section = {
  id: number;
  title: string;
  locked?: boolean;
  topics: Topic[];
};

const colors = {
  primary_Bg: "bg-[#121313]",
  secondary_Bg: "bg-[#1E1E1E]",
  primary_Hero: "bg-[#129274]",
  primary_Hero_Faded: "bg-[rgb(18, 146, 116, 0.24)]",
  secondary_Hero: "bg-[#64ACFF]",
  secondary_Hero_Faded: "bg-[rgb(100, 172, 255, 0.56)]",
  primary_Font: "text-[#FFFFFF]",
  secondary_Font: "text-[#B1AAA6]",
  special_Font: "text-[#64ACFF]",
  accent: "#B1AAA6",
  accent_Faded: "bg-[rgb(177, 170, 166, 0.41)]",
  primary_Icon: "white",
  secondary_Icon: "black",
  special_Icon: "#64ACFF",
};

export default function SectionCard({
  section,
  sectionIndex,
  onToggleTopic,
}: {
  section: Section;
  sectionIndex: number;
  onToggleTopic: (sectionIndex: number, topicIndex: number) => void;
}) {
  const [open, setOpen] = useState(false);

  const completedCount = section.topics.filter((t) => t.completed).length;
  const totalCount = section.topics.length;

  return (
    <div className={`${colors.primary_Bg} rounded-xl p-3`}>
      <button
        disabled={section.locked}
        onClick={() => setOpen(!open)}
        className={`w-full flex items-center justify-between px-1 py-2
          ${section.locked ? "cursor-not-allowed" : ""}`}
      >
        <div className="flex gap-1">
          <span className={`${colors.special_Font} font-bold`}>
            Section {section.id}:
          </span>
          <span className={colors.primary_Font}>{section.title}</span>
        </div>

        <div className="flex items-center gap-3">
          {section.locked ? (
            <Lock size={18} color={colors.special_Icon} />
          ) : (
            <ProgressRing completed={completedCount} total={totalCount} />
          )}
          <ChevronDown
            size={18}
            color={colors.accent}
            className={`transition ${open ? "rotate-180" : ""}`}
          />
        </div>
      </button>

      {open && (
        <div className="mt-3 flex flex-col gap-2">
          {section.topics.map((topic, i) => (
            <TopicRow
              key={i}
              topic={topic}
              onToggle={() => onToggleTopic(sectionIndex, i)}
            />
          ))}
        </div>
      )}
    </div>
  );
}

function TopicRow({ topic, onToggle }: { topic: Topic; onToggle: () => void }) {
  return (
    <div
      className={`${colors.secondary_Bg} px-4 py-3 rounded-lg flex justify-between`}
    >
      <span className="text-white font-mono">{topic.title}</span>

      <div className="flex gap-3 items-center">
        <span className={`${colors.secondary_Font} text-sm`}>
          {topic.duration}
        </span>
        <input
          type="checkbox"
          checked={!!topic.completed}
          onChange={onToggle}
          className="accent-[#64ACFF]"
        />
      </div>
    </div>
  );
}

function ProgressRing({
  completed,
  total,
}: {
  completed: number;
  total: number;
}) {
  const radius = 10;
  const stroke = 3;
  const normalizedRadius = radius - stroke * 0.5;
  const circumference = normalizedRadius * 2 * Math.PI;
  const progress = total === 0 ? 0 : completed / total;
  const strokeDashoffset = circumference - progress * circumference;

  return (
    <svg height="20" width="20">
      <circle
        stroke="#2a2a2a"
        fill="transparent"
        strokeWidth={stroke}
        r={normalizedRadius}
        cx="10"
        cy="10"
      />
      <circle
        stroke="#64ACFF"
        fill="transparent"
        strokeWidth={stroke}
        strokeDasharray={`${circumference} ${circumference}`}
        strokeDashoffset={strokeDashoffset}
        strokeLinecap="round"
        r={normalizedRadius}
        cx="10"
        cy="10"
        style={{ transition: "stroke-dashoffset 0.3s ease" }}
      />
    </svg>
  );
}
