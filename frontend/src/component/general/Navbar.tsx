"use client";

import { getColors } from "@/component/general/(Color Manager)/useColors";
import Link from "next/link";
import { useState } from "react";

const navLinks = [
  { href: "/#platform-preview", label: "FEATURES" },
  { href: "/listed-courses", label: "COURSES" },
  { href: "/our-services", label: "SERVICES" },
  { href: "/contact", label: "CONTACT US" },
];

type NavbarProps = {
  /** "light" renders a white-background nav for pages using the light theme
   * (currently only the landing page). Other pages keep the default dark nav. */
  theme?: "dark" | "light";
};

export function Navbar({ theme = "dark" }: NavbarProps) {
  const [menuOpen, setMenuOpen] = useState(false);
  const Colors = getColors();
  const isLight = theme === "light";

  const pillClasses = isLight
    ? "border-neutral-200 bg-white"
    : "border-neutral-700 bg-neutral-900";
  const textClasses = isLight ? "text-black" : "text-white";
  const mutedTextClasses = isLight ? "text-neutral-500" : "text-neutral-400";
  const outlineButtonClasses = isLight
    ? "border-neutral-300 text-black hover:border-black"
    : "border-neutral-500 text-white hover:border-white";
  const mobileMenuClasses = isLight
    ? "border-neutral-200 bg-white"
    : "border-neutral-700 bg-neutral-900";
  const mobileLinkClasses = isLight
    ? "text-black hover:opacity-70"
    : "text-neutral-400 hover:text-white";

  return (
    <header className="sticky top-0 z-50 px-4 pt-5 pb-2 sm:px-6">
      <div className="mx-auto max-w-[75%] min-w-[320px] sm:max-w-2xl md:max-w-4xl lg:max-w-5xl">
        <div className={`flex h-14 items-center justify-between rounded-3xl border px-4 shadow-sm sm:px-5 ${pillClasses}`}>
          <Link
            href="/"
            className={`flex shrink-0 items-center justify-center hover:opacity-80 cursor-pointer font-bold ${textClasses}`}
            aria-label="Home"
          >
            <span className={`${Colors.text.special}`}>B</span>{" "}
            <span>itwise Learn</span>
          </Link>

          <nav className="hidden items-center gap-8 md:flex">
            {navLinks.map((link) => (
              <Link
                key={link.href}
                href={link.href}
                className={`text-sm font-medium tracking-wide hover:opacity-80 ${textClasses}`}
              >
                {link.label}
              </Link>
            ))}
          </nav>

          <div className="flex items-center gap-4">
            <Link
              href="/student-login"
              className={`hidden text-sm font-medium hover:opacity-80 lg:inline-block ${textClasses}`}
            >
              STUDENT LOGIN
            </Link>
            <Link
              href="/multi-login"
              className={`hidden rounded-full border px-4 py-1.5 text-sm font-medium sm:inline-block ${outlineButtonClasses}`}
            >
              INSTITUTE LOGIN
            </Link>

            <button
              type="button"
              onClick={() => setMenuOpen((o) => !o)}
              className={`rounded p-2 md:hidden ${mutedTextClasses}`}
              aria-label="Menu"
            >
              {menuOpen ? "✕" : "☰"}
            </button>
          </div>
        </div>
      </div>

      {menuOpen && (
        <div className={`rounded-b-3xl border border-t-0 md:hidden ${mobileMenuClasses}`}>
          <nav className="flex flex-col gap-2 px-4 py-4">
            {navLinks.map((link) => (
              <Link
                key={link.href}
                href={link.href}
                className={`py-2 ${mobileLinkClasses}`}
                onClick={() => setMenuOpen(false)}
              >
                {link.label}
              </Link>
            ))}
            <Link
              href="/student-login"
              className={`py-2 ${mobileLinkClasses}`}
              onClick={() => setMenuOpen(false)}
            >
              STUDENT LOGIN
            </Link>
            <Link
              href="/multi-login"
              className={`py-2 ${mobileLinkClasses}`}
              onClick={() => setMenuOpen(false)}
            >
              INSTITUTE LOGIN
            </Link>
          </nav>
        </div>
      )}
    </header>
  );
}
