"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { SubmitResponse } from "@/lib/api";
import { CompanionHero } from "./CompanionHero";
import { GameButton } from "./GameButton";

/** ChallengeResult is the full-screen result overlay after submission. */
export function ChallengeResult({
  result,
  companionName,
  companionSpecies,
  hasNext,
  onNext,
}: {
  result: SubmitResponse;
  companionName?: string;
  companionSpecies?: string;
  hasNext: boolean;
  onNext: () => void;
}) {
  const [xpDisplay, setXpDisplay] = useState(0);
  const [gemDisplay, setGemDisplay] = useState(0);

  useEffect(() => {
    if (!result.correct) return;
    const xpTarget = result.xp_awarded;
    const gemTarget = result.gems_awarded;
    let frame = 0;
    const steps = 24;
    const timer = setInterval(() => {
      frame++;
      setXpDisplay(Math.round((xpTarget * frame) / steps));
      setGemDisplay(Math.round((gemTarget * frame) / steps));
      if (frame >= steps) clearInterval(timer);
    }, 40);
    return () => clearInterval(timer);
  }, [result]);

  const accent = result.correct ? "#34D399" : "#F5B942";

  return (
    <div
      className="animate-result fixed inset-0 z-50 flex flex-col"
      style={{ background: "rgba(15,17,23,0.98)" }}
    >
      <div className="flex flex-1 flex-col items-center justify-center px-6 pb-32 pt-12">
        <CompanionHero
          name={companionName}
          species={companionSpecies}
          size="lg"
          reaction={result.correct ? "correct" : "wrong"}
        />

        <p className="font-display mt-6 text-4xl font-extrabold" style={{ color: accent }}>
          {result.correct ? "Solid." : "Not quite."}
        </p>
        <p className="mt-2 text-center text-lg font-semibold" style={{ color: "#8B92A8" }}>
          Score: <span className="font-mono text-2xl font-extrabold" style={{ color: accent }}>{result.score}%</span>
        </p>
        <p className="mt-2 max-w-md text-center font-medium" style={{ color: "#C8CCDA" }}>{result.feedback}</p>

        {result.readiness_delta > 0 && (
          <p className="font-mono mt-3 rounded-full px-4 py-2 text-sm font-bold" style={{ background: "rgba(96,165,250,0.12)", color: "#60A5FA" }}>
            Readiness +{result.readiness_delta}%
          </p>
        )}

        {result.correct && (
          <div className="mt-6 flex gap-4">
            <span className="animate-xp font-mono rounded-full px-5 py-2 text-lg font-bold" style={{ background: "rgba(96,165,250,0.12)", color: "#60A5FA" }}>
              +{xpDisplay} XP
            </span>
            <span className="animate-gems font-mono rounded-full px-5 py-2 text-lg font-bold" style={{ background: "rgba(52,211,153,0.12)", color: "#34D399" }}>
              +{gemDisplay} 💎
            </span>
          </div>
        )}

        <div className="mt-8 w-full max-w-md rounded-2xl p-5 text-left" style={{ background: "#1A1D27", border: "1px solid #2E3347" }}>
          {result.strengths.length > 0 && (
            <div className="mb-4">
              <p className="font-mono text-[10px] font-bold uppercase tracking-widest" style={{ color: "#4A5068" }}>Covered</p>
              <ul className="mt-2 space-y-1">
                {result.strengths.map((s) => (
                  <li key={s} className="text-sm" style={{ color: "#C8CCDA" }}>✓ {s}</li>
                ))}
              </ul>
            </div>
          )}
          {result.gaps.length > 0 && (
            <div>
              <p className="font-mono text-[10px] font-bold uppercase tracking-widest" style={{ color: "#4A5068" }}>Missing</p>
              <ul className="mt-2 space-y-1">
                {result.gaps.map((g) => (
                  <li key={g} className="text-sm" style={{ color: "#C8CCDA" }}>→ {g}</li>
                ))}
              </ul>
            </div>
          )}
        </div>
      </div>

      <div className="fixed bottom-0 left-0 right-0 px-6 pb-8 pt-4" style={{ background: "rgba(15,17,23,0.9)" }}>
        {hasNext ? (
          <GameButton type="button" onClick={onNext}>
            Next Question →
          </GameButton>
        ) : (
          <Link href="/dashboard">
            <GameButton type="button" variant="secondary">
              Back to dashboard
            </GameButton>
          </Link>
        )}
      </div>
    </div>
  );
}
