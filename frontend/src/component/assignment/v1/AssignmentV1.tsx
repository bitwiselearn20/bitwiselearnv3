"use client";

import React, { useState, useEffect } from "react";
import Link from "next/link";
import BG from "./images/BG.png";
import { dummyAssignmentData } from "../DummyData/dummyData";

const colors = {
  primary_Bg: "bg-[#121313]",
  secondary_Bg: "bg-[#1E1E1E]",
  special_Bg: "bg-[#64ACFF]",
  primary_Hero: "bg-[#129274]",
  primary_Font: "text-[#FFFFFF]",
  secondary_Font: "text-[#B1AAA6]",
  special_Font: "text-[#64ACFF]",
};

export default function AssignmentV1({
  assignmentId,
}: {
  assignmentId: string;
}) {
  const [loading, setLoading] = useState(true);

  const data = dummyAssignmentData[assignmentId];

  if (!data) {
    return <div>Assignment not found</div>;
  }

  useEffect(() => {
    const timer = setTimeout(() => setLoading(false), 500);
    return () => clearTimeout(timer);
  }, []);

  return (
    <div
      className={`min-h-screen w-full bg-cover bg-center bg-no-repeat flex justify-center items-center font-mono ${colors.primary_Font}`}
      style={{ backgroundImage: `url(${BG.src})` }}
    >
      {loading ? (
        // Skeleton Content
        <div
          className={`${colors.secondary_Bg} p-4 rounded-lg flex flex-col justify-center items-center gap-8 w-96 animate-pulse`}
        >
          <div className="h-8 w-3/4 bg-gray-700/50 rounded-md"></div>

          <div className="flex items-center justify-between gap-6 w-full px-2">
            <div className="h-5 w-1/3 bg-gray-700/50 rounded-md"></div>
            <div className="h-5 w-1/3 bg-gray-700/50 rounded-md"></div>
          </div>

          <div className="h-10 w-32 bg-gray-700/50 rounded-md"></div>
        </div>
      ) : (
        // Real Content
        <div
          className={`${colors.secondary_Bg} p-4 rounded-lg flex flex-col justify-center items-center gap-8 w-96`}
        >
          <h1 className={`text-2xl font-bold`}>
            <span className={`${colors.special_Font}`}>
              {data.name.split(" ")[0]}{" "}
            </span>
            <span>{data.name.slice(data.name.indexOf(" "))}</span>
          </h1>
          <div className="flex items-center justify-between gap-6 w-full px-2">
            <div>
              <span className={`${colors.special_Font}`}>Duration: </span>{" "}
              <span>{data.durationInMinutes}</span> <span> mins</span>
            </div>
            <div>
              <span className={`${colors.special_Font}`}>Questions: </span>
              <span>{data.totalQuestions}</span>
            </div>
          </div>
          <Link
            href={`${data.id}/attempt`}
            className={`${colors.special_Bg} px-8 py-2 rounded-md hover:scale-105 hover:opacity-90 text-white font-semibold`}
          >
            Start Now
          </Link>
        </div>
      )}
    </div>
  );
}
