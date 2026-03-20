"use client";

import { createContext, useContext, useEffect, useState } from "react";
import { ThemeService, ThemeMode } from "@/services/theme-service";

type ThemeContextType = {
  theme: ThemeMode;
  setTheme: (theme: ThemeMode) => void;
  resolvedTheme: ThemeMode;
};

const ThemeContext = createContext<ThemeContextType | undefined>(undefined);

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const [theme, setThemeState] = useState<ThemeMode>(() => {
    return ThemeService.getStoredTheme() ?? ThemeService.getSystemTheme();
  });

  const resolvedTheme = theme;

  useEffect(() => {
    ThemeService.applyTheme(theme);

    if (!ThemeService.getStoredTheme()) {
      ThemeService.clearStoredTheme();
    }

    const media = window.matchMedia("(prefers-color-scheme: dark)");

    const listener = (e: MediaQueryListEvent) => {
      if (!ThemeService.getStoredTheme()) {
        setThemeState(e.matches ? "dark" : "light");
      }
    };

    media.addEventListener("change", listener);
    return () => media.removeEventListener("change", listener);
  }, [theme]);

  const setTheme = (newTheme: ThemeMode) => {
    setThemeState(newTheme);
    ThemeService.setTheme(newTheme);
  };

  return (
    <ThemeContext.Provider value={{ theme, setTheme, resolvedTheme }}>
      {children}
    </ThemeContext.Provider>
  );
}

export function useTheme() {
  const context = useContext(ThemeContext);

  if (!context) {
    throw new Error("useTheme must be used within ThemeProvider");
  }

  return context;
}