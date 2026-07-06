const modules = [
  {
    icon: "code",
    title: "Programming Languages",
    tagline: "Master the Fundamentals",
    description: "Build a rock-solid base in Python, Java, and JavaScript before moving to advanced topics.",
  },
  {
    icon: "tree",
    title: "Data Structures & Algorithms",
    tagline: "DSA & Problem Solving",
    description: "Structured practice across arrays, trees, graphs, and dynamic programming with real interview questions.",
  },
  {
    icon: "brain",
    title: "Aptitude & Reasoning",
    tagline: "Crack the Aptitude Round",
    description: "Quantitative, logical, and verbal reasoning practice tuned to how top recruiters actually test.",
  },
  {
    icon: "flag",
    title: "Company-Specific Training",
    tagline: "Target Top Recruiters",
    description: "Focused prep tracks modeled on the hiring patterns of the companies your students want to join.",
  },
];

function ModuleIcon({ name }: { name: string }) {
  const className = "h-6 w-6 shrink-0 text-black";
  switch (name) {
    case "code":
      return (
        <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 20l4-16m4 4l4 4-4 4M6 8l-4 4 4 4" />
        </svg>
      );
    case "tree":
      return (
        <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 3v4m0 0a3 3 0 100 6 3 3 0 000-6zm0 6v8m-5-3a3 3 0 106 0m-6 0a3 3 0 11-3-3m9 3a3 3 0 106 0m0 0a3 3 0 10-3-3" />
        </svg>
      );
    case "brain":
      return (
        <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9.663 17h4.673M12 3a4 4 0 00-4 4c0 1.5.5 2-.5 3.5S6 12.5 6 14a6 6 0 006 6 6 6 0 006-6c0-1.5-.5-2-1.5-3.5S16 8.5 16 7a4 4 0 00-4-4z" />
        </svg>
      );
    case "flag":
      return (
        <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 21V4m0 1h14l-2.5 3.5L17 12H3" />
        </svg>
      );
    default:
      return null;
  }
}

export function LearningRoadmap() {
  return (
    <section className="py-16 sm:py-24">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <h2 className="title-react mb-3 text-center text-4xl font-bold text-black sm:text-5xl">
          Comprehensive Learning Roadmap
        </h2>
        <p className="mx-auto mb-4 max-w-2xl text-center text-lg text-neutral-600 sm:text-xl">
          A clear path from fundamentals to placement day — no guesswork about what to study next.
        </p>
        <p className="mx-auto mb-12 max-w-2xl text-center text-sm text-neutral-500">
          Need to move faster? We also offer intensive crash courses for students on a tighter placement timeline.
        </p>
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
          {modules.map((m) => (
            <div
              key={m.title}
              className="rounded-2xl border border-neutral-200 bg-white p-6 shadow-lg"
            >
              <ModuleIcon name={m.icon} />
              <p className="mt-4 text-xs font-semibold uppercase tracking-wider text-neutral-500">
                {m.tagline}
              </p>
              <h4 className="title-react mt-1 text-lg font-bold text-black sm:text-xl">
                {m.title}
              </h4>
              <p className="mt-3 text-sm leading-relaxed text-neutral-600">
                {m.description}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
