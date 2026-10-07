/** Prepio design tokens — Career RPG for engineers, not children's language app. */

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
