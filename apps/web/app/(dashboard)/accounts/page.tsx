"use client";

import { SectionCard } from "@/components/dashboard/section-card";
import { usePlatformSummary } from "@/hooks/queries/use-social-dashboard";

export default function AccountsPage() {
  const summary = usePlatformSummary();

  return (
    <div className="space-y-5 p-4 md:p-6">
      <SectionCard title="Social Accounts" subtitle="Connect, disconnect and monitor account health">
        <div className="space-y-3">
          {summary.data?.map((item) => (
            <article key={item.platform} className="rounded-xl border border-[var(--line)] bg-[var(--surface)] p-4">
              <div className="flex flex-wrap items-center justify-between gap-2">
                <div>
                  <p className="text-sm font-semibold capitalize">{item.platform}</p>
                  <p className="text-xs text-[var(--muted)]">Permissions: publish, analytics, comments</p>
                </div>
                <span className={`rounded-full px-2 py-1 text-xs capitalize ${item.health === "healthy" ? "bg-emerald-500/15 text-emerald-500" : item.health === "warning" ? "bg-amber-500/20 text-amber-600 dark:text-amber-400" : "bg-red-500/15 text-red-500"}`}>{item.health}</span>
              </div>
              <div className="mt-3 flex flex-wrap gap-2 text-xs">
                <button className="rounded-full border border-[var(--line)] px-3 py-1.5">{item.connected ? "Disconnect" : "Connect"}</button>
                <button className="rounded-full border border-[var(--line)] px-3 py-1.5">Re-authenticate</button>
              </div>
            </article>
          ))}
        </div>
      </SectionCard>
    </div>
  );
}
