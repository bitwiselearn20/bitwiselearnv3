export default function LoadingAssessmentBuilder() {
  return (
    <div className="w-full animate-pulse">
      {/* Header */}
      <div className="flex items-center justify-between px-1 py-4">
        <div className="h-6 w-64 rounded bg-white/10" />
        <div className="flex gap-3">
          <div className="h-10 w-32 rounded bg-white/10" />
          <div className="h-10 w-24 rounded bg-white/10" />
        </div>
      </div>

      <div className="mb-8 h-px w-full bg-white/10" />

      {/* Sections */}
      <div className="space-y-4">
        {[1, 2, 3].map((i) => (
          <div
            key={i}
            className="rounded-xl border border-white/10 bg-white/5"
          >
            {/* Section Header */}
            <div className="flex items-center justify-between p-5">
              <div>
                <div className="h-5 w-48 rounded bg-white/10 mb-2" />
                <div className="h-4 w-32 rounded bg-white/10" />
              </div>

              <div className="flex gap-2">
                <div className="h-9 w-28 rounded bg-white/10" />
                <div className="h-9 w-16 rounded bg-white/10" />
                <div className="h-9 w-9 rounded-full bg-white/10" />
              </div>
            </div>

            {/* Fake expanded content */}
            <div className="border-t border-white/10 px-5 py-4 space-y-3">
              {[1, 2].map((q) => (
                <div
                  key={q}
                  className="rounded-lg border border-white/10 bg-white/5 p-4 space-y-3"
                >
                  <div className="h-4 w-3/4 rounded bg-white/10" />
                  <div className="grid grid-cols-2 gap-2">
                    <div className="h-8 rounded bg-white/10" />
                    <div className="h-8 rounded bg-white/10" />
                    <div className="h-8 rounded bg-white/10" />
                    <div className="h-8 rounded bg-white/10" />
                  </div>
                  <div className="h-3 w-24 rounded bg-white/10" />
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>

      {/* Floating delete button */}
      <div className="fixed bottom-6 right-6 h-12 w-44 rounded-full bg-white/10" />
    </div>
  );
}
