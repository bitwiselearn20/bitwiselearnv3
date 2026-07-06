"use client";

import { X } from "lucide-react";

type EntityData = {
  id: string;
  [key: string]: any;
};

type Props = {
  entityType: string;
  data: EntityData;
  onClose: () => void;
};

function formatValue(value: any) {
  if (!value) return "-";
  if (typeof value === "string" && value.includes("T")) {
    const date = new Date(value);
    if (!isNaN(date.getTime())) {
      return date.toLocaleString();
    }
  }
  return String(value);
}

export default function EntityDetailsModal({
  entityType,
  data,
  onClose,
}: Props) {
  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm">
      <div className="relative w-full max-w-xl rounded-2xl bg-divBg border border-white/10 shadow-2xl p-6">
        {/* Close */}
        <button
          onClick={onClose}
          className="absolute top-4 right-4 text-white/50 hover:text-white transition"
        >
          <X size={20} />
        </button>

        {/* Header */}
        <div className="mb-6">
          <h2 className="text-xl font-semibold text-white">
            {entityType} Details
          </h2>
          <p className="text-sm text-white/40 mt-1">ID: {data.id}</p>
        </div>

        {/* Content */}
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-x-6 gap-y-5 max-h-[60vh] overflow-y-auto pr-1">
          {Object.entries(data).map(([key, value]) => (
            <div key={key}>
              <p className="text-[11px] uppercase tracking-wide text-primaryBlue mb-1">
                {key.replace(/_/g, " ")}
              </p>
              <p className="text-sm text-white break-words">
                {formatValue(value)}
              </p>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
