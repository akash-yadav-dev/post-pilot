import type { Metadata } from "next";
import { Space_Grotesk } from "next/font/google";
import "./globals.css";
import { ThemeProvider } from "@/components/providers/theme-provider";
import { SiteHeader } from "@/components/layout/site-header";
import { SiteFooter } from "@/components/layout/site-footer";

const spaceGrotesk = Space_Grotesk({
  variable: "--font-space-grotesk",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "PostPilot",
  description: "Social scheduling command center",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
  }) {
  const initialThemeScript = `(function() {
        try {
          const key = "postpilot-theme";
          const stored = localStorage.getItem(key);
          const systemDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
          const theme = stored || (systemDark ? "dark" : "light");

          if (theme === "dark") {
            document.documentElement.classList.add("dark");
          } else {
            document.documentElement.classList.remove("dark");
          }
        } catch (e) {}
      })()`;
  return (
    <html lang="en" suppressHydrationWarning>
      <head>
        {/* ✅ Prevent hydration mismatch + theme flicker */}
        <script
          dangerouslySetInnerHTML={{
            __html: initialThemeScript,
          }}
        />
      </head>

      <body className={`${spaceGrotesk.variable} antialiased`}>
        <ThemeProvider>
          <div className="min-h-screen bg-[var(--bg)] text-[var(--primary)]">
            <SiteHeader />
            <main>{children}</main>
            <SiteFooter />
          </div>
        </ThemeProvider>
      </body>
    </html>
  );
}