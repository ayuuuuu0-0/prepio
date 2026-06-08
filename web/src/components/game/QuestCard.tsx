/** QuestCard displays a daily quest with progress and reward preview. */
export function QuestCard({
  title,
  icon,
  progress,
  target,
  completed,
  rewardXp,
  rewardGems,
}: {
  title: string;
  icon: string;
  progress: number;
  target: number;
  completed: boolean;
  rewardXp: number;
  rewardGems: number;
}) {
  const pct = Math.min(100, (progress / target) * 100);

  return (
    <div
      className={`rounded-3xl p-4 shadow-md transition ${
        completed ? "bg-lime-100 ring-2 ring-[#58CC02]" : "bg-white"
      }`}
    >
      <div className="flex items-start gap-3">
        <span className="text-2xl">{completed ? "✅" : icon}</span>
        <div className="flex-1">
          <p className={`font-display font-bold ${completed ? "text-lime-700 line-through" : "text-[#3C3C3C]"}`}>
            {title}
          </p>
          <div className="mt-2 h-3 overflow-hidden rounded-full bg-gray-200">
            <div
              className="h-full rounded-full bg-gradient-to-r from-[#58CC02] to-[#7CB342] transition-all duration-700"
              style={{ width: `${pct}%` }}
            />
          </div>
          <p className="mt-1 text-xs text-[#777]">
            {progress}/{target} · ⚡{rewardXp} XP · 💎{rewardGems} gems
          </p>
        </div>
      </div>
    </div>
  );
}
