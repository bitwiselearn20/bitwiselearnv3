export default function Loading() {
  return (
    <div className="flex h-screen overflow-hidden bg-bg animate-pulse">
      {/* Sidebar Skeleton */}
      <div className="w-20 h-full bg-white/5" />

      {/* Main Content */}
      <main className="flex-1 overflow-y-auto px-10 py-10 space-y-8">
        {/* Hero Header */}
        <div className="space-y-3">
          <div className="h-8 w-64 bg-white/10 rounded-md" />
          <div className="h-4 w-96 bg-white/10 rounded-md" />
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {[1, 2, 3].map((i) => (
            <div
              key={i}
              className="h-28 rounded-xl bg-white/10"
            />
          ))}
        </div>

        {/* Charts / Large Panels */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div className="h-80 rounded-xl bg-white/10" />
          <div className="h-80 rounded-xl bg-white/10" />
        </div>

        {/* Table / Activity Section */}
        <div className="space-y-4">
          <div className="h-6 w-48 bg-white/10 rounded-md" />
          {[1, 2, 3, 4].map((i) => (
            <div
              key={i}
              className="h-16 rounded-lg bg-white/10"
            />
          ))}
        </div>
      </main>
    </div>
  );
}