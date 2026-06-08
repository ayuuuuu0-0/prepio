"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { FormEvent, useState } from "react";
import { GameBackground } from "@/components/game/GameBackground";
import { CompanionHero } from "@/components/game/CompanionHero";
import { SpeechBubble } from "@/components/game/SpeechBubble";
import { GameButton } from "@/components/game/GameButton";
import { api } from "@/lib/api";

const inputStyle = {
  background: "#1A1D27",
  border: "1px solid #2E3347",
  color: "#E8EAED",
};

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
      api.setAuthTokens(res.access_token, res.refresh_token);
      router.push("/onboarding");
    } catch (err) {
      setError(err instanceof Error ? err.message : "registration failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <GameBackground>
      <div className="flex min-h-screen">
        <div
          className="hidden lg:flex lg:w-1/2 flex-col justify-between p-12"
          style={{ borderRight: "1px solid #2E3347" }}
        >
          <div>
            <p className="font-display text-xl font-bold" style={{ color: "#7C6EF5" }}>
              PREPIO
            </p>
            <p className="font-mono text-xs font-medium mt-1" style={{ color: "#4A5068" }}>
              Career RPG
            </p>
          </div>

          <div className="flex flex-col items-center gap-6 py-12">
            <CompanionHero name="Pip" species="red_panda" size="lg" />
            <SpeechBubble speakerName="Pip" className="max-w-xs text-center">
              Pick a companion. Start tracking your readiness from day one.
            </SpeechBubble>
          </div>

          <div className="grid grid-cols-3 gap-4">
            {[
              { label: "Engineers in prep", value: "12k+" },
              { label: "Avg readiness gain", value: "34%" },
              { label: "Offers received", value: "2.1k" },
            ].map((s) => (
              <div key={s.label}>
                <p className="font-mono text-2xl font-bold" style={{ color: "#7C6EF5" }}>
                  {s.value}
                </p>
                <p className="font-body text-xs mt-1" style={{ color: "#4A5068" }}>
                  {s.label}
                </p>
              </div>
            ))}
          </div>
        </div>

        <div className="flex w-full lg:w-1/2 flex-col items-center justify-center px-6 py-12">
          <div className="mb-8 text-center lg:hidden">
            <p className="font-display text-2xl font-bold" style={{ color: "#7C6EF5" }}>
              PREPIO
            </p>
            <CompanionHero name="Pip" species="red_panda" size="md" />
          </div>

          <div className="w-full max-w-sm">
            <h1 className="font-display text-2xl font-bold mb-1" style={{ color: "#E8EAED" }}>
              Create account
            </h1>
            <p className="font-body text-sm mb-8" style={{ color: "#4A5068" }}>
              Pick a companion. Start tracking your readiness.
            </p>

            <form onSubmit={onSubmit} className="space-y-4">
              {error && (
                <p
                  className="rounded-xl px-4 py-3 text-sm font-medium"
                  style={{
                    background: "rgba(248,113,113,0.1)",
                    border: "1px solid rgba(248,113,113,0.3)",
                    color: "#F87171",
                  }}
                >
                  {error}
                </p>
              )}

              <input
                className="w-full rounded-xl px-4 py-3 text-sm font-medium outline-none transition-all"
                style={inputStyle}
                onFocus={(e) => (e.target.style.borderColor = "#7C6EF5")}
                onBlur={(e) => (e.target.style.borderColor = "#2E3347")}
                type="email"
                placeholder="Email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
              <input
                className="w-full rounded-xl px-4 py-3 text-sm font-medium outline-none transition-all"
                style={inputStyle}
                onFocus={(e) => (e.target.style.borderColor = "#7C6EF5")}
                onBlur={(e) => (e.target.style.borderColor = "#2E3347")}
                placeholder="Username"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                required
              />
              <input
                className="w-full rounded-xl px-4 py-3 text-sm font-medium outline-none transition-all"
                style={inputStyle}
                onFocus={(e) => (e.target.style.borderColor = "#7C6EF5")}
                onBlur={(e) => (e.target.style.borderColor = "#2E3347")}
                type="password"
                placeholder="Password (8+ chars)"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                minLength={8}
                required
              />

              <GameButton type="submit" variant="gold" disabled={loading}>
                {loading ? "Creating..." : "Create Account"}
              </GameButton>
            </form>

            <p className="mt-6 text-center text-sm" style={{ color: "#4A5068" }}>
              Already have an account?{" "}
              <Link href="/login" className="font-semibold" style={{ color: "#7C6EF5" }}>
                Sign in
              </Link>
            </p>
          </div>
        </div>
      </div>
    </GameBackground>
  );
}
