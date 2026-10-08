import { describe, expect, it } from "vitest";
import { CUES, readSoundPref } from "./sound";

describe("readSoundPref", () => {
  it("is off by default and only on when stored as on", () => {
    expect(readSoundPref(undefined)).toBe(false);
    expect(readSoundPref({ getItem: () => null })).toBe(false);
    expect(readSoundPref({ getItem: () => "off" })).toBe(false);
    expect(readSoundPref({ getItem: () => "on" })).toBe(true);
  });
  it("treats unreadable storage as off", () => {
    expect(
      readSoundPref({
        getItem: () => {
          throw new Error("blocked");
        },
      }),
    ).toBe(false);
  });
});

describe("CUES", () => {
  it("keeps every cue short and audible", () => {
    for (const notes of Object.values(CUES)) {
      expect(notes.length).toBeGreaterThan(0);
      for (const [freq, at, dur] of notes) {
        expect(freq).toBeGreaterThan(100);
        expect(at + dur).toBeLessThan(1);
      }
    }
  });
});
