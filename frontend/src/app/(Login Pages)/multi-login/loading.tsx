export default function MultiLoginLoading() {
  return (
    <div className="bg-bg min-h-screen w-screen flex flex-col md:flex-row animate-pulse">
      {/* LEFT */}
      <div className="flex-1 flex flex-col px-6 py-10 md:p-16">
        {/* BRAND */}
        <div className="h-8 w-44 bg-white/10 rounded mb-10" />

        {/* LOGIN CARD */}
        <div className="w-full md:w-[60%] bg-divBg rounded-3xl p-8 space-y-6">
          {/* TITLE */}
          <div className="h-6 w-24 bg-white/10 rounded" />

          {/* ROLE SELECT */}
          <div className="flex gap-2 bg-bg p-2 rounded-xl">
            <div className="h-9 flex-1 bg-white/10 rounded-lg" />
            <div className="h-9 flex-1 bg-white/10 rounded-lg" />
            <div className="h-9 flex-1 bg-white/10 rounded-lg" />
          </div>

          {/* EMAIL */}
          <div className="space-y-2">
            <div className="h-4 w-32 bg-white/10 rounded" />
            <div className="h-10 w-full bg-white/10 rounded-lg" />
          </div>

          {/* PASSWORD */}
          <div className="space-y-2">
            <div className="h-4 w-32 bg-white/10 rounded" />
            <div className="h-10 w-full bg-white/10 rounded-lg" />
          </div>

          {/* REMEMBER + FORGOT */}
          <div className="flex justify-between">
            <div className="h-4 w-28 bg-white/10 rounded" />
            <div className="h-4 w-32 bg-white/10 rounded" />
          </div>

          {/* SUBMIT */}
          <div className="h-12 w-full bg-white/10 rounded-lg" />
        </div>
      </div>

      {/* RIGHT IMAGE */}
      <div className="hidden lg:block lg:w-[38%] lg:h-screen bg-white/5" />
    </div>
  );
}
