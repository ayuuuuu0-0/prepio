# PREPIO_REDESIGN.md

## Purpose

This document contains every change needed to give Prepio a distinct visual
identity that matches its audience: engineers and CS students chasing real
career outcomes, not children learning a hobby.

The changes are organized into three iterations. Do them in order.
Each iteration is independently shippable.

Current problem in one line:
**The app looks like Duolingo with a capybara instead of an owl.**

Target identity in one line:
**A career RPG built by engineers for engineers — ambitious, dark, precise,
and alive.**

Duolingo's audience is a 35-year-old learning French for vacation.
Prepio's audience is a 22-year-old grinding LeetCode to get into Google.
The UI must reflect the difference in stakes.

---

## What to Keep

These things are correct and should not change:

- Companion system — emotionally right, keep the idle animation
- Rounded corners throughout — appropriate, not childish
- Bottom tab navigation on mobile — right pattern
- `ReadinessRing` SVG rings — genuinely distinctive
- `QuestCard` progress bar pattern — clean
- `GameButton` pill shape — genre-appropriate, just recolor
- The 3-step onboarding flow — correct structure
- Font pairing concept — display + body is right, swap the specific choices
- All backend logic — zero backend changes in any iteration

---

## Iteration 1 — Token Swap

**Time:** 2–3 hours
**What changes:** CSS variables, colors, fonts, component styles
**What does not change:** Layout, structure, logic, backend
**Visual impact:** Immediately looks like a different product

This is pure style. Every file listed below gets exact replacement content.

---

### 1.1 — `web/next.config.ts`

No change needed.

---

### 1.2 — `web/src/app/layout.tsx`

**Replace entire file:**

```tsx
import type { Metadata } from "next";
import { Plus_Jakarta_Sans, Nunito, JetBrains_Mono } from "next/font/google";
import "./globals.css";

const plusJakarta = Plus_Jakarta_Sans({
  variable: "--font-display",
  subsets: ["latin"],
  weight: ["500", "600", "700", "800"],
});

const nunito = Nunito({
  variable: "--font-body",
  subsets: ["latin"],
  weight: ["400", "500", "600", "700"],
});

const jetbrainsMono = JetBrains_Mono({
  variable: "--font-mono",
  subsets: ["latin"],
  weight: ["400", "500", "600"],
});

export const metadata: Metadata = {
  title: "Prepio — Level Up Your Career",
  description: "A career RPG where interview prep is the progression mechanic",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className={`${plusJakarta.variable} ${nunito.variable} ${jetbrainsMono.variable} antialiased`}
      >
        {children}
      </body>
    </html>
  );
}
```

**Why:** Fredoka is a children's learning app font — it's what Duolingo uses.
Plus Jakarta Sans is modern, rounded but professional. JetBrains Mono on
numbers signals "developer product" immediately. Engineers recognize it.

---

### 1.3 — `web/src/app/globals.css`

**Replace entire file:**

```css
@import "tailwindcss";

/* ─── Design Tokens ──────────────────────────────────────────────── */
:root {
  /* Backgrounds */
  --prepio-bg:          #0F1117;
  --prepio-surface:     #1A1D27;
  --prepio-raised:      #242836;
  --prepio-border:      #2E3347;

  /* Accent — electric purple, not Duolingo green */
  --prepio-accent:      #7C6EF5;
  --prepio-accent-dim:  rgba(124, 110, 245, 0.15);
  --prepio-accent-glow: rgba(124, 110, 245, 0.35);

  /* Game indicators */
  --prepio-streak:      #FF6B35;   /* burnt orange — serious, not cartoon yellow */
  --prepio-gems:        #34D399;   /* emerald — wealth, not Duolingo lime */
  --prepio-gold:        #F5B942;   /* achievement gold */
  --prepio-xp:          #60A5FA;   /* cool blue for XP */

  /* Semantic */
  --prepio-success:     #34D399;
  --prepio-warning:     #F5B942;
  --prepio-danger:      #F87171;

  /* Text */
  --prepio-text:        #E8EAED;
  --prepio-text-muted:  #8B92A8;
  --prepio-text-dim:    #4A5068;

  /* Fonts */
  --font-display: var(--font-display);
  --font-body:    var(--font-body);
  --font-mono:    var(--font-mono);
}

@theme inline {
  --font-display: var(--font-display);
  --font-body: var(--font-body);
  --font-mono: var(--font-mono);
}

/* ─── Base ───────────────────────────────────────────────────────── */
body {
  background-color: var(--prepio-bg);
  color: var(--prepio-text);
  font-family: var(--font-body), system-ui, sans-serif;
}

.font-display {
  font-family: var(--font-display), system-ui, sans-serif;
}

.font-mono {
  font-family: var(--font-mono), monospace;
}

/* ─── Game backgrounds ───────────────────────────────────────────── */

/* Default: dark with purple + teal ambient glow */
.game-bg {
  background:
    radial-gradient(ellipse 60% 40% at 10% 10%, rgba(124, 110, 245, 0.10) 0%, transparent 70%),
    radial-gradient(ellipse 50% 40% at 90% 85%, rgba(52, 211, 153, 0.06) 0%, transparent 70%),
    var(--prepio-bg);
  min-height: 100vh;
}

/* Forest world: deep forest night, not cartoon green */
.game-bg-forest {
  background:
    radial-gradient(ellipse 60% 50% at 50% 0%, rgba(52, 211, 153, 0.12) 0%, transparent 70%),
    linear-gradient(180deg, #0D1F14 0%, #0F1A10 50%, #0F1117 100%);
  min-height: 100vh;
}

/* Challenge/focus: concentrated dark blue-purple */
.game-bg-challenge {
  background:
    radial-gradient(ellipse 70% 50% at 50% 0%, rgba(124, 110, 245, 0.12) 0%, transparent 70%),
    linear-gradient(180deg, #0F1020 0%, #0F1117 60%);
  min-height: 100vh;
}

/* ─── Subtle noise texture overlay (applied to .game-bg variants) ── */
.game-bg::before,
.game-bg-forest::before,
.game-bg-challenge::before {
  content: "";
  position: fixed;
  inset: 0;
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 200 200' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23n)' opacity='0.03'/%3E%3C/svg%3E");
  pointer-events: none;
  z-index: 0;
}

/* ─── Companion animations ───────────────────────────────────────── */
@keyframes companion-idle {
  0%, 100% { transform: translateY(0) scale(1); }
  50%       { transform: translateY(-6px) scale(1.02); }
}

.animate-companion-idle {
  animation: companion-idle 3s ease-in-out infinite;
}

@keyframes companion-react-correct {
  0%   { transform: translateY(0) scale(1) rotate(0deg); }
  25%  { transform: translateY(-14px) scale(1.1) rotate(-5deg); }
  50%  { transform: translateY(-10px) scale(1.08) rotate(5deg); }
  75%  { transform: translateY(-14px) scale(1.1) rotate(-3deg); }
  100% { transform: translateY(0) scale(1) rotate(0deg); }
}

.animate-companion-correct {
  animation: companion-react-correct 0.7s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
}

@keyframes companion-react-wrong {
  0%   { transform: translateX(0) rotate(0deg); }
  20%  { transform: translateX(-6px) rotate(-4deg); }
  40%  { transform: translateX(6px) rotate(4deg); }
  60%  { transform: translateX(-4px) rotate(-2deg); }
  80%  { transform: translateX(4px) rotate(2deg); }
  100% { transform: translateX(0) rotate(0deg); }
}

.animate-companion-wrong {
  animation: companion-react-wrong 0.6s ease-out forwards;
}

/* ─── Journey node animations ────────────────────────────────────── */
@keyframes node-pulse {
  0%, 100% { box-shadow: 0 0 0 0 var(--prepio-accent-glow); }
  50%       { box-shadow: 0 0 0 14px transparent; }
}

.animate-node-pulse {
  animation: node-pulse 2.2s ease-in-out infinite;
}

/* ─── Readiness ring ─────────────────────────────────────────────── */
@keyframes ring-fill {
  from { stroke-dashoffset: 283; }
}

.animate-ring-fill {
  animation: ring-fill 1.4s cubic-bezier(0.4, 0, 0.2, 1) forwards;
}

/* ─── Result screen ──────────────────────────────────────────────── */
@keyframes result-enter {
  0%   { transform: scale(0.92) translateY(12px); opacity: 0; }
  100% { transform: scale(1) translateY(0); opacity: 1; }
}

.animate-result {
  animation: result-enter 0.45s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
}

@keyframes xp-count {
  from { opacity: 0; transform: translateY(8px) scale(0.9); }
  to   { opacity: 1; transform: translateY(0) scale(1); }
}

.animate-xp {
  animation: xp-count 0.4s 0.2s cubic-bezier(0.34, 1.56, 0.64, 1) both;
}

.animate-gems {
  animation: xp-count 0.4s 0.35s cubic-bezier(0.34, 1.56, 0.64, 1) both;
}

/* ─── Button ─────────────────────────────────────────────────────── */
.game-btn {
  transition: transform 0.1s ease, box-shadow 0.15s ease, opacity 0.15s ease;
  box-shadow: 0 4px 0 var(--btn-shadow, #5B50D4);
}

.game-btn:hover:not(:disabled) {
  transform: translateY(-1px) scale(1.02);
  box-shadow: 0 6px 0 var(--btn-shadow, #5B50D4);
}

.game-btn:active:not(:disabled) {
  transform: translateY(2px) scale(0.98);
  box-shadow: 0 2px 0 var(--btn-shadow, #5B50D4);
}

.game-btn:disabled {
  opacity: 0.45;
  box-shadow: none;
  cursor: not-allowed;
}

/* ─── Speech bubble tail ─────────────────────────────────────────── */
.speech-bubble::after {
  content: "";
  position: absolute;
  bottom: -10px;
  left: 24px;
  border-width: 10px 10px 0;
  border-style: solid;
  border-color: var(--prepio-surface) transparent transparent;
}

/* ─── XP bar ─────────────────────────────────────────────────────── */
@keyframes xp-bar-fill {
  from { width: 0%; }
}

.animate-xp-bar {
  animation: xp-bar-fill 1.2s cubic-bezier(0.4, 0, 0.2, 1) forwards;
}

/* ─── Stagger utility ────────────────────────────────────────────── */
.stagger-1 { animation-delay: 0.05s; }
.stagger-2 { animation-delay: 0.10s; }
.stagger-3 { animation-delay: 0.15s; }
.stagger-4 { animation-delay: 0.20s; }
```

