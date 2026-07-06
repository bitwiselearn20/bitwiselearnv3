"use client";

import React, { useState } from "react";
import { X } from "lucide-react";

type Props = {
  openForm: (value: boolean) => void;
  onSubmit?: (data: AuthFormData) => void;
};

type AuthFormData = {
  name: string;
  email: string;
};

export default function AdminForm({ openForm, onSubmit }: Props) {
  const [formData, setFormData] = useState<AuthFormData>({ name: "", email: "" });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit?.(formData);
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm">
      <div className="relative w-full max-w-sm rounded-2xl border border-white/10 bg-divBg p-6 shadow-2xl">
        <button
          onClick={() => openForm(false)}
          className="absolute right-4 top-4 text-white/50 hover:text-white transition"
        >
          <X size={20} />
        </button>

        <div className="mb-6">
          <h2 className="mt-1 text-lg font-semibold text-white">Create User</h2>
        </div>

        <form onSubmit={handleSubmit} className="space-y-5">
          <Input label="Name" name="name" value={formData.name} onChange={handleChange} />
          <Input
            label="Email"
            name="email"
            type="email"
            value={formData.email}
            onChange={handleChange}
          />

          <div className="flex justify-between pt-4">
            <span />
            <button
              type="submit"
              className="rounded-md bg-primaryBlue px-4 py-2 text-sm font-semibold text-white transition hover:bg-primaryBlue/90"
            >
              Create User
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

/* ---------- UI Primitives ---------- */

function Label({ children }: { children: React.ReactNode }) {
  return (
    <label className="text-[11px] uppercase tracking-wide text-primaryBlue">
      {children}
    </label>
  );
}

function Input({
  label,
  ...props
}: React.InputHTMLAttributes<HTMLInputElement> & { label: string }) {
  return (
    <div>
      <Label>{label}</Label>
      <input
        {...props}
        className="mt-1 w-full rounded-lg border border-white/10 bg-black/30 px-3 py-2 text-sm text-white focus:ring-2 focus:ring-primaryBlue"
      />
    </div>
  );
}
