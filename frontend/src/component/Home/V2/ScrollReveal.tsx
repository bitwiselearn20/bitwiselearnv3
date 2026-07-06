"use client";

import { useInView } from "@/component/Home/V2/useInView";

type Variant = "fade-up" | "fade-left" | "fade-right" | "scale";
type Props = {
  children: React.ReactNode;
  variant?: Variant;
  stagger?: number;
  className?: string;
};

const variantClass: Record<Variant, string> = {
  "fade-up": "scroll-reveal",
  "fade-left": "scroll-reveal-left",
  "fade-right": "scroll-reveal-right",
  scale: "scroll-reveal-scale",
};

export function ScrollReveal({
  children,
  variant = "fade-up",
  stagger = 0,
  className = "",
}: Props) {
  const { ref, inView } = useInView({ threshold: 0.05, rootMargin: "0px 0px -80px 0px" });
  const staggerClass = stagger > 0 ? `scroll-stagger-${Math.min(stagger, 6)}` : "";
  const baseClass = variantClass[variant];

  return (
    <div
      ref={ref}
      className={`${baseClass} ${staggerClass} ${inView ? "is-visible" : ""} ${className}`}
    >
      {children}
    </div>
  );
}
