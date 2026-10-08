import { describe, expect, it } from "vitest";
import {
  clearSlot,
  correctAnswerLines,
  emptyDraft,
  isReady,
  placeChip,
  promptOf,
  recordMistake,
  splitCode,
  toAnswer,
  toggleArrange,
  type Draft,
} from "./answers";
import type { ClientStep } from "./types";

const mcq: ClientStep = { id: "q1", type: "mcq", position: 1, mcq: { prompt: "P", options: ["a", "b", "c"] } };
const tf: ClientStep = { id: "q2", type: "true_false", position: 2, true_false: { statement: "S" } };
const fill: ClientStep = {
  id: "q3",
  type: "fill_blank",
  position: 3,
  fill_blank: { code: "SET ___1___ EX ___2___", blanks: 2, bank: ["key", "60", "PX", "value"] },
};
const arrange: ClientStep = {
  id: "q4",
  type: "arrange",
  position: 4,
  arrange: { prompt: "Order", items: ["write", "read", "ack"] },
};
const prose: ClientStep = { id: "q5", type: "prose", position: 5, prose: { prompt: "Explain", min_chars: 100 } };
const unknown = { id: "q6", type: "type_line", position: 6 } as unknown as ClientStep;

describe("emptyDraft", () => {
  it("starts an empty draft for each supported type", () => {
    expect(emptyDraft(mcq)).toEqual({ type: "mcq", choice: null });
    expect(emptyDraft(tf)).toEqual({ type: "true_false", value: null });
    expect(emptyDraft(fill)).toEqual({ type: "fill_blank", slots: [null, null] });
    expect(emptyDraft(arrange)).toEqual({ type: "arrange", order: [] });
  });
  it("starts a prose draft empty", () => {
    expect(emptyDraft(prose)).toEqual({ type: "prose", text: "" });
  });
  it("returns null for steps the player can't render, or a missing payload", () => {
    expect(emptyDraft(unknown)).toBeNull();
    expect(emptyDraft({ id: "x", type: "mcq", position: 1 })).toBeNull();
  });
});

describe("isReady and toAnswer", () => {
  it("mcq needs a choice", () => {
    expect(isReady(mcq, { type: "mcq", choice: null })).toBe(false);
    const d: Draft = { type: "mcq", choice: 2 };
    expect(isReady(mcq, d)).toBe(true);
    expect(toAnswer(mcq, d)).toEqual({ choice: 2 });
  });
  it("true_false needs a value, and false is a value", () => {
    expect(isReady(tf, { type: "true_false", value: null })).toBe(false);
    const d: Draft = { type: "true_false", value: false };
    expect(isReady(tf, d)).toBe(true);
    expect(toAnswer(tf, d)).toEqual({ value: false });
  });
  it("fill_blank needs every blank filled and sends the chip text", () => {
    expect(isReady(fill, { type: "fill_blank", slots: [0, null] })).toBe(false);
    const d: Draft = { type: "fill_blank", slots: [0, 1] };
    expect(isReady(fill, d)).toBe(true);
    expect(toAnswer(fill, d)).toEqual({ blanks: ["key", "60"] });
  });
  it("arrange needs every item placed and sends item indexes in order", () => {
    expect(isReady(arrange, { type: "arrange", order: [1, 0] })).toBe(false);
    const d: Draft = { type: "arrange", order: [1, 0, 2] };
    expect(isReady(arrange, d)).toBe(true);
    expect(toAnswer(arrange, d)).toEqual({ order: [1, 0, 2] });
  });
  it("prose needs at least min_chars of real text, ignoring surrounding spaces", () => {
    expect(isReady(prose, { type: "prose", text: "x".repeat(99) })).toBe(false);
    expect(isReady(prose, { type: "prose", text: "   " + "x".repeat(99) + "   " })).toBe(false);
    const d: Draft = { type: "prose", text: "y".repeat(100) };
    expect(isReady(prose, d)).toBe(true);
    expect(toAnswer(prose, d)).toEqual({ text: "y".repeat(100) });
  });
  it("nothing is ready without a draft", () => {
    expect(isReady(mcq, null)).toBe(false);
  });
});

describe("fill_blank chips", () => {
  it("places chips into the first empty blank, never twice, and stops when full", () => {
    let s = placeChip([null, null], 2);
    expect(s).toEqual([2, null]);
    expect(placeChip(s, 2)).toBe(s);
    s = placeChip(s, 0);
    expect(s).toEqual([2, 0]);
    expect(placeChip(s, 1)).toBe(s);
  });
  it("clearing a blank frees it for the next chip", () => {
    const s = clearSlot([2, 0], 0);
    expect(s).toEqual([null, 0]);
    expect(placeChip(s, 3)).toEqual([3, 0]);
    expect(clearSlot(s, 0)).toBe(s);
    expect(clearSlot(s, 9)).toBe(s);
  });
});

describe("toggleArrange", () => {
  it("appends untapped items and removes tapped ones", () => {
    expect(toggleArrange([], 2)).toEqual([2]);
    expect(toggleArrange([2, 0], 1)).toEqual([2, 0, 1]);
    expect(toggleArrange([2, 0, 1], 0)).toEqual([2, 1]);
  });
});

describe("splitCode", () => {
  it("splits code around numbered blanks", () => {
    expect(splitCode("SET ___1___ EX ___2___")).toEqual([
      { kind: "text", text: "SET " },
      { kind: "blank", blank: 0 },
      { kind: "text", text: " EX " },
      { kind: "blank", blank: 1 },
    ]);
  });
  it("keeps code without blanks as one text part", () => {
    expect(splitCode("no blanks")).toEqual([{ kind: "text", text: "no blanks" }]);
  });
});

describe("correctAnswerLines", () => {
  it("names the correct answer for each type", () => {
    expect(correctAnswerLines(mcq, { choice: 1 })).toEqual(["b"]);
    expect(correctAnswerLines(tf, { value: true })).toEqual(["True"]);
    expect(correctAnswerLines(fill, { blanks: ["key", "60"] })).toEqual(["Blank 1: key", "Blank 2: 60"]);
    expect(correctAnswerLines(arrange, { order: [1, 0, 2] })).toEqual(["1. read", "2. write", "3. ack"]);
  });
  it("shows a single blank without a label", () => {
    const one: ClientStep = { ...fill, fill_blank: { code: "___1___", blanks: 1, bank: ["x", "y"] } };
    expect(correctAnswerLines(one, { blanks: ["x"] })).toEqual(["x"]);
  });
  it("shows nothing without a correct answer (a correct try)", () => {
    expect(correctAnswerLines(mcq, undefined)).toEqual([]);
  });
});

describe("mistakes review", () => {
  it("describes what each step asked", () => {
    expect(promptOf(mcq)).toBe("P");
    expect(promptOf(tf)).toBe("True or false: S");
    expect(promptOf(fill)).toBe("Fill in the blanks: SET [1] EX [2]");
    expect(promptOf(arrange)).toBe("Order");
    expect(promptOf(prose)).toBe("Explain");
    expect(promptOf(unknown)).toBe("");
  });
  it("keeps only the first miss of each step, in order", () => {
    let list = recordMistake([], { stepId: "a", prompt: "A", correct: ["x"] });
    list = recordMistake(list, { stepId: "b", prompt: "B", correct: ["y"] });
    list = recordMistake(list, { stepId: "a", prompt: "A again", correct: ["z"] });
    expect(list.map((m) => [m.stepId, m.prompt])).toEqual([
      ["a", "A"],
      ["b", "B"],
    ]);
  });
});
