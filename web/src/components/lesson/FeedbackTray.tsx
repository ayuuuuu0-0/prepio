"use client";

import { useEffect, useRef } from "react";
import type { AnswerResult } from "@/lib/lesson/types";

/**
 * FeedbackTray slides up after a step is graded. It always carries an icon and
 * words (never colour alone), the explanation, and, after a miss, why the chosen
 * option was wrong. Enter (handled by the page) or the button continues.
 */
export function FeedbackTray({
  result,
  correctAnswer = [],
  companionName,
  onContinue,
  continueLabel,
}: {
  result: AnswerResult;
  /** The server's correct answer in words, one line per part (sent only after a miss). */
  correctAnswer?: string[];
  companionName?: string;
  onContinue: () => void;
  continueLabel: string;
}) {
  const ref = useRef<HTMLButtonElement>(null);
  useEffect(() => {
    ref.current?.focus({ preventScroll: true });
  }, []);

  const ok = result.correct;
  // Written answers come back with a rubric score instead of a single correct option.
  const prose = typeof result.score === "number";
  const tone = ok ? "#34D399" : "#F87171";
  const heading = ok ? "Nice, that's right" : "Not quite";
  const intro = ok
    ? companionName
      ? `${companionName} is impressed.`
      : null
    : companionName
      ? `${companionName}: no worries, this one comes back at the end.`
      : "No worries, this one comes back at the end.";

  return (
    <div
      className="animate-tray sticky bottom-0 z-20 w-full"
      style={{
        background: ok ? "#10231C" : "#26151A",
        borderTop: `2px solid ${tone}`,
        boxShadow: "0 -12px 40px rgba(0,0,0,0.45)",
      }}
    >
      <div className="mx-auto w-full max-w-2xl px-4 pb-6 pt-4">
        <div role="status" aria-live="polite" className="max-h-[55vh] overflow-y-auto">
          <p className="font-display flex items-center gap-2 text-lg font-extrabold" style={{ color: tone }}>
            <span
              aria-hidden
              className="flex h-7 w-7 items-center justify-center rounded-full text-sm"
              style={{ background: tone, color: "#0F1117" }}
            >
              {ok ? "✓" : "✕"}
            </span>
            {heading}
          </p>
          {intro && (
            <p className="mt-1 text-sm font-semibold" style={{ color: "#8B92A8" }}>
              {intro}
            </p>
          )}

          {!ok && correctAnswer.length > 0 && (
            <div className="mt-3 text-sm" style={{ color: "#E8EAED" }}>
              <p className="font-mono text-[11px] font-bold uppercase tracking-widest" style={{ color: "#8B92A8" }}>
                Correct answer
              </p>
              {correctAnswer.map((line) => (
                <p key={line} className="font-semibold">
                  {line}
                </p>
              ))}
            </div>
          )}

          {!ok && result.why_not && (
            <p className="mt-3 text-sm leading-relaxed" style={{ color: "#C8CCDA" }}>
              <span className="font-bold">Why not that one: </span>
              {result.why_not}
            </p>
          )}

          {prose && (
            <div className="mt-3 text-sm leading-relaxed" style={{ color: "#C8CCDA" }}>
              {result.feedback && <p>{result.feedback}</p>}
              <p className="font-mono mt-2 text-[11px] font-bold uppercase tracking-widest" style={{ color: "#8B92A8" }}>
                Rubric score {result.score} / 100
              </p>
              {(result.strengths?.length ?? 0) > 0 && (
                <>
                <p className="font-mono mt-2 text-[11px] font-bold uppercase tracking-widest" style={{ color: "#8B92A8" }}>
                  You covered
                </p>
                <ul aria-label="What you covered">
                  {result.strengths!.map((s) => (
                    <li key={s}>
                      <span aria-hidden style={{ color: "#34D399" }}>✓ </span>
                                            {s}
                    </li>
                  ))}
                </ul>
                </>
              )}
              {(result.gaps?.length ?? 0) > 0 && (
                <>
                <p className="font-mono mt-2 text-[11px] font-bold uppercase tracking-widest" style={{ color: "#8B92A8" }}>
                  Worth adding
                </p>
                <ul aria-label="Worth adding">
                  {result.gaps!.map((g) => (
                    <li key={g}>
                      <span aria-hidden style={{ color: "#F5B942" }}>+ </span>
                                            {g}
                    </li>
                  ))}
                </ul>
                </>
              )}
            </div>
          )}

          {result.explanation && (
            <p className="mt-3 text-sm leading-relaxed" style={{ color: "#C8CCDA" }}>
              {prose ? <span className="font-bold">A strong answer: </span> : !ok && <span className="font-bold">Why it works: </span>}
              {result.explanation}
            </p>
          )}
        </div>

        <button
          ref={ref}
          type="button"
          onClick={onContinue}
          className="font-display mt-5 w-full rounded-full px-8 py-4 text-base font-bold tracking-wide focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
          style={{ background: ok ? "#34D399" : "#F87171", color: "#0F1117" }}
        >
          {continueLabel}
        </button>
      </div>
    </div>
  );
}
