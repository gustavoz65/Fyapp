import type { Metadata } from "next";
import "./globals.css";
import { ThemeProvider } from "@/providers/theme-provider";
import { AuthProvider } from "@/providers/auth-provider";
import { Toaster } from "@/components/ui/sonner";
import { Playfair_Display, Inter } from "next/font/google";

// FiNext Typography: Playfair Display for headings + Inter for body
const playfairDisplay = Playfair_Display({
  weight: ["400", "700", "800"],
  subsets: ["latin"],
  variable: "--font-playfair",
  display: "swap",
});

const inter = Inter({
  subsets: ["latin"],
  variable: "--font-inter",
  display: "swap",
});

export const metadata: Metadata = {
  title: "FiNext - Gestão Financeira Inteligente",
  description: "Transforme a forma como você gerencia suas finanças com FiNext",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="pt-BR" suppressHydrationWarning className="light">
      <body className={`antialiased ${playfairDisplay.variable} ${inter.variable}`}>
        <ThemeProvider defaultTheme="light" forcedTheme="light">
          <AuthProvider>
            {children}
            <Toaster />
          </AuthProvider>
        </ThemeProvider>
      </body>
    </html>
  );
}
