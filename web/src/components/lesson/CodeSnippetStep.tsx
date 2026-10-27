"use client";

import { KeyboardEvent, useId, useRef } from "react";

function formatLanguageName(lang: string): string {
  switch (lang.toLowerCase()) {
    case "go":
      return "Go";
    case "python":
      return "Python";
    case "node":
    case "nodejs":
      return "Node.js";
    case "javascript":
    case "js":
      return "JavaScript";
    case "typescript":
    case "ts":
      return "TypeScript";
    case "sql":
      return "SQL";
    default:
      return lang;
  }
}

export function CodeSnippetStep({
  prompt,
  language,
  starter,
  code,
  result,
  disabled,
  availableLanguages,
  onChange,
  onSelectLanguage,
}: {
  prompt: string;
  language: string;
  starter: string;
  code: string;
  result: { correct: boolean } | null;
  disabled: boolean;
  availableLanguages?: string[];
  onChange: (code: string) => void;
  onSelectLanguage?: (language: string) => void;
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

  const langs = availableLanguages && availableLanguages.length > 0 ? availableLanguages : [language];

  return (
    <div className="flex flex-1 flex-col">
      <div className="flex flex-wrap items-center justify-between gap-2">
        <p className="font-mono text-[11px] font-bold uppercase tracking-[0.2em]" style={{ color: "#7C6EF5" }}>
          Code Exercise
        </p>

        {/* Multi-Language Switcher Tabs */}
        {langs.length > 1 ? (
          <div
            className="flex items-center gap-1 rounded-xl p-1"
            style={{ background: "#151821", border: "1px solid #2E3347" }}
            role="tablist"
            aria-label="Language selection"
          >
            {langs.map((l) => {
              const active = l.toLowerCase() === language.toLowerCase();
              return (
                <button
                  key={l}
                  type="button"
                  role="tab"
                  aria-selected={active}
                  onClick={() => onSelectLanguage && onSelectLanguage(l)}
                  disabled={disabled || shown}
                  className={`rounded-lg px-2.5 py-1 text-xs font-mono font-bold transition-all ${
                    active
                      ? "bg-[#7C6EF5] text-white shadow-sm"
                      : "text-[#8B92A8] hover:text-[#E8EAED] hover:bg-white/5"
                  } disabled:opacity-50`}
                >
                  {formatLanguageName(l)}
                </button>
              );
            })}
          </div>
        ) : (
          <span
            className="font-mono rounded-md px-2 py-0.5 text-[10px] font-bold uppercase tracking-wider"
            style={{ background: "#1F2330", color: "#8B92A8", border: "1px solid #2E3347" }}
          >
            {formatLanguageName(language)}
          </span>
        )}
      </div>

      <h1 id={promptId} className="font-display mt-3 text-xl font-extrabold leading-snug sm:text-2xl" style={{ color: "#E8EAED" }}>
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
