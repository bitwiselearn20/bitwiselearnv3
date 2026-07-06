export default function AllCoursesLoading() {
  return (
    <div className="flex h-screen bg-[#0f0f0f] text-white animate-pulse">
      {/* ================= SIDEBAR ================= */}
      <aside className="w-64 bg-[#121212] border-r border-white/10" />

      {/* ================= MAIN ================= */}
      <main className="flex-1 p-6 overflow-y-auto">
        {/* HEADER */}
        <header className="flex items-center gap-10 mb-6">
          {/* Search */}
          <div className="w-1/2 h-10 bg-[#1e1e1e] rounded-lg" />

          {/* Dropdown */}
          <div className="w-36 h-10 bg-[#1e1e1e] rounded-xl" />
        </header>

        {/* BREADCRUMB */}
        <div className="flex items-center gap-3 mb-6">
          <div className="h-5 w-28 bg-white/10 rounded" />
          <div className="h-8 w-6 bg-white/5 rounded" />
          <div className="h-7 w-32 bg-white/10 rounded" />
        </div>

        {/* COURSES GRID */}
        <section className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {Array.from({ length: 6 }).map((_, i) => (
            <div
              key={i}
              className="rounded-xl p-4 bg-[#161616] flex flex-col gap-4"
            >
              {/* Thumbnail */}
              <div className="h-40 rounded-lg bg-[#0f0f0f]" />

              {/* Title */}
              <div className="h-5 w-3/4 bg-white/10 rounded" />

              {/* Level + Duration */}
              <div className="flex justify-between">
                <div className="h-4 w-20 bg-white/10 rounded" />
                <div className="h-4 w-16 bg-white/10 rounded" />
              </div>

              {/* Description */}
              <div className="space-y-2">
                <div className="h-4 bg-white/5 rounded" />
                <div className="h-4 w-5/6 bg-white/5 rounded" />
              </div>

              {/* Instructor */}
              <div className="flex justify-end items-center gap-2 mt-auto">
                <div className="w-7 h-7 rounded-full bg-white/10" />
                <div className="h-4 w-20 bg-white/10 rounded" />
              </div>
            </div>
          ))}
        </section>
      </main>
    </div>
  );
}