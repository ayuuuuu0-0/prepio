import type { ReactNode } from "react";
import { MasteryRing } from "@/components/dashboard/MasteryRing";
import { CompanionHero } from "@/components/game/CompanionHero";
import { topicTheme } from "@/lib/topics";
import s from "./landing.module.css";

/** Phone is a device frame for a looping product illustration. */
export function Phone({ label, status, children }: { label: string; status: string; children: ReactNode }) {
  return (
    <figure className="flex shrink-0 flex-col items-center">
      <figcaption className="mb-3 text-center">
        <span className="font-display block text-sm font-bold" style={{ color: "#E8EAED" }}>
          {label}
        </span>
        <span className="font-mono mt-0.5 flex items-center justify-center gap-1.5 text-[10px]" style={{ color: "#8B92A8" }}>
          <span aria-hidden className="h-1.5 w-1.5 rounded-full" style={{ background: "#34D399" }} />
          {status}
        </span>
      </figcaption>
      <div
        aria-hidden
        className="relative h-[430px] w-[212px] overflow-hidden rounded-[2.2rem] p-[7px]"
        style={{ background: "linear-gradient(160deg, #3a3f55, #1a1d27)", boxShadow: "0 30px 60px rgba(0,0,0,0.5)" }}
      >
        <div className="relative h-full w-full overflow-hidden rounded-[1.75rem]" style={{ background: "#0F1117" }}>
          <div className="absolute left-1/2 top-2 z-10 h-[18px] w-[66px] -translate-x-1/2 rounded-full bg-black" />
          <div className="font-mono flex justify-between px-5 pt-2.5 text-[9px] font-semibold" style={{ color: "#E8EAED" }}>
            <span>9:41</span>
            <span>●●● ▮</span>
          </div>
          <div className="h-[calc(100%-24px)]">{children}</div>
        </div>
      </div>
    </figure>
  );
}

const journeyNodes = [
  { label: "Caching basics", state: "done", x: 30 },
  { label: "Invalidation", state: "done", x: 62 },
  { label: "Load balancing", state: "current", x: 40 },
  { label: "Queues", state: "locked", x: 18 },
  { label: "Boss: a feed", state: "locked", x: 50 },
] as const;

/** JourneyScreen illustrates the path: finished nodes, the current one pulsing, the rest ahead. */
export function JourneyScreen() {
  return (
    <div className="flex h-full flex-col px-3 pt-4">
      <p className="font-mono text-[8px] font-bold uppercase tracking-[0.2em]" style={{ color: "#7C6EF5" }}>
        World 1
      </p>
      <p className="font-display text-[15px] font-extrabold" style={{ color: "#E8EAED" }}>
        Caching &amp; scale
      </p>
      <div className="relative mt-3 flex-1">
        {journeyNodes.map((n, i) => (
          <div key={n.label} className="absolute flex flex-col items-center" style={{ left: `${n.x}%`, top: `${i * 19 + 2}%` }}>
            <span
              className={`flex h-10 w-10 items-center justify-center rounded-full text-sm font-extrabold ${n.state === "current" ? s.nodeCurrent : ""}`}
              style={{
                background: n.state === "done" ? "#34D399" : n.state === "current" ? "#7C6EF5" : "#242836",
                color: n.state === "locked" ? "#4A5068" : "#0F1117",
                border: n.state === "locked" ? "1px solid #2E3347" : "none",
                boxShadow: n.state === "done" ? "0 4px 0 #1F9C70" : n.state === "current" ? "0 4px 0 #5B50D4" : "none",
              }}
            >
              {n.state === "done" ? "✓" : n.state === "current" ? "★" : "🔒"}
            </span>
            <span className="font-mono mt-1 whitespace-nowrap text-[8px]" style={{ color: n.state === "locked" ? "#4A5068" : "#8B92A8" }}>
              {n.label}
            </span>
          </div>
        ))}
      </div>
      <div className="mb-3 rounded-xl px-3 py-2 text-[10px]" style={{ background: "#1A1D27", border: "1px solid #2E3347", color: "#8B92A8" }}>
        <span className="font-semibold" style={{ color: "#E8EAED" }}>
          Next:
        </span>{" "}
        Load balancing · 4 min · +20 XP
      </div>
    </div>
  );
}

const options = ["Add a read-through cache", "Shard the users table", "Rewrite it in Rust"];

