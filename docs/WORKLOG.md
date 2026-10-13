# Prepio Engineering Work Log

A comprehensive, plain-language record of architectural decisions, implementation details, testing verifications, and operational status for the Prepio platform.

Current Branch: `feat/l9-leagues-landing`

---

## Executive Summary of Completed Phases

| Phase | Title | Primary Scope | Status & Commits |
|---|---|---|---|
| **L1** | Company Purge | Eliminated company-specific readiness, target companies, and legacy weighting tables | Done (`f52f6f8`) |
| **L2** | Lesson Platform Backend | Schemas, pure Go deterministic grader, content-as-code CLI, runtime APIs | Done (PR #3, `394efff`) |
| **L3** | Web Lesson Player | Interactive player shell, celebration screen, legacy daily loop dropped | Done (PR #2, `05f01a0`, `a10eb9c`) |
| **L4** | Onboarding & Dashboard | Focus topics (1–3), per-topic readiness calculation, Continue card | Done (PR #4, `297a956`, `332b222`) |
| **L5** | Intro Motion Stage | Beat-based animated conceptual intro, SVGs, audio/pause/replay controls | Done (PR #5, `b6e5ac6`) |
| **L5.1**| Post-Merge Hardening | Concurrency row locks on XP, cookie auth persistence, Kafka retries, DB constraints | Done (`7b4d384`) |
| **L6** | Exercise Types & Prose | True/False, Fill-in-Blank, Arrange, and rubric-graded Prose exercises | Done (`59fdf69`, `1d3596e`) |
| **L7** | Retention & Polish | Mistakes review, synthesized WebAudio sound cues, defensive exit navigation | Done (`3d4e786`) |
| **L8** | Working Professionals Content | Worlds 1–5, 40 lessons, 175 exercises in YAML, 100% topic skill coverage | Done (`e4dc964`, `2d3f1c0`) |
| **L9** | Weekly Leagues & Dashboard | Cohorts of 30, promotion/demotion tiers, leaderboard, dashboard league card | Done (`cb961e2`, `d5cd5af`) |
| **L10**| Public Landing Page | Dark-themed, high-conversion landing page at `/` with interactive previews | Done (`d247e21`) |
| **L11**| Focus-Ordered Journey | Worlds declare topics; learner's focus topics prioritize the learning path | Done (`3af9771`) |
| **L12**| Achievements Collection | Milestones in DB transaction with XP, profile badges, celebration cards | Done (`c4bd745`) |

---

## Phase L1 · Purge Company Targeting & Readiness

### 1. Decision & Rationale
* **The Problem:** Prepio was originally built around company-specific readiness (e.g., "Google Readiness: 42%"). Skill weights were tied to arbitrary company formulas in `company_skill_weights`, questions were tagged with company names, and users selected target employers during onboarding. This was a demo gimmick: no engineering organization tests solely against a static formula, and working software engineers care about concrete topic mastery (System Design, Backend, LLD, DSA). It introduced dead tables, distorted progression calculations, and added unnecessary cognitive overhead.
* **The Decision:** Completely purge company targeting, company-based readiness V1/V2, and all company database tables. Shift the entire product to **per-topic mastery** (System Design, Backend & Production, Low-Level Design, DSA Refresher).
* **Alternatives Rejected:** Keeping company readiness alongside topic mastery was rejected because maintaining two divergent progression models doubles maintenance, confuses learners, and preserves technical debt.

### 2. Plan
* Add database Migration 000034 to drop `company_skill_weights`, `user_targets`, and `question_tags`.
* Strip company endpoints and handlers from User, Question, Progress, and Gateway services.
* Replace company readiness calculations in the Progress service with pure skill mastery vectors and topic aggregations.
* Simplify the onboarding flow from 4 steps to 2 clean steps (experience level and companion selection).
* Remove company tags, filters, colors, and readiness rings from the web frontend.

### 3. Implementation
* **Database:** Created `migrations/000034_drop_company_tables.up.sql` (`DROP TABLE IF EXISTS company_skill_weights, user_targets, question_tags;`) and corresponding `.down.sql`.
* **User Service:** Deleted `services/user/internal/store/targets.go`; simplified onboarding DTOs and handlers to accept only `experienceLevel` and `companionId`.
* **Question Service:** Removed routes `/questions/stats/readiness` and `/questions/companies`; purged `CompanyTags` and `ReadinessDelta` from models, queries, and XP reward multipliers.
* **Progress Service:** Purged `CompanySkillWeight` stores; rewrote analysis logic with `BuildSkillGaps` and `BuildSkillMasteryExplanation` strictly based on skill scores.
* **Gateway Service:** Cleaned dashboard handlers of all company readiness cards and fetchers; removed `/dashboard/readiness` and `/companies/readiness`.
* **Shared & Config:** Deleted `shared/readiness` (`types.go`, `v1.go`); deleted `constants/companies.go`; removed company feature flags and rewards multipliers.
* **Web Client:** Cleaned `web/src/lib/api.ts`; deleted `ReadinessRing.tsx`; redesigned onboarding and dashboard to remove all company branding and rings.

### 4. Testing & Verification
* **Grep Audit:** Targeted grep `git grep -i -E "company|companies|target_compan|readiness_v2|question_tags|user_targets"` confirmed zero occurrences outside historical migrations and docs.
* **Database Rollback:** Applied Migration 000034 down and up cleanly on PostgreSQL.
* **Service Builds & Tests:** `go build ./...` and `go test ./...` passed across all services.
* **Frontend Build:** `npm run build` executed cleanly with 14 static pages generated without type or lint errors.

---

## Phase L2 · Lesson Platform Backend & Content-as-Code

### 1. Decision & Rationale
* **The Problem:** Previously, questions were loose records scattered across SQL seed files and pulled into dynamic "daily papers". This made curriculum progression unpredictable and impossible to version-control or review.
* **The Decision:** Implement a formal learning hierarchy: `World → Node → Lesson → Step`. Content must be authored as code in version-controlled YAML files under `content/`, validated in CI, and synchronized idempotently during deployment. Grading must happen **100% server-side**—the client must never receive answer keys, solution guides, or grading criteria.
* **Alternatives Rejected:** Client-side grading was rejected because shipping answers to the browser invites trivial cheating (damaging future competitive leagues) and causes business logic drift between Go and TypeScript.

### 2. Plan
* Design schema for lessons (`lessons`, `lesson_steps`, `lesson_skills`, `lesson_attempts`, `lesson_step_results`, `mastery_ledger`) in Migration 000035.
* Implement a pure Go deterministic grading engine (`engine/grader`) supporting multiple exercise types.
* Create CLI tools: `content-validator` for linting/integrity, and `content-sync` for idempotent deployment upserts.
* Build the lesson runtime API (`GET /lessons/{id}`, `POST /lessons/{id}/submit`, `POST /lessons/{id}/complete`).
* Wire Kafka event `lesson.completed` to the Progress service to update `user_skill_scores` and record ledger deltas.

### 3. Implementation
* **Database:** Created `migrations/000035_create_lesson_platform.up.sql` creating tables, foreign keys, unique slugs, and performance indexes.
* **Deterministic Grader:** Created pure Go evaluation functions for Multiple Choice, True/False, Fill-in-Blank, and Arrange, handling whitespace normalization and case insensitivity.
* **Payload Sanitizer:** Stripped `answer`, `explanation`, and grading rubrics from `GET /lessons/{id}` responses, ensuring learners only see prompts, choices, and structure.
* **Runtime API:** Handlers track ongoing attempts in `lesson_attempts`, grade individual steps, record first-attempt vs. retry outcomes in `lesson_step_results`, and re-queue mistakes without life/heart penalties.
* **Progress Consumer:** Consumes `lesson.completed`, computes skill deltas based on lesson skill weights, writes audit records to `mastery_ledger`, and updates user XP.
* **Content Tooling:** Authored `tools/content-validator` (validates YAML schema, acyclic node references, skill weights summing to 1.0) and `tools/content-sync` (executes single-transaction idempotent upserts).
* **Seed Content:** Authored the initial production lesson `why-caches-exist` in `content/lessons/backend/`.

### 4. Testing & Verification
* **Grader Table Tests:** Extensive table-driven unit tests verified grading accuracy, edge cases, and whitespace tolerance.
* **Zero-Leakage Security Test:** Integration tests inspected JSON output of `GET /lessons/{id}` to guarantee no answers or explanations leaked.
* **Sync Idempotency:** Verified that running `content-sync` multiple times produces identical database states with zero duplicates.
* **Attempt Lifecycle:** Verified full flow: start attempt → submit wrong answer (re-queued) → submit right answer → complete lesson → verify `mastery_ledger` and XP balance.

---

## Phase L3 · Web Lesson Player & Legacy Loop Removal

### 1. Decision & Rationale
* **The Problem:** The web frontend had disjointed challenge screens (`/challenge`, `/daily-paper`) and lacked a cohesive, distraction-free player.
* **The Decision:** Build a full-screen, dedicated web lesson player at `/lesson/[lessonId]` inspired by Duolingo's flow and Boot.dev's focus. Completely delete the legacy daily question loop and its associated database tables to eliminate dead code and technical debt.
* **Alternatives Rejected:** Retaining `/challenge` as an alternate game mode was rejected because dual learning paths confuse learners and preserve hundreds of lines of legacy code.

### 2. Plan
* Create a full-screen player shell with top progress bar, combo multiplier, step renderer, and responsive bottom feedback tray.
* Build a celebration screen displaying XP earned, streak updates, and skill mastery breakdowns.
* Implement Migration 000036 to drop legacy tables (`questions`, `daily_papers`, `question_pools`, `user_question_history`, `node_pools`, `node_skills`).
* Purge legacy endpoints across Question and Gateway services.

### 3. Implementation
* **Web Components:** Built `PlayerShell.tsx`, `FeedbackTray.tsx`, `MCQStep.tsx`, and `Celebration.tsx`.
* **State Management:** Handled step queues, mistake re-queuing (re-inserting failed steps at the end of the queue), combo streaks, and completion transitions.
* **Keyboard Navigation:** Mapped number keys 1–4 to options, Enter to check/advance, and Escape to open the exit modal.
* **Database:** Executed `migrations/000036_drop_legacy_questions.up.sql` removing all legacy question tables.
* **Cleanup:** Deleted `/challenge`, daily paper components, legacy evaluation services, and old question DTOs.

### 4. Testing & Verification
* **End-to-End Playthrough:** Played lessons in a real browser from journey node click through feedback trays to the celebration screen.
* **Mistake Re-queuing:** Confirmed that missed exercises re-appear at the end of the lesson queue until answered correctly.
* **Dead Code Sweep:** Confirmed zero references to dropped tables or deleted routes across the codebase.
* **Responsive Layout:** Verified proper scaling across mobile (390×844) and desktop (1280×800) viewports.

---

## Phase L4 · Onboarding & Dashboard (Per-Topic Mastery)

### 1. Decision & Rationale
* **The Problem:** Onboarding previously asked for general company targets. The dashboard was cluttered with "coming soon" placeholders and lacked actionable guidance.
* **The Decision:** Working engineers have specific learning goals. Onboarding must prompt for 1–3 focus topics (System Design, Backend & Production, Low-Level Design, DSA Refresher). The dashboard must answer two questions immediately: *"What do I do next?"* (prominent Continue card) and *"How am I progressing?"* (Topic mastery rings).
* **Alternatives Rejected:** Asking generic questions like "What role are you targeting?" was rejected in favor of actionable, topic-driven choices.

### 2. Plan
* Add `user_focus_topics` table in Migration 000037.
* Rebuild onboarding as a 3-step sequence: Topics (1–3) → Experience level → Companion selection.
* Implement `GET /progress/topics` to compute 0–100% normalized readiness scores per topic.
* Update Gateway's `GET /dashboard` to return focus topics, topic readiness, and the next active lesson node.
* Remove placeholder UI cards (leagues/quests) until their formal implementations.

### 3. Implementation
* **Database:** Created `migrations/000037_user_focus_topics.up.sql`.
* **User Service:** Added endpoints `PUT /users/focus-topics` and `GET /users/focus-topics`.
* **Question Service:** Added `GET /topics` catalog endpoint returning slug, title, and description.
* **Progress Service:** Implemented `GET /progress/topics` calculating topic readiness from underlying skill masteries.
* **Gateway Service:** Updated dashboard aggregation to fetch profile, focus topics, next journey node, and readiness metrics.
* **Web Frontend:** Redesigned `/onboarding` with topic chips; built `/dashboard` with Continue card and topic mastery rings; updated celebration screen to display exact topic readiness deltas.

### 4. Testing & Verification
* **HTTP Smoke Test:** Built and ran `scripts/lesson-smoke.mjs` verifying the flow: register → select topics → verify dashboard topic readiness.
* **Mathematical Verification:** Tested topic readiness aggregation formulas against edge cases (0% progress, 100% mastery, unselected topics).
* **UI Polish:** Verified responsive layouts and fixed text-wrapping issues on long topic titles (`57d28d8`).

---

## Phase L5 & L5.1 · Intro Motion Stage & Post-Merge Hardening

### 1. Decision & Rationale
* **L5 Decision:** Jumping directly into exercises feels like an exam. Every lesson must begin with an engaging, interactive conceptual intro stage that teaches the concept before testing it. To make reading compelling, we implemented beat-based motion (words appearing on beat with glowing emphasis, SVG visuals, and auto-advance).
* **L5.1 Decision:** Reviewing the merged L2–L5 codebase identified several critical production defects: session logouts on browser refresh, lost XP under concurrent events, content-sync crash loops, and unhandled migration failures. These required immediate hardening before feature expansion.

### 2. Plan
* **L5:** Define intro step schema (`beats`, `visual`, `emphasis`, `mediaUrl`). Build web intro renderer with play/pause/replay and `prefers-reduced-motion` support.
* **L5.1:** Add exclusion constraint to ensure at most one live lesson per node (Migration 000040); add row-level locking (`SELECT ... FOR UPDATE`) in Progress balance updates; implement Kafka retry with exponential backoff; configure secure cookie auth with `credentials: "include"`.

### 3. Implementation
* **Web Intro Stage:** Built `IntroStep.tsx` rendering beat sequences, glowing text spans, SVG visuals (`bars`, `flow`, `pulse`), and header replay controls.
* **Content Validation:** Extended `content-validator` to enforce beat timing limits, visual names, and secure HTTPS/relative media URLs.
* **Database (L5.1):** Migration `000040_lessons_node_unique_live.up.sql` added a deferrable exclusion constraint preventing duplicate live lessons on a node during sync.
* **Auth Persistence (L5.1):** Updated `shared/auth/cookie.go` and Next.js `api.ts` to use `SameSite=None`, `Secure`, and `credentials: "include"`, resolving reload logouts.
* **Concurrency (L5.1):** Wrapped XP, gem, and streak balance writes in `SELECT ... FOR UPDATE` transactions to prevent lost updates under concurrent Kafka delivery.
* **Kafka Resilience (L5.1):** Added `shared/kafka/retry.go` with exponential backoff before logging and skipping poison-pill messages.

### 4. Testing & Verification
* **Motion Verification:** Tested intro playback, pause/resume, replay, and verified fallback behavior under `prefers-reduced-motion`.
* **Concurrent Balance Test:** Simulated concurrent reward events; verified zero lost XP or gems.
* **Auth Persistence:** Verified session persistence across page reloads in Chrome and Safari.
* **Migration Failure Handling:** Confirmed deployment scripts halt immediately on migration failure.

---

## Phase L6 & L6.1 · Interactive Exercise Types & Prose

### 1. Decision & Rationale
* **The Problem:** Relying solely on multiple-choice questions induces cognitive fatigue and encourages passive elimination guessing.
* **The Decision:** Support diverse, active exercise types: `true_false` for rapid concept validation, `fill_blank` for syntax/keyword precision, `arrange` for architectural or algorithmic ordering, and `prose` (L6.1) for deep rubric-graded written explanations.
* **Constraint:** All exercise types must be evaluated deterministically on the server; the client never evaluates correctness or holds answers.

### 2. Plan
* Build dedicated client step renderers for `true_false`, `fill_blank`, and `arrange` with pure tap and keyboard interactions.
* Implement `prose` exercise support: text area with minimum character validation, Ctrl/Cmd+Enter check, and rich feedback trays.
* Standardize draft-to-wire answer serialization in `web/src/lib/lesson/answers.ts`.

### 3. Implementation
* **True/False (`TrueFalseStep.tsx`):** Two prominent cards with keyboard shortcuts (1/2 or T/F).
* **Fill-in-the-Blank (`FillBlankStep.tsx`):** Chip bank that slots into inline sentence or code blanks; tapping a placed chip returns it to the bank.
* **Arrange (`ArrangeStep.tsx`):** Tap-to-order sequence tiles with numbered position slots.
* **Prose (`ProseStep.tsx`):** Text area with live character counter ("X more characters to go" → "✓ Long enough to check"), keyboard submission, and feedback tray displaying rubric score, "You covered", "Worth adding", and "A strong answer".
* **Grader:** Pure Go evaluators handle whitespace trimming, casing, and position matching; rubric evaluator scores keyword coverage and structural requirements without LLM latency.

### 4. Testing & Verification
* **Unit Tests:** 68 unit tests in `answers.test.ts` verifying draft states, chip placement, removals, and wire payload formatting.
* **Prose Test:** Live full-stack test with test-only prose lesson: short submission rejected (<100 chars); partial answer scored 33; strong answer scored 100.
* **Zero Leakage:** Confirmed rubrics and correct answers are never sent to the browser before submission.
* **Real Browser:** Verified tap-only interactions on mobile viewports; zero console errors.

---

## Phase L7 · Retention & Player Polish

### 1. Decision & Rationale
* **The Problem:** Abrupt lesson exits cause frustration, silent navigation drops progress unexpectedly, and silent interfaces feel sterile.
* **The Decision:** Implement a mistakes review summary, synthesized WebAudio sound cues (opt-in), and defensive navigation guards.

### 2. Plan
* Add end-of-lesson mistakes review card to the summary screen.
* Synthesize audio cues dynamically via Web Audio API (zero audio file downloads).
* Provide sound mute toggle in player header, remembered per device.
* Add accessibility guards: background `inert` during exit dialogs and browser Back confirmation.

### 3. Implementation
* **Mistakes Review:** Summary card records all exercises missed on the first attempt during the session, displaying prompt, correct answer, and explanation.
* **Web Audio (`web/src/lib/sound.ts`):** Synthesized frequencies via `AudioContext` for correct chime, miss cue, and lesson completion fanfare.
* **Sound Toggle:** Header button persisting mute preference in `localStorage` (default muted).
* **Defensive Navigation:** Intercepts `popstate` and navigation links to prompt "Leave lesson? Progress will be lost" during active sessions.
* **Accessibility:** Added `inert` attribute to player background when leave modal is open; trapped focus within modal.

### 4. Testing & Verification
* **Audio Tests:** Unit tests in `sound.test.ts` verifying tone generation parameters and muted state handling.
* **Accessibility:** Confirmed Tab key cannot escape open dialogs and background shortcuts are ignored.
* **Browser Navigation:** Confirmed clicking browser Back button opens confirmation prompt rather than dropping session.

---

## Phase L8 & L8.1 · Working Professionals Path Content (Worlds 1–5, 40 Lessons)

### 1. Decision & Rationale
* **The Problem:** The platform required real, production-grade curriculum tailored to experienced software engineers, not introductory tutorials.
* **The Decision:** Author 5 complete worlds with 40 lessons and 175 exercises across all 4 topics. Content authored strictly in YAML, validated via CI, and checked against anti-bias standards (balanced option lengths, varied answer positions, zero single-option tells).

### 2. Plan
* Establish content curriculum in `content/PLAN.md` covering all topic skills, chained nodes, and boss encounters.
* Author World 1: First Ascent (System Design — 9 nodes)
* Author World 2: The Engine Room (Backend & Production — 8 nodes)
* Author World 3: Blueprint Workshop (Low-Level Design — 8 nodes)
* Author World 4: Algorithm Gardens (DSA Refresher Part 1 — 8 nodes)
* Author World 5: Deep Roots (DSA Refresher Part 2, L8.1 — 7 nodes)

### 3. Implementation
* **World 1 (System Design):** Caching, load balancing, sharding, replication, CAP theorem, boss.
* **World 2 (Backend & Production):** Concurrency, connection pools, idempotency, rate limiting, retries, metrics, boss.
* **World 3 (Low-Level Design):** SOLID, Factory, Strategy, Observer, State, Decorator, boss.
* **World 4 (DSA Part 1):** Two pointers, binary search, sliding window, fast/slow pointers, intervals, boss.
* **World 5 (DSA Part 2, L8.1):** Strings, linked lists, heaps/priority queues, greedy, memoization, DP table, boss.
* **Coverage:** 100% of skills covered across all 4 topics (System Design 4/4, Backend 4/4, LLD 3/3, DSA 14/14).

### 4. Testing & Verification
* **CI Validator:** `content-validator` passed with 0 errors across all 40 lesson files.
* **Live Full-Stack Playthrough:** Scripted agent played through all 40 lessons and 175 exercises against running services: all answers verified, unlocks fired as designed.
* **Audit Metrics:** Option positions evenly distributed (A: 25%, B: 25%, C: 25%, D: 25%); True/False balanced 50/50; correct option is longest in <25% of exercises.

---

## Phase L9 & L9.1 · Weekly Leagues, Leaderboards & Dashboard Card

### 1. Decision & Rationale
* **The Problem:** Learning platforms require a weekly social retention mechanic that encourages consistency without toxic anxiety.
* **The Decision:** Implement weekly leagues (Bronze, Silver, Gold, Platinum, Diamond) with cohorts of 30 learners. XP earned exclusively from lessons counts toward weekly standing. Demotion is gentle (safe zones, no demotion from lowest tier, skipped weeks maintain tier). Results derived lazily on user access at week boundaries—no heavy cron jobs.

### 2. Plan
* Add `league_cohorts` and `league_memberships` in Migration 000041.
* Progress service manages cohort placement and standings.
* Gateway aggregates leaderboard with public user cards (username + companion).
* Build `/league` screen with tier emblem, countdown timer, standings, and promotion/demotion zones.
* Add dashboard league card (L9.1) displaying rank, zone, weekly XP, and time remaining.

### 3. Implementation
* **Database:** Created `migrations/000041_create_leagues.up.sql`.
* **Progress Service:** `services/progress/internal/store/leagues.go` assigns learners to cohorts upon first weekly XP earn; updates league XP transactionally with lesson rewards.
* **Gateway Service:** `services/gateway/internal/leagues/service.go` enriches leaderboard with public user info (`POST /users/public-cards`).
* **Web Client:** Created `/league/page.tsx`, `TierEmblem.tsx`, and `LeagueCard.tsx` on the dashboard.

### 4. Testing & Verification
* **Rule Table Tests (`leagues_test.go`):** Cohort overflow at 30, tie-breaking by time achieved, promotion/demotion logic at week end, and skipped week retention.
* **Gateway Aggregation Tests:** Verified graceful fail-closed behavior when upstream user service is unreachable.
* **Frontend Tests:** Unit tests in `leagues.test.ts` verifying countdown timers and zone formatting.

---

## Phase L10 · Public Landing Page

### 1. Decision & Rationale
* **The Problem:** Prepio lacked a public front door explaining its value proposition to prospective engineers.
* **The Decision:** Build a high-converting public landing page at `/` adopting Duolingo's structural clarity (one promise, one obvious call-to-action) combined with Zeptap-inspired design aesthetics (floating pills, two-tone headlines, illustrated phone screens, dark theme).

### 2. Implementation & Testing
* **Components:** Built `Landing.tsx`, `PhoneScreens.tsx`, and `SessionRedirect.tsx` under `web/src/components/landing/`.
* **Session Redirect:** Authenticated visitors landing on `/` are automatically redirected to `/dashboard`.
* **Testing:** Tested responsive rendering across 390px, 768px, and 1280px screen widths; verified smooth reduced-motion behaviors.

---

## Phase L11 · Focus Topics Order the Journey

### 1. Decision & Rationale
* **The Problem:** Onboarding promised that selecting focus topics would customize the user's path, but worlds were hardcoded in fixed order.
* **The Decision:** Worlds must declare their topic (`worlds.topic_id` via Migration 000042). When querying `/path?focus=...`, the journey dynamically reorders worlds so the user's focus topics appear first. Unlocks still depend on finishing prerequisites—nothing is locked or hidden by focus.

### 2. Implementation & Testing
* **Database:** Migration `000042_world_topics.up.sql` added nullable `topic_id` to `worlds`.
* **Question Service:** `services/question/internal/service/lesson.go` orders worlds by user focus parameters before computing active node status.
* **Dashboard & Journey:** Dashboard passes user focus topics to `/path`; Journey UI displays "Your focus" badges on prioritized worlds.
* **Testing:** Integration tests verified ordering logic across single focus, multi-focus, and unknown topic parameters.

---

## Phase L12 · Achievements Collection

### 1. Decision & Rationale
* **The Problem:** Long-term milestones were missing; XP alone does not provide lasting collectible gratification.
* **The Decision:** Implement an achievements system owned by the Progress service, storing earned badges in `user_achievements` (Migration 000043). Achievements are awarded in the exact same transaction as lesson XP, ensuring zero duplicate or dropped unlocks.

### 2. Implementation & Testing
* **Configuration:** Defined 8 core achievements in `config/achievements.go` (First Steps, Warming Up, Getting Serious, Full Week, Level Five, Halfway There, Moving Up, Path Finder).
* **Progress Service:** Checks eligibility and inserts into `user_achievements` during lesson completion and streak processing.
* **Web Client:** Profile displays achievements grid (earned in gold with unlock date, locked in grey with clear requirements); celebration screen shows unlock banner.
* **Testing:** Verified with rule unit tests and live full-stack playthrough.

---

## Phase L13 · Interactive Code Snippet Exercises (Boot.dev Style)

### 1. Decision & Rationale
* **The Problem:** Learners could answer multiple-choice, ordering, fill-in-the-blank, and written questions, but lacked short, hands-on code writing / bug-fixing exercises where they directly edit snippets and see instant per-test pass/fail results.
* **The Decision:** Introduce `code_snippet` exercise type inspired by Boot.dev. Learners write code in an in-browser editor with line numbers, tab indentation, and language badges. Grading is 100% deterministic, server-side, and sandboxed without heavy Docker containers—using structured string, pattern (regex), and structural assertions authored in YAML. The server returns per-test results (`passed`, `name`), while tests and solution explanations are never leaked in client payloads.

### 2. Implementation
* **Database Migration:** Added `migrations/000044_allow_code_snippet_steps.up.sql` and `down.sql` to expand `lesson_steps.type` CHECK constraint to include `code_snippet`.
* **Domain Model & Validator:** Updated `services/question/internal/lesson/` with `StepCodeSnippet`, `CodeSnippetPayload`, and `validateCodeSnippet()` ensuring supported languages (`go`, `python`, `sql`, `typescript`), 2–8 test assertions (`contains`, `not_contains`, `pattern`, `structure`), and maximum 60 starter lines.
* **Grading Engine:** Added `GradeCodeSnippet()` in `services/question/internal/grading/grading.go` running tests sequentially and returning `test_results` (`name`, `passed`, `detail`).
* **Web Client:**
  - `CodeSnippetStep.tsx`: Code editor with monospaced font, gutter line numbers, 2-space tab indent, Enter key auto-indent, reset to starter button, and reactive focus/error borders.
  - `FeedbackTray.tsx`: Boot.dev-style test suite card displaying each test case with pass/fail badges (`✓` / `✕`) and names.
  - `answers.ts` & `types.ts`: Extended wire drafts, readiness checks (requires changing starter code), and unit tests (`vitest`).
* **Content Authoring:** Added real `code_snippet` exercises across 4 lessons covering 3 languages:
  - System Design (`why-caches-exist.yaml`): Cache-Aside read pattern in Go (`db.GetUser` miss handling + `cache.Set`).
  - Backend & Production (`retry-without-fear.yaml`): Payment idempotency store handler in Go.
  - Low-Level Design (`when-things-happen.yaml`): Observer pattern `subscribe` and `notify` in TypeScript.
  - DSA Refresher (`halving-the-search.yaml`): Lower-bound binary search integer division and `lo = mid + 1` advance in Python.

### 3. Testing & Verification
* **Backend Unit & Integration Tests:** `grading_test.go` and `validate_test.go` passed; integration test `TestAuthoredContent` and migration rollback tests passed cleanly.
* **Content Validation:** `content-sync --validate` confirmed 5 worlds and 40 lessons valid.
* **Frontend Vitest:** 73/73 tests passed in `web/src/lib/lesson/answers.test.ts` and related modules.
* **Frontend Build:** `npm run build` compiled with Turbopack with 0 errors across all 12 routes.

---

## Operational Verification & Build Status

* **Backend Compilation:** `go test ./...` passes across all microservices, integration tests, and tools.
* **Frontend Compilation:** Next.js production build (`npm run build` with Turbopack) passes with exit code 0 (all 12 static/dynamic routes generate cleanly).
* **Database Migrations:** All migrations 000001 through 000044 apply up cleanly; migrations 000035 through 000044 roll back and reapply with zero errors.
* **Self-Hosted Web Fonts:** Sora, Manrope, and IBM Plex Mono fonts are self-hosted via `web/src/fonts/` (WOFF2, ~104 KB total), removing any runtime or build-time dependency on Google Fonts.
* **Working Tree:** Branch `feat/l13-code-snippet`.

