/**
 * Pure timing helpers for the intro motion stage. They decide what is on screen at a given
 * moment; the renderer only draws it.
 */
import type { Beat } from "./types";

/** The share of a beat's duration spent revealing its words; the rest is time to read. */
const REVEAL_SHARE = 0.55;

/** Longest a reveal may take, so a long beat never makes the first words crawl. */
const MAX_REVEAL_MS = 2600;

/**
 * words splits a beat's text into the units that animate in. Whitespace stays attached to the
 * word it surrounds (including leading space), so the pieces always join back into the original.
 */
export function words(text: string): string[] {
  return text.match(/\s*\S+\s*/g) ?? [];
}

/** revealCount is how many of `total` words are visible `elapsedMs` into a beat of `durationMs`. */
export function revealCount(elapsedMs: number, durationMs: number, total: number): number {
  if (total <= 0) return 0;
  const revealMs = Math.min(MAX_REVEAL_MS, Math.max(1, durationMs * REVEAL_SHARE));
  const progress = Math.max(0, Math.min(1, elapsedMs / revealMs));
  return Math.min(total, Math.ceil(progress * total));
}

/** beatProgress is the fraction (0-1) of a beat that has played. */
export function beatProgress(elapsedMs: number, durationMs: number): number {
  if (durationMs <= 0) return 1;
  return Math.max(0, Math.min(1, elapsedMs / durationMs));
}

export type Segment = { before: string; emphasis: string; after: string };

/**
 * splitEmphasis finds the emphasised phrase inside the beat text. When the phrase is not
 * part of the text (an authoring mistake the validator also rejects) it falls back to plain text.
 */
export function splitEmphasis(beat: Pick<Beat, "text" | "emphasis">): Segment {
  const { text, emphasis } = beat;
  if (!emphasis) return { before: text, emphasis: "", after: "" };
  const at = text.indexOf(emphasis);
  if (at < 0) return { before: text, emphasis: "", after: "" };
  return { before: text.slice(0, at), emphasis, after: text.slice(at + emphasis.length) };
}

/** totalDurationMs is how long the whole intro plays when left alone. */
export function totalDurationMs(beats: Beat[]): number {
  return beats.reduce((sum, b) => sum + Math.max(0, b.duration_ms), 0);
}

/** mediaKind decides how an optional intro media URL is shown. Anything else is not rendered. */
export function mediaKind(url: string | undefined): "video" | "image" | null {
  if (!url) return null;
  let path: string;
  try {
    const parsed = new URL(url, "https://example.invalid");
    const siteRelative = url.startsWith("/") && !url.startsWith("//");
    if (parsed.protocol !== "https:" && !siteRelative) return null;
    if (url.startsWith("//")) return null;
    path = parsed.pathname.toLowerCase();
  } catch {
    return null;
  }
  if (/\.(mp4|webm)$/.test(path)) return "video";
  if (/\.(png|jpe?g|gif|webp|avif|svg)$/.test(path)) return "image";
  return null;
}
