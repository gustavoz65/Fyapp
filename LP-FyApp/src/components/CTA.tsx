"use client";

import { Sparkles, ArrowRight, TrendingUp } from "lucide-react";

const APP_URL = process.env.NEXT_PUBLIC_APP_URL || "https://fyapp-production.up.railway.app";

export default function CTA() {
  return (
    <section className="relative py-24 sm:py-32 overflow-hidden">
      {/* Animated Gradient Background */}
      <div className="absolute inset-0 bg-gradient-to-br from-accent via-secondary to-blue-600 animated-gradient" />

      {/* Decorative Elements */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-white/10 rounded-full blur-3xl animate-float" />
        <div className="absolute bottom-1/4 right-1/4 w-96 h-96 bg-white/10 rounded-full blur-3xl animate-float" style={{ animationDelay: "1s" }} />
      </div>

      {/* Pattern Overlay */}
      <div className="absolute inset-0 bg-[linear-gradient(to_right,#ffffff0a_1px,transparent_1px),linear-gradient(to_bottom,#ffffff0a_1px,transparent_1px)] bg-[size:24px_24px]" />

      <div className="max-w-4xl mx-auto px-6 sm:px-8 lg:px-12 text-center relative z-10">
        <div className="space-y-8 animate-fade-in-up">
          {/* Badge */}
          <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-white/20 backdrop-blur-sm border border-white/30">
            <Sparkles className="w-4 h-4 text-white" />
            <span className="text-sm font-semibold text-white">Comece Hoje</span>
          </div>

          {/* Heading */}
          <h2 className="text-4xl sm:text-5xl lg:text-6xl font-display font-bold text-white leading-tight">
            Pronto para transformar
            <br />
            suas finanças?
          </h2>

          {/* Description */}
          <p className="text-xl text-white/90 leading-relaxed max-w-2xl mx-auto">
            Junte-se a <span className="font-bold text-white">centenas de usuários</span> que já estão
            controlando suas finanças com FyApp. Comece grátis hoje mesmo.
          </p>

          {/* Stats */}
          <div className="flex flex-wrap justify-center gap-8 py-6">
            <div className="text-center">
              <div className="text-4xl font-bold text-white">100%</div>
              <div className="text-sm text-white/80">Grátis</div>
            </div>
            <div className="w-px h-16 bg-white/20" />
            <div className="text-center">
              <div className="text-4xl font-bold text-white">2min</div>
              <div className="text-sm text-white/80">Para começar</div>
            </div>
            <div className="w-px h-16 bg-white/20" />
            <div className="text-center">
              <div className="text-4xl font-bold text-white">24/7</div>
              <div className="text-sm text-white/80">Suporte</div>
            </div>
          </div>

          {/* CTA Buttons */}
          <div className="flex flex-col sm:flex-row gap-4 justify-center pt-4">
            <a
              href={`${APP_URL}/login`}
              className="group inline-flex items-center justify-center gap-2 px-10 py-4 bg-white text-accent rounded-xl font-bold text-lg shadow-2xl transition-all duration-300 hover:bg-white/90 hover:-translate-y-1 hover:shadow-[0_20px_60px_rgba(255,255,255,0.4)]"
            >
              Começar Grátis Agora
              <ArrowRight className="w-5 h-5 transition-transform group-hover:translate-x-1" />
            </a>
            <a
              href="#features"
              onClick={(e) => {
                e.preventDefault();
                document.getElementById("features")?.scrollIntoView({ behavior: "smooth" });
              }}
              className="group inline-flex items-center justify-center gap-2 px-10 py-4 bg-white/10 backdrop-blur-sm text-white rounded-xl font-bold text-lg border-2 border-white/30 transition-all duration-300 hover:bg-white/20 hover:-translate-y-1 hover:shadow-xl cursor-pointer"
            >
              Ver Recursos
              <TrendingUp className="w-5 h-5 transition-transform group-hover:-translate-y-0.5" />
            </a>
          </div>

          {/* Trust Badge */}
          <p className="text-sm text-white/70 pt-4">
            🔒 Seus dados estão seguros com criptografia de nível bancário
          </p>
        </div>
      </div>
    </section>
  );
}
