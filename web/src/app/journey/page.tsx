"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useCallback, useEffect, useState } from "react";
import { api, JourneyData } from "@/lib/api";
import { GameBackground } from "@/components/game/GameBackground";
import { CompanionHero } from "@/components/game/CompanionHero";
import { GameButton } from "@/components/game/GameButton";
import { BottomNav } from "@/components/game/BottomNav";

function NodeCircle({ status, label, nodeType }: { status: string; label: string; nodeType: string }) {
  const styles: Record<string, string> = {
    done: "bg-[#34D399] text-[#0F1117]",
    current: "animate-node-pulse bg-[#7C6EF5] text-white ring-2 ring-[#7C6EF5]/40",
    locked: "bg-[#242836] text-[#4A5068] border border-[#2E3347]",
  };
  const icons: Record<string, string> = {
    done: "✓",
    current: "⚡",
    locked: "🔒",
  };
  const isBoss = nodeType === "boss";

  return (
    <div className="flex flex-col items-center gap-2">
      <div
        className={`flex items-center justify-center rounded-full font-display text-xl font-extrabold shadow-lg ${styles[status] ?? styles.locked} ${
          isBoss ? "h-20 w-20" : "h-16 w-16"
        }`}
      >
        {isBoss && status !== "locked" ? "👑" : icons[status] ?? "🔒"}
      </div>
      <span className="font-display max-w-[110px] text-center text-xs font-bold" style={{ color: "#C8CCDA" }}>{label}</span>
    </div>
  );
}

export default function JourneyPage() {
  const router = useRouter();
  const [journey, setJourney] = useState<JourneyData | null>(null);
  const [loading, setLoading] = useState(true);

  const load = useCallback(async () => {
    try {
      setJourney(await api.getJourney());
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    (async () => {
      const ok = await api.ensureSession();
      if (!ok) {
        router.replace("/login");
        return;
      }
      load();
    })();
  }, [router, load]);

  const bgVariant = journey?.world.theme === "forest" ? "forest" : "default";

  return (
    <GameBackground variant={bgVariant as "forest" | "default"}>
      <main className="relative mx-auto min-h-screen max-w-lg px-4 pb-28 pt-6">
        <div className="text-center">
          <h1 className="font-display text-2xl font-extrabold" style={{ color: "#E8EAED" }}>
            🌲 {journey?.world.name ?? "Journey"}
          </h1>
          <p className="mt-1 text-sm font-semibold" style={{ color: "#8B92A8" }}>
            {journey?.world.description ?? "Loading..."}
          </p>
        </div>

        {loading ? (
          <p className="mt-12 text-center font-semibold text-white">Loading journey...</p>
        ) : (
          <div className="relative mt-12 flex flex-col items-center gap-8">
            <div className="absolute left-1/2 top-0 h-full w-1 -translate-x-1/2 rounded-full bg-white/30" />
            {journey?.nodes.map((node, i) => (
              <div key={node.id} className={`relative z-10 ${i % 2 === 0 ? "-ml-16" : "ml-16"}`}>
                {node.status === "current" && (
                  <div className="absolute -right-14 -top-2">
                    <CompanionHero size="sm" reaction="idle" />
                  </div>
                )}
                <NodeCircle status={node.status} label={node.label} nodeType={node.node_type} />
              </div>
            ))}
          </div>
        )}

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
