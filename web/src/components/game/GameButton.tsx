import { ButtonHTMLAttributes, ReactNode } from "react";

type Variant = "primary" | "secondary" | "gold";

const styles: Record<Variant, { bg: string; shadow: string }> = {
  primary: { bg: "bg-[#58CC02] hover:bg-[#61D802]", shadow: "#46A302" },
  secondary: { bg: "bg-[#1CB0F6] hover:bg-[#2EB8FF]", shadow: "#1899D6" },
  gold: { bg: "bg-[#FFC800] hover:bg-[#FFD633]", shadow: "#E5B000" },
};

/** GameButton is a large pill CTA — never a flat rectangle. */
export function GameButton({
  children,
  variant = "primary",
  className = "",
  ...props
}: ButtonHTMLAttributes<HTMLButtonElement> & {
  children: ReactNode;
  variant?: Variant;
}) {
  const s = styles[variant];
  return (
    <button
      className={`game-btn font-display w-full rounded-full px-8 py-4 text-lg font-bold uppercase tracking-wide text-white ${s.bg} ${className}`}
      style={{ ["--btn-shadow" as string]: s.shadow }}
      {...props}
    >
      {children}
    </button>
  );
}
