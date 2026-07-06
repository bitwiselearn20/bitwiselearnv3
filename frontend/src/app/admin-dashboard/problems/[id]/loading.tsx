export default function LoadingProblem() {
  return (
    <div className="flex h-screen gap-4 bg-[#0f0f0f] animate-pulse p-4">
      {/* LEFT: Problem Description */}
      <div className="w-[55%] rounded-xl border border-white/10 bg-white/5 p-6 flex flex-col gap-4">
        {/* Title */}
        <div className="h-7 w-3/4 bg-white/10 rounded-md" />

        {/* Meta info */}
        <div className="flex gap-4">
          <div className="h-4 w-24 bg-white/10 rounded-md" />
          <div className="h-4 w-20 bg-white/10 rounded-md" />
          <div className="h-4 w-28 bg-white/10 rounded-md" />
        </div>

        {/* Description paragraphs */}
        <div className="space-y-3 mt-4">
          <div className="h-4 w-full bg-white/10 rounded-md" />
          <div className="h-4 w-full bg-white/10 rounded-md" />
          <div className="h-4 w-5/6 bg-white/10 rounded-md" />
          <div className="h-4 w-4/6 bg-white/10 rounded-md" />
        </div>

        {/* Example / constraints */}
        <div className="mt-6 space-y-3">
          <div className="h-5 w-40 bg-white/10 rounded-md" />
          <div className="h-24 w-full bg-white/10 rounded-lg" />
        </div>
      </div>

      {/* RIGHT: Info / Trial Panel */}
      <div className="flex-1 rounded-xl border border-white/10 bg-white/5 p-6 flex flex-col gap-4">
        {/* Section title */}
        <div className="h-6 w-40 bg-white/10 rounded-md" />

        {/* Content blocks */}
        <div className="space-y-3">
          <div className="h-4 w-full bg-white/10 rounded-md" />
          <div className="h-4 w-5/6 bg-white/10 rounded-md" />
          <div className="h-4 w-4/6 bg-white/10 rounded-md" />
        </div>

        {/* Code editor / output placeholder */}
        <div className="mt-6 flex-1 bg-white/10 rounded-xl" />

        {/* Buttons */}
        <div className="flex gap-3 mt-4">
          <div className="h-10 w-28 bg-white/10 rounded-lg" />
          <div className="h-10 w-32 bg-white/10 rounded-lg" />
        </div>
      </div>
    </div>
  );
}