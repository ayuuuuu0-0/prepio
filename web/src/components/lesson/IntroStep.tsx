"use client";

import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import type { Beat } from "@/lib/lesson/types";
import { beatProgress, mediaKind, revealCount, splitEmphasis, words } from "@/lib/lesson/intro";
import { useReducedMotion } from "@/lib/useReducedMotion";
import { IntroVisual } from "./visuals";

type Piece = { text: string; emphasised: boolean };

/** pieces splits a beat into words in reading order, flagging the emphasised phrase. */
function pieces(beat: Beat): Piece[] {
  const seg = splitEmphasis(beat);
  return [
    ...words(seg.before).map((text) => ({ text, emphasised: false })),
    ...words(seg.emphasis).map((text) => ({ text, emphasised: true })),
    ...words(seg.after).map((text) => ({ text, emphasised: false })),
  ];
}

function Media({ url }: { url: string }) {
  const [failed, setFailed] = useState(false);
  const kind = mediaKind(url);
  if (!kind || failed) return null;
  const cls = "max-h-48 w-auto max-w-full rounded-2xl object-contain";
  return kind === "video" ? (
    <video src={url} className={cls} autoPlay muted loop playsInline aria-hidden onError={() => setFailed(true)} />
  ) : (
    // eslint-disable-next-line @next/next/no-img-element -- authored media of unknown size; decorative
    <img src={url} alt="" className={cls} aria-hidden onError={() => setFailed(true)} />
  );
}

/**
 * IntroStep plays the lesson's beats as a short motion piece: words reveal as each beat plays,
 * the emphasised phrase glows, an optional named visual or media sits above, and beats advance
 * on their own timing. Tap, Enter, or Space advances; "Skip intro" is always available; it can be
 * paused; and at the end it can be replayed. With reduced motion it never auto-advances and
 * shows each beat's full text at once.
 */
