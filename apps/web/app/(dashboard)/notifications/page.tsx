"use client";

import { SectionCard } from "@/components/dashboard/section-card";
import {
  useMarkNotificationReadMutation,
  useNotifications,
} from "@/hooks/queries/use-social-dashboard";

export default function NotificationsPage() {
  const notifications = useNotifications();
  const markReadMutation = useMarkNotificationReadMutation();

  return (
    <div className="space-y-5 p-4 md:p-6">
      <SectionCard title="Notifications" subtitle="In-app alerts for publishing and approval flow">
        <ul className="space-y-3">
          {notifications.data?.map((item) => (
            <li key={item.id} className="rounded-xl border border-[var(--line)] bg-[var(--surface)] p-4">
              <div className="flex items-start justify-between gap-3">
                <div>
                  <p className="text-sm font-semibold">{item.title}</p>
                  <p className="mt-1 text-sm text-[var(--muted)]">{item.description}</p>
                  <p className="mt-1 text-xs text-[var(--muted)]">{item.createdAt}</p>
                </div>
                {!item.read ? (
                  <button
                    onClick={() => markReadMutation.mutate(item.id)}
                    className="rounded-full border border-[var(--line)] px-3 py-1.5 text-xs"
                  >
                    Mark as read
                  </button>
                ) : (
                  <span className="rounded-full bg-emerald-500/15 px-2.5 py-1 text-xs text-emerald-500">Read</span>
                )}
              </div>
            </li>
          ))}
        </ul>
      </SectionCard>

      <SectionCard title="Email Notifications" subtitle="Optional alerts by email">
        <label className="flex items-center justify-between rounded-xl border border-[var(--line)] bg-[var(--surface)] px-4 py-3 text-sm">
          Post published and post failed alerts
          <input type="checkbox" defaultChecked className="h-4 w-4" />
        </label>
      </SectionCard>
    </div>
  );
}
