"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useCallback, useEffect, useState } from "react";
import Image from "next/image";
import { api, DashboardHome } from "@/lib/api";
import { companyColors, leagueThemes } from "@/lib/design/tokens";
import { GameBackground } from "@/components/game/GameBackground";
import { CompanionHero } from "@/components/game/CompanionHero";
import { SpeechBubble } from "@/components/game/SpeechBubble";
import { ReadinessRing } from "@/components/game/ReadinessRing";
import { GameCard } from "@/components/game/GameCard";
import { QuestCard } from "@/components/game/QuestCard";
import { GameButton } from "@/components/game/GameButton";
import { BottomNav } from "@/components/game/BottomNav";
import { HUDBar } from "@/components/game/HUDBar";

const questIcons: Record<string, string> = {
  daily_question: "⚡",
  maintain_streak: "🔥",
  score_high: "🎯",
};

export default function DashboardPage() {
  const router = useRouter();
  const [home, setHome] = useState<DashboardHome | null>(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);

  const load = useCallback(async () => {
    setError("");
    try {
      const data = await api.getDashboardHome();
      if (data.onboarding_needed) {
        router.replace("/onboarding");
        return;
      }
      setHome(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to load dashboard");
    } finally {
      setLoading(false);
    }
  }, [router]);

  useEffect(() => {
    (async () => {
      const ok = await api.ensureSession();
      if (!ok) {
        router.replace("/login");
        return;
      }
      load();
    })();
  }, [router, load]);

  if (loading) {
    return (
      <GameBackground>
        <main className="flex min-h-screen items-center justify-center pb-20">
          <CompanionHero name="Byte" size="md" />
          <p className="font-mono mt-4 animate-pulse text-sm font-semibold" style={{ color: "#7C6EF5" }}>
            Loading...
          </p>
        </main>
      </GameBackground>
    );
  }

  if (!home) {
    return (
      <GameBackground>
        <main className="flex min-h-screen items-center justify-center p-6">
          <p style={{ color: "#F87171" }}>{error || "Something went wrong"}</p>
        </main>
      </GameBackground>
    );
  }

  const league = leagueThemes[home.league.tier] ?? leagueThemes.bronze;

  return (
    <GameBackground>
      <main className="mx-auto max-w-lg px-4 pb-28 pt-6">
        <div className="flex items-center gap-2.5 mb-6 justify-start border-b pb-4" style={{ borderColor: "#2E3347" }}>
          <Image src="/logo.png" alt="Prepio Logo" width={28} height={28} className="rounded-lg" />
          <span className="font-display text-base font-extrabold tracking-wide" style={{ color: "#7C6EF5" }}>
            PREPIO
          </span>
          <span className="font-mono text-[9px] font-bold px-2 py-0.5 rounded-md" style={{ color: "#8B92A8", background: "#1A1D27", border: "1px solid #2E3347" }}>
            Beta
          </span>
        </div>

        <section className="flex items-start gap-4">
          <CompanionHero name={home.companion?.name} species={home.companion?.species} size="md" />
          <SpeechBubble className="mt-2 flex-1" speakerName={home.companion?.name ?? "Byte"}>
            {home.companion_message}
          </SpeechBubble>
        </section>

        <div className="mt-4">
          <HUDBar home={home} />
        </div>

        <GameCard className="mt-5" icon="🧭" accentColor="#60A5FA">
          <h2 className="font-display mb-4 text-base font-bold" style={{ color: "#E8EAED" }}>
            Career Readiness
          </h2>
          <div className="flex flex-wrap justify-center gap-4">
            {home.readiness.map((r, i) => (
              <ReadinessRing
                key={r.company}
                company={r.company}
                score={r.score}
                color={companyColors[r.company]?.ring ?? "#7C6EF5"}
                delay={i * 150}
              />
            ))}
          </div>
        </GameCard>

        <GameCard className="mt-4" icon={league.icon} accentColor={league.border}>
          {home.league.available ? (
            <>
              <div
                className={`-mx-5 -mt-5 mb-4 rounded-t-2xl bg-gradient-to-r ${league.gradient} px-5 pt-4 pb-3`}
              >
                <p className="font-mono text-[11px] font-bold uppercase tracking-widest text-white/70">
                  {home.league.label}
                </p>
                <p className="font-display mt-0.5 text-2xl font-extrabold text-white">
                  Rank #{home.league.rank}
                </p>
              </div>
              <p className="text-sm" style={{ color: "#8B92A8" }}>
                Keep climbing — promotion zone ahead!
              </p>
            </>
          ) : (
            <>
              <p className="font-display text-lg font-extrabold" style={{ color: "#E8EAED" }}>
                {home.league.label}
              </p>
              <p className="mt-1 text-sm" style={{ color: "#8B92A8" }}>
                Your league rank is calculating — check back soon!
              </p>
            </>
          )}
        </GameCard>

        <div className="mt-5 space-y-3">
          <h2 className="font-mono text-xs font-bold uppercase tracking-widest" style={{ color: "#4A5068" }}>
            Daily Quests
          </h2>
          {home.daily_quests.map((q) => (
            <QuestCard
              key={q.id}
              title={q.title}
              icon={questIcons[q.id] ?? "📋"}
              progress={q.progress}
              target={q.target}
              completed={q.completed}
              rewardXp={q.reward_xp}
              rewardGems={q.reward_gems}
              comingSoon={q.coming_soon}
            />
          ))}
        </div>

        <Link href="/challenge" className="mt-6 block">
          <GameButton type="button">Continue Prep →</GameButton>
        </Link>

        <button
          onClick={() => {
            api.setAuthTokens(null, null);
            router.push("/login");
          }}
          className="mt-4 w-full text-center font-mono text-xs font-semibold transition-colors"
          style={{ color: "#4A5068" }}
        >
          Sign out
        </button>
      </main>
      <BottomNav />
    </GameBackground>
  );
}
