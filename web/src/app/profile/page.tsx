"use client";

import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { api, Profile, TopicInfo } from "@/lib/api";
import { topicTheme } from "@/lib/topics";
import { achievementIcon, earnedLabel } from "@/lib/achievements";
import type { Achievement } from "@/lib/lesson/types";
import { GameBackground } from "@/components/game/GameBackground";
import { CompanionHero } from "@/components/game/CompanionHero";
import { GameCard } from "@/components/game/GameCard";
import { BottomNav } from "@/components/game/BottomNav";

export default function ProfilePage() {
  const router = useRouter();
  const [profile, setProfile] = useState<Profile | null>(null);
  const [topics, setTopics] = useState<TopicInfo[]>([]);
  const [achievements, setAchievements] = useState<Achievement[] | null>(null);

  useEffect(() => {
    (async () => {
      const ok = await api.ensureSession();
      if (!ok) {
        router.replace("/login");
        return;
      }
      try {
        const [me, catalog] = await Promise.all([api.getProfile(), api.getTopics().catch(() => [] as TopicInfo[])]);
        setProfile(me);
        setTopics(catalog);
        // Achievements are a section of the profile; if they fail to load the section says so.
        api.getAchievements().then(setAchievements, () => setAchievements([]));
      } catch {
        router.replace("/login");
      }
    })();
  }, [router]);

  if (!profile) {
    return (
      <GameBackground>
        <main className="flex min-h-screen items-center justify-center">
          <p className="font-mono animate-pulse text-sm font-semibold" style={{ color: "#7C6EF5" }}>
            Loading profile...
          </p>
        </main>
      </GameBackground>
    );
  }

  return (
    <GameBackground>
      <main className="mx-auto max-w-lg px-4 pb-28 pt-8">
        <div className="flex flex-col items-center">
          <CompanionHero name={profile.companion?.name} species={profile.companion?.species} size="lg" />
          <h1 className="font-display mt-4 text-2xl font-extrabold" style={{ color: "#E8EAED" }}>
            {profile.username}
          </h1>
          <p className="font-mono text-sm font-semibold" style={{ color: "#8B92A8" }}>
            {profile.email}
          </p>
        </div>

        <GameCard className="mt-6" icon="📚" accentColor="#60A5FA">
          <p className="font-display font-bold" style={{ color: "#E8EAED" }}>
            Experience
          </p>
          <p className="mt-2 text-sm font-semibold capitalize" style={{ color: "#8B92A8" }}>
            {profile.experience_level ?? "Not set"}
          </p>
        </GameCard>

        {profile.focus_topics && profile.focus_topics.length > 0 && (
          <GameCard className="mt-4" icon="🎯" accentColor="#7C6EF5">
            <p className="font-display font-bold" style={{ color: "#E8EAED" }}>
              Focus topics
            </p>
            <ul className="mt-3 flex flex-wrap gap-2">
              {profile.focus_topics.map((slug) => {
                const theme = topicTheme(slug);
                const name = topics.find((t) => t.slug === slug)?.name ?? slug;
                return (
                  <li
                    key={slug}
                    className="font-mono rounded-full px-3 py-1 text-xs font-bold"
                    style={{ background: `${theme.tint}22`, border: `1px solid ${theme.tint}55`, color: theme.tint }}
                  >
                    {theme.icon} {name}
                  </li>
                );
              })}
            </ul>
          </GameCard>
        )}

        <section className="mt-6" aria-labelledby="achievements-title">
          <div className="flex items-end justify-between">
            <h2 id="achievements-title" className="font-display text-lg font-extrabold" style={{ color: "#E8EAED" }}>
              Achievements
            </h2>
            {achievements && achievements.length > 0 && (
              <p className="font-mono text-xs font-bold" style={{ color: "#8B92A8" }}>
                {achievements.filter((a) => a.unlocked).length} of {achievements.length}
              </p>
            )}
          </div>
          {achievements === null ? (
            <p className="font-mono mt-3 animate-pulse text-xs" style={{ color: "#8B92A8" }} role="status">
              Loading achievements…
            </p>
          ) : achievements.length === 0 ? (
            <p className="mt-3 text-sm" style={{ color: "#8B92A8" }} role="alert">
              Couldn&apos;t load achievements right now.
            </p>
          ) : (
            <ul className="mt-3 grid grid-cols-2 gap-3">
              {achievements.map((a) => (
                <li
                  key={a.slug}
                  className="rounded-2xl p-4"
                  style={{
                    background: a.unlocked ? "rgba(245,185,66,0.1)" : "#1A1D27",
                    border: `1px solid ${a.unlocked ? "rgba(245,185,66,0.45)" : "#2E3347"}`,
                  }}
                >
                  <span aria-hidden className="text-2xl" style={{ filter: a.unlocked ? "none" : "grayscale(1)", opacity: a.unlocked ? 1 : 0.45 }}>
                    {achievementIcon(a.slug)}
                  </span>
                  <p className="font-display mt-2 text-sm font-extrabold" style={{ color: a.unlocked ? "#E8EAED" : "#8B92A8" }}>
                    {a.name}
                    <span className="sr-only">{a.unlocked ? ", earned" : ", not earned yet"}</span>
                  </p>
                  <p className="mt-1 text-xs leading-snug" style={{ color: "#8B92A8" }}>
                    {a.unlocked ? earnedLabel(a.unlocked_at) : a.description}
                  </p>
                </li>
              ))}
            </ul>
          )}
        </section>

        <button
          onClick={async () => {
            await api.logout();
            router.push("/login");
          }}
          className="mt-8 w-full text-center font-mono text-sm font-semibold transition-colors"
          style={{ color: "#4A5068" }}
        >
          Sign out
        </button>
      </main>
      <BottomNav />
    </GameBackground>
  );
}
