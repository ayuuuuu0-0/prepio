"use client";

import Link from "next/link";
import { useParams, useRouter } from "next/navigation";
import { useCallback, useEffect, useMemo, useReducer, useRef, useState } from "react";
import { api, ApiRequestError, Profile } from "@/lib/api";
import { currentStepId, nextTry, progress, reduce, start } from "@/lib/lesson/session";
import type { AnswerResult, AttemptData, ClientStep, CompletionData } from "@/lib/lesson/types";
import { CompanionHero } from "@/components/game/CompanionHero";
import { GameButton } from "@/components/game/GameButton";
import { PlayerShell } from "@/components/lesson/PlayerShell";
import { IntroStep } from "@/components/lesson/IntroStep";
import { McqStep } from "@/components/lesson/McqStep";
import { FeedbackTray } from "@/components/lesson/FeedbackTray";
import { Celebration } from "@/components/lesson/Celebration";

type Load =
  | { status: "loading" }
  | { status: "ready"; attempt: AttemptData }
  | { status: "error"; code: string; message: string };

const MAX_REWARD_POLLS = 4;

export default function LessonPlayerPage() {
  const router = useRouter();
  const { lessonId } = useParams<{ lessonId: string }>();

  const [load, setLoad] = useState<Load>({ status: "loading" });
  const [profile, setProfile] = useState<Profile | null>(null);
  const [attempt, setAttempt] = useState<AttemptData | null>(null);

  const bootstrap = useCallback(async () => {
    setLoad({ status: "loading" });
    try {
      const ok = await api.ensureSession();
      if (!ok) {
        router.replace("/login");
        return;
      }
      const [data, me] = await Promise.all([api.startAttempt(lessonId), api.getProfile().catch(() => null)]);
      setProfile(me);
      setAttempt(data);
      setLoad({ status: "ready", attempt: data });
    } catch (err) {
      if (err instanceof ApiRequestError) {
        setLoad({ status: "error", code: err.code, message: err.message });
      } else {
        setLoad({ status: "error", code: "network", message: "Can't reach Prepio right now." });
      }
    }
  }, [lessonId, router]);

  useEffect(() => {
    bootstrap();
  }, [bootstrap]);

  if (load.status === "loading") {
    return (
      <div className="game-bg-challenge flex min-h-dvh flex-col items-center justify-center gap-4 px-6">
        <CompanionHero name={profile?.companion?.name} species={profile?.companion?.species} size="md" />
        <p className="font-mono animate-pulse text-sm font-semibold" style={{ color: "#7C6EF5" }} role="status">
          Getting your lesson ready…
        </p>
      </div>
    );
  }

  if (load.status === "error" || !attempt) {
    return <LoadError error={load.status === "error" ? load : { code: "network", message: "" }} onRetry={bootstrap} />;
  }

  return <Player attempt={attempt} profile={profile} onResync={bootstrap} />;
}

function LoadError({ error, onRetry }: { error: { code: string; message: string }; onRetry: () => void }) {
  const locked = error.code === "lesson_locked";
  const missing = error.code === "lesson_not_found";
  const title = locked ? "This lesson is still locked" : missing ? "We can't find that lesson" : "Something went wrong";
  const body = locked
    ? "Finish the lesson before it on the journey to unlock this one."
    : missing
      ? "It may have been retired. Head back to the journey to pick another."
      : "Check your connection and try again. Your progress is safe.";

  return (
    <div className="game-bg-challenge flex min-h-dvh flex-col items-center justify-center px-6 text-center">
      <p className="text-5xl" aria-hidden>
        {locked ? "🔒" : "🧭"}
      </p>
      <h1 className="font-display mt-4 text-xl font-extrabold" style={{ color: "#E8EAED" }}>
        {title}
      </h1>
      <p className="mt-2 max-w-sm text-sm leading-relaxed" style={{ color: "#8B92A8" }} role="alert">
        {body}
      </p>
      <div className="mt-6 flex w-full max-w-xs flex-col gap-3">
        {!locked && !missing && (
          <GameButton type="button" onClick={onRetry}>
            Try again
          </GameButton>
        )}
        <Link href="/journey">
          <GameButton type="button" variant={locked || missing ? "primary" : "ghost"}>
            Back to journey
          </GameButton>
        </Link>
      </div>
    </div>
  );
}

