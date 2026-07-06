"use client";

import { useMemo, useState } from "react";
import { Input } from "./Input";

type StudentForm = {
  name: string;
  rollNumber: string;
  email: string;
};

interface CreateStudentV1Props {
  onClose?: () => void;
  onSubmit?: (data: StudentForm) => void;
}

const isValidEmail = (email: string) =>
  /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email);

export default function CreateStudentV1({
  onClose,
  onSubmit,
}: CreateStudentV1Props) {
  const [form, setForm] = useState<StudentForm>({
    name: "",
    rollNumber: "",
    email: "",
  });

  const [shake, setShake] = useState(false);
  const [submitted, setSubmitted] = useState(false);

  const isPrimaryEmailValid = isValidEmail(form.email);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setForm((prev) => ({ ...prev, [name]: value }));
  };

  const isValid = useMemo(() => {
    return form.name.trim() && form.rollNumber.trim() && isPrimaryEmailValid;
  }, [form, isPrimaryEmailValid]);

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
        <h2 className="text-lg font-semibold text-white mb-4">New Student</h2>

        <div className="space-y-4">
          <Input
            required
            label="Student Name*"
            name="name"
            placeholder="e.g. John Doe"
            value={form.name}
            onChange={handleChange}
            showError={submitted}
            errorMessage="Student name is required"
          />

          <Input
            required
            label="Roll Number*"
            name="rollNumber"
            placeholder="e.g. 123456"
            value={form.rollNumber}
            onChange={handleChange}
            showError={submitted}
            errorMessage="Roll number is required"
          />

          <Input
            required
            label="Primary Email*"
            name="email"
            placeholder="contact@Student.com"
            value={form.email}
            onChange={handleChange}
            isInvalid={!!form.email && !isPrimaryEmailValid}
            showError={submitted}
            errorMessage="Please enter a valid email address"
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
            Create Student
          </button>
        </div>
      </div>
    </div>
  );
}
