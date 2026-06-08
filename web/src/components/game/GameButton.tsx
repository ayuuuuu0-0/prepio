import { ButtonHTMLAttributes, ReactNode } from "react";

type Variant = "primary" | "secondary" | "gold" | "ghost";

interface BtnStyle {
  bg: string;
  shadow: string;
  text: string;
}

const styles: Record<Variant, BtnStyle> = {
  primary: {
    bg: "bg-gradient-to-r from-[#7C6EF5] to-[#9D8FF7] hover:from-[#8B7FF7] hover:to-[#AFA0FF]",
    shadow: "#5B50D4",
    text: "text-white",
  },
  secondary: {
    bg: "bg-[#242836] hover:bg-[#2E3347] border border-[#2E3347] hover:border-[#7C6EF5]/50",
    shadow: "#1A1D27",
    text: "text-[#C8CCDA]",
  },
  gold: {
    bg: "bg-gradient-to-r from-[#F5B942] to-[#FFCF60] hover:from-[#FFCF60] hover:to-[#FFD97A]",
    shadow: "#C8941A",
    text: "text-[#1A1500]",
  },
  ghost: {
    bg: "bg-transparent hover:bg-white/5 border border-white/10 hover:border-white/20",
    shadow: "transparent",
    text: "text-[#8B92A8] hover:text-[#C8CCDA]",
  },
};

/** GameButton is the primary pill-shaped CTA for game screens. */
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
      className={`game-btn font-display w-full rounded-full px-8 py-4 text-base font-bold tracking-wide ${s.bg} ${s.text} ${className}`}
      style={{ ["--btn-shadow" as string]: s.shadow }}
      {...props}
    >
      {children}
    </button>
  );
}
