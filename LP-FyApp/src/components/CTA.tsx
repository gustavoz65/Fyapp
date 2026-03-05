const APP_URL = process.env.NEXT_PUBLIC_APP_URL || "https://fyapp-production.up.railway.app";

export default function CTA() {
  return (
    <section className="py-32 px-8 bg-gradient-to-br from-[#7E8C54] to-[#6B7844] text-white relative overflow-hidden">
      <div className="absolute top-0 right-0 w-96 h-96 bg-gradient-radial from-white/10 to-transparent rounded-full -mr-48 -mt-48" />
      <div className="absolute bottom-0 left-0 w-96 h-96 bg-gradient-radial from-white/10 to-transparent rounded-full -ml-48 -mb-48" />

      <div className="max-w-3xl mx-auto text-center relative z-10">
        <h2 className="text-5xl font-bold mb-6 font-playfair">
          Pronto para Transformar suas Finanças?
        </h2>
        <p className="text-xl mb-8 opacity-95 leading-relaxed">
          Junte-se a milhares de usuários que já estão controlando suas finanças
          com FyApp. Comece grátis hoje mesmo.
        </p>
        <a
          href={`${APP_URL}/login`}
          className="inline-block bg-white text-[#7E8C54] px-10 py-4 rounded-lg font-bold text-lg transition-all duration-300 hover:bg-transparent hover:text-white hover:-translate-y-1 hover:shadow-lg border-2 border-white"
        >
          Começar Grátis Agora
        </a>
      </div>
    </section>
  );
}
