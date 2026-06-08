import { ReactNode } from "react";

/** GameCard is a colorful rounded card with illustration slot — never plain text-only. */
export function GameCard({
  children,
  gradient,
  icon,
  className = "",
}: {
  children: ReactNode;
  gradient?: string;
  icon?: string;
  className?: string;
}) {
  return (
    <div
      className={`overflow-hidden rounded-3xl p-5 shadow-lg ${gradient ?? "bg-white"} ${className}`}
    >
      {icon && <div className="mb-3 text-3xl">{icon}</div>}
      {children}
    </div>
  );
}
