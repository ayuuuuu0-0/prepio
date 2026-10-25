"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import Image from "next/image";
import { FormEvent, useEffect, useState } from "react";
import { GameBackground } from "@/components/game/GameBackground";
import { CompanionHero } from "@/components/game/CompanionHero";
import { SpeechBubble } from "@/components/game/SpeechBubble";
import { GameButton } from "@/components/game/GameButton";
import { api } from "@/lib/api";

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
  const [method, setMethod] = useState<"otp" | "password">("otp");
  const [otpStep, setOtpStep] = useState<"send" | "verify">("send");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [code, setCode] = useState("");
  const [countdown, setCountdown] = useState(0);
  const [error, setError] = useState("");
  const [message, setMessage] = useState("");
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (countdown <= 0) return;
    const timer = setInterval(() => {
      setCountdown((prev) => (prev > 0 ? prev - 1 : 0));
    }, 1000);
    return () => clearInterval(timer);
  }, [countdown]);

  async function handleSendOTP(e?: FormEvent) {
    if (e) e.preventDefault();
    if (!email) return;
    setLoading(true);
    setError("");
    setMessage("");
    try {
      const res = await api.sendOTP(email);
      setMessage(res.message || "Verification code sent to your email!");
      setOtpStep("verify");
      setCountdown(60);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to send code");
    } finally {
      setLoading(false);
    }
  }

  async function handleVerifyOTP(e: FormEvent) {
    e.preventDefault();
    if (!code || code.length !== 6) {
      setError("Please enter the 6-digit verification code");
      return;
    }
    setLoading(true);
    setError("");
    try {
      const res = await api.verifyOTP(email, code);
      api.setAuthTokens(res.access_token);
      const profile = await api.getProfile();
      router.push(profile.onboarding_completed ? "/dashboard" : "/onboarding");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Verification failed");
    } finally {
      setLoading(false);
    }
  }

  async function handlePasswordLogin(e: FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");
    try {
      const res = await api.login(email, password);
      api.setAuthTokens(res.access_token);
      const profile = await api.getProfile();
      router.push(profile.onboarding_completed ? "/dashboard" : "/onboarding");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Login failed");
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
              Byte is waiting. Your streak won&apos;t hold itself.
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
              Continue prep
            </h1>
            <p className="font-body text-sm mb-6" style={{ color: "#4A5068" }}>
              Sign in to your account
            </p>

            {/* Auth Method Selector */}
            <div
              className="flex rounded-xl p-1 mb-6 text-xs font-semibold"
              style={{ background: "#1A1D27", border: "1px solid #2E3347" }}
            >
              <button
                type="button"
                onClick={() => {
                  setMethod("otp");
                  setError("");
                  setMessage("");
                }}
                className={`flex-1 py-2 rounded-lg transition-all ${
                  method === "otp"
                    ? "bg-[#7C6EF5] text-white shadow-sm"
                    : "text-[#8C93A8] hover:text-[#E8EAED]"
                }`}
              >
                Email Code (OTP)
              </button>
              <button
                type="button"
                onClick={() => {
                  setMethod("password");
                  setError("");
                  setMessage("");
                }}
                className={`flex-1 py-2 rounded-lg transition-all ${
                  method === "password"
                    ? "bg-[#7C6EF5] text-white shadow-sm"
                    : "text-[#8C93A8] hover:text-[#E8EAED]"
                }`}
              >
                Password
              </button>
            </div>

            {error && (
              <p
                className="rounded-xl px-4 py-3 text-sm font-medium mb-4"
                style={{
                  background: "rgba(248,113,113,0.1)",
                  border: "1px solid rgba(248,113,113,0.3)",
                  color: "#F87171",
                }}
              >
                {error}
              </p>
            )}

            {message && (
              <p
                className="rounded-xl px-4 py-3 text-sm font-medium mb-4"
                style={{
                  background: "rgba(52,211,153,0.1)",
                  border: "1px solid rgba(52,211,153,0.3)",
                  color: "#34D399",
                }}
              >
                {message}
              </p>
            )}

            {method === "otp" ? (
              otpStep === "send" ? (
                <form onSubmit={handleSendOTP} className="space-y-4">
                  <div>
                    <label className="block text-xs font-medium mb-1.5" style={{ color: "#8C93A8" }}>
                      Email Address
                    </label>
                    <input
                      className="w-full rounded-xl px-4 py-3 text-sm font-medium outline-none transition-all"
                      style={inputStyle}
                      onFocus={(e) => (e.target.style.borderColor = "#7C6EF5")}
                      onBlur={(e) => (e.target.style.borderColor = "#2E3347")}
                      type="email"
                      placeholder="you@example.com"
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                      required
                    />
                  </div>
                  <GameButton type="submit" disabled={loading}>
                    {loading ? "Sending Code..." : "Send Verification Code"}
                  </GameButton>
                </form>
              ) : (
                <form onSubmit={handleVerifyOTP} className="space-y-4">
                  <div>
                    <div className="flex items-center justify-between mb-1.5">
                      <label className="text-xs font-medium" style={{ color: "#8C93A8" }}>
                        6-Digit Code
                      </label>
                      <button
                        type="button"
                        onClick={() => {
                          setOtpStep("send");
                          setCode("");
                          setError("");
                        }}
                        className="text-xs font-semibold hover:underline"
                        style={{ color: "#7C6EF5" }}
                      >
                        Change email
                      </button>
                    </div>
                    <input
                      className="w-full rounded-xl px-4 py-3 text-center text-xl font-bold tracking-[0.5em] outline-none transition-all"
                      style={inputStyle}
                      onFocus={(e) => (e.target.style.borderColor = "#7C6EF5")}
                      onBlur={(e) => (e.target.style.borderColor = "#2E3347")}
                      type="text"
                      maxLength={6}
                      placeholder="------"
                      value={code}
                      onChange={(e) => setCode(e.target.value.replace(/\D/g, "").slice(0, 6))}
                      required
                      autoFocus
                    />
                    <p className="text-[11px] mt-1 text-center" style={{ color: "#8C93A8" }}>
                      Sent to <span className="font-medium text-[#E8EAED]">{email}</span>
                    </p>
                  </div>

                  <GameButton type="submit" disabled={loading || code.length !== 6}>
                    {loading ? "Verifying..." : "Verify & Sign In"}
                  </GameButton>

                  <div className="pt-2 text-center">
                    {countdown > 0 ? (
                      <p className="text-xs font-mono" style={{ color: "#8C93A8" }}>
                        Resend code in {countdown}s
                      </p>
                    ) : (
                      <button
                        type="button"
                        onClick={() => handleSendOTP()}
                        disabled={loading}
                        className="text-xs font-semibold hover:underline"
                        style={{ color: "#7C6EF5" }}
                      >
                        Resend verification code
                      </button>
                    )}
                  </div>
                </form>
              )
            ) : (
              <form onSubmit={handlePasswordLogin} className="space-y-4">
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
                  type="password"
                  placeholder="Password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                />

                <GameButton type="submit" disabled={loading}>
                  {loading ? "Signing in..." : "Continue Prep"}
                </GameButton>
              </form>
            )}

            <p className="mt-6 text-center text-sm" style={{ color: "#4A5068" }}>
              No account?{" "}
              <Link href="/register" className="font-semibold" style={{ color: "#7C6EF5" }}>
                Start your learning path
              </Link>
            </p>
          </div>
        </div>
      </div>
    </GameBackground>
  );
}
