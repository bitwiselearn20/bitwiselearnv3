"use client";
import StudentDashboard from "@/component/Student-Dashboard/StudentDashboard";
import { useStudent } from "@/store/studentStore";
import React from "react";

export default function page() {
  const studentInfo = useStudent().info;
  return <StudentDashboard />;
}
