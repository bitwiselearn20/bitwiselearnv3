export default function AssignmentLoading() {
  return (
    <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3 animate-pulse">
      {Array.from({ length: 6 }).map((_, i) => (
        <div
          key={i}
          className="
            rounded-2xl p-5
            bg-white/5 border border-white/10
            flex flex-col gap-4
          "
        >
          {/* Title */}
          <div className="h-5 w-2/3 rounded bg-white/10" />

          {/* Description */}
          <div className="h-4 w-full rounded bg-white/10" />
          <div className="h-4 w-5/6 rounded bg-white/10" />

          {/* Meta info */}
          <div className="flex justify-between mt-2">
            <div className="h-4 w-20 rounded bg-white/10" />
            <div className="h-4 w-16 rounded bg-white/10" />
          </div>

          {/* Action button */}
          <div className="h-9 w-full rounded-lg bg-white/10 mt-3" />
        </div>
      ))}
    </div>
  );
}