---

### 1.4 — `web/src/lib/design/tokens.ts`

**Replace entire file:**

```ts
/** Prepio design tokens — Career RPG for engineers, not children's language app. */

export const MIN_ANSWER_LENGTH = 100;

export const colors = {
  /* Backgrounds */
  bg:           "#0F1117",
  surface:      "#1A1D27",
  raised:       "#242836",
  border:       "#2E3347",

  /* Accent */
  accent:       "#7C6EF5",
  accentDim:    "rgba(124,110,245,0.15)",
  accentGlow:   "rgba(124,110,245,0.35)",

  /* Game indicators */
  streak:       "#FF6B35",
  gems:         "#34D399",
  xp:           "#60A5FA",
  gold:         "#F5B942",

  /* Semantic */
  success:      "#34D399",
  warning:      "#F5B942",
  danger:       "#F87171",

  /* Text */
  textPrimary:  "#E8EAED",
  textMuted:    "#8B92A8",
  textDim:      "#4A5068",
} as const;

export const companyColors: Record<string, { ring: string; bg: string; text: string }> = {
  google:    { ring: "#4285F4", bg: "rgba(66,133,244,0.12)",   text: "#7EB3FF" },
  amazon:    { ring: "#FF9900", bg: "rgba(255,153,0,0.12)",    text: "#FFB84D" },
  meta:      { ring: "#7C6EF5", bg: "rgba(124,110,245,0.12)",  text: "#A99EFA" },
  uber:      { ring: "#34D399", bg: "rgba(52,211,153,0.12)",   text: "#6EE7B7" },
  atlassian: { ring: "#0052CC", bg: "rgba(0,82,204,0.12)",     text: "#5B9BD5" },
  netflix:   { ring: "#E50914", bg: "rgba(229,9,20,0.12)",     text: "#FF6B6B" },
};

export const leagueThemes: Record<string, { gradient: string; icon: string; border: string }> = {
  bronze:   { gradient: "from-amber-900/80 to-amber-700/80",   icon: "🥉", border: "#92400E" },
  silver:   { gradient: "from-slate-700/80 to-slate-500/80",   icon: "🥈", border: "#64748B" },
  gold:     { gradient: "from-yellow-700/80 to-amber-500/80",  icon: "🥇", border: "#D97706" },
  sapphire: { gradient: "from-blue-800/80 to-blue-600/80",     icon: "💎", border: "#2563EB" },
  ruby:     { gradient: "from-rose-800/80 to-rose-600/80",     icon: "♦️",  border: "#E11D48" },
  emerald:  { gradient: "from-emerald-800/80 to-emerald-600/80", icon: "💚", border: "#059669" },
  diamond:  { gradient: "from-sky-700/80 to-cyan-600/80",      icon: "✨", border: "#0891B2" },
  legend:   { gradient: "from-yellow-900/80 to-amber-800/80",  icon: "👑", border: "#B45309" },
};

export const roundTypeColors: Record<string, string> = {
  dsa:           "#7C6EF5",
  system_design: "#60A5FA",
  lld:           "#34D399",
  behavioral:    "#FF6B35",
  aptitude:      "#F5B942",
  fundamentals:  "#A99EFA",
};
```

**Why `MIN_ANSWER_LENGTH = 100`:** The current value of 3 means users can
type three characters and evaluate. A meaningful technical answer requires
at least a sentence. 100 characters is a short paragraph — the right floor.

---

### 1.5 — `web/src/components/game/GameBackground.tsx`

**Replace entire file:**

```tsx
import { ReactNode } from "react";

type Variant = "default" | "forest" | "challenge";

const variants: Record<Variant, string> = {
  default:   "game-bg",
  forest:    "game-bg-forest",
  challenge: "game-bg-challenge",
};

export function GameBackground({
  children,
  variant = "default",
  className = "",
}: {
  children: ReactNode;
  variant?: Variant;
  className?: string;
}) {
  return (
    <div className={`${variants[variant]} relative ${className}`}>
      {/* Ambient orbs — subtle depth, not Duolingo blobs */}
      <div className="pointer-events-none fixed inset-0 overflow-hidden" aria-hidden>
        <div
          className="absolute -left-20 -top-20 h-64 w-64 rounded-full opacity-20 blur-3xl"
          style={{ background: "radial-gradient(circle, #7C6EF5 0%, transparent 70%)" }}
        />
        <div
          className="absolute -right-20 bottom-1/3 h-48 w-48 rounded-full opacity-10 blur-3xl"
          style={{ background: "radial-gradient(circle, #34D399 0%, transparent 70%)" }}
        />
      </div>
      <div className="relative z-10">{children}</div>
    </div>
  );
}
```

