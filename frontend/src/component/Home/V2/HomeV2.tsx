import { Hero } from "@/component/Home/V2/Hero";
import { Companies } from "@/component/Home/V2/Companies";
import { WhyBitwise } from "@/component/Home/V2/WhyBitwise";
import { HowItWorks } from "@/component/Home/V2/HowItWorks";
import { Testimonials } from "@/component/Home/V2/Testimonials";
import { CTA } from "@/component/Home/V2/CTA";
import { FAQ } from "@/component/Home/V2/FAQ";
import { Newsletter } from "@/component/Home/V2/Newsletter";
import { ScrollReveal } from "@/component/Home/V2/ScrollReveal";
import Footer from "@/component/general/Footer";
import { Navbar } from "@/component/general/Navbar";

export default function HomeV2() {
  return (
    <>
      <ScrollReveal variant="fade-up">
        <Navbar />
      </ScrollReveal>
      <ScrollReveal variant="fade-up">
        <Hero />
      </ScrollReveal>
      <ScrollReveal variant="fade-up">
        <Companies />
      </ScrollReveal>
      <ScrollReveal variant="fade-up">
        <WhyBitwise />
      </ScrollReveal>
      <ScrollReveal variant="fade-up">
        <HowItWorks />
      </ScrollReveal>
      <ScrollReveal variant="fade-up">
        <Testimonials />
      </ScrollReveal>
      <ScrollReveal variant="fade-up">
        <CTA />
      </ScrollReveal>
      <ScrollReveal variant="fade-up">
        <FAQ />
      </ScrollReveal>
      <ScrollReveal variant="fade-up">
        <Newsletter />
      </ScrollReveal>
      <ScrollReveal variant="fade-up">
        <Footer />
      </ScrollReveal>
    </>
  );
}
