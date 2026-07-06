export default function Loading() {
  return (
    <div className="h-screen flex flex-col bg-black animate-pulse">
      {/* TOP BAR */}
      <div className="flex justify-end items-center px-4 py-3 border-b border-white/10">
        <div className="h-8 w-20 bg-white/10 rounded-md" />
      </div>

      {/* MAIN CONTENT */}
      <div className="flex-1 grid grid-cols-2 gap-4 p-4 min-h-0">
        {/* LEFT SECTION */}
        <div className="flex flex-col gap-4 h-full rounded-xl bg-white/5 p-4">
          {/* Section Header */}
          <div className="flex justify-between items-center">
            <div className="h-5 w-32 bg-white/10 rounded-md" />
            <div className="h-4 w-24 bg-white/10 rounded-md" />
          </div>

          {/* Question Content */}
          <div className="flex-1 space-y-3 mt-4">
            <div className="h-4 w-full bg-white/10 rounded-md" />
            <div className="h-4 w-[90%] bg-white/10 rounded-md" />
            <div className="h-4 w-[85%] bg-white/10 rounded-md" />
            <div className="h-4 w-[70%] bg-white/10 rounded-md" />
          </div>

          {/* Navigation */}
          <div className="flex justify-between mt-4">
            <div className="h-10 w-24 bg-white/10 rounded-md" />
            <div className="h-10 w-24 bg-white/10 rounded-md" />
          </div>
        </div>

        {/* RIGHT SECTION */}
        <div className="flex flex-col gap-4 h-full rounded-xl bg-white/5 p-4">
          {/* Assignment Title */}
          <div className="h-5 w-48 bg-white/10 rounded-md" />

          {/* Options / Editor */}
          <div className="flex-1 space-y-4 mt-4">
            {[1, 2, 3, 4].map((i) => (
              <div
                key={i}
                className="h-12 w-full bg-white/10 rounded-lg"
              />
            ))}
          </div>

          {/* Action Buttons */}
          <div className="flex justify-between mt-4">
            <div className="h-10 w-32 bg-white/10 rounded-md" />
            <div className="h-10 w-32 bg-white/10 rounded-md" />
          </div>
        </div>
      </div>
    </div>
  );
}