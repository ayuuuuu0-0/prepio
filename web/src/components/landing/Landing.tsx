import Image from "next/image";
import Link from "next/link";
import type { CSSProperties, ReactNode } from "react";
import { MasteryRing } from "@/components/dashboard/MasteryRing";
import { TierEmblem } from "@/components/league/TierEmblem";
import { chipRows, faqs, finalCta, hero, steps, topics } from "@/lib/landing/content";
import { companionVisual } from "@/lib/design/companions";
import { topicTheme } from "@/lib/topics";
import { CelebrationScreen, JourneyScreen, LessonScreen, Phone } from "./PhoneScreens";
import s from "./landing.module.css";

/** Landing is the public marketing page: one promise, real product screens, one call to action. */
export function Landing() {
  return (
    <div className="min-h-dvh overflow-x-hidden" style={{ background: "#0F1117", color: "#E8EAED" }}>
      <a
        href="#main"
        className="sr-only focus:not-sr-only focus:fixed focus:left-4 focus:top-4 focus:z-[60] focus:rounded-full focus:bg-white focus:px-4 focus:py-2 focus:text-black"
      >
        Skip to content
      </a>
      <Nav />
      <main id="main">
        <Hero />
        <Chips />
        <Topics />
        <Features />
        <HowItWorks />
        <Faq />
        <FinalCta />
      </main>
      <Footer />
    </div>
  );
}

function Nav() {
  const links = [
    { href: "#topics", label: "Topics" },
    { href: "#features", label: "Features" },
    { href: "#how", label: "How it works" },
    { href: "#faq", label: "FAQ" },
  ];
  return (
    <header className="fixed inset-x-0 top-0 z-50 px-4 pt-4" style={{ paddingTop: "max(1rem, env(safe-area-inset-top))" }}>
      <nav
        aria-label="Primary"
        className="mx-auto flex max-w-5xl items-center justify-between gap-3 rounded-full py-2 pl-4 pr-2"
        style={{ background: "rgba(26,29,39,0.72)", border: "1px solid rgba(255,255,255,0.08)", backdropFilter: "blur(18px)" }}
      >
        <Link href="/" className="flex items-center gap-2">
          <Image src="/logo.png" alt="" width={26} height={26} className="rounded-lg" />
          <span className="font-display text-sm font-extrabold tracking-wide" style={{ color: "#E8EAED" }}>
            Prepio
          </span>
        </Link>
        <ul className="hidden items-center gap-1 md:flex">
          {links.map((l) => (
            <li key={l.href}>
              <a href={l.href} className="rounded-full px-3 py-1.5 text-sm font-semibold transition-colors hover:bg-white/5" style={{ color: "#B9C0D0" }}>
                {l.label}
              </a>
            </li>
          ))}
        </ul>
        <div className="flex items-center gap-1">
          <Link href="/login" className="hidden rounded-full px-3 py-1.5 text-sm font-semibold sm:block" style={{ color: "#B9C0D0" }}>
            Log in
          </Link>
          <CtaLink href="/register" small>
            Get started
          </CtaLink>
        </div>
      </nav>
    </header>
  );
}

function CtaLink({ href, children, small = false }: { href: string; children: ReactNode; small?: boolean }) {
  return (
    <Link
      href={href}
      className={`game-btn font-display inline-flex items-center justify-center rounded-full font-extrabold text-white ${
        small ? "px-4 py-2 text-sm" : "px-7 py-3.5 text-base"
      }`}
      style={{ background: "linear-gradient(135deg, #7C6EF5, #9D8FF7)" }}
    >
      {children}
    </Link>
  );
}

type Sticker = { emoji: string; label: string; className: string; tilt: string; delay: string; bg: string };

const heroStickers: Sticker[] = [
  { emoji: "🔥", label: "streak", className: "left-[5%] top-[200px]", tilt: "-10deg", delay: "0s", bg: "#FF6B35" },
  { emoji: "💎", label: "gems", className: "right-[6%] top-[170px]", tilt: "8deg", delay: "1.2s", bg: "#34D399" },
  { emoji: companionVisual("pip").emoji, label: "companion", className: "left-[7%] top-[440px]", tilt: "6deg", delay: "0.6s", bg: "#FF6B35" },
  { emoji: "🏆", label: "league trophy", className: "right-[8%] top-[470px]", tilt: "-7deg", delay: "1.8s", bg: "#F5B942" },
];

