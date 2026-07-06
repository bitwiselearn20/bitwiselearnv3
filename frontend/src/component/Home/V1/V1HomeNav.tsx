"use client";

import React, { useState } from "react";
import { Menu, X } from "lucide-react";
import Link from "next/link";
import { Klee_One } from "next/font/google";

const kleeOne = Klee_One({
  subsets: ["latin"],
  weight: ["400", "600"],
});

const V1HomeNav = () => {
  const [open, setOpen] = useState(false);

  return (
    <nav
      className={`${kleeOne.className} w-full text-white fixed top-0 left-0 z-50`}
    >
      <div className="max-w-7xl mx-auto px-6 py-4 flex items-center justify-around">
        {/* Logo */}
        <div className="text-4xl font-semibold">
          <Link href="/">
            <button className="cursor-pointer">
              <span className="text-primaryBlue">B</span>
              <span>itwise</span> Learn
            </button>
          </Link>
        </div>

        {/* Desktop Links */}
        <div className="hidden md:flex items-center gap-6 text-xl font-light">
          <Link href="/" className="hover:text-blue-400 transition">
            Home
          </Link>
          <Link href="/about" className="hover:text-blue-400 transition">
            About
          </Link>
          <Link href="/contact" className="hover:text-blue-400 transition">
            Contact
          </Link>
          <Link href="/our-services" className="hover:text-blue-400 transition">
            Services
          </Link>
        </div>

        {/* Sign in button (desktop) */}
        <div className="hidden md:block">
          <Link
            href="/student-login"
            className="px-6 py-2.5 rounded-full bg-primaryBlue hover:bg-blue-500 transition text-base"
          >
            Sign in
          </Link>
        </div>

        {/* Mobile menu button */}
        <button className="md:hidden" onClick={() => setOpen(!open)}>
          {open ? <X size={22} /> : <Menu size={22} />}
        </button>
      </div>

      {/* Mobile Menu */}
      {open && (
        <div className="md:hidden bg-black/80 backdrop-blur-md px-6 pb-6">
          <div className="flex flex-col gap-4 text-sm">
            <a href="#" className="hover:text-blue-400">
              Home
            </a>
            <a href="#" className="hover:text-blue-400">
              About
            </a>
            <a href="#" className="hover:text-blue-400">
              Contact
            </a>
            <button className="mt-2 px-5 py-2 rounded-full bg-blue-500 hover:bg-blue-600 transition">
              Sign in
            </button>
          </div>
        </div>
      )}
    </nav>
  );
};

export default V1HomeNav;
