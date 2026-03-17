export type ThemeMode = "light" | "dark";

const THEME_STORAGE_KEY = "postpilot-theme";

export const ThemeService = {
  getStoredTheme(): ThemeMode | null {
    if (typeof window === "undefined") {
      return null;
    }

    const value = window.localStorage.getItem(THEME_STORAGE_KEY);
    if (value === "dark" || value === "light") {
      return value;
    }

    return null;
  },

  getSystemTheme(): ThemeMode {
    if (typeof window === "undefined") {
      return "light";
    }

    return window.matchMedia("(prefers-color-scheme: dark)").matches
      ? "dark"
      : "light";
  },

  getInitialTheme(): ThemeMode {
    return this.getStoredTheme() ?? this.getSystemTheme();
  },

  applyTheme(theme: ThemeMode): void {
    if (typeof document === "undefined") {
      return;
    }

    const root = document.documentElement;
    root.classList.toggle("dark", theme === "dark");
  },

  setTheme(theme: ThemeMode): void {
    if (typeof window === "undefined") {
      return;
    }

    window.localStorage.setItem(THEME_STORAGE_KEY, theme);
    this.applyTheme(theme);
  },
};
