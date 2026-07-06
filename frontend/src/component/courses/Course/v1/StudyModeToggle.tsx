export default function StudyModeToggle({
  enabled,
  onToggle,
}: {
  enabled: boolean;
  onToggle: () => void;
}) {
  return (
    <div className="relative group">
      <button
        onClick={onToggle}
        className={`relative w-10 h-6 rounded-full transition-colors duration-300 cursor-pointer
          ${enabled ? "bg-[#64ACFF]" : "bg-[#2a2a2a]"}`}
      >
        <span
          className={`absolute top-0.5 left-0.5 h-5 w-5 rounded-full bg-white transition-transform duration-300 cursor-pointer
            ${enabled ? "translate-x-4" : "translate-x-0"}`}
        />
      </button>

      <div
        className="absolute -bottom-8 left-1 -translate-x-1/2
        opacity-0 group-hover:opacity-100
        pointer-events-none
        transition-opacity duration-200
        bg-[#1E1E1E] text-white text-xs px-2 py-1 rounded-md whitespace-nowrap"
      >
        Study mode
      </div>
    </div>
  );
}
