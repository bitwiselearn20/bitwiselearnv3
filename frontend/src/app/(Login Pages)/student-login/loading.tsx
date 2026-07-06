export default function Loading() {
  return (
    <div className="bg-bg min-h-screen w-screen flex flex-col md:flex-row animate-pulse">
      {/* LEFT */}
      <div className="flex-1 flex flex-col px-6 py-10 md:p-16">
        {/* Logo */}
        <div className="h-8 w-40 bg-white/10 rounded-md mb-10" />

        {/* Card */}
        <div className="relative w-full md:w-[60%] bg-divBg rounded-3xl p-8">
          {/* Title */}
          <div className="h-6 w-28 bg-white/10 rounded-md mb-8" />

          {/* Email */}
          <div className="space-y-2 mb-6">
            <div className="h-4 w-32 bg-white/10 rounded-md" />
            <div className="h-11 w-full bg-white/10 rounded-lg" />
          </div>

          {/* Password */}
          <div className="space-y-2 mb-8">
            <div className="h-4 w-32 bg-white/10 rounded-md" />
            <div className="h-11 w-full bg-white/10 rounded-lg" />
          </div>

          {/* Button */}
          <div className="h-12 w-full bg-primaryBlue/40 rounded-lg" />

          {/* Forgot password */}
          <div className="h-4 w-32 bg-white/10 rounded-md mt-6 mx-auto" />
        </div>

        {/* Typewriter text */}
        <div className="h-5 w-80 bg-white/10 rounded-md mt-10" />
      </div>

      {/* RIGHT IMAGE */}
      <div className="hidden lg:block lg:w-[38%] relative">
        <div className="absolute inset-0 bg-white/10" />
      </div>
    </div>
  );
}