export default function Loading() {
  return (
    <section className="flex w-full flex-col gap-8 p-6">
      {/* Search Skeleton */}
      <div className="flex flex-col gap-2 max-w-md animate-pulse">
        <div className="h-10 w-full rounded-lg bg-white/10" />
        <div className="h-6 w-48 rounded bg-white/10" />
        <div className="h-3 w-72 rounded bg-white/10" />
      </div>

      {/* Grid Skeleton */}
      <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {Array.from({ length: 6 }).map((_, i) => (
          <div
            key={i}
            className="
              rounded-2xl p-5 flex flex-col gap-4
              bg-[#1E1E1E]
              border border-white/10
              animate-pulse
            "
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