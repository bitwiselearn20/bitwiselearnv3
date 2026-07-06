export default function CourseStudyModeLoading() {
  return (
    <div className="relative min-h-screen bg-[#121313] p-4 flex flex-col gap-4 overflow-hidden animate-pulse">
      {/* TOP BAR */}
      <div className="flex items-center justify-between gap-5 rounded-xl p-1 flex-wrap">
        <div className="flex items-center gap-3">
          <div className="h-8 w-10 rounded bg-white/10" />
          <div className="h-6 w-52 rounded bg-white/10" />
        </div>

        <div className="h-8 w-24 rounded bg-white/10" />
      </div>

      {/* MAIN CONTENT */}
      <div className="flex flex-col lg:flex-row gap-4 flex-1">
        {/* VIDEO SECTION */}
        <div className="flex-3 bg-[#1E1E1E] rounded-xl p-4 flex items-center justify-center">
          <div className="w-full h-full rounded-lg bg-black/30 flex items-center justify-center">
            <div className="h-12 w-12 rounded-full bg-white/10" />
          </div>
        </div>

        {/* NOTES PANEL */}
        <div className="flex-[1.3] bg-[#1E1E1E] rounded-xl p-4 flex flex-col gap-3">
          <div className="h-5 w-32 rounded bg-white/10" />
          <div className="h-4 w-full rounded bg-white/10" />
          <div className="h-4 w-5/6 rounded bg-white/10" />
          <div className="h-4 w-4/6 rounded bg-white/10" />
        </div>
      </div>

      {/* META SECTION */}
      <div className="w-full lg:w-2/3 bg-[#1E1E1E] rounded-xl p-4 flex flex-col gap-3">
        <div className="h-5 w-40 rounded bg-white/10" />
        <div className="h-4 w-full rounded bg-white/10" />
        <div className="h-4 w-5/6 rounded bg-white/10" />
        <div className="h-4 w-4/6 rounded bg-white/10" />
      </div>
    </div>
  );
}