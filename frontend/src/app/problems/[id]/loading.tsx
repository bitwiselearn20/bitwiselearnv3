export default function ProblemLoading() {
  return (
    <div className="h-screen flex bg-neutral-900 text-white overflow-hidden animate-pulse">
      {/* ================= LEFT SIDEBAR ================= */}
      <div
        className="border-r border-neutral-700 flex flex-col"
        style={{ width: 720 }}
      >
        {/* Tabs */}
        <div className="border-b border-neutral-700 bg-neutral-900 px-4 py-3 flex gap-6">
          <div className="h-4 w-24 bg-neutral-700 rounded" />
          <div className="h-4 w-20 bg-neutral-700 rounded" />
          <div className="h-4 w-24 bg-neutral-700 rounded" />
        </div>

        {/* Content */}
        <div className="flex-1 px-6 py-4 space-y-4 overflow-hidden">
          <div className="h-6 w-2/3 bg-neutral-700 rounded" />
          <div className="h-4 w-full bg-neutral-800 rounded" />
          <div className="h-4 w-11/12 bg-neutral-800 rounded" />
          <div className="h-4 w-10/12 bg-neutral-800 rounded" />

          <div className="mt-6 space-y-3">
            {Array.from({ length: 6 }).map((_, i) => (
              <div
                key={i}
                className="h-4 w-full bg-neutral-800 rounded"
              />
            ))}
          </div>
        </div>
      </div>

      {/* ================= RIGHT PANEL ================= */}
      <div className="flex-1 flex flex-col min-w-0">
        {/* Code Editor Skeleton */}
        <div className="flex-60 bg-neutral-800 border-b border-neutral-700 relative">
          {/* Editor Header */}
          <div className="absolute top-3 left-3 flex gap-3">
            <div className="h-4 w-20 bg-neutral-700 rounded" />
            <div className="h-4 w-16 bg-neutral-700 rounded" />
          </div>
        </div>

        {/* Resize Handle */}
        <div className="h-1 bg-neutral-700" />

        {/* Test Cases */}
        <div className="flex-40 p-4 space-y-3 overflow-hidden">
          <div className="h-5 w-32 bg-neutral-700 rounded" />
          {Array.from({ length: 4 }).map((_, i) => (
            <div
              key={i}
              className="h-12 bg-neutral-800 rounded"
            />
          ))}
        </div>
      </div>
    </div>
  );
}