"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { ThemeToggle } from "@/components/ui/theme-toggle";
import { AuthService } from "@/services/auth-service";

export function SiteHeader() {
  const pathname = usePathname();
  const router = useRouter();
  const authenticated = AuthService.isAuthenticated();

  const linkClass = (path: string) =>
    `rounded-full px-4 py-2 text-sm font-medium transition ${
      pathname === path
        ? "bg-[var(--primary)] text-[var(--bg)]"
        : "text-[var(--primary)] hover:bg-[var(--primary)]/10"
    }`;

  return (
    <header className="sticky top-0 z-50 border-b border-[var(--line)] bg-[var(--bg)]/85 backdrop-blur">
      <div className="mx-auto flex w-full max-w-6xl items-center justify-between px-6 py-4">
        <Link href="/" className="group flex items-center gap-2">
          <span className="inline-block h-7 w-7 rotate-45 rounded-md bg-[var(--primary)] transition group-hover:rotate-[225deg]" />
          <span className="text-lg font-semibold tracking-tight text-[var(--primary)]">PostPilot</span>
        </Link>

        <nav className="hidden items-center gap-2 md:flex">
          <Link href="/" className={linkClass("/")}>Home</Link>
          {authenticated && <Link href="/dashboard" className={linkClass("/dashboard")}>Dashboard</Link>}
          {!authenticated && <Link href="/login" className={linkClass("/login")}>Login</Link>}
          {!authenticated && <Link href="/signin" className={linkClass("/signin")}>Sign Up</Link>}
        </nav>

        <div className="flex items-center gap-3">
          <ThemeToggle />
          {authenticated && (
            <button
              type="button"
              onClick={async () => {
                await AuthService.logout();
                router.push("/");
              }}
              className="rounded-full border border-[var(--line)] px-4 py-2 text-xs font-semibold uppercase tracking-[0.2em] text-[var(--primary)] transition hover:bg-[var(--primary)] hover:text-[var(--bg)]"
            >
              Logout
            </button>
          )}
        </div>
      </div>
    </header>
  );
}
