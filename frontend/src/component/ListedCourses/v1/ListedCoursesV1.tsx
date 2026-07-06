"use client";

import React, { useEffect, useState } from "react";
import { Navbar } from "@/component/general/Navbar";
import Footer from "@/component/general/Footer";
import { getAllListedCourses } from "@/api/listed-courses/listed-courses";
import { CoursesHeader } from "./CoursesHeader";
import toast from "react-hot-toast";

type Course = {
  id: string;
  title: string;
  description: string;
  thumbnail?: string;
  price?: number;
};

function CourseCard({ course }: { course: Course }) {
  return (
    <div className="group rounded-2xl border border-neutral-700 bg-neutral-900/50 overflow-hidden transition-all duration-300 hover:border-neutral-500 hover:shadow-lg hover:shadow-black/30">
      {course.thumbnail ? (
        <div className="h-48 w-full overflow-hidden">
          <img
            src={course.thumbnail}
            alt={course.title}
            className="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
          />
        </div>
      ) : (
        <div className="flex h-48 items-center justify-center bg-neutral-800 text-neutral-500">
          No Thumbnail
        </div>
      )}

      <div className="p-6">
        <h3 className="text-lg font-semibold text-white line-clamp-1">
          {course.title}
        </h3>

        <p className="mt-2 text-sm text-neutral-400 line-clamp-2">
          {course.description}
        </p>

        <div className="mt-4 flex items-center justify-between">
          {course.price ? (
            <span className="text-white font-medium">
              ₹{course.price}
            </span>
          ) : (
            <span className="text-neutral-400 text-sm">
              Free
            </span>
          )}

          <button className="rounded-full bg-white px-4 py-2 text-sm font-medium text-black transition hover:bg-neutral-200">
            View Course
          </button>
        </div>
      </div>
    </div>
  );
}

const ListedCoursesV1 = () => {
  const [courses, setCourses] = useState<Course[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let isMounted = true;

    const fetchCourses = async () => {
      try {
        setLoading(true);

        const data = await getAllListedCourses();

        if (isMounted) {
          setCourses(data || []);
        }
      } catch (error: any) {
        toast.error("Failed to load courses");
      } finally {
        if (isMounted) {
          setLoading(false);
        }
      }
    };

    fetchCourses();

    return () => {
      isMounted = false;
    };
  }, []);

  return (
    <>
      <div className="relative w-full">
        <Navbar />
      </div>

      <CoursesHeader />

      <section className="py-16 sm:py-24">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          {loading ? (
            <div className="text-center text-neutral-400 animate-pulse">
              Loading courses...
            </div>
          ) : courses.length === 0 ? (
            <div className="text-center text-neutral-400">
              No published courses available yet.
            </div>
          ) : (
            <div className="grid gap-8 sm:grid-cols-2 lg:grid-cols-3">
              {courses.map((course) => (
                <CourseCard key={course.id} course={course} />
              ))}
            </div>
          )}
        </div>
      </section>

      <Footer />
    </>
  );
};

export default ListedCoursesV1;