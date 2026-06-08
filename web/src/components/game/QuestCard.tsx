/** QuestCard displays a daily quest with progress and reward preview. */
export function QuestCard({
  title,
  icon,
  progress,
  target,
  completed,
  rewardXp,
  rewardGems,
  comingSoon = false,
}: {
  title: string;
  icon: string;
  progress: number;
  target: number;
  completed: boolean;
  rewardXp: number;
  rewardGems: number;
  comingSoon?: boolean;
}) {
  const pct = Math.min(100, (progress / target) * 100);

  if (comingSoon) {
    return (
      <div
        className="rounded-2xl p-4 opacity-60"
        style={{
          background: "#1A1D27",
          border: "1px dashed #2E3347",
        }}
      >
        <div className="flex items-start gap-3">
          <span className="mt-0.5 text-xl opacity-50">{icon}</span>
          <div className="flex-1 min-w-0">
            <p className="font-display text-sm font-bold" style={{ color: "#4A5068" }}>
              {title}
            </p>
            <p className="font-mono mt-1 text-xs" style={{ color: "#4A5068" }}>
              Coming soon
            </p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div
      className="rounded-2xl p-4 transition-all"
      style={{
        background: completed ? "rgba(52,211,153,0.08)" : "#1A1D27",
        border: `1px solid ${completed ? "rgba(52,211,153,0.3)" : "#2E3347"}`,
      }}
    >
      <div className="flex items-start gap-3">
        <span className="mt-0.5 text-xl">{completed ? "✅" : icon}</span>
        <div className="flex-1 min-w-0">
          <p
            className="font-display text-sm font-bold"
            style={{
              color: completed ? "#34D399" : "#E8EAED",
              textDecoration: completed ? "line-through" : "none",
              opacity: completed ? 0.7 : 1,
            }}
          >
            {title}
          </p>
          <div
            className="mt-2 h-1.5 overflow-hidden rounded-full"
            style={{ background: "#2E3347" }}
          >
            <div
              className="h-full rounded-full transition-all duration-700"
              style={{
                width: `${pct}%`,
                background: completed
                  ? "#34D399"
                  : "linear-gradient(90deg, #7C6EF5, #9D8FF7)",
              }}
            />
          </div>
          <p className="font-mono mt-1.5 text-xs" style={{ color: "#4A5068" }}>
            {progress}/{target} · ⚡ {rewardXp} XP · 💎 {rewardGems}
          </p>
        </div>
      </div>
    </div>
  );
}
