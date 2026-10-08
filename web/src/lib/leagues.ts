import type { LeagueResult, LeagueTier, LeagueZone } from "@/lib/api";

/** Tier tints by slug. Presentation only: the ladder itself comes from the server. */
const tierTints: Record<string, string> = {
  bronze: "#D08A5B",
  silver: "#B8C2D6",
  gold: "#F5B942",
  platinum: "#8FE3D8",
  sapphire: "#60A5FA",
  ruby: "#F2667A",
  emerald: "#34D399",
  amethyst: "#A78BFA",
  obsidian: "#9AA3B8",
  diamond: "#7DD3FC",
};

/** tierTint returns the tint for a tier, falling back to the accent for tiers added later. */
export function tierTint(tier: Pick<LeagueTier, "slug">): string {
  return tierTints[tier.slug] ?? "#7C6EF5";
}

/** timeLeft formats the time until the league week ends: "3d 4h", "5h 12m", "12m", or "ending now". */
export function timeLeft(endsAt: string, now: Date = new Date()): string {
  const ms = new Date(endsAt).getTime() - now.getTime();
  if (!Number.isFinite(ms) || ms <= 60_000) return "ending now";
  const minutes = Math.floor(ms / 60_000);
  const days = Math.floor(minutes / 1440);
  const hours = Math.floor((minutes % 1440) / 60);
  const mins = minutes % 60;
  if (days > 0) return `${days}d ${hours}h`;
  if (hours > 0) return `${hours}h ${mins}m`;
  return `${mins}m`;
}

/** zoneLabel names a standing's zone in words, so the zone never relies on colour alone. */
export function zoneLabel(zone: LeagueZone): string {
  switch (zone) {
    case "promote":
      return "Promotion zone";
    case "demote":
      return "Drop zone";
    default:
      return "Safe";
  }
}

/** resultMessage tells the learner how last week ended, encouragingly. */
export function resultMessage(last: LeagueResult): string {
  switch (last.outcome) {
    case "promoted":
      return `You finished #${last.rank} and moved up to ${last.to_tier.name}.`;
    case "demoted":
      return `You finished #${last.rank}. This week you're in ${last.to_tier.name}, a good place to win it back.`;
    default:
      return `You finished #${last.rank} and held your place in ${last.to_tier.name}.`;
  }
}

/** zoneBoundaries returns the indexes after which a zone divider is drawn (promote→stay, stay→demote). */
export function zoneBoundaries(zones: LeagueZone[]): number[] {
  const out: number[] = [];
  for (let i = 0; i < zones.length - 1; i++) {
    if (zones[i] !== zones[i + 1]) out.push(i);
  }
  return out;
}
