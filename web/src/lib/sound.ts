"use client";

import { useCallback, useEffect, useState } from "react";

/**
 * Lesson sound effects, synthesized with WebAudio (no audio files). Sound is optional,
 * off by default, and the choice is remembered on this device.
 */

export type Cue = "correct" | "wrong" | "complete";

const STORAGE_KEY = "prepio_sound";

/** Notes per cue: [frequency Hz, start offset s, duration s]. */
export const CUES: Record<Cue, [number, number, number][]> = {
  correct: [
    [660, 0, 0.12],
    [880, 0.09, 0.18],
  ],
  wrong: [
    [330, 0, 0.16],
    [262, 0.12, 0.22],
  ],
  complete: [
    [523, 0, 0.14],
    [659, 0.12, 0.14],
    [784, 0.24, 0.14],
    [1047, 0.36, 0.32],
  ],
};

/** readSoundPref returns the stored preference; anything but "on" means off. */
export function readSoundPref(storage: Pick<Storage, "getItem"> | undefined): boolean {
  try {
    return storage?.getItem(STORAGE_KEY) === "on";
  } catch {
    return false;
  }
}

let ctx: AudioContext | null = null;

function play(cue: Cue) {
  if (typeof window === "undefined" || !("AudioContext" in window)) return;
  ctx ??= new AudioContext();
  if (ctx.state === "suspended") void ctx.resume();
  const now = ctx.currentTime;
  for (const [freq, at, dur] of CUES[cue]) {
    const osc = ctx.createOscillator();
    const gain = ctx.createGain();
    osc.type = cue === "wrong" ? "triangle" : "sine";
    osc.frequency.value = freq;
    gain.gain.setValueAtTime(0.0001, now + at);
    gain.gain.exponentialRampToValueAtTime(0.18, now + at + 0.02);
    gain.gain.exponentialRampToValueAtTime(0.0001, now + at + dur);
    osc.connect(gain).connect(ctx.destination);
    osc.start(now + at);
    osc.stop(now + at + dur + 0.05);
  }
}

/** useSound exposes the remembered on/off preference and a player that is silent when off. */
export function useSound() {
  const [enabled, setEnabled] = useState(false);

  useEffect(() => {
    setEnabled(readSoundPref(typeof window === "undefined" ? undefined : window.localStorage));
  }, []);

  const toggle = useCallback(() => {
    setEnabled((on) => {
      const next = !on;
      try {
        window.localStorage.setItem(STORAGE_KEY, next ? "on" : "off");
      } catch {
        // Storage can be unavailable (private mode); the toggle still works for this visit.
      }
      if (next) play("correct"); // a short confirmation that sound is on
      return next;
    });
  }, []);

  const cue = useCallback(
    (c: Cue) => {
      if (enabled) play(c);
    },
    [enabled],
  );

  return { enabled, toggle, cue };
}
