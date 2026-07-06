const certifications = [
  "AWS Solutions Architect",
  "Azure Fundamentals",
  "GCP Cloud Engineer",
  "DevOps Engineering",
  "Full-Stack Development",
  "Data Engineering",
];

export function Certifications() {
  return (
    <section className="bg-neutral-50 py-16 sm:py-20">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <h2 className="title-react mb-3 text-center text-4xl font-bold text-black sm:text-5xl">
          Industry-Level Certifications
        </h2>
        <p className="mx-auto mb-10 max-w-2xl text-center text-lg text-neutral-600 sm:text-xl">
          Certification tracks mapped to the roles recruiters are actually hiring for.
        </p>
        <div className="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-6">
          {certifications.map((cert) => (
            <div
              key={cert}
              className="flex h-24 items-center justify-center rounded-2xl border border-neutral-200 bg-white p-4 text-center text-sm font-medium text-neutral-700 shadow-sm"
            >
              {cert}
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
