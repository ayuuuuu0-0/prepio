"use client";

import Link from "next/link";
import { GameBackground } from "./GameBackground";
import { CompanionHero } from "./CompanionHero";
import { SpeechBubble } from "./SpeechBubble";
import { BottomNav } from "./BottomNav";

/** ComingSoonPage is a placeholder for features launching in a future phase. */
export function ComingSoonPage({
  title,
  icon,
  message,
}: {
  title: string;
  icon: string;
  message: string;
}) {
  return (
    <GameBackground>
      <main className="mx-auto flex min-h-screen max-w-lg flex-col items-center justify-center px-6 pb-28 pt-12 text-center">
        <p className="text-5xl">{icon}</p>
        <h1 className="font-display mt-4 text-2xl font-extrabold" style={{ color: "#E8EAED" }}>
          {title}
        </h1>
        <CompanionHero size="md" />
        <SpeechBubble className="mt-6 max-w-sm">{message}</SpeechBubble>
        <Link href="/dashboard" className="font-mono mt-8 text-sm font-semibold" style={{ color: "#7C6EF5" }}>
          ← Back to Home
        </Link>
      </main>
      <BottomNav />
    </GameBackground>
  );
}
