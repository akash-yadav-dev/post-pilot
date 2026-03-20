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
  // ✅ Initialize state lazily (runs once, no effect needed)
  const [theme, setThemeState] = useState<ThemeMode>(() => {
    return ThemeService.getStoredTheme() ?? ThemeService.getSystemTheme();
  });
  console.log("Initial theme:", ThemeService.getStoredTheme());

  const [resolvedTheme, setResolvedTheme] = useState<ThemeMode>(theme);
  const [mounted, setMounted] = useState(false);


  useEffect(() => {
    ThemeService.applyTheme(resolvedTheme);
    
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setMounted(true);

    const media = window.matchMedia("(prefers-color-scheme: dark)");

    const listener = (e: MediaQueryListEvent) => {
      if (!ThemeService.getStoredTheme()) {
        const newTheme = e.matches ? "dark" : "light";
        setResolvedTheme(newTheme);
        ThemeService.applyTheme(newTheme);
      }
    };

    media.addEventListener("change", listener);
    return () => media.removeEventListener("change", listener);
  }, [resolvedTheme]);

  const setTheme = (newTheme: ThemeMode) => {
    setThemeState(newTheme);
    setResolvedTheme(newTheme);
    ThemeService.setTheme(newTheme);
  };

  if (!mounted) return null;

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