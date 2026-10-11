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
      <PhoneFrame>{children}</PhoneFrame>
    </figure>
  );
}

/** PhoneFrame is the bare device frame, usable without a label. */
export function PhoneFrame({ children }: { children: ReactNode }) {
  return (
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
  );
}

/* ─── Hero canvas phone screens ────────────────────────────────── */

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
      <CompanionHero name="Milo" species="cat" size="sm" reaction="correct" />
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

/* ─── Sticky showcase phone screens ────────────────────────────── */

/** IntroPhoneScreen — beat-by-beat lesson intro stage illustration. */
export function IntroPhoneScreen() {
  const beats = [
    { text: "A cache sits between your app", em: false },
    { text: "and the database.", em: false },
    { text: "When the same data is requested twice,", em: false },
    { text: "the second read is free.", em: true },
  ];
  return (
    <div className="flex h-full flex-col px-3 pt-3">
      <div className="flex items-center gap-2 mb-3">
        <span className="text-[10px]" style={{ color: "#8B92A8" }}>✕</span>
        <div className="h-1.5 flex-1 overflow-hidden rounded-full" style={{ background: "#242836" }}>
          <div className="h-full w-1/5 rounded-full" style={{ background: "linear-gradient(90deg, #7C6EF5, #9D8FF7)" }} />
        </div>
      </div>
      {/* Visual bar */}
      <div className="mb-3 rounded-2xl p-3 flex items-center justify-center" style={{ background: "#1A1D27", border: "1px solid #2E3347", height: 80 }}>
        <div className="flex items-end gap-1.5">
          {[40, 70, 55, 88, 62, 44, 76].map((h, i) => (
            <div key={i} className="w-3 rounded-sm" style={{ height: h * 0.55, background: i === 4 ? "#7C6EF5" : "rgba(124,110,245,0.3)" }} />
          ))}
        </div>
      </div>
      <p className="font-mono text-[8px] font-bold uppercase tracking-[0.18em] mb-2" style={{ color: "#7C6EF5" }}>
        intro · 1 of 1
      </p>
      <div className="flex-1 space-y-1.5">
        {beats.map((b, i) => (
          <p
            key={i}
            className="font-display text-[12px] font-bold leading-snug"
            style={{ color: b.em ? "#7C6EF5" : "#E8EAED", opacity: i < 3 ? 1 : 0.6 }}
          >
            {b.text}
          </p>
        ))}
      </div>
      <div className="mt-auto mb-4 flex items-center justify-between">
        <span className="text-[10px]" style={{ color: "#8B92A8" }}>⏸</span>
        <div className="font-display rounded-full px-5 py-2 text-[11px] font-extrabold" style={{ background: "#7C6EF5", color: "#fff" }}>
          Next →
        </div>
      </div>
    </div>
  );
}

/** GradingPhoneScreen — answer submitted, graded result returned. */
export function GradingPhoneScreen() {
  return (
    <div className="flex h-full flex-col px-3 pt-3">
      <div className="flex items-center gap-2 mb-4">
        <span className="text-[10px]" style={{ color: "#8B92A8" }}>✕</span>
        <div className="h-1.5 flex-1 overflow-hidden rounded-full" style={{ background: "#242836" }}>
          <div className="h-full w-2/3 rounded-full" style={{ background: "linear-gradient(90deg, #7C6EF5, #9D8FF7)" }} />
        </div>
      </div>
      <p className="font-display text-[13px] font-bold leading-snug mb-3" style={{ color: "#E8EAED" }}>
        Which guarantee prevents double-charging a customer on retry?
      </p>
      {["Idempotency key", "Distributed lock", "Saga pattern"].map((o, i) => (
        <div
          key={o}
          className="flex items-center gap-2 rounded-xl px-2.5 py-2 text-[10px] font-semibold mb-2"
          style={{
            background: i === 0 ? "rgba(52,211,153,0.1)" : "#1A1D27",
            border: `1.5px solid ${i === 0 ? "#34D399" : "#2E3347"}`,
            color: "#E8EAED",
          }}
        >
          <span className="font-mono flex h-4 w-4 items-center justify-center rounded text-[8px]" style={{ background: i === 0 ? "#34D399" : "#242836", color: i === 0 ? "#0F1117" : "#8B92A8" }}>
            {i === 0 ? "✓" : i + 1}
          </span>
          {o}
        </div>
      ))}
      {/* Feedback tray */}
      <div className="mt-auto rounded-2xl px-3 py-3" style={{ background: "#10231C", border: "2px solid #34D399" }}>
        <p className="font-display flex items-center gap-1.5 text-[11px] font-extrabold" style={{ color: "#34D399" }}>
          <span className="flex h-4 w-4 items-center justify-center rounded-full text-[9px]" style={{ background: "#34D399", color: "#0F1117" }}>✓</span>
          Nice, that&apos;s right
        </p>
        <p className="mt-1 text-[9px] leading-snug" style={{ color: "#B9C0D0" }}>
          An idempotency key lets the server detect a replay and return the same result, so the client can retry safely.
        </p>
      </div>
    </div>
  );
}

