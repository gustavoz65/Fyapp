import Benefits from "@/components/Benefits";
import CTA from "@/components/CTA";
import Comparison from "@/components/Comparison";
import FAQ from "@/components/FAQ";
import Features from "@/components/Features";
import Footer from "@/components/Footer";
import GrowthSection from "@/components/GrowthSection";
import Header from "@/components/Header";
import Hero from "@/components/Hero";
import HowItWorks from "@/components/HowItWorks";
import ProblemSection from "@/components/ProblemSection";
import SecuritySection from "@/components/SecuritySection";
import Testimonials from "@/components/Testimonials";

export default function LandingPage() {
  return (
    <div className="w-full overflow-x-hidden">
      <Header />
      <Hero />
      <ProblemSection />
      <Features />
      <Comparison />
      <Benefits />
      <Testimonials />
      <GrowthSection />
      <HowItWorks />
      <SecuritySection />
      <FAQ />
      <CTA />
      <Footer />
    </div>
  );
}
