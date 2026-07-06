export default function TeacherDashboardLoading() {
  return (
    <div className="flex h-screen overflow-hidden bg-[#0f0f0f] animate-pulse">
      {/* ================= SIDEBAR ================= */}
      <div className="w-64 border-r border-white/10 p-4 space-y-4">
        {Array.from({ length: 7 }).map((_, i) => (
          <div
            key={i}
            className="h-10 bg-white/10 rounded-lg"
          />
        ))}
      </div>

      {/* ================= MAIN CONTENT ================= */}
      <main className="flex-1 overflow-y-auto px-10 py-10 space-y-8">
        {/* Hero title */}
        <div className="space-y-3">
          <div className="h-8 w-64 bg-white/10 rounded-lg" />
          <div className="h-4 w-96 bg-white/10 rounded-lg" />
        </div>

        {/* Stat cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {Array.from({ length: 3 }).map((_, i) => (
            <div
              key={i}
              className="h-32 bg-white/10 rounded-xl"
            />
          ))}
        </div>

        {/* Main dashboard section */}
        <div className="space-y-4">
          <div className="h-6 w-48 bg-white/10 rounded" />
          <div className="h-64 bg-white/10 rounded-xl" />
        </div>
      </main>
    </div>
  );
}