function StickerTile({ s: st, size = 76 }: { s: Sticker; size?: number }) {
  return (
    <span
      aria-hidden
      className={`absolute hidden items-center justify-center rounded-[22px] lg:flex ${s.sticker} ${st.className}`}
      style={
        {
          width: size,
          height: size,
          fontSize: size * 0.5,
          "--tilt": st.tilt,
          "--delay": st.delay,
          background: `linear-gradient(145deg, ${st.bg}55, #1A1D27 75%)`,
          border: "3px solid rgba(255,255,255,0.9)",
          boxShadow: `0 14px 30px rgba(0,0,0,0.45), 0 0 30px ${st.bg}44`,
        } as CSSProperties
      }
    >
      {st.emoji}
    </span>
  );
}

function Hero() {
  return (
    <section className={`relative px-4 pb-16 pt-32 sm:pt-40 ${s.sky}`}>
      {heroStickers.map((st) => (
        <StickerTile key={st.label} s={st} />
      ))}
      <div className="relative mx-auto flex max-w-5xl flex-col items-center text-center">
        <a
          href="#features"
          className="inline-flex items-center gap-2 rounded-full py-1 pl-1 pr-3 text-sm font-semibold"
          style={{ background: "rgba(255,255,255,0.06)", border: "1px solid rgba(255,255,255,0.12)", color: "#E8EAED" }}
        >
          <span className="rounded-full px-2 py-0.5 text-xs font-bold" style={{ background: "#7C6EF5", color: "#fff" }}>
            {hero.announcement.tag}
          </span>
          {hero.announcement.text}
          <span aria-hidden>→</span>
        </a>
        <h1 className={`font-display mt-6 max-w-5xl text-[2.6rem] font-extrabold sm:text-6xl md:text-7xl lg:text-[5rem] ${s.headline}`}>
          <span className="block">{hero.headline}</span>
          <span className="block" style={{ color: "#6F7690" }}>
            {hero.headlineMuted}
          </span>
        </h1>
        <p className="mt-6 max-w-xl text-base leading-relaxed sm:text-lg" style={{ color: "#B9C0D0" }}>
          {hero.sub}
        </p>
        <div className="mt-8 flex flex-col items-center gap-3 sm:flex-row sm:gap-5">
          <CtaLink href="/register">{hero.primaryCta}</CtaLink>
          <Link href="/login" className="text-sm font-semibold underline-offset-4 hover:underline" style={{ color: "#E8EAED" }}>
            {hero.secondaryCta} <span aria-hidden>→</span>
          </Link>
        </div>
      </div>

      <div className="relative mx-auto mt-16 max-w-5xl">
        <div
          className="overflow-hidden rounded-[1.75rem]"
          style={{ background: "#16181F", border: "1px solid #2A2E3D", boxShadow: "0 40px 120px rgba(124,110,245,0.22)" }}
        >
          <div className="flex items-center gap-2 px-5 py-3.5" style={{ borderBottom: "1px solid #232736" }}>
            <span className="h-3 w-3 rounded-full" style={{ background: "#FF5F57" }} />
            <span className="h-3 w-3 rounded-full" style={{ background: "#FEBC2E" }} />
            <span className="h-3 w-3 rounded-full" style={{ background: "#28C840" }} />
            <span className="font-mono ml-3 text-[11px]" style={{ color: "#6F7690" }}>
              prepio · one lesson, start to finish
            </span>
          </div>
          <div
            className="flex snap-x snap-mandatory gap-10 overflow-x-auto px-6 py-10 md:justify-center md:gap-14"
            style={{ backgroundImage: "radial-gradient(#232736 1px, transparent 1px)", backgroundSize: "18px 18px" }}
          >
            <div className="snap-center">
              <Phone label="The journey" status="Next node unlocked">
                <JourneyScreen />
              </Phone>
            </div>
            <div className="snap-center">
              <Phone label="A lesson" status="Graded instantly">
                <LessonScreen />
              </Phone>
            </div>
            <div className="snap-center">
              <Phone label="The payoff" status="Mastery moved">
                <CelebrationScreen />
              </Phone>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

function Chips() {
  return (
    <section aria-label="What you'll practice" className="py-14">
      <p className="font-mono text-center text-[11px] font-bold uppercase tracking-[0.22em]" style={{ color: "#6F7690" }}>
        Skills engineers sharpen on Prepio
      </p>
      <div className={`mt-6 flex flex-col gap-3 overflow-hidden ${s.marquee}`}>
        {chipRows.map((row, i) => (
          <div key={i} className={i % 2 === 0 ? s.track : s.trackReverse}>
            {[...row, ...row].map((chip, j) => (
              <span
                key={`${chip}-${j}`}
                aria-hidden={j >= row.length}
                className="mx-1.5 flex items-center gap-2 whitespace-nowrap rounded-full px-4 py-2.5 text-sm font-semibold"
                style={{ background: "#1A1D27", border: "1px solid #2E3347", color: "#D5DAE6" }}
              >
                <span aria-hidden style={{ color: "#7C6EF5" }}>
                  ◆
                </span>
                {chip}
              </span>
            ))}
          </div>
        ))}
      </div>
    </section>
  );
}

function SectionHeading({ id, title, muted }: { id: string; title: string; muted: string }) {
  return (
    <h2 id={id} className={`font-display mx-auto max-w-3xl text-center text-4xl font-extrabold sm:text-5xl md:text-6xl ${s.headline}`}>
      {title}
      <span className="block" style={{ color: "#6F7690" }}>
        {muted}
      </span>
    </h2>
  );
}

function Topics() {
  return (
    <section id="topics" aria-labelledby="topics-title" className="scroll-mt-24 px-4 py-20">
      <SectionHeading id="topics-title" title="Four topics." muted="One mastery score each." />
      <p className="mx-auto mt-5 max-w-xl text-center" style={{ color: "#B9C0D0" }}>
        Mastery is per topic, from 0 to 100, and it always tells you why it moved. Never from XP, never from streaks.
      </p>
      <div className="mx-auto mt-12 grid max-w-5xl gap-4 sm:grid-cols-2">
        {topics.map((t) => {
          const theme = topicTheme(t.slug);
          return (
            <article
              key={t.slug}
              className={`flex items-center gap-5 rounded-[1.75rem] p-6 ${s.card}`}
              style={{
                background: `linear-gradient(140deg, ${theme.tint}26 0%, #1A1D27 55%)`,
                border: `1px solid ${theme.tint}33`,
              }}
            >
              <div className="min-w-0 flex-1">
                <p className="text-2xl" aria-hidden>
                  {theme.icon}
                </p>
                <h3 className="font-display mt-2 text-xl font-extrabold tracking-tight">{t.name}</h3>
                <p className="mt-1.5 text-sm leading-relaxed" style={{ color: "#B9C0D0" }}>
                  {t.blurb}
                </p>
              </div>
              <MasteryRing mastery={t.demoMastery} tint={theme.tint} size={92} />
            </article>
          );
        })}
      </div>
    </section>
  );
}

function FeatureCard({
  eyebrow,
  title,
  muted,
  body,
  tint,
  wide = false,
  children,
}: {
  eyebrow: string;
  title: string;
  muted: string;
  body: string;
  tint: string;
  wide?: boolean;
  children: ReactNode;
}) {
  return (
    <article
      className={`flex flex-col overflow-hidden rounded-[1.75rem] ${s.card} ${wide ? "md:col-span-2 md:flex-row md:items-center" : ""}`}
      style={{ background: `linear-gradient(160deg, #1A1D27 40%, ${tint}22 100%)`, border: "1px solid #2E3347" }}
    >
      <div className={`flex items-center justify-center px-6 pt-8 ${wide ? "md:order-2 md:w-1/2 md:py-8" : ""}`}>{children}</div>
      <div className={`p-7 ${wide ? "md:w-1/2" : ""}`}>
        <p className="font-mono text-[11px] font-bold uppercase tracking-[0.2em]" style={{ color: tint }}>
          {eyebrow}
        </p>
        <h3 className={`font-display mt-3 text-3xl font-extrabold sm:text-4xl ${s.headline}`}>
          {title}
          <span className="block" style={{ color: "#6F7690" }}>
            {muted}
          </span>
        </h3>
        <p className="mt-4 text-sm leading-relaxed" style={{ color: "#B9C0D0" }}>
          {body}
        </p>
      </div>
    </article>
  );
}

const leagueRows = [
  { name: "mira.dev", species: "axolotl", xp: 340 },
  { name: "you", species: "red_panda", xp: 285, me: true },
  { name: "kenji_b", species: "owl", xp: 260 },
  { name: "sam.ops", species: "capybara", xp: 190 },
];

function Features() {
  return (
    <section id="features" aria-labelledby="features-title" className="scroll-mt-24 px-4 py-20">
      <SectionHeading id="features-title" title="Built like a game." muted="Graded like an engineer." />
      <div className="mx-auto mt-12 grid max-w-5xl gap-4 md:grid-cols-2">
        <FeatureCard
          wide
          eyebrow="Instant feedback"
          title="Every answer graded."
          muted="Every miss explained."
          body="Pick, check, and see why right away. Get one wrong and you'll see the right answer and the reasoning, then it comes back at the end of the lesson. No hearts. No lives. No penalties."
          tint="#34D399"
        >
          <div className="w-full max-w-sm rounded-2xl p-4" style={{ background: "#10231C", border: "1px solid #34D39955" }}>
            <p className="font-display flex items-center gap-2 font-extrabold" style={{ color: "#34D399" }}>
              <span className="flex h-6 w-6 items-center justify-center rounded-full text-xs" style={{ background: "#34D399", color: "#0F1117" }}>
                ✓
              </span>
              Nice, that&apos;s right
            </p>
            <p className="mt-2 text-sm" style={{ color: "#B9C0D0" }}>
              Idempotency keys let the client retry safely: the server recognises the repeat and returns the first result.
            </p>
          </div>
        </FeatureCard>

        <FeatureCard
          eyebrow="Leagues"
          title="Climb the ladder."
          muted="Every week."
          body="Your first lesson each week puts you on a leaderboard with up to 30 learners. The top ranks move up a league on Monday."
          tint="#F5B942"
        >
          <div className="w-full max-w-xs overflow-hidden rounded-2xl" style={{ background: "#151821", border: "1px solid #2E3347" }}>
            <div className="flex items-center gap-2 px-3 py-2.5" style={{ borderBottom: "1px solid #232736" }}>
              <TierEmblem tier={{ index: 2, slug: "gold", name: "Gold" }} size={28} />
              <span className="font-display text-sm font-extrabold">Gold League</span>
              <span className="font-mono ml-auto text-[10px]" style={{ color: "#8B92A8" }}>
                2d 6h left
              </span>
            </div>
            {leagueRows.map((r, i) => (
              <div
                key={r.name}
                className="flex items-center gap-2.5 px-3 py-2 text-sm"
                style={{ background: r.me ? "rgba(124,110,245,0.14)" : "transparent" }}
              >
                <span className="font-display w-4 font-extrabold" style={{ color: "#34D399" }}>
                  {i + 1}
                </span>
                <span aria-hidden>{companionVisual(undefined, r.species).emoji}</span>
                <span className="flex-1 font-semibold">{r.name}</span>
                <span className="font-mono text-xs font-bold" style={{ color: "#60A5FA" }}>
                  {r.xp} XP
                </span>
              </div>
            ))}
          </div>
        </FeatureCard>

        <FeatureCard
          eyebrow="Streaks"
          title="Show up."
          muted="That's the whole rule."
          body="One lesson keeps your streak alive. Life happens, so a streak freeze covers a missed day."
          tint="#FF6B35"
        >
          <div className="flex items-end gap-2" aria-hidden>
            {["M", "T", "W", "T", "F", "S", "S"].map((d, i) => (
              <div key={i} className="flex flex-col items-center gap-1.5">
                <span
                  className="flex h-10 w-10 items-center justify-center rounded-xl text-lg"
                  style={{
                    background: i < 5 ? "linear-gradient(160deg, #FF6B35, #F5B942)" : "#242836",
                    boxShadow: i < 5 ? "0 3px 0 #C2491F" : "none",
                    opacity: i < 5 ? 1 : 0.6,
                  }}
                >
                  {i < 5 ? "🔥" : ""}
                </span>
                <span className="font-mono text-[10px]" style={{ color: "#8B92A8" }}>
                  {d}
                </span>
              </div>
            ))}
          </div>
        </FeatureCard>

        <FeatureCard
          wide
          eyebrow="Companions"
          title="Someone in your corner."
          muted="Who notices."
          body="Pick a companion at the start. It cheers your wins, shrugs off your misses, and grows with you."
          tint="#A78BFA"
        >
          <div className="flex items-center gap-3 text-4xl" aria-hidden>
            {["capybara", "red_panda", "pangolin", "axolotl", "owl"].map((sp) => {
              const v = companionVisual(undefined, sp);
              return (
                <span
                  key={sp}
                  className="flex h-14 w-14 items-center justify-center rounded-full text-3xl"
                  style={{ background: `linear-gradient(135deg, ${v.glow}40, #1A1D27)`, border: "1px solid #2E3347", boxShadow: `0 0 20px ${v.glow}33` }}
                >
                  {v.emoji}
                </span>
              );
            })}
          </div>
        </FeatureCard>
      </div>
    </section>
  );
}

function HowItWorks() {
  return (
    <section id="how" aria-labelledby="how-title" className="scroll-mt-24 px-4 py-20" style={{ background: "#13151C" }}>
      <SectionHeading id="how-title" title="Pick. Play." muted="Level up." />
      <ol className="mx-auto mt-12 grid max-w-5xl gap-4 md:grid-cols-3">
        {steps.map((st, i) => (
          <li key={st.title} className="rounded-[1.75rem] p-7" style={{ background: "#1A1D27", border: "1px solid #2E3347" }}>
            <span
              className="font-display flex h-10 w-10 items-center justify-center rounded-full text-lg font-extrabold text-white"
              style={{ background: "#7C6EF5", boxShadow: "0 3px 0 #5B50D4" }}
            >
              {i + 1}
            </span>
            <h3 className="font-display mt-5 text-xl font-extrabold tracking-tight">{st.title}</h3>
            <p className="mt-2 text-sm leading-relaxed" style={{ color: "#B9C0D0" }}>
              {st.body}
            </p>
          </li>
        ))}
      </ol>
    </section>
  );
}

function Faq() {
  return (
    <section id="faq" aria-labelledby="faq-title" className="scroll-mt-24 px-4 py-20">
      <SectionHeading id="faq-title" title="Questions?" muted="Answers." />
      <div className="mx-auto mt-10 flex max-w-2xl flex-col gap-3">
        {faqs.map((f) => (
          <details key={f.q} className="group rounded-2xl px-5 py-4" style={{ background: "#1A1D27", border: "1px solid #2E3347" }}>
            <summary className="font-display flex cursor-pointer list-none items-center justify-between gap-4 font-bold">
              {f.q}
              <span aria-hidden className="text-xl transition-transform group-open:rotate-45" style={{ color: "#7C6EF5" }}>
                +
              </span>
            </summary>
            <p className="mt-3 text-sm leading-relaxed" style={{ color: "#B9C0D0" }}>
              {f.a}
            </p>
          </details>
        ))}
      </div>
    </section>
  );
}

function FinalCta() {
  return (
    <section className={`relative overflow-hidden px-4 py-28 text-center ${s.sky}`}>
      <StickerTile s={{ ...heroStickers[0], className: "left-[10%] top-[20%]" }} size={64} />
      <StickerTile s={{ ...heroStickers[3], className: "right-[10%] bottom-[18%]" }} size={64} />
      <h2 className={`font-display mx-auto max-w-3xl text-4xl font-extrabold sm:text-6xl ${s.headline}`}>
        {finalCta.headline}
        <span className="block" style={{ color: "#8B92A8" }}>
          {finalCta.headlineMuted}
        </span>
      </h2>
      <div className="mt-10">
        <CtaLink href="/register">{finalCta.cta}</CtaLink>
      </div>
    </section>
  );
}

function Footer() {
  return (
    <footer className="px-4 py-10" style={{ borderTop: "1px solid #1F2230", paddingBottom: "max(2.5rem, env(safe-area-inset-bottom))" }}>
      <div className="mx-auto flex max-w-5xl flex-col items-center justify-between gap-4 text-sm sm:flex-row" style={{ color: "#6F7690" }}>
        <span className="flex items-center gap-2">
          <Image src="/logo.png" alt="" width={20} height={20} className="rounded-md" />
          Prepio, a progression game for working engineers.
        </span>
        <span className="flex gap-5">
          <Link href="/login" className="hover:text-white">
            Log in
          </Link>
          <Link href="/register" className="hover:text-white">
            Sign up
          </Link>
        </span>
      </div>
    </footer>
  );
}
