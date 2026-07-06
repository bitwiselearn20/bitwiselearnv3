export default function LoadingBatchPage() {
  return (
    <div className="min-h-screen bg-[#0f0f0f] p-6 animate-pulse">
      <div className="flex gap-6 max-w-screen mx-auto">
        {/* Sidebar */}
        <aside className="w-72 rounded-xl bg-white/5 border border-white/10 p-4 space-y-4">
          <div className="h-6 w-3/4 rounded bg-white/10" />
          <div className="h-4 w-1/2 rounded bg-white/10" />
          <div className="h-px bg-white/10 my-3" />
          <div className="space-y-2">
            <div className="h-4 w-full rounded bg-white/10" />
            <div className="h-4 w-5/6 rounded bg-white/10" />
            <div className="h-4 w-4/6 rounded bg-white/10" />
          </div>
        </aside>

        {/* Main content */}
        <main className="flex-1 space-y-6">
          {/* Tabs */}
          <div className="flex gap-3">
            {[1, 2, 3].map((i) => (
              <div
                key={i}
                className="h-10 w-28 rounded-md bg-white/10"
              />
            ))}
          </div>

          {/* Entity list */}
          <div className="rounded-xl border border-white/10 bg-white/5">
            {[1, 2, 3, 4, 5].map((i) => (
              <div
                key={i}
                className="flex items-center justify-between px-4 py-4 border-b border-white/10 last:border-none"
              >
                <div className="space-y-2">
                  <div className="h-4 w-48 rounded bg-white/10" />
                  <div className="h-3 w-32 rounded bg-white/10" />
                </div>

                <div className="h-8 w-20 rounded bg-white/10" />
              </div>
            ))}
          </div>
        </main>
      </div>
    </div>
  );
}