---

### 1.6 — `web/src/components/game/CompanionHero.tsx`

**Replace entire file:**

```tsx
import { companionVisual } from "@/lib/design/companions";

type Reaction = "idle" | "correct" | "wrong";

export function CompanionHero({
  name,
  species,
  size = "lg",
  reaction = "idle",
}: {
  name?: string;
  species?: string;
  size?: "sm" | "md" | "lg";
  reaction?: Reaction;
}) {
  const visual = companionVisual(name, species);

  const sizeClass =
    size === "lg" ? "h-28 w-28 text-5xl" :
    size === "md" ? "h-20 w-20 text-4xl" :
    "h-12 w-12 text-2xl";

  const animClass =
    reaction === "correct" ? "animate-companion-correct" :
    reaction === "wrong"   ? "animate-companion-wrong" :
    "animate-companion-idle";

  return (
    <div className="relative inline-flex flex-col items-center">
      {/* Glow ring */}
      <div
        className={`absolute inset-0 rounded-full blur-md opacity-40 ${sizeClass}`}
        style={{ background: visual.glow, transform: "scale(1.2)" }}
        aria-hidden
      />
      {/* Companion circle */}
      <div
        className={`relative flex items-center justify-center rounded-full border border-white/10 ${sizeClass} ${animClass}`}
        style={{
          background: `linear-gradient(135deg, ${visual.glow}33 0%, #1A1D27 100%)`,
          boxShadow: `0 0 24px ${visual.glow}44, inset 0 1px 0 rgba(255,255,255,0.08)`,
        }}
      >
        <span role="img" aria-label={name ?? "companion"} className="select-none">
          {visual.emoji}
        </span>
      </div>
    </div>
  );
}
```

**Why:** The current version is a gradient circle with an emoji — visually a
Duolingo owl in a different shape. The new version has a glow ring, a dark
surface treatment, and a subtle inner border. It reads "premium game character"
not "children's app mascot."

---

### 1.7 — `web/src/components/game/SpeechBubble.tsx`

**Replace entire file:**

```tsx
import { ReactNode } from "react";

export function SpeechBubble({
  children,
  className = "",
  speakerName,
}: {
  children: ReactNode;
  className?: string;
  speakerName?: string;
}) {
  return (
    <div
      className={`speech-bubble relative rounded-2xl border border-white/8 px-4 py-3 text-sm leading-relaxed ${className}`}
      style={{
        background: "#1A1D27",
        boxShadow: "0 4px 20px rgba(0,0,0,0.4), inset 0 1px 0 rgba(255,255,255,0.05)",
      }}
    >
      {speakerName && (
        <p
          className="font-mono mb-1 text-xs font-semibold uppercase tracking-widest"
          style={{ color: "#7C6EF5" }}
        >
          {speakerName}
        </p>
      )}
      <p className="font-body font-medium" style={{ color: "#C8CCDA" }}>
        {children}
      </p>
    </div>
  );
}
```

**Why:** White speech bubble is Duolingo's signature pattern. Dark surface
with a subtle accent speaker label reads as an in-game HUD, not a children's
messaging interface. Fredoka font is removed from here — Plus Jakarta Sans is
more appropriate for the tone.

---

### 1.8 — `web/src/components/game/GameButton.tsx`

**Replace entire file:**

```tsx
import { ButtonHTMLAttributes, ReactNode } from "react";

type Variant = "primary" | "secondary" | "gold" | "ghost";

interface BtnStyle {
  bg: string;
  shadow: string;
  text: string;
}

const styles: Record<Variant, BtnStyle> = {
  primary: {
    bg: "bg-gradient-to-r from-[#7C6EF5] to-[#9D8FF7] hover:from-[#8B7FF7] hover:to-[#AFA0FF]",
    shadow: "#5B50D4",
    text: "text-white",
  },
  secondary: {
    bg: "bg-[#242836] hover:bg-[#2E3347] border border-[#2E3347] hover:border-[#7C6EF5]/50",
    shadow: "#1A1D27",
    text: "text-[#C8CCDA]",
  },
  gold: {
    bg: "bg-gradient-to-r from-[#F5B942] to-[#FFCF60] hover:from-[#FFCF60] hover:to-[#FFD97A]",
    shadow: "#C8941A",
    text: "text-[#1A1500]",
  },
  ghost: {
    bg: "bg-transparent hover:bg-white/5 border border-white/10 hover:border-white/20",
    shadow: "transparent",
    text: "text-[#8B92A8] hover:text-[#C8CCDA]",
  },
};

export function GameButton({
  children,
  variant = "primary",
  className = "",
  ...props
}: ButtonHTMLAttributes<HTMLButtonElement> & {
  children: ReactNode;
  variant?: Variant;
}) {
  const s = styles[variant];
  return (
    <button
      className={`game-btn font-display w-full rounded-full px-8 py-4 text-base font-bold tracking-wide ${s.bg} ${s.text} ${className}`}
      style={{ ["--btn-shadow" as string]: s.shadow }}
      {...props}
    >
      {children}
    </button>
  );
}
```

---

### 1.9 — `web/src/components/game/GameCard.tsx`

**Replace entire file:**

```tsx
import { ReactNode } from "react";

