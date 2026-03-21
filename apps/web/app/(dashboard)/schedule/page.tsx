"use client";

import { useMemo, useState } from "react";
import { SectionCard } from "@/components/dashboard/section-card";
import { usePosts, useQueueSlots, useUpdatePostStatusMutation } from "@/hooks/queries/use-social-dashboard";

type CalendarMode = "day" | "week" | "month";

export default function SchedulePage() {
  const [mode, setMode] = useState<CalendarMode>("week");
  const [timezone, setTimezone] = useState("Asia/Kolkata");
  const [bulkText, setBulkText] = useState("");
  const scheduledPosts = usePosts({ status: "scheduled" });
  const slots = useQueueSlots();
  const updateMutation = useUpdatePostStatusMutation();

  const parsedBulkRows = useMemo(
    () => bulkText.split("\n").map((line) => line.trim()).filter(Boolean),
    [bulkText]
  );

  return (
    <div className="space-y-5 p-4 md:p-6">
      <SectionCard title="Scheduling Calendar" subtitle="Day, week, month views with queue slots">
        <div className="flex flex-wrap items-center gap-2">
          {(["day", "week", "month"] as CalendarMode[]).map((item) => (
            <button
              key={item}
              onClick={() => setMode(item)}
              className={`rounded-full border px-3 py-1.5 text-xs capitalize ${mode === item ? "border-[var(--primary)] bg-[var(--primary)] text-[var(--bg)]" : "border-[var(--line)]"}`}
            >
              {item}
            </button>
          ))}
          <input
            value={timezone}
            onChange={(e) => setTimezone(e.target.value)}
            className="ml-auto rounded-xl border border-[var(--line)] bg-[var(--surface)] px-3 py-2 text-sm outline-none"
            placeholder="Timezone"
          />
        </div>
        <div className="mt-4 grid gap-3 md:grid-cols-3">
          {(scheduledPosts.data ?? []).slice(0, 6).map((post) => (
            <article
              key={post.id}
              draggable
              onDragStart={(e) => e.dataTransfer.setData("post-id", post.id)}
              className="cursor-grab rounded-xl border border-[var(--line)] bg-[var(--surface)] p-4 active:cursor-grabbing"
            >
              <p className="text-sm font-medium">{post.title}</p>
              <p className="mt-1 text-xs text-[var(--muted)]">{post.scheduledAt ? new Date(post.scheduledAt).toLocaleString() : "No date"}</p>
            </article>
          ))}
        </div>
      </SectionCard>

      <div className="grid gap-5 xl:grid-cols-2">
        <SectionCard title="Queue Time Slots" subtitle="Drop a post into a slot to reschedule quickly">
          <div className="space-y-2">
            {slots.data?.map((slot) => (
              <div
                key={slot.id}
                onDragOver={(e) => e.preventDefault()}
                onDrop={(e) => {
                  const id = e.dataTransfer.getData("post-id");
                  if (id) updateMutation.mutate({ id, status: "scheduled" });
                }}
                className="rounded-xl border border-dashed border-[var(--line)] bg-[var(--surface)] px-4 py-3"
              >
                <p className="text-sm font-medium">{slot.label}</p>
                <p className="text-xs text-[var(--muted)]">{slot.time} • {timezone}</p>
              </div>
            ))}
          </div>
        </SectionCard>

        <SectionCard title="Bulk Scheduling" subtitle="Paste CSV rows: title,datetime,platforms">
          <textarea
            value={bulkText}
            onChange={(e) => setBulkText(e.target.value)}
            rows={10}
            className="w-full rounded-xl border border-[var(--line)] bg-[var(--surface)] px-4 py-3 text-sm outline-none"
            placeholder="Launch teaser,2026-04-02T09:00,twitter|linkedin"
          />
          <div className="mt-3 rounded-xl border border-[var(--line)] bg-[var(--surface)] p-3 text-sm">
            Parsed rows: {parsedBulkRows.length}
          </div>
        </SectionCard>
      </div>
    </div>
  );
}
