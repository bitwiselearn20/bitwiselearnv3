export function CoursesHeader() {
  return (
    <section className="border-b border-neutral-800 bg-black/40 py-16 sm:py-20">
      {" "}
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        {" "}
        <h1 className="text-center text-3xl font-semibold tracking-tight text-white sm:text-4xl">
          {" "}
          Courses{" "}
        </h1>{" "}
        <p className="mx-auto mt-3 max-w-2xl text-center text-neutral-400">
          {" "}
          Explore our catalog and find the right course for your goals.{" "}
        </p>{" "}
      </div>{" "}
    </section>
  );
}
