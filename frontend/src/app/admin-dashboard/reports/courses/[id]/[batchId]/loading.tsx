export default function Loading() {
  return (
    <div className="min-h-screen bg-zinc-950 p-6 animate-pulse">
      <div className="max-w-7xl mx-auto space-y-6">

        {/* Back button + title */}
        <div className="h-6 w-48 bg-zinc-800 rounded" />
        <div className="flex items-center gap-3">
          <div className="h-8 w-1 bg-zinc-800 rounded-full" />
          <div className="h-7 w-64 bg-zinc-800 rounded" />
        </div>

        {/* Stats */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          {[1, 2, 3].map((i) => (
            <div
              key={i}
              className="h-24 rounded-lg border border-zinc-800 bg-zinc-900 p-5"
            />
          ))}
        </div>

        {/* Charts */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div className="h-90 rounded-lg border border-zinc-800 bg-zinc-900" />
          <div className="h-90 rounded-lg border border-zinc-800 bg-zinc-900" />
        </div>

        {/* Table */}
        <div className="rounded-lg border border-zinc-800 bg-zinc-900 overflow-hidden">
          <div className="h-14 border-b border-zinc-800 bg-zinc-900" />
          {[1, 2, 3, 4, 5].map((i) => (
            <div
              key={i}
              className="h-12 border-t border-zinc-800 bg-zinc-900"
            />
          ))}
        </div>

      </div>
    </div>
  );
}