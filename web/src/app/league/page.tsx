"use client";

import { useRouter } from "next/navigation";
import { Fragment, useCallback, useEffect, useState } from "react";
import { api, type League, type LeagueStanding } from "@/lib/api";
import { companionVisual } from "@/lib/design/companions";
import { resultMessage, tierTint, timeLeft, zoneBoundaries, zoneLabel } from "@/lib/leagues";
import { GameBackground } from "@/components/game/GameBackground";
import { GameButton } from "@/components/game/GameButton";
import { BottomNav } from "@/components/game/BottomNav";
import { CompanionHero } from "@/components/game/CompanionHero";
import { TierEmblem } from "@/components/league/TierEmblem";

export default function LeaguePage() {
  const router = useRouter();
  const [league, setLeague] = useState<League | null>(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);
  const [now, setNow] = useState(() => new Date());

  const load = useCallback(async () => {
    setError("");
    setLoading(true);
    try {
      setLeague(await api.getLeague());
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to load league");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    let cancelled = false;
    (async () => {
      const ok = await api.ensureSession();
      if (cancelled) return;
      if (!ok) {
        router.replace("/login");
        return;
      }
      load();
    })();
    return () => {
      cancelled = true;
    };
  }, [router, load]);

  // The countdown only needs minute precision.
  useEffect(() => {
    const id = setInterval(() => setNow(new Date()), 30_000);
    return () => clearInterval(id);
  }, []);

  if (loading) {
    return (
      <GameBackground>
        <main className="flex min-h-dvh flex-col items-center justify-center pb-20">
          <p className="font-mono animate-pulse text-sm font-semibold" style={{ color: "#7C6EF5" }} role="status">
            Loading your league…
          </p>
        </main>
        <BottomNav />
      </GameBackground>
    );
  }

  if (!league) {
    return (
      <GameBackground>
        <main className="flex min-h-dvh flex-col items-center justify-center gap-5 p-6 text-center">
          <p className="font-display text-lg font-extrabold" style={{ color: "#E8EAED" }}>
            Couldn&apos;t load your league
          </p>
          <p style={{ color: "#F87171" }} role="alert" className="text-sm font-semibold">
            {error || "Something went wrong"}
          </p>
          <div className="w-full max-w-xs">
            <GameButton type="button" onClick={load}>
              Try again
            </GameButton>
          </div>
        </main>
        <BottomNav />
      </GameBackground>
    );
  }

  const tint = tierTint(league.tier);
  const dividers = new Set(zoneBoundaries(league.standings.map((s) => s.zone)));

  return (
    <GameBackground>
      <main className="mx-auto max-w-2xl px-4 pb-28 pt-6">
        <header className="flex flex-col items-center text-center">
          <p className="font-mono text-[11px] font-bold uppercase tracking-[0.2em]" style={{ color: "#7C6EF5" }}>
            Weekly league
          </p>
          <div className="mt-4">
            <TierEmblem tier={league.tier} size={88} />
          </div>
          <h1 className="font-display mt-3 text-3xl font-extrabold tracking-tight" style={{ color: "#E8EAED" }}>
            {league.tier.name} League
          </h1>
          <p
            className="font-mono mt-3 rounded-full px-3 py-1 text-xs font-bold"
            style={{ color: tint, background: `${tint}1A`, border: `1px solid ${tint}40` }}
          >
            <span aria-hidden>⏱ </span>
            Ends in {timeLeft(league.ends_at, now)}
          </p>
        </header>

        <ol className="mt-6 flex items-center justify-center gap-1.5 overflow-x-auto pb-1" aria-label="League ladder">
          {league.tiers.map((t) => (
            <li key={t.slug}>
              <TierEmblem
                tier={t}
                size={t.index === league.tier.index ? 34 : 26}
                state={t.index === league.tier.index ? "current" : t.index < league.tier.index ? "reached" : "ahead"}
              />
            </li>
          ))}
        </ol>

        {league.last_result && (
          <p
            className="mt-6 rounded-2xl px-4 py-3 text-sm font-semibold"
            style={{ background: "#1A1D27", border: "1px solid #2E3347", color: "#E8EAED" }}
          >
            <span className="font-mono mr-2 text-[10px] uppercase tracking-[0.18em]" style={{ color: "#8B92A8" }}>
              Last week
            </span>
            {resultMessage(league.last_result)}
          </p>
        )}

        {!league.joined ? (
          <section
            className="mt-6 flex flex-col items-center gap-4 rounded-3xl px-6 py-8 text-center"
            style={{ background: "#1A1D27", border: "1px solid #2E3347" }}
          >
            <CompanionHero size="md" />
            <h2 className="font-display text-xl font-extrabold" style={{ color: "#E8EAED" }}>
              Finish a lesson to join this week
            </h2>
            <p className="max-w-sm text-sm" style={{ color: "#8B92A8" }}>
              Your first lesson this week puts you on a leaderboard with up to 30 learners in {league.tier.name}. The top
              ranks move up when the week ends.
            </p>
            <div className="w-full max-w-xs">
              <GameButton type="button" onClick={() => router.push("/journey")}>
                Go to the journey
              </GameButton>
            </div>
          </section>
        ) : (
          <section className="mt-6" aria-labelledby="standings-heading">
            <h2 id="standings-heading" className="sr-only">
              Standings
            </h2>
            <ol className="overflow-hidden rounded-3xl" style={{ background: "#1A1D27", border: "1px solid #2E3347" }}>
              {league.standings.map((s, i) => (
                <Fragment key={s.user_id}>
                  <StandingRow standing={s} />
                  {dividers.has(i) && <ZoneDivider below={league.standings[i + 1].zone} />}
                </Fragment>
              ))}
            </ol>
          </section>
        )}
      </main>
      <BottomNav />
    </GameBackground>
  );
}

function StandingRow({ standing: s }: { standing: LeagueStanding }) {
  const visual = companionVisual(s.companion_name, s.companion_species);
  const rankTint = s.zone === "promote" ? "#34D399" : s.zone === "demote" ? "#F87171" : "#8B92A8";
  return (
    <li
      className="flex items-center gap-3 px-4 py-3"
      style={{
        background: s.is_me ? "rgba(124,110,245,0.14)" : "transparent",
        borderTop: "1px solid #232736",
      }}
      aria-current={s.is_me ? "true" : undefined}
    >
      <span className="font-display w-7 text-center text-base font-extrabold tabular-nums" style={{ color: rankTint }}>
        {s.rank}
      </span>
      <span
        aria-hidden
        className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full text-xl"
        style={{ background: `linear-gradient(135deg, ${visual.glow}33 0%, #242836 100%)`, border: "1px solid #2E3347" }}
      >
        {visual.emoji}
      </span>
      <span className="min-w-0 flex-1">
        <span className="block truncate font-semibold" style={{ color: "#E8EAED" }}>
          {s.username}
          {s.is_me && (
            <span className="font-mono ml-2 text-[10px] font-bold uppercase tracking-wider" style={{ color: "#7C6EF5" }}>
              You
            </span>
          )}
        </span>
        <span className="sr-only">{zoneLabel(s.zone)}</span>
      </span>
      <span className="font-mono text-sm font-bold tabular-nums" style={{ color: "#60A5FA" }}>
        {s.weekly_xp} XP
      </span>
    </li>
  );
}

function ZoneDivider({ below }: { below: LeagueStanding["zone"] }) {
  const isDrop = below === "demote";
  const tint = isDrop ? "#F87171" : "#34D399";
  return (
    <li
      aria-hidden
      className="font-mono flex items-center gap-2 px-4 py-1.5 text-[10px] font-bold uppercase tracking-[0.18em]"
      style={{ color: tint, background: `${tint}0F` }}
    >
      <span>{isDrop ? "▼" : "▲"}</span>
      {isDrop ? "Drop zone" : "Promotion zone above"}
    </li>
  );
}
