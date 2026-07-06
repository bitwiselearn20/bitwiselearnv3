"use client";

interface ConfirmExitProps {
  open: boolean;
  onCancel: () => void;
  onConfirm: () => void;
}

export default function ConfirmExit({
  open,
  onCancel,
  onConfirm,
}: ConfirmExitProps) {
  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm">
      <div className="w-95 rounded-xl bg-neutral-900 p-6 text-white shadow-xl">
        <h2 className="text-lg font-semibold text-red-500">Exit Assignment?</h2>

        <p className="mt-2 text-sm text-neutral-400">
          Your progress will be lost and the assignment may be auto-submitted.
        </p>

        <div className="mt-6 flex justify-end gap-3">
          <button
            onClick={onCancel}
            className="rounded-md px-4 py-2 text-sm text-neutral-300 hover:bg-neutral-800"
          >
            Stay
          </button>

          <button
            onClick={onConfirm}
            className="rounded-md bg-red-600 px-4 py-2 text-sm font-medium hover:bg-red-700"
          >
            Exit
          </button>
        </div>
      </div>
    </div>
  );
}
