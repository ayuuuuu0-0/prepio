"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { FormEvent, useCallback, useEffect, useMemo, useState } from "react";
import { api, DailyPaper, Question, SubmitResponse } from "@/lib/api";
import { MIN_ANSWER_LENGTH, roundTypeColors } from "@/lib/design/tokens";
import { GameBackground } from "@/components/game/GameBackground";
import { CompanionHero } from "@/components/game/CompanionHero";
import { SpeechBubble } from "@/components/game/SpeechBubble";
import { GameButton } from "@/components/game/GameButton";
import { BottomNav } from "@/components/game/BottomNav";

export default function ChallengePage() {
  const router = useRouter();
  const [paper, setPaper] = useState<DailyPaper | null>(null);
  const [answeredIds, setAnsweredIds] = useState<Set<string>>(new Set());
  const [currentIndex, setCurrentIndex] = useState(0);
  const [answer, setAnswer] = useState("");
  const [result, setResult] = useState<SubmitResponse | null>(null);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [companion, setCompanion] = useState<{ name?: string; species?: string }>({});

  const trimmedLength = answer.trim().length;
  const canSubmit = trimmedLength >= MIN_ANSWER_LENGTH;
  const question = paper?.questions[currentIndex] ?? null;

  const load = useCallback(async () => {
    try {
      const [daily, profile] = await Promise.all([api.getDailyPaper(), api.getProfile()]);
      setPaper(daily);
      setCompanion({ name: profile.companion?.name, species: profile.companion?.species });

      const history = await api.getQuestionHistory(daily.session_id);
      const done = new Set(history.map((h) => h.question_id));
      setAnsweredIds(done);

      const firstOpen = daily.questions.findIndex((q) => !done.has(q.id));
      setCurrentIndex(firstOpen >= 0 ? firstOpen : 0);
    } catch (err) {
      setError(err instanceof Error ? err.message : "failed to load challenge");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    (async () => {
      const ok = await api.ensureSession();
      if (!ok) {
        router.replace("/login");
        return;
      }
      load();
    })();
  }, [router, load]);

  const hasNext = useMemo(() => {
    if (!paper) return false;
    return paper.questions.some((q, i) => i > currentIndex && !answeredIds.has(q.id));
  }, [paper, currentIndex, answeredIds]);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    if (!paper || !question) return;
    if (!canSubmit) {
      setError(`Write at least ${MIN_ANSWER_LENGTH} characters before submitting.`);
      return;
    }
    setSubmitting(true);
    setError("");
    try {
      const res = await api.submitAnswer(question.id, paper.session_id, answer.trim());
      setResult(res);
      setAnswer("");
      setAnsweredIds((prev) => new Set(prev).add(question.id));
    } catch (err) {
      const message = err instanceof Error ? err.message : "submit failed";
      setError(message.includes("already submitted") ? "You already answered this question today." : message);
    } finally {
      setSubmitting(false);
    }
  }

  function handleNext() {
    if (!paper) return;
    setResult(null);
    const next = paper.questions.findIndex((q, i) => i > currentIndex && !answeredIds.has(q.id));
    if (next >= 0) setCurrentIndex(next);
  }

  if (loading) {
    return (
      <GameBackground variant="challenge">
        <main className="flex min-h-screen items-center justify-center">
          <p className="font-mono animate-pulse text-sm font-semibold" style={{ color: "#7C6EF5" }}>
            Loading challenge...
          </p>
        </main>
      </GameBackground>
    );
  }

  return (
    <GameBackground variant="challenge">
      <div
        className="sticky top-0 z-20 px-4 py-3"
        style={{
          background: "rgba(15,17,23,0.9)",
          borderBottom: "1px solid #2E3347",
          backdropFilter: "blur(20px)",
        }}
      >
        <div className="mx-auto flex max-w-4xl items-center justify-between">
          <Link
            href="/dashboard"
            className="font-mono text-sm font-semibold transition-colors"
            style={{ color: "#4A5068" }}
          >
            ← Home
          </Link>
          <div className="flex items-center gap-3">
            <span
              className="font-mono rounded-lg px-3 py-1 text-xs font-bold"
              style={{
                background: "#1A1D27",
                border: "1px solid #2E3347",
                color: "#8B92A8",
              }}
            >
              Q {currentIndex + 1}/{paper?.questions.length ?? 1}
            </span>
            <CompanionHero
              name={companion.name}
              species={companion.species}
              size="sm"
              reaction={result ? (result.correct ? "correct" : "wrong") : "idle"}
            />
          </div>
        </div>
      </div>

      <main className="mx-auto grid max-w-4xl grid-cols-1 gap-6 px-4 pb-28 pt-6 md:grid-cols-2">
        <div className="space-y-4">
          <SpeechBubble speakerName={companion.name ?? "Byte"}>
            {getQuestionHint(question)}
          </SpeechBubble>

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

          {question && (
            <div
              className="rounded-2xl p-6"
              style={{
                background: "#1A1D27",
                border: "1px solid #2E3347",
                borderLeft: `3px solid ${roundTypeColors[question.round_type] ?? "#7C6EF5"}`,
              }}
            >
              <div className="mb-4 flex flex-wrap gap-2">
                <span
                  className="font-mono rounded-lg px-2.5 py-1 text-[11px] font-bold uppercase tracking-wider"
                  style={{
                    background: `${roundTypeColors[question.round_type] ?? "#7C6EF5"}18`,
                    color: roundTypeColors[question.round_type] ?? "#7C6EF5",
                  }}
                >
                  {question.round_type.replace("_", " ")}
                </span>
                <span
                  className="font-mono rounded-lg px-2.5 py-1 text-[11px] font-bold uppercase tracking-wider"
                  style={{
                    background:
                      question.difficulty === "hard"
                        ? "rgba(248,113,113,0.12)"
                        : question.difficulty === "medium"
                          ? "rgba(245,185,66,0.12)"
                          : "rgba(52,211,153,0.12)",
                    color:
                      question.difficulty === "hard"
                        ? "#F87171"
                        : question.difficulty === "medium"
                          ? "#F5B942"
                          : "#34D399",
                  }}
                >
                  {question.difficulty}
                </span>
              </div>
              <p className="text-base font-medium leading-relaxed" style={{ color: "#C8CCDA" }}>
                {question.body}
              </p>
            </div>
          )}
        </div>

        <div className="space-y-4">
          {result ? (
            <ResultCard result={result} onNext={handleNext} hasNext={hasNext} />
          ) : (
            <form onSubmit={onSubmit} className="space-y-3">
              <textarea
                className="w-full resize-none rounded-2xl p-4 text-sm font-medium leading-relaxed outline-none transition-all"
                style={{
                  background: "#1A1D27",
                  border: "1px solid #2E3347",
                  color: "#E8EAED",
                  minHeight: "220px",
                }}
                onFocus={(e) => (e.target.style.borderColor = "#7C6EF5")}
                onBlur={(e) => (e.target.style.borderColor = "#2E3347")}
                placeholder="Explain your approach..."
                value={answer}
                onChange={(e) => setAnswer(e.target.value)}
              />
              <p
                className="font-mono text-xs"
                style={{ color: canSubmit ? "#34D399" : "#4A5068" }}
              >
                {canSubmit
                  ? "Ready to submit"
                  : `${trimmedLength}/${MIN_ANSWER_LENGTH} chars — add more detail`}
              </p>
              <GameButton type="submit" disabled={submitting || !canSubmit}>
                {submitting ? "Evaluating..." : "Submit"}
              </GameButton>
            </form>
          )}
        </div>
      </main>
      <BottomNav />
    </GameBackground>
  );
}

