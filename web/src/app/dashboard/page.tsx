"use client";

import { useRouter } from "next/navigation";
import { FormEvent, useCallback, useEffect, useState } from "react";
import {
  api,
  DailyPaper,
  Progress,
  Question,
  Streak,
  SubmitResponse,
} from "@/lib/api";

export default function DashboardPage() {
  const router = useRouter();
  const [paper, setPaper] = useState<DailyPaper | null>(null);
  const [streak, setStreak] = useState<Streak | null>(null);
  const [progress, setProgress] = useState<Progress | null>(null);
  const [activeQuestion, setActiveQuestion] = useState<Question | null>(null);
  const [answer, setAnswer] = useState("");
  const [submitResult, setSubmitResult] = useState<SubmitResponse | null>(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);

  const refresh = useCallback(async () => {
    setError("");
    try {
      const [daily, streakData, progressData] = await Promise.all([
        api.getDailyPaper(),
        api.getStreak(),
        api.getProgress(),
      ]);
      setPaper(daily);
      setStreak(streakData);
      setProgress(progressData);
      if (!activeQuestion && daily.questions.length > 0) {
        setActiveQuestion(daily.questions[0]);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to load data");
    } finally {
      setLoading(false);
    }
  }, [activeQuestion]);

  useEffect(() => {
    if (!api.loadToken()) {
      router.replace("/login");
      return;
    }
    refresh();
  }, [router, refresh]);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    if (!paper || !activeQuestion || len(answer) === 0) return;
    setSubmitting(true);
    setSubmitResult(null);
    try {
      const result = await api.submitAnswer(activeQuestion.id, paper.session_id, answer);
      setSubmitResult(result);
      setAnswer("");
      await new Promise((r) => setTimeout(r, 800));
      await refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : "submit failed");
    } finally {
      setSubmitting(false);
    }
  }

  function len(s: string) {
    return s.length;
  }

  if (loading) {
    return (
      <main className="flex min-h-screen items-center justify-center">
        <p className="text-slate-500">Loading your daily paper...</p>
      </main>
    );
  }

  return (
    <main className="min-h-screen bg-slate-50">
      <header className="border-b border-slate-200 bg-white">
        <div className="mx-auto flex max-w-5xl items-center justify-between px-4 py-4">
          <h1 className="text-xl font-bold text-emerald-700">Prepio</h1>
          <button
            onClick={() => {
              api.setToken(null);
              router.push("/login");
            }}
            className="text-sm text-slate-500 hover:text-slate-800"
          >
            Sign out
          </button>
        </div>
      </header>

      <div className="mx-auto grid max-w-5xl gap-6 px-4 py-8 md:grid-cols-3">
        <section className="rounded-2xl bg-white p-6 shadow-sm md:col-span-2">
          <h2 className="text-lg font-semibold text-slate-900">Today&apos;s question</h2>
          {error && <p className="mt-3 text-sm text-red-600">{error}</p>}

          {activeQuestion && paper && (
            <>
              <div className="mt-4 flex gap-2 text-xs">
                <span className="rounded-full bg-slate-100 px-2 py-1 uppercase">{activeQuestion.difficulty}</span>
                <span className="rounded-full bg-slate-100 px-2 py-1 uppercase">{activeQuestion.round_type}</span>
                {activeQuestion.company_tags.map((tag) => (
                  <span key={tag} className="rounded-full bg-emerald-50 px-2 py-1 text-emerald-700">
                    {tag}
                  </span>
                ))}
              </div>
              <p className="mt-4 text-slate-800 leading-relaxed">{activeQuestion.body}</p>

              <form onSubmit={onSubmit} className="mt-6">
                <textarea
                  className="w-full rounded-xl border border-slate-200 p-3 text-sm"
                  rows={5}
                  placeholder="Write your answer..."
                  value={answer}
                  onChange={(e) => setAnswer(e.target.value)}
                />
                <button
                  type="submit"
                  disabled={submitting || len(answer) < 10}
                  className="mt-3 rounded-lg bg-emerald-600 px-4 py-2 text-sm font-semibold text-white hover:bg-emerald-700 disabled:opacity-50"
                >
                  {submitting ? "Submitting..." : "Submit answer"}
                </button>
              </form>

              {submitResult && (
                <div
                  className={`mt-4 rounded-xl p-4 text-sm ${
                    submitResult.correct ? "bg-emerald-50 text-emerald-800" : "bg-amber-50 text-amber-800"
                  }`}
                >
                  <p className="font-medium">{submitResult.correct ? "Correct!" : "Keep practicing"}</p>
                  <p className="mt-1">{submitResult.feedback}</p>
                  <p className="mt-2 text-xs opacity-75">
                    XP and gems update shortly via the event pipeline.
                  </p>
                </div>
              )}
            </>
          )}
        </section>

        <aside className="space-y-4">
          <div className="rounded-2xl bg-white p-6 shadow-sm">
            <h3 className="font-semibold text-slate-900">Streak</h3>
            {streak && (
              <div className="mt-3 space-y-2 text-sm text-slate-600">
                <p>
                  <span className="text-3xl font-bold text-orange-500">{streak.current_streak}</span> days
                </p>
                <p>Longest: {streak.longest_streak}</p>
                <p>Freezes: {streak.freeze_count}</p>
                <p>{streak.streak_active_today ? "Active today" : "Not active yet today"}</p>
              </div>
            )}
          </div>

          <div className="rounded-2xl bg-white p-6 shadow-sm">
            <h3 className="font-semibold text-slate-900">Progress</h3>
            {progress && (
              <div className="mt-3 space-y-2 text-sm text-slate-600">
                <p>
                  Level <span className="text-xl font-bold text-emerald-600">{progress.current_level}</span>
                </p>
                <p>{progress.total_xp} XP</p>
                <p>{progress.gem_balance} gems</p>
                <p>{progress.xp_to_next_level} XP to next level</p>
              </div>
            )}
          </div>
        </aside>
      </div>
    </main>
  );
}
