export default function LoadingInstitution() {
  return (
    <div className="min-h-screen bg-[#0f0f0f] p-6 animate-pulse">
      <div className="flex gap-6 max-w-screen mx-auto">
        {/* Institution Sidebar */}
        <aside className="w-72 rounded-2xl border border-white/10 bg-white/5 p-5 space-y-4">
          <div className="h-6 w-3/4 rounded bg-white/10" />
          <div className="h-4 w-1/2 rounded bg-white/10" />
          <div className="h-px bg-white/10 my-2" />
          <div className="h-10 w-full rounded-xl bg-white/10" />
          <div className="h-10 w-full rounded-xl bg-white/10" />
        </aside>

        {/* Main Content */}
        <main className="flex-1 space-y-6">
          {/* Tabs */}
          <div className="flex gap-3">
            {[1, 2, 3].map((i) => (
              <div
                key={i}
                className="h-10 w-28 rounded-xl bg-white/10"
              />
            ))}
          </div>

          {/* Entity List */}
          <div className="rounded-xl border border-white/10 bg-white/5">
            {[1, 2, 3, 4].map((i) => (
              <div
                key={i}
                className="flex items-center justify-between px-5 py-4 border-b border-white/10 last:border-none"
              >
                <div className="space-y-2">
                  <div className="h-4 w-52 rounded bg-white/10" />
                  <div className="h-3 w-32 rounded bg-white/10" />
                </div>
                <div className="h-8 w-24 rounded bg-white/10" />
              </div>
            ))}
          </div>
        </main>
      </div>
    </div>
  );
}