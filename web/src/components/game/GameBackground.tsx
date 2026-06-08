import { ReactNode } from "react";

type Variant = "default" | "forest" | "challenge";

const variants: Record<Variant, string> = {
  default: "game-bg",
  forest: "game-bg-forest",
  challenge: "game-bg-challenge",
};

/** GameBackground wraps screens with ambient dark game backgrounds. */
export function GameBackground({
  children,
  variant = "default",
  className = "",
}: {
  children: ReactNode;
  variant?: Variant;
  className?: string;
}) {
  return (
    <div className={`${variants[variant]} relative ${className}`}>
      <div className="pointer-events-none fixed inset-0 overflow-hidden" aria-hidden>
        <div
          className="absolute -left-20 -top-20 h-64 w-64 rounded-full opacity-20 blur-3xl"
          style={{ background: "radial-gradient(circle, #7C6EF5 0%, transparent 70%)" }}
        />
        <div
          className="absolute -right-20 bottom-1/3 h-48 w-48 rounded-full opacity-10 blur-3xl"
          style={{ background: "radial-gradient(circle, #34D399 0%, transparent 70%)" }}
        />
      </div>
      <div className="relative z-10">{children}</div>
    </div>
  );
}