/** LessonScreen illustrates one exercise: pick, Check, and the feedback tray. */
export function LessonScreen() {
  return (
    <div className="relative flex h-full flex-col px-3 pt-3">
      <div className="flex items-center gap-2">
        <span className="text-[10px]" style={{ color: "#8B92A8" }}>
          ✕
        </span>
        <div className="h-2 flex-1 overflow-hidden rounded-full" style={{ background: "#242836" }}>
          <div className="h-full w-3/5 rounded-full" style={{ background: "linear-gradient(90deg, #7C6EF5, #9D8FF7)" }} />
        </div>
        <span className="font-mono text-[9px] font-bold" style={{ color: "#F5B942" }}>
          🔥 3
        </span>
      </div>
      <p className="font-mono mt-4 text-[8px] font-bold uppercase tracking-[0.18em]" style={{ color: "#7C6EF5" }}>
        Choose one
      </p>
      <p className="font-display mt-1 text-[13px] font-bold leading-snug" style={{ color: "#E8EAED" }}>
        A read-heavy endpoint is slow and the data rarely changes. What do you reach for first?
      </p>
      <div className="mt-3 flex flex-col gap-2">
        {options.map((o, i) => (
          <div
            key={o}
            className={`flex items-center gap-2 rounded-xl px-2.5 py-2 text-[10px] font-semibold ${i === 0 ? s.optionPicked : ""}`}
            style={{ background: "#1A1D27", border: "1.5px solid #2E3347", color: "#E8EAED" }}
          >
            <span className="font-mono flex h-4 w-4 items-center justify-center rounded text-[8px]" style={{ background: "#242836", color: "#8B92A8" }}>
              {i + 1}
            </span>
            {o}
          </div>
        ))}
      </div>
      <div
        className={`font-display mt-auto mb-4 rounded-full py-2 text-center text-[11px] font-extrabold ${s.checkButton}`}
        style={{ background: "#7C6EF5", color: "#fff", boxShadow: "0 3px 0 #5B50D4" }}
      >
        Check
      </div>
      <div className={`absolute inset-x-0 bottom-0 px-3 pb-4 pt-3 ${s.tray}`} style={{ background: "#10231C", borderTop: "2px solid #34D399" }}>
        <p className="font-display flex items-center gap-1.5 text-[12px] font-extrabold" style={{ color: "#34D399" }}>
          <span className="flex h-4 w-4 items-center justify-center rounded-full text-[9px]" style={{ background: "#34D399", color: "#0F1117" }}>
            ✓
          </span>
          Nice, that&apos;s right
        </p>
        <p className="mt-1 text-[9px] leading-snug" style={{ color: "#B9C0D0" }}>
          Rarely-changing reads are exactly what a cache is for: the database stops doing the same work twice.
        </p>
        <div className="font-display mt-2 rounded-full py-1.5 text-center text-[10px] font-extrabold" style={{ background: "#34D399", color: "#0F1117" }}>
          Continue
        </div>
      </div>
    </div>
  );
}

const masteryMoves = [
  { skill: "Caching", delta: 3, width: "72%" },
  { skill: "Load Balancing", delta: 1, width: "38%" },
];

/** CelebrationScreen illustrates the finish: XP, the companion's reaction, and mastery that moved. */
export function CelebrationScreen() {
  const sd = topicTheme("system-design");
  return (
    <div className="flex h-full flex-col items-center px-3 pt-5 text-center">
      <CompanionHero name="Pip" size="sm" reaction="correct" />
      <p className="font-display mt-2 text-[15px] font-extrabold" style={{ color: "#E8EAED" }}>
        Lesson complete!
      </p>
      <div className="mt-3 flex gap-2">
        <span className={`font-mono rounded-lg px-2 py-1 text-[10px] font-bold ${s.pop}`} style={{ background: "rgba(96,165,250,0.14)", color: "#60A5FA" }}>
          +20 XP
        </span>
        <span className={`font-mono rounded-lg px-2 py-1 text-[10px] font-bold ${s.pop}`} style={{ background: "rgba(52,211,153,0.14)", color: "#34D399" }}>
          💎 +10
        </span>
        <span className={`font-mono rounded-lg px-2 py-1 text-[10px] font-bold ${s.pop}`} style={{ background: "rgba(255,107,53,0.14)", color: "#FF6B35" }}>
          🔥 4
        </span>
      </div>
      <div className="mt-4 flex w-full items-center gap-3 rounded-2xl p-3 text-left" style={{ background: "#1A1D27", border: "1px solid #2E3347" }}>
        <MasteryRing mastery={72} tint={sd.tint} size={52} />
        <div className="min-w-0 flex-1">
          <p className="font-display text-[11px] font-bold" style={{ color: "#E8EAED" }}>
            System Design
          </p>
          {masteryMoves.map((m) => (
            <div key={m.skill} className="mt-1.5">
              <div className="flex justify-between text-[9px]" style={{ color: "#8B92A8" }}>
                <span>{m.skill}</span>
                <span style={{ color: "#34D399" }}>+{m.delta}</span>
              </div>
              <div className="mt-0.5 h-1.5 overflow-hidden rounded-full" style={{ background: "#242836" }}>
                <div className={`h-full rounded-full ${s.grow}`} style={{ width: m.width, background: sd.tint }} />
              </div>
            </div>
          ))}
        </div>
      </div>
      <p className="mt-3 text-[9px]" style={{ color: "#8B92A8" }}>
        Unlocked: <span style={{ color: "#E8EAED" }}>Queues</span>
      </p>
    </div>
  );
}
