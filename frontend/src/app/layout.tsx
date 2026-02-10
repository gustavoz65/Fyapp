import type { Metadata } from "next";
import "./globals.css";
import { ThemeProvider } from "@/providers/theme-provider";
import { AuthProvider } from "@/providers/auth-provider";
import { Toaster } from "@/components/ui/sonner";
import { Inclusive_Sans } from "next/font/google";

const inclusiveSans = Inclusive_Sans({
  weight: ["400"],
  subsets: ["latin"],
  variable: "--font-inclusive-sans",
});

export const metadata: Metadata = {
  title: "Cashing - Gestao Financeira",
  description: "Gerencie suas financas pessoais de forma inteligente",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="pt-BR" suppressHydrationWarning>
      <body className={`antialiased ${inclusiveSans.variable}`}>
        <ThemeProvider>
          <AuthProvider>
            {children}
            <Toaster />
          </AuthProvider>
        </ThemeProvider>
      </body>
    </html>
  );
}
