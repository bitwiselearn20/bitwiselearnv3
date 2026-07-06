export default function LoadingInstitutions() {
  return (
    <div className="flex min-h-screen bg-[#0f0f0f] animate-pulse">
      {/* Sidebar */}
      <div className="h-screen w-64 bg-white/5 border-r border-white/10" />

      {/* Main Content */}
      <div className="ml-10 mt-10 w-full">
        {/* Header */}
        <div className="w-[80%] mx-auto mb-6 flex justify-between items-center">
          <div className="h-8 w-56 rounded bg-white/10" />
          <div className="h-11 w-44 rounded-xl bg-white/10" />
        </div>

        {/* Filter */}
        <div className="w-[80%] mx-auto mb-6">
          <div className="h-12 w-full rounded-xl bg-white/10" />
        </div>

        {/* Institutions List */}
        <div className="w-[80%] mx-auto space-y-4">
          {[1, 2, 3, 4, 5].map((i) => (
            <div
              key={i}
              className="flex items-center justify-between rounded-xl border border-white/10 bg-white/5 px-6 py-4"
            >
              <div className="space-y-2">
                <div className="h-5 w-64 rounded bg-white/10" />
                <div className="h-4 w-40 rounded bg-white/10" />
              </div>
              <div className="h-9 w-28 rounded-lg bg-white/10" />
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}