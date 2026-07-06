export default function AboutLoading() {
  return (
    <section className="w-full py-24 bg-[#121313]">
      <div className="max-w-7xl mx-auto px-6">
        {/* Header skeleton */}
        <div className="text-center max-w-5xl mx-auto mb-16">
          <div className="h-10 w-3/4 mx-auto rounded-lg bg-white/10 animate-pulse" />
          <div className="h-5 w-2/3 mx-auto mt-6 rounded-lg bg-white/10 animate-pulse" />
        </div>

        {/* Cards skeleton */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          {[1, 2, 3].map((i) => (
            <div
              key={i}
              className="rounded-3xl bg-[#1E1E1E] border border-white/10 p-8 animate-pulse"
            >
              <div className="flex items-center justify-between">
                <div className="h-14 w-14 rounded-2xl bg-white/10" />
                <div className="h-6 w-20 rounded-full bg-white/10" />
              </div>

              <div className="h-6 w-3/4 mt-8 rounded bg-white/10" />
              <div className="h-4 w-full mt-4 rounded bg-white/10" />
              <div className="h-4 w-5/6 mt-2 rounded bg-white/10" />

              <div className="mt-8 space-y-4">
                {[1, 2, 3, 4].map((j) => (
                  <div key={j} className="flex gap-3">
                    <div className="h-5 w-5 rounded-full bg-white/10" />
                    <div className="h-4 w-full rounded bg-white/10" />
                  </div>
                ))}
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
