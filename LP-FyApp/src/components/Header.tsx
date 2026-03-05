"use client";

import { useEffect, useState } from "react";

const APP_URL =
  process.env.NEXT_PUBLIC_APP_URL || "https://fyapp-production.up.railway.app";

export default function Header() {
  const [isScrolled, setIsScrolled] = useState(false);

  useEffect(() => {
    const handleScroll = () => {
      setIsScrolled(window.scrollY > 50);
    };
    window.addEventListener("scroll", handleScroll);
    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

  const scrollToSection = (id: string) => {
    const element = document.getElementById(id);
    if (element) {
      element.scrollIntoView({ behavior: "smooth" });
    }
  };

  return (
    <header
      className={`fixed top-0 left-0 right-0 z-50 bg-white/95 backdrop-blur-md border-b border-[#E5DDD0] py-4 transition-all duration-300 ${
        isScrolled ? "shadow-lg" : ""
      }`}
    >
      <div className="max-w-7xl mx-auto px-8 flex justify-between items-center">
        <a href="#" className="flex items-center gap-3 no-underline">
          <img
            src="/logoFy.png"
            alt="FyApp"
            className="h-8 w-8 object-contain"
          />
          <span className="text-2xl font-bold text-[#7E8C54]">FyApp</span>
        </a>

        <nav className="hidden md:flex gap-12 items-center">
          <button
            onClick={() => scrollToSection("features")}
            className="text-[#2C2C2C] text-sm font-medium transition-all duration-300 relative group cursor-pointer"
          >
            Recursos
            <span className="absolute bottom-0 left-0 w-0 h-0.5 bg-[#7E8C54] group-hover:w-full transition-all duration-300" />
          </button>
          <button
            onClick={() => scrollToSection("how")}
            className="text-[#2C2C2C] text-sm font-medium transition-all duration-300 relative group cursor-pointer"
          >
            Como Funciona
            <span className="absolute bottom-0 left-0 w-0 h-0.5 bg-[#7E8C54] group-hover:w-full transition-all duration-300" />
          </button>
          <button
            onClick={() => scrollToSection("security")}
            className="text-[#2C2C2C] text-sm font-medium transition-all duration-300 relative group cursor-pointer"
          >
            Segurança
            <span className="absolute bottom-0 left-0 w-0 h-0.5 bg-[#7E8C54] group-hover:w-full transition-all duration-300" />
          </button>
          <a
            href={`${APP_URL}/login`}
            className="bg-[#7E8C54] text-white px-6 py-2.5 rounded-lg text-sm font-semibold transition-all duration-300 border-2 border-[#7E8C54] hover:bg-transparent hover:text-[#7E8C54] hover:-translate-y-1 hover:shadow-lg"
          >
            Começar Grátis
          </a>
        </nav>

        {/* Mobile CTA */}
        <a
          href={`${APP_URL}/login`}
          className="md:hidden bg-[#7E8C54] text-white px-4 py-2 rounded-lg text-sm font-semibold"
        >
          Entrar
        </a>
      </div>
    </header>
  );
}
