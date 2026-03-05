import React from 'react';

export default function Hero() {
  return (
    <section className="mt-20 py-32 px-8 bg-gradient-to-br from-[#FAFAF8] to-[#F5F1E8] relative overflow-hidden">
      {/* Background floating element */}
      <div className="absolute top-0 right-0 w-96 h-96 bg-gradient-radial from-[rgba(126,140,84,0.1)] to-transparent rounded-full -mr-48 -mt-48 animate-pulse"></div>

      <div className="max-w-7xl mx-auto grid grid-cols-1 lg:grid-cols-2 gap-16 items-center relative z-10">
        {/* Text Content */}
        <div className="space-y-8 animate-in fade-in slide-in-from-left-8 duration-800">
          <h1 className="text-5xl lg:text-7xl font-bold text-[#2C2C2C] leading-tight">
            Seu Dinheiro, <span className="text-[#7E8C54] relative">Melhor<span className="absolute -bottom-3 left-0 right-0 h-1 bg-[#D4A574] rounded-full"></span></span> Controlado
          </h1>
          
          <p className="text-lg text-[#8B8B8B] leading-relaxed max-w-xl">
            FyApp é a plataforma de gestão financeira moderna que torna simples o controle de suas transações, orçamentos e metas. Segurança e simplicidade em um só lugar.
          </p>

          <div className="flex flex-col sm:flex-row gap-4 pt-4">
            <button className="bg-[#7E8C54] text-white px-8 py-3.5 rounded-lg font-semibold transition-all duration-300 hover:bg-[#6B7844] hover:-translate-y-1 hover:shadow-lg border-2 border-[#7E8C54]">
              Começar Agora
            </button>
            <button className="bg-transparent text-[#7E8C54] px-8 py-3.5 rounded-lg font-semibold transition-all duration-300 hover:bg-[#7E8C54] hover:text-white hover:-translate-y-1 hover:shadow-lg border-2 border-[#7E8C54]">
              Ver Demo
            </button>
          </div>
        </div>

        {/* Image */}
        <div className="relative animate-in fade-in slide-in-from-right-8 duration-800">
          <img
            src="https://d2xsxph8kpxj0f.cloudfront.net/310419663031400654/9rw8cUK9yoafNAfP3CRw5n/fyapp-hero-illustration-DEpLkYQmzhkmUiCqwRERat.webp"
            alt="FyApp Dashboard"
            className="w-full rounded-2xl shadow-2xl hover:scale-105 transition-transform duration-300"
          />
        </div>
      </div>
    </section>
  );
}
