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
  correctOptionText,
  companionName,
  onContinue,
  continueLabel,
}: {
  result: AnswerResult;
  correctOptionText?: string;
  companionName?: string;
  onContinue: () => void;
  continueLabel: string;
}) {
  const ref = useRef<HTMLButtonElement>(null);
  useEffect(() => {
    ref.current?.focus({ preventScroll: true });
  }, []);

  const ok = result.correct;
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
        <div role="status" aria-live="polite">
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

          {!ok && correctOptionText && (
            <p className="mt-3 text-sm" style={{ color: "#E8EAED" }}>
              <span className="font-mono text-[11px] font-bold uppercase tracking-widest" style={{ color: "#8B92A8" }}>
                Correct answer
              </span>
              <br />
              <span className="font-semibold">{correctOptionText}</span>
            </p>
          )}

          {!ok && result.why_not && (
            <p className="mt-3 text-sm leading-relaxed" style={{ color: "#C8CCDA" }}>
              <span className="font-bold">Why not that one: </span>
              {result.why_not}
            </p>
          )}

          {result.explanation && (
            <p className="mt-3 text-sm leading-relaxed" style={{ color: "#C8CCDA" }}>
              {!ok && <span className="font-bold">Why it works: </span>}
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
