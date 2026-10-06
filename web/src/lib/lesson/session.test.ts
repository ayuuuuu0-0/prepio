import { describe, expect, it } from "vitest";
import { currentStepId, nextTry, progress, reduce, start, type Session } from "./session";
import type { ClientStep } from "./types";

const step = (id: string, type: ClientStep["type"] = "mcq"): ClientStep => ({ id, type, position: 0 });
const steps = [step("intro", "intro"), step("a"), step("b"), step("c")];

type Event = "intro" | "right" | "wrong";

/** run plays events: each right/wrong grades the current step and presses Continue. */
const run = (s: Session, ...events: Event[]): Session => {
  for (const e of events) {
    if (e === "intro") {
      s = reduce(s, { type: "finishIntro" });
      continue;
    }
    const id = currentStepId(s)!;
    s = reduce(s, { type: "graded", stepId: id, correct: e === "right" });
    s = reduce(s, { type: "continue" });
  }
  return s;
};

describe("start", () => {
  it("queues graded steps only and begins at the intro", () => {
    const s = start(steps);
    expect(s.queue).toEqual(["a", "b", "c"]);
    expect(s.total).toBe(3);
    expect(s.phase).toBe("intro");
  });

  it("goes straight to exercises when there is no intro", () => {
    expect(start([step("a"), step("b")]).phase).toBe("answering");
  });

  it("resumes: skips solved steps and the intro, keeps try counts", () => {
    const s = start(steps, [
      { step_id: "a", tries: 2, done: true },
      { step_id: "b", tries: 1, done: false },
    ]);
    expect(s.queue).toEqual(["b", "c"]);
    expect(s.solved).toEqual(["a"]);
    expect(s.phase).toBe("answering");
    expect(nextTry(s, "b")).toBe(2);
    expect(nextTry(s, "c")).toBe(1);
    expect(progress(s)).toBeCloseTo(1 / 3);
  });

  it("is finished when everything is already solved", () => {
    const s = start([step("a")], [{ step_id: "a", tries: 1, done: true }]);
    expect(s.phase).toBe("finished");
    expect(progress(s)).toBe(1);
  });
});

describe("answering", () => {
  it("advances through steps that are answered correctly", () => {
    const s = run(start(steps), "intro", "right", "right", "right");
    expect(s.phase).toBe("finished");
    expect(progress(s)).toBe(1);
    expect(s.firstTrySolved).toBe(3);
  });

  it("returns a missed step at the end of the lesson", () => {
    let s = run(start(steps), "intro", "wrong"); // a missed
    expect(s.queue).toEqual(["b", "c", "a"]);
    s = run(s, "right", "right");
    expect(currentStepId(s)).toBe("a");
    expect(s.phase).toBe("answering");
    s = run(s, "right");
    expect(s.phase).toBe("finished");
  });

  it("keeps re-queueing a step until it is solved, counting every try", () => {
    let s = run(start([step("a"), step("b")]), "wrong", "right", "wrong");
    expect(s.queue).toEqual(["a"]);
    expect(nextTry(s, "a")).toBe(3);
    s = run(s, "right");
    expect(s.phase).toBe("finished");
    expect(s.tries.a).toBe(3);
    expect(s.firstTrySolved).toBe(1);
  });

  it("a lone missed step comes straight back", () => {
    const s = run(start([step("a")]), "wrong");
    expect(s.queue).toEqual(["a"]);
    expect(s.phase).toBe("answering");
  });

  it("tracks the combo and resets it on a miss or a retry success", () => {
    let s = start([step("a"), step("b"), step("c"), step("d")]);
    s = run(s, "right", "right");
    expect(s.combo).toBe(2);
    s = run(s, "wrong"); // c missed
    expect(s.combo).toBe(0);
    s = run(s, "right"); // d, first try
    expect(s.combo).toBe(1);
    s = run(s, "right"); // c on its second try does not extend the combo
    expect(s.combo).toBe(0);
  });

  it("progress only counts solved steps, never missed ones", () => {
    const s = run(start(steps), "intro", "wrong", "right");
    expect(progress(s)).toBeCloseTo(1 / 3);
  });
});

describe("replaying the intro", () => {
  it("returns to the intro from an exercise without losing any progress, then back", () => {
    let s = run(start(steps), "intro", "right"); // a solved, now on b
    const before = { queue: s.queue, solved: s.solved, tries: s.tries };
    s = reduce(s, { type: "replayIntro" });
    expect(s.phase).toBe("intro");
    expect({ queue: s.queue, solved: s.solved, tries: s.tries }).toEqual(before);
    s = reduce(s, { type: "finishIntro" });
    expect(s.phase).toBe("answering");
    expect(currentStepId(s)).toBe("b");
  });

  it("is ignored when the lesson has no intro, during feedback, and once finished", () => {
    const noIntro = start([step("a"), step("b")]);
    expect(reduce(noIntro, { type: "replayIntro" })).toBe(noIntro);

    const feedback = reduce(start(steps), { type: "finishIntro" });
    const graded = reduce(feedback, { type: "graded", stepId: "a", correct: true });
    expect(graded.phase).toBe("feedback");
    expect(reduce(graded, { type: "replayIntro" })).toBe(graded);

    const done = run(start(steps), "intro", "right", "right", "right");
    expect(reduce(done, { type: "replayIntro" })).toBe(done);
  });
});

describe("guards", () => {
  it("ignores answers for a step that is not current", () => {
    const s = reduce(start([step("a"), step("b")]), { type: "graded", stepId: "b", correct: true });
    expect(s.solved).toEqual([]);
  });

  it("ignores continue before feedback and intro actions outside the intro", () => {
    const s = start([step("a")]);
    expect(reduce(s, { type: "continue" })).toBe(s);
    expect(reduce(s, { type: "finishIntro" })).toBe(s);
  });

  it("does not grade twice while feedback is showing", () => {
    let s = reduce(start([step("a"), step("b")]), { type: "graded", stepId: "a", correct: false });
    expect(s.phase).toBe("feedback");
    s = reduce(s, { type: "graded", stepId: "a", correct: true });
    expect(s.solved).toEqual([]);
  });
});
