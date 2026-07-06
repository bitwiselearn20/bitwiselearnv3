export default function LoadingAllProblems() {
  return (
    <div className="relative flex gap-4 h-screen bg-[#0f0f0f] animate-pulse">
      {/* SIDEBAR */}
      <div className="w-65 h-screen bg-white/5 border-r border-white/10" />

      {/* MAIN CONTENT */}
      <div className="flex-1 p-6 space-y-6">
        {/* DASHBOARD HERO */}
        <div className="rounded-xl border border-white/10 bg-white/5 p-6 space-y-4">
          <div className="h-6 w-48 bg-white/10 rounded-md" />
          <div className="h-4 w-72 bg-white/10 rounded-md" />
        </div>

        {/* FILTER / ACTION BAR */}
        <div className="flex justify-between items-center">
          <div className="h-10 w-64 bg-white/10 rounded-lg" />
          <div className="h-10 w-36 bg-white/10 rounded-lg" />
        </div>

        {/* PROBLEMS LIST */}
        <div className="space-y-4">
          {[...Array(6)].map((_, i) => (
            <div
              key={i}
              className="h-20 rounded-xl border border-white/10 bg-white/5"
            />
          ))}
        </div>
      </div>
    </div>
  );
}