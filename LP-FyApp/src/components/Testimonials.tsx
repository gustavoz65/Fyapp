"use client";

import { Star, Quote } from "lucide-react";
import { useState } from "react";

const testimonials = [
  {
    name: "Maria Silva",
    role: "Empreendedora",
    initials: "MS",
    gradient: "from-emerald-500 to-teal-500",
    content: "FyApp mudou completamente como gerencio as finanças do meu negócio. Economizei mais de R$ 5.000 em 3 meses só identificando gastos desnecessários.",
    rating: 5,
    savings: "R$ 5.000",
  },
  {
    name: "João Santos",
    role: "Desenvolvedor",
    initials: "JS",
    gradient: "from-violet-500 to-purple-500",
    content: "Finalmente consegui organizar minha vida financeira. A interface é tão simples que uso todos os dias sem esforço. Recomendo muito!",
    rating: 5,
    savings: "40% economia",
  },
  {
    name: "Ana Costa",
    role: "Designer",
    initials: "AC",
    gradient: "from-blue-500 to-cyan-500",
    content: "O melhor app de finanças que já usei. As análises inteligentes me ajudaram a alcançar minha meta de comprar um carro em 8 meses.",
    rating: 5,
    savings: "Meta alcançada",
  },
];

export default function Testimonials() {
  const [activeIndex, setActiveIndex] = useState(0);

  return (
    <section className="relative py-24 sm:py-32 overflow-hidden">
      {/* Background */}
      <div className="absolute inset-0 bg-gradient-to-br from-primary/5 via-background to-accent/5" />

      {/* Decorative blur orbs */}
      <div className="absolute top-1/4 -left-32 w-96 h-96 bg-accent/10 rounded-full blur-3xl animate-pulse-slow" />
      <div className="absolute bottom-1/4 -right-32 w-96 h-96 bg-secondary/10 rounded-full blur-3xl animate-pulse-slow" style={{ animationDelay: "1s" }} />

      <div className="max-w-7xl mx-auto px-6 sm:px-8 lg:px-12 relative z-10">
        {/* Header */}
        <div className="text-center mb-16 space-y-4 animate-fade-in-up">
          <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-accent/10 border border-accent/20">
            <Quote className="w-4 h-4 text-accent" />
            <span className="text-sm font-semibold text-accent">Depoimentos</span>
          </div>
          <h2 className="text-4xl sm:text-5xl lg:text-display-lg font-display font-bold text-primary">
            O que nossos usuários{" "}
            <span className="text-gradient-accent">estão dizendo</span>
          </h2>
          <p className="text-lg text-text-secondary max-w-2xl mx-auto">
            Milhares de pessoas já transformaram suas finanças com FyApp
          </p>
        </div>

        {/* Testimonials Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mb-12">
          {testimonials.map((testimonial, index) => (
            <div
              key={index}
              className="group relative animate-fade-in-up"
              style={{ animationDelay: `${index * 150}ms` }}
            >
              <div className="relative h-full bg-surface rounded-2xl border border-border p-8 transition-all duration-500 hover:shadow-2xl hover:-translate-y-2 hover:border-accent/50">
                {/* Quote Icon */}
                <div className="absolute -top-4 -right-4 w-12 h-12 bg-accent/10 rounded-full flex items-center justify-center backdrop-blur-sm border border-accent/20">
                  <Quote className="w-6 h-6 text-accent" />
                </div>

                {/* Rating */}
                <div className="flex gap-1 mb-4">
                  {[...Array(testimonial.rating)].map((_, i) => (
                    <Star key={i} className="w-5 h-5 fill-amber-400 text-amber-400" />
                  ))}
                </div>

                {/* Content */}
                <p className="text-text-secondary leading-relaxed mb-6 italic">
                  "{testimonial.content}"
                </p>

                {/* Savings Badge */}
                <div className="mb-6 inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-accent/10 border border-accent/20">
                  <div className="w-2 h-2 rounded-full bg-accent animate-pulse" />
                  <span className="text-sm font-semibold text-accent">{testimonial.savings}</span>
                </div>

                {/* Author */}
                <div className="flex items-center gap-4">
                  <div className="relative">
                    <div className={`w-12 h-12 rounded-full bg-gradient-to-br ${testimonial.gradient} flex items-center justify-center ring-2 ring-border group-hover:ring-accent transition-all duration-300 shadow-lg`}>
                      <span className="text-white font-bold text-sm">{testimonial.initials}</span>
                    </div>
                    <div className="absolute -bottom-1 -right-1 w-5 h-5 bg-accent rounded-full flex items-center justify-center border-2 border-surface">
                      <svg className="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 20 20">
                        <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                      </svg>
                    </div>
                  </div>
                  <div>
                    <div className="font-display font-bold text-primary group-hover:text-accent transition-colors">
                      {testimonial.name}
                    </div>
                    <div className="text-sm text-text-tertiary">{testimonial.role}</div>
                  </div>
                </div>

                {/* Hover gradient effect */}
                <div className="absolute inset-0 rounded-2xl bg-gradient-to-r from-accent/0 via-accent/5 to-accent/0 opacity-0 group-hover:opacity-100 transition-opacity duration-500 pointer-events-none" />
              </div>
            </div>
          ))}
        </div>

        {/* Trust Indicators */}
        <div className="flex flex-wrap items-center justify-center gap-8 pt-8 border-t border-border/50">
          <div className="text-center">
            <div className="text-3xl font-bold text-primary mb-1">4.9/5</div>
            <div className="text-sm text-text-tertiary">Avaliação média</div>
          </div>
          <div className="w-px h-12 bg-border" />
          <div className="text-center">
            <div className="text-3xl font-bold text-primary mb-1">1K+</div>
            <div className="text-sm text-text-tertiary">Usuários ativos</div>
          </div>
          <div className="w-px h-12 bg-border" />
          <div className="text-center">
            <div className="text-3xl font-bold text-primary mb-1">R$ 2M+</div>
            <div className="text-sm text-text-tertiary">Economizados</div>
          </div>
        </div>
      </div>
    </section>
  );
}
