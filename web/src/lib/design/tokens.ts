/** Prepio design tokens — Career RPG for engineers, not children's language app. */

export const MIN_ANSWER_LENGTH = 100;

/** LevelThresholds mirrors config/levels.go cumulative XP per level. */
export const LEVEL_THRESHOLDS = [0, 100, 250, 500, 800, 1200, 1700, 2300, 3000, 3800];

export const colors = {
  bg: "#0F1117",
  surface: "#1A1D27",
  raised: "#242836",
  border: "#2E3347",
  accent: "#7C6EF5",
  accentDim: "rgba(124,110,245,0.15)",
  accentGlow: "rgba(124,110,245,0.35)",
  streak: "#FF6B35",
  gems: "#34D399",
  xp: "#60A5FA",
  gold: "#F5B942",
  success: "#34D399",
  warning: "#F5B942",
  danger: "#F87171",
  textPrimary: "#E8EAED",
  textMuted: "#8B92A8",
  textDim: "#4A5068",
} as const;

export const companyColors: Record<string, { ring: string; bg: string; text: string }> = {
  google: { ring: "#4285F4", bg: "rgba(66,133,244,0.12)", text: "#7EB3FF" },
  amazon: { ring: "#FF9900", bg: "rgba(255,153,0,0.12)", text: "#FFB84D" },
  meta: { ring: "#7C6EF5", bg: "rgba(124,110,245,0.12)", text: "#A99EFA" },
  uber: { ring: "#34D399", bg: "rgba(52,211,153,0.12)", text: "#6EE7B7" },
  atlassian: { ring: "#0052CC", bg: "rgba(0,82,204,0.12)", text: "#5B9BD5" },
  netflix: { ring: "#E50914", bg: "rgba(229,9,20,0.12)", text: "#FF6B6B" },
};

export const leagueThemes: Record<string, { gradient: string; icon: string; border: string }> = {
  bronze: { gradient: "from-amber-900/80 to-amber-700/80", icon: "🥉", border: "#92400E" },
  silver: { gradient: "from-slate-700/80 to-slate-500/80", icon: "🥈", border: "#64748B" },
  gold: { gradient: "from-yellow-700/80 to-amber-500/80", icon: "🥇", border: "#D97706" },
  sapphire: { gradient: "from-blue-800/80 to-blue-600/80", icon: "💎", border: "#2563EB" },
  ruby: { gradient: "from-rose-800/80 to-rose-600/80", icon: "♦️", border: "#E11D48" },
  emerald: { gradient: "from-emerald-800/80 to-emerald-600/80", icon: "💚", border: "#059669" },
  diamond: { gradient: "from-sky-700/80 to-cyan-600/80", icon: "✨", border: "#0891B2" },
  legend: { gradient: "from-yellow-900/80 to-amber-800/80", icon: "👑", border: "#B45309" },
};

export const roundTypeColors: Record<string, string> = {
  dsa: "#7C6EF5",
  system_design: "#60A5FA",
  lld: "#34D399",
  behavioral: "#FF6B35",
  aptitude: "#F5B942",
  fundamentals: "#A99EFA",
};
