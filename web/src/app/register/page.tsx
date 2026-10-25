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

const inputStyle = {
  background: "#1A1D27",
  border: "1px solid #2E3347",
  color: "#E8EAED",
};

export default function RegisterPage() {
  const router = useRouter();
  const [method, setMethod] = useState<"otp" | "password">("otp");
  const [otpStep, setOtpStep] = useState<"send" | "verify">("send");
  const [email, setEmail] = useState("");
  const [username, setUsername] = useState("");
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
      const res = await api.verifyOTP(email, code, username.trim() || undefined);
      api.setAuthTokens(res.access_token);
      router.push("/onboarding");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Verification failed");
    } finally {
      setLoading(false);
    }
  }

  async function handlePasswordRegister(e: FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");
    try {
      const res = await api.register(email, username, password);
      api.setAuthTokens(res.access_token);
      router.push("/onboarding");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Registration failed");
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
            <CompanionHero name="Pip" species="red_panda" size="lg" />
            <SpeechBubble speakerName="Pip" className="max-w-xs text-center">
              Pick a companion. Start mastering engineering topics today.
            </SpeechBubble>
          </div>

          <div className="grid grid-cols-3 gap-4">
            {[
              { label: "Core Topics", value: "System Design, DSA, LLD" },
              { label: "Lessons captured", value: "Misses, patterns, next steps" },
              { label: "Mastery focus", value: "Skills over streaks" },
            ].map((s) => (
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
            <CompanionHero name="Pip" species="red_panda" size="md" />
          </div>

          <div className="w-full max-w-sm">
            <h1 className="font-display text-2xl font-bold mb-1" style={{ color: "#E8EAED" }}>
              Create account
            </h1>
            <p className="font-body text-sm mb-6" style={{ color: "#4A5068" }}>
              Pick a companion. Start mastering topics from day one.
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
                Instant Email OTP
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
                With Password
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
                  <div>
                    <label className="block text-xs font-medium mb-1.5" style={{ color: "#8C93A8" }}>
                      Username <span style={{ color: "#4A5068" }}>(optional)</span>
                    </label>
                    <input
                      className="w-full rounded-xl px-4 py-3 text-sm font-medium outline-none transition-all"
                      style={inputStyle}
                      onFocus={(e) => (e.target.style.borderColor = "#7C6EF5")}
                      onBlur={(e) => (e.target.style.borderColor = "#2E3347")}
                      placeholder="e.g. devhero"
                      value={username}
                      onChange={(e) => setUsername(e.target.value)}
                    />
                  </div>
                  <GameButton type="submit" variant="gold" disabled={loading}>
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

                  <GameButton type="submit" variant="gold" disabled={loading || code.length !== 6}>
                    {loading ? "Verifying..." : "Verify & Create Account"}
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
              <form onSubmit={handlePasswordRegister} className="space-y-4">
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
            )}

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
