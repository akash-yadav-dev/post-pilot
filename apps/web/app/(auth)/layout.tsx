import Link from "next/link";
import { Leckerli_One } from "next/font/google";

const leckerliOne = Leckerli_One({
  subsets: ["latin"],
  weight: "400",
  variable: "--font-leckerli-one",
});

export default function AuthLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className={`${leckerliOne.variable} min-h-screen bg-[var(--bg)] text-[var(--primary)]`}>
      <div className="min-h-screen md:grid md:grid-cols-[1fr_1.15fr]">

        {/* ─── Left: Brand Panel ──────────────────────────────── */}
        <div className="hidden md:flex md:flex-col md:justify-between bg-[var(--surface)] border-r border-[var(--line)] px-12 py-10">

          {/* Logo */}
          <Link href="/" className="group inline-flex items-center gap-2 w-fit">
            <span className="inline-flex h-8 w-8 rotate-45 items-center justify-center rounded-lg bg-[var(--icon-bg)] transition-transform duration-300 ease-in-out group-hover:rotate-[405deg]">
              <span className="-rotate-45 text-sm font-semibold text-[var(--icon-text)] [font-family:var(--font-leckerli-one)]">
                P
              </span>
            </span>
            <span className="ps-1 text-base font-semibold tracking-tight text-[var(--primary)]">
              PostPilot
            </span>
          </Link>

          {/* Hero copy */}
          <div className="space-y-5">
            <span className="inline-flex items-center gap-2 rounded-full border border-[var(--line)] px-3 py-1.5">
              <span className="h-2 w-2 animate-pulse rounded-full bg-emerald-500" />
              <span className="text-xs text-[var(--muted)]">Live scheduling platform</span>
            </span>

            <h2 className="text-4xl font-semibold leading-[1.2] tracking-tight">
              Schedule once.
              <br />
              Publish everywhere.
            </h2>

            <p className="max-w-sm leading-relaxed text-[var(--muted)]">
              PostPilot gives creators and teams a sharp command center for
              content planning, scheduling, and cross-platform publishing.
            </p>

            {/* Stats row */}
            <div className="grid grid-cols-3 gap-5 border-t border-[var(--line)] pt-6">
              {[
                { value: "99.2%", label: "Publish rate" },
                { value: "12+",   label: "Platforms"    },
                { value: "10k+",  label: "Posts sent"   },
              ].map((stat) => (
                <div key={stat.label}>
                  <p className="text-2xl font-semibold">{stat.value}</p>
                  <p className="mt-1 text-xs text-[var(--muted)]">{stat.label}</p>
                </div>
              ))}
            </div>
          </div>

          <p className="text-xs text-[var(--muted)]">© 2026 PostPilot. All rights reserved.</p>
        </div>

        {/* ─── Right: Form Area ──────────────────────────────── */}
        <div className="flex min-h-screen flex-col items-center justify-center px-6 py-12 md:px-14">
          {/* Mobile logo */}
          <Link href="/" className="group mb-10 inline-flex items-center gap-2 md:hidden">
            <span className="inline-flex h-7 w-7 rotate-45 items-center justify-center rounded-md bg-[var(--icon-bg)]">
              <span className="-rotate-45 text-sm font-semibold text-[var(--icon-text)] [font-family:var(--font-leckerli-one)]">
                P
              </span>
            </span>
            <span className="ps-1 text-base font-semibold text-[var(--primary)]">PostPilot</span>
          </Link>

          <div className="w-full max-w-md">
            {children}
          </div>
        </div>

      </div>
    </div>
  );
}
