"use client";

import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { api, Profile, TopicInfo } from "@/lib/api";
import { topicTheme } from "@/lib/topics";
import { GameBackground } from "@/components/game/GameBackground";
import { CompanionHero } from "@/components/game/CompanionHero";
import { GameCard } from "@/components/game/GameCard";
import { BottomNav } from "@/components/game/BottomNav";

export default function ProfilePage() {
  const router = useRouter();
  const [profile, setProfile] = useState<Profile | null>(null);
  const [topics, setTopics] = useState<TopicInfo[]>([]);

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

        <button
          onClick={() => {
            api.setAuthTokens(null, null);
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
