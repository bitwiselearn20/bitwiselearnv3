"use client";

import Footer from "@/component/general/Footer";
import { Navbar } from "@/component/general/Navbar";
import React, { useState } from "react";
import { sendComplaint } from "@/api/contact/sendComplaint";
import toast from "react-hot-toast";

const contactEmail = process.env.NEXT_PUBLIC_CONTACT_EMAIL || "";

function EnvelopeIcon() {
  return (
    <svg className="h-5 w-5 shrink-0 text-neutral-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden>
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
    </svg>
  );
}

function PhoneIcon() {
  return (
    <svg className="h-5 w-5 shrink-0 text-neutral-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden>
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
    </svg>
  );
}

function SendIcon() {
  return (
    <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden>
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
    </svg>
  );
}

const ContactV2 = () => {
  const [status, setStatus] = useState<"idle" | "sending" | "sent" | "error">("idle");

  async function handleSubmit(e:React.FormEvent<HTMLFormElement>){
    e.preventDefault();
    setStatus("sending");
    const toastId = toast.loading("Sending Complaint...");

    const formData = new FormData(e.currentTarget);

    const payload = {
        name : formData.get("name") as string,
        email : formData.get("email") as string,
        phone : formData.get("phone") as string,
        message : formData.get("message") as string,
    };

    try {
        await sendComplaint(payload);
        setStatus("sent");
        toast.success("Complaint sent successfully !",{id:toastId});
        e.currentTarget.reset();
    } catch (error) {
        setStatus("error");
        toast.error("Failed to send message. Please try again",{id:toastId});
    }
    finally{
        setStatus("idle");
    }
  }

  return (
    <>
     <div className="relative w-full">
        <Navbar />
      </div>
      {/* Header Section */}
      <section className="border-b border-neutral-800 bg-black/40 py-16 sm:py-20">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <h1 className="text-center text-3xl font-semibold tracking-tight text-white sm:text-4xl">
            Contact Us
          </h1>
          <p className="mx-auto mt-3 max-w-2xl text-center text-neutral-400">
            Get in touch — we&apos;re here to help with any questions.
          </p>
        </div>
      </section>

      {/* Content Section */}
      <section className="py-16 sm:py-24">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <span className="inline-block rounded-full border border-neutral-600 px-4 py-1.5 text-center text-xs font-medium uppercase tracking-wider text-neutral-400">
            Get in touch
          </span>

          <h2 className="title-react mt-4 text-3xl font-bold text-white sm:text-4xl">
            Let&apos;s build the future of EdTech together.
          </h2>

          <p className="mx-auto mt-4 max-w-2xl text-neutral-400">
            We&apos;re here to help! Whether you&apos;re a student with a question, a faculty member needing support, or an institution interested in partnering with us, feel free to reach out.
          </p>

          <div className="mt-12 grid gap-8 lg:grid-cols-2">
            {/* Left: Contact info */}
            <div className="rounded-2xl border border-neutral-700 bg-neutral-900/50 p-6 sm:p-8">
              <div className="space-y-6">
                <div>
                  <div className="flex items-center gap-3">
                    <EnvelopeIcon />
                    <span className="text-sm font-medium text-neutral-300">Email Us</span>
                  </div>
                  <a
                    href={`mailto:${contactEmail}`}
                    className="mt-2 block font-medium text-white hover:underline"
                  >
                    {contactEmail}
                  </a>
                </div>

                <div>
                  <div className="flex items-center gap-3">
                    <PhoneIcon />
                    <span className="text-sm font-medium text-neutral-300">Phone</span>
                  </div>
                  <p className="mt-2 text-white">+91 9787777547</p>
                </div>

                <div className="rounded-xl border border-neutral-700 bg-neutral-800/50 p-4">
                  <h3 className="font-semibold text-white">Looking for career opportunities?</h3>
                  <p className="mt-1 text-sm text-neutral-400">
                    Send us your profile at{" "}
                    <a href={`mailto:${contactEmail}`} className="text-white hover:underline">
                      {contactEmail}
                    </a>
                  </p>
                </div>
              </div>
            </div>

            {/* Right: Form */}
            <div className="rounded-2xl border border-neutral-700 bg-neutral-900/50 p-6 sm:p-8">
              <form onSubmit={handleSubmit} className="space-y-5">
                <div>
                  <label htmlFor="name" className="block text-sm font-medium text-neutral-300">
                    Full Name
                  </label>
                  <input
                    id="name"
                    name="name"
                    type="text"
                    required
                    placeholder="John Doe"
                    className="mt-1.5 w-full rounded-lg border border-neutral-600 bg-neutral-800 px-4 py-3 text-white placeholder-neutral-500 focus:border-neutral-500 focus:outline-none focus:ring-1 focus:ring-neutral-500"
                  />
                </div>

                <div>
                  <label htmlFor="email" className="block text-sm font-medium text-neutral-300">
                    Email Address
                  </label>
                  <input
                    id="email"
                    name="email"
                    type="email"
                    required
                    placeholder="john@example.com"
                    className="mt-1.5 w-full rounded-lg border border-neutral-600 bg-neutral-800 px-4 py-3 text-white placeholder-neutral-500 focus:border-neutral-500 focus:outline-none focus:ring-1 focus:ring-neutral-500"
                  />
                </div>

                <div>
                  <label htmlFor="phone" className="block text-sm font-medium text-neutral-300">
                    Phone Number
                  </label>
                  <div className="mt-1.5 flex gap-2">
                    <input
                      id="country"
                      name="country"
                      type="text"
                      placeholder="+1"
                      className="w-20 rounded-lg border border-neutral-600 bg-neutral-800 px-4 py-3 text-white placeholder-neutral-500 focus:border-neutral-500 focus:outline-none focus:ring-1 focus:ring-neutral-500"
                    />
                    <input
                      id="phone"
                      name="phone"
                      type="tel"
                      placeholder="9876543210"
                      className="flex-1 rounded-lg border border-neutral-600 bg-neutral-800 px-4 py-3 text-white placeholder-neutral-500 focus:border-neutral-500 focus:outline-none focus:ring-1 focus:ring-neutral-500"
                    />
                  </div>
                </div>

                <div>
                  <label htmlFor="message" className="block text-sm font-medium text-neutral-300">
                    Message
                  </label>
                  <textarea
                    id="message"
                    name="message"
                    rows={4}
                    required
                    placeholder="Tell us about your requirements..."
                    className="mt-1.5 w-full resize-y rounded-lg border border-neutral-600 bg-neutral-800 px-4 py-3 text-white placeholder-neutral-500 focus:border-neutral-500 focus:outline-none focus:ring-1 focus:ring-neutral-500"
                  />
                </div>

                <button
                  type="submit"
                  disabled={status === "sending"}
                  className="inline-flex w-full items-center justify-center gap-2 rounded-full bg-white px-6 py-3 font-medium text-black transition-colors hover:bg-neutral-200 focus:outline-none focus:ring-2 focus:ring-neutral-500 focus:ring-offset-2 focus:ring-offset-black disabled:opacity-70 sm:w-auto"
                >
                  <SendIcon />
                  {status === "sending"
                    ? "Sending..."
                    : status === "sent"
                    ? "Message sent"
                    : "Send Message"}
                </button>
              </form>
            </div>
          </div>
        </div>
      </section>
      <Footer />
    </>
  );
};

export default ContactV2;