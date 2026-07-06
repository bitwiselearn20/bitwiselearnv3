export default function Loading() {
  return (
    <div className="flex bg-bg min-h-screen animate-pulse">
      {/* SIDEBAR */}
      <div className="h-screen w-64 bg-divBg p-6 space-y-6">
        <div className="h-8 w-32 bg-white/10 rounded-md" />
        {[...Array(6)].map((_, i) => (
          <div key={i} className="h-4 w-full bg-white/10 rounded-md" />
        ))}
      </div>

      {/* MAIN CONTENT */}
      <div className="ml-10 mt-10 w-full">
        {/* HEADER */}
        <div className="w-[80%] mx-auto mb-6 flex justify-between items-center">
          <div className="h-8 w-56 bg-white/10 rounded-md" />
          <div className="h-10 w-40 bg-primaryBlue/40 rounded-xl" />
        </div>

        {/* FILTER BAR */}
        <div className="w-[80%] mx-auto mb-6">
          <div className="h-12 w-full bg-white/10 rounded-xl" />
        </div>

        {/* DASHBOARD CARDS / TABLE */}
        <div className="w-[80%] mx-auto space-y-4">
          {[...Array(5)].map((_, i) => (
            <div
              key={i}
              className="h-20 w-full bg-divBg rounded-xl flex items-center px-6 gap-6"
            >
              <div className="h-10 w-10 bg-white/10 rounded-full" />
              <div className="flex-1 space-y-2">
                <div className="h-4 w-1/3 bg-white/10 rounded-md" />
                <div className="h-3 w-1/4 bg-white/10 rounded-md" />
              </div>
              <div className="h-8 w-20 bg-white/10 rounded-lg" />
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}