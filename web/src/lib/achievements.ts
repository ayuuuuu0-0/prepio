/** Achievement visuals. The catalog and who earned what come from the server; only icons live here. */

const icons: Record<string, string> = {
  "first-steps": "🌱",
  "warming-up": "🔥",
  "ten-lessons": "📚",
  "full-week": "🗓️",
  "level-five": "⭐",
  "halfway-there": "🎯",
  "moving-up": "🏆",
  "path-finder": "🧭",
};

/** achievementIcon returns the icon for an achievement, with a neutral badge for new ones. */
export function achievementIcon(slug: string): string {
  return icons[slug] ?? "🏅";
}

/** earnedLabel says when an achievement was earned, in the viewer's locale. */
export function earnedLabel(unlockedAt: string | undefined, locale?: string): string {
  if (!unlockedAt) return "";
  const d = new Date(unlockedAt);
  if (Number.isNaN(d.getTime())) return "";
  return "Earned " + d.toLocaleDateString(locale, { day: "numeric", month: "short", year: "numeric" });
}
