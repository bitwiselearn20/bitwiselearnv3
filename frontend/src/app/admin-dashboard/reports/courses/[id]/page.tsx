"use client";
import CourseEnrollment from "@/component/report/courseEnrollment/CourseEnrollment";
import { useParams } from "next/navigation";
import React from "react";

function page() {
  const params = useParams();
  return <CourseEnrollment courseId={params.id as string} />;
}

export default page;
