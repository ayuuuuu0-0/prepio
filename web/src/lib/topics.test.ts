import { describe, expect, it } from "vitest";
import { coverageText, formatDelta, groupMasteryChanges, ringOffset, topicTheme } from "./topics";
import type { MasteryChange } from "./lesson/types";

const change = (over: Partial<MasteryChange>): MasteryChange => ({
  skill_slug: "s",
  skill_name: "Skill",
  topic_slug: "system-design",
  topic_name: "System Design",
  before: 0,
  after: 0,
  delta: 0,
  accuracy: 1,
  ...over,
});

describe("ringOffset", () => {
  it("draws nothing for not started, everything for 100", () => {
    expect(ringOffset(null, 200)).toBe(200);
    expect(ringOffset(100, 200)).toBe(0);
    expect(ringOffset(0, 200)).toBe(200);
    expect(ringOffset(25, 200)).toBe(150);
  });

  it("clamps out-of-range values", () => {
    expect(ringOffset(140, 200)).toBe(0);
    expect(ringOffset(-5, 200)).toBe(200);
  });
});

describe("coverageText", () => {
  it("never frames an unstarted topic as a failure", () => {
    expect(coverageText(0, 4)).toBe("Not started yet");
    expect(coverageText(1, 4)).toBe("1 of 4 skills started");
    expect(coverageText(4, 4)).toBe("All 4 skills started");
    expect(coverageText(0, 0)).toBe("No skills yet");
  });
});

describe("groupMasteryChanges", () => {
  it("groups skills under their topic and sums the movement", () => {
    const groups = groupMasteryChanges([
      change({ skill_slug: "caching", skill_name: "Caching", delta: 2 }),
      change({ skill_slug: "lb", skill_name: "Load Balancing", delta: 1 }),
      change({ skill_slug: "api", skill_name: "API Design", topic_slug: "backend-production", topic_name: "Backend & Production", delta: 3 }),
    ]);
    expect(groups).toHaveLength(2);
    expect(groups[0]).toMatchObject({ topicName: "System Design", total: 3 });
    expect(groups[0].skills.map((s) => s.name)).toEqual(["Caching", "Load Balancing"]);
    expect(groups[1]).toMatchObject({ topicName: "Backend & Production", total: 3 });
  });

  it("leaves out skills that did not move and handles skills with no topic", () => {
    const groups = groupMasteryChanges([
      change({ skill_slug: "still", delta: 0 }),
      change({ skill_slug: "loose", skill_name: "Loose", topic_slug: undefined, topic_name: undefined, delta: 2 }),
    ]);
    expect(groups).toHaveLength(1);
    expect(groups[0]).toMatchObject({ topicName: "Other skills", total: 2 });
  });

  it("returns nothing when nothing moved", () => {
    expect(groupMasteryChanges([])).toEqual([]);
    expect(groupMasteryChanges([change({ delta: 0 })])).toEqual([]);
  });
});

describe("formatDelta and topicTheme", () => {
  it("signs deltas explicitly", () => {
    expect(formatDelta(3)).toBe("+3");
    expect(formatDelta(-2)).toBe("−2");
    expect(formatDelta(0)).toBe("0");
  });

  it("falls back to a neutral theme for unknown topics", () => {
    expect(topicTheme("system-design").tint).toBe("#7C6EF5");
    expect(topicTheme("something-new").icon).toBe("◆");
  });
});
