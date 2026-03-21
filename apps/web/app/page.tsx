import Link from "next/link";

export default function Home() {
  return (
    <div className="mx-auto flex w-full max-w-6xl flex-col gap-16 px-6 pb-20 pt-12">
      <section className="grid items-center gap-8 rounded-3xl border border-[var(--line)] bg-[var(--surface)] p-8 shadow-[0_25px_80px_-35px_rgba(0,0,0,0.45)] md:grid-cols-2">
        <div>
          <p className="text-xs font-semibold uppercase tracking-[0.26em] text-[var(--muted)]">Phase 1 Product</p>
          <h1 className="mt-4 text-4xl font-semibold tracking-tight md:text-5xl">Schedule once. Publish everywhere.</h1>
          <p className="mt-5 max-w-xl text-base leading-7 text-[var(--muted)] md:text-lg">
            PostPilot gives creators and teams a sharp command center for content planning, scheduling, and cross-platform publishing.
          </p>
          <div className="mt-8 flex flex-wrap gap-3">
            <Link
              href="/signup"
              className="rounded-full bg-[var(--primary)] px-6 py-3 text-sm font-semibold uppercase tracking-[0.16em] text-[var(--bg)] transition hover:-translate-y-0.5"
            >
              Start Free
            </Link>
            <Link
              href="/login"
              className="rounded-full border border-[var(--line)] px-6 py-3 text-sm font-semibold uppercase tracking-[0.16em] text-[var(--primary)] transition hover:bg-[var(--primary)] hover:text-[var(--bg)]"
            >
              Log In
            </Link>
          </div>
        </div>

        <div className="grid gap-4">
          {["Queue Builder", "Smart Retry", "Analytics Signals"].map((item) => (
            <article key={item} className="rounded-2xl border border-[var(--line)] bg-[var(--bg)] p-5">
              <h3 className="text-lg font-semibold">{item}</h3>
              <p className="mt-2 text-sm text-[var(--muted)]">
                Built for high signal workflows, with production-grade reliability from draft to publish.
              </p>
            </article>
          ))}
        </div>
      </section>

      <section>
        <h2 className="text-2xl font-semibold tracking-tight md:text-3xl">Features for modern social teams</h2>
        <div className="mt-6 grid gap-4 md:grid-cols-3">
          {[
            "Multi-platform scheduling",
            "Approval workflow",
            "Failure and retry visibility",
            "Timezone aware planner",
            "Post performance snapshots",
            "Secure token management",
          ].map((feature) => (
            <div key={feature} className="rounded-2xl border border-[var(--line)] bg-[var(--surface)] p-5 text-sm font-medium">
              {feature}
            </div>
          ))}
        </div>
      </section>

      <section>
        <h2 className="text-2xl font-semibold tracking-tight md:text-3xl">What early users say</h2>
        <div className="mt-6 grid gap-4 md:grid-cols-3">
          {[
            {
              quote: "We moved from spreadsheets to a real publishing system in one week.",
              author: "Sanaa, Growth Lead",
            },
            {
              quote: "The scheduling and retry flow feels safe enough for production campaigns.",
              author: "Mark, Marketing Ops",
            },
            {
              quote: "A clean dashboard and exactly the controls we needed for phase one.",
              author: "Elina, Founder",
            },
          ].map((item) => (
            <figure key={item.author} className="rounded-2xl border border-[var(--line)] bg-[var(--surface)] p-6">
              <blockquote className="text-sm leading-6 text-[var(--muted)]">&ldquo;{item.quote}&rdquo;</blockquote>
              <figcaption className="mt-4 text-sm font-semibold">{item.author}</figcaption>
            </figure>
          ))}
        </div>
      </section>
    </div>
  );
}
