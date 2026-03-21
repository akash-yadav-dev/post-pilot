"use client";

import { useMemo, useState } from "react";
import { SectionCard } from "@/components/dashboard/section-card";
import { useMediaAssets } from "@/hooks/queries/use-social-dashboard";

export default function LibraryPage() {
  const assets = useMediaAssets();
  const [tagFilter, setTagFilter] = useState("");

  const filtered = useMemo(() => {
    if (!tagFilter) return assets.data ?? [];
    return (assets.data ?? []).filter((asset) => asset.tags.some((tag) => tag.toLowerCase().includes(tagFilter.toLowerCase())));
  }, [assets.data, tagFilter]);

  return (
    <div className="space-y-5 p-4 md:p-6">
      <SectionCard title="Content Library" subtitle="Manage reusable images and videos">
        <input
          value={tagFilter}
          onChange={(e) => setTagFilter(e.target.value)}
          placeholder="Filter by tag"
          className="w-full max-w-sm rounded-xl border border-[var(--line)] bg-[var(--surface)] px-3 py-2 text-sm outline-none"
        />
        <div className="mt-4 grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
          {filtered.map((asset) => (
            <article key={asset.id} className="overflow-hidden rounded-xl border border-[var(--line)] bg-[var(--surface)]">
              <img src={asset.url} alt={asset.name} className="h-40 w-full object-cover" />
              <div className="p-3">
                <p className="text-sm font-medium">{asset.name}</p>
                <p className="mt-1 text-xs text-[var(--muted)]">{asset.type} • {asset.uploadedAt}</p>
                <div className="mt-2 flex flex-wrap gap-1">
                  {asset.tags.map((tag) => (
                    <span key={tag} className="rounded-full border border-[var(--line)] px-2 py-0.5 text-xs">#{tag}</span>
                  ))}
                </div>
              </div>
            </article>
          ))}
        </div>
      </SectionCard>
    </div>
  );
}
