"use client";

import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { api, Profile } from "@/lib/api";
import { GameBackground } from "@/components/game/GameBackground";
import { CompanionHero } from "@/components/game/CompanionHero";
import { GameCard } from "@/components/game/GameCard";
import { BottomNav } from "@/components/game/BottomNav";

export default function ProfilePage() {
  const router = useRouter();
  const [profile, setProfile] = useState<Profile | null>(null);

  useEffect(() => {
    (async () => {
      const ok = await api.ensureSession();
      if (!ok) {
        router.replace("/login");
        return;
      }
      try {
        setProfile(await api.getProfile());
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

        <GameCard className="mt-6" icon="🎯" accentColor="#7C6EF5">
          <p className="font-display font-bold" style={{ color: "#E8EAED" }}>
            Target Companies
          </p>
          <p className="mt-2 text-sm font-semibold" style={{ color: "#8B92A8" }}>
            {profile.target_companies.length > 0 ? profile.target_companies.join(", ") : "Not set"}
          </p>
        </GameCard>

        <GameCard className="mt-4" icon="📚" accentColor="#60A5FA">
          <p className="font-display font-bold" style={{ color: "#E8EAED" }}>
            Experience
          </p>
          <p className="mt-2 text-sm font-semibold capitalize" style={{ color: "#8B92A8" }}>
            {profile.experience_level ?? "Not set"}
          </p>
        </GameCard>

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
