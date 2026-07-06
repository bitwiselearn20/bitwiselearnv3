"use client";

import { Menu, X } from "lucide-react";
import { useEffect, useState } from "react";
import SectionCard from "./SectionCard";
import Breadcrumbs from "./BreadCrumbs";

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

const initialSections: Section[] = [
  {
    id: 1,
    title: "Learn the basics",
    topics: [
      { title: "Topic - 1", duration: "0:56", completed: true },
      { title: "Topic - 2", duration: "1:23" },
      { title: "Topic - 3", duration: "1:45" },
      { title: "Topic - 4", duration: "2:00" },
    ],
  },
  {
    id: 2,
    title: "Learn the basics - 2",
    locked: true,
    topics: [
      { title: "Topic - 1", duration: "1:10" },
      { title: "Topic - 2", duration: "2:00" },
    ],
  },
  {
    id: 3,
    title: "Learn the basics - 3",
    locked: true,
    topics: [],
  },
];

export default function LeftSection() {
  const [isOpen, setIsOpen] = useState(false);
  const [sections, setSections] = useState(initialSections);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const timer = setTimeout(() => setLoading(false), 800);
    return () => clearTimeout(timer);
  }, []);

  const updateTopic = (sectionIndex: number, topicIndex: number) => {
    setSections((prev) => {
      const updated = structuredClone(prev);
      const topic = updated[sectionIndex].topics[topicIndex];
      topic.completed = !topic.completed;

      const allCompleted = updated[sectionIndex].topics.every(
        (t) => t.completed,
      );

      if (allCompleted && updated[sectionIndex + 1]) {
        updated[sectionIndex + 1].locked = false;
      }

      return updated;
    });
  };

  return (
    <div className="relative bg-[#1E1E1E] p-4 w-full h-full rounded-xl flex flex-col gap-3 overflow-hidden">
      {/* HEADER */}
      <div className="flex gap-10 items-center">
        <button
          onClick={() => setIsOpen(true)}
          className="p-2 rounded-sm cursor-pointer"
        >
          <Menu size={25} color="white" />
        </button>

        <Breadcrumbs />
      </div>

      {/* MAIN CONTENT */}
      <div className="flex-1 w-full flex flex-col gap-4">
        {loading ? (
          <>
            {[1, 2, 3].map((i) => (
              <div
                key={i}
                className="animate-pulse bg-[#121313] rounded-lg p-4 space-y-3"
              >
                <div className="h-4 w-1/3 bg-white/10 rounded" />
                <div className="h-3 w-full bg-white/10 rounded" />
                <div className="h-3 w-5/6 bg-white/10 rounded" />
                <div className="h-3 w-2/3 bg-white/10 rounded" />
              </div>
            ))}
          </>
        ) : (
          sections.map((section, i) => (
            <SectionCard
              key={section.id}
              section={section}
              sectionIndex={i}
              onToggleTopic={updateTopic}
            />
          ))
        )}
      </div>

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
    </div>
  );
}
