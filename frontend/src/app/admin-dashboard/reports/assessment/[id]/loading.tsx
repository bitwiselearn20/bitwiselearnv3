export default function LoadingAssessmentReport() {
  return (
    <div className="p-6 space-y-8 animate-pulse">
      {/* HEADER */}
      <div>
        <div className="h-6 w-48 bg-neutral-800 rounded-md mb-2" />
        <div className="h-4 w-72 bg-neutral-800 rounded-md" />
      </div>

      {/* STATS */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {[...Array(3)].map((_, i) => (
          <div
            key={i}
            className="bg-neutral-900 border border-neutral-800 rounded-lg p-4 space-y-3"
          >
            <div className="h-4 w-32 bg-neutral-800 rounded" />
            <div className="h-8 w-20 bg-neutral-800 rounded" />
          </div>
        ))}
      </div>

      {/* CHARTS */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {[...Array(2)].map((_, i) => (
          <div
            key={i}
            className="bg-neutral-900 border border-neutral-800 rounded-lg p-4"
          >
            <div className="h-64 w-full bg-neutral-800 rounded-md" />
          </div>
        ))}
      </div>

      {/* FILTER BAR */}
      <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
        <div className="h-10 w-full md:w-1/2 bg-neutral-800 rounded-md" />
        <div className="h-10 w-full md:w-1/4 bg-neutral-800 rounded-md" />
      </div>

      {/* TABLE */}
      <div className="bg-neutral-900 border border-neutral-800 rounded-lg overflow-hidden">
        {/* TABLE HEADER */}
        <div className="grid grid-cols-6 gap-4 bg-neutral-800 p-4">
          {[...Array(6)].map((_, i) => (
            <div key={i} className="h-4 bg-neutral-700 rounded" />
          ))}
        </div>

        {/* TABLE ROWS */}
        <div className="divide-y divide-neutral-800">
          {[...Array(6)].map((_, i) => (
            <div key={i} className="grid grid-cols-6 gap-4 p-4">
              {[...Array(6)].map((__, j) => (
                <div
                  key={j}
                  className="h-4 bg-neutral-800 rounded"
                />
              ))}
            </div>
          ))}
        </div>
      </div>

      {/* PAGINATION */}
      <div className="flex justify-between items-center">
        <div className="h-9 w-20 bg-neutral-800 rounded-md" />
        <div className="h-4 w-16 bg-neutral-800 rounded" />
        <div className="h-9 w-20 bg-neutral-800 rounded-md" />
      </div>
    </div>
  );
}