import type { TopicCard } from "@/lib/api";
import { coverageText, topicTheme } from "@/lib/topics";
import { MasteryRing } from "./MasteryRing";

/** TopicMasteryCard shows readiness in one topic: a tinted card with a mastery ring and coverage. */
export function TopicMasteryCard({ topic }: { topic: TopicCard }) {
  const theme = topicTheme(topic.slug);

  return (
    <article
      className="relative overflow-hidden rounded-3xl p-5"
      style={{
        background: `linear-gradient(150deg, ${theme.tint}2B 0%, #1A1D27 58%)`,
        border: `1px solid ${theme.tint}40`,
        boxShadow: "0 8px 32px rgba(0,0,0,0.3)",
      }}
      aria-label={`${topic.name}: ${topic.mastery === null ? "not started" : `mastery ${topic.mastery} out of 100`}`}
    >
      <div
        aria-hidden
        className="pointer-events-none absolute -right-8 -top-8 h-28 w-28 rounded-full blur-2xl"
        style={{ background: theme.tint, opacity: 0.22 }}
      />

      <div className="relative flex items-center gap-4">
        <MasteryRing mastery={topic.mastery} tint={theme.tint} />
        <div className="min-w-0 flex-1">
          <div className="flex items-center gap-2">
            <span aria-hidden className="text-lg">
              {theme.icon}
            </span>
            <h3 className="font-display truncate text-lg font-extrabold leading-tight" style={{ color: "#E8EAED" }}>
              {topic.name}
            </h3>
          </div>
          <p className="mt-1 text-sm font-semibold" style={{ color: "#8B92A8" }}>
            {coverageText(topic.skills_started, topic.skills_total)}
          </p>
          {topic.focused && (
            <span
              className="font-mono mt-2 inline-block rounded-full px-2.5 py-0.5 text-[10px] font-bold uppercase tracking-[0.18em]"
              style={{ background: `${theme.tint}26`, color: theme.tint, border: `1px solid ${theme.tint}55` }}
            >
              Your focus
            </span>
          )}
        </div>
      </div>
    </article>
  );
}
