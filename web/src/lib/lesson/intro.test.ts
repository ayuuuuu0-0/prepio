import { describe, expect, it } from "vitest";
import { beatProgress, mediaKind, revealCount, splitEmphasis, totalDurationMs, words } from "./intro";

describe("words", () => {
  it("keeps spacing attached so the reassembled text is unchanged", () => {
    const parts = words("Reading from memory is fast.");
    expect(parts).toEqual(["Reading ", "from ", "memory ", "is ", "fast."]);
    expect(parts.join("")).toBe("Reading from memory is fast.");
  });

  it("keeps leading space so text after an emphasised phrase still reads correctly", () => {
    expect(words(" is fast.")).toEqual([" is ", "fast."]);
    expect(words(" is fast.").join("")).toBe(" is fast.");
  });

  it("handles empty and padded text", () => {
    expect(words("")).toEqual([]);
    expect(words("   ")).toEqual([]);
  });
});

describe("revealCount", () => {
  it("reveals nothing at the start and everything once the reveal window has passed", () => {
    expect(revealCount(0, 4000, 5)).toBe(0);
    expect(revealCount(10, 4000, 5)).toBe(1);
    expect(revealCount(4000, 4000, 5)).toBe(5);
  });

  it("reveals words progressively and never more than exist", () => {
    const seen = [0, 300, 600, 900, 1200, 1500, 1800, 2200].map((t) => revealCount(t, 4000, 8));
    for (let i = 1; i < seen.length; i++) expect(seen[i]).toBeGreaterThanOrEqual(seen[i - 1]);
    expect(Math.max(...seen)).toBe(8);
  });

  it("caps the reveal time so a long beat does not make the first words crawl", () => {
    expect(revealCount(2600, 60000, 10)).toBe(10);
  });

  it("is safe for degenerate input", () => {
    expect(revealCount(100, 4000, 0)).toBe(0);
    expect(revealCount(-50, 4000, 4)).toBe(0);
    expect(revealCount(100, 0, 3)).toBe(3);
  });
});

describe("beatProgress", () => {
  it("is a clamped fraction", () => {
    expect(beatProgress(0, 4000)).toBe(0);
    expect(beatProgress(1000, 4000)).toBe(0.25);
    expect(beatProgress(9000, 4000)).toBe(1);
    expect(beatProgress(-5, 4000)).toBe(0);
    expect(beatProgress(5, 0)).toBe(1);
  });
});

describe("splitEmphasis", () => {
  it("splits around the emphasised phrase", () => {
    expect(splitEmphasis({ text: "Reading from memory is fast.", emphasis: "memory" })).toEqual({
      before: "Reading from ",
      emphasis: "memory",
      after: " is fast.",
    });
  });

  it("falls back to plain text when there is no emphasis or it is not in the text", () => {
    expect(splitEmphasis({ text: "Hello", emphasis: undefined })).toEqual({ before: "Hello", emphasis: "", after: "" });
    expect(splitEmphasis({ text: "Hello", emphasis: "nope" })).toEqual({ before: "Hello", emphasis: "", after: "" });
  });
});

describe("totalDurationMs", () => {
  it("sums beat durations and ignores negative ones", () => {
    expect(totalDurationMs([{ text: "a", duration_ms: 3000 }, { text: "b", duration_ms: 4500 }])).toBe(7500);
    expect(totalDurationMs([{ text: "a", duration_ms: -5 }])).toBe(0);
    expect(totalDurationMs([])).toBe(0);
  });
});

describe("mediaKind", () => {
  it("accepts https media and site-relative paths by extension", () => {
    expect(mediaKind("https://cdn.example.com/intro.mp4")).toBe("video");
    expect(mediaKind("https://cdn.example.com/intro.webm?v=2")).toBe("video");
    expect(mediaKind("https://cdn.example.com/diagram.png")).toBe("image");
    expect(mediaKind("/media/diagram.SVG")).toBe("image");
  });

  it("refuses anything that is not plainly https or site-relative media", () => {
    expect(mediaKind(undefined)).toBeNull();
    expect(mediaKind("")).toBeNull();
    expect(mediaKind("http://cdn.example.com/a.png")).toBeNull();
    expect(mediaKind("javascript:alert(1)")).toBeNull();
    expect(mediaKind("data:image/png;base64,AAAA")).toBeNull();
    expect(mediaKind("//evil.example.com/a.png")).toBeNull();
    expect(mediaKind("https://cdn.example.com/page.html")).toBeNull();
  });
});
