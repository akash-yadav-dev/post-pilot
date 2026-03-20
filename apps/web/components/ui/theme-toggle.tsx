"use client";

import { useTheme } from "@/components/providers/theme-provider";
import { Moon, Sun } from "lucide-react";

export function ThemeToggle() {
  const { resolvedTheme, setTheme } = useTheme();

  return (
    <button
      type="button"
      onClick={() => setTheme(resolvedTheme === "dark" ? "light" : "dark")}
      className="inline-flex items-center gap-2 rounded-full border border-[var(--line)] bg-[var(--surface)] px-4 py-2 text-xs font-semibold uppercase tracking-[0.2em] text-[var(--primary)] transition hover:-translate-y-0.5 cursor-pointer"
    >
      {resolvedTheme === "dark" ? (
        <Sun height={18} width={18} />
      ) : (
        <Moon height={18} width={18} />
      )}
    </button>
  );
}
