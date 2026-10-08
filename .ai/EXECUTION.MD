# EXECUTION.md

## AGENT INSTRUCTIONS

This document is the source of truth for implementation priorities.

Before implementing any feature, read: PRODUCT.MD, ARCHITECTURE.MD, CONTENT_SYSTEM.MD, EXECUTION.MD.

* If a feature is not listed in the **current phase**: do not build it.
* If uncertain: ask for clarification.
* If a conflict exists between code and these documents: the documents are updated **in the same change** (see Change Process). They are never left stale.

---

## Change Process

If you discover a better design, or a decision here is wrong, do not implement it immediately. Document:

1. Problem
2. Proposed solution
3. Benefits
4. Risks
5. Migration impact

Then request approval. When approved, update every affected document in the same change so no conflicting statement survives.

---

## Phase History

| Phase | Outcome |
|---|---|
| **Phase A, Foundation Rebuild** | **Superseded.** A1 skill graph, A2 content architecture, A3 question schema: done and kept (skills, subskills, `user_skill_scores`). A4 company-based Readiness V2: **superseded** by per-topic mastery and removed in L1. A5 content management: **replaced** by content-as-code (L2). |

---

# CURRENT PHASE

## Phase L: Lesson Platform

**Status: ACTIVE.** Goal: turn Prepio into a production lesson-based learning platform for working engineers: short lessons, server-graded exercises, per-topic mastery, a journey that means something.

### Standing rules (apply to every task in this phase)

* **Production quality only.** No hardcoded content, mock data, fixtures, or stubbed paths in shipped code.
* **Dead-code sweep at the end of each task:** orphaned symbols, routes, env vars, feature flags, tests, components, tables, docs.
* **Gates before a task is done:**
  1. `go build ./... && go vet ./... && go test ./...`
  2. Web: lint, typecheck, build, plus unit tests for non-trivial logic
  3. Migrations apply **up and down** on a fresh database
  4. UI work: screenshots reviewed at **390×844 (mobile) and 1280×800 (desktop)**
  5. Docs updated in the same change
* Historical migrations are never edited. Work happens on feature branches.
* The mobile (Flutter) client is **paused**: not updated or gated in this phase. It will be rebuilt against the new contract after the web experience is validated.

### Tasks (in order)

**L1, Purge company targeting and company-based readiness**
**Status: done** (commit f52f6f8).
Remove: `user_targets`, `question_tags`, `company_skill_weights` (new migration, with down); company logic in user, question, progress, and gateway services; `shared/readiness`; company constants, config, and events fields; the `READINESS_V2` flag; web and dashboard company UI; the dashboard readiness card (reintroduced as topic mastery in L4). Keep skills, subskills, `question_skills`, `user_skill_scores`, XP, gems, streaks.
*Done when:* grep for `company|companies|target_compan|readiness_v2|question_tags|user_targets` is clean outside historical migrations; migration up/down passes; gates pass; login → onboarding (experience, companion) → dashboard works.

