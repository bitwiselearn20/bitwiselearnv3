"use client";

import { useState } from "react";
import { faqItems } from "@/lib/content/faq";

export function FAQ() {
  const [openId, setOpenId] = useState<string | null>(null);

  return (
    <section className="py-16 sm:py-24">
      <div className="mx-auto max-w-3xl px-4 sm:px-6 lg:px-8">
        <h2 className="title-react mb-4 text-center text-4xl font-bold text-white sm:text-5xl">
          Frequently asked questions
        </h2>
        <p className="mb-10 text-center text-neutral-400">
          Our learners have used Bitwise to transition careers, secure
          promotions, and break into competitive industries.
        </p>
        <div className="space-y-3">
          {faqItems.map((item) => {
            const isOpen = openId === item.id;
            return (
              <div
                key={item.id}
                className="overflow-hidden rounded-2xl border border-neutral-700 bg-neutral-900 shadow-lg"
              >
                <button
                  type="button"
                  onClick={() =>
                    setOpenId((prev) => (prev === item.id ? null : item.id))
                  }
                  className="flex w-full items-center gap-4 px-5 py-4 text-left font-medium text-white hover:bg-neutral-800"
                >
                  <span
                    className={`flex h-8 w-8 shrink-0 items-center justify-center rounded-full transition-colors ${
                      isOpen ? "bg-white text-black" : "bg-neutral-600 text-white"
                    }`}
                  >
                    {isOpen ? (
                      <svg className="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                      </svg>
                    ) : (
                      <svg className="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                      </svg>
                    )}
                  </span>
                  <span className="flex-1">{item.question}</span>
                </button>
                {isOpen && (
                  <div className="border-t border-neutral-700 px-5 pb-5 pl-[3.25rem] pt-1 text-sm leading-relaxed text-neutral-400">
                    {item.answer}
                  </div>
                )}
              </div>
            );
          })}
        </div>
      </div>
    </section>
  );
}
