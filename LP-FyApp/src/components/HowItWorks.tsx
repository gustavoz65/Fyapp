const steps = [
  {
    number: 1,
    title: "Crie sua Conta",
    description:
      "Registre-se em segundos com seu email ou conta Google. Personalize suas preferências e configure suas notificações.",
    image:
      "https://d2xsxph8kpxj0f.cloudfront.net/310419663031400654/9rw8cUK9yoafNAfP3CRw5n/fyapp-onboarding-illustration-WgVHrmWdMcQTz6pQ2GqxJK.webp",
  },
  {
    number: 2,
    title: "Adicione suas Contas",
    description:
      "Cadastre seus bancos, cartões e carteiras. FyApp organiza tudo em um único painel centralizado.",
    image:
      "https://d2xsxph8kpxj0f.cloudfront.net/310419663031400654/9rw8cUK9yoafNAfP3CRw5n/fyapp-onboarding-illustration-WgVHrmWdMcQTz6pQ2GqxJK.webp",
  },
  {
    number: 3,
    title: "Comece a Controlar",
    description:
      "Visualize, analise e otimize suas finanças. Defina metas e acompanhe seu progresso em tempo real.",
    image:
      "https://d2xsxph8kpxj0f.cloudfront.net/310419663031400654/9rw8cUK9yoafNAfP3CRw5n/fyapp-onboarding-illustration-WgVHrmWdMcQTz6pQ2GqxJK.webp",
  },
];

export default function HowItWorks() {
  return (
    <section id="how" className="py-32 px-8 bg-white border-t border-b border-[#E5DDD0]">
      <div className="max-w-7xl mx-auto">
        <div className="text-center mb-16">
          <h2 className="text-5xl font-bold text-[#2C2C2C] mb-4 font-playfair">
            Como Funciona
          </h2>
          <p className="text-lg text-[#8B8B8B] max-w-2xl mx-auto">
            Começar com FyApp é simples e rápido — em 3 passos você está pronto
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          {steps.map((step) => (
            <div key={step.number} className="text-center">
              <div className="w-20 h-20 bg-[#7E8C54] text-white rounded-full flex items-center justify-center text-4xl font-bold mx-auto mb-6 transition-all duration-300 hover:scale-125 hover:shadow-lg">
                {step.number}
              </div>
              <h3 className="text-2xl font-bold text-[#2C2C2C] mb-3 font-playfair">
                {step.title}
              </h3>
              <p className="text-[#8B8B8B] leading-relaxed mb-6">
                {step.description}
              </p>
              <img
                src={step.image}
                alt={`Passo ${step.number}`}
                className="w-full rounded-xl shadow-lg hover:scale-105 transition-transform duration-300"
              />
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
