export default function AttemptLoading() {
  return (
    <div className="h-screen w-full bg-[#0f0f0f] flex flex-col animate-pulse">
      {/* TOP BAR */}
      <div className="flex justify-end items-center px-4 py-3 border-b border-white/10 gap-3">
        <div className="h-6 w-6 rounded-full bg-white/10" />
        <div className="h-9 w-24 rounded-md bg-white/10" />
      </div>

      {/* MAIN CONTENT */}
      <div className="flex-1 grid grid-cols-2 gap-4 p-4 min-h-0">
        {/* LEFT SECTION */}
        <div className="h-full rounded-xl bg-white/5 border border-white/10 p-4 flex flex-col gap-4">
          <div className="h-6 w-40 bg-white/10 rounded" />
          <div className="h-4 w-64 bg-white/10 rounded" />

          <div className="flex-1 rounded-lg bg-white/10" />

          <div className="flex justify-between">
            <div className="h-9 w-24 bg-white/10 rounded-md" />
            <div className="h-9 w-24 bg-white/10 rounded-md" />
          </div>
        </div>

        {/* RIGHT SECTION */}
        <div className="h-full rounded-xl bg-white/5 border border-white/10 p-4 flex flex-col gap-4">
          <div className="h-6 w-48 bg-white/10 rounded" />

          {/* MCQ / Code placeholder */}
          <div className="flex-1 rounded-lg bg-white/10" />

          {/* Actions */}
          <div className="flex gap-3">
            <div className="h-9 w-28 bg-white/10 rounded-md" />
            <div className="h-9 w-28 bg-white/10 rounded-md" />
          </div>
        </div>
      </div>
    </div>
  );
}