/** StatChip shows XP, gems, streak as colorful game chips. */
export function StatChip({ icon, label, value, color }: { icon: string; label: string; value: string | number; color: string }) {
  return (
    <div
      className="flex items-center gap-2 rounded-full px-4 py-2 shadow-md"
      style={{ backgroundColor: `${color}22`, border: `2px solid ${color}` }}
    >
      <span className="text-lg">{icon}</span>
      <div>
        <p className="font-display text-xs font-bold uppercase" style={{ color }}>
          {label}
        </p>
        <p className="font-display text-sm font-extrabold text-[#3C3C3C]">{value}</p>
      </div>
    </div>
  );
}
