const APP_URL = process.env.NEXT_PUBLIC_APP_URL || "https://fyapp.up.railway.app";

export default function GrowthSection() {
  return (
    <section className="py-32 px-8 max-w-7xl mx-auto">
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-16 items-center">
        <div className="relative order-2 lg:order-1">
          <img
            src="https://d2xsxph8kpxj0f.cloudfront.net/310419663031400654/9rw8cUK9yoafNAfP3CRw5n/fyapp-growth-illustration-J647hEMMMwFwwztfBg6Kws.webp"
            alt="Crescimento Financeiro"
            className="w-full rounded-2xl shadow-2xl hover:scale-105 transition-transform duration-300"
          />
        </div>

        <div className="space-y-6 order-1 lg:order-2">
          <h2 className="text-4xl font-bold text-[#2C2C2C] font-playfair">
            Veja Seu Dinheiro Crescer
          </h2>

          <p className="text-lg text-[#8B8B8B] leading-relaxed">
            Com FyApp, você não apenas controla seus gastos, mas também planeja
            seu futuro financeiro. Nossas ferramentas de análise ajudam você a
            identificar oportunidades de economia e investimento.
          </p>

          <p className="text-lg text-[#8B8B8B] leading-relaxed">
            Milhares de usuários já aumentaram sua saúde financeira em mais de
            40% usando FyApp. Você pode ser o próximo.
          </p>

          <a
            href={`${APP_URL}/login`}
            className="inline-block bg-[#7E8C54] text-white px-8 py-3.5 rounded-lg font-semibold transition-all duration-300 hover:bg-[#6B7844] hover:-translate-y-1 hover:shadow-lg border-2 border-[#7E8C54]"
          >
            Começar Sua Jornada
          </a>
        </div>
      </div>
    </section>
  );
}
