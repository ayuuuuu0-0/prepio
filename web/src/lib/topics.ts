/**
 * Presentation helpers for topics and mastery. These only format what the server computed:
 * mastery, readiness, and deltas always come from the Progress service.
 */
import type { MasteryChange } from "@/lib/lesson/types";

export type TopicTheme = { tint: string; icon: string };

const themes: Record<string, TopicTheme> = {
  "system-design": { tint: "#7C6EF5", icon: "🏛️" },
  "backend-production": { tint: "#34D399", icon: "⚙️" },
  "low-level-design": { tint: "#F472B6", icon: "🧩" },
  "dsa-refresher": { tint: "#F5B942", icon: "🌲" },
};

const fallbackTheme: TopicTheme = { tint: "#60A5FA", icon: "◆" };

/** topicTheme returns the accent colour and icon for a topic; unknown topics get a neutral theme. */
export function topicTheme(slug: string): TopicTheme {
  return themes[slug] ?? fallbackTheme;
}

/**
 * ringOffset is the SVG stroke-dashoffset that draws `mastery` (0-100) of a ring.
 * A null mastery ("not started") draws an empty ring.
 */
export function ringOffset(mastery: number | null, circumference: number): number {
  if (mastery === null) return circumference;
  const pct = Math.max(0, Math.min(100, mastery));
  return circumference * (1 - pct / 100);
}

/** coverageText says how much of a topic the learner has started, without any failure framing. */
export function coverageText(started: number, total: number): string {
  if (total <= 0) return "No skills yet";
  if (started === 0) return "Not started yet";
  if (started >= total) return `All ${total} skills started`;
  return `${started} of ${total} skills started`;
}

export type MasteryGroup = {
  topicSlug: string;
  topicName: string;
  skills: { slug: string; name: string; delta: number }[];
};

/**
 * groupMasteryChanges files the server-reported skill deltas under their topic name, e.g.
 * "System Design · Caching +1 · Load Balancing +1". It deliberately produces no topic-level
 * number: topic readiness is computed by the server, never summed here. Unmoved skills are left out.
 */
export function groupMasteryChanges(changes: MasteryChange[]): MasteryGroup[] {
  const groups: MasteryGroup[] = [];
  for (const c of changes) {
    if (c.delta === 0) continue;
    const slug = c.topic_slug || "other";
    let group = groups.find((g) => g.topicSlug === slug);
    if (!group) {
      group = { topicSlug: slug, topicName: c.topic_name || "Other skills", skills: [] };
      groups.push(group);
    }
    group.skills.push({ slug: c.skill_slug, name: c.skill_name, delta: c.delta });
  }
  return groups;
}

/** formatDelta renders a change with an explicit sign: +3, −2, 0. */
export function formatDelta(delta: number): string {
  if (delta > 0) return `+${delta}`;
  if (delta < 0) return `−${Math.abs(delta)}`;
  return "0";
}
