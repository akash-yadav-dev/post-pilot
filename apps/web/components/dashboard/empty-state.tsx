type EmptyStateProps = {
  title: string;
  message: string;
};

export function EmptyState({ title, message }: EmptyStateProps) {
  return (
    <div className="rounded-2xl border border-dashed border-[var(--line)] bg-[var(--surface)] p-8 text-center">
      <p className="text-lg font-semibold tracking-tight">{title}</p>
      <p className="mx-auto mt-2 max-w-xl text-sm text-[var(--muted)]">{message}</p>
    </div>
  );
}
