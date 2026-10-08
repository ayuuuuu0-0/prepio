"use client";

import Link from "next/link";
import type { DashboardHome } from "@/lib/api";
import { tierTint, timeLeft, zoneLabel } from "@/lib/leagues";
import { TierEmblem } from "@/components/league/TierEmblem";

/**
 * LeagueCard puts this week's league on the dashboard: where you rank and what it would earn,
 * or, before your first lesson of the week, that one lesson joins you. It links to /league.
 */
export function LeagueCard({ league }: { league: DashboardHome["league"] }) {
  const tier = { index: league.tier_index, slug: league.tier_slug, name: league.tier_name };
  const tint = tierTint(tier);
  const zoneTint = league.zone === "promote" ? "#34D399" : league.zone === "demote" ? "#F5B942" : "#8B92A8";

  return (
    <Link
      href="/league"
      className="flex items-center gap-4 rounded-3xl px-5 py-4 transition-transform hover:-translate-y-0.5 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7C6EF5]"
      style={{ background: `linear-gradient(135deg, ${tint}1F 0%, #1A1D27 60%)`, border: `1px solid ${tint}40` }}
    >
      <TierEmblem tier={tier} size={48} />
      <div className="min-w-0 flex-1">
        <p className="font-mono text-[11px] font-bold uppercase tracking-[0.2em]" style={{ color: tint }}>
          {league.tier_name} League
        </p>
        {league.joined ? (
          <>
            <p className="font-display mt-0.5 text-lg font-extrabold" style={{ color: "#E8EAED" }}>
              #{league.rank} of {league.cohort_size}
              <span className="ml-2 text-sm font-bold" style={{ color: zoneTint }}>
                {zoneLabel(league.zone)}
              </span>
            </p>
            <p className="text-xs" style={{ color: "#8B92A8" }}>
              {league.weekly_xp} XP this week · ends in {timeLeft(league.ends_at)}
            </p>
          </>
        ) : (
          <>
            <p className="font-display mt-0.5 text-base font-extrabold" style={{ color: "#E8EAED" }}>
              One lesson joins you this week
            </p>
            <p className="text-xs" style={{ color: "#8B92A8" }}>
              Week ends in {timeLeft(league.ends_at)}
            </p>
          </>
        )}
      </div>
      <span aria-hidden className="text-lg" style={{ color: "#8B92A8" }}>
        →
      </span>
    </Link>
  );
}
