export default function AdminLoginLoading() {
  return (
    <div className="bg-bg min-h-screen w-screen flex flex-col md:flex-row animate-pulse">
      {/* LEFT SECTION */}
      <div className="flex-1 flex flex-col px-6 py-10 md:p-16">
        {/* LOGO */}
        <div className="h-8 w-40 bg-white/10 rounded mb-10" />

        {/* LOGIN CARD */}
        <div className="w-full md:w-[60%] bg-divBg mt-10 md:mt-16 md:ml-16 rounded-3xl p-8 space-y-6">
          <div className="h-6 w-24 bg-white/10 rounded" />

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

          {/* REMEMBER */}
          <div className="h-4 w-40 bg-white/10 rounded" />

          {/* BUTTON */}
          <div className="h-12 w-full bg-white/10 rounded-lg" />

          {/* FORGOT */}
          <div className="h-4 w-32 bg-white/10 rounded mx-auto" />
        </div>

        {/* TYPEWRITER PLACEHOLDER */}
        <div className="mt-10 h-5 w-56 bg-white/10 rounded mx-auto md:mx-0" />
      </div>

      {/* RIGHT IMAGE SKELETON */}
      <div className="hidden lg:block lg:w-[38%] lg:h-screen bg-white/5" />
    </div>
  );
}
