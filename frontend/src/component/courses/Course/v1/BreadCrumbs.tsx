import Link from "next/link";

export default function Breadcrumbs() {
  return (
    <div className="text-sm text-gray-400 flex items-center gap-2">
      <Link
        href="/courses"
        className="hover:text-white cursor-pointer text-2xl"
      >
        Courses
      </Link>
      <span className="text-blue-400 text-3xl">{">"}</span>
      <span className="text-xl mt-1">HTML Basics</span>
    </div>
  );
}
