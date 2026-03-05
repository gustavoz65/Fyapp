import React from 'react';

const footerSections = [
  {
    title: 'Produto',
    links: ['Recursos', 'Preços', 'Segurança', 'Roadmap'],
  },
  {
    title: 'Empresa',
    links: ['Sobre', 'Blog', 'Carreiras', 'Contato'],
  },
  {
    title: 'Legal',
    links: ['Privacidade', 'Termos', 'Cookies', 'Compliance'],
  },
  {
    title: 'Redes Sociais',
    links: ['Twitter', 'LinkedIn', 'Instagram', 'GitHub'],
  },
];

export default function Footer() {
  return (
    <footer className="bg-[#2C2C2C] text-white py-16 px-8">
      <div className="max-w-7xl mx-auto">
        {/* Footer Content */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8 mb-8 pb-8 border-b border-white/10">
          {footerSections.map((section, index) => (
            <div key={index}>
              <h4 className="text-sm font-bold mb-4 uppercase tracking-wider">{section.title}</h4>
              <ul className="space-y-3">
                {section.links.map((link, linkIndex) => (
                  <li key={linkIndex}>
                    <a
                      href="#"
                      className="text-white/70 text-sm transition-all duration-300 hover:text-white hover:translate-x-1 inline-block"
                    >
                      {link}
                    </a>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>

        {/* Footer Bottom */}
        <div className="flex flex-col md:flex-row justify-between items-center gap-4 text-sm text-white/60">
          <p>&copy; 2026 FyApp. Todos os direitos reservados.</p>
          <p>Feito com ❤️ para você</p>
        </div>
      </div>
    </footer>
  );
}
