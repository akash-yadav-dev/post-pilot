export type ThemeMode = "light" | "dark";

const THEME_STORAGE_KEY = "postpilot-theme";

export const ThemeService = {
  getStoredTheme(): ThemeMode | null {
    if (typeof window === "undefined") return null;

    const value = localStorage.getItem(THEME_STORAGE_KEY);
    return value === "dark" || value === "light" ? value : null;
  },

  getSystemTheme(): ThemeMode {
    if (typeof window === "undefined") return "light";

    return window.matchMedia("(prefers-color-scheme: dark)").matches
      ? "dark"
      : "light";
  },

  applyTheme(theme: ThemeMode) {
    const root = document.documentElement;
    root.classList.toggle("dark", theme === "dark");
  },

  clearStoredTheme() {
    if (typeof window === "undefined") return;

    localStorage.removeItem(THEME_STORAGE_KEY);
  },

  setTheme(theme: ThemeMode) {
    localStorage.setItem(THEME_STORAGE_KEY, theme);
    this.applyTheme(theme);
  },
};
