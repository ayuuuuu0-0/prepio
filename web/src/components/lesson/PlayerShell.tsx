"use client";

import { ReactNode, useEffect, useRef } from "react";

/** PlayerShell is the focus-mode frame: no bottom navigation, a close button, progress, and combo. */
export function PlayerShell({
  progress,
  combo,
  confirmingExit,
  onRequestExit,
  onCancelExit,
  onConfirmExit,
  children,
  tray,
  companion,
  onReplayIntro,
  soundOn,
  onToggleSound,
}: {
  progress: number;
  combo: number;
  confirmingExit: boolean;
  onRequestExit: () => void;
  onCancelExit: () => void;
  onConfirmExit: () => void;
  children: ReactNode;
  tray?: ReactNode;
  companion?: ReactNode;
  /** When set, shows a button that replays the lesson intro. */
  onReplayIntro?: () => void;
  soundOn: boolean;
  onToggleSound: () => void;
}) {
  const pct = Math.round(Math.max(0, Math.min(1, progress)) * 100);

  return (
    <div className="game-bg-challenge relative flex min-h-dvh flex-col">
      <header inert={confirmingExit} className="relative z-10 mx-auto flex w-full max-w-2xl items-center gap-3 px-4 pb-2 pt-4">
        <button
          type="button"
          onClick={onRequestExit}
          aria-label="Leave lesson"
          className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full text-lg transition-colors hover:bg-white/10 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7C6EF5]"
          style={{ color: "#8B92A8" }}
        >
          <span aria-hidden>✕</span>
        </button>

        <div
          className="h-3 flex-1 overflow-hidden rounded-full"
          style={{ background: "#242836" }}
          role="progressbar"
          aria-label="Lesson progress"
          aria-valuemin={0}
          aria-valuemax={100}
          aria-valuenow={pct}
        >
          <div
            className="h-full rounded-full transition-[width] duration-500 ease-out motion-reduce:transition-none"
            style={{
              width: `${pct}%`,
              background: "linear-gradient(90deg, #7C6EF5, #9D8FF7)",
              boxShadow: pct > 0 ? "0 0 12px rgba(124,110,245,0.5)" : "none",
            }}
          />
        </div>

        <div
          className="font-mono flex h-10 min-w-[3.25rem] shrink-0 items-center justify-center gap-1 rounded-full px-3 text-sm font-bold"
          style={{
            background: combo >= 2 ? "rgba(255,107,53,0.15)" : "#1A1D27",
            border: `1px solid ${combo >= 2 ? "#FF6B35" : "#2E3347"}`,
            color: combo >= 2 ? "#FF6B35" : "#4A5068",
          }}
          aria-label={`Combo ${combo}`}
        >
          <span aria-hidden>⚡</span>
          {combo}
        </div>

        {onReplayIntro && (
          <button
            type="button"
            onClick={onReplayIntro}
            aria-label="Replay the intro"
            className="font-display flex h-10 shrink-0 items-center gap-1 rounded-full px-3 text-xs font-bold transition-colors hover:bg-white/10 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7C6EF5]"
            style={{ color: "#8B92A8", border: "1px solid #2E3347" }}
          >
            <span aria-hidden>↺</span>
            <span aria-hidden className="hidden sm:inline">Intro</span>
          </button>
        )}

        <button
          type="button"
          onClick={onToggleSound}
          aria-pressed={soundOn}
          aria-label={soundOn ? "Sound on" : "Sound off"}
          title={soundOn ? "Sound on" : "Sound off"}
          className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full text-base transition-colors hover:bg-white/10 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7C6EF5]"
          style={{ color: soundOn ? "#E8EAED" : "#4A5068", border: "1px solid #2E3347" }}
        >
          <span aria-hidden>{soundOn ? "🔊" : "🔇"}</span>
        </button>

        {companion}
      </header>

      <main inert={confirmingExit} className="relative z-10 mx-auto flex w-full max-w-2xl flex-1 flex-col px-4 pb-6 pt-4">{children}</main>

      {tray && (
        <div inert={confirmingExit} className="contents">
          {tray}
        </div>
      )}

      {confirmingExit && <ExitDialog onCancel={onCancelExit} onConfirm={onConfirmExit} />}
    </div>
  );
}

function ExitDialog({ onCancel, onConfirm }: { onCancel: () => void; onConfirm: () => void }) {
  const stayRef = useRef<HTMLButtonElement>(null);

  useEffect(() => {
    stayRef.current?.focus();
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") onCancel();
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [onCancel]);

  return (
    <div className="fixed inset-0 z-50 flex items-end justify-center bg-black/60 p-4 sm:items-center" role="presentation">
      <div
        role="alertdialog"
        aria-modal="true"
        aria-labelledby="exit-title"
        aria-describedby="exit-desc"
        className="w-full max-w-sm rounded-2xl p-6"
        style={{ background: "#1A1D27", border: "1px solid #2E3347", boxShadow: "0 12px 48px rgba(0,0,0,0.6)" }}
      >
        <h2 id="exit-title" className="font-display text-lg font-extrabold" style={{ color: "#E8EAED" }}>
          Leave this lesson?
        </h2>
        <p id="exit-desc" className="mt-2 text-sm leading-relaxed" style={{ color: "#8B92A8" }}>
          Your progress is saved. You can pick up right where you left off.
        </p>
        <div className="mt-5 flex flex-col gap-3">
          <button
            ref={stayRef}
            type="button"
            onClick={onCancel}
            className="font-display w-full rounded-full px-6 py-3 text-sm font-bold text-white focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
            style={{ background: "linear-gradient(90deg, #7C6EF5, #9D8FF7)" }}
          >
            Keep going
          </button>
          <button
            type="button"
            onClick={onConfirm}
            className="font-display w-full rounded-full px-6 py-3 text-sm font-bold transition-colors hover:bg-white/5 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7C6EF5]"
            style={{ color: "#8B92A8", border: "1px solid #2E3347" }}
          >
            Leave lesson
          </button>
        </div>
      </div>
    </div>
  );
}
