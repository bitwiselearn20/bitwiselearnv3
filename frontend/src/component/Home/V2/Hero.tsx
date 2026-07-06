"use client";

import Image from "next/image";
import { useRef, useEffect, useState } from "react";
import { useInView } from "@/component/Home/V2/useInView";

const SAMPLE_VIDEO_SRC = "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerBlazes.mp4";

const reviewAvatars = [
  "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=80&h=80&fit=crop&crop=face",
  "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=80&h=80&fit=crop&crop=face",
  "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=80&h=80&fit=crop&crop=face",
];

const TYPING_WORDS = ["Coding", "Placements", "Interviews", "Your Career"];

function useTypewriter(words: string[]) {
  const [text, setText] = useState("");
  const [wordIndex, setWordIndex] = useState(0);
  const [deleting, setDeleting] = useState(false);

  useEffect(() => {
    const current = words[wordIndex % words.length];
    const speed = deleting ? 45 : 90;
    const pauseAtFull = 1400;
    const pauseAtEmpty = 300;

    if (!deleting && text === current) {
      const t = setTimeout(() => setDeleting(true), pauseAtFull);
      return () => clearTimeout(t);
    }
    if (deleting && text === "") {
      const t = setTimeout(() => {
        setDeleting(false);
        setWordIndex((i) => (i + 1) % words.length);
      }, pauseAtEmpty);
      return () => clearTimeout(t);
    }
    const t = setTimeout(() => {
      setText(current.slice(0, deleting ? text.length - 1 : text.length + 1));
    }, speed);
    return () => clearTimeout(t);
  }, [text, deleting, wordIndex, words]);

  return text;
}

