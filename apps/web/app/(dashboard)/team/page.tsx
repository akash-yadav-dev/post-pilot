"use client";

import { useState } from "react";
import { SectionCard } from "@/components/dashboard/section-card";
import { useActivityFeed, useTeamMembers } from "@/hooks/queries/use-social-dashboard";

export default function TeamPage() {
  const members = useTeamMembers();
  const activity = useActivityFeed();
  const [comment, setComment] = useState("");

  return (
    <div className="space-y-5 p-4 md:p-6">
      <SectionCard title="Team Collaboration" subtitle="Roles, approvals, comments, and logs">
        <div className="space-y-3">
          {members.data?.map((member) => (
            <div key={member.id} className="flex items-center justify-between rounded-xl border border-[var(--line)] bg-[var(--surface)] px-4 py-3">
              <div>
                <p className="text-sm font-medium">{member.name}</p>
                <p className="text-xs text-[var(--muted)]">{member.email}</p>
              </div>
              <select defaultValue={member.role} className="rounded-lg border border-[var(--line)] bg-[var(--bg)] px-2 py-1 text-xs capitalize">
                <option value="admin">Admin</option>
                <option value="editor">Editor</option>
                <option value="viewer">Viewer</option>
              </select>
            </div>
          ))}
        </div>
      </SectionCard>

      <div className="grid gap-5 xl:grid-cols-2">
        <SectionCard title="Approval Workflow" subtitle="Draft → Review → Approved → Scheduled">
          <div className="flex flex-wrap gap-2 text-xs">
            {[
              "Draft",
              "Review",
              "Approved",
              "Scheduled",
            ].map((stage) => (
              <span key={stage} className="rounded-full border border-[var(--line)] bg-[var(--surface)] px-3 py-1.5">{stage}</span>
            ))}
          </div>
          <textarea
            value={comment}
            onChange={(e) => setComment(e.target.value)}
            placeholder="Add review comment on a post"
            rows={4}
            className="mt-4 w-full rounded-xl border border-[var(--line)] bg-[var(--surface)] px-4 py-3 text-sm outline-none"
          />
        </SectionCard>

        <SectionCard title="Activity Logs" subtitle="Audit trail for team actions">
          <ul className="space-y-2">
            {activity.data?.map((item) => (
              <li key={item.id} className="rounded-xl border border-[var(--line)] bg-[var(--surface)] px-3 py-2">
                <p className="text-sm font-medium">{item.action}</p>
                <p className="text-xs text-[var(--muted)]">{item.actor} • {item.time}</p>
              </li>
            ))}
          </ul>
        </SectionCard>
      </div>
    </div>
  );
}
