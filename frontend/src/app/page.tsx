"use client";

import { Button } from "@/components/ui/button";
import { useAuth } from "@/providers/auth-provider";
import {
  ArrowRight,
  BarChart3,
  CheckCircle2,
  Lock,
  Shield,
  Smartphone,
  TrendingUp,
  Zap,
} from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

export default function Home() {
  const router = useRouter();
  const { isAuthenticated, isLoading, user } = useAuth();

  useEffect(() => {
    if (!isLoading && isAuthenticated && user) {
      const onboardingCompleted = localStorage.getItem(
        `onboarding_completed_${user.id}`,
      );
      router.replace(onboardingCompleted ? "/dashboard" : "/onboarding");
    }
  }, [isAuthenticated, isLoading, user, router]);

  if (isLoading) {
    return (
      <div className="flex h-screen items-center justify-center">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      <nav className="sticky top-0 z-50 bg-background/80 backdrop-blur-sm border-b border-border">
        <div className="container flex items-center justify-between h-16">
          <div className="flex items-center gap-2">
            <span className="brand-text text-2xl text-primary">Fy</span>
          </div>
          <div className="hidden md:flex items-center gap-8">
            <a
              href="#features"
              className="text-foreground/70 hover:text-foreground transition-colors"
            >
              Recursos
            </a>
            <a
              href="#how-it-works"
              className="text-foreground/70 hover:text-foreground transition-colors"
            >
              Como Funciona
            </a>
            <a
              href="#security"
              className="text-foreground/70 hover:text-foreground transition-colors"
            >
              Segurança
            </a>
          </div>
          <Link href="/login">
            <Button className="bg-primary hover:bg-primary/90 text-primary-foreground">
              Começar
            </Button>
          </Link>
        </div>
      </nav>

      <section className="relative overflow-hidden">
        <div
          className="absolute inset-0 opacity-40"
          style={{
            backgroundImage: `url('https://private-us-east-1.manuscdn.com/sessionFile/BWgvnPXo93dkkKVKYxsfpD/sandbox/oLKk66aY310sPCKnXoVGp2-img-1_1771876510000_na1fn_aGVyby1iYWNrZ3JvdW5k.png?x-oss-process=image/resize,w_1920,h_1920/format,webp/quality,q_80&Expires=1798761600&Policy=eyJTdGF0ZW1lbnQiOlt7IlJlc291cmNlIjoiaHR0cHM6Ly9wcml2YXRlLXVzLWVhc3QtMS5tYW51c2Nkbi5jb20vc2Vzc2lvbkZpbGUvQldndm5QWG85M2Rra0tWS1l4c2ZwRC9zYW5kYm94L29MS2s2NmFZMzEwc1BDS25Yb1ZHcDItaW1nLTFfMTc3MTg3NjUxMDAwMF9uYTFmbl9hR1Z5YnkxaVlXTnJaM0p2ZFc1ay5wbmc~eC1vc3MtcHJvY2Vzcz1pbWFnZS9yZXNpemUsd18xOTIwLGhfMTkyMC9mb3JtYXQsd2VicC9xdWFsaXR5LHFfODAiLCJDb25kaXRpb24iOnsiRGF0ZUxlc3NUaGFuIjp7IkFXUzpFcG9jaFRpbWUiOjE3OTg3NjE2MDB9fX1dfQ__&Key-Pair-Id=K2HSFNDJXOU9YS&Signature=YRWkz6j5QDPTpsR9n2Orwm3OvFE~OVXyH3JnigWpJCmIKiMey6KuVxyQNNQQMlqn2UrDpQc5Dkwl4V2sa845HPwkl~wElXABv9Daxwy8NaYipP9bCbv0gIb5LRcT7iv-6W~OYGAWCSdFZdV4pxP77Riir4R7VvCAexUuyQI7hcOXeWPLvLLoOqn2FY-O~yk5jijzXVms5F32wQ0XnHmygEhFTm~zW~Lbbk-ECJ~aDjs18lJA7ATmqLIiGKhGLCe5mdsD4EzSIonHLAMrgv10-ZqOfsWDToaZf3ACePAM6ZAw5jQeJzR~pCu2jfpHg3f7FGR7ocNlrYdU27GQpkLXlQ__')`,
            backgroundSize: "cover",
            backgroundPosition: "center",
          }}
        />

        <div className="container relative z-10 py-24 md:py-32">
          <div className="max-w-2xl">
            <h1 className="text-5xl md:text-7xl font-bold text-foreground mb-6 leading-tight">
              Seu Dinheiro, <span className="text-primary">Melhor</span>{" "}
              Controlado
            </h1>
            <p className="text-lg md:text-xl text-foreground/70 mb-8 leading-relaxed">
              Fy é a plataforma de gestão financeira moderna que torna simples o
              controle de suas transações, orçamentos e metas. Segurança e
              simplicidade em um só lugar.
            </p>
            <div className="flex flex-col sm:flex-row gap-4">
              <Link href="/register">
                <Button
                  size="lg"
                  className="bg-primary hover:bg-primary/90 text-primary-foreground text-base w-full sm:w-auto"
                >
                  Começar Gratuitamente <ArrowRight className="ml-2 w-5 h-5" />
                </Button>
              </Link>
              <Link href="/login">
                <Button
                  size="lg"
                  variant="outline"
                  className="border-border hover:bg-secondary w-full sm:w-auto"
                >
                  Ver Demo
                </Button>
              </Link>
            </div>
          </div>
        </div>
      </section>

      <section className="py-16 md:py-24 bg-card">
        <div className="container">
          <div className="max-w-2xl">
            <h2 className="text-3xl md:text-4xl font-bold text-foreground mb-6">
              Controle financeiro nunca foi tão fácil
            </h2>
            <p className="text-lg text-foreground/70 mb-8 leading-relaxed">
              A maioria das pessoas luta para entender seus gastos e planejar o
              futuro financeiro. Fy muda isso oferecendo uma visão clara e em
              tempo real de suas finanças, com ferramentas intuitivas que
              qualquer um pode usar.
            </p>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {[
                "Visualize todos seus gastos em um único lugar",
                "Crie orçamentos e acompanhe seu progresso",
                "Defina metas financeiras e alcance-as",
                "Receba alertas inteligentes sobre sua saúde financeira",
              ].map((item, idx) => (
                <div key={idx} className="flex items-start gap-3">
                  <CheckCircle2 className="w-5 h-5 text-primary mt-1 flex-shrink-0" />
                  <span className="text-foreground/80">{item}</span>
                </div>
              ))}
            </div>
          </div>
        </div>
      </section>

      <section id="features" className="py-16 md:py-24 bg-background">
        <div className="container">
          <div className="text-center mb-16">
            <h2 className="text-4xl md:text-5xl font-bold text-foreground mb-4">
              Recursos Poderosos
            </h2>
            <p className="text-lg text-foreground/70 max-w-2xl mx-auto">
              Tudo que você precisa para gerenciar suas finanças com confiança
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {[
              {
                icon: BarChart3,
                title: "Dashboard Inteligente",
                description:
                  "Visualize suas finanças em tempo real com gráficos interativos",
              },
              {
                icon: TrendingUp,
                title: "Análise de Gastos",
                description:
                  "Entenda seus padrões de gastos e identifique oportunidades",
              },
              {
                icon: Smartphone,
                title: "Acesso em Qualquer Lugar",
                description:
                  "Gerencie suas finanças no desktop, tablet ou smartphone",
              },
              {
                icon: Lock,
                title: "Segurança de Banco",
                description:
                  "Seus dados são protegidos com criptografia de nível militar",
              },
            ].map((feature, idx) => {
              const Icon = feature.icon;
              return (
                <div
                  key={idx}
                  className="p-6 bg-card rounded-lg border border-border hover:border-primary/50 transition-all duration-300 hover:shadow-lg hover:-translate-y-1"
                >
                  <Icon className="w-8 h-8 text-primary mb-4" />
                  <h3 className="text-lg font-bold text-foreground mb-2">
                    {feature.title}
                  </h3>
                  <p className="text-foreground/70">{feature.description}</p>
                </div>
              );
            })}
          </div>
        </div>
      </section>

      <section id="how-it-works" className="py-16 md:py-24 bg-card">
        <div className="container">
          <div className="text-center mb-16">
            <h2 className="text-4xl md:text-5xl font-bold text-foreground mb-4">
              Como Funciona
            </h2>
            <p className="text-lg text-foreground/70 max-w-2xl mx-auto">
              Três passos simples para começar a controlar suas finanças
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-8 md:gap-12">
            {[
              {
                step: "1",
                title: "Crie sua Conta",
                description:
                  "Cadastro rápido e seguro. Leva menos de 2 minutos.",
              },
              {
                step: "2",
                title: "Conecte suas Contas",
                description:
                  "Integre suas contas bancárias de forma segura e automática.",
              },
              {
                step: "3",
                title: "Comece a Controlar",
                description:
                  "Visualize, analise e controle suas finanças em tempo real.",
              },
            ].map((item, idx) => (
              <div key={idx} className="text-center">
                <div className="w-16 h-16 rounded-full bg-primary text-white flex items-center justify-center font-bold text-2xl mx-auto mb-4">
                  {item.step}
                </div>
                <h3 className="text-xl font-bold text-foreground mb-2">
                  {item.title}
                </h3>
                <p className="text-foreground/70">{item.description}</p>
                {idx < 2 && (
                  <div className="hidden md:flex justify-center mt-6">
                    <ArrowRight className="w-6 h-6 text-primary/50 rotate-90" />
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>
      </section>

      <section id="security" className="py-16 md:py-24 bg-background">
        <div className="container">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-12 items-center">
            <div>
              <h2 className="text-4xl md:text-5xl font-bold text-foreground mb-6">
                Segurança que você pode confiar
              </h2>
              <p className="text-lg text-foreground/70 mb-8 leading-relaxed">
                Seus dados financeiros são os mais sensíveis. Por isso,
                implementamos os mais altos padrões de segurança da indústria.
              </p>
              <div className="space-y-4">
                {[
                  "Criptografia end-to-end de todos os dados",
                  "Autenticação de dois fatores (2FA)",
                  "Conformidade com LGPD e regulamentações internacionais",
                  "Auditorias de segurança regulares por terceiros",
                ].map((item, idx) => (
                  <div key={idx} className="flex items-start gap-3">
                    <Shield className="w-5 h-5 text-primary mt-1 flex-shrink-0" />
                    <span className="text-foreground/80">{item}</span>
                  </div>
                ))}
              </div>
            </div>
            <div className="bg-card p-8 rounded-lg border border-border">
              <div className="space-y-6">
                <div className="flex items-center gap-4">
                  <Lock className="w-8 h-8 text-primary" />
                  <div>
                    <h4 className="font-bold text-foreground">
                      Proteção SSL/TLS
                    </h4>
                    <p className="text-sm text-foreground/70">
                      Todas as conexões são criptografadas
                    </p>
                  </div>
                </div>
                <div className="flex items-center gap-4">
                  <Zap className="w-8 h-8 text-primary" />
                  <div>
                    <h4 className="font-bold text-foreground">
                      Backup Automático
                    </h4>
                    <p className="text-sm text-foreground/70">
                      Seus dados são salvos continuamente
                    </p>
                  </div>
                </div>
                <div className="flex items-center gap-4">
                  <Shield className="w-8 h-8 text-primary" />
                  <div>
                    <h4 className="font-bold text-foreground">
                      Monitoramento 24/7
                    </h4>
                    <p className="text-sm text-foreground/70">
                      Detecção de atividades suspeitas
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section className="py-16 md:py-24 bg-primary text-primary-foreground relative overflow-hidden">
        <div
          className="absolute inset-0 opacity-20"
          style={{
            backgroundImage: `url('https://private-us-east-1.manuscdn.com/sessionFile/BWgvnPXo93dkkKVKYxsfpD/sandbox/oLKk66aY310sPCKnXoVGp2-img-4_1771876505000_na1fn_Y3RhLWJhY2tncm91bmQ.png?x-oss-process=image/resize,w_1920,h_1920/format,webp/quality,q_80&Expires=1798761600&Policy=eyJTdGF0ZW1lbnQiOlt7IlJlc291cmNlIjoiaHR0cHM6Ly9wcml2YXRlLXVzLWVhc3QtMS5tYW51c2Nkbi5jb20vc2Vzc2lvbkZpbGUvQldndm5QWG85M2Rra0tWS1l4c2ZwRC9zYW5kYm94L29MS2s2NmFZMzEwc1BDS25Yb1ZHcDItaW1nLTRfMTc3MTg3NjUwNTAwMF9uYTFmbl9ZM1JoTFdKaFkydG5jbTkxYm1RLnBuZz94LW9zcy1wcm9jZXNzPWltYWdlL3Jlc2l6ZSx3XzE5MjAsaF8xOTIwL2Zvcm1hdCx3ZWJwL3F1YWxpdHkscV84MCIsIkNvbmRpdGlvbiI6eyJEYXRlTGVzc1RoYW4iOnsiQVdTOkVwb2NoVGltZSI6MTc5ODc2MTYwMH19fV19&Key-Pair-Id=K2HSFNDJXOU9YS&Signature=KvuBxA0mmS~XpNXfv-Fo9PNoeOByFSFSei-kLkUXBCUkJlfcPuD1wHHrAae-Rgb3MREInWB8CRb8KY63s5d~u8EPlW3~Bx1ovrnHQ5yOYeSntwu0bFKxFopoOhONSZ0Re7iuu3Lwkd2gAiSfGX5z0ZRJtOTs7ClCS-LtKnhEmC-5aU01hJPvIXqHZZYiOWX5PGPM9pE1iMJmc90Mio~x9FIXd5MkpwhFdnp33Vjl4sFSlZ3FZuw6Mn7OqZzS4CMER529lg-RPWlvJ~E6ZW4xEsnH~XpqN8kv3QhMypozSlHrye2TgwyVHIEXL3wUU545emACZXThVwdeYQJg2ED~Fw__')`,
            backgroundSize: "cover",
            backgroundPosition: "center",
          }}
        />
        <div className="container relative z-10 text-center">
          <h2 className="text-4xl md:text-5xl font-bold mb-4">
            Pronto para Transformar suas Finanças?
          </h2>
          <p className="text-lg mb-8 opacity-90 max-w-2xl mx-auto">
            Junte-se a milhares de usuários que já estão controlando suas
            finanças com Fy
          </p>
          <Link href="/register">
            <Button
              size="lg"
              className="bg-white hover:bg-gray-100 text-primary font-bold text-base"
            >
              Começar Agora <ArrowRight className="ml-2 w-5 h-5" />
            </Button>
          </Link>
        </div>
      </section>

      <footer className="bg-card border-t border-border py-12">
        <div className="container">
          <div className="grid grid-cols-1 md:grid-cols-4 gap-8 mb-8">
            <div>
              <div className="mb-4">
                <span className="brand-text text-xl text-primary">Fy</span>
              </div>
              <p className="text-sm text-foreground/70">
                Gestão financeira moderna e segura
              </p>
            </div>
            <div>
              <h4 className="font-bold text-foreground mb-4">Produto</h4>
              <ul className="space-y-2 text-sm text-foreground/70">
                <li>
                  <a
                    href="#features"
                    className="hover:text-primary transition-colors"
                  >
                    Recursos
                  </a>
                </li>
                <li>
                  <a
                    href="#how-it-works"
                    className="hover:text-primary transition-colors"
                  >
                    Como Funciona
                  </a>
                </li>
                <li>
                  <a
                    href="#security"
                    className="hover:text-primary transition-colors"
                  >
                    Segurança
                  </a>
                </li>
              </ul>
            </div>
            <div>
              <h4 className="font-bold text-foreground mb-4">Empresa</h4>
              <ul className="space-y-2 text-sm text-foreground/70">
                <li>
                  <a href="#" className="hover:text-primary transition-colors">
                    Sobre
                  </a>
                </li>
                <li>
                  <a href="#" className="hover:text-primary transition-colors">
                    Blog
                  </a>
                </li>
                <li>
                  <a href="#" className="hover:text-primary transition-colors">
                    Contato
                  </a>
                </li>
              </ul>
            </div>
            <div>
              <h4 className="font-bold text-foreground mb-4">Legal</h4>
              <ul className="space-y-2 text-sm text-foreground/70">
                <li>
                  <a href="#" className="hover:text-primary transition-colors">
                    Privacidade
                  </a>
                </li>
                <li>
                  <a href="#" className="hover:text-primary transition-colors">
                    Termos
                  </a>
                </li>
                <li>
                  <a href="#" className="hover:text-primary transition-colors">
                    Cookies
                  </a>
                </li>
              </ul>
            </div>
          </div>
          <div className="border-t border-border pt-8 text-center text-sm text-foreground/70">
            <p>&copy; 2025 Fy. Todos os direitos reservados.</p>
          </div>
        </div>
      </footer>
    </div>
  );
}
