"use client";

import { X, Trash, PenLine, Check } from "lucide-react";
import { useState } from "react";

type Role = "SUPER_ADMIN" | "ADMIN" | "INSTITUTION" | "TEACHER";

type EntityData = {
  id: number;
  name?: string;
  batchname?: string;
  [key: string]: any;
};

type Props = {
  entityType: string;
  data: EntityData;
  role: Role;
  onClose: () => void;
  onDelete: (id: number) => void;
};

const permissions: Record<Role, ("EDIT" | "DELETE")[]> = {
  SUPER_ADMIN: ["EDIT", "DELETE"],
  ADMIN: ["EDIT"],
  INSTITUTION: [],
  TEACHER: [],
};

export default function EntityDetailsModal({
  entityType,
  data,
  role,
  onClose,
  onDelete,
}: Props) {
  const [isEditing, setIsEditing] = useState(false);
  const [editedData, setEditedData] = useState<EntityData>({ ...data });

  const canEdit = permissions[role].includes("EDIT");
  const canDelete = permissions[role].includes("DELETE");

  function handleSave() {
    // 🔹 Hook API PATCH here later
    setIsEditing(false);
  }

  const entries = Object.entries(editedData).filter(
    ([key]) => key !== "id" // 👈 prevent editing ID
  );

  return (
    <div className="fixed inset-0 z-50 bg-black/60 backdrop-blur-sm flex items-center justify-center">
      <div className="bg-divBg rounded-2xl p-6 w-full max-w-xl relative shadow-xl">
        {/* Close */}
        <button
          onClick={onClose}
          className="absolute top-4 right-4 text-white/70 hover:text-white"
        >
          <X />
        </button>

        {/* Header */}
        <div className="mb-6">
          <h2 className="text-white text-2xl font-semibold">
            {entityType} Details
          </h2>
          <p className="text-white/60 text-sm mt-1">
            ID: #{data.id} • Role: {role.replace("_", " ")}
          </p>
        </div>

        {/* Details */}
        <div className="max-h-72 overflow-y-auto pr-1">
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
            {entries.map(([key, value]) => (
              <div key={key} className="flex flex-col">
                <label className="text-primaryBlue text-xs uppercase tracking-wide mb-1">
                  {key.replace(/_/g, " ")}
                </label>

                {isEditing ? (
                  <input
                    value={value ?? ""}
                    onChange={(e) =>
                      setEditedData((prev) => ({
                        ...prev,
                        [key]: e.target.value,
                      }))
                    }
                    className="w-full bg-transparent border-b border-primaryBlue text-white outline-none"
                  />
                ) : (
                  <p className="text-white wrap-break-word">{String(value)}</p>
                )}
              </div>
            ))}
          </div>

          {!entries.length && (
            <p className="text-white/50 text-sm text-center mt-4">
              No matching fields
            </p>
          )}
        </div>

        {/* Actions */}
        <div className="flex justify-end gap-4 mt-6">
          {canEdit && (
            <button
              onClick={isEditing ? handleSave : () => setIsEditing(true)}
              className="flex items-center gap-2 text-primaryBlue hover:opacity-80"
            >
              {isEditing ? <Check size={18} /> : <PenLine size={18} />}
              {isEditing ? "Save" : "Edit"}
            </button>
          )}

          {canDelete && (
            <button
              onClick={() => onDelete(data.id)}
              className="flex items-center gap-2 text-red-500 hover:opacity-80"
            >
              <Trash size={18} />
              Delete
            </button>
          )}
        </div>
      </div>
    </div>
  );
}
