"use client";

import Image from "next/image";
import Link from "next/link";
import { useEffect, useRef, useState, type CSSProperties, type ReactNode } from "react";
import { MasteryRing } from "@/components/dashboard/MasteryRing";
import { TierEmblem } from "@/components/league/TierEmblem";
import {
  chipRows,
  faqs,
  finalCta,
  hero,
  stickyFeatures,
  steps,
  topics,
} from "@/lib/landing/content";
import { companionVisual } from "@/lib/design/companions";
import { topicTheme } from "@/lib/topics";
import {
  CelebrationScreen,
  GradingPhoneScreen,
  IntroPhoneScreen,
  JourneyScreen,
  LessonScreen,
  MasteryPhoneScreen,
  MistakePhoneScreen,
  Phone,
  PhoneFrame,
} from "./PhoneScreens";
import s from "./landing.module.css";

/** Scroll reveal hook using IntersectionObserver for Zeptap-style animations */
function useScrollReveal() {
  const ref = useRef<HTMLDivElement>(null);
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    const el = ref.current;
    if (!el) return;

    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setIsVisible(true);
          observer.unobserve(entry.target);
        }
      },
      { threshold: 0.12, rootMargin: "0px 0px -40px 0px" }
    );

    observer.observe(el);
    return () => observer.disconnect();
  }, []);

  return { ref, isVisible };
}

