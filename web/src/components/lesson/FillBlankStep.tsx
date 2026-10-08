"use client";

import { Fragment } from "react";
import { clearSlot, placeChip, splitCode } from "@/lib/lesson/answers";

/**
 * FillBlankStep renders code with numbered blanks and a word bank. Tapping a chip places it
 * in the first empty blank; tapping a filled blank sends its chip back. No dragging needed.
 */
export function FillBlankStep({
  code,
  bank,
  slots,
  result,
  disabled,
  onChange,
}: {
  code: string;
  bank: string[];
  slots: (number | null)[];
  result: { correct: boolean } | null;
  disabled: boolean;
  onChange: (slots: (number | null)[]) => void;
}) {
  const shown = result !== null;
  const locked = disabled || shown;
  const tone = !shown ? "#7C6EF5" : result.correct ? "#34D399" : "#F87171";
  const placed = new Set(slots.filter((s): s is number => s !== null));

  return (
    <div className="flex flex-1 flex-col">
      <p className="font-mono mb-3 text-[11px] font-bold uppercase tracking-[0.2em]" style={{ color: "#7C6EF5" }}>
        Fill in the blanks
      </p>
      <h1 className="sr-only">Fill in the blanks in the code</h1>

      <pre
        className="font-mono overflow-x-auto whitespace-pre-wrap rounded-2xl p-4 text-sm leading-8"
        style={{ background: "#13151C", border: `2px solid ${shown ? tone : "#2E3347"}`, color: "#C8CCDA" }}
      >
        {splitCode(code).map((part, i) =>
          part.kind === "text" ? (
            <Fragment key={i}>{part.text}</Fragment>
          ) : (
            <button
              key={i}
              type="button"
              disabled={locked || slots[part.blank] === null}
              onClick={() => onChange(clearSlot(slots, part.blank))}
              aria-label={
                slots[part.blank] === null
                  ? `Blank ${part.blank + 1}, empty`
                  : `Blank ${part.blank + 1}: ${bank[slots[part.blank]!]}. Tap to remove`
              }
              className="mx-0.5 inline-flex min-w-[4.5rem] items-center justify-center rounded-lg px-2 align-middle font-semibold leading-7 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7C6EF5] disabled:cursor-default"
              style={
                slots[part.blank] === null
                  ? { border: "2px dashed #4A5068", color: "#4A5068" }
                  : { background: `${tone}26`, border: `2px solid ${tone}`, color: "#E8EAED" }
              }
            >
              {slots[part.blank] === null ? part.blank + 1 : bank[slots[part.blank]!]}
            </button>
          ),
        )}
      </pre>
      {shown && (
        <p className="mt-2 text-sm font-semibold" style={{ color: tone }}>
          <span aria-hidden>{result.correct ? "✓ " : "✕ "}</span>
          {result.correct ? "Every blank is right" : "Not every blank is right"}
        </p>
      )}

      <p className="font-mono mb-2 mt-6 text-[11px] font-bold uppercase tracking-[0.2em]" style={{ color: "#8B92A8" }}>
        Word bank
      </p>
      <ul className="flex flex-wrap gap-2" aria-label="Word bank">
        {bank.map((word, i) => {
          const used = placed.has(i);
          return (
            <li key={i}>
              <button
                type="button"
                disabled={locked || used}
                onClick={() => onChange(placeChip(slots, i))}
                aria-label={used ? `${word}, placed` : `Place ${word}`}
                className="font-mono rounded-xl px-3.5 py-2.5 text-sm font-semibold transition-[transform,opacity] duration-100 active:scale-[0.97] disabled:cursor-default motion-reduce:transition-none focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7C6EF5]"
                style={{
                  background: used ? "transparent" : "#1A1D27",
                  border: used ? "2px dashed #2E3347" : "2px solid #2E3347",
                  color: used ? "#4A5068" : "#E8EAED",
                  boxShadow: used ? "none" : "0 3px 0 #242836",
                }}
              >
                {word}
              </button>
            </li>
          );
        })}
      </ul>
    </div>
  );
}
