export default function Loading() {
  return (
    <div className="flex gap-6 h-screen bg-zinc-950 text-zinc-200 animate-pulse">

      {/* Left Course Card */}
      <aside className="w-[320px] ml-4 mt-4 shrink-0 border border-zinc-800 bg-zinc-900 rounded-xl overflow-hidden h-fit">
        <div className="h-40 w-full bg-zinc-800" />

        <div className="p-5 space-y-4 mt-4">
          <div className="space-y-2">
            <div className="h-5 w-3/4 bg-zinc-800 rounded" />
            <div className="h-4 w-1/2 bg-zinc-800 rounded" />
          </div>

          <div className="h-6 w-24 bg-zinc-800 rounded-md" />

          <div className="space-y-3">
            {[1, 2, 3].map((i) => (
              <div key={i} className="space-y-1">
                <div className="h-4 w-24 bg-zinc-800 rounded" />
                <div className="h-4 w-full bg-zinc-800 rounded" />
              </div>
            ))}
          </div>
        </div>
      </aside>

      {/* Right Section */}
      <section className="flex-1 mt-4">
        <div className="mx-auto w-full max-w-[90%] space-y-4">

          {/* Filters */}
          <div className="flex items-center justify-between gap-3 rounded-lg p-4 border border-zinc-800 bg-zinc-900">
            <div className="flex gap-3">
              <div className="h-9 w-64 bg-zinc-800 rounded-md" />
              <div className="h-9 w-40 bg-zinc-800 rounded-md" />
            </div>
            <div className="h-4 w-32 bg-zinc-800 rounded" />
          </div>

          {/* Table */}
          <div className="border border-zinc-800 rounded-xl overflow-hidden">
            <div className="h-12 bg-zinc-900 border-b border-zinc-800" />

            {[1, 2, 3, 4, 5].map((i) => (
              <div
                key={i}
                className="h-12 bg-zinc-900 border-b border-zinc-800"
              />
            ))}
          </div>

        </div>
      </section>
    </div>
  );
}