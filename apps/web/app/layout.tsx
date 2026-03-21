import type { Metadata } from "next";
import { Space_Grotesk, Leckerli_One } from "next/font/google";
import "./globals.css";
import ClientLayoutShell from "@/components/layout/client-layout-shell";

export const leckerliOne = Leckerli_One({
  subsets: ["latin"],
  weight: "400",
  variable: "--font-leckerli-one", 
});

const spaceGrotesk = Space_Grotesk({
  variable: "--font-space-grotesk",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "PostPilot - Social scheduling command center",
  description: "Social scheduling command center",

  // Favicon + Icons
  icons: {
    icon: [
      { url: "./favicon/favicon.ico" }, // fallback (important)
      { url: "./favicon/favicon-16x16.png", sizes: "16x16", type: "image/png" },
      { url: "./favicon/favicon-32x32.png", sizes: "32x32", type: "image/png" },
    ],
    apple: [{ url: "./favicon/apple-touch-icon.png", sizes: "180x180" }],
    shortcut: "./favicon/favicon.ico",
  },

  // PWA support
  manifest: "./favicon/site.webmanifest",
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

      <body
        className={`${spaceGrotesk.variable} ${leckerliOne.variable} antialiased`}
      >
        <ClientLayoutShell>{children}</ClientLayoutShell>
      </body>
    </html>
  );
}