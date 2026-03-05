import React from 'react';
import Header from './components/Header';
import Hero from './components/Hero';
import ProblemSection from './components/ProblemSection';
import Features from './components/Features';
import Benefits from './components/Benefits';
import GrowthSection from './components/GrowthSection';
import HowItWorks from './components/HowItWorks';
import SecuritySection from './components/SecuritySection';
import CTA from './components/CTA';
import Footer from './components/Footer';

function App() {
  return (
    <div className="w-full overflow-x-hidden">
      <Header />
      <Hero />
      <ProblemSection />
      <Features />
      <Benefits />
      <GrowthSection />
      <HowItWorks />
      <SecuritySection />
      <CTA />
      <Footer />
    </div>
  );
}

export default App;
