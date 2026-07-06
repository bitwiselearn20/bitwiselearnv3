export default function Loading() {
  return (
    <div className="min-h-screen bg-[#0f0f0f] p-6 animate-pulse">
      <div className="flex gap-6 max-w-screen mx-auto">

        {/* Sidebar Skeleton */}
        <aside className="w-75 shrink-0 rounded-xl bg-zinc-900 border border-zinc-800 p-5 space-y-4">
          <div className="h-10 w-10 rounded-full bg-zinc-800" />
          <div className="h-5 w-40 bg-zinc-800 rounded" />
          <div className="h-4 w-32 bg-zinc-800 rounded" />

          <div className="space-y-3 mt-6">
            {[1, 2, 3].map((i) => (
              <div key={i} className="h-4 w-full bg-zinc-800 rounded" />
            ))}
          </div>
        </aside>

        {/* Main Content Skeleton */}
        <main className="flex-1 space-y-4">
          {/* Header */}
          <div className="h-7 w-64 bg-zinc-800 rounded" />

          {/* Institution Cards / Table Rows */}
          <div className="space-y-3 mt-4">
            {[1, 2, 3, 4].map((i) => (
              <div
                key={i}
                className="h-20 rounded-xl bg-zinc-900 border border-zinc-800"
              />
            ))}
          </div>
        </main>

      </div>
    </div>
  );
}