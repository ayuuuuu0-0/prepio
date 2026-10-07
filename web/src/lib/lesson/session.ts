/**
 * Pure sequencing for the lesson player.
 *
 * The client only decides presentation order (including re-queueing missed
 * steps). Correctness, completion, XP, and mastery are all server decisions;
 * nothing here can mark a step correct on its own.
 */
import type { ClientStep, StepProgress } from "./types";

export type Phase = "intro" | "answering" | "feedback" | "finished";

export type Session = {
  /** Step ids still to present, current first. The intro is never in the queue. */
  queue: string[];
  /** Graded steps the server has accepted as correct. */
  solved: string[];
  /** Tries the server has recorded per graded step. */
  tries: Record<string, number>;
  /** Total graded steps in the lesson. */
  total: number;
  /** Consecutive steps solved on their first try (UI flair only). */
  combo: number;
  /** Steps answered correctly on try 1, for the live accuracy hint. */
  firstTrySolved: number;
  phase: Phase;
  /** Result of the step currently in the feedback tray. */
  lastCorrect: boolean | null;
  hasIntro: boolean;
};

export type Action =
  | { type: "finishIntro" }
  | { type: "graded"; stepId: string; correct: boolean }
  | { type: "continue" };

const isGraded = (s: ClientStep) => s.type !== "intro";

/** start builds the initial session, resuming from the server's per-step progress. */
export function start(steps: ClientStep[], progress: StepProgress[] = []): Session {
  const graded = steps.filter(isGraded);
  const done = new Set(progress.filter((p) => p.done).map((p) => p.step_id));
  const tries: Record<string, number> = {};
  for (const p of progress) tries[p.step_id] = p.tries;

  const queue = graded.map((s) => s.id).filter((id) => !done.has(id));
  const hasIntro = steps.some((s) => s.type === "intro");
  // An intro is only shown on a fresh attempt; a resumed one goes straight to the exercises.
  const showIntro = hasIntro && progress.length === 0;

  return {
    queue,
    solved: graded.map((s) => s.id).filter((id) => done.has(id)),
    tries,
    total: graded.length,
    combo: 0,
    firstTrySolved: 0,
    phase: queue.length === 0 ? "finished" : showIntro ? "intro" : "answering",
    lastCorrect: null,
    hasIntro,
  };
}

/** nextTry is the 1-based try number to send for a step. */
export function nextTry(s: Session, stepId: string): number {
  return (s.tries[stepId] ?? 0) + 1;
}

export function currentStepId(s: Session): string | undefined {
  return s.queue[0];
}

/** progress is the share of graded steps solved, from 0 to 1. */
export function progress(s: Session): number {
  return s.total === 0 ? 1 : s.solved.length / s.total;
}

export function reduce(s: Session, a: Action): Session {
  switch (a.type) {
    case "finishIntro":
      return s.phase === "intro" ? { ...s, phase: "answering" } : s;

    case "graded": {
      if (s.phase !== "answering" || s.queue[0] !== a.stepId) return s;
      const tries = { ...s.tries, [a.stepId]: (s.tries[a.stepId] ?? 0) + 1 };
      if (!a.correct) {
        return { ...s, tries, combo: 0, phase: "feedback", lastCorrect: false };
      }
      const first = tries[a.stepId] === 1;
      return {
        ...s,
        tries,
        solved: [...s.solved, a.stepId],
        combo: first ? s.combo + 1 : 0,
        firstTrySolved: s.firstTrySolved + (first ? 1 : 0),
        phase: "feedback",
        lastCorrect: true,
      };
    }

    case "continue": {
      if (s.phase !== "feedback") return s;
      const [head, ...rest] = s.queue;
      // A missed step goes to the back of the line: it returns at the end of the lesson.
      const queue = s.lastCorrect ? rest : [...rest, head];
      return {
        ...s,
        queue,
        phase: queue.length === 0 ? "finished" : "answering",
        lastCorrect: null,
      };
    }
  }
}
