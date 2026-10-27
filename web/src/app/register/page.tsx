"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import Image from "next/image";
import { FormEvent, useEffect, useState } from "react";
import { GameBackground } from "@/components/game/GameBackground";
import { CompanionHero } from "@/components/game/CompanionHero";
import { SpeechBubble } from "@/components/game/SpeechBubble";
import { GameButton } from "@/components/game/GameButton";
import { api, ApiRequestError } from "@/lib/api";

const inputStyle = {
  background: "#1A1D27",
  border: "1px solid #2E3347",
  color: "#E8EAED",
};

export default function RegisterPage() {
  const router = useRouter();
  const [step, setStep] = useState<"details" | "otp">("details");
  const [email, setEmail] = useState("");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [code, setCode] = useState("");
  const [countdown, setCountdown] = useState(0);
  const [devCode, setDevCode] = useState("");
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
    const cleanEmail = email.trim();
    const cleanUsername = username.trim();

    if (!cleanEmail || !cleanEmail.includes("@") || !cleanEmail.includes(".")) {
      setError("Please enter a valid email address.");
      return;
    }
    if (cleanUsername.length < 3) {
      setError("Username must be at least 3 characters long.");
      return;
    }
    if (password.length < 8) {
      setError("Password must be at least 8 characters long.");
      return;
    }

    setLoading(true);
    setError("");
    setMessage("");
    try {
      const res = await api.sendOTP(cleanEmail, cleanUsername, true);
      setMessage(res.message || "Verification code sent to your email!");
      if (res.dev_code) {
        setDevCode(res.dev_code);
      }
      setStep("otp");
      setCountdown(60);
    } catch (err) {
      if (err instanceof ApiRequestError) {
        setError(err.message);
      } else {
        setError(err instanceof Error ? err.message : "Failed to send verification code");
      }
    } finally {
      setLoading(false);
    }
  }

  async function handleVerifyAndRegister(e: FormEvent) {
    e.preventDefault();
    const cleanCode = code.trim();
    if (!cleanCode || cleanCode.length !== 6) {
      setError("Please enter the complete 6-digit verification code");
      return;
    }
    setLoading(true);
    setError("");
    try {
      const res = await api.verifyOTP(email.trim(), cleanCode, username.trim(), password);
      api.setAuthTokens(res.access_token);
      router.push("/onboarding");
    } catch (err) {
      if (err instanceof ApiRequestError) {
        setError(err.message);
      } else {
        setError(err instanceof Error ? err.message : "Verification failed");
      }
    } finally {
      setLoading(false);
    }
  }

  return (
    <GameBackground>
      <div className="flex min-h-screen">
        {/* Left marketing panel on desktop */}
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
              Create your account to unlock your personal companion and level up!
            </SpeechBubble>
          </div>

          <div className="text-xs font-mono" style={{ color: "#8C93A8" }}>
            Real-world system design · Multi-language coding · Gamified mastery
          </div>
        </div>

        {/* Right form panel */}
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
              {step === "details" ? "Create your account" : "Verify your email"}
            </h1>
            <p className="font-body text-sm mb-6" style={{ color: "#8B92A8" }}>
              {step === "details"
                ? "Enter your details to receive your 6-digit signup OTP"
                : `We sent a 6-digit verification code to ${email}`}
            </p>

            {error && (
              <p
                className="rounded-xl px-4 py-3 text-sm font-medium mb-4"
                role="alert"
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
                role="status"
                style={{
                  background: "rgba(52,211,153,0.1)",
                  border: "1px solid rgba(52,211,153,0.3)",
                  color: "#34D399",
                }}
              >
                {message}
              </p>
            )}

            {/* Step 1: User enters email, username, password */}
            {step === "details" && (
              <form onSubmit={handleSendOTP} className="space-y-4">
                <div>
                  <label className="font-mono block text-xs font-semibold mb-1" style={{ color: "#8C93A8" }}>
                    Email
                  </label>
                  <input
                    type="email"
                    autoComplete="email"
                    required
                    placeholder="you@domain.com"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    className="w-full rounded-xl px-4 py-3 text-sm outline-none transition focus:border-[#7C6EF5]"
                    style={inputStyle}
                  />
                </div>

                <div>
                  <label className="font-mono block text-xs font-semibold mb-1" style={{ color: "#8C93A8" }}>
                    Username
                  </label>
                  <input
                    type="text"
                    autoComplete="username"
                    required
                    placeholder="e.g. ayush_coder"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    className="w-full rounded-xl px-4 py-3 text-sm outline-none transition focus:border-[#7C6EF5]"
                    style={inputStyle}
                  />
                </div>

                <div>
                  <label className="font-mono block text-xs font-semibold mb-1" style={{ color: "#8C93A8" }}>
                    Password (min. 8 characters)
                  </label>
                  <input
                    type="password"
                    autoComplete="new-password"
                    required
                    minLength={8}
                    placeholder="••••••••"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className="w-full rounded-xl px-4 py-3 text-sm outline-none transition focus:border-[#7C6EF5]"
                    style={inputStyle}
                  />
                </div>

                <div className="pt-2">
                  <GameButton
                    type="submit"
                    disabled={loading || !email.trim() || !username.trim() || password.length < 8}
                    className="w-full"
                  >
                    {loading ? "Sending verification code…" : "Send Verification OTP →"}
                  </GameButton>
                </div>
              </form>
            )}

            {/* Step 2: User enters OTP */}
            {step === "otp" && (
              <form onSubmit={handleVerifyAndRegister} className="space-y-4">
                {devCode && (
                  <div className="flex items-center justify-between rounded-xl p-3 bg-amber-500/10 border border-amber-500/30 text-xs">
                    <span className="text-amber-300 font-mono">Dev OTP: <strong>{devCode}</strong></span>
                    <button
                      type="button"
                      onClick={() => setCode(devCode)}
                      className="px-2.5 py-1 bg-amber-500 text-black font-bold rounded-lg hover:bg-amber-400"
                    >
                      Fill Code
                    </button>
                  </div>
                )}

                <div>
                  <label className="font-mono block text-xs font-semibold mb-1" style={{ color: "#8C93A8" }}>
                    Enter 6-Digit Code
                  </label>
                  <input
                    type="text"
                    inputMode="numeric"
                    pattern="[0-9]*"
                    maxLength={6}
                    autoFocus
                    required
                    placeholder="123456"
                    value={code}
                    onChange={(e) => setCode(e.target.value.replace(/\D/g, "").slice(0, 6))}
                    className="w-full text-center tracking-[0.4em] font-mono text-xl font-bold rounded-xl px-4 py-3 outline-none transition focus:border-[#7C6EF5]"
                    style={inputStyle}
                  />
                </div>

                <div className="pt-2">
                  <GameButton
                    type="submit"
                    disabled={loading || code.trim().length !== 6}
                    className="w-full"
                  >
                    {loading ? "Verifying & Creating Account…" : "Verify & Start Journey"}
                  </GameButton>
                </div>

                <div className="flex items-center justify-between text-xs pt-2" style={{ color: "#8B92A8" }}>
                  <button
                    type="button"
                    onClick={() => setStep("details")}
                    className="hover:underline text-[#7C6EF5]"
                  >
                    ← Edit details
                  </button>
                  <button
                    type="button"
                    disabled={countdown > 0 || loading}
                    onClick={() => handleSendOTP()}
                    className="hover:underline text-[#7C6EF5] disabled:opacity-40 disabled:hover:no-underline"
                  >
                    {countdown > 0 ? `Resend in ${countdown}s` : "Resend code"}
                  </button>
                </div>
              </form>
            )}

            <div className="mt-6 text-center text-sm" style={{ color: "#8B92A8" }}>
              <span>Already have an account? </span>
              <Link href="/login" className="font-semibold text-[#7C6EF5] hover:underline">
                Sign in with password →
              </Link>
            </div>
          </div>
        </div>
      </div>
    </GameBackground>
  );
}
