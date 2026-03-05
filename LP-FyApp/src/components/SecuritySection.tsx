import { Lock } from "lucide-react";

const APP_URL = process.env.NEXT_PUBLIC_APP_URL || "https://fyapp.up.railway.app";

const securityFeatures = [
  "Criptografia de nível bancário (AES-256)",
  "Autenticação de dois fatores (2FA)",
  "Conformidade com LGPD, GDPR e PCI-DSS",
  "Backup automático e redundância de dados",
  "Monitoramento 24/7 de atividades suspeitas",
];

export default function SecuritySection() {
  return (
    <section
      id="security"
      className="py-32 px-8 bg-gradient-to-br from-[#FAFAF8] to-[#F5F1E8]"
    >
      <div className="max-w-7xl mx-auto grid grid-cols-1 lg:grid-cols-2 gap-16 items-center">
        <div className="space-y-8">
          <h2 className="text-4xl font-bold text-[#2C2C2C] font-playfair">
            Segurança em Primeiro Lugar
          </h2>

          <p className="text-lg text-[#8B8B8B] leading-relaxed">
            Sua segurança financeira é nossa prioridade. FyApp utiliza as
            tecnologias mais avançadas de proteção de dados.
          </p>

          <ul className="space-y-3">
            {securityFeatures.map((feature, index) => (
              <li key={index} className="flex items-start gap-3 text-lg text-[#8B8B8B]">
                <Lock className="w-6 h-6 text-[#7E8C54] shrink-0 mt-1" />
                <span>{feature}</span>
              </li>
            ))}
          </ul>

          <a
            href={`${APP_URL}/login`}
            className="inline-block bg-[#7E8C54] text-white px-8 py-3.5 rounded-lg font-semibold transition-all duration-300 hover:bg-[#6B7844] hover:-translate-y-1 hover:shadow-lg border-2 border-[#7E8C54]"
          >
            Saiba Mais Sobre Segurança
          </a>
        </div>

        <div className="relative">
          <img
            src="https://d2xsxph8kpxj0f.cloudfront.net/310419663031400654/9rw8cUK9yoafNAfP3CRw5n/fyapp-security-illustration-acUv2UXhPFgFD7SXqPsEPX.webp"
            alt="Segurança FyApp"
            className="w-full rounded-2xl shadow-2xl hover:scale-105 transition-transform duration-300"
          />
        </div>
      </div>
    </section>
  );
}
