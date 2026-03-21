"use client";

import { useMemo, useState } from "react";
import { motion } from "framer-motion";
import type { PostStatus, SocialPlatform } from "@/lib/social/types";
import { SectionCard } from "@/components/dashboard/section-card";
import { useCreatePostMutation } from "@/hooks/queries/use-social-dashboard";

const platformLimits: Record<SocialPlatform, number> = {
  twitter: 280,
  linkedin: 3000,
  facebook: 63206,
  instagram: 2200,
  pinterest: 500,
};

const hashtagSuggestions = ["#socialmedia", "#growth", "#marketing", "#contentops", "#creator"];
const mentionSuggestions = ["@productteam", "@marketing", "@community"];

export default function CreatePostPage() {
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [mediaUrls, setMediaUrls] = useState<string[]>([]);
  const [scheduleAt, setScheduleAt] = useState("");
  const [urlInput, setUrlInput] = useState("");
  const [selectedPlatforms, setSelectedPlatforms] = useState<SocialPlatform[]>(["twitter", "linkedin"]);
  const [activePlatform, setActivePlatform] = useState<SocialPlatform>("twitter");
  const [perPlatformCopy, setPerPlatformCopy] = useState<Partial<Record<SocialPlatform, string>>>({});

  const createMutation = useCreatePostMutation();

  const effectiveContent = perPlatformCopy[activePlatform] ?? content;
  const remaining = platformLimits[activePlatform] - effectiveContent.length;

  const linkPreview = useMemo(() => {
    if (!urlInput.startsWith("http")) return null;
    return {
      title: "Generated Link Preview",
      description: "This simulates how your link appears on social platforms.",
      url: urlInput,
    };
  }, [urlInput]);

  function togglePlatform(platform: SocialPlatform) {
    setSelectedPlatforms((prev) =>
      prev.includes(platform) ? prev.filter((p) => p !== platform) : [...prev, platform]
    );
  }

  function createByStatus(status: PostStatus) {
    return async () => {
      if (!title || !content || selectedPlatforms.length === 0) return;
      await createMutation.mutateAsync({
        title,
        content,
        mediaUrls,
        platforms: selectedPlatforms,
        scheduledAt: scheduleAt || undefined,
        status,
      });
      setTitle("");
      setContent("");
      setMediaUrls([]);
      setScheduleAt("");
      setUrlInput("");
      setPerPlatformCopy({});
    };
  }

  return (
    <div className="grid gap-5 p-4 md:grid-cols-12 md:p-6">
      <motion.div className="space-y-5 md:col-span-8" initial={{ opacity: 0, y: 14 }} animate={{ opacity: 1, y: 0 }}>
        <SectionCard title="Composer" subtitle="Rich post editor with platform customization">
          <div className="space-y-4">
            <input
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="Post title"
              className="w-full rounded-xl border border-[var(--line)] bg-[var(--surface)] px-4 py-3 text-sm outline-none focus:border-[var(--primary)]"
            />
            <textarea
              value={content}
              onChange={(e) => setContent(e.target.value)}
              placeholder="Write your post..."
              rows={7}
              className="w-full rounded-xl border border-[var(--line)] bg-[var(--surface)] px-4 py-3 text-sm outline-none focus:border-[var(--primary)]"
            />
            <div className="flex flex-wrap gap-2 text-xs">
              <button onClick={() => setContent((c) => `${c} 😊`)} className="rounded-full border border-[var(--line)] px-3 py-1.5">Add Emoji</button>
              {hashtagSuggestions.map((tag) => (
                <button key={tag} onClick={() => setContent((c) => `${c} ${tag}`)} className="rounded-full border border-[var(--line)] px-3 py-1.5">
                  {tag}
                </button>
              ))}
              {mentionSuggestions.map((mention) => (
                <button key={mention} onClick={() => setContent((c) => `${c} ${mention}`)} className="rounded-full border border-[var(--line)] px-3 py-1.5">
                  {mention}
                </button>
              ))}
            </div>
          </div>
        </SectionCard>

        <SectionCard title="Platform Customization" subtitle="Tailor copy and respect character limits">
          <div className="flex flex-wrap gap-2">
            {(Object.keys(platformLimits) as SocialPlatform[]).map((platform) => (
              <button
                key={platform}
                onClick={() => {
                  togglePlatform(platform);
                  setActivePlatform(platform);
                }}
                className={`rounded-full border px-3 py-1.5 text-xs capitalize ${selectedPlatforms.includes(platform) ? "border-[var(--primary)] bg-[var(--primary)] text-[var(--bg)]" : "border-[var(--line)]"}`}
              >
                {platform}
              </button>
            ))}
          </div>
          <div className="mt-4">
            <textarea
              value={perPlatformCopy[activePlatform] ?? content}
              onChange={(e) => setPerPlatformCopy((prev) => ({ ...prev, [activePlatform]: e.target.value }))}
              rows={4}
              className="w-full rounded-xl border border-[var(--line)] bg-[var(--surface)] px-4 py-3 text-sm outline-none focus:border-[var(--primary)]"
              placeholder={`Custom copy for ${activePlatform}`}
            />
            <p className={`mt-2 text-xs ${remaining < 0 ? "text-red-500" : "text-[var(--muted)]"}`}>
              {remaining} characters remaining for {activePlatform}
            </p>
          </div>
        </SectionCard>
      </motion.div>

      <motion.aside className="space-y-5 md:col-span-4" initial={{ opacity: 0, y: 14 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.1 }}>
        <SectionCard title="Media + Links" subtitle="Upload assets and generate preview">
          <input
            type="url"
            value={urlInput}
            onChange={(e) => setUrlInput(e.target.value)}
            placeholder="Paste URL for preview"
            className="w-full rounded-xl border border-[var(--line)] bg-[var(--surface)] px-4 py-3 text-sm outline-none focus:border-[var(--primary)]"
          />
          <input
            type="url"
            placeholder="Paste image/video URL"
            className="mt-3 w-full rounded-xl border border-[var(--line)] bg-[var(--surface)] px-4 py-3 text-sm outline-none focus:border-[var(--primary)]"
            onKeyDown={(e) => {
              if (e.key === "Enter") {
                e.preventDefault();
                const value = (e.target as HTMLInputElement).value;
                if (value.startsWith("http")) {
                  setMediaUrls((prev) => [...prev, value]);
                  (e.target as HTMLInputElement).value = "";
                }
              }
            }}
          />
          {linkPreview ? (
            <div className="mt-3 rounded-xl border border-[var(--line)] bg-[var(--surface)] p-3">
              <p className="text-sm font-medium">{linkPreview.title}</p>
              <p className="mt-1 text-xs text-[var(--muted)]">{linkPreview.description}</p>
              <p className="mt-2 text-xs text-[var(--muted)]">{linkPreview.url}</p>
            </div>
          ) : null}
          {mediaUrls.length ? (
            <ul className="mt-3 space-y-2">
              {mediaUrls.map((url) => (
                <li key={url} className="rounded-lg border border-[var(--line)] px-3 py-2 text-xs text-[var(--muted)]">{url}</li>
              ))}
            </ul>
          ) : null}
        </SectionCard>

        <SectionCard title="Publishing" subtitle="Post now, schedule, queue, or draft">
          <input
            type="datetime-local"
            value={scheduleAt}
            onChange={(e) => setScheduleAt(e.target.value)}
            className="w-full rounded-xl border border-[var(--line)] bg-[var(--surface)] px-4 py-3 text-sm outline-none focus:border-[var(--primary)]"
          />
          <div className="mt-3 grid gap-2">
            <button onClick={createByStatus("draft")} className="rounded-xl border border-[var(--line)] px-4 py-2.5 text-sm">Save Draft</button>
            <button onClick={createByStatus("published")} className="rounded-xl border border-[var(--line)] px-4 py-2.5 text-sm">Post Now</button>
            <button onClick={createByStatus("scheduled")} className="rounded-xl border border-[var(--line)] px-4 py-2.5 text-sm">Schedule</button>
            <button onClick={createByStatus("queued")} className="rounded-xl bg-[var(--primary)] px-4 py-2.5 text-sm font-semibold text-[var(--bg)]">Queue Post</button>
          </div>
          {createMutation.isError ? <p className="mt-3 text-sm text-red-500">Unable to save post.</p> : null}
          {createMutation.isSuccess ? <p className="mt-3 text-sm text-emerald-500">Post saved successfully.</p> : null}
        </SectionCard>
      </motion.aside>
    </div>
  );
}