/** Animated Reveal Component that triggers when the user scrolls to it */
function RevealOnScroll({
  children,
  className = "",
  delay = 0,
}: {
  children: ReactNode;
  className?: string;
  delay?: number;
}) {
  const { ref, isVisible } = useScrollReveal();

  return (
    <div
      ref={ref}
      className={`${s.revealItem} ${isVisible ? s.revealed : ""} ${className}`}
      style={{ transitionDelay: `${delay}ms` }}
    >
      {children}
    </div>
  );
}

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
        <StickyShowcase />
        <AlternatingTopics />
        <Features />
        <CompanionPillars />
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
    { href: "#experience", label: "Experience" },
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
      className={`game-btn font-display inline-flex items-center justify-center rounded-full font-extrabold text-white transition-transform active:scale-95 ${
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
  { emoji: companionVisual("milo").emoji, label: "companion", className: "left-[7%] top-[440px]", tilt: "6deg", delay: "0.6s", bg: "#F59E0B" },
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
    <section className={`relative px-4 pb-20 pt-32 sm:pt-40 ${s.sky}`}>
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

      {/* Mac Window Phone Canvas */}
      <div className="relative mx-auto mt-20 max-w-5xl">
        <RevealOnScroll>
          <div
            className="overflow-hidden rounded-[1.75rem]"
            style={{ background: "#16181F", border: "1px solid #2A2E3D", boxShadow: "0 40px 120px rgba(124,110,245,0.22)" }}
          >
            <div className="flex items-center gap-2 px-5 py-3.5" style={{ borderBottom: "1px solid #232736" }}>
              <span className="h-3 w-3 rounded-full" style={{ background: "#FF5F57" }} />
              <span className="h-3 w-3 rounded-full" style={{ background: "#FEBC2E" }} />
              <span className="h-3 w-3 rounded-full" style={{ background: "#28C840" }} />
              <span className="font-mono ml-3 text-[11px]" style={{ color: "#6F7690" }}>
                prepio · live phone canvas
              </span>
            </div>
            <div
              className="flex snap-x snap-mandatory gap-10 overflow-x-auto px-6 py-12 md:justify-center md:gap-14"
              style={{ backgroundImage: "radial-gradient(#232736 1px, transparent 1px)", backgroundSize: "18px 18px" }}
            >
              <div className="snap-center">
                <Phone label="The Journey" status="Next node unlocked">
                  <JourneyScreen />
                </Phone>
              </div>
              <div className="snap-center">
                <Phone label="The Lesson" status="Server graded">
                  <LessonScreen />
                </Phone>
              </div>
              <div className="snap-center">
                <Phone label="The Payoff" status="Mastery moved">
                  <CelebrationScreen />
                </Phone>
              </div>
            </div>
          </div>
        </RevealOnScroll>
      </div>
    </section>
  );
}

function Chips() {
  return (
    <section aria-label="What you'll practice" className="py-20">
      <RevealOnScroll>
        <p className="font-mono text-center text-[11px] font-bold uppercase tracking-[0.22em]" style={{ color: "#6F7690" }}>
          Proven engineering skills · Five minutes at a time
        </p>
        <div className={`mt-8 flex flex-col gap-4 overflow-hidden ${s.marquee}`}>
          {chipRows.map((row, i) => (
            <div key={i} className={i % 2 === 0 ? s.track : s.trackReverse}>
              {[...row, ...row].map((chip, j) => (
                <span
                  key={`${chip}-${j}`}
                  aria-hidden={j >= row.length}
                  className="mx-2 flex items-center gap-2 whitespace-nowrap rounded-full px-5 py-2.5 text-sm font-semibold"
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
      </RevealOnScroll>
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

/** Sticky interactive showcase (Zeptap-style): Tabbed features switching live phone preview */
function StickyShowcase() {
  const [activeTab, setActiveTab] = useState<string>(stickyFeatures[0].id);

  return (
    <section id="experience" aria-labelledby="experience-title" className="scroll-mt-24 px-4 py-28">
      <RevealOnScroll>
        <SectionHeading id="experience-title" title="The lesson loop." muted="Step by step." />
        <p className="mx-auto mt-4 max-w-xl text-center text-sm sm:text-base leading-relaxed" style={{ color: "#B9C0D0" }}>
          Inspired by the fastest feedback loops in software. No passive video watching.
        </p>
      </RevealOnScroll>

      <div className="mx-auto mt-16 grid max-w-5xl gap-12 lg:grid-cols-12 lg:items-start">
        {/* Left Side: Selectable Feature Cards with generous spacing and reveal animations */}
        <div className="flex flex-col gap-6 lg:col-span-7">
          {stickyFeatures.map((feat, index) => {
            const isActive = activeTab === feat.id;
            return (
              <RevealOnScroll key={feat.id} delay={index * 80}>
                <div
                  onClick={() => setActiveTab(feat.id)}
                  className={`cursor-pointer rounded-[1.75rem] p-7 transition-all duration-300 ${s.card}`}
                  style={{
                    background: isActive ? `linear-gradient(135deg, ${feat.tint}18, #1A1D27)` : "#141720",
                    border: `1.5px solid ${isActive ? feat.tint : "#242838"}`,
                    boxShadow: isActive ? `0 14px 34px ${feat.tint}1F` : "none",
                  }}
                >
                  <div className="flex items-center justify-between">
                    <span className="font-mono text-[11px] font-bold uppercase tracking-wider" style={{ color: feat.tint }}>
                      {feat.eyebrow}
                    </span>
                    {isActive && (
                      <span className="font-mono text-xs font-semibold px-2.5 py-0.5 rounded-full" style={{ background: `${feat.tint}25`, color: feat.tint }}>
                        Active Preview
                      </span>
                    )}
                  </div>
                  <h3 className="font-display mt-3 text-xl sm:text-2xl font-extrabold text-white">
                    {feat.title} <span style={{ color: feat.tint }}>{feat.titleAccent}</span>
                  </h3>
                  <p className="mt-2.5 text-sm leading-relaxed" style={{ color: "#A8B0C4" }}>
                    {feat.body}
                  </p>

                  {/* Zeptap-Style Step Log */}
                  {isActive && (
                    <div className="mt-5 rounded-xl p-3.5" style={{ background: "#0F1117", border: "1px solid #232736" }}>
                      <div className="flex flex-col gap-2">
                        {feat.log.map((entry, idx) => (
                          <div key={idx} className={`flex items-center justify-between font-mono text-[11px] ${s.logRow}`}>
                            <span className="flex items-center gap-2" style={{ color: "#E8EAED" }}>
                              <span className="h-1.5 w-1.5 rounded-full" style={{ background: feat.tint }} />
                              {entry.step}
                            </span>
                            <span style={{ color: "#7F889E" }}>{entry.detail}</span>
                          </div>
                        ))}
                      </div>
                    </div>
                  )}
                </div>
              </RevealOnScroll>
            );
          })}
        </div>

        {/* Right Side: Sticky Responsive Phone Preview */}
        <div className="lg:sticky lg:top-28 flex justify-center lg:col-span-5">
          <RevealOnScroll delay={100}>
            <div
              className="rounded-[2.5rem] p-4"
              style={{
                background: "radial-gradient(circle at center, rgba(124,110,245,0.15) 0%, transparent 70%)",
              }}
            >
              <PhoneFrame>
                {activeTab === "lesson" && <IntroPhoneScreen />}
                {activeTab === "grading" && <GradingPhoneScreen />}
                {activeTab === "mistakes" && <MistakePhoneScreen />}
                {activeTab === "mastery" && <MasteryPhoneScreen />}
              </PhoneFrame>
            </div>
          </RevealOnScroll>
        </div>
      </div>
    </section>
  );
}

