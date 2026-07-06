import { Hero } from "@/component/Home/V2/Hero";
import { ValueProposition } from "@/component/Home/V2/ValueProposition";
import { DashboardPreview } from "@/component/Home/V2/DashboardPreview";
import { About } from "@/component/Home/V2/About";
import { Certifications } from "@/component/Home/V2/Certifications";
import { LearningRoadmap } from "@/component/Home/V2/LearningRoadmap";
import { Teams } from "@/component/Home/V2/Teams";
import { Testimonials } from "@/component/Home/V2/Testimonials";
import { CTA } from "@/component/Home/V2/CTA";
import { ScrollReveal } from "@/component/Home/V2/ScrollReveal";
import Footer from "@/component/general/Footer";
import { Navbar } from "@/component/general/Navbar";

export default function HomeV2() {
  return (
    <div className="landing-light-theme min-h-screen bg-white">
      <ScrollReveal variant="fade-up">
        <Navbar theme="light" />
      </ScrollReveal>
      <ScrollReveal variant="fade-up">
        <Hero />
      </ScrollReveal>
      <ScrollReveal variant="fade-up">
        <ValueProposition />
      </ScrollReveal>
      <ScrollReveal variant="fade-up">
        <DashboardPreview />
      </ScrollReveal>
      <ScrollReveal variant="fade-up">
        <About />
      </ScrollReveal>
      <ScrollReveal variant="fade-up">
        <Certifications />
      </ScrollReveal>
      <ScrollReveal variant="fade-up">
        <LearningRoadmap />
      </ScrollReveal>
      <ScrollReveal variant="fade-up">
        <Teams />
      </ScrollReveal>
      <ScrollReveal variant="fade-up">
        <Testimonials />
      </ScrollReveal>
      <ScrollReveal variant="fade-up">
        <CTA />
      </ScrollReveal>
      <ScrollReveal variant="fade-up">
        <Footer />
      </ScrollReveal>
    </div>
  );
}
