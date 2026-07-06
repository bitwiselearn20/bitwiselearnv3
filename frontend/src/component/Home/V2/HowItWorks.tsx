import Image from "next/image";

const steps = [
  {
    icon: "target",
    title: "Choose a course",
    description: "Pick from career-focused learning programs.",
  },
  {
    icon: "book",
    title: "Learn through real lessons",
    description: "High-quality videos, resources, quizzes & projects.",
  },
  {
    icon: "certificate",
    title: "Get certified",
    description: "Earn certificates to prove your skills and move forward.",
  },
  {
    icon: "rocket",
    title: "Apply your skills",
    description: "Earn certificates to prove your skills and move forward.",
  },
];

function StepIcon({ name }: { name: string }) {
  const className = "h-6 w-6 shrink-0 text-white";
  switch (name) {
    case "target":
      return (
        <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      );
    case "book":
      return (
        <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
        </svg>
      );
    case "certificate":
      return (
        <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
        </svg>
      );
    case "rocket":
      return (
        <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
        </svg>
      );
    default:
      return null;
  }
}

export function HowItWorks() {
  return (
    <section className="py-16 sm:py-24">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <h2 className="title-react mb-3 text-center text-4xl font-bold text-white sm:text-5xl">
          How to begin learning with Bitwise Learn
        </h2>
        <p className="mx-auto mb-12 max-w-2xl text-center text-lg text-neutral-400 sm:text-xl">
          Pick from career-focused learning programs. High-quality videos, resources, quizzes & projects.
        </p>

        <div className="grid grid-cols-1 gap-4 lg:grid-cols-[1fr_2.2fr_1fr] lg:grid-rows-2 lg:gap-6">
          {/* Row 1 left */}
          <div className="rounded-2xl bg-neutral-900 p-6 shadow-lg lg:row-start-1">
            <StepIcon name={steps[0].icon} />
            <h4 className="title-react mt-4 text-xl font-bold text-white sm:text-2xl">{steps[0].title}</h4>
            <p className="mt-3 text-base leading-relaxed text-neutral-400">{steps[0].description}</p>
          </div>
          {/* Center: large image tile (spans both rows) */}
          <div className="relative min-h-[320px] overflow-hidden rounded-2xl bg-neutral-800 shadow-lg sm:min-h-[400px] lg:col-start-2 lg:row-span-2 lg:row-start-1 lg:min-h-[500px]">
            <Image
              src="https://images.unsplash.com/photo-1522202176988-66273c2fd55f?w=800&h=600&fit=crop"
              alt="Learning with Bitwise Learn"
              fill
              className="object-cover"
              sizes="(max-width: 1024px) 100vw, 45vw"
            />
          </div>
          {/* Row 1 right */}
          <div className="rounded-2xl bg-neutral-900 p-6 shadow-lg lg:row-start-1">
            <StepIcon name={steps[2].icon} />
            <h4 className="title-react mt-4 text-xl font-bold text-white sm:text-2xl">{steps[2].title}</h4>
            <p className="mt-3 text-base leading-relaxed text-neutral-400">{steps[2].description}</p>
          </div>
          {/* Row 2 left */}
          <div className="rounded-2xl bg-neutral-900 p-6 shadow-lg lg:row-start-2">
            <StepIcon name={steps[1].icon} />
            <h4 className="title-react mt-4 text-xl font-bold text-white sm:text-2xl">{steps[1].title}</h4>
            <p className="mt-3 text-base leading-relaxed text-neutral-400">{steps[1].description}</p>
          </div>
          {/* Row 2 right (col 3) */}
          <div className="rounded-2xl bg-neutral-900 p-6 shadow-lg lg:col-start-3 lg:row-start-2">
            <StepIcon name={steps[3].icon} />
            <h4 className="title-react mt-4 text-xl font-bold text-white sm:text-2xl">{steps[3].title}</h4>
            <p className="mt-3 text-base leading-relaxed text-neutral-400">{steps[3].description}</p>
          </div>
        </div>
      </div>
    </section>
  );
}
