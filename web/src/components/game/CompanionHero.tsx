import { companionVisual } from "@/lib/design/companions";

type Reaction = "idle" | "correct" | "wrong";

/** CompanionHero shows the animated companion with optional reaction state. */
export function CompanionHero({
  name,
  species,
  size = "lg",
  reaction = "idle",
}: {
  name?: string;
  species?: string;
  size?: "sm" | "md" | "lg";
  reaction?: Reaction;
}) {
  const visual = companionVisual(name, species);

  const sizeClass =
    size === "lg" ? "h-28 w-28 text-5xl" :
    size === "md" ? "h-20 w-20 text-4xl" :
    "h-12 w-12 text-2xl";

  const animClass =
    reaction === "correct" ? "animate-companion-correct" :
    reaction === "wrong" ? "animate-companion-wrong" :
    "animate-companion-idle";

  return (
    <div className="relative inline-flex flex-col items-center">
      <div
        className={`absolute inset-0 rounded-full blur-md opacity-40 ${sizeClass}`}
        style={{ background: visual.glow, transform: "scale(1.2)" }}
        aria-hidden
      />
      <div
        className={`relative flex items-center justify-center rounded-full border border-white/10 ${sizeClass} ${animClass}`}
        style={{
          background: `linear-gradient(135deg, ${visual.glow}33 0%, #1A1D27 100%)`,
          boxShadow: `0 0 24px ${visual.glow}44, inset 0 1px 0 rgba(255,255,255,0.08)`,
        }}
      >
        <span role="img" aria-label={name ?? "companion"} className="select-none">
          {visual.emoji}
        </span>
      </div>
    </div>
  );
}
