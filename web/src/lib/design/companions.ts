/** Companion visual identity from UI_ASSETS.md */

export type CompanionVisual = {
  emoji: string;
  gradient: string;
  glow: string;
  accent: string;
};

const bySpecies: Record<string, CompanionVisual> = {
  capybara: { emoji: "🦫", gradient: "from-amber-600 via-amber-400 to-lime-400", glow: "#7CB342", accent: "#8B5E3C" },
  red_panda: { emoji: "🐼", gradient: "from-orange-500 via-red-400 to-orange-300", glow: "#FF6B35", accent: "#E63946" },
  pangolin: { emoji: "🦔", gradient: "from-purple-600 via-violet-500 to-amber-400", glow: "#7B5CFF", accent: "#9B59B6" },
  axolotl: { emoji: "🦎", gradient: "from-pink-400 via-rose-300 to-sky-300", glow: "#FF6B9D", accent: "#FF85C0" },
  snow_leopard: { emoji: "🐆", gradient: "from-slate-300 via-blue-200 to-white", glow: "#1CB0F6", accent: "#64748B" },
  owl: { emoji: "🦉", gradient: "from-emerald-500 to-lime-400", glow: "#58CC02", accent: "#46A302" },
};

const byName: Record<string, CompanionVisual> = {
  byte: bySpecies.capybara,
  pip: bySpecies.red_panda,
  nova: bySpecies.pangolin,
  kodo: bySpecies.axolotl,
  zara: bySpecies.snow_leopard,
  prep: bySpecies.owl,
};

/** companionVisual returns art direction for a companion by name or species. */
export function companionVisual(name?: string, species?: string): CompanionVisual {
  if (name && byName[name.toLowerCase()]) {
    return byName[name.toLowerCase()];
  }
  if (species && bySpecies[species]) {
    return bySpecies[species];
  }
  return bySpecies.capybara;
}
