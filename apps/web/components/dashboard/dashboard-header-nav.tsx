"use client";

import { Bell, Menu } from "lucide-react";

type DashboardHeaderNavProps = {
  activeLabel: string;
  activeDescription: string;
  userName?: string;
  onOpenMobileNav: () => void;
};

export function DashboardHeaderNav({
  activeLabel,
  activeDescription,
  userName,
  onOpenMobileNav,
}: DashboardHeaderNavProps) {
  return (
    <header className="sticky top-0 z-20 flex items-center justify-between gap-3 border-b border-[var(--line)] bg-[var(--surface)]/95 px-4 py-4 backdrop-blur md:px-6">
      <div className="flex items-center gap-3">
        <button
          type="button"
          aria-label="Open mobile navigation"
          onClick={onOpenMobileNav}
          className="inline-flex h-10 w-10 items-center justify-center rounded-xl border border-[var(--line)] text-[var(--primary)] md:hidden"
        >
          <Menu size={18} />
        </button>

        <div>
          <p className="text-xs uppercase tracking-[0.2em] text-[var(--muted)]">{activeLabel}</p>
          <h1 className="mt-1 text-xl font-semibold tracking-tight md:text-2xl">
            Welcome{userName ? `, ${userName}` : ""}
          </h1>
        </div>
      </div>

      <div className="flex items-center gap-2">
        <button
          type="button"
          className="inline-flex h-10 w-10 items-center justify-center rounded-xl border border-[var(--line)] text-[var(--primary)] transition hover:bg-[var(--bg)]"
          aria-label="Notifications"
        >
          <Bell size={18} />
        </button>
        <div className="hidden rounded-xl border border-[var(--line)] px-3 py-2 text-xs text-[var(--muted)] md:block">
          {activeDescription}
        </div>
      </div>
    </header>
  );
}
