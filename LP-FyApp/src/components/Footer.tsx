import { Heart, Github, Twitter, Linkedin, Mail } from "lucide-react";

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

const socialLinks = [
  { icon: Github, href: "#", label: "Github" },
  { icon: Twitter, href: "#", label: "Twitter" },
  { icon: Linkedin, href: "#", label: "LinkedIn" },
  { icon: Mail, href: "#", label: "Email" },
];

export default function Footer() {
  return (
    <footer className="relative bg-primary text-white overflow-hidden">
      {/* Background Pattern */}
      <div className="absolute inset-0 bg-[linear-gradient(to_right,#ffffff05_1px,transparent_1px),linear-gradient(to_bottom,#ffffff05_1px,transparent_1px)] bg-[size:24px_24px]" />

      <div className="max-w-7xl mx-auto px-6 sm:px-8 lg:px-12 py-16 relative z-10">
        {/* Main Footer Content */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-12 mb-12 pb-12 border-b border-white/10">
          {/* Brand Section */}
          <div className="lg:col-span-2">
            <div className="flex items-center gap-3 mb-6">
              <div className="relative">
                <img
                  src="/logoFy.png"
                  alt="FyApp"
                  className="h-10 w-10 object-contain"
                />
                <div className="absolute -inset-1 bg-accent/20 rounded-lg blur opacity-50" />
              </div>
              <span className="text-2xl font-display font-bold">
                Fy<span className="text-accent">App</span>
              </span>
            </div>
            <p className="text-white/70 text-sm leading-relaxed mb-6 max-w-sm">
              Gestão financeira inteligente para quem quer mais controle e
              menos preocupação. Transforme sua relação com o dinheiro.
            </p>

            {/* Social Links */}
            <div className="flex items-center gap-3">
              {socialLinks.map((social) => {
                const Icon = social.icon;
                return (
                  <a
                    key={social.label}
                    href={social.href}
                    aria-label={social.label}
                    className="w-10 h-10 rounded-lg bg-white/5 border border-white/10 flex items-center justify-center transition-all duration-300 hover:bg-accent hover:border-accent hover:-translate-y-1 group"
                  >
                    <Icon className="w-4 h-4 text-white/70 group-hover:text-white transition-colors" />
                  </a>
                );
              })}
            </div>
          </div>

          {/* Footer Links */}
          {footerSections.map((section) => (
            <div key={section.title}>
              <h4 className="text-sm font-display font-bold mb-4 uppercase tracking-wider text-white">
                {section.title}
              </h4>
              <ul className="space-y-3">
                {section.links.map((link) => (
                  <li key={link.label}>
                    <a
                      href={link.href}
                      className="text-white/70 text-sm transition-all duration-300 hover:text-accent hover:translate-x-1 inline-block group"
                    >
                      <span className="relative">
                        {link.label}
                        <span className="absolute bottom-0 left-0 w-0 h-px bg-accent group-hover:w-full transition-all duration-300" />
                      </span>
                    </a>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>

        {/* Bottom Bar */}
        <div className="flex flex-col md:flex-row justify-between items-center gap-4 text-sm">
          <p className="text-white/60">
            &copy; {new Date().getFullYear()} FyApp. Todos os direitos reservados.
          </p>
          <div className="flex items-center gap-2 text-white/60">
            <span>Feito com</span>
            <Heart className="w-4 h-4 text-red-500 fill-red-500 animate-pulse" />
            <span>para o seu bolso</span>
          </div>
        </div>
      </div>

      {/* Decorative Bottom Gradient */}
      <div className="absolute bottom-0 left-0 right-0 h-1 bg-gradient-to-r from-accent via-secondary to-blue-500" />
    </footer>
  );
}
