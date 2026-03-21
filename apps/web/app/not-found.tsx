import Link from "next/link";

export default function NotFoundPage() {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-[var(--bg)] text-[var(--primary)] px-6">
      <div className="w-28 h-28 rounded-2xl border border-white/20 flex items-center justify-center mb-6 rotate-45 bg-[var(--icon-bg)]">
        <span className="text-4xl font-bold -rotate-45 [font-family:var(--font-leckerli-one)] text-[var(--icon-text)]">
          P
        </span>
      </div>

      <h1 className="text-5xl font-bold mb-4">404</h1>
      <p className="text-gray-400 mb-6 text-center">
        Page not found. It might have been moved or deleted.
      </p>

      <Link
        href="/"
        className="rounded-full px-4 py-2 text-sm font-medium transition bg-[var(--primary)] text-[var(--bg)] hover:bg-[var(--primary)]/90"
      >
        Go Home
      </Link>
    </div>
  );
}
