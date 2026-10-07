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
**Status: implemented** on branch `feat/l2-lesson-platform`, pending review. The first lesson (`why-caches-exist`, world `First Ascent`) uses only an intro and `mcq` steps so the L3 player, which starts with MCQ, can run it; world and node names are provisional until the outline is designed with the product owner (L8). Skills added: the four Backend & Production skills. The new lesson journey is served at `GET /path`; the legacy `GET /journey` stays until L3 removes it.
Schema (`lessons`, `lesson_steps`, `lesson_skills`, `lesson_attempts`, `lesson_step_results`, `mastery_ledger`); topic catalog migration; `content/` with the validator and `content-sync` (wired into the deploy entrypoint after migrations); grading package (pure, table-tested) and prose via the rubric evaluator; lesson runtime API; `lesson.completed` event; Progress consumer writes mastery + ledger + XP; Journey binds nodes to lessons and derives status from attempts; the first real lesson of the Working Professionals path authored in `content/`.
*Done when:* integration tests cover the attempt lifecycle, idempotency, invalid answers, and completion rules; a test proves client payloads contain no answers; `content-sync` is idempotent; gates pass.

**L3, Lesson player (web) and legacy removal**
**Status: implemented** on branch `feat/l3-lesson-player` (player, journey, legacy removal), pending review. Dropped with the legacy loop, beyond the list below: `question_skills`, `node_pools`, `node_skills`, and `user_journey_progress` (all reference `questions` or the removed pool/index builders, and nothing reads them). The skill catalog, `user_skill_scores`, XP, gems, and streaks are kept. The web exercise types beyond `mcq` arrive in L6; a lesson step the player cannot render shows a leave-lesson message rather than a placeholder.
Player shell (focus mode, top bar, progress, combo), feedback tray, MCQ step, miss re-queue, celebration, summary card, journey preview sheet, node-completion animation, unlock messaging; companion in the player. Then **remove the legacy loop**: `/challenge`, daily papers, the legacy question/answer flow and `question.answered` event, journey index/pool builders and the `JourneyPoolSelection` flag, `question_pools*`, `questions`, `daily_papers*`, `user_question_history` (new migration). Keep the rubric evaluator for `prose`.
*Done when:* a user completes a real lesson end to end from the journey; legacy grep is clean; screenshots reviewed; gates pass.

**L4, Onboarding and dashboard**
Topic picker (1–3) → experience → companion; `focus_topics` stored by the User domain. Dashboard: per-topic mastery rings, Continue → next lesson, explained mastery changes on the celebration screen. Remove "coming soon" quest/league cards and dead nav tabs until their phases.
*Done when:* new users pick topics and land on a dashboard that shows real topic mastery and the next lesson; no placeholders remain; gates pass.

**L5, Intro motion stage:** beat-based motion renderer, tap to advance, skip, replay, optional `mediaUrl`.

**L6, Exercise types:** `true_false`, `fill_blank` (tap-to-place), `arrange` (tap-to-order); validator rules and player support.

**L7, Retention polish:** mistakes review, replay intro, sound toggle, keyboard and accessibility audit.

**L8, Working Professionals path content:** worlds, nodes, and lessons authored with the product owner. Quality over quantity; first world complete before the second begins.

### Phase L Success Criteria

* A user can complete lessons end to end with instant, server-graded feedback.
* Topic mastery is real, explained, and visible on the dashboard.
* Content is authored as code, validated in CI, and synced on deploy.
* No company logic, legacy question loop, or dead flags remain.

---

## OUT OF SCOPE (forbidden before launch)

Arena · Multiplayer · Guilds · Voice interviews · AI interviewers · AI resume review · Recruiter marketplace · Season pass · Advanced recommendation engine · Community content · Complex social features · Company-specific prep · Hearts/lives · Code-execution sandbox · AI-generated lessons · Leagues and quests (until a later phase) · Mobile client updates · Additional audience paths · Marketing/landing page

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
