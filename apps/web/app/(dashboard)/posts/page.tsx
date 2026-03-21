"use client";

import { useMemo, useState } from "react";
import { SectionCard } from "@/components/dashboard/section-card";
import { EmptyState } from "@/components/dashboard/empty-state";
import { PostStatusPill } from "@/components/dashboard/post-status-pill";
import {
  useDeletePostMutation,
  useDuplicatePostMutation,
  usePosts,
  useUpdatePostStatusMutation,
} from "@/hooks/queries/use-social-dashboard";
import type { PostStatus } from "@/lib/social/types";

const statusTabs: Array<{ label: string; value?: PostStatus }> = [
  { label: "All" },
  { label: "Drafts", value: "draft" },
  { label: "Scheduled", value: "scheduled" },
  { label: "Published", value: "published" },
  { label: "Failed", value: "failed" },
];

export default function PostsPage() {
  const [tab, setTab] = useState<PostStatus | undefined>();
  const [platform, setPlatform] = useState("");
  const [search, setSearch] = useState("");

  const filters = useMemo(() => ({ status: tab, platform: platform || undefined, search: search || undefined }), [tab, platform, search]);
  const posts = usePosts(filters);
  const deleteMutation = useDeletePostMutation();
  const duplicateMutation = useDuplicatePostMutation();
  const updateStatusMutation = useUpdatePostStatusMutation();

  return (
    <div className="space-y-5 p-4 md:p-6">
      <SectionCard title="Posts Management" subtitle="Drafts, scheduled, published, failed with quick actions">
        <div className="flex flex-wrap items-center gap-2">
          {statusTabs.map((item) => (
            <button
              key={item.label}
              onClick={() => setTab(item.value)}
              className={`rounded-full border px-3 py-1.5 text-xs ${tab === item.value ? "border-[var(--primary)] bg-[var(--primary)] text-[var(--bg)]" : "border-[var(--line)]"}`}
            >
              {item.label}
            </button>
          ))}
        </div>
        <div className="mt-3 grid gap-2 md:grid-cols-3">
          <input value={search} onChange={(e) => setSearch(e.target.value)} placeholder="Search posts" className="rounded-xl border border-[var(--line)] bg-[var(--surface)] px-3 py-2 text-sm outline-none" />
          <select value={platform} onChange={(e) => setPlatform(e.target.value)} className="rounded-xl border border-[var(--line)] bg-[var(--surface)] px-3 py-2 text-sm outline-none">
            <option value="">All platforms</option>
            <option value="twitter">Twitter/X</option>
            <option value="linkedin">LinkedIn</option>
            <option value="facebook">Facebook</option>
            <option value="instagram">Instagram</option>
            <option value="pinterest">Pinterest</option>
          </select>
          <input type="date" className="rounded-xl border border-[var(--line)] bg-[var(--surface)] px-3 py-2 text-sm outline-none" />
        </div>
      </SectionCard>

      <SectionCard title="Results">
        {!posts.data?.length ? (
          <EmptyState title="No posts found" message="Try changing the filter or create a new post." />
        ) : (
          <div className="space-y-3">
            {posts.data.map((post) => (
              <article key={post.id} className="rounded-xl border border-[var(--line)] bg-[var(--bg)] p-4">
                <div className="flex flex-wrap items-center justify-between gap-2">
                  <div>
                    <h3 className="text-sm font-semibold">{post.title}</h3>
                    <p className="mt-1 text-xs text-[var(--muted)]">{post.platforms.join(", ")} • Updated {new Date(post.updatedAt).toLocaleString()}</p>
                  </div>
                  <PostStatusPill status={post.status} />
                </div>
                <p className="mt-2 line-clamp-2 text-sm text-[var(--muted)]">{post.content}</p>
                <div className="mt-3 flex flex-wrap gap-2 text-xs">
                  <button className="rounded-full border border-[var(--line)] px-3 py-1.5">Edit</button>
                  <button onClick={() => duplicateMutation.mutate(post.id)} className="rounded-full border border-[var(--line)] px-3 py-1.5">Duplicate</button>
                  <button onClick={() => deleteMutation.mutate(post.id)} className="rounded-full border border-red-400/50 px-3 py-1.5 text-red-500">Delete</button>
                  <button onClick={() => updateStatusMutation.mutate({ id: post.id, status: "scheduled" })} className="rounded-full border border-[var(--line)] px-3 py-1.5">Reschedule</button>
                </div>
              </article>
            ))}
          </div>
        )}
      </SectionCard>
    </div>
  );
}
