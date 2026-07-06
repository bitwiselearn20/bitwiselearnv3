"use client";

import { useMemo, useState } from "react";
import { Input } from "./Input";

type VendorForm = {
  name: string;
  email: string;
  secondaryEmail?: string;
  tagline: string;
  phoneNumber: string;
  secondaryPhoneNumber?: string;
  websiteLink: string;
};

interface CreateVendorV1Props {
  onClose?: () => void;
  onSubmit?: (data: VendorForm) => void;
}

const isValidEmail = (email: string) =>
  /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);

export default function CreateVendorV1({
  onClose,
  onSubmit,
}: CreateVendorV1Props) {
  const [form, setForm] = useState<VendorForm>({
    name: "",
    email: "",
    secondaryEmail: "",
    tagline: "",
    phoneNumber: "",
    secondaryPhoneNumber: "",
    websiteLink: "",
  });

  const [shake, setShake] = useState(false);
  const [submitted, setSubmitted] = useState(false);

  const isPrimaryEmailValid = isValidEmail(form.email);
  const isSecondaryEmailValid =
    !form.secondaryEmail || isValidEmail(form.secondaryEmail);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setForm((prev) => ({ ...prev, [name]: value }));
  };

  const isValid = useMemo(() => {
    return (
      form.name.trim() &&
      form.tagline.trim() &&
      form.phoneNumber.trim() &&
      form.websiteLink.trim() &&
      isPrimaryEmailValid &&
      isSecondaryEmailValid
    );
  }, [form, isPrimaryEmailValid, isSecondaryEmailValid]);

  const handleSubmit = () => {
    setSubmitted(true);

    if (!isValid) {
      setShake(true);
      setTimeout(() => setShake(false), 400);
      return;
    }

    onSubmit?.(form);
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-sm">
      <div
        className={`w-full max-w-lg rounded-2xl bg-[#0B1324] p-6 shadow-xl
        animate-modal-in ${shake ? "animate-shake" : ""}`}
      >
        <h2 className="text-lg font-semibold text-white mb-4">New vendor</h2>

        <div className="space-y-4">
          <Input
            required
            label="Vendor Name*"
            name="name"
            placeholder="e.g. Acme Technologies"
            value={form.name}
            onChange={handleChange}
            showError={submitted}
            errorMessage="Vendor name is required"
          />

          <Input
            required
            label="Tagline*"
            name="tagline"
            placeholder="Brief description"
            value={form.tagline}
            onChange={handleChange}
            showError={submitted}
            errorMessage="Tagline is required"
          />

          <Input
            required
            label="Primary Email*"
            name="email"
            placeholder="contact@vendor.com"
            value={form.email}
            onChange={handleChange}
            isInvalid={!!form.email && !isPrimaryEmailValid}
            showError={submitted}
            errorMessage="Please enter a valid email address"
          />

          <Input
            label="Secondary Email"
            name="secondaryEmail"
            placeholder="Optional"
            value={form.secondaryEmail}
            onChange={handleChange}
            isInvalid={!!form.secondaryEmail && !isSecondaryEmailValid}
            errorMessage="Invalid email format"
          />

          <Input
            required
            label="Phone Number*"
            name="phoneNumber"
            placeholder="+91 98765 43210"
            value={form.phoneNumber}
            onChange={handleChange}
            showError={submitted}
            errorMessage="Phone number is required"
          />

          <Input
            label="Secondary Phone Number"
            name="secondaryPhoneNumber"
            placeholder="Optional"
            value={form.secondaryPhoneNumber}
            onChange={handleChange}
          />

          <Input
            required
            label="Website*"
            name="websiteLink"
            placeholder="https://vendor.com"
            value={form.websiteLink}
            onChange={handleChange}
            showError={submitted}
            errorMessage="Website is required"
          />
        </div>

        <div className="mt-6 flex justify-end gap-3">
          <button
            onClick={onClose}
            className="rounded-lg bg-white/5 px-4 py-2 text-sm text-white hover:bg-white/10 transition"
          >
            Cancel
          </button>

          <button
            onClick={handleSubmit}
            disabled={!isValid}
            className={`rounded-lg px-4 py-2 text-sm font-medium transition
            ${
              isValid
                ? "bg-sky-500 text-black hover:bg-sky-400 active:scale-95"
                : "bg-sky-500/40 text-black/50 cursor-not-allowed"
            }`}
          >
            Create Vendor
          </button>
        </div>
      </div>
    </div>
  );
}
