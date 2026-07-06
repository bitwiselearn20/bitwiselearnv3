export default function CourseV2Loading() {
  return (
    <div className="min-h-screen bg-[#0f0f0f] p-4 animate-pulse">
      <div className="flex gap-4 h-[calc(100vh-2rem)]">
        {/* ================= LEFT SIDEBAR ================= */}
        <aside className="w-[320px] bg-[#161616] rounded-lg p-4 space-y-4">
          {/* Mode Switch */}
          <div className="flex gap-3 justify-center">
            <div className="h-8 w-24 rounded bg-white/10" />
            <div className="h-8 w-24 rounded bg-white/10" />
          </div>

          {/* Sections */}
          {[1, 2, 3].map((i) => (
            <div
              key={i}
              className="rounded-lg bg-[#1e1e1e] p-3 space-y-2"
            >
              <div className="h-4 w-3/4 bg-white/10 rounded" />
              <div className="h-3 w-full bg-white/5 rounded" />
              <div className="h-3 w-5/6 bg-white/5 rounded" />
            </div>
          ))}
        </aside>

        {/* ================= RIGHT CONTENT ================= */}
        <div className="flex-1 bg-[#161616] rounded-xl flex flex-col overflow-hidden">
          {/* HEADER */}
          <div className="p-6 flex justify-between items-center">
            <div className="h-6 w-64 bg-white/10 rounded" />

            <div className="flex gap-3">
              <div className="h-10 w-28 bg-white/10 rounded-lg" />
              <div className="h-10 w-28 bg-white/10 rounded-lg" />
              <div className="h-10 w-32 bg-white/10 rounded-lg" />
            </div>
          </div>

          {/* CONTENT */}
          <div className="flex-1 p-6">
            <div className="flex gap-6 h-full">
              {/* VIDEO + TRANSCRIPT */}
              <div className="flex-1 flex flex-col gap-6">
                {/* Video */}
                <div className="aspect-video rounded-xl bg-black/30 flex items-center justify-center">
                  <div className="h-12 w-12 rounded-full bg-white/10" />
                </div>

                {/* Transcript */}
                <div className="flex-1 rounded-xl bg-[#1e1e1e] p-4 space-y-3">
                  <div className="h-4 w-40 bg-white/10 rounded" />
                  <div className="h-3 w-full bg-white/5 rounded" />
                  <div className="h-3 w-5/6 bg-white/5 rounded" />
                  <div className="h-3 w-4/6 bg-white/5 rounded" />
                </div>
              </div>

              {/* PDF PANEL */}
              <div className="w-[40%] rounded-xl bg-black/20 flex items-center justify-center">
                <div className="h-10 w-10 rounded bg-white/10" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}