import { motion } from "framer-motion";

type MetricCardProps = {
  label: string;
  value: string | number;
  note: string;
};

export function MetricCard({ label, value, note }: MetricCardProps) {
  return (
    <motion.article
      initial={{ opacity: 0, y: 16 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.35, ease: "easeOut" }}
      className="rounded-2xl border border-[var(--line)] bg-[var(--bg)] p-5"
    >
      <p className="text-xs uppercase tracking-[0.18em] text-[var(--muted)]">{label}</p>
      <p className="mt-3 text-3xl font-semibold">{value}</p>
      <p className="mt-2 text-sm text-[var(--muted)]">{note}</p>
    </motion.article>
  );
}
