"use client";

import { KeyboardEvent, useId, useRef } from "react";

export function CodeSnippetStep({
  prompt,
  language,
  starter,
  code,
  result,
  disabled,
  onChange,
}: {
  prompt: string;
  language: "go" | "python" | "sql" | "typescript";
  starter: string;
  code: string;
  result: { correct: boolean } | null;
  disabled: boolean;
  onChange: (code: string) => void;
}) {
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const promptId = useId();
  const countId = useId();

  const shown = result !== null;
  const changed = code.trim().length > 0 && code.trim() !== starter.trim();
  const tone = shown ? (result.correct ? "#34D399" : "#F87171") : changed ? "#7C6EF5" : "#2E3347";

  const lines = code.split("\n");
  const lineCount = Math.max(lines.length, 6);

  const handleKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
    if (disabled || shown) return;

    // Tab key inserts 2 spaces
    if (e.key === "Tab") {
      e.preventDefault();
      const target = e.currentTarget;
      const start = target.selectionStart;
      const end = target.selectionEnd;
      const val = target.value;

      const next = val.substring(0, start) + "  " + val.substring(end);
      onChange(next);

      // Restore cursor position after state update
      requestAnimationFrame(() => {
        if (textareaRef.current) {
          textareaRef.current.selectionStart = textareaRef.current.selectionEnd = start + 2;
        }
      });
      return;
    }

    // Enter preserves indentation of current line
    if (e.key === "Enter" && !e.ctrlKey && !e.metaKey) {
      e.preventDefault();
      const target = e.currentTarget;
      const start = target.selectionStart;
      const val = target.value;

      // Find indentation of current line
      const currentLineStart = val.lastIndexOf("\n", start - 1) + 1;
      const currentLine = val.substring(currentLineStart, start);
      const match = currentLine.match(/^(\s+)/);
      const indent = match ? match[1] : "";

      const next = val.substring(0, start) + "\n" + indent + val.substring(target.selectionEnd);
      onChange(next);

      requestAnimationFrame(() => {
        if (textareaRef.current) {
          textareaRef.current.selectionStart = textareaRef.current.selectionEnd = start + 1 + indent.length;
        }
      });
    }
  };

  const handleReset = () => {
    if (disabled || shown) return;
    onChange(starter);
    if (textareaRef.current) {
      textareaRef.current.focus();
    }
  };

  return (
    <div className="flex flex-1 flex-col">
      <div className="flex items-center justify-between">
        <p className="font-mono mb-2 text-[11px] font-bold uppercase tracking-[0.2em]" style={{ color: "#7C6EF5" }}>
          Code Exercise
        </p>
        <span
          className="font-mono rounded-md px-2 py-0.5 text-[10px] font-bold uppercase tracking-wider"
          style={{ background: "#1F2330", color: "#8B92A8", border: "1px solid #2E3347" }}
        >
          {language}
        </span>
      </div>

      <h1 id={promptId} className="font-display text-xl font-extrabold leading-snug sm:text-2xl" style={{ color: "#E8EAED" }}>
        {prompt}
      </h1>

      {/* Editor Box */}
      <div
        className="mt-5 flex overflow-hidden rounded-2xl"
        style={{
          background: "#0D0F14",
          border: `2px solid ${tone}`,
          boxShadow: shown
            ? result.correct
              ? "0 0 24px rgba(52,211,153,0.18)"
              : "0 0 24px rgba(248,113,113,0.18)"
            : "none",
        }}
      >
        {/* Line Numbers Column */}
        <div
          aria-hidden="true"
          className="font-mono select-none px-3 py-3.5 text-right text-xs leading-6"
          style={{ background: "#0B0C10", color: "#4A5068", borderRight: "1px solid #1E2230" }}
        >
          {Array.from({ length: lineCount }).map((_, i) => (
            <div key={i}>{i + 1}</div>
          ))}
        </div>

        {/* Textarea Code Buffer */}
        <textarea
          ref={textareaRef}
          aria-labelledby={promptId}
          aria-describedby={countId}
          value={code}
          onChange={(e) => onChange(e.target.value)}
          onKeyDown={handleKeyDown}
          readOnly={disabled || shown}
          spellCheck={false}
          autoCapitalize="off"
          autoComplete="off"
          autoCorrect="off"
          rows={Math.max(lineCount, 8)}
          className="font-mono w-full resize-none p-3.5 text-xs sm:text-sm leading-6 outline-none"
          style={{
            background: "transparent",
            color: "#E8EAED",
            tabSize: 2,
          }}
        />
      </div>

      {/* Bottom status bar */}
      <div id={countId} className="font-mono mt-2.5 flex items-center justify-between text-xs" style={{ color: "#8B92A8" }}>
        <div className="flex items-center gap-3">
          <span>{lines.length} lines</span>
          {!shown && (
            <button
              type="button"
              onClick={handleReset}
              disabled={disabled || code === starter}
              className="text-xs underline-offset-2 hover:underline disabled:opacity-40 disabled:hover:no-underline"
              style={{ color: "#7C6EF5" }}
            >
              Reset to starter
            </button>
          )}
        </div>
        <span className="hidden sm:inline">Tab indents · Ctrl+Enter to check</span>
      </div>
    </div>
  );
}
