import { Shield, Lock, Eye, Server, Bell, CheckCircle2, ArrowRight } from "lucide-react";

const APP_URL = process.env.NEXT_PUBLIC_APP_URL || "https://fyapp-production.up.railway.app";

const securityFeatures = [
  {
    icon: Lock,
    text: "Criptografia de nível bancário (AES-256)",
  },
  {
    icon: Shield,
    text: "Autenticação de dois fatores (2FA)",
  },
  {
    icon: CheckCircle2,
    text: "Conformidade com LGPD, GDPR e PCI-DSS",
  },
  {
    icon: Server,
    text: "Backup automático e redundância de dados",
  },
  {
    icon: Eye,
    text: "Monitoramento 24/7 de atividades suspeitas",
  },
];

export default function SecuritySection() {
  return (
    <section
      id="security"
      className="relative py-24 sm:py-32 overflow-hidden"
    >
      {/* Background */}
      <div className="absolute inset-0 bg-gradient-to-br from-background via-accent/5 to-background" />

      <div className="max-w-7xl mx-auto px-6 sm:px-8 lg:px-12 relative z-10">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 lg:gap-20 items-center">
          {/* Left Content */}
          <div className="space-y-8 animate-fade-in-up">
            {/* Badge */}
            <div className="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-red-500/10 border border-red-500/20">
              <Shield className="w-4 h-4 text-red-500" />
              <span className="text-sm font-semibold text-red-500">Segurança Máxima</span>
            </div>

            {/* Heading */}
            <h2 className="text-4xl sm:text-5xl lg:text-display-lg font-display font-bold text-primary leading-tight">
              Segurança em{" "}
              <span className="relative inline-block">
                <span className="text-gradient-accent">primeiro lugar</span>
                <svg className="absolute -bottom-2 left-0 w-full" height="8" viewBox="0 0 200 8" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M1 5.5C50 2.5 100 2.5 199 5.5" stroke="rgb(16, 185, 129)" strokeWidth="3" strokeLinecap="round"/>
                </svg>
              </span>
            </h2>

            {/* Description */}
            <p className="text-lg text-text-secondary leading-relaxed">
              Sua segurança financeira é nossa <span className="font-semibold text-primary">prioridade máxima</span>.
              FyApp utiliza as tecnologias mais avançadas de proteção de dados
              para garantir que suas informações estejam sempre protegidas.
            </p>

            {/* Security Features List */}
            <ul className="space-y-4">
              {securityFeatures.map((feature, index) => {
                const Icon = feature.icon;
                return (
                  <li
                    key={index}
                    className="flex items-start gap-4 group"
                    style={{ animationDelay: `${index * 100}ms` }}
                  >
                    <div className="flex-shrink-0 w-10 h-10 rounded-lg bg-accent/10 flex items-center justify-center group-hover:bg-accent group-hover:scale-110 transition-all duration-300">
                      <Icon className="w-5 h-5 text-accent group-hover:text-white transition-colors duration-300" />
                    </div>
                    <span className="text-base text-text-secondary group-hover:text-primary transition-colors duration-300 pt-2">
                      {feature.text}
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
              Começar com Segurança
              <ArrowRight className="w-5 h-5 transition-transform group-hover:translate-x-1" />
            </a>
          </div>

          {/* Right Content - Image */}
          <div className="relative animate-fade-in-up lg:animate-slide-in-right" style={{ animationDelay: "200ms" }}>
            <div className="relative">
              {/* Glow Effect */}
              <div className="absolute -inset-4 bg-gradient-to-r from-red-500/20 via-accent/20 to-red-500/20 rounded-3xl blur-3xl opacity-50 animate-pulse-slow" />

              {/* Image Container */}
              <div className="relative rounded-2xl overflow-hidden border border-border/50 shadow-2xl bg-surface/30 backdrop-blur-sm">
                <img
                  src="https://d2xsxph8kpxj0f.cloudfront.net/310419663031400654/9rw8cUK9yoafNAfP3CRw5n/fyapp-security-illustration-acUv2UXhPFgFD7SXqPsEPX.webp"
                  alt="Segurança FyApp - Proteção de Dados"
                  className="w-full h-auto"
                />

                {/* Overlay Gradient */}
                <div className="absolute inset-0 bg-gradient-to-t from-primary/10 via-transparent to-transparent pointer-events-none" />
              </div>

              {/* Security Badge */}
              <div className="absolute -bottom-6 left-1/2 -translate-x-1/2 glass rounded-xl px-6 py-3 shadow-xl animate-fade-in-up" style={{ animationDelay: "600ms" }}>
                <div className="flex items-center gap-3">
                  <div className="w-10 h-10 rounded-lg bg-accent/20 flex items-center justify-center">
                    <Shield className="w-5 h-5 text-accent" />
                  </div>
                  <div>
                    <div className="text-xs text-text-tertiary">Certificado</div>
                    <div className="text-sm font-bold text-primary">Nível Bancário</div>
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
