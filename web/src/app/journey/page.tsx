"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { api } from "@/lib/api";
import { GameBackground } from "@/components/game/GameBackground";
import { CompanionHero } from "@/components/game/CompanionHero";
import { GameButton } from "@/components/game/GameButton";
import { BottomNav } from "@/components/game/BottomNav";

/** Journey nodes for Foundation Forest — hero screen placeholder per SCREEN_SPECS. */
const nodes = [
  { id: 1, type: "done", label: "Arrays" },
  { id: 2, type: "done", label: "Strings" },
  { id: 3, type: "current", label: "Hashing" },
  { id: 4, type: "locked", label: "Sorting" },
  { id: 5, type: "boss", label: "Boss" },
];

function NodeCircle({ type, label }: { type: string; label: string }) {
  const styles: Record<string, string> = {
    done: "bg-[#58CC02] text-white",
    current: "animate-node-pulse bg-[#1CB0F6] text-white ring-4 ring-sky-200",
    locked: "bg-gray-300 text-gray-500",
    boss: "h-20 w-20 bg-gradient-to-br from-[#FFC800] to-[#FF9600] text-white text-2xl",
  };
  const icons: Record<string, string> = {
    done: "✓",
    current: "⚡",
    locked: "🔒",
    boss: "👑",
  };
  const isBoss = type === "boss";

  return (
    <div className="flex flex-col items-center gap-2">
      <div
        className={`flex items-center justify-center rounded-full font-display text-xl font-extrabold shadow-lg ${styles[type]} ${
          isBoss ? "h-20 w-20" : "h-16 w-16"
        }`}
      >
        {icons[type]}
      </div>
      <span className="font-display text-xs font-bold text-white drop-shadow">{label}</span>
    </div>
  );
}

export default function JourneyPage() {
  const router = useRouter();

  useEffect(() => {
    if (!api.loadToken()) router.replace("/login");
  }, [router]);

  return (
    <GameBackground variant="forest">
      <main className="relative mx-auto min-h-screen max-w-lg px-4 pb-28 pt-6">
        <div className="text-center">
          <h1 className="font-display text-2xl font-extrabold text-white drop-shadow-lg">🌲 Foundation Forest</h1>
          <p className="mt-1 text-sm font-semibold text-white/80">World 1 · Beginner Journey</p>
        </div>

        {/* Winding path of nodes */}
        <div className="relative mt-12 flex flex-col items-center gap-8">
          <div className="absolute left-1/2 top-0 h-full w-1 -translate-x-1/2 rounded-full bg-white/30" />
          {nodes.map((node, i) => (
            <div key={node.id} className={`relative z-10 ${i % 2 === 0 ? "-ml-16" : "ml-16"}`}>
              {node.type === "current" && (
                <div className="absolute -right-14 -top-2">
                  <CompanionHero size="sm" />
                </div>
              )}
              <NodeCircle type={node.type} label={node.label} />
            </div>
          ))}
        </div>

        <div className="mt-10 px-4">
          <Link href="/challenge">
            <GameButton type="button">Play Current Node →</GameButton>
          </Link>
        </div>
      </main>
      <BottomNav />
    </GameBackground>
  );
}
