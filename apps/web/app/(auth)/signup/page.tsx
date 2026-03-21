"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { FormEvent, useState } from "react";
import { motion } from "framer-motion";
import { AuthService } from "@/services/auth-service";

const ease = [0.25, 0.46, 0.45, 0.94] as const;

const container = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: { staggerChildren: 0.09, delayChildren: 0.05 },
  },
};

const item = {
  hidden: { opacity: 0, y: 18 },
  visible: { opacity: 1, y: 0, transition: { duration: 0.4, ease } },
};

export default function SignupPage() {
  const router = useRouter();
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      await AuthService.register({ name, email, password });
      router.push("/onboarding");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unable to create account");
    } finally {
      setLoading(false);
    }
  }

  return (
    <motion.div variants={container} initial="hidden" animate="visible">

      {/* Heading */}
      <motion.div variants={item}>
        <p className="text-xs font-semibold uppercase tracking-[0.22em] text-[var(--muted)]">
          Get started free
        </p>
        <h1 className="mt-2 text-3xl font-semibold tracking-tight">
          Create your account
        </h1>
        <p className="mt-2 text-sm text-[var(--muted)]">
          Start scheduling social posts in minutes.
        </p>
      </motion.div>

      <form onSubmit={handleSubmit} className="mt-8 space-y-4">

        <motion.label variants={item} className="block text-sm">
          <span className="mb-1.5 block font-medium">Full name</span>
          <input
            type="text"
            required
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="Jane Doe"
            className="w-full rounded-xl border border-[var(--line)] bg-[var(--bg)] px-4 py-3 text-sm outline-none transition focus:border-[var(--primary)] focus:ring-2 focus:ring-[var(--primary)]/15"
          />
        </motion.label>

        <motion.label variants={item} className="block text-sm">
          <span className="mb-1.5 block font-medium">Email</span>
          <input
            type="email"
            required
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="you@company.com"
            className="w-full rounded-xl border border-[var(--line)] bg-[var(--bg)] px-4 py-3 text-sm outline-none transition focus:border-[var(--primary)] focus:ring-2 focus:ring-[var(--primary)]/15"
          />
        </motion.label>

        <motion.label variants={item} className="block text-sm">
          <span className="mb-1.5 block font-medium">Password</span>
          <input
            type="password"
            required
            minLength={8}
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            placeholder="At least 8 characters"
            className="w-full rounded-xl border border-[var(--line)] bg-[var(--bg)] px-4 py-3 text-sm outline-none transition focus:border-[var(--primary)] focus:ring-2 focus:ring-[var(--primary)]/15"
          />
        </motion.label>

        {error && (
          <motion.p
            initial={{ opacity: 0, scale: 0.97 }}
            animate={{ opacity: 1, scale: 1 }}
            className="rounded-xl border border-red-400/40 bg-red-500/10 px-3 py-2 text-sm text-red-500"
          >
            {error}
          </motion.p>
        )}

        <motion.div variants={item}>
          <motion.button
            type="submit"
            disabled={loading}
            whileHover={{ y: -2 }}
            whileTap={{ scale: 0.98 }}
            className="w-full rounded-full bg-[var(--primary)] px-5 py-3 text-sm font-semibold uppercase tracking-[0.18em] text-[var(--bg)] transition disabled:opacity-50 cursor-pointer"
          >
            {loading ? "Creating Account…" : "Create Account"}
          </motion.button>
        </motion.div>
      </form>

      <motion.div variants={item} className="mt-5">
        <div className="mb-3 flex items-center gap-3 text-xs text-[var(--muted)]">
          <span className="h-px flex-1 bg-[var(--line)]" />
          Continue with OAuth
          <span className="h-px flex-1 bg-[var(--line)]" />
        </div>
        <div className="grid grid-cols-2 gap-2 text-xs">
          {[
            "Google",
            "Twitter/X",
            "Facebook",
            "LinkedIn",
          ].map((provider) => (
            <button
              key={provider}
              type="button"
              className="rounded-xl border border-[var(--line)] bg-[var(--surface)] px-3 py-2 text-left transition hover:bg-[var(--line)]/30"
            >
              {provider}
            </button>
          ))}
        </div>
      </motion.div>

      <motion.p variants={item} className="mt-6 text-sm text-[var(--muted)]">
        Already have an account?{" "}
        <Link
          href="/login"
          className="font-semibold text-[var(--primary)] underline underline-offset-4 transition hover:opacity-75"
        >
          Log in
        </Link>
      </motion.p>

    </motion.div>
  );
}
