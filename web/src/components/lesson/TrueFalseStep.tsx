"use client";

import { useEffect } from "react";

export type TrueFalseResult = { correct: boolean; chosen: boolean; correctValue?: boolean };

const choices = [
  { value: true, label: "True", key: "1" },
  { value: false, label: "False", key: "2" },
] as const;

/**
 * TrueFalseStep renders a statement with True and False. Keys 1/2 (or T/F) select; the parent
 * handles Enter. After grading the choice is marked with an icon and text, never colour alone.
 */
export function TrueFalseStep({
  statement,
  selected,
  result,
  disabled,
  onSelect,
}: {
  statement: string;
  selected: boolean | null;
  result: TrueFalseResult | null;
  disabled: boolean;
  onSelect: (value: boolean) => void;
}) {
  useEffect(() => {
    if (disabled || result) return;
    const onKey = (e: KeyboardEvent) => {
      const k = e.key.toLowerCase();
      if (k === "1" || k === "t") onSelect(true);
      if (k === "2" || k === "f") onSelect(false);
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [disabled, result, onSelect]);

  return (
    <div className="flex flex-1 flex-col">
      <p className="font-mono mb-3 text-[11px] font-bold uppercase tracking-[0.2em]" style={{ color: "#7C6EF5" }}>
        True or false?
      </p>
      <h1 className="font-display text-xl font-extrabold leading-snug sm:text-2xl" style={{ color: "#E8EAED" }}>
        {statement}
      </h1>

      <div role="radiogroup" aria-label="True or false" className="mt-8 grid grid-cols-2 gap-3">
        {choices.map((c) => {
          const shown = result !== null;
          const isChosen = selected === c.value;
          const isRight = shown && result.chosen === c.value && result.correct;
          const isWrong = shown && result.chosen === c.value && !result.correct;
          const isRevealed = shown && !result.correct && result.correctValue === c.value;

          let border = "#2E3347";
          let bg = "#1A1D27";
          if (isChosen && !shown) {
            border = "#7C6EF5";
            bg = "rgba(124,110,245,0.15)";
          }
          if (isRight || isRevealed) {
            border = "#34D399";
            bg = "rgba(52,211,153,0.12)";
          }
          if (isWrong) {
            border = "#F87171";
            bg = "rgba(248,113,113,0.12)";
          }

          return (
            <button
              key={c.label}
              type="button"
              role="radio"
              aria-checked={isChosen}
              disabled={disabled || shown}
              onClick={() => onSelect(c.value)}
              className="flex flex-col items-center justify-center gap-2 rounded-2xl px-4 py-7 transition-[transform,border-color,background-color] duration-100 active:scale-[0.98] disabled:cursor-default motion-reduce:transition-none focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7C6EF5]"
              style={{ background: bg, border: `2px solid ${border}` }}
            >
              <span
                aria-hidden
                className="font-mono flex h-8 w-8 items-center justify-center rounded-lg text-sm font-bold"
                style={{
                  background: isChosen || isRight || isRevealed || isWrong ? border : "#242836",
                  color: isChosen || isRight || isRevealed || isWrong ? "#0F1117" : "#8B92A8",
                }}
              >
                {isRight || isRevealed ? "✓" : isWrong ? "✕" : c.key}
              </span>
              <span className="font-display text-lg font-extrabold" style={{ color: "#E8EAED" }}>
                {c.label}
              </span>
              {isRight && <span className="sr-only">Correct</span>}
              {isWrong && <span className="sr-only">Incorrect</span>}
              {isRevealed && <span className="sr-only">Correct answer</span>}
            </button>
          );
        })}
      </div>
    </div>
  );
}