export function GameCard({
  children,
  icon,
  className = "",
  accentColor,
}: {
  children: ReactNode;
  icon?: string;
  className?: string;
  /** Left border accent color — makes each card type visually distinct */
  accentColor?: string;
}) {
  return (
    <div
      className={`overflow-hidden rounded-2xl p-5 ${className}`}
      style={{
        background: "#1A1D27",
        border: "1px solid #2E3347",
        boxShadow: "0 4px 24px rgba(0,0,0,0.3)",
        borderLeft: accentColor ? `3px solid ${accentColor}` : "1px solid #2E3347",
      }}
    >
      {icon && (
        <div className="mb-3 text-2xl">{icon}</div>
      )}
      {children}
    </div>
  );
}
```

---

### 1.10 — `web/src/components/game/StatChip.tsx`

**Replace entire file:**

```tsx
export function StatChip({
  icon,
  label,
  value,
  color,
}: {
  icon: string;
  label: string;
  value: string | number;
  color: string;
}) {
  return (
    <div
      className="flex items-center gap-2.5 rounded-xl px-3 py-2"
      style={{
        background: `${color}14`,
        border: `1px solid ${color}30`,
      }}
    >
      <span className="text-base leading-none">{icon}</span>
      <div>
        <p
          className="font-display text-[10px] font-bold uppercase tracking-widest"
          style={{ color: `${color}99` }}
        >
          {label}
        </p>
        <p className="font-mono text-sm font-bold" style={{ color }}>
          {value}
        </p>
      </div>
    </div>
  );
}
```

**Why JetBrains Mono here:** Numbers in monospace is the single fastest signal
that this is a developer product. Streak: `12` in JetBrains Mono looks like a
metric. In Fredoka it looks like a children's score.

---

### 1.11 — `web/src/components/game/QuestCard.tsx`

**Replace entire file:**

```tsx
export function QuestCard({
  title,
  icon,
  progress,
  target,
  completed,
  rewardXp,
  rewardGems,
}: {
  title: string;
  icon: string;
  progress: number;
  target: number;
  completed: boolean;
  rewardXp: number;
  rewardGems: number;
}) {
  const pct = Math.min(100, (progress / target) * 100);

  return (
    <div
      className="rounded-2xl p-4 transition-all"
      style={{
        background: completed ? "rgba(52,211,153,0.08)" : "#1A1D27",
        border: `1px solid ${completed ? "rgba(52,211,153,0.3)" : "#2E3347"}`,
      }}
    >
      <div className="flex items-start gap-3">
        <span className="mt-0.5 text-xl">{completed ? "✅" : icon}</span>
        <div className="flex-1 min-w-0">
          <p
            className="font-display text-sm font-bold"
            style={{
              color: completed ? "#34D399" : "#E8EAED",
              textDecoration: completed ? "line-through" : "none",
              opacity: completed ? 0.7 : 1,
            }}
          >
            {title}
          </p>
          <div
            className="mt-2 h-1.5 overflow-hidden rounded-full"
            style={{ background: "#2E3347" }}
          >
            <div
              className="h-full rounded-full transition-all duration-700"
              style={{
                width: `${pct}%`,
                background: completed
                  ? "#34D399"
                  : "linear-gradient(90deg, #7C6EF5, #9D8FF7)",
              }}
            />
          </div>
          <p className="font-mono mt-1.5 text-xs" style={{ color: "#4A5068" }}>
            {progress}/{target} · ⚡ {rewardXp} XP · 💎 {rewardGems}
          </p>
        </div>
      </div>
    </div>
  );
}
```

---

### 1.12 — `web/src/components/game/ReadinessRing.tsx`

**Replace entire file:**

```tsx
export function ReadinessRing({
  company,
  score,
  color,
  delay = 0,
}: {
  company: string;
  score: number;
  color: string;
  delay?: number;
}) {
  const radius = 44;
  const circumference = 2 * Math.PI * radius;
  const offset = circumference - (score / 100) * circumference;

  return (
    <div className="flex flex-col items-center gap-2">
      <div className="relative h-24 w-24">
        {/* Glow behind the ring */}
        <div
          className="absolute inset-2 rounded-full blur-md opacity-20"
          style={{ background: color }}
          aria-hidden
        />
        <svg className="relative h-24 w-24 -rotate-90" viewBox="0 0 100 100">
          {/* Track */}
          <circle
            cx="50" cy="50" r={radius}
            fill="none"
            stroke="#2E3347"
            strokeWidth="7"
          />
          {/* Progress */}
          <circle
            cx="50" cy="50" r={radius}
            fill="none"
            stroke={color}
            strokeWidth="7"
            strokeLinecap="round"
            strokeDasharray={circumference}
            strokeDashoffset={offset}
            className="animate-ring-fill"
            style={{ animationDelay: `${delay}ms` }}
          />
        </svg>
        <div className="absolute inset-0 flex flex-col items-center justify-center">
          <span
            className="font-mono text-base font-bold leading-none"
            style={{ color }}
          >
            {score}%
          </span>
        </div>
      </div>
      <span
        className="font-display text-[11px] font-semibold uppercase tracking-wider capitalize"
        style={{ color: "#8B92A8" }}
      >
        {company}
      </span>
    </div>
  );
}
```

---

### 1.13 — `web/src/components/game/BottomNav.tsx`

**Replace entire file:**

```tsx
"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

const tabs = [
  { href: "/dashboard", label: "Home",    icon: "⊞" },
  { href: "/journey",   label: "Journey", icon: "◈" },
  { href: "/challenge", label: "Play",    icon: "▶" },
  { href: "#",          label: "League",  icon: "⬡", soon: true },
  { href: "#",          label: "Profile", icon: "◉", soon: true },
];

export function BottomNav() {
  const path = usePathname();

  return (
    <nav
      className="fixed bottom-0 left-0 right-0 z-50"
      style={{
        background: "rgba(15,17,23,0.95)",
        borderTop: "1px solid #2E3347",
        backdropFilter: "blur(20px)",
      }}
    >
      <div className="mx-auto flex max-w-lg items-center justify-around px-2 py-2">
        {tabs.map((tab) => {
          const active = path === tab.href;
          return (
            <Link
              key={tab.href}
              href={tab.href}
              className={`relative flex flex-col items-center gap-0.5 rounded-xl px-3 py-2 transition-all ${
                tab.soon ? "pointer-events-none opacity-30" : ""
              }`}
              style={{
                color: active ? "#7C6EF5" : "#4A5068",
                background: active ? "rgba(124,110,245,0.12)" : "transparent",
              }}
            >
              <span className="text-lg leading-none">{tab.icon}</span>
              <span className="font-display text-[10px] font-semibold tracking-wide">
                {tab.label}
              </span>
              {active && (
                <span
                  className="absolute -top-px left-1/4 right-1/4 h-px rounded-full"
                  style={{ background: "#7C6EF5" }}
                />
              )}
            </Link>
          );
        })}
      </div>
    </nav>
  );
}
```

---

## Iteration 2 — Layout Differentiation

**Time:** 1–2 days
**What changes:** Login split-screen, dashboard HUD bar + section reorder,
two-column challenge layout, new components
**What does not change:** Backend logic, API calls, state management

---

### 2.1 — New component: `web/src/components/game/HUDBar.tsx`

Replace the three `StatChip` row on the dashboard with a single HUD bar that
also includes the XP progress bar. This is the most impactful structural change.

**Create new file:**

```tsx
"use client";

import { DashboardHome } from "@/lib/api";

export function HUDBar({ home }: { home: DashboardHome }) {
  const { streak, progress, companion } = home;
  const xpPct = Math.min(
    100,
    Math.round(
      ((progress.total_xp % (progress.total_xp + progress.xp_to_next_level || 100)) /
        (progress.total_xp + progress.xp_to_next_level || 100)) *
        100
    )
  );

  return (
    <div
      className="flex items-center gap-3 rounded-2xl px-4 py-3"
      style={{
        background: "#1A1D27",
        border: "1px solid #2E3347",
      }}
    >
      {/* Companion badge */}
      <div className="flex items-center gap-2 shrink-0">
        <span className="text-xl">
          {companion?.name ? getCompanionEmoji(companion.species) : "🦫"}
        </span>
        <div>
          <p className="font-mono text-[10px] font-bold uppercase tracking-widest" style={{ color: "#7C6EF5" }}>
            {companion?.name ?? "Byte"}
          </p>
          <p className="font-mono text-xs font-bold" style={{ color: "#60A5FA" }}>
            Lv.{progress.current_level}
          </p>
        </div>
      </div>

      {/* Divider */}
      <div className="h-8 w-px shrink-0" style={{ background: "#2E3347" }} />

      {/* Streak */}
      <div className="flex items-center gap-1.5 shrink-0">
        <span className="text-sm">🔥</span>
        <span className="font-mono text-sm font-bold" style={{ color: "#FF6B35" }}>
          {streak.current_streak}
        </span>
      </div>

      {/* Gems */}
      <div className="flex items-center gap-1.5 shrink-0">
        <span className="text-sm">💎</span>
        <span className="font-mono text-sm font-bold" style={{ color: "#34D399" }}>
          {progress.gem_balance}
        </span>
      </div>

      {/* Divider */}
      <div className="h-8 w-px shrink-0" style={{ background: "#2E3347" }} />

      {/* XP bar */}
      <div className="flex-1 min-w-0">
        <div className="flex justify-between items-center mb-1">
          <span className="font-mono text-[10px] font-semibold" style={{ color: "#4A5068" }}>
            XP
          </span>
          <span className="font-mono text-[10px] font-semibold" style={{ color: "#60A5FA" }}>
            {progress.xp_to_next_level} to next
          </span>
        </div>
        <div
          className="h-1.5 w-full overflow-hidden rounded-full"
          style={{ background: "#2E3347" }}
        >
          <div
            className="h-full rounded-full animate-xp-bar"
            style={{
              width: `${xpPct}%`,
              background: "linear-gradient(90deg, #60A5FA, #7C6EF5)",
            }}
          />
        </div>
      </div>
    </div>
  );
}

