export function CTA() {
  return (
    <section className="border-y border-neutral-200 bg-white py-16 sm:py-24">
      <div className="mx-auto max-w-3xl px-4 text-center sm:px-6 lg:px-8">
        <h2 className="title-react text-4xl font-bold text-black sm:text-5xl">
          Ready to Launch Your Tech Career?
        </h2>
        <p className="mx-auto mt-5 max-w-xl text-lg text-neutral-600">
          Join thousands of students who have transformed their lives with BitwiseLearn.
        </p>
        <div className="mt-8 flex flex-col items-center justify-center gap-3 sm:flex-row">
          <a
            href="/multi-login"
            className="rounded-full bg-black px-7 py-3 font-medium text-white hover:bg-neutral-800"
          >
            Get Started Now
          </a>
          <a
            href="/contact"
            className="rounded-full border border-neutral-300 px-7 py-3 font-medium text-black hover:border-black"
          >
            Contact Sales
          </a>
        </div>
      </div>
    </section>
  );
}
