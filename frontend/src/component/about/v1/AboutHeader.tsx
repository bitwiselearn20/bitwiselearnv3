export function AboutHeader() {
  return (
    <section className="border-b border-neutral-800 bg-black/40 py-16 sm:py-20">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <span className="inline-block rounded-full border border-neutral-600 px-4 py-1.5 text-center text-xs font-medium uppercase tracking-wider text-neutral-400">
          Bitwise-Learn focus
        </span>
        <h1 className="title-react mt-4 text-center text-3xl font-bold text-white sm:text-4xl">
          Bitwise-Learn is built for placement & industry-ready courses
        </h1>
        <p className="mx-auto mt-4 max-w-3xl text-center text-lg text-neutral-400">
          We combine learning technology with placement support and industry-aligned curricula so students don’t just complete courses—they become job-ready and get the right opportunities.
        </p>
        <div className="mx-auto mt-8 max-w-3xl rounded-2xl border border-neutral-700 bg-neutral-900/60 px-6 py-5 sm:px-8 sm:py-6">
          <h2 className="text-sm font-semibold uppercase tracking-wider text-neutral-300">
            About our modules
          </h2>
          <p className="mt-3 text-neutral-400">
            Our learning is modular: choose <strong className="text-neutral-300">placement modules</strong> (resume, interviews, drives), <strong className="text-neutral-300">skill-based course modules</strong> (DSA, full-stack, cloud, soft skills), and <strong className="text-neutral-300">assessment modules</strong> to build a path that fits your goals. Each module is designed to be industry-aligned and placement-ready, so institutes and learners can mix and match what they need.
          </p>
          <ul className="mt-4 space-y-2 text-sm text-neutral-400">
            <li className="flex items-center gap-2">
              <span className="h-1.5 w-1.5 shrink-0 rounded-full bg-neutral-500" aria-hidden />
              Placement modules: resume builder, mock interviews, company drives, aptitude prep
            </li>
            <li className="flex items-center gap-2">
              <span className="h-1.5 w-1.5 shrink-0 rounded-full bg-neutral-500" aria-hidden />
              Skill modules: DSA, programming languages, full-stack, cloud, DevOps, soft skills
            </li>
            <li className="flex items-center gap-2">
              <span className="h-1.5 w-1.5 shrink-0 rounded-full bg-neutral-500" aria-hidden />
              Assessment modules: quizzes, coding tests, project reviews, certification exams
            </li>
          </ul>
        </div>
      </div>
    </section>
  );
}
