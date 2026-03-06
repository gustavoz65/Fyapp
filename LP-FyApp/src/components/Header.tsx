"use client";

import { useEffect, useState } from "react";
import { Sparkles } from "lucide-react";

const APP_URL =
  process.env.NEXT_PUBLIC_APP_URL || "https://fyapp-production.up.railway.app";

export default function Header() {
  const [isScrolled, setIsScrolled] = useState(false);

  useEffect(() => {
    const handleScroll = () => {
      setIsScrolled(window.scrollY > 20);
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
      className={`fixed top-0 left-0 right-0 z-50 transition-all duration-300 ${
        isScrolled
          ? "bg-surface/80 backdrop-blur-xl border-b border-border shadow-lg"
          : "bg-transparent"
      }`}
    >
      <div className="max-w-7xl mx-auto px-6 sm:px-8 lg:px-12">
        <div className="flex justify-between items-center py-4">
          {/* Logo */}
          <a href="#" className="flex items-center gap-3 group no-underline">
            <div className="relative">
              <img
                src="/logoFy.png"
                alt="FyApp"
                className="h-9 w-9 object-contain transition-transform duration-300 group-hover:scale-110"
              />
              <div className="absolute -inset-1 bg-accent/20 rounded-lg blur opacity-0 group-hover:opacity-100 transition-opacity duration-300" />
            </div>
            <span className="text-2xl font-display font-bold text-primary">
              Fy<span className="text-accent">App</span>
            </span>
          </a>

          {/* Desktop Navigation */}
          <nav className="hidden md:flex items-center gap-1">
            <button
              onClick={() => scrollToSection("features")}
              className="px-4 py-2 text-sm font-medium text-text-secondary hover:text-primary rounded-lg hover:bg-background-alt transition-all duration-200 cursor-pointer"
            >
              Recursos
            </button>
            <button
              onClick={() => scrollToSection("how")}
              className="px-4 py-2 text-sm font-medium text-text-secondary hover:text-primary rounded-lg hover:bg-background-alt transition-all duration-200 cursor-pointer"
            >
              Como Funciona
            </button>
            <button
              onClick={() => scrollToSection("security")}
              className="px-4 py-2 text-sm font-medium text-text-secondary hover:text-primary rounded-lg hover:bg-background-alt transition-all duration-200 cursor-pointer"
            >
              Segurança
            </button>

            <div className="ml-4 h-6 w-px bg-border" />

            <a
              href={`${APP_URL}/login`}
              className="ml-4 group inline-flex items-center gap-2 px-6 py-2.5 bg-accent text-white rounded-lg text-sm font-semibold shadow-md shadow-accent/20 transition-all duration-300 hover:bg-accent-dark hover:-translate-y-0.5 hover:shadow-lg hover:shadow-accent/30"
            >
              <Sparkles className="w-4 h-4 transition-transform group-hover:rotate-12" />
              Começar Grátis
            </a>
          </nav>

          {/* Mobile CTA */}
          <a
            href={`${APP_URL}/login`}
            className="md:hidden inline-flex items-center gap-2 px-5 py-2.5 bg-accent text-white rounded-lg text-sm font-semibold shadow-md shadow-accent/20"
          >
            <Sparkles className="w-4 h-4" />
            Entrar
          </a>
        </div>
      </div>
    </header>
  );
}
