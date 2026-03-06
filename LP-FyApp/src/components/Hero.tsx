"use client";

import { ArrowRight, TrendingUp, Shield, Zap, Sparkles } from "lucide-react";
import { useEffect, useRef, useState } from "react";

const APP_URL = process.env.NEXT_PUBLIC_APP_URL || "https://fyapp-production.up.railway.app";

export default function Hero() {
  const [mousePosition, setMousePosition] = useState({ x: 0, y: 0 });
  const heroRef = useRef<HTMLElement>(null);

  useEffect(() => {
    const handleMouseMove = (e: MouseEvent) => {
      if (heroRef.current) {
        const rect = heroRef.current.getBoundingClientRect();
        const x = (e.clientX - rect.left) / rect.width;
        const y = (e.clientY - rect.top) / rect.height;
        setMousePosition({ x, y });
      }
    };

    window.addEventListener("mousemove", handleMouseMove);
    return () => window.removeEventListener("mousemove", handleMouseMove);
  }, []);

  return (
    <section
      ref={heroRef}
      className="relative mt-20 min-h-[95vh] flex items-center overflow-hidden"
    >
      {/* Advanced Animated Background */}
      <div className="absolute inset-0 bg-gradient-mesh" />

      {/* Animated Grid Pattern */}
      <div className="absolute inset-0 bg-[linear-gradient(to_right,#8080800a_1px,transparent_1px),linear-gradient(to_bottom,#8080800a_1px,transparent_1px)] bg-[size:32px_32px] [mask-image:radial-gradient(ellipse_80%_50%_at_50%_50%,#000_70%,transparent_100%)]" />

      {/* Parallax Orbs with mouse tracking */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div
          className="absolute top-1/4 left-1/4 w-[600px] h-[600px] bg-accent/10 rounded-full blur-3xl animate-float transition-transform duration-1000 ease-out"
          style={{
            transform: `translate(${mousePosition.x * 30}px, ${mousePosition.y * 30}px)`,
          }}
        />
        <div
          className="absolute bottom-1/4 right-1/4 w-[500px] h-[500px] bg-secondary/10 rounded-full blur-3xl animate-float transition-transform duration-1000 ease-out"
          style={{
            transform: `translate(${mousePosition.x * -20}px, ${mousePosition.y * -20}px)`,
            animationDelay: "1s",
          }}
        />
        <div
          className="absolute top-1/2 right-1/3 w-[400px] h-[400px] bg-blue-500/10 rounded-full blur-3xl animate-float transition-transform duration-1000 ease-out"
          style={{
            transform: `translate(${mousePosition.x * 15}px, ${mousePosition.y * 15}px)`,
            animationDelay: "2s",
          }}
        />
      </div>

      {/* Spotlight effect */}
      <div
        className="absolute inset-0 opacity-30 pointer-events-none"
        style={{
          background: `radial-gradient(circle 600px at ${mousePosition.x * 100}% ${mousePosition.y * 100}%, rgba(16, 185, 129, 0.15), transparent 80%)`,
        }}
      />

      <div className="max-w-7xl mx-auto px-6 sm:px-8 lg:px-12 py-20 relative z-10">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 lg:gap-20 items-center">
          {/* Left Content */}
          <div className="space-y-8 animate-fade-in-up">
            {/* Badge */}
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-surface/80 backdrop-blur-sm border border-border shadow-sm">
              <Zap className="w-4 h-4 text-accent" />
              <span className="text-sm font-medium text-text-secondary">
                Finanças simples para todos
              </span>
            </div>

            {/* Heading */}
            <h1 className="text-5xl sm:text-6xl lg:text-display-xl font-display font-bold text-primary leading-tight">
              Transforme sua{" "}
              <span className="relative inline-block">
                <span className="text-gradient-multi">relação</span>
                <svg className="absolute -bottom-2 left-0 w-full" height="8" viewBox="0 0 200 8" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M1 5.5C50 2.5 100 2.5 199 5.5" stroke="url(#gradient)" strokeWidth="3" strokeLinecap="round"/>
                  <defs>
                    <linearGradient id="gradient" x1="0%" y1="0%" x2="100%" y2="0%">
                      <stop offset="0%" stopColor="rgb(139, 92, 246)" />
                      <stop offset="50%" stopColor="rgb(16, 185, 129)" />
                      <stop offset="100%" stopColor="rgb(59, 130, 246)" />
                    </linearGradient>
                  </defs>
                </svg>
              </span>
              {" "}com o dinheiro
            </h1>

            <p className="text-lg sm:text-xl text-text-secondary leading-relaxed max-w-xl">
              Gerencie suas finanças de forma intuitiva e inteligente.
              FyApp torna o controle financeiro acessível, seguro e descomplicado.
            </p>

            {/* Stats */}
            <div className="flex flex-wrap gap-8 py-4">
              <div className="flex items-center gap-3">
                <div className="w-12 h-12 rounded-xl bg-accent/10 flex items-center justify-center">
                  <TrendingUp className="w-6 h-6 text-accent" />
                </div>
                <div>
                  <div className="text-2xl font-bold text-primary">100%</div>
                  <div className="text-sm text-text-tertiary">Gratuito</div>
                </div>
              </div>
              <div className="flex items-center gap-3">
                <div className="w-12 h-12 rounded-xl bg-secondary/10 flex items-center justify-center">
                  <Shield className="w-6 h-6 text-secondary" />
                </div>
                <div>
                  <div className="text-2xl font-bold text-primary">100%</div>
                  <div className="text-sm text-text-tertiary">Seguro</div>
                </div>
              </div>
            </div>

            {/* CTA Buttons */}
            <div className="flex flex-col sm:flex-row gap-4 pt-4">
              <a
                href={`${APP_URL}/login`}
                className="group inline-flex items-center justify-center gap-2 px-8 py-4 bg-accent text-white rounded-xl font-semibold shadow-lg shadow-accent/30 transition-all duration-300 hover:bg-accent-dark hover:-translate-y-1 hover:shadow-xl hover:shadow-accent/40"
              >
                Começar Agora
                <ArrowRight className="w-5 h-5 transition-transform group-hover:translate-x-1" />
              </a>
              <a
                href="#how"
                onClick={(e) => {
                  e.preventDefault();
                  document.getElementById("how")?.scrollIntoView({ behavior: "smooth" });
                }}
                className="group inline-flex items-center justify-center gap-2 px-8 py-4 bg-surface border-2 border-border text-primary rounded-xl font-semibold transition-all duration-300 hover:border-accent hover:bg-accent/5 hover:-translate-y-1 hover:shadow-lg cursor-pointer"
              >
                Ver Como Funciona
              </a>
            </div>
          </div>

          {/* Right Content - Dashboard Preview */}
          <div className="relative animate-fade-in-up lg:animate-slide-in-right" style={{ animationDelay: "200ms" }}>
            <div className="relative">
              {/* Glow Effect */}
              <div className="absolute -inset-4 bg-gradient-to-r from-accent/20 via-secondary/20 to-blue-500/20 rounded-3xl blur-2xl opacity-50" />

              {/* Image Container */}
              <div className="relative rounded-2xl overflow-hidden border border-border/50 shadow-2xl bg-surface/50 backdrop-blur-sm">
                <img
                  src="https://d2xsxph8kpxj0f.cloudfront.net/310419663031400654/9rw8cUK9yoafNAfP3CRw5n/fyapp-hero-illustration-DEpLkYQmzhkmUiCqwRERat.webp"
                  alt="FyApp Dashboard - Gerencie suas finanças"
                  className="w-full h-auto"
                />
              </div>

              {/* Floating Cards */}
              <div className="absolute -right-4 top-1/4 glass rounded-xl p-4 shadow-xl animate-float hidden lg:block" style={{ animationDelay: "0.5s" }}>
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-lg bg-accent/20 flex items-center justify-center">
                    <TrendingUp className="w-5 h-5 text-accent" />
                  </div>
                  <div>
                    <div className="text-xs text-text-tertiary">Economia</div>
                    <div className="text-sm font-bold text-primary">+R$ 1.245</div>
                  </div>
                </div>
              </div>

              <div className="absolute -left-4 bottom-1/4 glass rounded-xl p-4 shadow-xl animate-float hidden lg:block" style={{ animationDelay: "1s" }}>
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-lg bg-secondary/20 flex items-center justify-center">
                    <Shield className="w-5 h-5 text-secondary" />
                  </div>
                  <div>
                    <div className="text-xs text-text-tertiary">Segurança</div>
                    <div className="text-sm font-bold text-primary">Ativa</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
