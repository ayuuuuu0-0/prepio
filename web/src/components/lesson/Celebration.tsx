"use client";

import { useState } from "react";
import type { CompletionData } from "@/lib/lesson/types";
import { CompanionHero } from "@/components/game/CompanionHero";
import { SpeechBubble } from "@/components/game/SpeechBubble";
import { GameButton } from "@/components/game/GameButton";
import { formatDelta, groupMasteryChanges, topicTheme } from "@/lib/topics";

function message(accuracy: number, firstCompletion: boolean | undefined): string {
  if (firstCompletion === false) return "Nice practice run. Repetition is how it sticks.";
  if (accuracy >= 0.99) return "Flawless. Every answer first try.";
  if (accuracy >= 0.6) return "Solid work. The misses are what you learn from.";
  return "You finished it, and that's what counts. Try it again tomorrow and it will feel easier.";
}

function Chip({ label, value, tone, delay }: { label: string; value: string; tone: string; delay: string }) {
  return (
    <div
      className="animate-xp flex flex-col items-center rounded-2xl px-4 py-3"
      style={{ background: "#1A1D27", border: "1px solid #2E3347", animationDelay: delay }}
    >
      <span className="font-mono text-[10px] font-bold uppercase tracking-[0.2em]" style={{ color: "#8B92A8" }}>
        {label}
      </span>
      <span className="font-display mt-1 text-2xl font-extrabold" style={{ color: tone }}>
        {value}
      </span>
    </div>
  );
}

/**
 * Celebration shows what the attempt earned (XP, accuracy, streak, companion
 * reaction), then a summary card of the lesson's takeaways, then returns to the
 * journey. All numbers come from the server.
 */
