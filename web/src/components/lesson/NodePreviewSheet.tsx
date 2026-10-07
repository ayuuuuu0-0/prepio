"use client";

import { useEffect, useRef } from "react";
import Link from "next/link";
import type { PathNode } from "@/lib/lesson/types";

const difficultyLabel = { easy: "Easy", medium: "Medium", hard: "Hard" } as const;

/**
 * NodePreviewSheet opens when a journey node is tapped: what you'll learn, how long it takes,
 * and the XP on offer. Locked nodes say exactly what unlocks them, so no tap is ever dead.
 */
export function NodePreviewSheet({ node, onClose }: { node: PathNode; onClose: () => void }) {
  const closeRef = useRef<HTMLButtonElement>(null);
  const locked = node.status === "locked";
  const done = node.status === "done";

  useEffect(() => {
    closeRef.current?.focus();
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose();
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [onClose]);

  const cta = done ? "Practice again" : node.in_progress ? "Continue lesson" : "Start lesson";

  return (
    <div
      className="fixed inset-0 z-[60] flex items-end justify-center bg-black/60 sm:items-center sm:p-4"
      onClick={onClose}
      role="presentation"
    >
      <div
        role="dialog"
        aria-modal="true"
        aria-labelledby="node-title"
        onClick={(e) => e.stopPropagation()}
        className="animate-sheet w-full max-w-md rounded-t-3xl p-6 pb-8 sm:rounded-3xl"
        style={{ background: "#1A1D27", border: "1px solid #2E3347", boxShadow: "0 -16px 60px rgba(0,0,0,0.6)" }}
      >
        <div className="flex items-start justify-between gap-4">
          <div>
            <p className="font-mono text-[11px] font-bold uppercase tracking-[0.2em]" style={{ color: "#7C6EF5" }}>
              {node.kind === "boss" ? "Boss lesson" : "Lesson"}
              {done ? " · Completed" : ""}
            </p>
            <h2 id="node-title" className="font-display mt-1 text-xl font-extrabold leading-snug" style={{ color: "#E8EAED" }}>
              {node.title}
            </h2>
          </div>
          <button
            ref={closeRef}
            type="button"
            onClick={onClose}
            aria-label="Close"
            className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full transition-colors hover:bg-white/10 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7C6EF5]"
            style={{ color: "#8B92A8" }}
          >
            <span aria-hidden>✕</span>
          </button>
        </div>

        <div className="mt-4 flex flex-wrap gap-2">
          <Stat label={`${node.est_minutes} min`} />
          <Stat label={difficultyLabel[node.difficulty]} />
          <Stat label={`Up to ${node.xp_preview} XP`} tone="#60A5FA" />
        </div>

        <p className="font-mono mt-5 text-[11px] font-bold uppercase tracking-[0.2em]" style={{ color: "#8B92A8" }}>
          What you&apos;ll learn
        </p>
        <ul className="mt-2 flex flex-col gap-2">
          {node.takeaways.map((t, i) => (
            <li key={i} className="flex gap-2 text-sm leading-snug" style={{ color: "#C8CCDA" }}>
              <span aria-hidden style={{ color: "#7C6EF5" }}>
                ◆
              </span>
              <span>{t}</span>
            </li>
          ))}
        </ul>

        {locked ? (
          <div
            className="mt-6 flex items-start gap-3 rounded-2xl px-4 py-3"
            style={{ background: "#242836", border: "1px solid #2E3347" }}
            role="status"
          >
            <span aria-hidden className="text-lg">
              🔒
            </span>
            <p className="text-sm font-semibold leading-snug" style={{ color: "#C8CCDA" }}>
              {node.unlock_hint || "Finish the earlier lessons to unlock this one."}
            </p>
          </div>
        ) : (
          <Link
            href={`/lesson/${node.lesson_id}`}
            className="game-btn font-display mt-6 block w-full rounded-full px-8 py-4 text-center text-base font-bold tracking-wide text-white"
            style={{
              background: done ? "#242836" : "linear-gradient(90deg, #7C6EF5, #9D8FF7)",
              border: done ? "1px solid #2E3347" : "none",
              ["--btn-shadow" as string]: done ? "#1A1D27" : "#5B50D4",
            }}
          >
            {cta}
          </Link>
        )}
      </div>
    </div>
  );
}

function Stat({ label, tone = "#C8CCDA" }: { label: string; tone?: string }) {
  return (
    <span
      className="font-mono rounded-full px-3 py-1 text-xs font-bold"
      style={{ background: "#242836", border: "1px solid #2E3347", color: tone }}
    >
      {label}
    </span>
  );
}
