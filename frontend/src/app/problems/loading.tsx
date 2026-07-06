export default function AllQuestionsLoading() {
  return (
    <div className="flex h-screen bg-[#0f0f0f] text-white animate-pulse">
      {/* ================= SIDEBAR ================= */}
      <div className="w-64 border-r border-white/10 p-4 space-y-4">
        {Array.from({ length: 6 }).map((_, i) => (
          <div
            key={i}
            className="h-10 bg-white/10 rounded-lg"
          />
        ))}
      </div>

      {/* ================= MAIN ================= */}
      <div className="flex-1 p-4 space-y-6 overflow-hidden">
        {/* Top Section */}
        <div className="flex gap-4">
          {/* Ongoing Courses */}
          <div className="flex-1 bg-white/5 rounded-xl p-4 space-y-4">
            <div className="h-5 w-40 bg-white/10 rounded" />
            {Array.from({ length: 3 }).map((_, i) => (
              <div
                key={i}
                className="h-14 bg-white/10 rounded-lg"
              />
            ))}
          </div>

          {/* Question Info Sidebar */}
          <div className="w-80 bg-white/5 rounded-xl p-4 space-y-4">
            <div className="h-5 w-32 bg-white/10 rounded" />
            {Array.from({ length: 4 }).map((_, i) => (
              <div
                key={i}
                className="h-10 bg-white/10 rounded-lg"
              />
            ))}
          </div>
        </div>

        {/* ================= QUESTIONS LIST ================= */}
        <div className="bg-white/5 rounded-xl p-4 space-y-4">
          <div className="h-6 w-48 bg-white/10 rounded" />

          {Array.from({ length: 6 }).map((_, i) => (
            <div
              key={i}
              className="h-14 bg-white/10 rounded-lg"
            />
          ))}
        </div>
      </div>
    </div>
  );
}