type Graded = { result: AnswerResult; chosen: number };

function Player({
  attempt,
  profile,
  onResync,
}: {
  attempt: AttemptData;
  profile: Profile | null;
  onResync: () => void;
}) {
  const router = useRouter();
  const [session, dispatch] = useReducer(reduce, undefined, () => start(attempt.steps, attempt.progress));
  const [selected, setSelected] = useState<number | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [graded, setGraded] = useState<Graded | null>(null);
  const [problem, setProblem] = useState<string | null>(null);
  const [confirmingExit, setConfirmingExit] = useState(false);
  const [completion, setCompletion] = useState<CompletionData | null>(null);
  const [completeFailed, setCompleteFailed] = useState(false);
  const [streak, setStreak] = useState<number | null>(null);
  const polls = useRef(0);

  const stepsById = useMemo(() => new Map(attempt.steps.map((s) => [s.id, s])), [attempt.steps]);
  const intro = attempt.steps.find((s) => s.type === "intro");
  const stepId = currentStepId(session);
  const step: ClientStep | undefined = stepId ? stepsById.get(stepId) : undefined;
  const companionName = profile?.companion?.name;
  const companionSpecies = profile?.companion?.species;

  const reaction = graded ? (graded.result.correct ? "correct" : "wrong") : "idle";

  const check = useCallback(async () => {
    if (!step || step.type !== "mcq" || selected === null || submitting || session.phase !== "answering") return;
    setSubmitting(true);
    setProblem(null);
    try {
      const result = await api.answerStep(attempt.attempt_id, step.id, nextTry(session, step.id), { choice: selected });
      setGraded({ result, chosen: selected });
      dispatch({ type: "graded", stepId: step.id, correct: result.correct });
    } catch (err) {
      if (err instanceof ApiRequestError && err.code === "invalid_try") {
        // Our try counter drifted from the server (for example after a reload in another tab): resync.
        onResync();
        return;
      }
      setProblem(
        err instanceof ApiRequestError && err.status < 500 ? err.message : "Couldn't check that. Check your connection and try again.",
      );
    } finally {
      setSubmitting(false);
    }
  }, [step, selected, submitting, session, attempt.attempt_id, onResync]);

  const next = useCallback(() => {
    setGraded(null);
    setSelected(null);
    setProblem(null);
    dispatch({ type: "continue" });
  }, []);

  // Enter checks the selected answer. (On the feedback tray, Enter presses its focused Continue button.)
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key !== "Enter" || confirmingExit) return;
      const target = e.target as HTMLElement | null;
      if (target?.getAttribute("role") === "radio" || target?.tagName === "BUTTON") return;
      if (session.phase === "answering") check();
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [check, session.phase, confirmingExit]);

  // When every graded step is solved, complete the attempt; the server decides if it really is complete.
  useEffect(() => {
    if (session.phase !== "finished" || completion || completeFailed) return;
    let cancelled = false;
    (async () => {
      try {
        const done = await api.completeAttempt(attempt.attempt_id);
        if (!cancelled) setCompletion(done);
      } catch {
        if (!cancelled) setCompleteFailed(true);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [session.phase, completion, completeFailed, attempt.attempt_id]);

  // Progress applies the completion asynchronously; re-ask (idempotent) until the rewards are in.
  useEffect(() => {
    if (!completion?.rewards_pending || polls.current >= MAX_REWARD_POLLS) return;
    const t = setTimeout(async () => {
      polls.current += 1;
      try {
        setCompletion(await api.completeAttempt(attempt.attempt_id));
      } catch {
        /* keep what we have; the celebration still shows accuracy */
      }
    }, 1200);
    return () => clearTimeout(t);
  }, [completion, attempt.attempt_id]);

  useEffect(() => {
    if (!completion) return;
    api
      .getDashboardHome()
      .then((home) => setStreak(home.streak.current_streak))
      .catch(() => setStreak(null));
  }, [completion?.attempt_id]); // eslint-disable-line react-hooks/exhaustive-deps

  const backToJourney = () => router.push(completion ? `/journey?done=${completion.node_id}` : "/journey");

  // Finished: celebration and summary.
  if (session.phase === "finished") {
    return (
      <div className="game-bg-challenge relative flex min-h-dvh flex-col">
        <div className="relative z-10 mx-auto flex w-full max-w-2xl flex-1 flex-col px-4 pb-8 pt-8">
          {completion ? (
            <Celebration
              completion={completion}
              streak={streak}
              companionName={companionName}
              companionSpecies={companionSpecies}
              onBackToJourney={backToJourney}
            />
          ) : completeFailed ? (
            <div className="flex flex-1 flex-col items-center justify-center text-center">
              <h1 className="font-display text-xl font-extrabold" style={{ color: "#E8EAED" }}>
                We couldn&apos;t finish saving that
              </h1>
              <p className="mt-2 text-sm" style={{ color: "#8B92A8" }} role="alert">
                Your answers are safe. Try again.
              </p>
              <div className="mt-6 w-full max-w-xs">
                <GameButton type="button" onClick={() => setCompleteFailed(false)}>
                  Try again
                </GameButton>
              </div>
            </div>
          ) : (
            <div className="flex flex-1 items-center justify-center">
              <p className="font-mono animate-pulse text-sm font-semibold" style={{ color: "#7C6EF5" }} role="status">
                Finishing up…
              </p>
            </div>
          )}
        </div>
      </div>
    );
  }

  const tray =
    session.phase === "feedback" && graded ? (
      <FeedbackTray
        result={graded.result}
        correctOptionText={
          step?.mcq && graded.result.correct_answer?.choice !== undefined
            ? step.mcq.options[graded.result.correct_answer.choice]
            : undefined
        }
        companionName={companionName}
        onContinue={next}
        continueLabel={session.queue.length === 1 && graded.result.correct ? "Finish" : "Continue"}
      />
    ) : undefined;

  return (
    <PlayerShell
      progress={progress(session)}
      combo={session.combo}
      confirmingExit={confirmingExit}
      onRequestExit={() => setConfirmingExit(true)}
      onCancelExit={() => setConfirmingExit(false)}
      onConfirmExit={() => router.push("/journey")}
      tray={tray}
      companion={<CompanionHero name={companionName} species={companionSpecies} size="sm" reaction={reaction} />}
    >
      {session.phase === "intro" && intro?.intro ? (
        <IntroStep beats={intro.intro.beats} onDone={() => dispatch({ type: "finishIntro" })} />
      ) : step?.type === "mcq" && step.mcq ? (
        <>
          <McqStep
            prompt={step.mcq.prompt}
            options={step.mcq.options}
            selected={selected}
            result={
              graded
                ? { correct: graded.result.correct, chosen: graded.chosen, correctIndex: graded.result.correct_answer?.choice }
                : null
            }
            disabled={submitting || session.phase !== "answering"}
            onSelect={setSelected}
          />
          {session.phase === "answering" && (
            <div className="mt-auto pt-6">
              {problem && (
                <p className="mb-3 text-center text-sm font-semibold" style={{ color: "#F87171" }} role="alert">
                  {problem}
                </p>
              )}
              <GameButton type="button" onClick={check} disabled={selected === null || submitting}>
                {submitting ? "Checking…" : "Check"}
              </GameButton>
            </div>
          )}
        </>
      ) : (
        <div className="flex flex-1 flex-col items-center justify-center text-center">
          <h1 className="font-display text-xl font-extrabold" style={{ color: "#E8EAED" }}>
            This exercise isn&apos;t available in this version
          </h1>
          <p className="mt-2 max-w-sm text-sm" style={{ color: "#8B92A8" }}>
            Your progress is saved. Leave the lesson and come back after the next update.
          </p>
          <div className="mt-6 w-full max-w-xs">
            <GameButton type="button" onClick={() => router.push("/journey")}>
              Back to journey
            </GameButton>
          </div>
        </div>
      )}
    </PlayerShell>
  );
}
