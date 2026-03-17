"use client";

import { useTheme } from "@/components/providers/theme-provider";

export function ThemeToggle() {
  const { theme, toggleTheme } = useTheme();

  return (
    <button
      type="button"
      onClick={toggleTheme}
      className="inline-flex items-center gap-2 rounded-full border border-[var(--line)] bg-[var(--surface)] px-4 py-2 text-xs font-semibold uppercase tracking-[0.2em] text-[var(--primary)] transition hover:-translate-y-0.5"
    >
      <span className="h-2 w-2 rounded-full bg-[var(--primary)]" aria-hidden />
      {theme === "dark" ? "Dark" : "Light"}
    </button>
  );
}
