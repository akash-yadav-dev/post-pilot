"use client";

import { SectionCard } from "@/components/dashboard/section-card";

export default function SettingsPage() {
  return (
    <div className="space-y-5 p-4 md:p-6">
      <SectionCard title="Profile Settings" subtitle="Identity and login preferences">
        <div className="grid gap-3 md:grid-cols-2">
          <input placeholder="Full name" className="rounded-xl border border-[var(--line)] bg-[var(--surface)] px-3 py-2 text-sm outline-none" />
          <input placeholder="Email" className="rounded-xl border border-[var(--line)] bg-[var(--surface)] px-3 py-2 text-sm outline-none" />
        </div>
      </SectionCard>

      <SectionCard title="Workspace Settings" subtitle="Team and workspace defaults">
        <div className="grid gap-3 md:grid-cols-2">
          <input placeholder="Workspace name" className="rounded-xl border border-[var(--line)] bg-[var(--surface)] px-3 py-2 text-sm outline-none" />
          <input defaultValue="Asia/Kolkata" placeholder="Timezone" className="rounded-xl border border-[var(--line)] bg-[var(--surface)] px-3 py-2 text-sm outline-none" />
        </div>
      </SectionCard>

      <SectionCard title="Billing" subtitle="Plan and usage overview">
        <div className="rounded-xl border border-[var(--line)] bg-[var(--surface)] p-4">
          <p className="text-sm font-medium">Pro Workspace Plan</p>
          <p className="mt-1 text-xs text-[var(--muted)]">Renews monthly • 8 team seats in use</p>
        </div>
      </SectionCard>
    </div>
  );
}
