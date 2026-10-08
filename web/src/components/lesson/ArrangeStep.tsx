"use client";

import { toggleArrange } from "@/lib/lesson/answers";

/**
 * ArrangeStep asks for items in order. Tapping an item adds it to the next position; tapping a
 * placed item takes it back out. No dragging needed, and every control is a button.
 */
export function ArrangeStep({
  prompt,
  items,
  order,
  result,
  disabled,
  onChange,
}: {
  prompt: string;
  items: string[];
  order: number[];
  result: { correct: boolean } | null;
  disabled: boolean;
  onChange: (order: number[]) => void;
}) {
  const shown = result !== null;
  const locked = disabled || shown;
  const tone = !shown ? "#7C6EF5" : result.correct ? "#34D399" : "#F87171";
  const remaining = items.map((_, i) => i).filter((i) => !order.includes(i));

  return (
    <div className="flex flex-1 flex-col">
      <p className="font-mono mb-3 text-[11px] font-bold uppercase tracking-[0.2em]" style={{ color: "#7C6EF5" }}>
        Put these in order
      </p>
      <h1 className="font-display text-xl font-extrabold leading-snug sm:text-2xl" style={{ color: "#E8EAED" }}>
        {prompt}
      </h1>

      <ol className="mt-6 flex flex-col gap-2" aria-label="Your order">
        {items.map((_, pos) => {
          const item = order[pos];
          const filled = item !== undefined;
          return (
            <li key={pos}>
              <button
                type="button"
                disabled={locked || !filled}
                onClick={() => filled && onChange(toggleArrange(order, item))}
                aria-label={filled ? `Position ${pos + 1}: ${items[item]}. Tap to remove` : `Position ${pos + 1}, empty`}
                className="flex w-full items-center gap-3 rounded-2xl px-4 py-3 text-left disabled:cursor-default focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7C6EF5]"
                style={
                  filled
                    ? { background: `${tone}1F`, border: `2px solid ${tone}`, minHeight: 56 }
                    : { background: "transparent", border: "2px dashed #2E3347", minHeight: 56 }
                }
              >
                <span
                  aria-hidden
                  className="font-mono flex h-7 w-7 shrink-0 items-center justify-center rounded-lg text-xs font-bold"
                  style={{ background: filled ? tone : "#242836", color: filled ? "#0F1117" : "#4A5068" }}
                >
                  {pos + 1}
                </span>
                <span className="flex-1 font-semibold leading-snug" style={{ color: filled ? "#E8EAED" : "#4A5068" }}>
                  {filled ? items[item] : "Tap an item below"}
                </span>
              </button>
            </li>
          );
        })}
      </ol>
      {shown && (
        <p className="mt-2 text-sm font-semibold" style={{ color: tone }}>
          <span aria-hidden>{result.correct ? "✓ " : "✕ "}</span>
          {result.correct ? "That's the right order" : "Not quite the right order"}
        </p>
      )}

      {!shown && remaining.length > 0 && (
        <>
          <p className="font-mono mb-2 mt-6 text-[11px] font-bold uppercase tracking-[0.2em]" style={{ color: "#8B92A8" }}>
            Tap in order
          </p>
          <ul className="flex flex-col gap-2" aria-label="Items to place">
            {remaining.map((i) => (
              <li key={i}>
                <button
                  type="button"
                  disabled={locked}
                  onClick={() => onChange(toggleArrange(order, i))}
                  className="w-full rounded-2xl px-4 py-3 text-left font-semibold transition-transform duration-100 active:scale-[0.99] motion-reduce:transition-none focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7C6EF5]"
                  style={{ background: "#1A1D27", border: "2px solid #2E3347", color: "#E8EAED", boxShadow: "0 3px 0 #242836" }}
                >
                  {items[i]}
                </button>
              </li>
            ))}
          </ul>
        </>
      )}
    </div>
  );
}
