const UNIVERSITIES = [
  "Panjab University",
  "GNDU Amritsar",
  "Thapar Institute",
  "Anna University",
  "VIT Vellore",
  "SRM Institute",
  "Osmania University",
  "JNTUH",
  "University of Hyderabad",
  "Andhra University",
  "JNTUK",
  "Sri Venkateswara University",
  "Bangalore University",
  "VTU",
  "Manipal Academy",
  "Christ University",
];

export function Companies() {
  return (
    <section className="border-y border-neutral-800 bg-neutral-900 py-12 overflow-hidden">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <p className="mb-8 text-center text-sm font-medium text-neutral-400">
          Trusted by universities across Punjab, Tamil Nadu, Telangana, Andhra Pradesh & Karnataka
        </p>
        <div className="relative flex w-full items-center">
          <div className="companies-marquee flex w-max items-center gap-8">
            {[...UNIVERSITIES, ...UNIVERSITIES].map((name, i) => (
              <div
                key={`${name}-${i}`}
                className="flex shrink-0 h-12 min-w-[8rem] max-w-[11rem] px-4 items-center justify-center rounded-lg bg-neutral-800 text-xs font-medium text-neutral-300 shadow-sm border border-neutral-700 text-center"
              >
                {name}
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
