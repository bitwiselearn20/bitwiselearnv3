"use client";

import { useState } from "react";
import Image from "next/image";

type TabKey = "student" | "institute" | "vendor";

const tabs: { key: TabKey; label: string }[] = [
  { key: "student", label: "Student" },
  { key: "institute", label: "Institute" },
  { key: "vendor", label: "Vendor" },
];

const tabContent: Record<
  TabKey,
  { highlights: string[]; cta: string; href: string; image: string }
> = {
  student: {
    highlights: [
      "Personalized Learning Experience",
      "Interactive Learning Path",
      "Live Class Participation",
      "Mentor Connect",
    ],
    cta: "Access Student Portal",
    href: "/student-login",
    image: "https://images.unsplash.com/photo-1516321318423-f06f85e504b3?w=900&h=700&fit=crop",
  },
  institute: {
    highlights: [
      "Batch & Student Management",
      "Real-Time Placement Analytics",
      "Assessment & Reporting Tools",
      "Multi-Branch Oversight",
    ],
    cta: "Access Institute Portal",
    href: "/multi-login",
    image: "https://images.unsplash.com/photo-1523240795612-9a054b0db644?w=900&h=700&fit=crop",
  },
  vendor: {
    highlights: [
      "Multi-Institution Management",
      "Consolidated Training Reports",
      "Centralized Batch Oversight",
      "Partner Performance Insights",
    ],
    cta: "Access Vendor Portal",
    href: "/multi-login",
    image: "https://images.unsplash.com/photo-1531482615713-2afd69097998?w=900&h=700&fit=crop",
  },
};

export function DashboardPreview() {
  const [active, setActive] = useState<TabKey>("student");
  const content = tabContent[active];

  return (
    <section id="platform-preview" className="bg-neutral-50 py-16 sm:py-24">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <h2 className="title-react mb-3 text-center text-4xl font-bold text-black sm:text-5xl">
          Powerful Dashboards for Everyone
        </h2>
        <p className="mx-auto mb-10 max-w-2xl text-center text-lg text-neutral-600 sm:text-xl">
          A dedicated experience for every role in the placement ecosystem.
        </p>

        <div className="mb-10 flex justify-center gap-2">
          {tabs.map((tab) => (
            <button
              key={tab.key}
              type="button"
              onClick={() => setActive(tab.key)}
              className={`rounded-full px-6 py-2 text-sm font-medium transition-colors ${
                active === tab.key
                  ? "bg-black text-white"
                  : "border border-neutral-300 text-neutral-600 hover:border-black"
              }`}
            >
              {tab.label}
            </button>
          ))}
        </div>

        <div className="grid grid-cols-1 items-center gap-10 lg:grid-cols-2">
          <div>
            <ul className="space-y-4">
              {content.highlights.map((h) => (
                <li key={h} className="flex items-center gap-3">
                  <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-black text-white">
                    <svg className="h-3.5 w-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
                    </svg>
                  </span>
                  <span className="font-medium text-black">{h}</span>
                </li>
              ))}
            </ul>
            <a
              href={content.href}
              className="mt-8 inline-block rounded-full bg-black px-7 py-3 font-medium text-white hover:bg-neutral-800"
            >
              {content.cta}
            </a>
          </div>
          <div className="relative min-h-[320px] overflow-hidden rounded-2xl bg-neutral-200 shadow-lg sm:min-h-[400px]">
            <Image
              key={active}
              src={content.image}
              alt={`${active} dashboard preview`}
              fill
              className="object-cover"
              sizes="(max-width: 1024px) 100vw, 50vw"
            />
            <div className="absolute inset-0 bg-linear-to-t from-black/60 via-transparent to-transparent" />
          </div>
        </div>
      </div>
    </section>
  );
}
