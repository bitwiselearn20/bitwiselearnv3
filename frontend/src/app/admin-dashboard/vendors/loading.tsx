export default function Loading() {
  return (
    <div className="flex bg-bg min-h-screen animate-pulse">
      {/* Sidebar */}
      <div className="h-screen">
        <div className="w-20 h-full bg-white/5" />
      </div>

      {/* Main Content */}
      <div className="ml-10 mt-10 w-full">
        {/* Header */}
        <div className="w-[80%] mx-auto mb-6 flex justify-between items-center">
          <div className="h-8 w-56 bg-white/10 rounded-md" />
          <div className="h-10 w-40 bg-primaryBlue/30 rounded-xl" />
        </div>

        {/* Filter */}
        <div className="w-[80%] mx-auto mb-6">
          <div className="h-12 w-full bg-white/10 rounded-lg" />
        </div>

        {/* Dashboard Cards / Table */}
        <div className="w-[80%] mx-auto space-y-4">
          {[1, 2, 3, 4, 5].map((i) => (
            <div
              key={i}
              className="h-20 w-full rounded-xl bg-white/10"
            />
          ))}
        </div>
      </div>
    </div>
  );
}