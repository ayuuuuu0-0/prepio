/**
 * Answer drafts: what the learner has built so far for one exercise, before Check.
 * These helpers only shape input for the server; grading happens on the server.
 */
import type { Answer, ClientStep } from "./types";

/** Draft is the in-progress answer for one step. */
export type Draft =
  | { type: "mcq"; choice: number | null }
  | { type: "true_false"; value: boolean | null }
  /** slots[i] is the bank index placed in blank i+1, or null while empty. */
  | { type: "fill_blank"; slots: (number | null)[] }
  /** order lists item indexes in the order the learner tapped them. */
  | { type: "arrange"; order: number[] };

/** GRADABLE lists the step types the player can render and submit. */
export const GRADABLE = ["mcq", "true_false", "fill_blank", "arrange"] as const;

/** emptyDraft starts a draft for a step, or returns null if the player can't render the step. */
export function emptyDraft(step: ClientStep): Draft | null {
  switch (step.type) {
    case "mcq":
      return step.mcq ? { type: "mcq", choice: null } : null;
    case "true_false":
      return step.true_false ? { type: "true_false", value: null } : null;
    case "fill_blank":
      return step.fill_blank ? { type: "fill_blank", slots: Array(step.fill_blank.blanks).fill(null) } : null;
    case "arrange":
      return step.arrange ? { type: "arrange", order: [] } : null;
    default:
      return null;
  }
}

/** isReady says whether the draft is complete enough to Check. */
export function isReady(step: ClientStep, d: Draft | null): boolean {
  if (!d) return false;
  switch (d.type) {
    case "mcq":
      return d.choice !== null;
    case "true_false":
      return d.value !== null;
    case "fill_blank":
      return d.slots.length > 0 && d.slots.every((s) => s !== null);
    case "arrange":
      return d.order.length === (step.arrange?.items.length ?? -1);
  }
}

/** toAnswer converts a ready draft into the wire answer. */
export function toAnswer(step: ClientStep, d: Draft): Answer {
  switch (d.type) {
    case "mcq":
      return { choice: d.choice ?? undefined };
    case "true_false":
      return { value: d.value ?? undefined };
    case "fill_blank": {
      const bank = step.fill_blank?.bank ?? [];
      return { blanks: d.slots.map((s) => (s === null ? "" : bank[s])) };
    }
    case "arrange":
      return { order: [...d.order] };
  }
}

/** placeChip puts a bank chip into the first empty blank. A chip already placed is not placed twice. */
export function placeChip(slots: (number | null)[], bankIndex: number): (number | null)[] {
  if (slots.includes(bankIndex)) return slots;
  const i = slots.indexOf(null);
  if (i === -1) return slots;
  const next = [...slots];
  next[i] = bankIndex;
  return next;
}

/** clearSlot empties one blank, returning its chip to the bank. */
export function clearSlot(slots: (number | null)[], slot: number): (number | null)[] {
  if (slot < 0 || slot >= slots.length || slots[slot] === null) return slots;
  const next = [...slots];
  next[slot] = null;
  return next;
}

/** toggleArrange adds an item to the end of the order, or takes it back out if already placed. */
export function toggleArrange(order: number[], item: number): number[] {
  return order.includes(item) ? order.filter((i) => i !== item) : [...order, item];
}

export type CodePart = { kind: "text"; text: string } | { kind: "blank"; blank: number };

/** splitCode splits fill-in-the-blank code into text and blanks (___1___ → blank 0). */
export function splitCode(code: string): CodePart[] {
  const parts: CodePart[] = [];
  const re = /___(\d+)___/g;
  let last = 0;
  for (let m = re.exec(code); m !== null; m = re.exec(code)) {
    if (m.index > last) parts.push({ kind: "text", text: code.slice(last, m.index) });
    parts.push({ kind: "blank", blank: Number(m[1]) - 1 });
    last = m.index + m[0].length;
  }
  if (last < code.length) parts.push({ kind: "text", text: code.slice(last) });
  return parts;
}

/**
 * correctAnswerLines describes the server's correct answer (sent only after a miss) in words
 * for the feedback tray. Returns an empty list when there is nothing to show.
 */
export function correctAnswerLines(step: ClientStep, correct: Answer | undefined): string[] {
  if (!correct) return [];
  switch (step.type) {
    case "mcq":
      return correct.choice !== undefined && step.mcq?.options[correct.choice] !== undefined
        ? [step.mcq.options[correct.choice]]
        : [];
    case "true_false":
      return correct.value === undefined ? [] : [correct.value ? "True" : "False"];
    case "fill_blank":
      return (correct.blanks ?? []).map((b, i) => (correct.blanks!.length > 1 ? `Blank ${i + 1}: ${b}` : b));
    case "arrange": {
      const items = step.arrange?.items ?? [];
      return (correct.order ?? []).filter((i) => items[i] !== undefined).map((i, pos) => `${pos + 1}. ${items[i]}`);
    }
    default:
      return [];
  }
}
