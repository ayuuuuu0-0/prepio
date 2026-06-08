"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { FormEvent, useCallback, useEffect, useState } from "react";
import { api, DailyPaper, Question, SubmitResponse } from "@/lib/api";
import { MIN_ANSWER_LENGTH } from "@/lib/design/tokens";
import { GameBackground } from "@/components/game/GameBackground";
import { CompanionHero } from "@/components/game/CompanionHero";
import { SpeechBubble } from "@/components/game/SpeechBubble";
import { GameButton } from "@/components/game/GameButton";
import { BottomNav } from "@/components/game/BottomNav";

export default function ChallengePage() {
  const router = useRouter();
  const [paper, setPaper] = useState<DailyPaper | null>(null);
  const [question, setQuestion] = useState<Question | null>(null);
  const [answer, setAnswer] = useState("");
  const [result, setResult] = useState<SubmitResponse | null>(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);

  const trimmedLength = answer.trim().length;
  const canSubmit = trimmedLength >= MIN_ANSWER_LENGTH;

  const load = useCallback(async () => {
    try {
      const daily = await api.getDailyPaper();
      setPaper(daily);
      if (daily.questions.length > 0) setQuestion(daily.questions[0]);
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to load challenge");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    if (!api.loadToken()) {
      router.replace("/login");
      return;
    }
    load();
  }, [router, load]);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    api.loadToken();
    if (!paper || !question) return;
    if (!canSubmit) {
      setError(`Write at least ${MIN_ANSWER_LENGTH} characters before submitting.`);
      return;
    }
    setSubmitting(true);
    setError("");
    setResult(null);
    try {
      const res = await api.submitAnswer(question.id, paper.session_id, answer.trim());
      setResult(res);
      setAnswer("");
    } catch (err) {
      const message = err instanceof Error ? err.message : "submit failed";
      if (message.includes("already submitted")) {
        setError("You already answered this question today. Come back tomorrow!");
      } else {
        setError(message);
      }
    } finally {
      setSubmitting(false);
    }
  }

  if (loading) {
    return (
      <GameBackground variant="challenge">
        <main className="flex min-h-screen items-center justify-center">
          <p className="font-display animate-pulse text-lg font-bold text-[#1CB0F6]">Preparing challenge...</p>
        </main>
      </GameBackground>
    );
  }

  return (
    <GameBackground variant="challenge">
      <main className="mx-auto max-w-lg px-4 pb-28 pt-4">
        <div className="flex items-center justify-between">
          <Link href="/dashboard" className="font-display text-sm font-bold text-[#1CB0F6]">
            ← Home
          </Link>
          <div className="flex items-center gap-2">
            <span className="font-display rounded-full bg-white px-3 py-1 text-sm font-bold shadow">
              Q 1/5
            </span>
            <CompanionHero size="sm" />
          </div>
        </div>

        {error && (
          <p className="mt-4 rounded-2xl bg-orange-100 px-4 py-3 text-sm font-semibold text-orange-700">{error}</p>
        )}

        {result && (
          <div
            className={`animate-confetti mt-4 rounded-3xl p-6 text-center shadow-lg ${
              result.correct
                ? "bg-gradient-to-br from-[#58CC02] to-[#7CB342] text-white"
                : "bg-gradient-to-br from-[#FFB84D] to-[#FF9600] text-white"
            }`}
          >
            <p className="text-4xl">{result.correct ? "🎉" : "💪"}</p>
            <p className="font-display mt-2 text-2xl font-extrabold">
              {result.correct ? "Amazing!" : "Almost there!"}
            </p>
            <p className="mt-2 font-semibold">{result.feedback}</p>
            {result.correct && (
              <div className="mt-4 flex justify-center gap-4">
                <span className="font-display rounded-full bg-white/25 px-4 py-2 font-bold">⚡ +{result.xp_awarded || "XP"}</span>
                <span className="font-display rounded-full bg-white/25 px-4 py-2 font-bold">💎 +{result.gems_awarded || "Gems"}</span>
              </div>
            )}
            <Link href="/dashboard" className="mt-4 inline-block font-display font-bold underline">
              Back to dashboard
            </Link>
          </div>
        )}

        {question && paper && !result && (
          <section className="mt-6">
            <SpeechBubble className="mb-4">You&apos;ve got this! Take your time and think it through.</SpeechBubble>

            <div className="rounded-3xl bg-white p-6 shadow-lg">
              <div className="flex gap-2">
                <span className="font-display rounded-full bg-[#1CB0F6]/20 px-3 py-1 text-xs font-bold uppercase text-[#1CB0F6]">
                  {question.difficulty}
                </span>
                <span className="font-display rounded-full bg-[#7B5CFF]/20 px-3 py-1 text-xs font-bold uppercase text-[#7B5CFF]">
                  {question.round_type}
                </span>
              </div>
              <p className="mt-4 text-base font-semibold leading-relaxed text-[#3C3C3C]">{question.body}</p>

              <form onSubmit={onSubmit} className="mt-6">
                <textarea
                  className="w-full rounded-2xl border-2 border-[#E5E5E5] bg-[#FAFAFA] p-4 text-sm font-medium outline-none focus:border-[#58CC02]"
                  rows={6}
                  placeholder="Write your answer..."
                  value={answer}
                  onChange={(e) => setAnswer(e.target.value)}
                />
                <p className={`mt-2 text-xs font-semibold ${canSubmit ? "text-[#58CC02]" : "text-[#777]"}`}>
                  {canSubmit
                    ? "Ready to submit!"
                    : `${trimmedLength}/${MIN_ANSWER_LENGTH} characters — keep writing...`}
                </p>
                <div className="mt-4">
                  <GameButton type="submit" disabled={submitting || !canSubmit}>
                    {submitting ? "Checking..." : "Submit Answer!"}
                  </GameButton>
                </div>
              </form>
            </div>
          </section>
        )}
      </main>
      <BottomNav />
    </GameBackground>
  );
}