function getCompanionEmoji(species?: string): string {
  const map: Record<string, string> = {
    capybara: "🦫", red_panda: "🐼", pangolin: "🦔",
    axolotl: "🦎", snow_leopard: "🐆", owl: "🦉",
  };
  return map[species ?? ""] ?? "🦫";
}
```

---

### 2.2 — `web/src/app/dashboard/page.tsx`

**Key structural changes:**
- Replace `<section className="flex items-end gap-4">` companion block with
  a tighter dark card that uses the `speakerName` prop on `SpeechBubble`
- Replace the three `StatChip` row with `<HUDBar home={home} />`
- Change the Career Readiness `GameCard` to use `accentColor="#60A5FA"`
- Change the League `GameCard` to use `accentColor` matching the tier
- Change the quests section heading to not be childish

**Specific replacements inside `dashboard/page.tsx`:**

Replace:
```tsx
{/* Top — Companion + speech */}
<section className="flex items-end gap-4">
  <CompanionHero name={home.companion?.name} species={home.companion?.species} size="lg" />
  <SpeechBubble className="flex-1">{home.companion_message}</SpeechBubble>
</section>

{/* Stats row */}
<div className="mt-5 flex flex-wrap gap-2">
  <StatChip icon="🔥" label="Streak" value={`${home.streak.current_streak} days`} color="#FF9600" />
  <StatChip icon="⚡" label="Level" value={home.progress.current_level} color="#1CB0F6" />
  <StatChip icon="💎" label="Gems" value={home.progress.gem_balance} color="#7B5CFF" />
</div>
```

With:
```tsx
{/* Companion + speech bubble */}
<section className="flex items-start gap-4">
  <CompanionHero
    name={home.companion?.name}
    species={home.companion?.species}
    size="md"
  />
  <SpeechBubble
    className="flex-1 mt-2"
    speakerName={home.companion?.name ?? "Byte"}
  >
    {home.companion_message}
  </SpeechBubble>
</section>

{/* HUD bar — replaces separate stat chips */}
<div className="mt-4">
  <HUDBar home={home} />
</div>
```

Replace the Career Readiness card opening:
```tsx
<GameCard className="mt-5 bg-gradient-to-br from-white to-sky-50" icon="🧭">
  <h2 className="font-display mb-4 text-lg font-bold text-[#3C3C3C]">Career Readiness</h2>
```
With:
```tsx
<GameCard className="mt-5" icon="🧭" accentColor="#60A5FA">
  <h2 className="font-display mb-4 text-base font-bold" style={{ color: "#E8EAED" }}>
    Career Readiness
  </h2>
```

Replace the League card:
```tsx
<GameCard
  className={`mt-4 bg-gradient-to-r ${league.gradient} text-white`}
  icon={league.icon}
>
  <p className="font-display text-sm font-bold uppercase opacity-90">{home.league.label}</p>
  <p className="font-display text-3xl font-extrabold">Rank #{home.league.rank}</p>
  <p className="mt-1 text-sm opacity-80">Keep climbing — promotion zone ahead!</p>
</GameCard>
```
With:
```tsx
<GameCard
  className="mt-4"
  icon={league.icon}
  accentColor={league.border}
>
  <div
    className={`-mx-5 -mt-5 mb-4 rounded-t-2xl bg-gradient-to-r ${league.gradient} px-5 pt-4 pb-3`}
  >
    <p className="font-mono text-[11px] font-bold uppercase tracking-widest text-white/70">
      {home.league.label}
    </p>
    <p className="font-display mt-0.5 text-2xl font-extrabold text-white">
      Rank #{home.league.rank}
    </p>
  </div>
  <p className="text-sm" style={{ color: "#8B92A8" }}>
    Keep climbing — promotion zone ahead!
  </p>
</GameCard>
```

Replace the quests heading:
```tsx
<h2 className="font-display text-lg font-bold text-[#3C3C3C]">⚡ Daily Quests</h2>
```
With:
```tsx
<h2 className="font-display text-sm font-bold uppercase tracking-widest" style={{ color: "#4A5068" }}>
  Daily Quests
</h2>
```

Also add to imports at top of `dashboard/page.tsx`:
```tsx
import { HUDBar } from "@/components/game/HUDBar";
```

And remove the `StatChip` import since it's replaced by `HUDBar`.

---

### 2.3 — `web/src/app/login/page.tsx`

**Replace entire file** — split-screen layout:

```tsx
"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { FormEvent, useState } from "react";
import { GameBackground } from "@/components/game/GameBackground";
import { CompanionHero } from "@/components/game/CompanionHero";
import { SpeechBubble } from "@/components/game/SpeechBubble";
import { GameButton } from "@/components/game/GameButton";
import { api } from "@/lib/api";

const socialProof = [
  { label: "Engineers in prep", value: "12k+" },
  { label: "Avg readiness gain", value: "34%" },
  { label: "Offers received", value: "2.1k" },
];

