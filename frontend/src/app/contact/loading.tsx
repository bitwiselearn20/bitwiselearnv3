export default function ContactLoading() {
  return (
    <div className="animate-pulse text-white max-w-7xl mx-auto px-6 pt-24">
      <div className="h-12 w-48 bg-white/10 rounded-xl mx-auto mb-6" />
      <div className="h-4 w-72 bg-white/10 rounded-lg mx-auto mb-12" />

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-16">
        <div className="space-y-6 max-w-md">
          <div className="h-4 w-20 bg-white/10 rounded" />
          <div className="h-12 w-full bg-white/10 rounded-xl" />

          <div className="h-4 w-20 bg-white/10 rounded" />
          <div className="h-12 w-full bg-white/10 rounded-xl" />

          <div className="h-4 w-24 bg-white/10 rounded" />
          <div className="h-24 w-full bg-white/10 rounded-xl" />

          <div className="h-12 w-40 bg-white/10 rounded-xl mx-auto mt-6" />
        </div>

        <div className="space-y-6">
          <div className="h-6 w-56 bg-white/10 rounded" />

          {[1, 2, 3, 4].map((i) => (
            <div key={i} className="h-4 w-full bg-white/10 rounded" />
          ))}

          <div className="h-px bg-white/10 my-8" />

          <div className="h-4 w-40 bg-white/10 rounded mx-auto" />

          <div className="flex justify-center gap-6 mt-4">
            <div className="h-4 w-40 bg-white/10 rounded" />
            <div className="h-4 w-40 bg-white/10 rounded" />
          </div>
        </div>
      </div>
    </div>
  );
}
