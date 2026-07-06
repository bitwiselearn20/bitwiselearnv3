const HIGHLIGHTS = [
  {
    title: "Placement-focused learning",
    description:
      "Bitwise-Learn is built to make you placement-ready. From resume-building and interview prep to company-specific mock rounds and placement drives, we align learning with what recruiters look for.",
    icon: "briefcase",
  },
  {
    title: "Industry-ready courses",
    description:
      "Courses are designed with industry inputs: real projects, live coding, and skills that match current job roles. Learn what companies use today—DSA, full-stack, cloud, DevOps, and soft skills.",
    icon: "code",
  },
  {
    title: "Outcomes that matter",
    description:
      "Track progress with assessments, certifications, and placement support. Institutes and learners get visibility into readiness and outcomes so education leads to employability.",
    icon: "chart",
  },
  {
    title: "Certifications & credentials",
    description:
      "Earn industry-recognized certificates and micro-credentials that signal readiness to employers. Our assessment modules are designed to validate skills in line with job requirements.",
    icon: "certificate",
  },
  {
    title: "Live practice & labs",
    description:
      "Hands-on coding labs, live doubt-solving, and practice environments so you can apply concepts in real time. Build a portfolio of projects that demonstrate your capabilities.",
    icon: "beaker",
  },
  {
    title: "Analytics for institutes",
    description:
      "Institutes get dashboards on batch performance, placement rates, and skill gaps. Use data to refine curriculum, identify at-risk learners, and improve outcomes at scale.",
    icon: "academic",
  },
];

function Icon({ name }: { name: string }) {
  const className = "h-6 w-6 shrink-0 text-white";
  switch (name) {
    case "briefcase":
      return (
        <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.62-9-1.745M16 6V4a2 2 0 00-2-2h-4a2 2 0 00-2 2v2m4 6h.01M5 20h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
        </svg>
      );
    case "code":
      return (
        <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
        </svg>
      );
    case "chart":
      return (
        <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
        </svg>
      );
    case "certificate":
      return (
        <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z" />
        </svg>
      );
    case "beaker":
      return (
        <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19.428 15.428a2 2 0 00-1.022-.547l-2.387-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 10.172V5L8 4z" />
        </svg>
      );
    case "academic":
      return (
        <svg className={className} fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 14l9-5-9-5-9 5 9 5z" />
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 14l6.16-3.422a12.083 12.083 0 01.665 6.479A11.952 11.952 0 0012 20.055a11.952 11.952 0 00-6.824-2.998 12.078 12.078 0 01.665-6.479L12 14z" />
        </svg>
      );
    default:
      return null;
  }
}

const STATS = [
  { value: "50+", label: "Partner institutes" },
  { value: "125k+", label: "Learners onboarded" },
  { value: "500+", label: "Industry-ready courses & modules" },
  { value: "92%", label: "Placement outcome rate" },
];

export function EdTechFocus() {
  return (
    <section className="py-16 sm:py-24">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="mb-16 grid grid-cols-2 gap-6 rounded-2xl border border-neutral-700 bg-neutral-900/60 px-6 py-8 sm:grid-cols-4">
          {STATS.map((stat) => (
            <div key={stat.label} className="text-center">
              <p className="text-2xl font-bold text-white sm:text-3xl">{stat.value}</p>
              <p className="mt-1 text-sm text-neutral-400">{stat.label}</p>
            </div>
          ))}
        </div>
        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {HIGHLIGHTS.map((item) => (
            <article
              key={item.title}
              className="rounded-2xl border border-neutral-700 bg-neutral-900 p-6 transition-colors hover:border-neutral-600"
            >
              <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-neutral-800">
                <Icon name={item.icon} />
              </div>
              <h3 className="title-react mt-4 text-xl font-semibold text-white">{item.title}</h3>
              <p className="mt-2 text-sm leading-relaxed text-neutral-400">{item.description}</p>
            </article>
          ))}
        </div>
      </div>
    </section>
  );
}
