import { ReactNode } from "react";

type Variant = "default" | "forest" | "challenge";

const variants: Record<Variant, string> = {
  default: "game-bg",
  forest: "game-bg-forest",
  challenge: "game-bg-challenge",
};

/** GameBackground wraps screens with gradient game worlds — never plain white. */
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
    <div className={`${variants[variant]} ${className}`}>
      <div className="pointer-events-none fixed inset-0 overflow-hidden opacity-30">
        <div className="absolute -left-10 top-20 h-32 w-32 rounded-full bg-lime-300 blur-3xl" />
        <div className="absolute right-0 top-40 h-40 w-40 rounded-full bg-sky-300 blur-3xl" />
        <div className="absolute bottom-20 left-1/3 h-24 w-24 rounded-full bg-amber-200 blur-2xl" />
      </div>
      <div className="relative z-10">{children}</div>
    </div>
  );
}
