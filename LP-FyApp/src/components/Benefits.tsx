const benefitsList = [
  {
    number: "50K+",
    title: "Usuários Ativos",
    description: "Pessoas confiando em FyApp para gerenciar suas finanças",
  },
  {
    number: "R$ 2B+",
    title: "Gerenciados",
    description: "Volume total de transações processadas com segurança",
  },
  {
    number: "99.9%",
    title: "Disponibilidade",
    description: "Plataforma sempre disponível quando você precisa",
  },
  {
    number: "24/7",
    title: "Suporte",
    description: "Equipe dedicada pronta para ajudar você",
  },
];

export default function Benefits() {
  return (
    <section className="py-32 px-8 bg-gradient-to-br from-[#FAFAF8] to-[#F0EAE0]">
      <div className="max-w-7xl mx-auto">
        <div className="text-center mb-16">
          <h2 className="text-5xl font-bold text-[#2C2C2C] mb-4 font-playfair">
            Por Que Escolher FyApp?
          </h2>
          <p className="text-lg text-[#8B8B8B]">Números que falam por si</p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {benefitsList.map((benefit, index) => (
            <div
              key={index}
              className="text-center p-8 bg-white rounded-2xl transition-all duration-300 hover:shadow-lg hover:-translate-y-2"
            >
              <div className="text-5xl font-bold text-[#7E8C54] mb-2">
                {benefit.number}
              </div>
              <h3 className="text-xl font-bold text-[#2C2C2C] mb-2">
                {benefit.title}
              </h3>
              <p className="text-sm text-[#8B8B8B] leading-relaxed">
                {benefit.description}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
