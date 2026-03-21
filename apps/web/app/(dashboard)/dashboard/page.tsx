"use client";

import { motion } from "framer-motion";
import { MetricCard } from "@/components/dashboard/metric-card";
import { SectionCard } from "@/components/dashboard/section-card";
import { SkeletonBlock } from "@/components/dashboard/skeleton-block";
import { EmptyState } from "@/components/dashboard/empty-state";
import { PostStatusPill } from "@/components/dashboard/post-status-pill";
import {
  useActivityFeed,
  useDashboardMetrics,
  usePlatformSummary,
  useUpcomingPosts,
} from "@/hooks/queries/use-social-dashboard";

export default function DashboardPage() {
  const metrics = useDashboardMetrics();
  const upcoming = useUpcomingPosts();
  const platforms = usePlatformSummary();
  const activity = useActivityFeed();

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      className="space-y-5 p-4 md:p-6"
    >
      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        <MetricCard label="Scheduled Posts" value={metrics.data?.scheduledCount ?? "--"} note="Queued and upcoming" />
        <MetricCard label="Published Posts" value={metrics.data?.publishedCount ?? "--"} note="Successfully delivered" />
        <MetricCard label="Engagement" value={metrics.data?.engagement ?? "--"} note="Likes, comments, shares" />
        <MetricCard label="Failed Posts" value={metrics.data?.failedCount ?? "--"} note="Needs attention" />
      </div>

      <div className="grid gap-5 xl:grid-cols-3">
        <SectionCard title="Upcoming Schedule" subtitle="Next items in publishing queue">
          {upcoming.isLoading ? (
            <div className="space-y-2">
              <SkeletonBlock className="h-14 w-full" />
              <SkeletonBlock className="h-14 w-full" />
              <SkeletonBlock className="h-14 w-full" />
            </div>
          ) : upcoming.data?.length ? (
            <ul className="space-y-2">
              {upcoming.data.slice(0, 4).map((post) => (
                <li key={post.id} className="rounded-xl border border-[var(--line)] bg-[var(--surface)] px-4 py-3">
                  <div className="flex flex-wrap items-center justify-between gap-2">
                    <p className="text-sm font-medium">{post.title}</p>
                    <PostStatusPill status={post.status} />
                  </div>
                  <p className="mt-1 text-xs text-[var(--muted)]">
                    {post.platforms.join(", ")} • {post.scheduledAt ? new Date(post.scheduledAt).toLocaleString() : "No date"}
                  </p>
                </li>
              ))}
            </ul>
          ) : (
            <EmptyState title="No scheduled posts" message="Create a post and choose schedule or queue to populate this list." />
          )}
        </SectionCard>

        <SectionCard title="Platform Summary" subtitle="Connected channels and health">
          {platforms.isLoading ? (
            <div className="space-y-2">
              <SkeletonBlock className="h-12 w-full" />
              <SkeletonBlock className="h-12 w-full" />
              <SkeletonBlock className="h-12 w-full" />
            </div>
          ) : (
            <ul className="space-y-2">
              {platforms.data?.map((item) => (
                <li key={item.platform} className="flex items-center justify-between rounded-xl border border-[var(--line)] bg-[var(--surface)] px-3 py-2">
                  <div>
                    <p className="text-sm font-medium capitalize">{item.platform}</p>
                    <p className="text-xs text-[var(--muted)]">{item.posts} posts • {item.engagementRate}% engagement</p>
                  </div>
                  <span className={`rounded-full px-2 py-1 text-xs capitalize ${item.health === "healthy" ? "bg-emerald-500/15 text-emerald-500" : item.health === "warning" ? "bg-amber-500/20 text-amber-600 dark:text-amber-400" : "bg-red-500/15 text-red-500"}`}>
                    {item.health}
                  </span>
                </li>
              ))}
            </ul>
          )}
        </SectionCard>

        <SectionCard title="Activity Feed" subtitle="Latest team and system actions">
          {activity.isLoading ? (
            <div className="space-y-2">
              <SkeletonBlock className="h-12 w-full" />
              <SkeletonBlock className="h-12 w-full" />
              <SkeletonBlock className="h-12 w-full" />
            </div>
          ) : (
            <ul className="space-y-2">
              {activity.data?.map((item) => (
                <li key={item.id} className="rounded-xl border border-[var(--line)] bg-[var(--surface)] px-3 py-2">
                  <p className="text-sm font-medium">{item.action}</p>
                  <p className="text-xs text-[var(--muted)]">{item.actor} • {item.time}</p>
                </li>
              ))}
            </ul>
          )}
        </SectionCard>
      </div>
    </motion.div>
  );
}
