"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useCallback, useEffect, useState } from "react";
import { api, DashboardHome } from "@/lib/api";
import { companyColors, leagueThemes } from "@/lib/design/tokens";
import { GameBackground } from "@/components/game/GameBackground";
import { CompanionHero } from "@/components/game/CompanionHero";
import { SpeechBubble } from "@/components/game/SpeechBubble";
import { ReadinessRing } from "@/components/game/ReadinessRing";
import { GameCard } from "@/components/game/GameCard";
import { QuestCard } from "@/components/game/QuestCard";
import { StatChip } from "@/components/game/StatChip";
import { GameButton } from "@/components/game/GameButton";
import { BottomNav } from "@/components/game/BottomNav";

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
    if (!api.loadToken()) {
      router.replace("/login");
      return;
    }
    load();
  }, [router, load]);

  if (loading) {
    return (
      <GameBackground>
        <main className="flex min-h-screen items-center justify-center pb-20">
          <CompanionHero name="Byte" size="md" />
          <p className="font-display mt-4 animate-pulse text-lg font-bold text-[#58CC02]">Loading your world...</p>
        </main>
      </GameBackground>
    );
  }

  if (!home) {
    return (
      <GameBackground>
        <main className="flex min-h-screen items-center justify-center p-6">
          <p className="text-orange-700">{error || "Something went wrong"}</p>
        </main>
      </GameBackground>
    );
  }

  const league = leagueThemes[home.league.tier] ?? leagueThemes.bronze;

  return (
    <GameBackground>
      <main className="mx-auto max-w-lg px-4 pb-28 pt-6">
        {/* Top — Companion + speech */}
        <section className="flex items-end gap-4">
          <CompanionHero name={home.companion?.name} species={home.companion?.species} size="lg" />
          <SpeechBubble className="flex-1">{home.companion_message}</SpeechBubble>
        </section>

        {/* Stats row */}
        <div className="mt-5 flex flex-wrap gap-2">
          <StatChip icon="🔥" label="Streak" value={`${home.streak.current_streak} days`} color="#FF9600" />
          <StatChip icon="⚡" label="Level" value={home.progress.current_level} color="#1CB0F6" />
          <StatChip icon="💎" label="Gems" value={home.progress.gem_balance} color="#7B5CFF" />
        </div>

        {/* Readiness rings */}
        <GameCard className="mt-5 bg-gradient-to-br from-white to-sky-50" icon="🧭">
          <h2 className="font-display mb-4 text-lg font-bold text-[#3C3C3C]">Career Readiness</h2>
          <div className="flex flex-wrap justify-center gap-4">
            {home.readiness.map((r, i) => (
              <ReadinessRing
                key={r.company}
                company={r.company}
                score={r.score}
                color={companyColors[r.company]?.ring ?? "#58CC02"}
                delay={i * 150}
              />
            ))}
          </div>
        </GameCard>

        {/* League */}
        <GameCard
          className={`mt-4 bg-gradient-to-r ${league.gradient} text-white`}
          icon={league.icon}
        >
          <p className="font-display text-sm font-bold uppercase opacity-90">{home.league.label}</p>
          <p className="font-display text-3xl font-extrabold">Rank #{home.league.rank}</p>
          <p className="mt-1 text-sm opacity-80">Keep climbing — promotion zone ahead!</p>
        </GameCard>

        {/* Daily quests */}
        <div className="mt-5 space-y-3">
          <h2 className="font-display text-lg font-bold text-[#3C3C3C]">⚡ Daily Quests</h2>
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
            />
          ))}
        </div>

        {/* Primary CTA */}
        <Link href="/challenge" className="mt-6 block">
          <GameButton type="button">Continue Journey →</GameButton>
        </Link>

        <button
          onClick={() => {
            api.setToken(null);
            router.push("/login");
          }}
          className="mt-4 w-full text-center text-xs font-semibold text-[#999] hover:text-[#777]"
        >
          Sign out
        </button>
      </main>
      <BottomNav />
    </GameBackground>
  );
}