export function Hero() {
  const { ref: topRef, inView: topInView } = useInView({ threshold: 0.2 });
  const { ref: cardsRef, inView: cardsInView } = useInView({ threshold: 0.15 });
  const videoRef = useRef<HTMLVideoElement>(null);
  const hasAutoplayedRef = useRef(false);
  const [isPlaying, setIsPlaying] = useState(false);
  const typedWord = useTypewriter(TYPING_WORDS);

  useEffect(() => {
    if (!cardsInView || hasAutoplayedRef.current || !videoRef.current) return;
    const timer = setTimeout(() => {
      videoRef.current?.play().then(() => {
        hasAutoplayedRef.current = true;
        setIsPlaying(true);
      }).catch(() => {});
    }, 800);
    return () => clearTimeout(timer);
  }, [cardsInView]);

  function toggleVideo() {
    if (!videoRef.current) return;
    if (isPlaying) {
      videoRef.current.pause();
      setIsPlaying(false);
    } else {
      videoRef.current.play();
      setIsPlaying(true);
    }
  }

  return (
    <section className="relative overflow-hidden py-12 sm:py-20">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        {/* Top: avatars + 125k+ student reviews */}
        <div
          ref={topRef}
          className={`flex flex-col items-center text-center scroll-reveal ${topInView ? "is-visible" : ""}`}
        >
          <div className="flex -space-x-3">
            {reviewAvatars.map((src, i) => (
              <div
                key={i}
                className="relative h-12 w-12 overflow-hidden rounded-full border-2 border-white ring-1 ring-neutral-200 sm:h-14 sm:w-14"
              >
                <Image
                  src={src}
                  alt=""
                  fill
                  className="object-cover"
                  sizes="48px"
                />
              </div>
            ))}
          </div>
          <p className="mt-3 text-sm font-medium uppercase tracking-wider text-black">
            Transforming tech education
          </p>
        </div>

        {/* Headline + description + CTA */}
        <div
          className={`mt-10 text-center scroll-reveal scroll-stagger-1 ${topInView ? "is-visible" : ""}`}
        >
          <h1 className="title-react text-4xl font-bold tracking-tight text-black sm:text-5xl md:text-6xl lg:text-7xl">
            Master the Art of
            <br />
            <span className="text-black">
              {typedWord}
              <span className="animate-pulse text-neutral-400">|</span>
            </span>
          </h1>
          <p className="mx-auto mt-6 max-w-2xl text-lg text-neutral-600 sm:text-xl">
            Bridge the gap between academic theory and industry reality. Join the growing community of learners turning classroom knowledge into real, job-ready skill.
          </p>
          <div className="mt-8 flex flex-col items-center justify-center gap-3 sm:flex-row">
            <a
              href="/multi-login"
              className="rounded-full bg-black px-7 py-3 font-medium text-white hover:bg-neutral-800"
            >
              Start Learning
            </a>
            <a
              href="#platform-preview"
              className="rounded-full border border-neutral-300 px-7 py-3 font-medium text-black hover:border-black"
            >
              Watch Demo
            </a>
          </div>
        </div>

        {/* Three separate tiles: left portrait | center video (bigger) | right portrait */}
        <div
          ref={cardsRef}
          className="mt-16 grid grid-cols-1 gap-6 md:grid-cols-5"
        >
          {/* Left tile – portrait */}
          <div
            className={`md:col-span-1 overflow-hidden rounded-2xl shadow-lg transition-all duration-600 ease-out ${
              cardsInView ? "opacity-100 translate-x-0" : "opacity-0 -translate-x-6"
            } transform`}
            style={{ transitionDelay: "0ms" }}
          >
            <div className="relative aspect-3/4 overflow-hidden rounded-2xl">
              <Image
                src="https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=400&h=533&fit=crop&crop=face"
                alt=""
                fill
                className="object-cover"
                sizes="(max-width: 768px) 100vw, 33vw"
              />
              <div className="absolute inset-0 bg-linear-to-t from-black/85 via-black/30 to-transparent" />
              <div className="absolute bottom-5 left-5 right-5 text-white drop-shadow">
                <p className="text-xl font-bold sm:text-2xl">95% Placement</p>
                <p className="text-lg font-medium opacity-95">Success Rate</p>
              </div>
            </div>
          </div>

          {/* Center tile – video (bigger, rectangular) */}
          <div
            className={`md:col-span-3 overflow-hidden rounded-2xl shadow-lg transition-all duration-700 ease-out tansform ${
              cardsInView ? "opacity-100 translate-y-0" : "opacity-0 -translate-y-8"
            }`}
            style={{ transitionDelay: "100ms" }}
          >
            <div className="relative aspect-[2.2/1] w-full overflow-hidden rounded-2xl bg-neutral-100">
              <video
                ref={videoRef}
                src={SAMPLE_VIDEO_SRC}
                className="absolute inset-0 h-full w-full object-cover"
                muted
                playsInline
                loop
                onPlay={() => setIsPlaying(true)}
                onPause={() => setIsPlaying(false)}
              />
              <div className="absolute inset-0 bg-linear-to-t from-black/80 via-black/20 to-transparent" />
              {!isPlaying && (
                <button
                  type="button"
                  onClick={toggleVideo}
                  className="absolute left-1/2 top-1/2 flex h-16 w-16 -translate-x-1/2 -translate-y-1/2 items-center justify-center rounded-full border-2 border-white bg-white/20 hover:bg-white/30"
                  aria-label="Play video"
                >
                  <svg className="ml-1 h-8 w-8 text-white" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M8 5v14l11-7z" />
                  </svg>
                </button>
              )}
              <div className="absolute bottom-5 left-5">
                <p className="text-2xl font-bold text-white sm:text-3xl">See BitwiseLearn in action</p>
                <p className="text-base text-white/90">A two-minute tour of the platform</p>
              </div>
              <button
                type="button"
                onClick={toggleVideo}
                className="absolute bottom-4 right-4 flex h-10 w-10 items-center justify-center rounded-full bg-white/20 hover:bg-white/30"
                aria-label={isPlaying ? "Pause" : "Play"}
              >
                {isPlaying ? (
                  <svg className="h-5 w-5 text-white" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M6 4h4v16H6V4zm8 0h4v16h-4V4z" />
                  </svg>
                ) : (
                  <svg className="ml-0.5 h-5 w-5 text-white" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M8 5v14l11-7z" />
                  </svg>
                )}
              </button>
            </div>
          </div>

          {/* Right tile – portrait */}
          <div
            className={`md:col-span-1 overflow-hidden rounded-2xl shadow-lg transition-all duration-600 ease-out transform ${
              cardsInView ? "opacity-100 translate-x-0" : "opacity-0 translate-x-6"
            }`}
            style={{ transitionDelay: "0ms" }}
          >
            <div className="relative aspect-3/4 overflow-hidden rounded-2xl">
              <Image
                src="https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=400&h=533&fit=crop&crop=face"
                alt=""
                fill
                className="object-cover"
                sizes="(max-width: 768px) 100vw, 33vw"
              />
              <div className="absolute inset-0 bg-linear-to-t from-black/80 via-transparent to-transparent" />
              <p className="absolute bottom-5 left-5 right-14 text-lg font-bold text-white drop-shadow sm:text-xl">
                100+ industry-experienced trainers
              </p>
              <div className="absolute bottom-4 right-4 flex h-9 w-9 items-center justify-center rounded-lg bg-white/20">
                <svg className="h-4 w-4 text-white" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M6 4h4v16H6V4zm8 0h4v16h-4V4z" />
                </svg>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
