"use client";

import { ThemeProvider as NextThemesProvider } from "next-themes";

export function ThemeProvider({ children, defaultTheme, forcedTheme }: {
  children: React.ReactNode;
  defaultTheme?: string;
  forcedTheme?: string;
}) {
  return (
    <NextThemesProvider
      attribute="class"
      defaultTheme={defaultTheme || "light"}
      forcedTheme={forcedTheme}
      enableSystem={false}
      disableTransitionOnChange
    >
      {children}
    </NextThemesProvider>
  );
}
