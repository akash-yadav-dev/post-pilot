"use client";

import { ThemeProvider } from "@/components/providers/theme-provider";
import { QueryProvider } from "@/components/providers/query-provider";
import { SiteShell } from "@/components/layout/site-shell";
import GlobalLoader from "@/components/ui/global-loader";
import { usePageLoader } from "@/hooks/usePageLoader";

type ClientLayoutShellProps = {
  children: React.ReactNode;
};

export default function ClientLayoutShell({ children }: ClientLayoutShellProps) {
  const loading = usePageLoader();

  return (
    <QueryProvider>
      <ThemeProvider>
        <div className="min-h-screen bg-[var(--bg)] text-[var(--primary)]">
          <SiteShell>
            {children}
            {loading && <GlobalLoader />}
          </SiteShell>
        </div>
      </ThemeProvider>
    </QueryProvider>
  );
}