export default function LoginPage() {
  const router = useRouter();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setLoading(true);
    setError("");
    try {
      const res = await api.login(email, password);
      api.setToken(res.access_token);
      const profile = await api.getProfile();
      router.push(profile.onboarding_completed ? "/dashboard" : "/onboarding");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Login failed");
    } finally {
      setLoading(false);
    }
  }

  return (
    <GameBackground>
      <div className="flex min-h-screen">

        {/* Left — brand panel (hidden on small screens) */}
        <div
          className="hidden lg:flex lg:w-1/2 flex-col justify-between p-12"
          style={{ borderRight: "1px solid #2E3347" }}
        >
          {/* Logo */}
          <div>
            <p className="font-display text-xl font-bold" style={{ color: "#7C6EF5" }}>
              PREPIO
            </p>
            <p className="font-mono text-xs font-medium mt-1" style={{ color: "#4A5068" }}>
              Career RPG
            </p>
          </div>

          {/* Companion + message */}
          <div className="flex flex-col items-center gap-6 py-12">
            <CompanionHero name="Byte" species="capybara" size="lg" />
            <SpeechBubble speakerName="Byte" className="max-w-xs text-center">
              Byte is waiting. Your streak won&apos;t hold itself.
            </SpeechBubble>
          </div>

          {/* Social proof */}
          <div className="grid grid-cols-3 gap-4">
            {socialProof.map((s) => (
              <div key={s.label}>
                <p className="font-mono text-2xl font-bold" style={{ color: "#7C6EF5" }}>
                  {s.value}
                </p>
                <p className="font-body text-xs mt-1" style={{ color: "#4A5068" }}>
                  {s.label}
                </p>
              </div>
            ))}
          </div>
        </div>

        {/* Right — form panel */}
        <div className="flex w-full lg:w-1/2 flex-col items-center justify-center px-6 py-12">
          {/* Mobile logo */}
          <div className="mb-8 text-center lg:hidden">
            <p className="font-display text-2xl font-bold" style={{ color: "#7C6EF5" }}>
              PREPIO
            </p>
            <CompanionHero name="Byte" species="capybara" size="md" />
          </div>

          <div className="w-full max-w-sm">
            <h1 className="font-display text-2xl font-bold mb-1" style={{ color: "#E8EAED" }}>
              Continue prep
            </h1>
            <p className="font-body text-sm mb-8" style={{ color: "#4A5068" }}>
              Sign in to your account
            </p>

            <form onSubmit={onSubmit} className="space-y-4">
              {error && (
                <p
                  className="rounded-xl px-4 py-3 text-sm font-medium"
                  style={{
                    background: "rgba(248,113,113,0.1)",
                    border: "1px solid rgba(248,113,113,0.3)",
                    color: "#F87171",
                  }}
                >
                  {error}
                </p>
              )}

              <input
                className="w-full rounded-xl px-4 py-3 text-sm font-medium outline-none transition-all"
                style={{
                  background: "#1A1D27",
                  border: "1px solid #2E3347",
                  color: "#E8EAED",
                }}
                onFocus={(e) => (e.target.style.borderColor = "#7C6EF5")}
                onBlur={(e) => (e.target.style.borderColor = "#2E3347")}
                type="email"
                placeholder="Email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
              <input
                className="w-full rounded-xl px-4 py-3 text-sm font-medium outline-none transition-all"
                style={{
                  background: "#1A1D27",
                  border: "1px solid #2E3347",
                  color: "#E8EAED",
                }}
                onFocus={(e) => (e.target.style.borderColor = "#7C6EF5")}
                onBlur={(e) => (e.target.style.borderColor = "#2E3347")}
                type="password"
                placeholder="Password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />

              <GameButton type="submit" disabled={loading}>
                {loading ? "Signing in..." : "Continue Prep"}
              </GameButton>
            </form>

            <p className="mt-6 text-center text-sm" style={{ color: "#4A5068" }}>
              No account?{" "}
              <Link
                href="/register"
                className="font-semibold transition-colors"
                style={{ color: "#7C6EF5" }}
                onMouseEnter={(e) => ((e.target as HTMLElement).style.color = "#9D8FF7")}
                onMouseLeave={(e) => ((e.target as HTMLElement).style.color = "#7C6EF5")}
              >
                Join 12k engineers in prep
              </Link>
            </p>
          </div>
        </div>
      </div>
    </GameBackground>
  );
}
```

---

### 2.4 — `web/src/app/register/page.tsx`

Apply the same split-screen treatment. Replace:
- `"Join the adventure!"` → `"Pick a companion. Start tracking your readiness."`
- `"Begin Adventure!"` → `"Create Account"`
- `"Already playing?"` → `"Already have an account?"`
- All white inputs → dark inputs (same pattern as login above)
- `GameBackground` without left panel stays mobile-centered
- Add the left brand panel matching login

The pattern is identical to 2.3. Apply the same structural change.

---

### 2.5 — `web/src/app/challenge/page.tsx`

**Key changes — two-column layout on desktop:**

Replace the outer `<main>` and card wrapper:

```tsx
{/* Was: single column card */}
{/* Now: two columns on md+ screens */}

return (
  <GameBackground variant="challenge">
    {/* Top bar */}
    <div
      className="sticky top-0 z-20 px-4 py-3"
      style={{
        background: "rgba(15,17,23,0.9)",
        borderBottom: "1px solid #2E3347",
        backdropFilter: "blur(20px)",
      }}
    >
      <div className="mx-auto flex max-w-4xl items-center justify-between">
        <Link
          href="/dashboard"
          className="font-mono text-sm font-semibold transition-colors"
          style={{ color: "#4A5068" }}
        >
          ← Home
        </Link>
        <div className="flex items-center gap-3">
          <span
            className="font-mono rounded-lg px-3 py-1 text-xs font-bold"
            style={{
              background: "#1A1D27",
              border: "1px solid #2E3347",
              color: "#8B92A8",
            }}
          >
            {/* Replace Q 1/5 with dynamic value: */}
            Q {currentIndex + 1}/{paper?.questions.length ?? 1}
          </span>
          <CompanionHero size="sm" reaction={result ? (result.correct ? "correct" : "wrong") : "idle"} />
        </div>
      </div>
    </div>

    {/* Two-column body */}
    <main className="mx-auto grid max-w-4xl grid-cols-1 gap-6 px-4 pb-28 pt-6 md:grid-cols-2">

      {/* Left col — question */}
      <div className="space-y-4">
        <SpeechBubble speakerName={home?.companion?.name ?? "Byte"}>
          {getQuestionHint(question)}
        </SpeechBubble>

        {question && (
          <div
            className="rounded-2xl p-6"
            style={{
              background: "#1A1D27",
              border: `1px solid #2E3347`,
              borderLeft: `3px solid ${roundTypeColors[question.round_type] ?? "#7C6EF5"}`,
            }}
          >
            <div className="flex flex-wrap gap-2 mb-4">
              <span
                className="font-mono rounded-lg px-2.5 py-1 text-[11px] font-bold uppercase tracking-wider"
                style={{
                  background: `${roundTypeColors[question.round_type] ?? "#7C6EF5"}18`,
                  color: roundTypeColors[question.round_type] ?? "#7C6EF5",
                }}
              >
                {question.round_type.replace("_", " ")}
              </span>
              <span
                className="font-mono rounded-lg px-2.5 py-1 text-[11px] font-bold uppercase tracking-wider"
                style={{
                  background:
                    question.difficulty === "hard"   ? "rgba(248,113,113,0.12)" :
                    question.difficulty === "medium" ? "rgba(245,185,66,0.12)"  :
                    "rgba(52,211,153,0.12)",
                  color:
                    question.difficulty === "hard"   ? "#F87171" :
                    question.difficulty === "medium" ? "#F5B942" :
                    "#34D399",
                }}
              >
                {question.difficulty}
              </span>
            </div>
            <p className="text-base font-medium leading-relaxed" style={{ color: "#C8CCDA" }}>
              {question.body}
            </p>
          </div>
        )}
      </div>

      {/* Right col — answer + result */}
      <div className="space-y-4">
        {result ? (
          <ResultCard result={result} onNext={handleNext} hasNext={hasNextQuestion} />
        ) : (
          <form onSubmit={onSubmit} className="space-y-3">
            <textarea
              className="w-full rounded-2xl p-4 text-sm font-medium leading-relaxed outline-none transition-all resize-none"
              style={{
                background: "#1A1D27",
                border: "1px solid #2E3347",
                color: "#E8EAED",
                minHeight: "220px",
              }}
              onFocus={(e) => (e.target.style.borderColor = "#7C6EF5")}
              onBlur={(e) => (e.target.style.borderColor = "#2E3347")}
              placeholder="Explain your approach..."
              value={answer}
              onChange={(e) => setAnswer(e.target.value)}
            />
            <p
              className="font-mono text-xs"
              style={{ color: canSubmit ? "#34D399" : "#4A5068" }}
            >
              {canSubmit
                ? "Ready to submit"
                : `${trimmedLength}/${MIN_ANSWER_LENGTH} chars — add more detail`}
            </p>
            <GameButton type="submit" disabled={submitting || !canSubmit}>
              {submitting ? "Evaluating..." : "Submit"}
            </GameButton>
          </form>
        )}
      </div>
    </main>
    <BottomNav />
  </GameBackground>
);
```

Add a `ResultCard` sub-component and `getQuestionHint` helper inside
`challenge/page.tsx`:

```tsx
function ResultCard({
  result,
  onNext,
  hasNext,
}: {
  result: SubmitResponse;
  onNext: () => void;
  hasNext: boolean;
}) {
  return (
    <div
      className="animate-result rounded-2xl p-6 space-y-4"
      style={{
        background: result.correct
          ? "rgba(52,211,153,0.08)"
          : "rgba(245,185,66,0.08)",
        border: `1px solid ${result.correct ? "rgba(52,211,153,0.3)" : "rgba(245,185,66,0.3)"}`,
      }}
    >
      <div className="flex items-center gap-3">
        <span className="text-3xl">{result.correct ? "✓" : "→"}</span>
        <div>
          <p className="font-display text-lg font-bold" style={{ color: result.correct ? "#34D399" : "#F5B942" }}>
            {result.correct ? "Solid." : "Not quite."}
          </p>
          <p className="text-sm" style={{ color: "#8B92A8" }}>
            {result.feedback}
          </p>
        </div>
      </div>

      {result.correct && (
        <div className="flex gap-3">
          <div
            className="animate-xp flex items-center gap-2 rounded-xl px-3 py-2"
            style={{ background: "rgba(96,165,250,0.12)", border: "1px solid rgba(96,165,250,0.2)" }}
          >
            <span className="font-mono text-sm font-bold" style={{ color: "#60A5FA" }}>
              +{result.xp_awarded} XP
            </span>
          </div>
          <div
            className="animate-gems flex items-center gap-2 rounded-xl px-3 py-2"
            style={{ background: "rgba(52,211,153,0.12)", border: "1px solid rgba(52,211,153,0.2)" }}
          >
            <span className="font-mono text-sm font-bold" style={{ color: "#34D399" }}>
              +{result.gems_awarded} 💎
            </span>
          </div>
        </div>
      )}

      {hasNext ? (
        <GameButton type="button" onClick={onNext}>
          Next Question →
        </GameButton>
      ) : (
        <Link href="/dashboard">
          <GameButton type="button" variant="secondary">
            Back to dashboard
          </GameButton>
        </Link>
      )}
    </div>
  );
}

