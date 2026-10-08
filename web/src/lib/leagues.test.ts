import { describe, expect, it } from "vitest";
import { resultMessage, tierTint, timeLeft, zoneBoundaries, zoneLabel } from "./leagues";
import type { LeagueResult } from "./api";

const now = new Date("2026-10-07T12:00:00Z");

describe("timeLeft", () => {
  it("shows days and hours when more than a day is left", () => {
    expect(timeLeft("2026-10-12T00:00:00Z", now)).toBe("4d 12h");
  });
  it("shows hours and minutes on the last day", () => {
    expect(timeLeft("2026-10-07T17:30:00Z", now)).toBe("5h 30m");
  });
  it("shows minutes in the last hour", () => {
    expect(timeLeft("2026-10-07T12:12:00Z", now)).toBe("12m");
  });
  it("says ending now at or past the end, or for a bad date", () => {
    expect(timeLeft("2026-10-07T12:00:30Z", now)).toBe("ending now");
    expect(timeLeft("2026-10-01T00:00:00Z", now)).toBe("ending now");
    expect(timeLeft("not a date", now)).toBe("ending now");
  });
});

describe("zones", () => {
  it("labels every zone in words", () => {
    expect(zoneLabel("promote")).toBe("Promotion zone");
    expect(zoneLabel("stay")).toBe("Safe");
    expect(zoneLabel("demote")).toBe("Drop zone");
  });
  it("draws dividers where the zone changes", () => {
    expect(zoneBoundaries(["promote", "promote", "stay", "stay", "demote"])).toEqual([1, 3]);
    expect(zoneBoundaries(["promote", "promote"])).toEqual([]);
    expect(zoneBoundaries([])).toEqual([]);
  });
});

describe("resultMessage", () => {
  const base: LeagueResult = {
    week_start: "2026-09-28",
    rank: 2,
    cohort_size: 20,
    outcome: "promoted",
    from_tier: { index: 0, slug: "bronze", name: "Bronze" },
    to_tier: { index: 1, slug: "silver", name: "Silver" },
  };
  it("celebrates a promotion", () => {
    expect(resultMessage(base)).toBe("You finished #2 and moved up to Silver.");
  });
  it("frames a demotion as a chance, never a punishment", () => {
    const msg = resultMessage({ ...base, rank: 19, outcome: "demoted", to_tier: base.from_tier });
    expect(msg).toContain("win it back");
    expect(msg.toLowerCase()).not.toMatch(/fail|lost|demoted/);
  });
  it("acknowledges holding a place", () => {
    expect(resultMessage({ ...base, rank: 9, outcome: "stayed", to_tier: base.from_tier })).toBe(
      "You finished #9 and held your place in Bronze.",
    );
  });
});

describe("tierTint", () => {
  it("falls back to the accent for unknown tiers", () => {
    expect(tierTint({ slug: "gold" })).toBe("#F5B942");
    expect(tierTint({ slug: "mythril" })).toBe("#7C6EF5");
  });
});
