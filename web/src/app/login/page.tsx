"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import Image from "next/image";
import { FormEvent, useState } from "react";
import { GameBackground } from "@/components/game/GameBackground";
import { CompanionHero } from "@/components/game/CompanionHero";
import { SpeechBubble } from "@/components/game/SpeechBubble";
import { GameButton } from "@/components/game/GameButton";
import { api, ApiRequestError } from "@/lib/api";

const socialProof = [
  { label: "Core Topics", value: "System Design, DSA, LLD" },
  { label: "Lessons captured", value: "Misses, patterns, next steps" },
  { label: "Mastery focus", value: "Skills over streaks" },
];

const inputStyle = {
  background: "#1A1D27",
  border: "1px solid #2E3347",
  color: "#E8EAED",
};

export default function LoginPage() {
  const router = useRouter();
  const [identifier, setIdentifier] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [isNotFound, setIsNotFound] = useState(false);
  const [loading, setLoading] = useState(false);

  async function handlePasswordLogin(e: FormEvent) {
    e.preventDefault();
    if (!identifier.trim() || !password) return;
    setLoading(true);
    setError("");
    setIsNotFound(false);
    try {
      const res = await api.login(identifier.trim(), password);
      api.setAuthTokens(res.access_token);
      const profile = await api.getProfile();
      router.push(profile.onboarding_completed ? "/dashboard" : "/onboarding");
    } catch (err) {
      if (err instanceof ApiRequestError) {
        if (err.code === "user_not_found" || err.status === 404 || err.message.toLowerCase().includes("not exist") || err.message.toLowerCase().includes("not found")) {
          setIsNotFound(true);
          setError("No account found with this email or username.");
        } else {
          setError(err.message || "Incorrect password. Please try again.");
        }
      } else {
        setError(err instanceof Error ? err.message : "Login failed. Please check your credentials.");
      }
    } finally {
      setLoading(false);
    }
  }

  return (
    <GameBackground>
      <div className="flex min-h-screen">
        {/* Left branding column on desktop */}
        <div
          className="hidden lg:flex lg:w-1/2 flex-col justify-between p-12"
          style={{ borderRight: "1px solid #2E3347" }}
        >
          <div className="flex items-center gap-3">
            <Image src="/logo.png" alt="Prepio Logo" width={36} height={36} className="rounded-xl" />
            <div>
              <p className="font-display text-xl font-bold leading-none" style={{ color: "#7C6EF5" }}>
                PREPIO
              </p>
              <p className="font-mono text-[10px] font-medium mt-0.5" style={{ color: "#4A5068" }}>
                Career RPG
              </p>
            </div>
          </div>

          <div className="flex flex-col items-center gap-6 py-12">
            <CompanionHero name="Byte" species="capybara" size="lg" />
            <SpeechBubble speakerName="Byte" className="max-w-xs text-center">
              Welcome back! Ready to sharpen your skills today?
            </SpeechBubble>
          </div>

          <div className="grid grid-cols-3 gap-4">
            {socialProof.map((s) => (
              <div key={s.label}>
                <p className="font-body text-sm font-semibold leading-snug lg:text-base" style={{ color: "#E8EAED" }}>
                  {s.value}
                </p>
                <p className="font-mono mt-1 text-[10px] font-semibold uppercase tracking-widest" style={{ color: "#7C6EF5" }}>
                  {s.label}
                </p>
              </div>
            ))}
          </div>
        </div>

        {/* Right login form column */}
        <div className="flex w-full lg:w-1/2 flex-col items-center justify-center px-6 py-12">
          <div className="mb-8 flex flex-col items-center gap-2 lg:hidden">
            <div className="flex items-center gap-2">
              <Image src="/logo.png" alt="Prepio Logo" width={32} height={32} className="rounded-lg" />
              <p className="font-display text-2xl font-bold" style={{ color: "#7C6EF5" }}>
                PREPIO
              </p>
            </div>
            <CompanionHero name="Byte" species="capybara" size="md" />
          </div>

          <div className="w-full max-w-sm">
            <h1 className="font-display text-2xl font-bold mb-1" style={{ color: "#E8EAED" }}>
              Sign in to Prepio
            </h1>
            <p className="font-body text-sm mb-6" style={{ color: "#8B92A8" }}>
              Enter your email or username and password
            </p>

            {error && (
              <div
                className="rounded-xl px-4 py-3 text-sm font-medium mb-4 space-y-1"
                role="alert"
                style={{
                  background: "rgba(248,113,113,0.1)",
                  border: "1px solid rgba(248,113,113,0.3)",
                  color: "#F87171",
                }}
              >
                <p>{error}</p>
                {isNotFound && (
                  <p className="text-xs text-amber-300">
                    Don&apos;t have an account yet?{" "}
                    <Link href="/register" className="font-bold underline hover:text-white">
                      Sign up here with verification code
                    </Link>
                  </p>
                )}
              </div>
            )}

            <form onSubmit={handlePasswordLogin} className="space-y-4">
              <div>
                <label className="font-mono block text-xs font-semibold mb-1" style={{ color: "#8C93A8" }}>
                  Email or Username
                </label>
                <input
                  type="text"
                  autoComplete="username"
                  required
                  placeholder="you@domain.com or username"
                  value={identifier}
                  onChange={(e) => setIdentifier(e.target.value)}
                  className="w-full rounded-xl px-4 py-3 text-sm outline-none transition focus:border-[#7C6EF5]"
                  style={inputStyle}
                />
              </div>

              <div>
                <label className="font-mono block text-xs font-semibold mb-1" style={{ color: "#8C93A8" }}>
                  Password
                </label>
                <input
                  type="password"
                  autoComplete="current-password"
                  required
                  placeholder="••••••••"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="w-full rounded-xl px-4 py-3 text-sm outline-none transition focus:border-[#7C6EF5]"
                  style={inputStyle}
                />
              </div>

              <div className="pt-2">
                <GameButton type="submit" disabled={loading || !identifier.trim() || !password} className="w-full">
                  {loading ? "Signing in…" : "Sign In"}
                </GameButton>
              </div>
            </form>

            <div className="mt-6 text-center text-sm" style={{ color: "#8B92A8" }}>
              <span>New to Prepio? </span>
              <Link href="/register" className="font-semibold text-[#7C6EF5] hover:underline">
                Create account with OTP →
              </Link>
            </div>
          </div>
        </div>
      </div>
    </GameBackground>
  );
}
