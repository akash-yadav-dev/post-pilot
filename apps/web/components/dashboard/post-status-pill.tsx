import type { PostStatus } from "@/lib/social/types";

const statusClass: Record<PostStatus, string> = {
  draft: "bg-[var(--line)] text-[var(--muted)]",
  scheduled: "bg-blue-500/15 text-blue-500",
  published: "bg-emerald-500/15 text-emerald-500",
  failed: "bg-red-500/15 text-red-500",
  queued: "bg-amber-500/20 text-amber-600 dark:text-amber-400",
  review: "bg-violet-500/15 text-violet-500",
  approved: "bg-cyan-500/15 text-cyan-500",
};

export function PostStatusPill({ status }: { status: PostStatus }) {
  return <span className={`rounded-full px-2.5 py-1 text-xs font-medium capitalize ${statusClass[status]}`}>{status}</span>;
}
