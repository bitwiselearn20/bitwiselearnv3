"use client";

interface ConfirmSubmitProps {
  open: boolean;
  onCancel: () => void;
  onConfirm: () => void;
}

export default function ConfirmSubmit({
  open,
  onCancel,
  onConfirm,
}: ConfirmSubmitProps) {
  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm">
      <div className="w-90 rounded-xl bg-neutral-900 p-6 text-white shadow-xl">
        <h2 className="text-lg font-semibold">Submit Assignment?</h2>

        <p className="mt-2 text-sm text-neutral-400">
          Once submitted, you won’t be able to change your answers.
        </p>

        <div className="mt-6 flex justify-end gap-3">
          <button
            onClick={onCancel}
            className="rounded-md px-4 py-2 text-sm text-neutral-300 hover:bg-neutral-800"
          >
            Cancel
          </button>

          <button
            onClick={onConfirm}
            className="rounded-md bg-red-600 px-4 py-2 text-sm font-medium hover:bg-red-700"
          >
            Submit
          </button>
        </div>
      </div>
    </div>
  );
}
