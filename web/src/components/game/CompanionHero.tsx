import { companionVisual } from "@/lib/design/companions";

/** CompanionHero shows a large animated companion — mandatory on major screens. */
export function CompanionHero({
  name,
  species,
  size = "lg",
}: {
  name?: string;
  species?: string;
  size?: "sm" | "md" | "lg";
}) {
  const visual = companionVisual(name, species);
  const sizeClass = size === "lg" ? "h-28 w-28 text-6xl" : size === "md" ? "h-20 w-20 text-4xl" : "h-12 w-12 text-2xl";

  return (
    <div
      className={`animate-companion-idle relative flex items-center justify-center rounded-full bg-gradient-to-br ${visual.gradient} shadow-lg ${sizeClass}`}
      style={{ boxShadow: `0 8px 24px ${visual.glow}55` }}
    >
      <span role="img" aria-label={name ?? "companion"}>
        {visual.emoji}
      </span>
    </div>
  );
}
