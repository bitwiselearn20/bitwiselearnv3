const pillars = [
  {
    icon: "target",
    title: "Placement Training",
    description: "Structured prep that gets students interview-ready, not just course-complete.",
    tags: ["Mock Interviews", "Resume Building", "Aptitude Tests", "Company-Specific Prep"],
  },
  {
    icon: "certificate",
    title: "Industry Certifications",
    description: "Credentials recruiters actually recognize, built into the learning path.",
    tags: ["AWS Solutions Architect", "GCP Cloud Engineer", "Azure Fundamentals"],
  },
  {
    icon: "link",
    title: "Academic Integration",
    description: "One ecosystem connecting coursework, practice, and placement — not three disconnected tools.",
    tags: ["Unified Ecosystem"],
  },
];

function PillarIcon({ name }: { name: string }) {
  const className = "h-7 w-7 shrink-0 text-black";
  switch (name) {
    case "target":
      return (
        <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      );
    case "certificate":
      return (
        <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z" />
        </svg>
      );
    case "link":
      return (
        <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.828 10.172a4 4 0 010 5.656l-3 3a4 4 0 01-5.656-5.656l1.5-1.5M10.172 13.828a4 4 0 010-5.656l3-3a4 4 0 015.656 5.656l-1.5 1.5" />
        </svg>
      );
    default:
      return null;
  }
}

export function ValueProposition() {
  return (
    <section className="py-16 sm:py-24">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <h2 className="title-react mb-3 text-center text-4xl font-bold text-black sm:text-5xl">
          We Build Careers, Not Just Skills
        </h2>
        <p className="mx-auto mb-12 max-w-2xl text-center text-lg text-neutral-600 sm:text-xl">
          A complete ecosystem designed to turn students into industry-ready professionals.
        </p>
        <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
          {pillars.map((pillar) => (
            <div
              key={pillar.title}
              className="rounded-2xl border border-neutral-200 bg-white p-8 shadow-lg"
            >
              <PillarIcon name={pillar.icon} />
              <h3 className="title-react mt-4 text-xl font-bold text-black sm:text-2xl">
                {pillar.title}
              </h3>
              <p className="mt-3 text-base leading-relaxed text-neutral-600">
                {pillar.description}
              </p>
              <div className="mt-5 flex flex-wrap gap-2">
                {pillar.tags.map((tag) => (
                  <span
                    key={tag}
                    className="rounded-full border border-neutral-300 px-3 py-1 text-xs font-medium text-neutral-700"
                  >
                    {tag}
                  </span>
                ))}
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
