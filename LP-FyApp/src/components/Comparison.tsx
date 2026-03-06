import { X, Check, ArrowRight } from "lucide-react";

const APP_URL = process.env.NEXT_PUBLIC_APP_URL || "https://fyapp-production.up.railway.app";

const comparisons = [
  {
    without: "Planilhas confusas e desorganizadas",
    with: "Dashboard intuitivo e visual",
  },
  {
    without: "Horas perdidas categorizando gastos",
    with: "Categorização automática inteligente",
  },
  {
    without: "Surpresas no fim do mês",
    with: "Alertas em tempo real",
  },
  {
    without: "Dificuldade em atingir metas",
    with: "Acompanhamento motivador de progresso",
  },
  {
    without: "Dados espalhados em vários apps",
    with: "Tudo centralizado em um só lugar",
  },
];

export default function Comparison() {
  return (
    <section className="relative py-24 sm:py-32 overflow-hidden">
      {/* Background */}
      <div className="absolute inset-0 bg-surface" />
      <div className="absolute inset-0 bg-[radial-gradient(circle_at_30%_50%,rgba(139,92,246,0.05),transparent_50%),radial-gradient(circle_at_70%_50%,rgba(16,185,129,0.05),transparent_50%)]" />

      <div className="max-w-6xl mx-auto px-6 sm:px-8 lg:px-12 relative z-10">
        {/* Header */}
        <div className="text-center mb-16 space-y-4 animate-fade-in-up">
          <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-secondary/10 border border-secondary/20">
            <ArrowRight className="w-4 h-4 text-secondary" />
            <span className="text-sm font-semibold text-secondary">Transformação</span>
          </div>
          <h2 className="text-4xl sm:text-5xl lg:text-display-lg font-display font-bold text-primary">
            Antes vs Depois do{" "}
            <span className="text-gradient-accent">FyApp</span>
          </h2>
          <p className="text-lg text-text-secondary max-w-2xl mx-auto">
            Veja a diferença de gerenciar suas finanças com e sem o FyApp
          </p>
        </div>

        {/* Comparison Cards */}
        <div className="relative">
          {/* Divider with VS */}
          <div className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 z-20 hidden lg:block">
            <div className="w-16 h-16 rounded-full bg-gradient-to-br from-accent to-secondary shadow-2xl flex items-center justify-center border-4 border-surface">
              <span className="text-white font-bold text-sm">VS</span>
            </div>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 lg:gap-12">
            {/* Sem FyApp */}
            <div className="relative animate-fade-in-up">
              <div className="sticky top-24">
                <div className="bg-gradient-to-br from-red-50 to-orange-50 dark:from-red-950/20 dark:to-orange-950/20 rounded-2xl border-2 border-red-200/50 p-8 shadow-lg">
                  <div className="flex items-center gap-3 mb-6">
                    <div className="w-12 h-12 rounded-xl bg-red-500/10 flex items-center justify-center">
                      <X className="w-6 h-6 text-red-500" />
                    </div>
                    <h3 className="text-2xl font-display font-bold text-red-900 dark:text-red-100">
                      Sem FyApp
                    </h3>
                  </div>
                  <ul className="space-y-4">
                    {comparisons.map((item, index) => (
                      <li key={index} className="flex items-start gap-3">
                        <X className="w-5 h-5 text-red-500 shrink-0 mt-0.5" />
                        <span className="text-red-900/80 dark:text-red-100/80">
                          {item.without}
                        </span>
                      </li>
                    ))}
                  </ul>
                </div>
              </div>
            </div>

            {/* Com FyApp */}
            <div className="relative animate-fade-in-up" style={{ animationDelay: "200ms" }}>
              <div className="sticky top-24">
                <div className="relative bg-gradient-to-br from-emerald-50 to-teal-50 dark:from-emerald-950/20 dark:to-teal-950/20 rounded-2xl border-2 border-emerald-500/50 p-8 shadow-2xl">
                  {/* Shine effect */}
                  <div className="absolute inset-0 rounded-2xl bg-gradient-to-r from-transparent via-white/20 to-transparent opacity-0 hover:opacity-100 transition-opacity duration-1000 pointer-events-none shimmer" />

                  <div className="flex items-center gap-3 mb-6 relative z-10">
                    <div className="w-12 h-12 rounded-xl bg-accent/20 flex items-center justify-center">
                      <Check className="w-6 h-6 text-accent" />
                    </div>
                    <h3 className="text-2xl font-display font-bold text-emerald-900 dark:text-emerald-100">
                      Com FyApp
                    </h3>
                  </div>
                  <ul className="space-y-4 relative z-10">
                    {comparisons.map((item, index) => (
                      <li key={index} className="flex items-start gap-3 group">
                        <div className="w-5 h-5 rounded-full bg-accent/20 flex items-center justify-center shrink-0 mt-0.5 group-hover:bg-accent group-hover:scale-110 transition-all duration-300">
                          <Check className="w-3 h-3 text-accent group-hover:text-white transition-colors" />
                        </div>
                        <span className="text-emerald-900 dark:text-emerald-100 font-medium">
                          {item.with}
                        </span>
                      </li>
                    ))}
                  </ul>

                  {/* CTA Button */}
                  <div className="mt-8 pt-6 border-t border-emerald-200/50 relative z-10">
                    <a
                      href={`${APP_URL}/login`}
                      className="group w-full inline-flex items-center justify-center gap-2 px-6 py-3.5 bg-accent text-white rounded-xl font-semibold shadow-lg shadow-accent/30 transition-all duration-300 hover:bg-accent-dark hover:-translate-y-1 hover:shadow-xl hover:shadow-accent/40"
                    >
                      Começar Agora
                      <ArrowRight className="w-5 h-5 transition-transform group-hover:translate-x-1" />
                    </a>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Bottom CTA */}
        <div className="mt-16 text-center animate-fade-in-up" style={{ animationDelay: "400ms" }}>
          <p className="text-text-secondary text-lg mb-2">
            Junte-se a <span className="font-bold text-primary">milhares de usuários</span> que já fizeram a escolha certa
          </p>
          <p className="text-sm text-text-tertiary">
            ⚡ Comece grátis em menos de 2 minutos
          </p>
        </div>
      </div>
    </section>
  );
}
