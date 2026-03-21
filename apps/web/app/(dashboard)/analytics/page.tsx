"use client";

import { SectionCard } from "@/components/dashboard/section-card";
import { useAnalytics, usePosts } from "@/hooks/queries/use-social-dashboard";

export default function AnalyticsPage() {
  const analytics = useAnalytics();
  const publishedPosts = usePosts({ status: "published" });
  const topPost = publishedPosts.data?.slice().sort((a, b) => (b.likes + b.comments + b.shares) - (a.likes + a.comments + a.shares))[0];

  return (
    <div className="space-y-5 p-4 md:p-6">
      <SectionCard title="Performance Metrics" subtitle="Engagement rate, reach, impressions and CTR trends">
        <div className="grid gap-4 md:grid-cols-3">
          <article className="rounded-xl border border-[var(--line)] bg-[var(--surface)] p-4">
            <p className="text-xs uppercase tracking-[0.16em] text-[var(--muted)]">Engagement Rate</p>
            <p className="mt-2 text-2xl font-semibold">{analytics.data?.at(-1)?.engagementRate ?? "--"}%</p>
          </article>
          <article className="rounded-xl border border-[var(--line)] bg-[var(--surface)] p-4">
            <p className="text-xs uppercase tracking-[0.16em] text-[var(--muted)]">Reach</p>
            <p className="mt-2 text-2xl font-semibold">{analytics.data?.at(-1)?.reach ?? "--"}</p>
          </article>
          <article className="rounded-xl border border-[var(--line)] bg-[var(--surface)] p-4">
            <p className="text-xs uppercase tracking-[0.16em] text-[var(--muted)]">CTR</p>
            <p className="mt-2 text-2xl font-semibold">{analytics.data?.at(-1)?.ctr ?? "--"}%</p>
          </article>
        </div>

        <div className="mt-5 grid gap-2">
          {analytics.data?.map((point) => (
            <div key={point.label} className="grid grid-cols-[40px_1fr_auto] items-center gap-3">
              <span className="text-xs text-[var(--muted)]">{point.label}</span>
              <div className="h-2 rounded-full bg-[var(--line)]">
                <div className="h-2 rounded-full bg-[var(--primary)]" style={{ width: `${Math.min(100, point.engagementRate * 14)}%` }} />
              </div>
              <span className="text-xs font-medium">{point.engagementRate}%</span>
            </div>
          ))}
        </div>
      </SectionCard>

      <SectionCard title="Top Performing Post" subtitle="Best engagement from published content">
        {topPost ? (
          <article className="rounded-xl border border-[var(--line)] bg-[var(--surface)] p-4">
            <p className="text-sm font-semibold">{topPost.title}</p>
            <p className="mt-2 text-sm text-[var(--muted)]">{topPost.content}</p>
            <p className="mt-2 text-xs text-[var(--muted)]">Likes {topPost.likes} • Comments {topPost.comments} • Shares {topPost.shares}</p>
          </article>
        ) : (
          <p className="text-sm text-[var(--muted)]">No published post data yet.</p>
        )}

        <div className="mt-4 flex flex-wrap gap-2 text-xs">
          <button className="rounded-full border border-[var(--line)] px-3 py-1.5">Export CSV</button>
          <button className="rounded-full border border-[var(--line)] px-3 py-1.5">Export PDF</button>
        </div>
      </SectionCard>
    </div>
  );
}
