import { ReactNode } from "react";

/** GameCard is a dark surface card with optional left accent border. */
export function GameCard({
  children,
  icon,
  className = "",
  accentColor,
}: {
  children: ReactNode;
  icon?: string;
  className?: string;
  accentColor?: string;
}) {
  return (
    <div
      className={`overflow-hidden rounded-2xl p-5 ${className}`}
      style={{
        background: "#1A1D27",
        border: "1px solid #2E3347",
        boxShadow: "0 4px 24px rgba(0,0,0,0.3)",
        borderLeft: accentColor ? `3px solid ${accentColor}` : "1px solid #2E3347",
      }}
    >
      {icon && <div className="mb-3 text-2xl">{icon}</div>}
      {children}
    </div>
  );
}
