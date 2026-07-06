"use client";

import { useState } from "react";
import Link from "next/link";
import { ServicesHeader } from "./ServiceHeader";
import Footer from "@/component/general/Footer";
import V1HomeNav from "@/component/Home/V1/V1HomeNav";

const SERVICE_FILTERS = ["SaaS", "Training", "Partnerships"] as const;

const SERVICE_CARDS = [
  {
    id: "saas",
    badge: "CORE OFFERING",
    title: "SaaS Platform Subscription",
    description:
      "Provide your students, faculty, and vendors with a unified LMS, Institute MS, and Vendor MS. Manage learning, assessments, batches, and analytics from one platform.",
    features: [
      "Student LMS with DSA, assessments & learning paths",
      "Institute MS for batch, student & admin management",
      "Vendor MS for partner & trainer management",
      "Centralized reporting & dashboards",
    ],
    icon: "cloud",
    category: "SaaS",
  },
  {
    id: "trainers",
    badge: "EXPERT-LED TRAINING",
    title: "Trainers & Faculty Services",
    description:
      "Deploy highly experienced trainers for coding, DSA, full-stack, cloud, DevOps, aptitude, and communication to upgrade your student outcomes.",
    features: [
      "Programming: C, C++, Java, Python, SQL",
      "DSA, full-stack, DevOps & cloud technologies",
      "Aptitude, reasoning & verbal ability",
      "Workshops, bootcamps & long-term engagements",
    ],
    icon: "mentor",
    category: "Training",
  },
  {
    id: "saas-trainers",
    badge: "FULL STACK TRAINING SOLUTION",
    title: "SaaS + Trainers (End-to-End)",
    description:
      "Get the platform as your training engine — platform, trainers, content, assessments and analytics all in one integrated solution.",
    features: [
      "Platform + trainers + curriculum + content",
      "End-to-end placement & certification programs",
      "Project-based learning & graded assessments",
      "Ideal for institutes without internal training teams",
    ],
    icon: "rocket",
    category: "Partnerships",
  },
  {
    id: "partnerships",
    badge: "LONG-TERM COLLABORATION",
    title: "Academic & Vendor Partnerships",
    description:
      "Collaborate as your academic or vendor partner to co-build programs, share revenue, and deliver long-term value.",
    features: [
      "Custom programs mapped to your curriculum",
      "Joint branding & marketing initiatives",
      "Shared analytics & outcome tracking",
      "Flexible engagement & commercial models",
    ],
    icon: "document",
    category: "Partnerships",
  },
];

const ENGAGEMENT_STEPS = [
  {
    step: 1,
    badge: "PLATFORM-FIRST",
    title: "Only SaaS",
    subtitle:
      "You use our LMS, Institute MS and Vendor MS with your own trainers and content.",
    icon: "monitor",
  },
  {
    step: 2,
    badge: "FACULTY-FIRST",
    title: "Only Trainers",
    subtitle: "You retain your own systems, and we provide expert trainers.",
    icon: "person",
  },
  {
    step: 3,
    badge: "COMPLETE SOLUTION",
    title: "SaaS + Trainers",
    subtitle:
      "We handle tech + trainers + content + assessments for a fully managed training experience.",
    icon: "graduation",
  },
  {
    step: 4,
    badge: "COLLABORATIVE MODEL",
    title: "Institute / Vendor Partnership",
    subtitle:
      "We co-create programs and share responsibilities for delivery, quality and outcomes.",
    icon: "building",
  },
];

const TARGET_SEGMENTS = [
  {
    badge: "ACADEMIC",
    title: "Institutes & Universities",
    description:
      "Deliver structured, measurable training outcomes for your students across semesters, labs and placement seasons.",
    features: [
      "Central admin & department dashboards",
      "Batch, course & assessment management",
      "Placement-focused training paths",
      "Data-driven reports for management",
    ],
    icon: "graduation",
  },
  {
    badge: "VENDOR",
    title: "Training Vendors & Companies",
    description:
      "Manage your students, trainers and partner institutes with one robust SaaS backbone.",
    features: [
      "Vendor-specific panels & branding",
      "Trainer & batch allocation tools",
      "Performance & revenue tracking",
      "White-label implementation support",
    ],
    icon: "building",
  },
  {
    badge: "STUDENT",
    title: "Individual Learners",
    description:
      "Students and job-seekers can access guided learning paths, practice content, and assessments.",
    features: [
      "DSA & problem-solving practice",
      "Placement & certification-oriented courses",
      "Assessments with detailed analytics",
      "Self-paced or mentor-led options",
    ],
    icon: "learner",
  },
];

