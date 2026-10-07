"use client";

import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { Suspense, useCallback, useEffect, useMemo, useState } from "react";
import { api } from "@/lib/api";
import type { NodeStatus, PathData, PathNode } from "@/lib/lesson/types";
import { GameBackground } from "@/components/game/GameBackground";
import { CompanionHero } from "@/components/game/CompanionHero";
import { GameButton } from "@/components/game/GameButton";
import { BottomNav } from "@/components/game/BottomNav";
import { NodePreviewSheet } from "@/components/lesson/NodePreviewSheet";

const statusText: Record<NodeStatus, string> = {
  done: "completed",
  current: "ready to start",
  available: "available",
  locked: "locked",
};

function NodeButton({
  node,
  animate,
  onOpen,
}: {
  node: PathNode;
  animate: "complete" | "unlock" | null;
  onOpen: () => void;
}) {
  const styles: Record<NodeStatus, string> = {
    done: "bg-[#34D399] text-[#0F1117]",
    current: "animate-node-pulse bg-[#7C6EF5] text-white ring-2 ring-[#7C6EF5]/40",
    available: "bg-[#242836] text-[#B6ACFF] border-2 border-[#7C6EF5]",
    locked: "bg-[#242836] text-[#4A5068] border border-[#2E3347]",
  };
  const icons: Record<NodeStatus, string> = { done: "✓", current: "⚡", available: "◆", locked: "🔒" };
  const isBoss = node.node_type === "boss";
  const anim = animate === "complete" ? "animate-node-complete" : animate === "unlock" ? "animate-node-unlock" : "";

  return (
    <button
      type="button"
      onClick={onOpen}
      aria-label={`${node.label}, ${statusText[node.status]}`}
      className="flex flex-col items-center gap-2 rounded-2xl p-1 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-[#7C6EF5]"
    >
      <span
        className={`flex items-center justify-center rounded-full font-display text-xl font-extrabold shadow-lg transition-transform active:scale-95 ${styles[node.status]} ${anim} ${
          isBoss ? "h-20 w-20" : "h-16 w-16"
        }`}
      >
        {isBoss && node.status !== "locked" ? "👑" : icons[node.status]}
      </span>
      <span className="font-display max-w-[120px] text-center text-xs font-bold" style={{ color: "#C8CCDA" }}>
        {node.label}
      </span>
    </button>
  );
}