export function Celebration({
  completion,
  streak,
  companionName,
  companionSpecies,
  onBackToJourney,
}: {
  completion: CompletionData;
  streak: number | null;
  companionName?: string;
  companionSpecies?: string;
  onBackToJourney: () => void;
}) {
  const [showSummary, setShowSummary] = useState(false);
  const rewards = completion.rewards;
  const pct = Math.round(completion.accuracy * 100);
  const moved = rewards?.first_completion ? groupMasteryChanges(rewards.mastery_changes) : [];

  if (showSummary) {
    return (
      <div className="animate-result flex flex-1 flex-col sm:justify-center">
        <p className="font-mono mb-3 text-[11px] font-bold uppercase tracking-[0.2em]" style={{ color: "#7C6EF5" }}>
          What you learned
        </p>
        <h1 className="font-display text-2xl font-extrabold" style={{ color: "#E8EAED" }}>
          {completion.lesson.title}
        </h1>
        <ul className="mt-5 flex flex-col gap-3">
          {completion.takeaways.map((t, i) => (
            <li
              key={i}
              className="flex gap-3 rounded-2xl px-4 py-4 text-base leading-snug"
              style={{ background: "#1A1D27", border: "1px solid #2E3347", color: "#E8EAED" }}
            >
              <span
                aria-hidden
                className="font-mono mt-0.5 flex h-6 w-6 shrink-0 items-center justify-center rounded-md text-xs font-bold"
                style={{ background: "rgba(124,110,245,0.2)", color: "#B6ACFF" }}
              >
                {i + 1}
              </span>
              <span className="font-semibold">{t}</span>
            </li>
          ))}
        </ul>

        {completion.unlocked_nodes.length > 0 && (
          <div
            className="mt-5 rounded-2xl px-4 py-4"
            style={{ background: "rgba(52,211,153,0.1)", border: "1px solid rgba(52,211,153,0.4)" }}
          >
            <p className="font-mono text-[11px] font-bold uppercase tracking-[0.2em]" style={{ color: "#34D399" }}>
              Unlocked
            </p>
            <p className="font-display mt-1 text-base font-extrabold" style={{ color: "#E8EAED" }}>
              {completion.unlocked_nodes.map((n) => n.label).join(" · ")}
            </p>
          </div>
        )}

        <div className="mt-auto w-full pt-8 sm:mt-10 sm:pt-0">
          <GameButton type="button" onClick={onBackToJourney} autoFocus>
            Back to journey
          </GameButton>
        </div>
      </div>
    );
  }

  return (
    <div className="animate-result flex flex-1 flex-col items-center text-center sm:justify-center">
      <p className="font-mono mt-2 text-[11px] font-bold uppercase tracking-[0.2em]" style={{ color: "#34D399" }}>
        Lesson complete
      </p>
      <h1 className="font-display mt-2 text-2xl font-extrabold sm:text-3xl" style={{ color: "#E8EAED" }}>
        {completion.lesson.title}
      </h1>

      <div className="mt-6 flex w-full items-start gap-4 text-left">
        <CompanionHero name={companionName} species={companionSpecies} size="md" reaction="correct" />
        <SpeechBubble className="flex-1" speakerName={companionName ?? "Byte"}>
          {message(completion.accuracy, rewards?.first_completion)}
        </SpeechBubble>
      </div>

      <div className="mt-6 grid w-full grid-cols-2 gap-3 sm:grid-cols-4">
        <Chip label="XP" value={rewards ? `+${rewards.xp_awarded}` : "…"} tone="#60A5FA" delay="0.2s" />
        <Chip label="Accuracy" value={`${pct}%`} tone="#34D399" delay="0.3s" />
        <Chip label="Streak" value={streak === null ? "…" : `${streak}🔥`} tone="#FF6B35" delay="0.4s" />
        {rewards && rewards.gems_awarded > 0 && (
          <Chip label="Gems" value={`+${rewards.gems_awarded}`} tone="#34D399" delay="0.5s" />
        )}
      </div>

      {moved.length > 0 && (
        <section className="mt-6 w-full text-left" aria-labelledby="mastery-heading">
          <h2 id="mastery-heading" className="font-mono text-[11px] font-bold uppercase tracking-[0.2em]" style={{ color: "#8B92A8" }}>
            Mastery moved
          </h2>
          <ul className="mt-3 flex flex-col gap-3">
            {moved.map((group) => {
              const theme = topicTheme(group.topicSlug);
              return (
                <li
                  key={group.topicSlug}
                  className="animate-xp rounded-2xl p-4"
                  style={{
                    background: `linear-gradient(150deg, ${theme.tint}26 0%, #1A1D27 60%)`,
                    border: `1px solid ${theme.tint}40`,
                    animationDelay: "0.55s",
                  }}
                >
                  <p className="font-display text-base font-extrabold" style={{ color: "#E8EAED" }}>
                    <span aria-hidden className="mr-2">{theme.icon}</span>
                    {group.topicName}
                  </p>
                  <ul className="mt-2 flex flex-wrap gap-2">
                    {group.skills.map((skill) => (
                      <li
                        key={skill.slug}
                        className="font-mono rounded-full px-3 py-1 text-xs font-bold"
                        style={{ background: "#242836", border: "1px solid #2E3347", color: "#C8CCDA" }}
                      >
                        {skill.name} <span style={{ color: theme.tint }}>{formatDelta(skill.delta)}</span>
                      </li>
                    ))}
                  </ul>
                </li>
              );
            })}
          </ul>
        </section>
      )}

      {completion.rewards_pending && (
        <p className="mt-4 text-xs font-semibold" style={{ color: "#8B92A8" }} role="status">
          Tallying your rewards…
        </p>
      )}
      {rewards && !rewards.first_completion && (
        <p className="mt-4 text-xs font-semibold" style={{ color: "#8B92A8" }}>
          Practice runs don&apos;t earn XP. You earn it the first time you finish a lesson.
        </p>
      )}

      <div className="mt-auto w-full pt-8 sm:mt-10 sm:pt-0">
        <GameButton type="button" onClick={() => setShowSummary(true)} autoFocus>
          Continue
        </GameButton>
      </div>
    </div>
  );
}
