"use client";

import { ChevronDown, HelpCircle } from "lucide-react";
import { useState } from "react";

const faqs = [
  {
    question: "O FyApp é realmente gratuito?",
    answer: "Sim! O FyApp é 100% gratuito para sempre. Você tem acesso a todos os recursos sem custos ocultos, sem limite de transações e sem necessidade de cartão de crédito para começar.",
  },
  {
    question: "Meus dados estão seguros?",
    answer: "Absolutamente. Usamos criptografia de nível bancário (AES-256), autenticação de dois fatores e somos totalmente conformes com LGPD, GDPR e PCI-DSS. Seus dados nunca são compartilhados com terceiros.",
  },
  {
    question: "Como funciona a sincronização automática?",
    answer: "Você conecta suas contas bancárias e cartões de forma segura. O FyApp sincroniza automaticamente suas transações em tempo real, categorizando-as inteligentemente para você não precisar fazer nada manualmente.",
  },
  {
    question: "Posso usar em vários dispositivos?",
    answer: "Sim! O FyApp funciona em qualquer dispositivo - computador, tablet e smartphone. Seus dados são sincronizados em tempo real na nuvem, então você sempre tem acesso às suas finanças onde estiver.",
  },
  {
    question: "E se eu não gostar?",
    answer: "Não há compromisso! Como é gratuito, você pode experimentar sem riscos. Se não atender suas necessidades, basta parar de usar. Mas acreditamos que você vai adorar! 😊",
  },
  {
    question: "Quanto tempo leva para configurar?",
    answer: "Menos de 2 minutos! Basta criar sua conta, conectar suas contas bancárias e pronto. O FyApp cuida do resto automaticamente.",
  },
];

export default function FAQ() {
  const [openIndex, setOpenIndex] = useState<number | null>(0);

  return (
    <section className="relative py-24 sm:py-32 overflow-hidden">
      {/* Background */}
      <div className="absolute inset-0 bg-gradient-to-b from-background to-background-alt" />

      <div className="max-w-4xl mx-auto px-6 sm:px-8 lg:px-12 relative z-10">
        {/* Header */}
        <div className="text-center mb-16 space-y-4 animate-fade-in-up">
          <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-accent/10 border border-accent/20">
            <HelpCircle className="w-4 h-4 text-accent" />
            <span className="text-sm font-semibold text-accent">FAQ</span>
          </div>
          <h2 className="text-4xl sm:text-5xl lg:text-display-lg font-display font-bold text-primary">
            Perguntas{" "}
            <span className="text-gradient-accent">Frequentes</span>
          </h2>
          <p className="text-lg text-text-secondary">
            Tire suas dúvidas sobre o FyApp
          </p>
        </div>

        {/* FAQ Items */}
        <div className="space-y-4">
          {faqs.map((faq, index) => (
            <div
              key={index}
              className="group animate-fade-in-up"
              style={{ animationDelay: `${index * 50}ms` }}
            >
              <button
                onClick={() => setOpenIndex(openIndex === index ? null : index)}
                className="w-full text-left bg-surface rounded-xl border border-border p-6 transition-all duration-300 hover:border-accent/50 hover:shadow-lg"
              >
                <div className="flex items-start justify-between gap-4">
                  <h3 className="text-lg font-display font-bold text-primary pr-8">
                    {faq.question}
                  </h3>
                  <div className={`flex-shrink-0 w-8 h-8 rounded-lg bg-accent/10 flex items-center justify-center transition-all duration-300 ${openIndex === index ? 'rotate-180 bg-accent' : ''}`}>
                    <ChevronDown className={`w-5 h-5 transition-colors ${openIndex === index ? 'text-white' : 'text-accent'}`} />
                  </div>
                </div>

                {/* Answer */}
                <div
                  className={`grid transition-all duration-300 ease-in-out ${
                    openIndex === index ? 'grid-rows-[1fr] opacity-100 mt-4' : 'grid-rows-[0fr] opacity-0'
                  }`}
                >
                  <div className="overflow-hidden">
                    <p className="text-text-secondary leading-relaxed">
                      {faq.answer}
                    </p>
                  </div>
                </div>
              </button>
            </div>
          ))}
        </div>

        {/* Bottom CTA */}
        <div className="mt-16 text-center p-8 rounded-2xl bg-gradient-to-br from-accent/5 to-secondary/5 border border-accent/20 animate-fade-in-up" style={{ animationDelay: "300ms" }}>
          <p className="text-lg text-text-secondary mb-4">
            Ainda tem dúvidas?
          </p>
          <a
            href="#"
            className="inline-flex items-center gap-2 text-accent hover:text-accent-dark font-semibold transition-colors group"
          >
            Entre em contato conosco
            <ChevronDown className="w-4 h-4 rotate-[-90deg] transition-transform group-hover:translate-x-1" />
          </a>
        </div>
      </div>
    </section>
  );
}
