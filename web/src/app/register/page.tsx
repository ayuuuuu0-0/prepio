"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { FormEvent, useState } from "react";
import { GameBackground } from "@/components/game/GameBackground";
import { CompanionHero } from "@/components/game/CompanionHero";
import { SpeechBubble } from "@/components/game/SpeechBubble";
import { GameButton } from "@/components/game/GameButton";
import { api } from "@/lib/api";

export default function RegisterPage() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");
    try {
      const res = await api.register(email, username, password);
      api.setToken(res.access_token);
      router.push("/onboarding");
    } catch (err) {
      setError(err instanceof Error ? err.message : "registration failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <GameBackground>
      <main className="mx-auto flex min-h-screen max-w-md flex-col items-center justify-center px-6 py-12">
        <CompanionHero name="Pip" species="red_panda" size="lg" />
        <SpeechBubble className="mt-6 max-w-sm text-center">
          Join the adventure! Pick your companion and start becoming interview-ready.
        </SpeechBubble>

        <form onSubmit={onSubmit} className="mt-8 w-full space-y-4">
          {error && (
            <p className="rounded-2xl bg-orange-100 px-4 py-3 text-center text-sm font-semibold text-orange-700">
              {error}
            </p>
          )}
          <input
            className="w-full rounded-2xl border-2 border-[#E5E5E5] bg-white px-4 py-3 font-semibold outline-none focus:border-[#FF9600]"
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
          <input
            className="w-full rounded-2xl border-2 border-[#E5E5E5] bg-white px-4 py-3 font-semibold outline-none focus:border-[#FF9600]"
            placeholder="Username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
          />
          <input
            className="w-full rounded-2xl border-2 border-[#E5E5E5] bg-white px-4 py-3 font-semibold outline-none focus:border-[#FF9600]"
            type="password"
            placeholder="Password (8+ chars)"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            minLength={8}
            required
          />
          <GameButton type="submit" variant="gold" disabled={loading}>
            {loading ? "Creating..." : "Begin Adventure!"}
          </GameButton>
        </form>

        <p className="mt-6 text-center text-sm font-semibold text-[#777]">
          Already playing?{" "}
          <Link href="/login" className="font-display font-bold text-[#58CC02] hover:underline">
            Sign in
          </Link>
        </p>
      </main>
    </GameBackground>
  );
}
