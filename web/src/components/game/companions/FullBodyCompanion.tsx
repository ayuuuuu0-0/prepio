"use client";

import { ByteCompanion } from "./ByteCompanion";
import { PipCompanion } from "./PipCompanion";
import { ZaraCompanion } from "./ZaraCompanion";
import { KodoCompanion } from "./KodoCompanion";
import { NovaCompanion } from "./NovaCompanion";

type Reaction = "idle" | "correct" | "wrong";

export function FullBodyCompanion({
  name,
  species,
  reaction = "idle",
  className = "",
}: {
  name?: string;
  species?: string;
  reaction?: Reaction;
  className?: string;
}) {
  const n = (name ?? "").toLowerCase();
  const s = (species ?? "").toLowerCase();

  if (n.includes("pip") || s.includes("red_panda") || s.includes("panda")) {
    return <PipCompanion reaction={reaction} className={className} />;
  }
  if (n.includes("zara") || s.includes("cat") || s.includes("orange_cat") || s.includes("snow_leopard")) {
    return <ZaraCompanion reaction={reaction} className={className} />;
  }
  if (n.includes("kodo") || s.includes("frog") || n.includes("prep") || s.includes("owl") || s.includes("axolotl")) {
    return <KodoCompanion reaction={reaction} className={className} />;
  }
  if (n.includes("nova") || s.includes("hedgehog") || s.includes("pangolin")) {
    return <NovaCompanion reaction={reaction} className={className} />;
  }

  // Default to Byte the Capybara in purple hoodie
  return <ByteCompanion reaction={reaction} className={className} />;
}
