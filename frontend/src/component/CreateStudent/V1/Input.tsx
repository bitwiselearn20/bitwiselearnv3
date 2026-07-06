interface InputProps {
  label: string;
  name: string;
  placeholder?: string;
  value?: string;
  required?: boolean;
  isInvalid?: boolean;
  showError?: boolean;
  errorMessage?: string;
  onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void;
}

export function Input({
  label,
  required,
  isInvalid,
  showError,
  errorMessage,
  value,
  ...props
}: InputProps) {
  const isEmpty = required && !value?.trim();
  const shouldShowError = (isEmpty && showError) || isInvalid;

  return (
    <div className="space-y-1">
      <label className="block text-xs font-medium text-white/70">{label}</label>

      <input
        {...props}
        value={value}
        className={`w-full rounded-lg border px-3 py-2 text-sm text-white
        bg-white/5 placeholder:text-white/30 transition-all
        focus:outline-none focus:scale-[1.01]
        ${
          shouldShowError
            ? "border-red-500/50 focus:border-red-400"
            : "border-white/10 focus:border-sky-500"
        }`}
      />

      {shouldShowError && errorMessage && (
        <p className="text-[11px] text-red-400 animate-fade-in">
          {errorMessage}
        </p>
      )}
    </div>
  );
}