function JourneyContent() {
  const router = useRouter();
  const params = useSearchParams();
  const doneNodeId = params.get("done");

  const [path, setPath] = useState<PathData | null>(null);
  const [error, setError] = useState("");
  const [open, setOpen] = useState<PathNode | null>(null);
  const [companion, setCompanion] = useState<{ name?: string; species?: string }>({});

  const load = useCallback(async () => {
    setError("");
    try {
      setPath(await api.getPath());
    } catch (err) {
      setError(err instanceof Error ? err.message : "Couldn't load your journey.");
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
      api
        .getProfile()
        .then((p) => setCompanion({ name: p.companion?.name, species: p.companion?.species }))
        .catch(() => undefined);
    })();
  }, [router, load]);

  // After a lesson, celebrate on the map once, then drop the query so a refresh stays calm.
  useEffect(() => {
    if (!doneNodeId) return;
    const t = setTimeout(() => router.replace("/journey", { scroll: false }), 4000);
    return () => clearTimeout(t);
  }, [doneNodeId, router]);

  const allNodes = useMemo(() => path?.worlds.flatMap((w) => w.nodes) ?? [], [path]);
  const next = allNodes.find((n) => n.status === "current");
  const justDone = doneNodeId ? allNodes.find((n) => n.id === doneNodeId) : undefined;

  const animationFor = (node: PathNode): "complete" | "unlock" | null => {
    if (!doneNodeId) return null;
    if (node.id === doneNodeId) return "complete";
    if (node.status === "current" || node.status === "available") return node === next ? "unlock" : null;
    return null;
  };

  return (
    <GameBackground variant="default">
      <main className="relative mx-auto min-h-dvh max-w-lg px-4 pb-32 pt-8">
        {error ? (
          <div className="mt-16 text-center" role="alert">
            <p className="font-display text-lg font-extrabold" style={{ color: "#E8EAED" }}>
              Couldn&apos;t load your journey
            </p>
            <p className="mt-1 text-sm" style={{ color: "#8B92A8" }}>
              {error}
            </p>
            <div className="mx-auto mt-5 max-w-xs">
              <GameButton type="button" onClick={load}>
                Try again
              </GameButton>
            </div>
          </div>
        ) : !path ? (
          <p className="mt-24 animate-pulse text-center font-mono text-sm font-semibold" style={{ color: "#7C6EF5" }} role="status">
            Loading your journey…
          </p>
        ) : path.worlds.length === 0 ? (
          <div className="mt-24 text-center">
            <p className="font-display text-lg font-extrabold" style={{ color: "#E8EAED" }}>
              Your journey starts here
            </p>
            <p className="mt-1 text-sm" style={{ color: "#8B92A8" }}>
              New lessons are being prepared. Check back soon.
            </p>
          </div>
        ) : (
          <>
            {justDone && (
              <div
                className="animate-sheet mb-6 rounded-2xl px-4 py-3 text-center"
                style={{ background: "rgba(52,211,153,0.1)", border: "1px solid rgba(52,211,153,0.4)" }}
                role="status"
              >
                <p className="font-display text-sm font-extrabold" style={{ color: "#34D399" }}>
                  {justDone.label} complete
                </p>
                {next && (
                  <p className="mt-0.5 text-xs font-semibold" style={{ color: "#C8CCDA" }}>
                    Unlocked: {next.label}
                  </p>
                )}
              </div>
            )}

            {path.worlds.map((w, wi) => (
              <section key={w.id} className={wi > 0 ? "mt-14" : ""} aria-labelledby={`world-${w.slug}`}>
                <div className="text-center">
                  <p className="font-mono text-[11px] font-bold uppercase tracking-[0.2em]" style={{ color: "#7C6EF5" }}>
                    World {wi + 1}
                  </p>
                  <h1 id={`world-${w.slug}`} className="font-display mt-1 text-2xl font-extrabold" style={{ color: "#E8EAED" }}>
                    {w.name}
                  </h1>
                  {w.description && (
                    <p className="mt-1 text-sm font-semibold" style={{ color: "#8B92A8" }}>
                      {w.description}
                    </p>
                  )}
                </div>

                <div className="relative mt-10 flex flex-col items-center gap-8">
                  <div className="absolute left-1/2 top-4 h-[calc(100%-2rem)] w-1 -translate-x-1/2 rounded-full bg-white/15" aria-hidden />
                  {w.nodes.map((node, i) => (
                    <div key={node.id} className={`relative z-10 ${i % 2 === 0 ? "-ml-16" : "ml-16"}`}>
                      {node === next && (
                        <div className="absolute -right-16 -top-1" aria-hidden>
                          <CompanionHero name={companion.name} species={companion.species} size="sm" reaction="idle" />
                        </div>
                      )}
                      <NodeButton node={node} animate={animationFor(node)} onOpen={() => setOpen(node)} />
                    </div>
                  ))}
                </div>
              </section>
            ))}

            {next && (
              <div className="fixed bottom-20 left-0 right-0 z-40 px-4">
                <div className="mx-auto max-w-lg">
                  <Link href={`/lesson/${next.lesson_id}`}>
                    <GameButton type="button">{next.in_progress ? "Continue" : "Start"}: {next.title}</GameButton>
                  </Link>
                </div>
              </div>
            )}
          </>
        )}
      </main>

      {open && <NodePreviewSheet node={open} onClose={() => setOpen(null)} />}
      <BottomNav />
    </GameBackground>
  );
}

export default function JourneyPage() {
  return (
    <Suspense fallback={null}>
      <JourneyContent />
    </Suspense>
  );
}