/** Returns a question-specific hint based on round type and difficulty. */
function getQuestionHint(question: Question | null): string {
  if (!question) return "Loading your challenge...";
  const hints: Record<string, string> = {
    dsa:            "Think about time and space complexity before you write.",
    system_design:  "Start with requirements, then components, then trade-offs.",
    lld:            "Focus on class design, responsibilities, and relationships.",
    behavioral:     "Use the STAR format: Situation, Task, Action, Result.",
    aptitude:       "Break the problem into smaller steps first.",
    fundamentals:   "Cover the core concept clearly before diving into detail.",
  };
  return hints[question.round_type] ?? "Take your time. Think it through.";
}
```

Also add multi-question state to `challenge/page.tsx`:

```tsx
// Add these state variables:
const [currentIndex, setCurrentIndex] = useState(0);

// Derive question from index instead of always questions[0]:
// Replace: if (daily.questions.length > 0) setQuestion(daily.questions[0]);
// With:
if (daily.questions.length > 0) setQuestion(daily.questions[currentIndex]);

// Add derived boolean:
const hasNextQuestion = paper !== null && currentIndex < (paper.questions.length - 1);

// Add handler:
function handleNext() {
  if (!paper || !hasNextQuestion) return;
  const next = currentIndex + 1;
  setCurrentIndex(next);
  setQuestion(paper.questions[next]);
  setResult(null);
  setAnswer("");
}
```

Also add import at top:
```tsx
import { roundTypeColors } from "@/lib/design/tokens";
```

---

## Iteration 3 — Voice and Microcopy

**Time:** 2–4 hours
**What changes:** All user-facing strings — buttons, labels, hints, companion
messages, error messages
**What does not change:** Layout, logic, backend

This is the highest ROI per hour after Iteration 1. The same pixel can feel
like a children's app or a serious career product depending solely on the words.

---

### 3.1 — `services/gateway/internal/dashboard/service.go`

**Replace `companionMessage` function:**

```go
func companionMessage(name string, progress ProgressCard) string {
  if len(name) == 0 {
    name = "Your companion"
  }
  xpNeeded := progress.XPToNextLevel
  level := progress.CurrentLevel
  challenges := 0
  if xpNeeded > 0 {
    challenges = (xpNeeded + config.XPByDifficulty["medium"] - 1) / config.XPByDifficulty["medium"]
  }

  switch {
  case progress.TotalXP == 0:
    return fmt.Sprintf("%s is ready. Let's see what you can do.", name)
  case challenges <= 1:
    return fmt.Sprintf("One challenge from Level %d. Don't stop now.", level+1)
  case challenges <= 3:
    return fmt.Sprintf("%d challenges from Level %d. Google Readiness is watching.", challenges, level+1)
  case progress.CurrentLevel < 5:
    return fmt.Sprintf("Level %d. The real prep starts around Level 10 — keep going.", level)
  default:
    return fmt.Sprintf("Level %d. %d challenges from Level %d. Consistency compounds.", level, challenges, level+1)
  }
}
```

**Replace `defaultDailyQuests` function:**

```go
func defaultDailyQuests(streakActive bool) []DailyQuestCard {
  return []DailyQuestCard{
    {
      ID: "daily_question", Title: "Complete today's challenge",
      Progress: 0, Target: 1, Completed: false, RewardXP: 50, RewardGems: 10,
    },
    {
      ID: "maintain_streak", Title: "Keep the streak alive",
      Progress: boolToInt(streakActive), Target: 1, Completed: streakActive,
      RewardXP: 20, RewardGems: 5,
    },
    {
      ID: "score_high", Title: "Score above 80% on a challenge",
      Progress: 0, Target: 1, Completed: false, RewardXP: 30, RewardGems: 5,
    },
  }
}
```

---

### 3.2 — `web/src/app/login/page.tsx`

Already handled in 2.3. Confirmed replacements:

| Was | Now |
|-----|-----|
| `"Welcome back! Ready to level up your career today?"` | `"Byte is waiting. Your streak won't hold itself."` |
| `"New adventurer? Start your journey"` | `"No account? Join 12k engineers in prep"` |
| `"LET'S GO!"` | `"Continue Prep"` |

---

### 3.3 — `web/src/app/register/page.tsx`

| Was | Now |
|-----|-----|
| `"Join the adventure! Pick your companion and start becoming interview-ready."` | `"Pick a companion. Start tracking your readiness from day one."` |
| `"Begin Adventure!"` | `"Create Account"` |
| `"Already playing?"` | `"Already have an account?"` |
| `"Sign in"` link | `"Sign in"` (unchanged, this is fine) |

---

### 3.4 — `web/src/app/onboarding/page.tsx`

Replace `stepMessages` array:

```tsx
const stepMessages = [
  "Which companies are you targeting? We'll personalise your prep.",
  "How much experience do you have? Sets your starting difficulty.",
  "Choose your companion — they'll grow with you throughout the journey.",
];
```

Replace the finish button:
```tsx
{loading ? "Setting up..." : "Start Prep"}
```

Replace the continue buttons:
```tsx
// Step 1 continue
"Next →"    // unchanged, fine

// Step 2 continue  
"Next →"    // unchanged, fine

