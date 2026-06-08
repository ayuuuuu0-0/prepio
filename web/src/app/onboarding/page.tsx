"use client";

import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { api, Companion, EXPERIENCE_LEVELS, TARGET_COMPANIES } from "@/lib/api";
import { companionVisual } from "@/lib/design/companions";
import { companyColors } from "@/lib/design/tokens";
import { GameBackground } from "@/components/game/GameBackground";
import { CompanionHero } from "@/components/game/CompanionHero";
import { SpeechBubble } from "@/components/game/SpeechBubble";
import { GameButton } from "@/components/game/GameButton";

const stepMessages = [
  "Which dream companies are we conquering?",
  "How much experience do you bring to the quest?",
  "Pick your companion — they'll guide your journey!",
];

export default function OnboardingPage() {
  const router = useRouter();
  const [step, setStep] = useState(1);
  const [companions, setCompanions] = useState<Companion[]>([]);
  const [targets, setTargets] = useState<string[]>([]);
  const [experience, setExperience] = useState("");
  const [companionId, setCompanionId] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const selectedCompanion = companions.find((c) => c.id === companionId);

  useEffect(() => {
    if (!api.loadToken()) {
      router.replace("/login");
      return;
    }
    api.getCompanions().then(setCompanions).catch((err) =>
      setError(err instanceof Error ? err.message : "Failed to load companions")
    );
  }, [router]);

  function toggleTarget(company: string) {
    setTargets((prev) =>
      prev.includes(company) ? prev.filter((c) => c !== company) : [...prev, company]
    );
  }

  async function finish() {
    setLoading(true);
    setError("");
    try {
      await api.completeOnboarding(targets, experience, companionId);
      router.replace("/dashboard");
    } catch (err) {
      setError(err instanceof Error ? err.message : "onboarding failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <GameBackground variant="forest">
      <main className="mx-auto min-h-screen max-w-lg px-4 py-8">
        <div className="flex items-center gap-2">
          {[1, 2, 3].map((s) => (
            <div
              key={s}
              className={`h-3 flex-1 rounded-full transition ${
                s <= step ? "bg-[#58CC02]" : "bg-white/50"
              }`}
            />
          ))}
        </div>

        <div className="mt-8 flex items-end gap-3">
          <CompanionHero
            name={selectedCompanion?.name ?? "Byte"}
            species={selectedCompanion?.species}
            size="md"
          />
          <SpeechBubble className="flex-1">{stepMessages[step - 1]}</SpeechBubble>
        </div>

        {error && (
          <p className="mt-4 rounded-2xl bg-orange-100 px-4 py-3 text-center text-sm font-semibold text-orange-700">
            {error}
          </p>
        )}

        {step === 1 && (
          <section className="mt-6 grid grid-cols-2 gap-3">
            {TARGET_COMPANIES.map((company) => {
              const selected = targets.includes(company);
              const cc = companyColors[company];
              return (
                <button
                  key={company}
                  type="button"
                  onClick={() => toggleTarget(company)}
                  className={`font-display rounded-3xl px-4 py-5 text-lg font-bold capitalize shadow-md transition ${
                    selected ? "scale-105 ring-4 ring-[#58CC02]" : "hover:scale-102"
                  }`}
                  style={{
                    backgroundColor: cc?.bg ?? "#f0f0f0",
                    color: cc?.ring ?? "#333",
                  }}
                >
                  {company}
                </button>
              );
            })}
            <div className="col-span-2 mt-4">
              <GameButton disabled={targets.length === 0} onClick={() => setStep(2)}>
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
                className={`font-display block w-full rounded-3xl px-5 py-4 text-left text-lg font-bold shadow-md transition ${
                  experience === level.id
                    ? "bg-[#7B5CFF] text-white ring-4 ring-purple-300"
                    : "bg-white text-[#3C3C3C] hover:bg-purple-50"
                }`}
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
                  className={`flex w-full items-center gap-4 rounded-3xl p-4 shadow-md transition ${
                    selected ? "ring-4 ring-[#58CC02] scale-[1.02]" : "bg-white"
                  } bg-gradient-to-r ${v.gradient}`}
                >
                  <span className="text-4xl">{v.emoji}</span>
                  <div className="text-left">
                    <p className="font-display text-xl font-bold text-white drop-shadow">{c.name}</p>
                    <p className="text-sm font-semibold capitalize text-white/90">
                      {c.species.replace("_", " ")}
                    </p>
                  </div>
                </button>
              );
            })}
            <div className="mt-4 flex gap-3">
              <GameButton variant="secondary" className="flex-1" onClick={() => setStep(2)}>
                Back
              </GameButton>
              <GameButton
                className="flex-1"
                variant="gold"
                disabled={companionId.length === 0 || loading}
                onClick={finish}
              >
                {loading ? "Starting..." : "Start My Journey!"}
              </GameButton>
            </div>
          </section>
        )}
      </main>
    </GameBackground>
  );
}