/** Duolingo-style alternating topic rows */
function AlternatingTopics() {
  return (
    <section id="topics" aria-labelledby="topics-title" className="scroll-mt-24 px-4 py-28" style={{ background: "#12141C" }}>
      <RevealOnScroll>
        <SectionHeading id="topics-title" title="Four core topics." muted="One mastery score each." />
        <p className="mx-auto mt-4 max-w-xl text-center text-sm sm:text-base leading-relaxed" style={{ color: "#B9C0D0" }}>
          Mastery is computed 0 to 100 based on verified lesson steps. No fluffy badges.
        </p>
      </RevealOnScroll>

      <div className="mx-auto mt-16 flex max-w-5xl flex-col gap-8">
        {topics.map((t, i) => {
          const theme = topicTheme(t.slug);
          const isEven = i % 2 === 0;

          return (
            <RevealOnScroll key={t.slug} delay={i * 80}>
              <article
                className={`flex flex-col md:flex-row items-center justify-between gap-10 rounded-[2.25rem] p-9 ${s.card}`}
                style={{
                  background: `linear-gradient(135deg, ${theme.tint}18, #181B24 60%)`,
                  border: `1px solid ${theme.tint}33`,
                }}
              >
                <div className={`flex flex-1 flex-col ${isEven ? "md:order-1" : "md:order-2"}`}>
                  <div className="flex items-center gap-3">
                    <span className="text-3xl" aria-hidden>{t.companion}</span>
                    <span className="font-mono text-xs font-bold uppercase tracking-wider" style={{ color: theme.tint }}>
                      {theme.icon} {t.name}
                    </span>
                  </div>
                  <h3 className="font-display mt-3 text-2xl sm:text-3xl font-extrabold text-white">
                    {t.name}
                  </h3>
                  <p className="mt-2.5 text-base leading-relaxed" style={{ color: "#C0C7D6" }}>
                    {t.blurb}
                  </p>
                  <div className="mt-5 flex items-center gap-2 font-mono text-xs" style={{ color: "#8E97AE" }}>
                    <span className="h-1.5 w-1.5 rounded-full" style={{ background: theme.tint }} />
                    {t.detail}
                  </div>
                </div>

                <div className={`flex items-center justify-center p-4 ${isEven ? "md:order-2" : "md:order-1"}`}>
                  <div
                    className="rounded-full p-4"
                    style={{
                      background: `radial-gradient(circle, ${theme.tint}22 0%, transparent 70%)`,
                    }}
                  >
                    <MasteryRing mastery={t.demoMastery} tint={theme.tint} size={110} />
                  </div>
                </div>
              </article>
            </RevealOnScroll>
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
      <div className={`p-8 ${wide ? "md:w-1/2" : ""}`}>
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
  { name: "mira.dev", species: "rabbit", xp: 340 },
  { name: "you", species: "cat", xp: 285, me: true },
  { name: "kenji_b", species: "dog", xp: 260 },
  { name: "sam.ops", species: "hamster", xp: 190 },
];

function Features() {
  return (
    <section id="features" aria-labelledby="features-title" className="scroll-mt-24 px-4 py-28">
      <RevealOnScroll>
        <SectionHeading id="features-title" title="Built like a game." muted="Graded like production." />
      </RevealOnScroll>
      <div className="mx-auto mt-16 grid max-w-5xl gap-6 md:grid-cols-2">
        <RevealOnScroll className="md:col-span-2">
          <FeatureCard
            wide
            eyebrow="Server-side grading"
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
        </RevealOnScroll>

        <RevealOnScroll delay={100}>
          <FeatureCard
            eyebrow="Weekly Leagues"
            title="Climb the ladder."
            muted="Every Monday reset."
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
        </RevealOnScroll>

        <RevealOnScroll delay={150}>
          <FeatureCard
            eyebrow="Daily Streaks"
            title="Show up daily."
            muted="Consistency over intensity."
            body="One lesson keeps your streak alive. Life happens, so a streak freeze covers a missed day without guilt."
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
        </RevealOnScroll>
      </div>
    </section>
  );
}

/** Duolingo-style companion avatars showcase */
function CompanionPillars() {
  const companionsList = [
    { id: "milo", name: "Milo", species: "cat", trait: "Playful & Curious" },
    { id: "buster", name: "Buster", species: "dog", trait: "Loyal & Encouraging" },
    { id: "luna", name: "Luna", species: "rabbit", trait: "Quiet & Focused" },
    { id: "mochi", name: "Mochi", species: "hamster", trait: "Energetic & Quick" },
    { id: "peanut", name: "Peanut", species: "corgi", trait: "Cheerful Companion" },
  ];

  return (
    <section className="px-4 py-24" style={{ background: "#11141D" }}>
      <div className="mx-auto max-w-5xl">
        <RevealOnScroll>
          <div className="text-center">
            <p className="font-mono text-xs font-bold uppercase tracking-widest text-[#7C6EF5]">
              Companion System
            </p>
            <h2 className="font-display mt-2 text-3xl sm:text-4xl font-extrabold text-white">
              Choose someone for your corner.
            </h2>
            <p className="mt-2 text-sm sm:text-base text-[#9FA8BC]">
              Companions react to your wins, shrug off your misses, and celebrate your mastery deltas.
            </p>
          </div>
        </RevealOnScroll>

        <div className="mt-14 grid grid-cols-2 gap-5 sm:grid-cols-3 lg:grid-cols-5">
          {companionsList.map((c, i) => {
            const v = companionVisual(c.id, c.species);
            return (
              <RevealOnScroll key={c.id} delay={i * 60}>
                <div
                  className={`flex flex-col items-center rounded-2xl p-6 text-center transition-all ${s.card}`}
                  style={{
                    background: "#161922",
                    border: `1px solid #232738`,
                  }}
                >
                  <div
                    className="flex h-16 w-16 items-center justify-center rounded-full text-3xl shadow-lg"
                    style={{
                      background: `linear-gradient(135deg, ${v.glow}30, #1A1D27)`,
                      border: `1.5px solid ${v.glow}55`,
                      boxShadow: `0 0 20px ${v.glow}22`,
                    }}
                  >
                    {v.emoji}
                  </div>
                  <h3 className="font-display mt-3 text-base font-bold text-white">{c.name}</h3>
                  <p className="font-mono mt-1 text-[11px]" style={{ color: "#8E97AE" }}>
                    {c.trait}
                  </p>
                </div>
              </RevealOnScroll>
            );
          })}
        </div>
      </div>
    </section>
  );
}

function HowItWorks() {
  return (
    <section id="how" aria-labelledby="how-title" className="scroll-mt-24 px-4 py-28" style={{ background: "#13151C" }}>
      <RevealOnScroll>
        <SectionHeading id="how-title" title="Pick. Play." muted="Level up." />
      </RevealOnScroll>
      <ol className="mx-auto mt-16 grid max-w-5xl gap-6 md:grid-cols-3">
        {steps.map((st, i) => (
          <RevealOnScroll key={st.title} delay={i * 100}>
            <li className="rounded-[1.75rem] p-8" style={{ background: "#1A1D27", border: "1px solid #2E3347" }}>
              <span
                className="font-display flex h-10 w-10 items-center justify-center rounded-full text-lg font-extrabold text-white"
                style={{ background: "#7C6EF5", boxShadow: "0 3px 0 #5B50D4" }}
              >
                {i + 1}
              </span>
              <h3 className="font-display mt-5 text-xl font-extrabold tracking-tight">{st.title}</h3>
              <p className="mt-2.5 text-sm leading-relaxed" style={{ color: "#B9C0D0" }}>
                {st.body}
              </p>
            </li>
          </RevealOnScroll>
        ))}
      </ol>
    </section>
  );
}

function Faq() {
  return (
    <section id="faq" aria-labelledby="faq-title" className="scroll-mt-24 px-4 py-28">
      <RevealOnScroll>
        <SectionHeading id="faq-title" title="Questions?" muted="Answers." />
      </RevealOnScroll>
      <div className="mx-auto mt-14 flex max-w-2xl flex-col gap-4">
        {faqs.map((f, i) => (
          <RevealOnScroll key={f.q} delay={i * 50}>
            <details className="group rounded-2xl px-6 py-5" style={{ background: "#1A1D27", border: "1px solid #2E3347" }}>
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
          </RevealOnScroll>
        ))}
      </div>
    </section>
  );
}

function FinalCta() {
  return (
    <section className={`relative overflow-hidden px-4 py-36 text-center ${s.sky}`}>
      <StickerTile s={{ ...heroStickers[0], className: "left-[10%] top-[20%]" }} size={64} />
      <StickerTile s={{ ...heroStickers[3], className: "right-[10%] bottom-[18%]" }} size={64} />
      <RevealOnScroll>
        <h2 className={`font-display mx-auto max-w-3xl text-4xl font-extrabold sm:text-6xl ${s.headline}`}>
          {finalCta.headline}
          <span className="block" style={{ color: "#8B92A8" }}>
            {finalCta.headlineMuted}
          </span>
        </h2>
        <div className="mt-12">
          <CtaLink href="/register">{finalCta.cta}</CtaLink>
        </div>
      </RevealOnScroll>
    </section>
  );
}

function Footer() {
  return (
    <footer className="px-4 py-12" style={{ borderTop: "1px solid #1F2230", paddingBottom: "max(2.5rem, env(safe-area-inset-bottom))" }}>
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
