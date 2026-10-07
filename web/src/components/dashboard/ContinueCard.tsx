import Link from "next/link";
import type { NextLesson } from "@/lib/api";

/** ContinueCard is the dashboard's one clear next step: the lesson the path says is next. */
export function ContinueCard({ lesson, companionName }: { lesson: NextLesson | null; companionName?: string }) {
  if (!lesson) {
    return (
      <section
        className="rounded-3xl p-6"
        style={{ background: "#1A1D27", border: "1px solid #2E3347" }}
        aria-labelledby="caught-up"
      >
        <p className="font-mono text-[11px] font-bold uppercase tracking-[0.2em]" style={{ color: "#34D399" }}>
          All caught up
        </p>
        <h2 id="caught-up" className="font-display mt-2 text-2xl font-extrabold leading-tight" style={{ color: "#E8EAED" }}>
          You&apos;ve finished every lesson available right now.
        </h2>
        <p className="mt-2 text-sm font-semibold" style={{ color: "#8B92A8" }}>
          {companionName ? `${companionName} will let you know` : "We'll let you know"} when new lessons arrive. Revisit
          the journey to practice any lesson again.
        </p>
        <Link
          href="/journey"
          className="font-display mt-5 inline-block rounded-full px-6 py-3 text-sm font-bold focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-[#7C6EF5]"
          style={{ background: "#242836", color: "#E8EAED", border: "1px solid #2E3347" }}
        >
          Open the journey
        </Link>
      </section>
    );
  }

  return (
    <section
      className="relative overflow-hidden rounded-3xl p-6"
      style={{
        background: "linear-gradient(135deg, #6354E8 0%, #8B7FF7 55%, #A99EFA 100%)",
        boxShadow: "0 16px 48px rgba(124,110,245,0.35)",
      }}
      aria-labelledby="up-next"
    >
      <div aria-hidden className="pointer-events-none absolute -right-10 -top-10 h-40 w-40 rounded-full bg-white/15 blur-2xl" />
      <div aria-hidden className="pointer-events-none absolute -bottom-12 left-1/3 h-32 w-32 rounded-full bg-black/15 blur-2xl" />

      <p className="font-mono relative text-[11px] font-bold uppercase tracking-[0.2em] text-white/80">
        {lesson.in_progress ? "Pick up where you left off" : "Up next"} · {lesson.world_name}
      </p>
      <h2 id="up-next" className="font-display relative mt-2 text-3xl font-extrabold leading-tight tracking-tight text-white">
        {lesson.title}
      </h2>

      <div className="relative mt-4 flex flex-wrap gap-2">
        {[`${lesson.est_minutes} min`, `Up to ${lesson.xp_preview} XP`, lesson.kind === "boss" ? "Boss lesson" : null]
          .filter(Boolean)
          .map((chip) => (
            <span
              key={chip}
              className="font-mono rounded-full bg-white/20 px-3 py-1 text-xs font-bold text-white backdrop-blur"
            >
              {chip}
            </span>
          ))}
      </div>

      <Link
        href={`/lesson/${lesson.lesson_id}`}
        className="font-display relative mt-6 flex w-full items-center justify-center gap-2 rounded-full bg-white px-8 py-4 text-base font-extrabold tracking-wide text-[#3D2FC4] shadow-lg transition-transform hover:scale-[1.01] active:scale-[0.99] focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white motion-reduce:transition-none"
      >
        {lesson.in_progress ? "Continue lesson" : "Start lesson"} <span aria-hidden>→</span>
      </Link>
    </section>
  );
}