function Icon({
  name,
  className = "h-6 w-6",
}: {
  name: string;
  className?: string;
}) {
  const cls = `${className} shrink-0 text-white`;
  switch (name) {
    case "cloud":
      return (
        <svg
          className={cls}
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M3 15a4 4 0 004 4h9a5 5 0 10-.1-9.999 5.002 5.002 0 10-9.78 2.096A4.001 4.001 0 003 15z"
          />
        </svg>
      );
    case "mentor":
      return (
        <svg
          className={cls}
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z"
          />
        </svg>
      );
    case "rocket":
      return (
        <svg
          className={cls}
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M15.59 14.37a6 6 0 01-5.84 7.38v-4.8m5.84-2.58a14.98 14.98 0 006.16-12.12A14.98 14.98 0 009.631 8.41m5.96 5.96a14.926 14.926 0 01-5.841 2.58m-.119-8.54a6 6 0 00-7.381 5.84h4.8m2.581-5.84a14.927 14.927 0 00-2.58 5.84m2.699 2.7c-.103.021-.207.041-.311.06a15.09 15.09 0 01-2.448-2.448 14.9 14.9 0 01.06-.312m-2.24 2.39a4.493 4.493 0 00-1.757 4.306 4.493 4.493 0 004.306-1.758M16.5 9a4.5 4.5 0 11-9 0 4.5 4.5 0 019 0z"
          />
        </svg>
      );
    case "document":
      return (
        <svg
          className={cls}
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
          />
        </svg>
      );
    case "monitor":
      return (
        <svg
          className={cls}
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"
          />
        </svg>
      );
    case "person":
      return (
        <svg
          className={cls}
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
          />
        </svg>
      );
    case "graduation":
      return (
        <svg
          className={cls}
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M12 14l9-5-9-5-9 5 9 5z"
          />
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M12 14l6.16-3.422a12.083 12.083 0 01.665 6.479A11.952 11.952 0 0012 20.055a11.952 11.952 0 00-6.824-2.998 12.078 12.078 0 01.665-6.479L12 14z"
          />
        </svg>
      );
    case "building":
      return (
        <svg
          className={cls}
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
          />
        </svg>
      );
    case "learner":
      return (
        <svg
          className={cls}
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253"
          />
        </svg>
      );
    default:
      return null;
  }
}

function CheckIcon() {
  return (
    <svg
      className="h-5 w-5 shrink-0 text-emerald-500"
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M5 13l4 4L19 7"
      />
    </svg>
  );
}