function ResultCard({
  result,
  onNext,
  hasNext,
}: {
  result: SubmitResponse;
  onNext: () => void;
  hasNext: boolean;
}) {
  return (
    <div
      className="animate-result space-y-4 rounded-2xl p-6"
      style={{
        background: result.correct ? "rgba(52,211,153,0.08)" : "rgba(245,185,66,0.08)",
        border: `1px solid ${result.correct ? "rgba(52,211,153,0.3)" : "rgba(245,185,66,0.3)"}`,
      }}
    >
      <div className="flex items-center gap-3">
        <span className="text-3xl">{result.correct ? "✓" : "→"}</span>
        <div>
          <p
            className="font-display text-lg font-bold"
            style={{ color: result.correct ? "#34D399" : "#F5B942" }}
          >
            {result.correct ? "Solid." : "Not quite."}
          </p>
          <p className="text-sm" style={{ color: "#8B92A8" }}>
            {result.feedback}
          </p>
          <p className="font-mono mt-1 text-xs" style={{ color: "#4A5068" }}>
            Score: {result.score}%
          </p>
        </div>
      </div>

      {result.strengths.length > 0 && (
        <div>
          <p className="font-mono text-[10px] font-bold uppercase tracking-widest" style={{ color: "#4A5068" }}>
            Covered
          </p>
          <ul className="mt-2 space-y-1">
            {result.strengths.map((s) => (
              <li key={s} className="text-sm" style={{ color: "#C8CCDA" }}>
                ✓ {s}
              </li>
            ))}
          </ul>
        </div>
      )}

      {result.gaps.length > 0 && (
        <div>
          <p className="font-mono text-[10px] font-bold uppercase tracking-widest" style={{ color: "#4A5068" }}>
            Missing
          </p>
          <ul className="mt-2 space-y-1">
            {result.gaps.map((g) => (
              <li key={g} className="text-sm" style={{ color: "#C8CCDA" }}>
                → {g}
              </li>
            ))}
          </ul>
        </div>
      )}

      {result.correct && (
        <div className="flex gap-3">
          <div
            className="animate-xp flex items-center gap-2 rounded-xl px-3 py-2"
            style={{ background: "rgba(96,165,250,0.12)", border: "1px solid rgba(96,165,250,0.2)" }}
          >
            <span className="font-mono text-sm font-bold" style={{ color: "#60A5FA" }}>
              +{result.xp_awarded} XP
            </span>
          </div>
          <div
            className="animate-gems flex items-center gap-2 rounded-xl px-3 py-2"
            style={{ background: "rgba(52,211,153,0.12)", border: "1px solid rgba(52,211,153,0.2)" }}
          >
            <span className="font-mono text-sm font-bold" style={{ color: "#34D399" }}>
              +{result.gems_awarded} 💎
            </span>
          </div>
        </div>
      )}

      {result.readiness_delta > 0 && (
        <p className="font-mono text-xs font-semibold" style={{ color: "#60A5FA" }}>
          Readiness +{result.readiness_delta}%
        </p>
      )}

      {hasNext ? (
        <GameButton type="button" onClick={onNext}>
          Next Question →
        </GameButton>
      ) : (
        <Link href="/dashboard">
          <GameButton type="button" variant="secondary">
            Back to dashboard
          </GameButton>
        </Link>
      )}
    </div>
  );
}

/** Returns a question-specific hint based on round type. */
function getQuestionHint(question: Question | null): string {
  if (!question) return "Loading your challenge...";
  const hints: Record<string, string> = {
    dsa: "Think about time and space complexity before you write.",
    system_design: "Start with requirements, then components, then trade-offs.",
    lld: "Focus on class design, responsibilities, and relationships.",
    behavioral: "Use the STAR format: Situation, Task, Action, Result.",
    aptitude: "Break the problem into smaller steps first.",
    fundamentals: "Cover the core concept clearly before diving into detail.",
  };
  return hints[question.round_type] ?? "Take your time. Think it through.";
}
