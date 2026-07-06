export default function InstitutionDashboardLoading() {
  return (
    <div className="flex h-screen overflow-hidden bg-[#0f0f0f] animate-pulse">
      {/* ================= SIDEBAR ================= */}
      <aside className="w-64 bg-[#121212] border-r border-white/10" />

      {/* ================= MAIN ================= */}
      <main className="flex-1 overflow-y-auto px-10 py-10">
        <div className="space-y-6">
          {/* Title */}
          <div className="h-8 w-1/3 bg-white/10 rounded" />

          {/* Subtitle */}
          <div className="h-4 w-1/2 bg-white/5 rounded" />

          {/* Stat Cards */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-8">
            {Array.from({ length: 3 }).map((_, i) => (
              <div
                key={i}
                className="h-32 rounded-xl bg-[#161616]"
              />
            ))}
          </div>

          {/* Main Content Block */}
          <div className="mt-10 space-y-4">
            <div className="h-6 w-40 bg-white/10 rounded" />
            <div className="h-48 rounded-xl bg-[#161616]" />
          </div>
        </div>
      </main>
    </div>
  );
}