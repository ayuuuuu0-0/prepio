import { tierTint } from "@/lib/leagues";
import type { LeagueTier } from "@/lib/api";

/**
 * TierEmblem draws a league tier as a faceted shield in the tier's tint. Tiers above the
 * learner's are shown dimmed, so the ladder reads as "where you could go".
 */
export function TierEmblem({
  tier,
  size = 56,
  state = "current",
}: {
  tier: LeagueTier;
  size?: number;
  state?: "current" | "reached" | "ahead";
}) {
  const tint = tierTint(tier);
  const dim = state === "ahead";
  const id = `tier-${tier.slug}-${size}`;
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 64 64"
      role="img"
      aria-label={`${tier.name} league${dim ? " (ahead)" : ""}`}
      className="shrink-0"
      style={{
        opacity: dim ? 0.32 : 1,
        filter: state === "current" ? `drop-shadow(0 0 14px ${tint}88)` : undefined,
      }}
    >
      <defs>
        <linearGradient id={id} x1="0" y1="0" x2="1" y2="1">
          <stop offset="0%" stopColor={tint} />
          <stop offset="100%" stopColor={tint} stopOpacity="0.55" />
        </linearGradient>
      </defs>
      <path d="M32 4 L56 13 V31 C56 45 45 55 32 60 C19 55 8 45 8 31 V13 Z" fill={`url(#${id})`} />
      <path d="M32 4 L56 13 V31 C56 45 45 55 32 60 Z" fill="#000" opacity="0.14" />
      <path d="M32 14 L46 19.5 V31 C46 39.5 40 45.5 32 49 C24 45.5 18 39.5 18 31 V19.5 Z" fill="#0F1117" opacity="0.28" />
      <text
        x="32"
        y="35"
        textAnchor="middle"
        dominantBaseline="central"
        className="font-display"
        style={{ fill: "#0F1117", fontSize: 18, fontWeight: 800 }}
      >
        {tier.index + 1}
      </text>
    </svg>
  );
}
