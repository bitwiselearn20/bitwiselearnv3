import Image from "next/image";

const keyFeatures = [
  "Industry-Vetted Curriculum",
  "1:1 Mentorship",
  "Placement Support Every Step of the Way",
];

export function About() {
  return (
    <section className="py-16 sm:py-24">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 items-center gap-10 lg:grid-cols-2">
          <div className="relative order-2 min-h-[320px] overflow-hidden rounded-2xl bg-neutral-900 shadow-lg sm:min-h-[400px] lg:order-1">
            <Image
              src="https://images.unsplash.com/photo-1509062522246-3755977927d7?w=900&h=700&fit=crop"
              alt="Students learning together"
              fill
              className="object-cover"
              sizes="(max-width: 1024px) 100vw, 50vw"
            />
            <div className="absolute inset-0 bg-linear-to-t from-black/60 via-transparent to-transparent" />
            <div className="absolute bottom-6 left-6 rounded-2xl bg-neutral-900/90 px-5 py-4 backdrop-blur">
              <p className="title-react text-3xl font-bold text-white">95%</p>
              <p className="text-sm text-neutral-300">Placement success rate</p>
            </div>
          </div>
          <div className="order-1 lg:order-2">
            <h2 className="title-react mb-4 text-4xl font-bold text-black sm:text-5xl">
              Bridging the Gap Between Academia &amp; Industry
            </h2>
            <p className="mb-6 text-lg text-neutral-600">
              We&apos;re not just another ed-tech platform. BitwiseLearn is a career accelerator built to turn classroom theory into skills recruiters actually pay for.
            </p>
            <ul className="space-y-3">
              {keyFeatures.map((f) => (
                <li key={f} className="flex items-center gap-3">
                  <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-black text-white">
                    <svg className="h-3.5 w-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={3} d="M5 13l4 4L19 7" />
                    </svg>
                  </span>
                  <span className="font-medium text-black">{f}</span>
                </li>
              ))}
            </ul>
          </div>
        </div>
      </div>
    </section>
  );
}
