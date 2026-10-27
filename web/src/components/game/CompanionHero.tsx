"use client";

import { companionVisual } from "@/lib/design/companions";
import { FullBodyCompanion } from "./companions/FullBodyCompanion";

type Reaction = "idle" | "correct" | "wrong";

/** CompanionHero shows the full-body dressed mascot with animated gestures and reactions. */
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

  const containerSize =
    size === "lg" ? "h-36 w-32" :
    size === "md" ? "h-28 w-24" :
    "h-16 w-14";

  const animClass =
    reaction === "correct" ? "animate-companion-correct" :
    reaction === "wrong" ? "animate-companion-wrong" :
    "animate-companion-idle";

  return (
    <div className={`relative inline-flex flex-col items-center justify-end ${containerSize}`}>
      {/* Background Soft Glow */}
      <div
        className="absolute bottom-1 h-14 w-24 rounded-full blur-xl opacity-50"
        style={{ background: visual.glow }}
        aria-hidden
      />
      {/* Full-Body Dressed Mascot Vector */}
      <div className={`relative h-full w-full flex items-center justify-center ${animClass}`}>
        <FullBodyCompanion
          name={name}
          species={species}
          reaction={reaction}
          className="h-full w-full object-contain filter drop-shadow-[0_8px_16px_rgba(0,0,0,0.4)]"
        />
      </div>
    </div>
  );
}
