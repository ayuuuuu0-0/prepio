"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { FormEvent, useState } from "react";
import { GameBackground } from "@/components/game/GameBackground";
import { CompanionHero } from "@/components/game/CompanionHero";
import { SpeechBubble } from "@/components/game/SpeechBubble";
import { GameButton } from "@/components/game/GameButton";
import { api } from "@/lib/api";

export default function LoginPage() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");
    try {
      const res = await api.login(email, password);
      api.setToken(res.access_token);
      const profile = await api.getProfile();
      router.push(profile.onboarding_completed ? "/dashboard" : "/onboarding");
    } catch (err) {
      setError(err instanceof Error ? err.message : "login failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <GameBackground>
      <main className="mx-auto flex min-h-screen max-w-md flex-col items-center justify-center px-6 py-12">
        <CompanionHero name="Byte" species="capybara" size="lg" />
        <SpeechBubble className="mt-6 max-w-sm text-center">
          Welcome back! Ready to level up your career today?
        </SpeechBubble>

        <form onSubmit={onSubmit} className="mt-8 w-full space-y-4">
          {error && (
            <p className="rounded-2xl bg-orange-100 px-4 py-3 text-center text-sm font-semibold text-orange-700">
              {error}
            </p>
          )}
          <input
            className="w-full rounded-2xl border-2 border-[#E5E5E5] bg-white px-4 py-3 font-semibold outline-none focus:border-[#58CC02]"
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
          <input
            className="w-full rounded-2xl border-2 border-[#E5E5E5] bg-white px-4 py-3 font-semibold outline-none focus:border-[#58CC02]"
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
          <GameButton type="submit" disabled={loading}>
            {loading ? "Signing in..." : "Let's Go!"}
          </GameButton>
        </form>

        <p className="mt-6 text-center text-sm font-semibold text-[#777]">
          New adventurer?{" "}
          <Link href="/register" className="font-display font-bold text-[#1CB0F6] hover:underline">
            Start your journey
          </Link>
        </p>
      </main>
    </GameBackground>
  );
}
