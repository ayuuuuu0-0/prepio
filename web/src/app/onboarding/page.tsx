"use client";

import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { api, Companion, EXPERIENCE_LEVELS, TopicInfo } from "@/lib/api";
import { companionVisual } from "@/lib/design/companions";
import { topicTheme } from "@/lib/topics";
import { GameBackground } from "@/components/game/GameBackground";
import { CompanionHero } from "@/components/game/CompanionHero";
import { SpeechBubble } from "@/components/game/SpeechBubble";
import { GameButton } from "@/components/game/GameButton";

const MAX_TOPICS = 3;

const stepMessages = [
  "What do you want to sharpen? Pick up to three. They move to the front of your path, and nothing gets locked.",
  "How much experience do you have? This sets your starting point.",
  "Choose your companion. They'll grow with you throughout the journey.",
];

export default function OnboardingPage() {
  const router = useRouter();
  const [step, setStep] = useState(1);
  const [topics, setTopics] = useState<TopicInfo[]>([]);
  const [companions, setCompanions] = useState<Companion[]>([]);
  const [focus, setFocus] = useState<string[]>([]);
  const [experience, setExperience] = useState("");
  const [companionId, setCompanionId] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const selectedCompanion = companions.find((c) => c.id === companionId);

  useEffect(() => {
    let cancelled = false;
    (async () => {
      const ok = await api.ensureSession();
      if (cancelled) return;
      if (!ok) {
        router.replace("/login");
        return;
      }
      try {
        const [c, t] = await Promise.all([api.getCompanions(), api.getTopics()]);
        if (cancelled) return;
        setCompanions(c);
        setTopics(t);
      } catch (err) {
        if (!cancelled) setError(err instanceof Error ? err.message : "Failed to load onboarding");
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [router]);

  function toggleTopic(slug: string) {
    setFocus((current) => {
      if (current.includes(slug)) return current.filter((s) => s !== slug);
      if (current.length >= MAX_TOPICS) return current;
      return [...current, slug];
    });
  }

  async function finish() {
    setLoading(true);
    setError("");
    try {
      await api.completeOnboarding(experience, companionId, focus);
      router.replace("/dashboard");
    } catch (err) {
      setError(err instanceof Error ? err.message : "onboarding failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <GameBackground variant="forest">
      <main className="mx-auto min-h-dvh max-w-lg px-4 py-8">
        <div className="flex items-center gap-2" role="progressbar" aria-label="Onboarding progress" aria-valuemin={1} aria-valuemax={3} aria-valuenow={step}>
          {[1, 2, 3].map((s) => (
            <div key={s} className={`h-3 flex-1 rounded-full transition ${s <= step ? "bg-[#7C6EF5]" : "bg-[#2E3347]"}`} />
          ))}
        </div>

        <div className="mt-8 flex items-end gap-3">
          <CompanionHero name={selectedCompanion?.name ?? "Byte"} species={selectedCompanion?.species} size="md" />
          <SpeechBubble className="flex-1">{stepMessages[step - 1]}</SpeechBubble>
        </div>

        {error && (
          <p
            className="mt-4 rounded-2xl px-4 py-3 text-center text-sm font-semibold"
            role="alert"
            style={{ background: "rgba(248,113,113,0.1)", border: "1px solid rgba(248,113,113,0.3)", color: "#F87171" }}
          >
            {error}
          </p>
        )}

        {step === 1 && (
          <section className="mt-6" aria-labelledby="topics-title">
            <h1 id="topics-title" className="font-display text-2xl font-extrabold tracking-tight" style={{ color: "#E8EAED" }}>
              What do you want to sharpen?
            </h1>
            <p className="font-mono mt-1 text-xs font-bold uppercase tracking-[0.18em]" style={{ color: "#8B92A8" }} aria-live="polite">
              {focus.length} of {MAX_TOPICS} picked
            </p>
            <div className="mt-4 space-y-3" role="group" aria-label="Topics">
              {topics.map((topic) => {
                const theme = topicTheme(topic.slug);
                const index = focus.indexOf(topic.slug);
                const picked = index >= 0;
                const full = !picked && focus.length >= MAX_TOPICS;
                return (
                  <button
                    key={topic.slug}
                    type="button"
                    aria-pressed={picked}
                    disabled={full}
                    onClick={() => toggleTopic(topic.slug)}
                    className="relative flex w-full items-center gap-4 overflow-hidden rounded-2xl p-4 text-left transition-[transform,border-color] active:scale-[0.99] disabled:cursor-not-allowed disabled:opacity-45 motion-reduce:transition-none focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7C6EF5]"
                    style={{
                      background: `linear-gradient(145deg, ${theme.tint}${picked ? "40" : "22"} 0%, #1A1D27 62%)`,
                      border: `2px solid ${picked ? theme.tint : "#2E3347"}`,
                    }}
                  >
                    <span aria-hidden className="text-3xl">
                      {theme.icon}
                    </span>
                    <span className="min-w-0 flex-1">
                      <span className="font-display block text-lg font-extrabold leading-tight" style={{ color: "#E8EAED" }}>
                        {topic.name}
                      </span>
                      <span className="mt-0.5 block text-sm font-semibold leading-snug" style={{ color: "#8B92A8" }}>
                        {topic.description}
                      </span>
                    </span>
                    <span
                      aria-hidden
                      className="font-mono flex h-7 w-7 shrink-0 items-center justify-center rounded-full text-xs font-bold"
                      style={{
                        background: picked ? theme.tint : "#242836",
                        color: picked ? "#0F1117" : "#4A5068",
                        border: picked ? "none" : "1px solid #2E3347",
                      }}
                    >
                      {picked ? index + 1 : ""}
                    </span>
                    {picked && <span className="sr-only">Picked, priority {index + 1}</span>}
                  </button>
                );
              })}
            </div>
            <div className="mt-6">
              <GameButton disabled={focus.length === 0} onClick={() => setStep(2)}>
                Continue →
              </GameButton>
            </div>
          </section>
        )}

        {step === 2 && (
          <section className="mt-6 space-y-3">
            {EXPERIENCE_LEVELS.map((level) => (
              <button
                key={level.id}
                type="button"
                onClick={() => setExperience(level.id)}
                className={`font-display block w-full rounded-2xl px-5 py-4 text-left text-lg font-bold transition ${
                  experience === level.id ? "ring-2 ring-[#7C6EF5]" : "hover:border-[#7C6EF5]/50"
                }`}
                style={{
                  background: experience === level.id ? "rgba(124,110,245,0.2)" : "#1A1D27",
                  border: "1px solid #2E3347",
                  color: "#E8EAED",
                }}
              >
                {level.label}
              </button>
            ))}
            <div className="mt-4 flex gap-3">
              <GameButton variant="secondary" className="flex-1" onClick={() => setStep(1)}>
                Back
              </GameButton>
              <GameButton className="flex-1" disabled={experience.length === 0} onClick={() => setStep(3)}>
                Continue →
              </GameButton>
            </div>
          </section>
        )}

        {step === 3 && (
          <section className="mt-6 space-y-3">
            {companions.map((c) => {
              const v = companionVisual(c.name, c.species);
              const selected = companionId === c.id;
              return (
                <button
                  key={c.id}
                  type="button"
                  onClick={() => setCompanionId(c.id)}
                  className={`flex w-full items-center gap-4 rounded-2xl p-4 transition ${selected ? "scale-[1.02] ring-2 ring-[#7C6EF5]" : ""}`}
                  style={{
                    background: `linear-gradient(135deg, ${v.glow}33 0%, #1A1D27 100%)`,
                    border: `1px solid ${selected ? "#7C6EF5" : "#2E3347"}`,
                  }}
                >
                  <span className="text-4xl">{v.emoji}</span>
                  <div className="text-left">
                    <p className="font-display text-xl font-bold text-white drop-shadow">{c.name}</p>
                    <p className="text-sm font-semibold capitalize text-white/90">{c.species.replace("_", " ")}</p>
                  </div>
                </button>
              );
            })}
            <div className="mt-4 flex gap-3">
              <GameButton variant="secondary" className="flex-1" onClick={() => setStep(2)}>
                Back
              </GameButton>
              <GameButton className="flex-1" variant="gold" disabled={companionId.length === 0 || loading} onClick={finish}>
                {loading ? "Setting up..." : "Start Learning"}
              </GameButton>
            </div>
          </section>
        )}
      </main>
    </GameBackground>
  );
}
