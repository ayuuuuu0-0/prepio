"use client";

import { useRouter } from "next/navigation";
import { useCallback, useEffect, useState } from "react";
import Image from "next/image";
import { api, DashboardHome } from "@/lib/api";
import { GameBackground } from "@/components/game/GameBackground";
import { CompanionHero } from "@/components/game/CompanionHero";
import { SpeechBubble } from "@/components/game/SpeechBubble";
import { GameButton } from "@/components/game/GameButton";
import { BottomNav } from "@/components/game/BottomNav";
import { HUDBar } from "@/components/game/HUDBar";
import { ContinueCard } from "@/components/dashboard/ContinueCard";
import { LeagueCard } from "@/components/dashboard/LeagueCard";
import { TopicMasteryCard } from "@/components/dashboard/TopicMasteryCard";

export default function DashboardPage() {
  const router = useRouter();
  const [home, setHome] = useState<DashboardHome | null>(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);

  const load = useCallback(async () => {
    setError("");
    setLoading(true);
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
        <main className="flex min-h-dvh flex-col items-center justify-center pb-20">
          <CompanionHero name="Byte" size="md" />
          <p className="font-mono mt-4 animate-pulse text-sm font-semibold" style={{ color: "#7C6EF5" }} role="status">
            Loading…
          </p>
        </main>
      </GameBackground>
    );
  }

  if (!home) {
    return (
      <GameBackground>
        <main className="flex min-h-dvh flex-col items-center justify-center gap-5 p-6 text-center">
          <p className="font-display text-lg font-extrabold" style={{ color: "#E8EAED" }}>
            Couldn&apos;t load your dashboard
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
      </GameBackground>
    );
  }

  return (
    <GameBackground>
      <main className="mx-auto max-w-3xl px-4 pb-28 pt-6">
        <div className="mb-6 flex items-center justify-start gap-2.5 border-b pb-4" style={{ borderColor: "#2E3347" }}>
          <Image src="/logo.png" alt="Prepio Logo" width={28} height={28} className="rounded-lg" />
          <span className="font-display text-base font-extrabold tracking-wide" style={{ color: "#7C6EF5" }}>
            PREPIO
          </span>
          <span
            className="font-mono rounded-md px-2 py-0.5 text-[9px] font-bold"
            style={{ color: "#8B92A8", background: "#1A1D27", border: "1px solid #2E3347" }}
          >
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

        <div className="mt-5">
          <ContinueCard lesson={home.next_lesson} companionName={home.companion?.name} />
        </div>

        <div className="mt-4">
          <LeagueCard league={home.league} />
        </div>

        <section className="mt-8" aria-labelledby="topics-heading">
          <div className="flex items-end justify-between gap-3">
            <div>
              <p className="font-mono text-[11px] font-bold uppercase tracking-[0.2em]" style={{ color: "#7C6EF5" }}>
                Your readiness
              </p>
              <h2 id="topics-heading" className="font-display mt-1 text-2xl font-extrabold leading-tight tracking-tight" style={{ color: "#E8EAED" }}>
                Where you stand, by topic
              </h2>
            </div>
          </div>
          <div className="mt-4 grid gap-4 sm:grid-cols-2">
            {home.topics.map((topic) => (
              <TopicMasteryCard key={topic.slug} topic={topic} />
            ))}
          </div>
        </section>

        <button
          onClick={async () => {
            await api.logout();
            router.push("/login");
          }}
          className="mt-10 w-full text-center font-mono text-xs font-semibold transition-colors"
          style={{ color: "#4A5068" }}
        >
          Sign out
        </button>
      </main>
      <BottomNav />
    </GameBackground>
  );
}
