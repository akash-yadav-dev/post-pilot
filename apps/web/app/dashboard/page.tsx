"use client";

import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { AuthService, type AuthUser } from "@/services/auth-service";

type DashboardState = {
  loading: boolean;
  error: string | null;
  user: AuthUser | null;
};

export default function DashboardPage() {
  const router = useRouter();
  const [state, setState] = useState<DashboardState>({
    loading: true,
    error: null,
    user: null,
  });

  useEffect(() => {
    const boot = async () => {
      if (!AuthService.isAuthenticated()) {
        router.replace("/login");
        return;
      }

      try {
        const user = await AuthService.me();
        setState({ loading: false, error: null, user });
      } catch (err) {
        AuthService.clearTokens();
        setState({
          loading: false,
          error: err instanceof Error ? err.message : "Unable to load dashboard",
          user: null,
        });
        router.replace("/login");
      }
    };

    void boot();
  }, [router]);

  if (state.loading) {
    return <div className="mx-auto max-w-6xl px-6 py-16 text-sm text-[var(--muted)]">Loading dashboard...</div>;
  }

  return (
    <div className="mx-auto flex w-full max-w-6xl flex-col gap-8 px-6 py-12">
      <section className="rounded-3xl border border-[var(--line)] bg-[var(--surface)] p-8">
        <h1 className="text-3xl font-semibold tracking-tight">Dashboard</h1>
        <p className="mt-2 text-sm text-[var(--muted)]">
          Welcome{state.user ? `, ${state.user.name}` : ""}. Your publishing operations start here.
        </p>
      </section>

      <section className="grid gap-4 md:grid-cols-3">
        {[
          { label: "Scheduled Posts", value: "12", note: "Across all connected platforms" },
          { label: "Pending Approvals", value: "4", note: "Ready for editorial review" },
          { label: "Publish Success", value: "99.2%", note: "Based on the last 30 days" },
        ].map((card) => (
          <article key={card.label} className="rounded-2xl border border-[var(--line)] bg-[var(--surface)] p-6">
            <p className="text-xs uppercase tracking-[0.2em] text-[var(--muted)]">{card.label}</p>
            <p className="mt-3 text-3xl font-semibold">{card.value}</p>
            <p className="mt-2 text-sm text-[var(--muted)]">{card.note}</p>
          </article>
        ))}
      </section>

      {state.error && <p className="rounded-xl border border-red-400/40 bg-red-500/10 px-3 py-2 text-sm text-red-500">{state.error}</p>}
    </div>
  );
}