export function OurServicesV1() {
  const [activeFilter, setActiveFilter] = useState<string | null>(null);
  const filteredCards =
    activeFilter === null
      ? SERVICE_CARDS
      : SERVICE_CARDS.filter(
          (c) => c.category.toLowerCase() === activeFilter.toLowerCase(),
        );

  return (
    <>
      {/* navabar  */}
      <V1HomeNav />
      {/* Header  */}
      <ServicesHeader />
      {/* Our Services */}
      <section className="py-16 sm:py-24">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <span className="inline-block rounded-full border border-neutral-600 px-4 py-1.5 text-center text-xs font-medium uppercase tracking-wider text-neutral-400">
            Our Services
          </span>
          <h2 className="title-react mt-4 text-3xl font-bold text-white sm:text-4xl">
            Services that adapt to your training ecosystem
          </h2>
          <p className="mx-auto mt-4 max-w-3xl text-neutral-400">
            Whether you are an institute, a vendor, or an individual learner, we
            can plug into your world with just the right mix of platform, people
            and processes.
          </p>V1HomeNav
          <div className="mt-8 flex flex-wrap justify-center gap-3">
            {SERVICE_FILTERS.map((f) => (
              <button
                key={f}
                type="button"
                onClick={() => setActiveFilter(activeFilter === f ? null : f)}
                className={`rounded-full px-5 py-2.5 text-sm font-medium uppercase tracking-wide transition-colors ${
                  activeFilter === f
                    ? "bg-white text-black"
                    : "border border-neutral-600 bg-neutral-900 text-neutral-300 hover:bg-neutral-800"
                }`}
              >
                {f}
              </button>
            ))}
          </div>
          <div className="mt-12 grid gap-6 sm:grid-cols-2">
            {filteredCards.map((card) => (
              <article
                key={card.id}
                className="relative rounded-2xl border border-neutral-700 bg-neutral-900 p-6 transition-colors hover:border-neutral-600 sm:p-8"
              >
                <span className="absolute right-4 top-4 rounded bg-neutral-800 px-2.5 py-1 text-xs font-medium uppercase tracking-wider text-neutral-400">
                  {card.badge}
                </span>
                <div className="flex gap-4 pr-28">
                  <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-neutral-800">
                    <Icon name={card.icon} />
                  </div>
                  <div>
                    <h3 className="title-react text-xl font-semibold text-white">
                      {card.title}
                    </h3>
                    <p className="mt-2 text-sm text-neutral-400">
                      {card.description}
                    </p>
                  </div>
                </div>
                <ul className="mt-6 space-y-2">
                  {card.features.map((f) => (
                    <li
                      key={f}
                      className="flex items-start gap-2 text-sm text-neutral-300"
                    >
                      <span className="mt-0.5 text-neutral-500">•</span>
                      {f}
                    </li>
                  ))}
                </ul>
                <Link
                  href="/contact"
                  className="mt-6 inline-flex items-center gap-1 text-sm font-medium text-white hover:underline"
                >
                  Talk to us about this
                  <span aria-hidden>→</span>
                </Link>
              </article>
            ))}
          </div>
        </div>
      </section>

      {/* Engagement Models */}
      <section className="border-t border-neutral-800 py-16 sm:py-24">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <span className="inline-block rounded-full border border-neutral-600 px-4 py-1.5 text-center text-xs font-medium uppercase tracking-wider text-neutral-400">
            Engagement Models
          </span>
          <h2 className="title-react mt-4 text-3xl font-bold text-white sm:text-4xl">
            Models that match <span className="italic">how you operate</span>
          </h2>
          <p className="mx-auto mt-4 max-w-2xl text-neutral-400">
            Start with what you need today — from platform-only to a fully
            managed training solution — and scale as your strategy evolves.
          </p>
          <p className="mt-4 flex flex-wrap justify-center gap-2 text-sm text-neutral-500">
            <span>Flexible</span>
            <span>•</span>
            <span>Scalable</span>
            <span>•</span>
            <span>Modular</span>
          </p>
          <div className="mt-12 grid gap-6 sm:grid-cols-2 lg:grid-cols-4">
            {ENGAGEMENT_STEPS.map((item, i) => (
              <article
                key={item.step}
                className="relative rounded-2xl border border-neutral-700 bg-neutral-900 p-6"
              >
                {i < ENGAGEMENT_STEPS.length - 1 && (
                  <span
                    className="absolute -right-3 top-1/2 z-10 hidden -translate-y-1/2 text-neutral-600 lg:inline"
                    aria-hidden
                  >
                    →
                  </span>
                )}
                <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-neutral-800">
                  <Icon name={item.icon} className="h-5 w-5" />
                </div>
                <span className="mt-3 block text-xs font-medium uppercase tracking-wider text-neutral-500">
                  Step {item.step} · {item.badge}
                </span>
                <h3 className="title-react mt-1 text-lg font-semibold text-white">
                  {item.title}
                </h3>
                <p className="mt-2 text-sm text-neutral-400">{item.subtitle}</p>
              </article>
            ))}
          </div>
        </div>
      </section>

      {/* Target Segments */}
      <section className="border-t border-neutral-800 py-16 sm:py-24">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 mb-10">
          <span className="inline-block rounded-full border border-neutral-600 px-4 py-1.5 text-center text-xs font-medium uppercase tracking-wider text-neutral-400">
            Target Segments
          </span>
          <h2 className="title-react mt-4 text-3xl font-bold text-white sm:text-4xl">
            Designed for every layer of your{" "}
            <span className="italic">ecosystem</span>
          </h2>
          <p className="mx-auto mt-4 max-w-2xl text-neutral-400">
            We bridge the expectations of institutes, vendors and learners
            through an integrated training platform that keeps everyone in sync.
          </p>
          <div className="mt-12 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {TARGET_SEGMENTS.map((seg) => (
              <article
                key={seg.title}
                className="rounded-2xl border border-neutral-700 bg-neutral-900 p-6 transition-colors hover:border-neutral-600 sm:p-8"
              >
                <div className="flex items-start gap-4">
                  <div className="flex h-12 w-12 shrink-0 items-center justify-center rounded-xl bg-neutral-800">
                    <Icon name={seg.icon} />
                  </div>
                  <div>
                    <span className="text-xs font-medium uppercase tracking-wider text-neutral-500">
                      {seg.badge}
                    </span>
                    <h3 className="title-react mt-1 text-xl font-semibold text-white">
                      {seg.title}
                    </h3>
                    <p className="mt-2 text-sm text-neutral-400">
                      {seg.description}
                    </p>
                  </div>
                </div>
                <ul className="mt-6 space-y-3">
                  {seg.features.map((f) => (
                    <li
                      key={f}
                      className="flex items-start gap-2 text-sm text-neutral-300"
                    >
                      <CheckIcon />
                      {f}
                    </li>
                  ))}
                </ul>
              </article>
            ))}
          </div>
        </div>

        {/* footer  */}
        <Footer />
      </section>
    </>
  );
}
