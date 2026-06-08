/** StatChip displays a single game metric with monospace value. */
export function StatChip({
  icon,
  label,
  value,
  color,
}: {
  icon: string;
  label: string;
  value: string | number;
  color: string;
}) {
  return (
    <div
      className="flex items-center gap-2.5 rounded-xl px-3 py-2"
      style={{
        background: `${color}14`,
        border: `1px solid ${color}30`,
      }}
    >
      <span className="text-base leading-none">{icon}</span>
      <div>
        <p
          className="font-display text-[10px] font-bold uppercase tracking-widest"
          style={{ color: `${color}99` }}
        >
          {label}
        </p>
        <p className="font-mono text-sm font-bold" style={{ color }}>
          {value}
        </p>
      </div>
    </div>
  );
}
