export default function LoadingAllBatches() {
  return (
    <div className="flex min-h-screen bg-[#0f0f0f] animate-pulse">
      {/* Sidebar */}
      <div className="h-screen w-64 border-r border-white/10 bg-white/5 p-4 space-y-4">
        <div className="h-6 w-3/4 rounded bg-white/10" />
        <div className="h-4 w-2/3 rounded bg-white/10" />
        <div className="h-px bg-white/10 my-3" />
        {[1, 2, 3, 4].map((i) => (
          <div key={i} className="h-4 w-full rounded bg-white/10" />
        ))}
      </div>

      {/* Main Content */}
      <div className="ml-10 mt-10 w-full">
        {/* Header */}
        <div className="w-[80%] mx-auto mb-5 flex justify-between items-center">
          <div className="h-8 w-56 rounded bg-white/10" />
          <div className="h-10 w-36 rounded-xl bg-white/10" />
        </div>

        {/* Filter */}
        <div className="w-[80%] mx-auto mb-6">
          <div className="h-12 w-full rounded-xl bg-white/10" />
        </div>

        {/* Dashboard List */}
        <div className="w-[80%] mx-auto rounded-xl border border-white/10 bg-white/5">
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
      </div>
    </div>
  );
}