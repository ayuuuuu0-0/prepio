/** Prepio game design tokens — Duolingo-inspired, never corporate. */

/** MinAnswerLength is the minimum trimmed characters required to submit an answer. */
export const MIN_ANSWER_LENGTH = 3;

export const colors = {
  green: "#58CC02",
  greenDark: "#46A302",
  purple: "#7B5CFF",
  purpleDark: "#5B3FE0",
  orange: "#FF9600",
  orangeSoft: "#FFB84D",
  blue: "#1CB0F6",
  gold: "#FFC800",
  pink: "#FF6B9D",
  bgTop: "#E8F5D8",
  bgBottom: "#C8E6FF",
  surface: "#FFFFFF",
  text: "#3C3C3C",
  textMuted: "#777777",
} as const;

export const companyColors: Record<string, { ring: string; bg: string }> = {
  google: { ring: "#1CB0F6", bg: "#E3F6FF" },
  amazon: { ring: "#FF9600", bg: "#FFF3E0" },
  meta: { ring: "#7B5CFF", bg: "#F0EBFF" },
  uber: { ring: "#58CC02", bg: "#E8F8E0" },
  atlassian: { ring: "#0052CC", bg: "#E8F0FF" },
  netflix: { ring: "#E50914", bg: "#FFE8E8" },
};

export const leagueThemes: Record<string, { gradient: string; icon: string }> = {
  bronze: { gradient: "from-amber-700 to-amber-500", icon: "🥉" },
  silver: { gradient: "from-slate-400 to-slate-300", icon: "🥈" },
  gold: { gradient: "from-yellow-500 to-amber-400", icon: "🥇" },
  sapphire: { gradient: "from-blue-500 to-cyan-400", icon: "💎" },
  ruby: { gradient: "from-red-500 to-rose-400", icon: "♦️" },
  emerald: { gradient: "from-green-500 to-emerald-400", icon: "💚" },
  diamond: { gradient: "from-sky-300 to-blue-200", icon: "✨" },
  legend: { gradient: "from-yellow-600 to-amber-900", icon: "👑" },
};

export const radii = {
  card: "1.5rem",
  pill: "9999px",
  bubble: "1.25rem",
} as const;
