/** Companion visual identity from UI_ASSETS.md */

export type CompanionVisual = {
  emoji: string;
  gradient: string;
  glow: string;
  accent: string;
};

const bySpecies: Record<string, CompanionVisual> = {
  // Classic / starter companions
  capybara: { emoji: "🦫", gradient: "from-amber-600 via-amber-400 to-lime-400", glow: "#7CB342", accent: "#8B5E3C" },
  red_panda: { emoji: "🐼", gradient: "from-orange-500 via-red-400 to-orange-300", glow: "#FF6B35", accent: "#E63946" },
  hedgehog: { emoji: "🦔", gradient: "from-purple-600 via-violet-500 to-amber-400", glow: "#7B5CFF", accent: "#9B59B6" },
  pangolin: { emoji: "🦔", gradient: "from-purple-600 via-violet-500 to-amber-400", glow: "#7B5CFF", accent: "#9B59B6" },
  frog: { emoji: "🐸", gradient: "from-emerald-500 to-lime-400", glow: "#58CC02", accent: "#46A302" },
  axolotl: { emoji: "🐸", gradient: "from-emerald-500 to-lime-400", glow: "#58CC02", accent: "#46A302" },
  orange_cat: { emoji: "🐱", gradient: "from-amber-500 via-orange-400 to-amber-300", glow: "#F59E0B", accent: "#D97706" },
  cat: { emoji: "🐱", gradient: "from-amber-500 via-orange-400 to-amber-300", glow: "#F59E0B", accent: "#D97706" },
  snow_leopard: { emoji: "🐱", gradient: "from-amber-500 via-orange-400 to-amber-300", glow: "#F59E0B", accent: "#D97706" },
  owl: { emoji: "🦉", gradient: "from-emerald-500 to-lime-400", glow: "#58CC02", accent: "#46A302" },

  // Cuter domestic & day-to-day companions
  dog: { emoji: "🐶", gradient: "from-amber-500 via-yellow-400 to-orange-300", glow: "#F59E0B", accent: "#B45309" },
  rabbit: { emoji: "🐰", gradient: "from-pink-300 via-rose-200 to-purple-200", glow: "#EC4899", accent: "#DB2777" },
  hamster: { emoji: "🐹", gradient: "from-orange-400 via-amber-300 to-yellow-200", glow: "#FB923C", accent: "#EA580C" },
  corgi: { emoji: "🐕", gradient: "from-amber-500 via-orange-400 to-yellow-300", glow: "#F59E0B", accent: "#B45309" },
  duck: { emoji: "🐥", gradient: "from-yellow-400 via-amber-300 to-lime-300", glow: "#EAB308", accent: "#CA8A04" },
};

const byName: Record<string, CompanionVisual> = {
  byte: bySpecies.capybara,
  pip: bySpecies.red_panda,
  nova: bySpecies.hedgehog,
  kodo: bySpecies.frog,
  zara: bySpecies.orange_cat,
  prep: bySpecies.owl,

  // Domestic cuties
  milo: bySpecies.orange_cat,
  luna: bySpecies.rabbit,
  buster: bySpecies.dog,
  mochi: bySpecies.hamster,
  pipi: bySpecies.duck,
  peanut: bySpecies.corgi,
  clover: bySpecies.rabbit,
};

/** companionVisual returns art direction for a companion by name or species. */
export function companionVisual(name?: string, species?: string): CompanionVisual {
  if (name && byName[name.toLowerCase()]) {
    return byName[name.toLowerCase()];
  }
  if (species && bySpecies[species.toLowerCase()]) {
    return bySpecies[species.toLowerCase()];
  }
  return bySpecies.orange_cat;
}
