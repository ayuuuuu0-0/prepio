import { describe, expect, it } from "vitest";
import { achievementIcon, earnedLabel } from "./achievements";

describe("achievementIcon", () => {
  it("has an icon for known achievements and a fallback for new ones", () => {
    expect(achievementIcon("first-steps")).toBe("🌱");
    expect(achievementIcon("something-new")).toBe("🏅");
  });
});

describe("earnedLabel", () => {
  it("formats the earned date", () => {
    expect(earnedLabel("2026-10-08T09:00:00Z", "en-GB")).toBe("Earned 8 Oct 2026");
  });
  it("is empty when not earned or unreadable", () => {
    expect(earnedLabel(undefined)).toBe("");
    expect(earnedLabel("nope")).toBe("");
  });
});
