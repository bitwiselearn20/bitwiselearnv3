export default function VendorDashboardLoading() {
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
        {/* Hero Section */}
        <div className="space-y-3">
          <div className="h-8 w-72 bg-white/10 rounded-lg" />
          <div className="h-4 w-md bg-white/10 rounded-lg" />
        </div>

        {/* Vendor cards / stats */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {Array.from({ length: 3 }).map((_, i) => (
            <div
              key={i}
              className="h-36 bg-white/10 rounded-xl"
            />
          ))}
        </div>

        {/* Main content area */}
        <div className="space-y-4">
          <div className="h-6 w-56 bg-white/10 rounded" />
          <div className="h-72 bg-white/10 rounded-xl" />
        </div>
      </main>
    </div>
  );
}