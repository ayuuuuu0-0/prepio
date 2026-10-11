/**
 * Marketing copy for the landing page. This is not lesson content: the phone screens on the
 * landing page are illustrations of the product, and real lessons live only in content/.
 */

export const hero = {
  announcement: { tag: "New", text: "Weekly leagues are live" },
  headline: "Level up as an engineer.",
  headlineMuted: "Five minutes a day.",
  sub: "Short, instantly graded lessons in system design, backend, low-level design, and DSA. Watch your mastery move, one topic at a time.",
  primaryCta: "Start learning free",
  secondaryCta: "I already have an account",
};

/** Stats that appear below the hero canvas. */
export const stats = [
  { value: "40", label: "lessons" },
  { value: "175", label: "exercises" },
  { value: "4", label: "topics" },
  { value: "5", label: "exercise types" },
] as const;

/**
 * Benefit-first phrases for the scrolling chip marquee.
 * Two rows scroll in opposite directions.
 */
export const chipRows: string[][] = [
  [
    "5-minute lessons",
    "instant feedback",
    "no lives lost",
    "server-graded",
    "system design",
    "weekly leagues",
    "zero guesswork",
    "master backend",
    "earn XP",
    "streak streaks",
  ],
  [
    "DSA refresher",
    "LLD patterns",
    "mistakes come back",
    "mastery rings",
    "beat-based intros",
    "topic mastery",
    "unlock nodes",
    "beat the ladder",
    "idempotent APIs",
    "level up",
  ],
];

/** The four features shown in the sticky scroll showcase section. */
export const stickyFeatures = [
  {
    id: "lesson",
    eyebrow: "The lesson",
    title: "Learn, then",
    titleAccent: "prove it.",
    body: "Every lesson starts with a beat-by-beat animated concept intro. Then exercises — multiple choice, fill-in-the-blank, arrange, true/false — one at a time, in focus mode.",
    log: [
      { step: "intro stage", detail: "concept explained" },
      { step: "exercise 1 / 5", detail: "server graded" },
      { step: "exercise 2 / 5", detail: "combo ×2" },
    ],
    tint: "#7C6EF5",
  },
  {
    id: "grading",
    eyebrow: "Instant feedback",
    title: "Wrong? Find out",
    titleAccent: "why, right now.",
    body: "Every answer is graded on the server in under 100ms. You see the verdict, the reasoning, and the concept reinforced — not just a red flash.",
    log: [
      { step: "answer submitted", detail: "POST /lessons/{id}/submit" },
      { step: "graded server-side", detail: "94ms" },
      { step: "explanation returned", detail: "answer key never sent" },
    ],
    tint: "#34D399",
  },
  {
    id: "mistakes",
    eyebrow: "No penalties",
    title: "Mistakes come",
    titleAccent: "back. That's it.",
    body: "Miss an exercise and it re-queues at the end of the lesson. No hearts, no lives, no penalty streaks. Just another chance to get it right before the lesson closes.",
    log: [
      { step: "wrong answer", detail: "re-queued" },
      { step: "lesson continues", detail: "no XP deducted" },
      { step: "retry at the end", detail: "correct ✓" },
    ],
    tint: "#FF6B35",
  },
  {
    id: "mastery",
    eyebrow: "Real progress",
    title: "Mastery, not",
    titleAccent: "points.",
    body: "Each topic has a 0–100 mastery score backed by the skills in each lesson. The celebration screen tells you exactly which skills moved and why — never just 'good job'.",
    log: [
      { step: "lesson complete", detail: "lesson.completed event" },
      { step: "skill deltas computed", detail: "mastery_ledger written" },
      { step: "System Design", detail: "+3 mastery → 72%" },
    ],
    tint: "#60A5FA",
  },
] as const;

/** The four topics shown in the alternating rows section. */
export const topics = [
  {
    slug: "system-design",
    name: "System Design",
    blurb: "Caching, scaling, queues, and the trade-offs you'll defend in design reviews.",
    detail: "9 nodes · 14 exercises · from caching basics to the boss challenge",
    demoMastery: 72,
    companion: "🦉",
  },
  {
    slug: "backend-production",
    name: "Backend & Production",
    blurb: "APIs, transactions, concurrency, and keeping things alive at 3 a.m.",
    detail: "8 nodes · idempotency, rate limits, retries, observability",
    demoMastery: 58,
    companion: "🦫",
  },
  {
    slug: "low-level-design",
    name: "Low-Level Design",
    blurb: "Classes, boundaries, and patterns that survive the next feature request.",
    detail: "8 nodes · SOLID, Factory, Strategy, Observer, Decorator",
    demoMastery: 41,
    companion: "🦝",
  },
  {
    slug: "dsa-refresher",
    name: "DSA Refresher",
    blurb: "The data structures and algorithms you actually need, without the grind.",
    detail: "15 nodes · two pointers to dynamic programming — 14 skills covered",
    demoMastery: 64,
    companion: "🦎",
  },
] as const;

export const steps = [
  { title: "Pick what to sharpen", body: "Choose one to three topics. They shape your path; nothing gets locked away." },
  { title: "Play a five-minute lesson", body: "A short intro, then a few exercises with instant feedback. Misses come back at the end." },
  { title: "Watch your mastery move", body: "Every lesson explains exactly which skills grew, and the next node unlocks." },
];

export const faqs = [
  {
    q: "Is this another LeetCode?",
    a: "No. Prepio is about understanding: short lessons that teach and measure real engineering skills, not a bank of problems to grind.",
  },
  {
    q: "Who is it for?",
    a: "Working software engineers who want to get sharper at system design, backend and production work, low-level design, and DSA, for the job or for interviews.",
  },
  {
    q: "How long does a lesson take?",
    a: "Under five minutes. One lesson a day keeps your streak alive.",
  },
  {
    q: "How do leagues work?",
    a: "Your first lesson each week puts you on a leaderboard with up to 30 learners. Earn XP from lessons; the top ranks move up a league when the week ends.",
  },
  {
    q: "Are there hearts or lives?",
    a: "Never. Mistakes come back at the end of the lesson so you can get them right. No penalties, no guilt.",
  },
  {
    q: "Does it cost anything?",
    a: "Prepio is free while in beta.",
  },
];

export const finalCta = {
  headline: "Your next level",
  headlineMuted: "is five minutes away.",
  cta: "Start learning free",
};
