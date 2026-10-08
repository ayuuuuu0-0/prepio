# PRODUCT.md

## AGENT INSTRUCTIONS

Before any product, UI, UX, or feature decision, read in this order:

1. PRODUCT.MD
2. ARCHITECTURE.MD
3. CONTENT_SYSTEM.MD
4. EXECUTION.MD

If a proposed implementation conflicts with any of them: **STOP. Do not implement. Explain the conflict.**
If you believe a document is wrong, use the change process in EXECUTION.MD. Never work around a document silently, and never leave a document contradicting the code.

---

## What Is Prepio?

Prepio is a progression game for working engineers. Getting measurably better at the skills that matter on the job and in interviews is the reward.

Most learning products fail the same way: users start motivated, then quit, not because they are incapable but because the experience is a grind. Prepio fixes the experience: short lessons, instant feedback, and visible mastery.

## Core Belief

People do not wake up wanting to do exercises. They want to be better engineers, grow in their careers, and walk into hard conversations (design reviews, interviews, incidents) with confidence. Lessons are the vehicle. **Career progression is the product.**

## Positioning

We are NOT:

* LeetCode or any problem-grinding site
* A question bank
* A video course library
* A company-specific interview-prep tool

We borrow **Duolingo's engagement model** (short daily lessons, instant feedback, a path, streaks) and **boot.dev's hands-on feel** (write and complete real code and designs), applied to engineering depth: system design, backend and production concepts, low-level design, and a DSA refresher.

## Target User

**Primary: working software engineers** who want to deepen system design, production/backend knowledge, low-level design, and refresh DSA, for their job or for interviews.

Other audiences (students, early-career) would be separate paths. They are out of scope until the primary path is validated.

The product should feel professional. Not childish. Not corporate.

## Core Product Loop

```
User opens app → companion greets
→ dashboard shows topic mastery + the next lesson
→ Continue → journey path → (tap node → preview) → lesson
→ lesson: intro → exercises with instant feedback → celebration → summary
→ mastery updates (explained) → next node unlocks → companion reacts
→ XP counts toward this week's league
→ user returns tomorrow
```

Every system exists to strengthen this loop.

## Product Pillars

Every feature must support at least one pillar. If it supports none, it should not exist.

### 1. Progression
Users must always see progress: journey, mastery, levels, achievements, companion evolution. Users never wonder "am I improving?"

### 2. Readiness (= topic mastery)
Readiness is the most important metric. It is **per-topic mastery from 0–100**, for the topics: **System Design, Backend & Production, Low-Level Design, DSA Refresher.**

* Derived from skill mastery. Never from XP, levels, or streaks.
* Always explainable: why it rose, which lessons and skills moved it, and what to do next.
* Not company-specific. There is no company selection and no company readiness.

### 3. Journey
The journey is the primary experience. Users move through worlds, not lists, tables, or folders. World and node names are authored in content and should be memorable.

### 4. Companions
Companions are emotional progression systems, not avatars. They appear on every major screen, **including inside the lesson player**, and react to results. If a companion could be removed without changing the experience, the implementation is wrong.

### 5. Collection
Achievements, badges, titles, world trophies, companion evolutions, skins.

---

## The Lesson Experience (UX contract)

A lesson is the unit of learning. These rules are product requirements, not suggestions.

**Shape**
* ≤ 5 minutes: one intro (30–40 s, always skippable) + 3–5 exercises.
* Focus mode: bottom navigation hidden; leaving asks for confirmation.

**Every exercise follows one flow:** select → **Check** → feedback tray → **Continue**.
* Feedback is instant and never relies on colour alone (icon + text).
* Correct: short reason it works. Incorrect: the correct answer + why, in an encouraging tone.
* Missed exercises **return at the end of the lesson**. There are no hearts, lives, or penalties.

**Exercise types:** multiple choice, true/false, fill-in-the-blank (word bank), arrange-in-order, and rare rubric-graded prose. Interaction is tap-based (no drag required), keyboard accessible (1–4 select, Enter checks), and respects reduced-motion settings.

**Finish**
1. Celebration: XP, accuracy, streak, companion reaction.
2. Summary card: key takeaways and the worked answer for design problems.
3. Back to the journey; the completed node animates and the next one unlocks.

**Journey interaction:** tapping a node opens a preview (what you'll learn, time, XP). Locked nodes state exactly what unlocks them. No dead taps.

**Sound:** optional, default off, remembered.

## Onboarding

1. **What do you want to sharpen?** Pick 1–3 topics. Topics reorder and highlight the path. They never lock content.
2. **Experience level.** Sets the starting point.
3. **Choose a companion.**

No company selection. No "what are you here for" quiz.

## Dashboard Philosophy

The dashboard is the most important screen. In seconds, the user must understand: where they are (topic mastery), what to do next (Continue → next lesson), how they're improving, and what reward is next. It motivates action; it does not display analytics. **It never shows placeholders for unbuilt features.**

## Visual Identity

Dark theme, purple accents, modern gradients, developer-oriented. Intentional.
Do not copy Duolingo's colours or mascot. We want its engagement model, not its look.

## Rewards Philosophy

Every action produces visible feedback: XP gained, mastery moved (and why), companion reacts, progress updates, an animation plays. Streaks reward showing up: completing one lesson keeps the streak.

## Leagues

A weekly reason to come back. Your first lesson of the week puts you on a leaderboard with up to 30 learners in your league tier; XP from lessons ranks you; when the week ends the top ranks move up a tier. Leagues are ranked by XP and never by mastery: they reward showing up, while mastery stays the honest measure of skill.

Leagues must follow the Psychology rules below: demotion is limited (never from the lowest tier, never in small cohorts, never for a week you skipped), zones are named in words as well as colour, and copy frames a lower finish as a chance to win it back, never as failure. Rules and numbers live in config, not code (see EXECUTION.MD, L9).

## Landing Page

The public front door at `/`: one promise, our companions, and one obvious call to action (Duolingo's structure), built with Zeptap-style components (floating pill nav, announcement pill, two-tone headlines, sticker icons, a window of live phone screens, a chip strip, tinted feature cards). It keeps Prepio's dark theme and purple accent. Its phone screens illustrate the product; they never show real lesson content or user data. Signed-in visitors go straight to the app.

## Psychology

Create motivation, curiosity, momentum, achievement. Avoid fear, punishment, guilt, shame. Users should feel encouraged to return, not pressured.

## Anti-Patterns (forbidden)

* Generic question lists or a "questions page"
* Dashboard analytics: charts, tables, and metrics everywhere
* Enterprise/admin UI: white backgrounds, plain cards, data-heavy screens
* Company-specific preparation features
* Punishment mechanics: hearts, lives, shaming, guilt copy
* "Coming soon" placeholders in shipped UI
* Hardcoded content or mock data in shipped paths
* Feature creep: AI interviewers, voice, guilds, marketplaces, real-time multiplayer (see EXECUTION.MD). Weekly leagues are asynchronous leaderboards, not multiplayer.

## Success Metrics

A successful user understands their mastery, understands the path, completes lessons, returns consistently, feels improvement, and asks for more content, without external motivation, manual tracking, or a complex onboarding.

## Long-Term Goal

A career progression RPG for engineers. Users should not feel like they are studying; they should feel like they are leveling up. That feeling is the product.
