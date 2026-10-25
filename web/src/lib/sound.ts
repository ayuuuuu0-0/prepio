"use client";

import { useCallback, useEffect, useState } from "react";

/**
 * Lesson sound effects, synthesized with WebAudio (no audio files). Sound is optional,
 * off by default, and the choice is remembered on this device.
 */

export type Cue = "correct" | "wrong" | "complete" | "poke";

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
  poke: [
    [440, 0, 0.08],
    [554, 0.06, 0.12],
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

/** speakCompanionLine uses Web Speech API to voice companion dialogue with customized personality pitch. */
export function speakCompanionLine(line: string, companionName?: string, enabled: boolean = true) {
  if (!enabled || typeof window === "undefined" || !("speechSynthesis" in window)) return;
  try {
    window.speechSynthesis.cancel();
    const utterance = new SpeechSynthesisUtterance(line);
    const name = (companionName ?? "").toLowerCase();
    if (name.includes("pip")) {
      utterance.pitch = 1.35;
      utterance.rate = 1.15;
    } else if (name.includes("byte")) {
      utterance.pitch = 0.85;
      utterance.rate = 0.95;
    } else if (name.includes("zara") || name.includes("cat")) {
      utterance.pitch = 1.25;
      utterance.rate = 1.1;
    } else if (name.includes("kodo") || name.includes("owl")) {
      utterance.pitch = 0.8;
      utterance.rate = 0.9;
    } else {
      utterance.pitch = 1.1;
      utterance.rate = 1.05;
    }
    window.speechSynthesis.speak(utterance);
  } catch {
    // Graceful fallback if speech synthesis is disabled in browser
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

  const speak = useCallback(
    (line: string, companionName?: string) => {
      if (enabled) speakCompanionLine(line, companionName, enabled);
    },
    [enabled],
  );

  return { enabled, toggle, cue, speak };
}
