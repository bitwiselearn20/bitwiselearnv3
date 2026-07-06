export default function Loading() {
  return (
    <div className="relative flex gap-1 h-screen bg-zinc-950 animate-pulse">

      {/* Sidebar */}
      <div className="w-20 h-full bg-zinc-900 border-r border-zinc-800" />

      {/* Main Content */}
      <div className="w-full p-6">
        {/* Hero */}
        <div className="space-y-4">
          <div className="h-8 w-64 bg-zinc-800 rounded" />
          <div className="h-4 w-96 bg-zinc-800 rounded" />

          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-6">
            {[1, 2, 3].map((i) => (
              <div
                key={i}
                className="h-28 rounded-xl bg-zinc-900 border border-zinc-800"
              />
            ))}
          </div>
        </div>
      </div>

    </div>
  );
}