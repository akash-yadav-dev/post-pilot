"use client";

import { useState } from "react";
import type { SocialPlatform } from "@/lib/social/types";
import { SectionCard } from "@/components/dashboard/section-card";

const platforms: SocialPlatform[] = ["twitter", "linkedin", "facebook", "instagram", "pinterest"];

export default function OnboardingPage() {
  const [selected, setSelected] = useState<SocialPlatform[]>(["twitter", "linkedin"]);

  return (
    <div className="space-y-5 p-4 md:p-6">
      <SectionCard title="Account Onboarding" subtitle="Connect accounts and choose your publishing channels">
        <div className="grid gap-3 md:grid-cols-2">
          {["Google", "Twitter/X", "Facebook", "LinkedIn"].map((provider) => (
            <button key={provider} className="rounded-xl border border-[var(--line)] bg-[var(--surface)] px-4 py-3 text-sm text-left">
              Continue with {provider}
            </button>
          ))}
        </div>

        <p className="mt-5 text-sm font-medium">Select platforms</p>
        <div className="mt-2 flex flex-wrap gap-2">
          {platforms.map((platform) => {
            const active = selected.includes(platform);
            return (
              <button
                key={platform}
                onClick={() => setSelected((prev) => prev.includes(platform) ? prev.filter((item) => item !== platform) : [...prev, platform])}
                className={`rounded-full border px-3 py-1.5 text-xs capitalize ${active ? "border-[var(--primary)] bg-[var(--primary)] text-[var(--bg)]" : "border-[var(--line)]"}`}
              >
                {platform}
              </button>
            );
          })}
        </div>
      </SectionCard>
    </div>
  );
}
