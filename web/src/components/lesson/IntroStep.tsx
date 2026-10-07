"use client";

import { useEffect, useState } from "react";
import type { Beat } from "@/lib/lesson/types";

/** emphasise wraps the emphasised phrase of a beat in a highlight, keeping the rest as plain text. */
function emphasise(beat: Beat) {
  if (!beat.emphasis || !beat.text.includes(beat.emphasis)) return beat.text;
  const [before, ...rest] = beat.text.split(beat.emphasis);
  return (
    <>
      {before}
      <span style={{ color: "#B6ACFF", textShadow: "0 0 24px rgba(124,110,245,0.6)" }}>{beat.emphasis}</span>
      {rest.join(beat.emphasis)}
    </>
  );
}

/**
 * IntroStep walks through the lesson's beats one at a time. Tap, Enter, or Space advances;
 * "Skip intro" is always available. (The animated motion stage arrives in a later phase.)
 */
export function IntroStep({ beats, onDone }: { beats: Beat[]; onDone: () => void }) {
  const [index, setIndex] = useState(0);
  const last = index >= beats.length - 1;

  const advance = () => (last ? onDone() : setIndex((i) => i + 1));

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Enter" || e.key === " ") {
        e.preventDefault();
        if (index >= beats.length - 1) onDone();
        else setIndex((i) => i + 1);
      }
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [index, beats.length, onDone]);

  return (
    <div className="flex flex-1 flex-col">
      <button
        type="button"
        onClick={advance}
        aria-label={last ? "Start the lesson" : "Next"}
        className="flex flex-1 flex-col items-center justify-center rounded-3xl px-6 text-center focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-[#7C6EF5]"
      >
        <p className="font-mono mb-6 text-[11px] font-bold uppercase tracking-[0.2em]" style={{ color: "#7C6EF5" }}>
          Quick intro
        </p>
        <p
          key={index}
          aria-live="polite"
          className="font-display animate-result text-2xl font-extrabold leading-snug sm:text-3xl"
          style={{ color: "#E8EAED" }}
        >
          {emphasise(beats[index])}
        </p>
        <div className="mt-10 flex gap-2" aria-hidden>
          {beats.map((_, i) => (
            <span
              key={i}
              className="h-1.5 rounded-full transition-all"
              style={{ width: i === index ? 28 : 8, background: i <= index ? "#7C6EF5" : "#2E3347" }}
            />
          ))}
        </div>
        <p className="mt-6 text-xs font-semibold" style={{ color: "#4A5068" }}>
          Tap to continue
        </p>
      </button>

      <button
        type="button"
        onClick={onDone}
        className="font-display mx-auto mt-4 rounded-full px-5 py-2 text-sm font-bold transition-colors hover:bg-white/5 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7C6EF5]"
        style={{ color: "#8B92A8" }}
      >
        Skip intro
      </button>
    </div>
  );
}