**L2, Lesson platform backend**
**Status: done** (merged in #3). The first lesson (`why-caches-exist`, world `First Ascent`) uses only an intro and `mcq` steps so the L3 player, which starts with MCQ, can run it; world and node names are provisional until the outline is designed with the product owner (L8). Skills added: the four Backend & Production skills. The lesson journey is served at `GET /path`.
Schema (`lessons`, `lesson_steps`, `lesson_skills`, `lesson_attempts`, `lesson_step_results`, `mastery_ledger`); topic catalog migration; `content/` with the validator and `content-sync` (wired into the deploy entrypoint after migrations); grading package (pure, table-tested) and prose via the rubric evaluator; lesson runtime API; `lesson.completed` event; Progress consumer writes mastery + ledger + XP; Journey binds nodes to lessons and derives status from attempts; the first real lesson of the Working Professionals path authored in `content/`.
*Done when:* integration tests cover the attempt lifecycle, idempotency, invalid answers, and completion rules; a test proves client payloads contain no answers; `content-sync` is idempotent; gates pass.

**L3, Lesson player (web) and legacy removal**
**Status: done** (merged in #2). Dropped with the legacy loop, beyond the list below: `question_skills`, `node_pools`, `node_skills`, and `user_journey_progress` (all reference `questions` or the removed pool/index builders, and nothing reads them). The skill catalog, `user_skill_scores`, XP, gems, and streaks are kept. The web exercise types beyond `mcq` arrive in L6; a lesson step the player cannot render shows a leave-lesson message rather than a placeholder.
Player shell (focus mode, top bar, progress, combo), feedback tray, MCQ step, miss re-queue, celebration, summary card, journey preview sheet, node-completion animation, unlock messaging; companion in the player. Then **remove the legacy loop**: `/challenge`, daily papers, the legacy question/answer flow and `question.answered` event, journey index/pool builders and the `JourneyPoolSelection` flag, `question_pools*`, `questions`, `daily_papers*`, `user_question_history` (new migration). Keep the rubric evaluator for `prose`.
*Done when:* a user completes a real lesson end to end from the journey; legacy grep is clean; screenshots reviewed; gates pass.

**L4, Onboarding and dashboard**
**Status: done** (merged in #4). Onboarding is topics (1-3, ordered) then experience then companion; `focus_topics` live in `user_focus_topics` (User domain) and come back on the profile. The dashboard shows a Continue card for the path's current node and one readiness card per topic (focus topics first, never hidden), and the league/quests placeholders, their routes and nav tabs are gone. The celebration screen explains each mastery change by topic and skill. Focus topics currently reorder and highlight the dashboard; reordering the journey itself waits until there is more than one world.
Topic picker (1–3) → experience → companion; `focus_topics` stored by the User domain. Dashboard: per-topic mastery rings, Continue → next lesson, explained mastery changes on the celebration screen. Remove "coming soon" quest/league cards and dead nav tabs until their phases.
*Done when:* new users pick topics and land on a dashboard that shows real topic mastery and the next lesson; no placeholders remain; gates pass.

**L5, Intro motion stage:** beat-based motion renderer, tap to advance, skip, replay, optional `mediaUrl`.
**Status: done** (merged in #5). Words reveal per beat with the emphasis phrase glowing, named visuals (`bars`, `flow`, `pulse`) sit above the text, beats auto-advance with a pause control, the intro can be replayed from its end and from the player header, and `mediaUrl` renders images or video only from https or site-relative URLs. Reduced motion shows full text with manual advance. The validator now checks beat timing, visual names, and media URLs.

**L5.1, Post-merge hardening** (review of L2–L5)
**Status: implemented** on branch `feat/l9-leagues-landing`, pending review.
* Content sync can replace or move the lesson on a node: a node holds at most one *live* lesson (migration 000040, a deferrable exclusion constraint), and removed lessons are deprecated before the upsert.
* Progress balance writes lock the user's row and change it relative to what is stored, so concurrent lesson and streak consumers never overwrite each other's XP or gems; the streak gem bonus is idempotent per event.
* A late or redelivered `lesson.completed` never resets a streak or moves its date backwards.
* Kafka consumers retry a failing handler with backoff, then log and skip the message instead of stopping for good.
* A failed migration stops the deploy (`prod-entrypoint.sh`, `start-dev.sh`).
* Web: the refresh cookie is sent (`credentials: "include"`; `Secure` + `SameSite=None` over HTTPS), so a reload no longer signs the user out, and sign-out revokes the cookie; API errors with non-JSON bodies are handled; the celebration screen no longer sums a topic number on the client; Enter checks after a mouse selection; the journey only claims unlocks the server reported.
* Removed: `GET /skills/readiness` and its analysis code, the unused `LedgerStore`, dead constants and errors, `XpProgressBar`, scaffold SVGs.
* Known limit: rolling migrations back below 000034 fails at historical migration 000028 (`uuid LIKE`); rollbacks are supported down to 000033.

**L9, Leagues and leaderboard**
**Status: implemented** on branch `feat/l9-leagues-landing`, pending review. Pulled ahead of L6–L8 by the product owner on 2026-10-07 (change record below).
Weekly leagues owned by Progress. A learner joins the current week's league the first time Progress grants them XP that week and is placed in the oldest cohort of their tier with room (`config.LeagueCohortSize`, 30). League XP is written in the same transaction that grants the XP, so it is counted exactly once, and it counts in the week it is granted. Standings rank by weekly XP, ties to whoever got there first. At week end (Monday 00:00 UTC) the top `LeaguePromoteCount` move up a tier and, only in cohorts larger than promote + demote, the bottom `LeagueDemoteCount` move down; there is no demotion from the lowest tier and none for weeks you skip (you keep the tier your last week earned). Results are derived from final standings when the learner next joins or views the league, so no job runs at the week boundary. Tiers, sizes, and counts live in `config/leagues.go`. `GET /progress/league` (Progress) serves standings; `GET /league` (gateway) attaches usernames and companions from the User domain (`POST /users/public-cards`) and fails closed. Web: `/league` (tier, ladder, countdown, standings with promotion and drop zones named in words, last week's result) and a League tab.
*Done when:* rules are table-tested; store tests cover joining, ties, cohort overflow, promotion, and skipped weeks; gateway aggregation is tested including failure; gates pass; screenshots reviewed.

**L10, Landing page**
**Status: implemented** on branch `feat/l9-leagues-landing`, pending review. Pulled ahead by the product owner on 2026-10-07 (change record below).
`/` is a public landing page; signed-in visitors are sent on to the dashboard or onboarding. Duolingo's structure (one promise, companions, one obvious call to action) with Zeptap-inspired components: floating pill nav, an announcement pill, a two-tone headline (bold claim, muted second half), floating sticker icons, a window holding three looping phone screens (journey, a graded exercise with its feedback tray, the celebration with mastery moving), a scrolling strip of skill chips, tinted topic cards with mastery rings, feature cards, three steps, FAQ, and a closing call to action. Dark theme and purple accent kept. Marketing copy lives in `web/src/lib/landing/content.ts`; the phone screens are illustrations, never real lesson content or user data. Motion stops under `prefers-reduced-motion`.

**Change record (L9, L10).** *Problem:* leagues and a landing page were out of scope until after validation, but the product owner wants a weekly social reason to return and a public front door for the beta. *Proposal:* the two tasks above. *Benefits:* a weekly retention loop that rewards showing up (XP comes only from lessons); a page that explains Prepio to new users. *Risks:* leagues can feel punishing, so demotion is limited (lowest tier exempt, small cohorts never demote, skipped weeks keep the tier, encouraging copy), and with few users cohorts are small (everyone near the top promotes). *Migration impact:* migration 000041 adds `league_cohorts` and `league_memberships`; no existing data changes. Quests stay out of scope.

**L6, Exercise types:** `true_false`, `fill_blank` (tap-to-place), `arrange` (tap-to-order); validator rules and player support.
**Status: implemented** on branch `feat/l9-leagues-landing`, pending review. Grading, validation, and the answer-free client payloads for all three types existed since L2; L6 adds the web player. True/false takes taps or keys 1/2 (T/F); fill-in-the-blank places a tapped bank chip in the first empty blank and a tapped blank returns its chip; arrange appends a tapped item to the next position and a tapped position takes it back. Enter checks once the answer is complete. After a miss the tray states the server's correct answer in words (the option, True/False, each blank, or the full order). Drafts and their conversion to wire answers live in `web/src/lib/lesson/answers.ts` (unit-tested); the client never compares an answer with the correct one. `prose` is still not rendered by the player and shows the leave-lesson message.

**L7, Retention polish:** mistakes review, replay intro, sound toggle, keyboard and accessibility audit.
**Status: implemented** on branch `feat/l9-leagues-landing`, pending review. Mistakes review: the summary card lists each exercise missed in this lesson (first miss only) with the server's correct answer and explanation; it is UI state from this session, not stored. Replay intro shipped in L5. Sound: three short cues (correct, miss, lesson complete) synthesized with WebAudio, so there are no audio files; off by default, toggled in the player header, remembered on this device (`localStorage`). Accessibility: while the leave-lesson dialog is open the player behind it is `inert` and exercise shortcut keys are ignored; browser Back asks before leaving an unfinished lesson; the header's intro button is icon-only on small screens.

**L8, Working Professionals path content:** worlds, nodes, and lessons authored with the product owner. Quality over quantity; first world complete before the second begins.
**Status: implemented** on branch `feat/l9-leagues-landing`, pending review. Planned in `content/PLAN.md` (outline, authoring standard, and the review that shaped it) with the product owner's go-ahead on 2026-10-08. Four worlds, one per topic, 33 lessons: First Ascent (System Design, 9 nodes, extending the existing world), The Engine Room (Backend & Production), Blueprint Workshop (Low-Level Design), and Algorithm Gardens (DSA Refresher), 8 nodes each. Each world is a chain ending in a boss; worlds open in parallel so a topic never locks another. World 1 was written first and became the template for the others. `why-caches-exist` was rebalanced (version 2) so no lesson gives its answer away by option position. No `prose` steps until the player renders them.

### Phase L Success Criteria

* A user can complete lessons end to end with instant, server-graded feedback.
* Topic mastery is real, explained, and visible on the dashboard.
* Content is authored as code, validated in CI, and synced on deploy.
* No company logic, legacy question loop, or dead flags remain.

---

## OUT OF SCOPE (forbidden before launch)

Arena · Multiplayer · Guilds · Voice interviews · AI interviewers · AI resume review · Recruiter marketplace · Season pass · Advanced recommendation engine · Community content · Complex social features · Company-specific prep · Hearts/lives · Code-execution sandbox · AI-generated lessons · Quests (until a later phase) · Mobile client updates · Additional audience paths

These solve scale or breadth, not validation.

## DEFINITION OF DONE

A task is **not** done because the code compiles, the UI renders, or the API responds. It is done when:

* Architecture and content model respected
* Documentation updated
* Tests pass and all gates above are met
* UI matches PRODUCT.MD
* No hardcoded content or assumptions
* No ownership violations
* Dead code removed

## AGENT DECISION RULES

1. Prefer configurable, content-driven systems.
2. Avoid hardcoding.
3. Avoid shortcuts that require rewrites.
4. Avoid speculative scale work.
5. Avoid creating new services.
6. Follow existing ownership.
7. Remove what you make dead.

## LAUNCH CRITERIA

* Users can progress through the path lesson by lesson.
* Users understand their topic mastery and why it changes.
* Users feel progression and return voluntarily.
* Content supports at least 30 days of daily use.

The goal is validation, not feature completeness.

## FINAL RULE

Do not optimize for 100,000 users. Optimize for 20 users who genuinely love the product. If 20 users love it, scaling can be solved later. If they do not, nothing else matters.
