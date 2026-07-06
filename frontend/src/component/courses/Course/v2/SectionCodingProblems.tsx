"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { CheckCircle2, Circle, Code2 } from "lucide-react";
import { getStudentSectionProblems } from "@/api/courses/section/get-student-section-problems";
import { getColors } from "@/component/general/(Color Manager)/useColors";

type CodingProblem = {
  id: string;
  name: string;
  difficulty: string;
  completed: boolean;
};

// Renders the coding problems attached to a section, with the student's
// completed / not-completed status. Self-contained: pass a sectionId.
export default function SectionCodingProblems({ sectionId }: { sectionId: string }) {
  const Colors = getColors();
  const [problems, setProblems] = useState<CodingProblem[]>([]);
  const [completed, setCompleted] = useState(0);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!sectionId) {
      setLoading(false);
      return;
    }
    let active = true;
    setLoading(true);
    getStudentSectionProblems(sectionId)
      .then((data) => {
        if (!active) return;
        const list = Array.isArray(data?.problems) ? data.problems : [];
        setProblems(list);
        setCompleted(data?.completed ?? 0);
      })
      .catch(() => active && setProblems([]))
      .finally(() => active && setLoading(false));
    return () => {
      active = false;
    };
  }, [sectionId]);

  if (loading) {
    return (
      <p className={`text-sm ${Colors.text.secondary} px-1 py-2`}>
        Loading coding problems…
      </p>
    );
  }

  if (problems.length === 0) return null; // nothing attached -> show nothing

  return (
    <div className="mt-4 space-y-2">
      <div className="flex items-center justify-between px-1">
        <div className={`flex items-center gap-2 text-sm font-semibold ${Colors.text.primary}`}>
          <Code2 size={16} className={Colors.text.special} />
          Coding Problems
        </div>
        <span className={`text-xs ${Colors.text.secondary}`}>
          {completed}/{problems.length} solved
        </span>
      </div>

      {problems.map((p) => (
        <Link
          key={p.id}
          href={`/problems/${p.id}`}
          className={`
            flex items-center justify-between gap-3 rounded-lg px-4 py-2.5
            ${Colors.border.defaultThin} ${Colors.background.primary}
            ${Colors.hover.special} transition
          `}
        >
          <span className="flex items-center gap-2">
            {p.completed ? (
              <CheckCircle2 size={16} className="text-green-400" />
            ) : (
              <Circle size={16} className={Colors.text.secondary} />
            )}
            <span className={`text-sm ${Colors.text.primary}`}>{p.name}</span>
          </span>
          <span className={`text-xs ${Colors.text.secondary}`}>{p.difficulty}</span>
        </Link>
      ))}
    </div>
  );
}
