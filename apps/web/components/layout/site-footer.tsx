import Link from "next/link";

const footerLinks = [
  { href: "/", label: "Product" },
  { href: "/signin", label: "Get Started" },
  { href: "/login", label: "Login" },
];

export function SiteFooter() {
  return (
    <footer className="border-t border-[var(--line)] bg-[var(--bg)]">
      <div className="mx-auto flex w-full max-w-6xl flex-col gap-6 px-6 py-10 md:flex-row md:items-center md:justify-between">
        <div>
          <p className="text-sm font-semibold uppercase tracking-[0.2em] text-[var(--primary)]">PostPilot</p>
          <p className="mt-2 max-w-xl text-sm text-[var(--muted)]">
            Launch social content faster with reliable scheduling, thoughtful automation, and a cleaner workflow.
          </p>
        </div>

        <div className="flex gap-2">
          {footerLinks.map((item) => (
            <Link
              key={item.label}
              href={item.href}
              className="rounded-full px-4 py-2 text-sm font-medium text-[var(--primary)] transition hover:bg-[var(--primary)]/10"
            >
              {item.label}
            </Link>
          ))}
        </div>
      </div>
    </footer>
  );
}
