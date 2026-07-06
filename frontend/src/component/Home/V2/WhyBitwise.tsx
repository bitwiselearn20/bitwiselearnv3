import Image from "next/image";

const reasons = [
  {
    icon: "person",
    title: "Expert instructors",
    description:
      "Our instructors are creative minds and strategic thinkers.",
  },
  {
    icon: "money",
    title: "Affordable pricing",
    description: "We're a team of creative minds and strategic thinkers.",
  },
  {
    icon: "trophy",
    title: "Awards",
    description: "But along the way, our work has been honored.",
  },
  {
    icon: "star",
    title: "Reviews",
    description: "Strategic placements for testimonials, student success.",
  },
];

function Icon({ name }: { name: string }) {
  const className = "h-6 w-6 shrink-0 text-white";
  switch (name) {
    case "person":
      return (
        <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
        </svg>
      );
    case "money":
      return (
        <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      );
    case "trophy":
      return (
        <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 21h8M12 17v4M7 4h10M7 4a2 2 0 00-2 2v2a2 2 0 002 2h10a2 2 0 002-2V6a2 2 0 00-2-2M7 4V3h10v1M9 7h6m-4 0v4m4-4v4" />
        </svg>
      );
    case "star":
      return (
        <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
        </svg>
      );
    default:
      return null;
  }
}

export function WhyBitwise() {
  return (
    <section className="py-16 sm:py-24">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <h2 className="title-react mb-3 text-center text-4xl font-bold text-white sm:text-5xl">
          Why choose Bitwise Learn
        </h2>
        <p className="mx-auto mb-12 max-w-2xl text-center text-lg text-neutral-400 sm:text-xl">
          Designed for better learning. Built for real success. Designed for better learning. Built for real success.
        </p>

        <div className="grid grid-cols-1 gap-4 lg:grid-cols-[1fr_2.2fr_1fr] lg:grid-rows-2 lg:gap-6">
          {/* Row 1 left – light grey card, lifted */}
          <div className="rounded-2xl bg-neutral-900 p-6 shadow-lg lg:row-start-1">
            <Icon name={reasons[0].icon} />
            <h4 className="title-react mt-4 text-xl font-bold text-white sm:text-2xl">{reasons[0].title}</h4>
            <p className="mt-3 text-base leading-relaxed text-neutral-400">{reasons[0].description}</p>
          </div>
          {/* Center: large image tile (spans both rows), rounded, no padding */}
          <div className="relative min-h-[320px] overflow-hidden rounded-2xl bg-neutral-800 shadow-lg sm:min-h-[400px] lg:col-start-2 lg:row-span-2 lg:row-start-1 lg:min-h-[500px]">
            <Image
              src="https://images.unsplash.com/photo-1498050108023-c5249f4df085?w=800&h=600&fit=crop"
              alt="Person using laptop"
              fill
              className="object-cover"
              sizes="(max-width: 1024px) 100vw, 45vw"
            />
          </div>
          {/* Row 1 right */}
          <div className="rounded-2xl bg-neutral-900 p-6 shadow-lg lg:row-start-1">
            <Icon name={reasons[2].icon} />
            <h4 className="title-react mt-4 text-xl font-bold text-white sm:text-2xl">{reasons[2].title}</h4>
            <p className="mt-3 text-base leading-relaxed text-neutral-400">{reasons[2].description}</p>
          </div>
          {/* Row 2 left */}
          <div className="rounded-2xl bg-neutral-900 p-6 shadow-lg lg:row-start-2">
            <Icon name={reasons[1].icon} />
            <h4 className="title-react mt-4 text-xl font-bold text-white sm:text-2xl">{reasons[1].title}</h4>
            <p className="mt-3 text-base leading-relaxed text-neutral-400">{reasons[1].description}</p>
          </div>
          {/* Row 2 right (col 3) */}
          <div className="rounded-2xl bg-neutral-900 p-6 shadow-lg lg:col-start-3 lg:row-start-2">
            <Icon name={reasons[3].icon} />
            <h4 className="title-react mt-4 text-xl font-bold text-white sm:text-2xl">{reasons[3].title}</h4>
            <p className="mt-3 text-base leading-relaxed text-neutral-400">{reasons[3].description}</p>
          </div>
        </div>
      </div>
    </section>
  );
}
