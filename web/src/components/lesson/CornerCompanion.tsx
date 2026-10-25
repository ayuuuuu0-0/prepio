"use client";

import { useEffect, useState } from "react";
import { CompanionHero } from "@/components/game/CompanionHero";
import { SpeechBubble } from "@/components/game/SpeechBubble";

type Reaction = "idle" | "correct" | "wrong";

export function CornerCompanion({
  name,
  species,
  reaction = "idle",
  dialogueLine,
  soundOn,
  onToggleSound,
  onPoke,
}: {
  name?: string;
  species?: string;
  reaction?: Reaction;
  dialogueLine?: string | null;
  soundOn: boolean;
  onToggleSound: () => void;
  onPoke?: () => void;
}) {
  const [bubbleText, setBubbleText] = useState<string | null>(null);
  const [poked, setPoked] = useState(false);

  // Update bubble text when a new dialogue line arrives
  useEffect(() => {
    if (dialogueLine) {
      setBubbleText(dialogueLine);
    }
  }, [dialogueLine]);

  // Handle companion click / poke
  const handlePoke = () => {
    setPoked(true);
    setTimeout(() => setPoked(false), 600);

    const idleLines = [
      `I'm watching your code! Keep going.`,
      `Stay focused — you got this!`,
      `Quality over speed. Think through edge cases.`,
      `Your streak is on the line, let's nail it!`,
    ];
    const randomLine = idleLines[Math.floor(Math.random() * idleLines.length)];
    setBubbleText(randomLine);
    if (onPoke) onPoke();
  };

  const displayName = name ?? "Byte";

  return (
    <aside
      aria-label="Companion co-pilot"
      className="fixed bottom-24 right-4 z-40 flex items-end sm:bottom-28 sm:right-6 pointer-events-none"
    >
      {/* Dynamic Floating Speech Bubble */}
      {bubbleText && (
        <div className="pointer-events-auto mr-3 mb-2 max-w-[220px] sm:max-w-xs transition-all duration-300 animate-in fade-in slide-in-from-right-4">
          <div className="relative">
            <SpeechBubble speakerName={displayName} className="text-xs sm:text-sm shadow-xl">
              {bubbleText}
            </SpeechBubble>
            {/* Small speech arrow pointing right toward companion */}
            <div
              className="absolute -right-2 bottom-4 h-0 w-0 border-y-4 border-y-transparent border-l-8 border-l-[#1A1D27]"
              aria-hidden
            />
          </div>
        </div>
      )}

      {/* Companion Avatar Button with Micro-interactions */}
      <div className="pointer-events-auto relative flex flex-col items-center">
        <button
          type="button"
          onClick={handlePoke}
          title={`Click to interact with ${displayName}`}
          className={`group relative rounded-full transition-transform active:scale-95 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-[#7C6EF5] ${
            poked ? "animate-bounce" : ""
          }`}
        >
          <CompanionHero
            name={name}
            species={species}
            size="md"
            reaction={reaction}
          />

          {/* Emotional Status Badge */}
          <div
            className="absolute -bottom-1 -right-1 flex h-6 w-6 items-center justify-center rounded-full text-xs font-bold border border-[#1A1D27] shadow-md"
            style={{
              background:
                reaction === "correct"
                  ? "#34D399"
                  : reaction === "wrong"
                  ? "#F87171"
                  : "#7C6EF5",
              color: "#0E1015",
            }}
          >
            {reaction === "correct" ? "✓" : reaction === "wrong" ? "!" : "★"}
          </div>
        </button>

        {/* Audio Mute/Unmute Pill */}
        <button
          type="button"
          onClick={(e) => {
            e.stopPropagation();
            onToggleSound();
          }}
          title={soundOn ? "Sound enabled (click to mute)" : "Sound muted (click to unmute)"}
          className="pointer-events-auto mt-1 flex items-center gap-1 rounded-full px-2 py-0.5 text-[10px] font-mono font-semibold transition-colors hover:bg-white/10"
          style={{
            background: "#1A1D27",
            border: "1px solid #2E3347",
            color: soundOn ? "#E8EAED" : "#6B7280",
          }}
        >
          <span>{soundOn ? "🔊" : "🔇"}</span>
          <span className="hidden sm:inline">{soundOn ? "On" : "Off"}</span>
        </button>
      </div>
    </aside>
  );
}
