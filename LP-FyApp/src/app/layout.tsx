import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "FyApp - Gestão Financeira Inteligente",
  description:
    "FyApp é a plataforma de gestão financeira moderna que torna simples o controle de suas transações, orçamentos e metas.",
  openGraph: {
    title: "FyApp - Gestão Financeira Inteligente",
    description:
      "Controle total das suas finanças com segurança e simplicidade em um só lugar.",
    type: "website",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="pt-BR">
      <body className="w-full overflow-x-hidden">{children}</body>
    </html>
  );
}