/** MistakePhoneScreen — wrong answer, miss flash, re-queue badge. */
export function MistakePhoneScreen() {
  return (
    <div className="flex h-full flex-col px-3 pt-3">
      <div className="flex items-center gap-2 mb-4">
        <span className="text-[10px]" style={{ color: "#8B92A8" }}>✕</span>
        <div className="h-1.5 flex-1 overflow-hidden rounded-full" style={{ background: "#242836" }}>
          <div className="h-full w-1/2 rounded-full" style={{ background: "linear-gradient(90deg, #7C6EF5, #9D8FF7)" }} />
        </div>
        <span className="font-mono text-[9px] font-bold" style={{ color: "#F5B942" }}>🔥 5</span>
      </div>
      <p className="font-display text-[13px] font-bold leading-snug mb-3" style={{ color: "#E8EAED" }}>
        What does the Strategy pattern primarily help with?
      </p>
      {["Swap algorithms at runtime", "Add logging globally", "Reduce class count"].map((o, i) => (
        <div
          key={o}
          className={`flex items-center gap-2 rounded-xl px-2.5 py-2 text-[10px] font-semibold mb-2 ${i === 2 ? s.optionWrong : ""}`}
          style={{
            background: "#1A1D27",
            border: "1.5px solid #2E3347",
            color: "#E8EAED",
          }}
        >
          <span className="font-mono flex h-4 w-4 items-center justify-center rounded text-[8px]" style={{ background: "#242836", color: "#8B92A8" }}>
            {i + 1}
          </span>
          {o}
        </div>
      ))}
      {/* Miss tray */}
      <div className="mt-auto rounded-2xl px-3 py-3" style={{ background: "#1E0F0F", border: "2px solid #EF4444" }}>
        <p className="font-display flex items-center gap-1.5 text-[11px] font-extrabold" style={{ color: "#EF4444" }}>
          <span className="flex h-4 w-4 items-center justify-center rounded-full text-[9px]" style={{ background: "#EF4444", color: "#fff" }}>✕</span>
          Not quite
        </p>
        <p className="mt-1 text-[9px]" style={{ color: "#B9C0D0" }}>Correct: <span style={{ color: "#E8EAED" }}>Swap algorithms at runtime</span></p>
        <div className={`mt-2 inline-flex items-center gap-1 rounded-full px-2.5 py-1 text-[9px] font-semibold ${s.requeueBadge}`} style={{ background: "rgba(124,110,245,0.18)", color: "#A78BFA" }}>
          ↺ Coming back at the end
        </div>
      </div>
    </div>
  );
}

/** MasteryPhoneScreen — celebration with mastery ring and skill delta bars. */
export function MasteryPhoneScreen() {
  const theme = topicTheme("system-design");
  return (
    <div className="flex h-full flex-col items-center px-3 pt-4 text-center">
      <CompanionHero name="Milo" species="cat" size="sm" reaction="correct" />
      <p className="font-display mt-1 text-[14px] font-extrabold" style={{ color: "#E8EAED" }}>Lesson complete!</p>
      <div className="mt-2 flex gap-1.5">
        {["+20 XP", "💎 +10", "🔥 6"].map((b, i) => (
          <span key={b} className={`font-mono rounded-lg px-2 py-1 text-[9px] font-bold ${s.pop}`} style={{ animationDelay: `${i * 0.15}s`, background: i === 0 ? "rgba(96,165,250,0.14)" : i === 1 ? "rgba(52,211,153,0.14)" : "rgba(255,107,53,0.14)", color: i === 0 ? "#60A5FA" : i === 1 ? "#34D399" : "#FF6B35" }}>
            {b}
          </span>
        ))}
      </div>
      {/* Topic ring card */}
      <div className="mt-3 w-full rounded-2xl p-3 text-left" style={{ background: "#1A1D27", border: "1px solid #2E3347" }}>
        <div className="flex items-center gap-3">
          <MasteryRing mastery={75} tint={theme.tint} size={48} />
          <div>
            <p className="font-display text-[10px] font-bold" style={{ color: "#E8EAED" }}>System Design</p>
            <p className="text-[9px]" style={{ color: "#34D399" }}>↑ 72% → 75%</p>
          </div>
        </div>
        {[{ skill: "Caching", w: "75%", d: "+3" }, { skill: "Load Balancing", w: "42%", d: "+2" }].map((m) => (
          <div key={m.skill} className="mt-2">
            <div className="flex justify-between text-[9px] mb-0.5" style={{ color: "#8B92A8" }}>
              <span>{m.skill}</span>
              <span style={{ color: "#34D399" }}>{m.d}</span>
            </div>
            <div className="h-1.5 overflow-hidden rounded-full" style={{ background: "#242836" }}>
              <div className={`h-full rounded-full ${s.grow}`} style={{ width: m.w, background: theme.tint }} />
            </div>
          </div>
        ))}
      </div>
      <p className="mt-2 text-[9px]" style={{ color: "#8B92A8" }}>
        Unlocked: <span style={{ color: "#E8EAED" }}>Queues</span>
      </p>
    </div>
  );
}
