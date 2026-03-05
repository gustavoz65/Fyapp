const footerSections = [
  {
    title: "Produto",
    links: [
      { label: "Recursos", href: "#features" },
      { label: "Segurança", href: "#security" },
      { label: "Como Funciona", href: "#how" },
    ],
  },
  {
    title: "Empresa",
    links: [
      { label: "Sobre", href: "#" },
      { label: "Blog", href: "#" },
      { label: "Contato", href: "#" },
    ],
  },
  {
    title: "Legal",
    links: [
      { label: "Privacidade", href: "#" },
      { label: "Termos", href: "#" },
      { label: "Cookies", href: "#" },
    ],
  },
];

export default function Footer() {
  return (
    <footer className="bg-[#2C2C2C] text-white py-16 px-8">
      <div className="max-w-7xl mx-auto">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8 mb-8 pb-8 border-b border-white/10">
          {/* Brand */}
          <div>
            <div className="flex items-center gap-3 mb-4">
              <div className="w-8 h-8 bg-gradient-to-br from-[#7E8C54] to-[#6B7844] rounded-md flex items-center justify-center text-white font-bold text-sm">
                Fy
              </div>
              <span className="text-xl font-bold text-[#7E8C54]">FyApp</span>
            </div>
            <p className="text-white/60 text-sm leading-relaxed">
              Gestão financeira inteligente para quem quer mais controle e
              menos preocupação.
            </p>
          </div>

          {footerSections.map((section) => (
            <div key={section.title}>
              <h4 className="text-sm font-bold mb-4 uppercase tracking-wider">
                {section.title}
              </h4>
              <ul className="space-y-3">
                {section.links.map((link) => (
                  <li key={link.label}>
                    <a
                      href={link.href}
                      className="text-white/70 text-sm transition-all duration-300 hover:text-white hover:translate-x-1 inline-block"
                    >
                      {link.label}
                    </a>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>

        <div className="flex flex-col md:flex-row justify-between items-center gap-4 text-sm text-white/60">
          <p>&copy; 2026 FyApp. Todos os direitos reservados.</p>
          <p>Feito com ❤️ para o seu bolso</p>
        </div>
      </div>
    </footer>
  );
}