export function IntroStep({
  beats,
  mediaUrl,
  onDone,
}: {
  beats: Beat[];
  mediaUrl?: string;
  onDone: () => void;
}) {
  const reduced = useReducedMotion();
  const [index, setIndex] = useState(0);
  const [elapsed, setElapsed] = useState(0);
  const [paused, setPaused] = useState(false);
  const [ended, setEnded] = useState(false);
  const last = index >= beats.length - 1;
  const beat = beats[index];

  const autoplay = !reduced && !paused && !ended;

  // One animation-frame clock per beat; it stops when paused, ended, or motion is reduced.
  const frame = useRef<number | null>(null);
  useEffect(() => {
    if (!autoplay) return;
    let previous = performance.now();
    const tick = (now: number) => {
      const dt = now - previous;
      previous = now;
      setElapsed((e) => e + dt);
      frame.current = requestAnimationFrame(tick);
    };
    frame.current = requestAnimationFrame(tick);
    return () => {
      if (frame.current !== null) cancelAnimationFrame(frame.current);
    };
  }, [autoplay, index]);

  const next = useCallback(() => {
    if (last) {
      setEnded(true);
      setElapsed(beat.duration_ms);
      return;
    }
    setIndex((i) => i + 1);
    setElapsed(0);
  }, [last, beat.duration_ms]);

  // Beats move on by themselves once their time is up.
  useEffect(() => {
    if (autoplay && elapsed >= beat.duration_ms) next();
  }, [autoplay, elapsed, beat.duration_ms, next]);

  const replay = useCallback(() => {
    setIndex(0);
    setElapsed(0);
    setEnded(false);
    setPaused(false);
  }, []);

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if ((e.key === "Enter" || e.key === " ") && !(e.target as HTMLElement | null)?.closest("button")) {
        e.preventDefault();
        if (ended) onDone();
        else next();
      }
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [ended, next, onDone]);

  const parts = useMemo(() => pieces(beat), [beat]);
  const shown = reduced || ended ? parts.length : revealCount(elapsed, beat.duration_ms, parts.length);
  const plainText = beat.text;
  const showMedia = index === 0 && mediaKind(mediaUrl) !== null;

  return (
    <div className="flex flex-1 flex-col">
      <button
        type="button"
        onClick={() => (ended ? undefined : next())}
        disabled={ended}
        aria-label={last ? "Finish the intro" : "Next beat"}
        className="flex flex-1 flex-col items-center justify-center rounded-3xl px-6 text-center focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-[#7C6EF5] disabled:cursor-default"
      >
        <p className="font-mono mb-5 text-[11px] font-bold uppercase tracking-[0.2em]" style={{ color: "#7C6EF5" }}>
          Quick intro
        </p>

        {(showMedia || beat.visual) && (
          <div className="mb-6 flex min-h-[7rem] w-full items-center justify-center" aria-hidden={showMedia}>
            {showMedia && mediaUrl ? <Media url={mediaUrl} /> : <IntroVisual name={beat.visual} />}
          </div>
        )}

        <p key={index} className="font-display text-2xl font-extrabold leading-snug sm:text-3xl" style={{ color: "#E8EAED" }}>
          <span className="sr-only" aria-live="polite" aria-atomic="true">
            {plainText}
          </span>
          <span aria-hidden>
            {parts.map((piece, i) => (
              <span
                key={i}
                className="inline-block whitespace-pre transition-[opacity,transform] duration-300 motion-reduce:transition-none"
                style={{
                  opacity: i < shown ? 1 : 0,
                  transform: i < shown ? "translateY(0)" : "translateY(0.4em)",
                  color: piece.emphasised ? "#B6ACFF" : undefined,
                  textShadow: piece.emphasised && i < shown ? "0 0 24px rgba(124,110,245,0.7)" : undefined,
                }}
              >
                {piece.text}
              </span>
            ))}
          </span>
        </p>

        <div className="mt-10 flex w-full max-w-xs gap-2" aria-hidden>
          {beats.map((b, i) => {
            const fill = i < index || ended ? 1 : i === index ? beatProgress(elapsed, b.duration_ms) : 0;
            return (
              <span key={i} className="h-1.5 flex-1 overflow-hidden rounded-full" style={{ background: "#2E3347" }}>
                <span className="block h-full rounded-full" style={{ width: `${fill * 100}%`, background: "#7C6EF5" }} />
              </span>
            );
          })}
        </div>

        {!ended && (
          <p className="mt-5 text-xs font-semibold" style={{ color: "#8B92A8" }}>
            Tap to continue
          </p>
        )}
      </button>

      {ended ? (
        <div className="mx-auto mt-4 flex w-full max-w-xs flex-col gap-3">
          <button
            type="button"
            onClick={onDone}
            autoFocus
            className="font-display w-full rounded-full px-8 py-4 text-base font-bold text-white focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
            style={{ background: "linear-gradient(90deg, #7C6EF5, #9D8FF7)" }}
          >
            Start the lesson
          </button>
          <button
            type="button"
            onClick={replay}
            className="font-display w-full rounded-full px-8 py-3 text-sm font-bold transition-colors hover:bg-white/5 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7C6EF5]"
            style={{ color: "#C8CCDA", border: "1px solid #2E3347" }}
          >
            Replay intro
          </button>
        </div>
      ) : (
        <div className="mx-auto mt-4 flex items-center gap-2">
          {!reduced && (
            <button
              type="button"
              onClick={() => setPaused((p) => !p)}
              aria-pressed={paused}
              className="font-display rounded-full px-4 py-2 text-sm font-bold transition-colors hover:bg-white/5 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7C6EF5]"
              style={{ color: "#8B92A8" }}
            >
              {paused ? "▶ Play" : "❚❚ Pause"}
            </button>
          )}
          <button
            type="button"
            onClick={onDone}
            className="font-display rounded-full px-5 py-2 text-sm font-bold transition-colors hover:bg-white/5 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7C6EF5]"
            style={{ color: "#8B92A8" }}
          >
            Skip intro
          </button>
        </div>
      )}
    </div>
  );
}
