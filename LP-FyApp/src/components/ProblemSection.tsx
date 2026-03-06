import { CheckCircle2, ArrowRight, LineChart, PiggyBank, Sparkles } from "lucide-react";

const APP_URL = process.env.NEXT_PUBLIC_APP_URL || "https://fyapp-production.up.railway.app";

const benefits = [
  {
    text: "Visualize todos os seus gastos em um único dashboard",
    icon: LineChart,
  },
  {
    text: "Receba alertas inteligentes sobre transações importantes",
    icon: Sparkles,
  },
  {
    text: "Defina e acompanhe suas metas financeiras",
    icon: PiggyBank,
  },
  {
    text: "Gerencie múltiplas contas e cartões de crédito",
    icon: CheckCircle2,
  },
  {
    text: "Relatórios detalhados e análises em tempo real",
    icon: LineChart,
  },
];

export default function ProblemSection() {
  return (
    <section className="relative py-24 sm:py-32 overflow-hidden">
      {/* Background */}
      <div className="absolute inset-0 bg-surface" />
      <div className="absolute inset-0 bg-gradient-to-br from-accent/5 via-transparent to-secondary/5" />

      <div className="max-w-7xl mx-auto px-6 sm:px-8 lg:px-12 relative z-10">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 lg:gap-20 items-center">
          {/* Left Content */}
          <div className="space-y-8 animate-fade-in-up">
            {/* Badge */}
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-accent/10 border border-accent/20">
              <div className="w-2 h-2 rounded-full bg-accent animate-pulse" />
              <span className="text-sm font-semibold text-accent">Por que FyApp?</span>
            </div>

            {/* Heading */}
            <h2 className="text-4xl sm:text-5xl lg:text-display-lg font-display font-bold text-primary leading-tight">
              Controle financeiro{" "}
              <span className="text-gradient-accent">nunca foi tão fácil</span>
            </h2>

            {/* Description */}
            <p className="text-lg text-text-secondary leading-relaxed">
              A maioria das pessoas luta para entender seus gastos e planejar o
              futuro financeiro. <span className="font-semibold text-primary">FyApp muda isso</span> oferecendo uma visão clara e em
              tempo real de suas finanças, com ferramentas intuitivas que{" "}
              <span className="font-semibold text-primary">qualquer um pode usar</span>.
            </p>

            {/* Benefits List */}
            <ul className="space-y-4">
              {benefits.map((benefit, index) => {
                const Icon = benefit.icon;
                return (
                  <li
                    key={index}
                    className="flex items-start gap-4 group"
                    style={{ animationDelay: `${index * 100}ms` }}
                  >
                    <div className="flex-shrink-0 w-10 h-10 rounded-lg bg-accent/10 flex items-center justify-center group-hover:bg-accent group-hover:scale-110 transition-all duration-300">
                      <CheckCircle2 className="w-5 h-5 text-accent group-hover:text-white transition-colors duration-300" />
                    </div>
                    <span className="text-base text-text-secondary group-hover:text-primary transition-colors duration-300 pt-2">
                      {benefit.text}
                    </span>
                  </li>
                );
              })}
            </ul>

            {/* CTA */}
            <a
              href={`${APP_URL}/login`}
              className="group inline-flex items-center gap-2 px-8 py-4 bg-primary text-white rounded-xl font-semibold shadow-lg shadow-primary/20 transition-all duration-300 hover:bg-primary-light hover:-translate-y-1 hover:shadow-xl hover:shadow-primary/30"
            >
              Explorar Recursos
              <ArrowRight className="w-5 h-5 transition-transform group-hover:translate-x-1" />
            </a>
          </div>

          {/* Right Content - Image */}
          <div className="relative lg:order-last animate-fade-in-up lg:animate-slide-in-right" style={{ animationDelay: "200ms" }}>
            <div className="relative">
              {/* Background Glow */}
              <div className="absolute -inset-4 bg-gradient-to-r from-accent/20 via-secondary/20 to-accent/20 rounded-3xl blur-3xl opacity-50 animate-pulse-slow" />

              {/* Image Container with Glass Effect */}
              <div className="relative rounded-2xl overflow-hidden border border-border/50 shadow-2xl bg-surface/30 backdrop-blur-sm">
                <img
                  src="https://d2xsxph8kpxj0f.cloudfront.net/310419663031400654/9rw8cUK9yoafNAfP3CRw5n/fyapp-features-illustration-X72PdFpSHkAv6RDpzsykco.webp"
                  alt="Dashboard FyApp - Recursos e Funcionalidades"
                  className="w-full h-auto"
                />

                {/* Overlay Gradient */}
                <div className="absolute inset-0 bg-gradient-to-t from-primary/10 via-transparent to-transparent pointer-events-none" />
              </div>

              {/* Decorative Elements */}
              <div className="absolute -top-6 -right-6 w-32 h-32 bg-accent/20 rounded-full blur-3xl animate-pulse-slow" />
              <div className="absolute -bottom-6 -left-6 w-32 h-32 bg-secondary/20 rounded-full blur-3xl animate-pulse-slow" style={{ animationDelay: "1s" }} />
            </div>

            {/* Stats Card */}
            <div className="absolute -bottom-8 left-8 right-8 glass rounded-xl p-6 shadow-xl hidden lg:block animate-fade-in-up" style={{ animationDelay: "600ms" }}>
              <div className="flex items-center justify-around gap-6">
                <div className="text-center">
                  <div className="text-3xl font-bold text-primary">100%</div>
                  <div className="text-sm text-text-tertiary">Grátis</div>
                </div>
                <div className="w-px h-12 bg-border" />
                <div className="text-center">
                  <div className="text-3xl font-bold text-accent">Fácil</div>
                  <div className="text-sm text-text-tertiary">de Usar</div>
                </div>
                <div className="w-px h-12 bg-border" />
                <div className="text-center">
                  <div className="text-3xl font-bold text-secondary">Seguro</div>
                  <div className="text-sm text-text-tertiary">100%</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
