import { LEVEL_THRESHOLDS } from "@/lib/design/tokens";

/** XpProgressBar shows progress toward the next level. */
export function XpProgressBar({
  level,
  totalXp,
  xpToNextLevel,
}: {
  level: number;
  totalXp: number;
  xpToNextLevel: number;
}) {
  const levelStart = LEVEL_THRESHOLDS[level - 1] ?? 0;
  const levelEnd = LEVEL_THRESHOLDS[level] ?? levelStart + 1000;
  const xpInLevel = totalXp - levelStart;
  const xpForLevel = levelEnd - levelStart;
  const pct = xpForLevel > 0 ? Math.min(100, (xpInLevel / xpForLevel) * 100) : 100;

  return (
    <div
      className="rounded-2xl p-4"
      style={{ background: "#1A1D27", border: "1px solid #2E3347" }}
    >
      <div className="flex items-center justify-between">
        <span className="font-mono text-sm font-bold" style={{ color: "#60A5FA" }}>
          Level {level}
        </span>
        <span className="font-mono text-sm font-bold" style={{ color: "#7C6EF5" }}>
          Level {level + 1}
        </span>
      </div>
      <div className="mt-2 h-2 overflow-hidden rounded-full" style={{ background: "#2E3347" }}>
        <div
          className="h-full rounded-full transition-all duration-700"
          style={{
            width: `${pct}%`,
            background: "linear-gradient(90deg, #60A5FA, #7C6EF5)",
          }}
        />
      </div>
      <p className="mt-2 text-center font-mono text-xs font-semibold" style={{ color: "#4A5068" }}>
        {xpToNextLevel > 0 ? `${xpToNextLevel} XP to next level` : "Max level reached"}
      </p>
    </div>
  );
}
