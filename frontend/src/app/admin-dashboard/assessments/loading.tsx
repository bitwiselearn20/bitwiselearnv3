export default function Loading() {
  return (
    <section className="flex w-full flex-col gap-6 p-4 animate-pulse">
      {/* HEADER */}
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        {/* Search */}
        <div className="h-10 w-full sm:w-96 rounded-md bg-white/10" />

        {/* Add button */}
        <div className="h-10 w-36 rounded-md bg-primaryBlue/40" />
      </div>

      {/* GRID */}
      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {Array.from({ length: 6 }).map((_, i) => (
          <div
            key={i}
            className="rounded-xl p-4 flex flex-col gap-4 bg-divBg border border-white/10"
          >
            <div className="flex justify-between gap-4">
              <div className="h-5 w-2/3 rounded bg-white/10" />
              <div className="h-5 w-16 rounded bg-white/10" />
            </div>

            <div className="h-4 w-full rounded bg-white/10" />
            <div className="h-4 w-5/6 rounded bg-white/10" />

            <div className="h-3 w-3/4 rounded bg-white/10" />

            <div className="h-9 w-full rounded bg-white/10 mt-auto" />
          </div>
        ))}
      </div>
    </section>
  );
}