// Final button (was "Start My Journey!")
"Start Prep"
```

---

### 3.5 — `web/src/app/challenge/page.tsx`

The `getQuestionHint` function in 2.5 already handles this. Confirm these
additional copy changes:

| Was | Now |
|-----|-----|
| `"Preparing challenge..."` | `"Loading challenge..."` |
| `"Write your answer..."` textarea placeholder | `"Explain your approach..."` |
| `"keep writing..."` | `"add more detail"` |
| `"Ready to submit!"` | `"Ready to submit"` |
| `"Checking..."` | `"Evaluating..."` |
| `"SUBMIT ANSWER!"` | `"Submit"` |
| `"Q 1/5"` | Dynamic: `Q {n}/{total}` |

---

### 3.6 — Dashboard section labels

In `dashboard/page.tsx`, replace heading:

```tsx
// Was:
<h2 className="font-display text-lg font-bold text-[#3C3C3C]">⚡ Daily Quests</h2>

// Now (already in 2.2 but confirming the label):
<h2 className="font-mono text-xs font-bold uppercase tracking-widest" style={{ color: "#4A5068" }}>
  Daily Quests
</h2>
```

The `Continue Journey →` CTA button copy:
```tsx
// Was: "Continue Journey →"
// Now:
"Continue Prep →"
```

---

## MD Files to Update

These go in the repo root alongside `agents.md` and `architecture.md`.
Replace the relevant sections in the files below.

---

### Update `DESIGN_SYSTEM.md` — append this section:

```md
# Target Audience

Users are 19–28 year old engineers and CS students.
They are chasing real career outcomes — job offers, salary bumps, FAANG placements.
They respect competence. They distrust condescension.

The UI should feel like:
A senior engineer built a game to help you get their job.

Not like:
A teacher praising a child for finishing homework.

---

# Dark Theme Mandate

Primary background: #0F1117
Primary surface:    #1A1D27
Primary accent:     #7C6EF5 (electric purple)

NEVER use:
- White or light grey backgrounds
- Duolingo green (#58CC02) as the primary CTA color
- Fredoka font outside companion dialogue

ALWAYS use:
- JetBrains Mono for all numbers (XP, gems, streak, levels, %scores)
- Plus Jakarta Sans for display headings
- #E8EAED for primary text on dark backgrounds

---

# What is and is not Duolingo

Prepio is in the same genre: gamified daily-habit learning.
Prepio is NOT the same product: children's language app.

Structural patterns we share (intentionally):
- Bottom tab navigation
- Companion with speech bubble
- XP, gems, streak mechanics
- Progress rings and bars

Patterns we do NOT share:
- Lime green as primary CTA
- White/pastel backgrounds
- Childish rounded display fonts
- Overly enthusiastic copy ("AMAZING! YOU DID IT!")
- Single-column centered form layouts
```

---

### Update `SCREEN_SPECS.md` — append this section:

```md
# Voice Guidelines

Every string visible to the user should pass this test:
"Would a senior engineer say this to a junior colleague?"

WRONG:  "You've got this! Take your time and think it through."
RIGHT:  "Think about hash maps for this one."

WRONG:  "Amazing! You're on a roll!"
RIGHT:  "Solid." or "Good — time complexity covered."

WRONG:  "Almost there! Keep going!"
RIGHT:  "Not quite. Here's what was missing: [specific bullet points]"

WRONG:  "LET'S GO!" (button copy)
RIGHT:  "Continue Prep" or "Submit" or "Start Session"

WRONG:  "New adventurer? Start your journey"
RIGHT:  "No account? Join 12k engineers in prep"

The companion can be warm and encouraging.
The companion cannot be condescending or childish.
There is a difference between "I believe in you!" and "One challenge from Level 12."
One is generic. One is specific and motivating.

---

# Input Styling

Dark background inputs, not white:
background: #1A1D27
border: 1px solid #2E3347
focus border: #7C6EF5
text: #E8EAED

No white or grey form fields anywhere.

---

# Result Screen Rules

Correct answer: Do NOT say "Amazing!" Say "Solid."
Wrong answer:   Do NOT say "Almost there!" Show what was missed.
Wrong answer:   Do NOT use red. Use orange (#FF6B35) or amber (#F5B942).
XP and gems:    Always animate with count-up. Never appear statically.
Feedback:       At least 2 bullet points. Never a single vague sentence.
```

---

### Update `vision.md` — append this section:

```md
# Audience Differentiation

Duolingo's user: Learning French for vacation, hobby, personal growth.
Prepio's user:   Trying to get their first job at Google or double their salary.

The stakes are categorically different.

Duolingo can be whimsical because nothing rides on a missed lesson.
Prepio users are under real pressure. The UI must honor that.

Playful: Yes — this is a game and games should feel rewarding.
Childish: Never — users are making career decisions that affect their income.

The emotional target for every session:
"I made measurable progress toward the job I want."

Not:
"I earned a fun star today!"

---

# Design Reference Points

Primary inspiration: Duolingo (genre, mechanics, companion)
Secondary inspiration: Linear, Vercel, Raycast (professionalism, dark theme)
Tertiary inspiration: Pokémon (collection, evolution, world progression)

The combination:
- Duolingo's daily habit loop
- Linear's dark professional aesthetic
- Pokémon's emotional attachment to companions

This combination does not currently exist in the market.
It is Prepio's moat.
```

---

## Checklist

**Iteration 1 — Token Swap**
- [ ] `layout.tsx` — fonts swapped to Plus Jakarta Sans + JetBrains Mono
- [ ] `globals.css` — full replacement with dark tokens
- [ ] `tokens.ts` — full replacement, `MIN_ANSWER_LENGTH` = 100
- [ ] `GameBackground.tsx` — dark with purple/teal ambient glow
- [ ] `CompanionHero.tsx` — dark circle, glow ring, reaction prop
- [ ] `SpeechBubble.tsx` — dark surface, accent speaker label
- [ ] `GameButton.tsx` — purple primary, dark secondary, ghost variant
- [ ] `GameCard.tsx` — dark surface, optional left accent border
- [ ] `StatChip.tsx` — dark surface, monospace value
- [ ] `QuestCard.tsx` — dark surface, purple progress bar
- [ ] `ReadinessRing.tsx` — dark track, glow behind ring
- [ ] `BottomNav.tsx` — dark background, 5 tabs, League/Profile greyed

**Iteration 2 — Layout**
- [ ] `HUDBar.tsx` — new component created
- [ ] `dashboard/page.tsx` — HUDBar replaces stat chips, section reorder
- [ ] `login/page.tsx` — split-screen layout
- [ ] `register/page.tsx` — split-screen layout (same pattern as login)
- [ ] `challenge/page.tsx` — two-column md+ layout, dynamic Q count, ResultCard
- [ ] `challenge/page.tsx` — multi-question state with `handleNext`

**Iteration 3 — Voice**
- [ ] `gateway/dashboard/service.go` — `companionMessage` replaced
- [ ] `gateway/dashboard/service.go` — `defaultDailyQuests` titles updated
- [ ] `login/page.tsx` — all copy updated (done in 2.3)
- [ ] `register/page.tsx` — all copy updated
- [ ] `onboarding/page.tsx` — `stepMessages` and button copy updated
- [ ] `challenge/page.tsx` — all copy updated (done in 2.5)
- [ ] `dashboard/page.tsx` — section labels and CTA copy updated

**MD files**
- [ ] `DESIGN_SYSTEM.md` — audience tone + dark mandate appended
- [ ] `SCREEN_SPECS.md` — voice guidelines + input styling appended
- [ ] `vision.md` — audience differentiation + reference points appended