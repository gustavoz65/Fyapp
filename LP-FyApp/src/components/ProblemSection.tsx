import { Check } from "lucide-react";

const APP_URL = process.env.NEXT_PUBLIC_APP_URL || "https://fyapp-production.up.railway.app";

const benefits = [
  "Visualize todos os seus gastos em um único dashboard",
  "Receba alertas inteligentes sobre transações importantes",
  "Defina e acompanhe suas metas financeiras",
  "Gerencie múltiplas contas e cartões de crédito",
  "Relatórios detalhados e análises em tempo real",
];

export default function ProblemSection() {
  return (
    <section className="py-32 px-8 bg-white border-t border-b border-[#E5DDD0]">
      <div className="max-w-7xl mx-auto grid grid-cols-1 lg:grid-cols-2 gap-16 items-center">
        <div className="space-y-8">
          <h2 className="text-4xl font-bold text-[#2C2C2C] font-playfair">
            Controle Financeiro Nunca Foi Tão Fácil
          </h2>

          <p className="text-lg text-[#8B8B8B] leading-relaxed">
            A maioria das pessoas luta para entender seus gastos e planejar o
            futuro financeiro. FyApp muda isso oferecendo uma visão clara e em
            tempo real de suas finanças, com ferramentas intuitivas que qualquer
            um pode usar.
          </p>

          <ul className="space-y-3">
            {benefits.map((benefit, index) => (
              <li key={index} className="flex items-start gap-3 text-lg text-[#8B8B8B]">
                <Check className="w-6 h-6 text-[#7E8C54] shrink-0 mt-1" />
                <span>{benefit}</span>
              </li>
            ))}
          </ul>

          <a
            href={`${APP_URL}/login`}
            className="inline-block bg-[#7E8C54] text-white px-8 py-3.5 rounded-lg font-semibold transition-all duration-300 hover:bg-[#6B7844] hover:-translate-y-1 hover:shadow-lg border-2 border-[#7E8C54]"
          >
            Explorar Recursos
          </a>
        </div>

        <div className="relative">
          <img
            src="https://d2xsxph8kpxj0f.cloudfront.net/310419663031400654/9rw8cUK9yoafNAfP3CRw5n/fyapp-features-illustration-X72PdFpSHkAv6RDpzsykco.webp"
            alt="Recursos FyApp"
            className="w-full rounded-2xl shadow-2xl"
          />
        </div>
      </div>
    </section>
  );
}
