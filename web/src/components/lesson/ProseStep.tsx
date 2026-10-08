"use client";

/**
 * ProseStep asks for a short written answer, graded by the server's rubric. Enter adds a new
 * line; Ctrl+Enter (or Cmd+Enter) checks, handled by the page. The counter shows how much is
 * still needed in words and numbers, never colour alone.
 */
export function ProseStep({
  prompt,
  minChars,
  text,
  result,
  disabled,
  onChange,
}: {
  prompt: string;
  minChars: number;
  text: string;
  result: { correct: boolean } | null;
  disabled: boolean;
  onChange: (text: string) => void;
}) {
  const shown = result !== null;
  const count = text.trim().length;
  const enough = count >= minChars;
  const tone = shown ? (result.correct ? "#34D399" : "#F87171") : enough ? "#34D399" : "#2E3347";

  return (
    <div className="flex flex-1 flex-col">
      <p className="font-mono mb-3 text-[11px] font-bold uppercase tracking-[0.2em]" style={{ color: "#7C6EF5" }}>
        In your own words
      </p>
      <h1 id="prose-prompt" className="font-display text-xl font-extrabold leading-snug sm:text-2xl" style={{ color: "#E8EAED" }}>
        {prompt}
      </h1>

      <textarea
        aria-labelledby="prose-prompt"
        aria-describedby="prose-count"
        value={text}
        onChange={(e) => onChange(e.target.value)}
        readOnly={disabled || shown}
        maxLength={4000}
        rows={7}
        placeholder="Explain it as you would to a teammate."
        className="mt-6 w-full resize-y rounded-2xl p-4 text-base leading-relaxed outline-none focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7C6EF5]"
        style={{ background: "#13151C", border: `2px solid ${tone}`, color: "#E8EAED" }}
      />
      <p id="prose-count" className="font-mono mt-2 flex justify-between gap-3 text-xs" style={{ color: "#8B92A8" }}>
        <span>
          {enough ? (
            <>
              <span aria-hidden style={{ color: "#34D399" }}>✓ </span>Long enough to check
            </>
          ) : (
            `${minChars - count} more characters to go`
          )}
        </span>
        <span className="hidden sm:inline">Ctrl+Enter to check</span>
      </p>
    </div>
  );
}
