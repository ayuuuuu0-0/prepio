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

/** Skills and lesson themes that scroll past in the chip strip. */
export const chipRows: string[][] = [
  [
    "why caches exist",
    "consistent hashing",
    "rate limiting",
    "idempotent APIs",
    "database indexes",
    "load balancing",
    "message queues",
    "read replicas",
  ],
  [
    "the outbox pattern",
    "isolation levels",
    "backpressure",
    "SOLID, for real",
    "observability basics",
    "two pointers",
    "graph traversal",
    "connection pools",
  ],
];

export const topics = [
  {
    slug: "system-design",
    name: "System Design",
    blurb: "Caching, scaling, queues, and the trade-offs you'll defend in design reviews.",
    demoMastery: 72,
  },
  {
    slug: "backend-production",
    name: "Backend & Production",
    blurb: "APIs, transactions, concurrency, and keeping things alive at 3 a.m.",
    demoMastery: 58,
  },
  {
    slug: "low-level-design",
    name: "Low-Level Design",
    blurb: "Classes, boundaries, and patterns that survive the next feature request.",
    demoMastery: 41,
  },
  {
    slug: "dsa-refresher",
    name: "DSA Refresher",
    blurb: "The data structures and algorithms you actually need, without the grind.",
    demoMastery: 64,
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
