"use client";

import { useEffect } from "react";

export type McqResult = { correctIndex?: number; correct: boolean; chosen: number };

/**
 * McqStep renders one multiple-choice exercise. Keys 1-4 select an option; the
 * parent handles Enter. After grading, the chosen option is marked with an icon
 * and text (never colour alone) and the correct option is revealed on a miss.
 */
export function McqStep({
  prompt,
  options,
  selected,
  result,
  disabled,
  onSelect,
}: {
  prompt: string;
  options: string[];
  selected: number | null;
  result: McqResult | null;
  disabled: boolean;
  onSelect: (index: number) => void;
}) {
  useEffect(() => {
    if (disabled || result) return;
    const onKey = (e: KeyboardEvent) => {
      const n = Number(e.key);
      if (Number.isInteger(n) && n >= 1 && n <= options.length) onSelect(n - 1);
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [disabled, result, options.length, onSelect]);

  return (
    <div className="flex flex-1 flex-col">
      <p className="font-mono mb-3 text-[11px] font-bold uppercase tracking-[0.2em]" style={{ color: "#7C6EF5" }}>
        Pick one
      </p>
      <h1 className="font-display text-xl font-extrabold leading-snug sm:text-2xl" style={{ color: "#E8EAED" }}>
        {prompt}
      </h1>

      <div role="radiogroup" aria-label="Answer options" className="mt-6 flex flex-col gap-3">
        {options.map((text, i) => {
          const isChosen = selected === i;
          const showResult = result !== null;
          const isWrongChoice = showResult && result.chosen === i && !result.correct;
          const isRightChoice = showResult && result.chosen === i && result.correct;
          const isRevealed = showResult && !result.correct && result.correctIndex === i;

          let border = "#2E3347";
          let bg = "#1A1D27";
          if (isChosen && !showResult) {
            border = "#7C6EF5";
            bg = "rgba(124,110,245,0.15)";
          }
          if (isRightChoice || isRevealed) {
            border = "#34D399";
            bg = "rgba(52,211,153,0.12)";
          }
          if (isWrongChoice) {
            border = "#F87171";
            bg = "rgba(248,113,113,0.12)";
          }

          return (
            <button
              key={i}
              type="button"
              role="radio"
              aria-checked={isChosen}
              disabled={disabled || showResult}
              onClick={() => onSelect(i)}
              className="flex w-full items-center gap-4 rounded-2xl px-4 py-4 text-left transition-[transform,border-color,background-color] duration-100 active:scale-[0.99] disabled:cursor-default motion-reduce:transition-none focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7C6EF5]"
              style={{ background: bg, border: `2px solid ${border}`, minHeight: 64 }}
            >
              <span
                className="font-mono flex h-8 w-8 shrink-0 items-center justify-center rounded-lg text-sm font-bold"
                style={{
                  background: isChosen || isRightChoice || isRevealed ? border : "#242836",
                  color: isChosen || isRightChoice || isRevealed || isWrongChoice ? "#0F1117" : "#8B92A8",
                }}
                aria-hidden
              >
                {isRightChoice || isRevealed ? "✓" : isWrongChoice ? "✕" : i + 1}
              </span>
              <span className="flex-1 text-base font-semibold leading-snug" style={{ color: "#E8EAED" }}>
                {text}
              </span>
              {isRightChoice && <span className="sr-only">Correct</span>}
              {isWrongChoice && <span className="sr-only">Incorrect</span>}
              {isRevealed && <span className="sr-only">Correct answer</span>}
            </button>
          );
        })}
      </div>
    </div>
  );
}
