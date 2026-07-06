const teams = [
  {
    title: "Training Team",
    description: "IIT and NIT alumni who design and deliver every coding, aptitude, and interview-prep session.",
  },
  {
    title: "Platform Development Team",
    description: "Engineers building the coding sandbox, assessments, and analytics that power every dashboard.",
  },
  {
    title: "Content Development Team",
    description: "Curriculum specialists keeping every course, question bank, and roadmap current with what recruiters ask for.",
  },
];

export function Teams() {
  return (
    <section className="bg-neutral-50 py-16 sm:py-24">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <h2 className="title-react mb-3 text-center text-4xl font-bold text-black sm:text-5xl">
          Our Teams
        </h2>
        <p className="mx-auto mb-12 max-w-2xl text-center text-lg text-neutral-600 sm:text-xl">
          The people behind the platform, the curriculum, and the coaching.
        </p>
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-3">
          {teams.map((team) => (
            <div
              key={team.title}
              className="rounded-2xl border border-neutral-200 bg-white p-8 text-center shadow-lg"
            >
              <h4 className="title-react text-xl font-bold text-black sm:text-2xl">
                {team.title}
              </h4>
              <p className="mt-3 text-sm leading-relaxed text-neutral-600">
                {team.description}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
