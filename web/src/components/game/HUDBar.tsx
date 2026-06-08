"use client";

import { DashboardHome } from "@/lib/api";
import { LEVEL_THRESHOLDS } from "@/lib/design/tokens";

/** HUDBar combines companion badge, streak, gems, and XP progress in one row. */
export function HUDBar({ home }: { home: DashboardHome }) {
  const { streak, progress, companion } = home;
  const levelStart = LEVEL_THRESHOLDS[progress.current_level - 1] ?? 0;
  const levelEnd = LEVEL_THRESHOLDS[progress.current_level] ?? levelStart + 1000;
  const xpInLevel = progress.total_xp - levelStart;
  const xpForLevel = levelEnd - levelStart;
  const xpPct = xpForLevel > 0 ? Math.min(100, Math.round((xpInLevel / xpForLevel) * 100)) : 100;

  return (
    <div
      className="flex items-center gap-3 rounded-2xl px-4 py-3"
      style={{
        background: "#1A1D27",
        border: "1px solid #2E3347",
      }}
    >
      <div className="flex shrink-0 items-center gap-2">
        <span className="text-xl">
          {companion?.name ? getCompanionEmoji(companion.species) : "🦫"}
        </span>
        <div>
          <p className="font-mono text-[10px] font-bold uppercase tracking-widest" style={{ color: "#7C6EF5" }}>
            {companion?.name ?? "Byte"}
          </p>
          <p className="font-mono text-xs font-bold" style={{ color: "#60A5FA" }}>
            Lv.{progress.current_level}
          </p>
        </div>
      </div>

      <div className="h-8 w-px shrink-0" style={{ background: "#2E3347" }} />

      <div className="flex shrink-0 items-center gap-1.5">
        <span className="text-sm">🔥</span>
        <span className="font-mono text-sm font-bold" style={{ color: "#FF6B35" }}>
          {streak.current_streak}
        </span>
      </div>

      <div className="flex shrink-0 items-center gap-1.5">
        <span className="text-sm">💎</span>
        <span className="font-mono text-sm font-bold" style={{ color: "#34D399" }}>
          {progress.gem_balance}
        </span>
      </div>

      <div className="h-8 w-px shrink-0" style={{ background: "#2E3347" }} />

      <div className="min-w-0 flex-1">
        <div className="mb-1 flex items-center justify-between">
          <span className="font-mono text-[10px] font-semibold" style={{ color: "#4A5068" }}>
            XP
          </span>
          <span className="font-mono text-[10px] font-semibold" style={{ color: "#60A5FA" }}>
            {progress.xp_to_next_level} to next
          </span>
        </div>
        <div className="h-1.5 w-full overflow-hidden rounded-full" style={{ background: "#2E3347" }}>
          <div
            className="animate-xp-bar h-full rounded-full"
            style={{
              width: `${xpPct}%`,
              background: "linear-gradient(90deg, #60A5FA, #7C6EF5)",
            }}
          />
        </div>
      </div>
    </div>
  );
}

function getCompanionEmoji(species?: string): string {
  const map: Record<string, string> = {
    capybara: "🦫",
    red_panda: "🐼",
    pangolin: "🦔",
    axolotl: "🦎",
    snow_leopard: "🐆",
    owl: "🦉",
  };
  return map[species ?? ""] ?? "🦫";
}
