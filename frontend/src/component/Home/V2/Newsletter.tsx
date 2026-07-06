"use client";

import { useState } from "react";

export function Newsletter() {
  const [status, setStatus] = useState<"idle" | "success" | "error">("idle");
  const [email, setEmail] = useState("");

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!email.trim()) {
      setStatus("error");
      return;
    }
    setStatus("success");
    setEmail("");
  }

  return (
    <section className="py-16 sm:py-24">
      <div className="mx-auto max-w-2xl px-4 text-center sm:px-6 lg:px-8">
        <h2 className="title-react text-3xl font-bold text-white sm:text-4xl">
          Subscribe to our newsletter for latest updates
        </h2>
        <ul className="mt-4 flex flex-wrap justify-center gap-6 text-sm text-neutral-400">
          <li className="flex items-center gap-2">
            <span className="text-white">✓</span> Daily design update
          </li>
          <li className="flex items-center gap-2">
            <span className="text-white">✓</span> Affiliate earning
            opportunity
          </li>
        </ul>
        <form onSubmit={handleSubmit} className="mt-8">
          <div className="flex flex-col gap-3 sm:flex-row sm:justify-center">
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="Enter your email"
              className="rounded-full border border-neutral-600 bg-neutral-900 px-4 py-3 text-white placeholder:text-neutral-500 focus:border-white focus:outline-none focus:ring-1 focus:ring-white sm:w-72"
              aria-label="Email"
            />
            <button
              type="submit"
              className="rounded-full bg-white px-6 py-3 font-medium text-black hover:bg-neutral-200"
            >
              Subscribe
            </button>
          </div>
          {status === "success" && (
            <p className="mt-3 text-sm text-neutral-300">
              Thank you! Your submission has been received!
            </p>
          )}
          {status === "error" && (
            <p className="mt-3 text-sm text-neutral-300">
              Oops! Something went wrong while submitting the form.
            </p>
          )}
        </form>
      </div>
    </section>
  );
}
