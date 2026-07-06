export default function Loading() {
  return (
    <div className="flex h-screen w-full bg-[#0f0f0f] text-white animate-pulse">
      {/* Sidebar Skeleton */}
      <div className="w-60 border-r border-white/10 p-4 flex flex-col gap-4">
        <div className="h-8 w-32 rounded bg-white/10" />
        <div className="h-4 w-full rounded bg-white/10" />
        <div className="h-4 w-5/6 rounded bg-white/10" />
        <div className="h-4 w-4/6 rounded bg-white/10" />
        <div className="h-4 w-3/6 rounded bg-white/10" />
      </div>

      {/* Main Compiler Area */}
      <div className="flex-1 p-6 flex flex-col gap-4">
        {/* Top Bar */}
        <div className="flex justify-between items-center">
          <div className="h-6 w-40 rounded bg-white/10" />
          <div className="flex gap-3">
            <div className="h-9 w-20 rounded bg-white/10" />
            <div className="h-9 w-24 rounded bg-white/10" />
          </div>
        </div>

        {/* Editor */}
        <div className="flex-1 rounded-xl bg-white/5 border border-white/10" />

        {/* Output Section */}
        <div className="h-40 rounded-xl bg-white/5 border border-white/10" />
      </div>
    </div>
  );
}