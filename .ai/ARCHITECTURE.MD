# ARCHITECTURE.md

## AGENT INSTRUCTIONS

Before any code change, read in this order: PRODUCT.MD, ARCHITECTURE.MD, CONTENT_SYSTEM.MD, EXECUTION.MD.

If a proposed implementation violates ownership boundaries or any rule below: **STOP. Do not implement. Explain the violation.**

---

## Purpose

Defines the system architecture, ownership boundaries, domain model, engineering standards, and decision rules, so Prepio stays maintainable as it grows. Architecture exists to support the product, not the other way around.

## Philosophy

Optimize for maintainability, extensibility, content growth, and product evolution. Do not optimize for premature scale, premature microservices, or theoretical performance. Build the simplest version of the correct architecture.

---

## Core Domain Model

```
World → Node → Lesson → Step
                 │
                 └── maps (weighted) to ──▶ Skill

Step result → Skill mastery → Topic readiness
```

This hierarchy is never bypassed.

| Concept | Definition |
|---|---|
| **World** | A progression destination. Contains nodes. Never contains lessons directly. |
| **Node** | A milestone on the path. Binds exactly one lesson. Does not own step content. |
| **Lesson** | The unit of learning: an ordered list of steps (≤ 5 min), mapped to one or more skills. A *boss* lesson is an ordinary authored lesson with `kind: boss` (longer, mixed review). There is no dynamic question selection. |
| **Step** | One screen in a lesson: `intro` or an exercise (`mcq`, `true_false`, `fill_blank`, `arrange`, `prose`). |
| **Skill** | The unit of mastery (0–100). Lessons teach and measure skills. |
| **Topic** | A user-facing grouping of skills (a skill category). Readiness is computed per topic. |
| **Attempt** | One user run through a lesson. Holds step results. Owned by the Journey domain. |

Skills drive readiness. Lessons exist to build and measure skills.

---

## Ownership Model

Ownership is absolute. **No service may mutate data owned by another service.**

| Domain | Owns | Hosted in |
|---|---|---|
| **User** | Authentication, profile, preferences (experience level, focus topics), account settings | `user` service |
| **Companion** | Companions, evolution, dialogue, cosmetics | `user` service (today) |
| **Content** | Skills, subskills, topics, lessons, steps and their payloads, authoring validation | `question` service |
| **Journey** | Worlds, nodes, unlock rules, journey progress, **lesson attempts and step results**, grading | `question` service |
| **Progress** | XP, levels, gems, skill mastery + mastery ledger, topic readiness, achievements, quest progress | `progress` service |
| **Streak** | Streak state | `streak` service |
| **Notification** | Dispatch | `notification` service |
| **Analytics** | Events, reports, funnels (not built yet). Never owns business state. | n/a |

> **Known naming debt:** the `question` service now hosts the Content and Journey domains. Renaming it is deferred until it is cheap; do not create a new service for this.

If ownership is unclear: do not implement. Ask.

---

## Data Flow

```
Answer a step ──▶ server grades (Journey) ──▶ step result stored
Attempt completes (all graded steps correct) ──▶ node marked done (Journey)
                                             └─▶ `lesson.completed` fact
`lesson.completed` ──▶ Progress: skill mastery + ledger, XP, gems
                   ├─▶ Streak: streak update
                   └─▶ (Progress emits) skill.progress.updated, readiness.updated, level.up
```

Progress flows downward, never backwards. The gateway may aggregate responses for the client (as it does for the dashboard); aggregation does not change ownership.

## Event-Driven Principles

Events communicate **facts**, not actions, and contain no business logic. XP and mastery calculation belong to Progress, never to the publisher.

Good: `lesson.completed`. Bad: `lesson.completed_and_grant_xp`.

**Approved event types:** `lesson.completed`, `skill.progress.updated`, `readiness.updated`, `journey.node.completed`, `world.completed`, `achievement.unlocked`, `quest.completed`, `companion.evolved`, `level.up`.

The `lesson.completed` payload carries performance facts only (user, lesson, attempt, per-skill weights, first-attempt correctness, step counts, timestamp).

---

## Grading Rules (the server is the only grader)

1. **The client never receives answers or rubrics.** Lesson payloads sent to clients have all answer fields stripped. This is enforced by a test.
2. Deterministic step types (`mcq`, `true_false`, `fill_blank`, `arrange`) are graded by **pure functions in one Go package**, table-tested. `prose` is graded by the rubric evaluator (structural stages now; LLM stage when wired). The LLM is never used for deterministic types.
3. A grading response contains: correctness, the explanation, and **only after an incorrect answer**, the correct answer.
4. **Attempt integrity:**
   * One in-progress attempt per user and lesson; resumable.
   * Step answers are idempotent per (attempt, step, try number).
   * An attempt completes only when **every graded step has a correct result**.
   * Mastery uses **first-attempt correctness**. XP and mastery deltas are capped server-side.
5. **Re-queueing missed steps is UI sequencing** (presentation order is frontend state). Correctness and completion are server decisions.
6. Latency budget: step-answer p95 < 150 ms server-side. The UI shows the pressed state immediately and never blocks on animation.

## Lesson Runtime Capabilities

APIs expose capabilities, not tables:

* `GET /journey`: path with node status plus preview fields (lesson title, takeaways, minutes, XP)
* `POST /lessons/{id}/attempts`: start or resume; returns steps **without answers**
* `POST /attempts/{id}/steps/{stepId}/answer`: server-graded feedback
* `POST /attempts/{id}/complete`: performance summary (the gateway adds rewards from Progress)
* `GET /progress/...`: topic mastery, explained mastery changes

---

## Frontend Architecture

UI is a consumer, not a source of truth.

The frontend **may**: render, animate, collect input, show optimistic pressed states, and sequence steps (including re-queue ordering).
The frontend **must not**: grade answers, hold answers, or calculate mastery, readiness, XP, or rewards. Never duplicate business logic in the client.

Frontend state is UI, screen, and temporary interaction state only. Progress, mastery, journey, rewards, and achievements are backend state.

## Content Configuration Rule

Content is configurable, not hardcoded in application code.

* Worlds, nodes, and lessons are **content-as-code**: files in `content/`, validated in CI, and loaded into the database by an **idempotent `content-sync` command that runs after migrations on every deploy**.
* The database is the runtime source of truth. The repo files are the authoring source of truth.
* Sync **never deletes user data**; removed content is deprecated, not dropped.
* The skill catalog (structure) changes through migrations; it changes rarely.
* No lesson text, options, or answers appear in application code. Test fixtures live only under test directories.

## Database Principles

Tables represent domains, not screens. (Good: `lessons`, `lesson_attempts`, `skills`. Bad: `dashboard_cards`.) Tables model reality, not UI.

## Migration Rules

Every schema change requires: a migration, a rollback, an index review, and a documentation update.
**Historical migrations are never edited.** Removals are new migrations. Never alter production schemas manually.

## Engineering Standards

* **Production quality only:** no demo shortcuts, hardcoded content, mock data, or stubbed shipped paths.
* **Dead code is removed in the phase that makes it dead:** unused routes, env vars, flags, tests, components, tables. Feature flags are deleted when their phase ends.
* Every phase ends with the sweep and gates defined in EXECUTION.MD.
* Prefer extending existing domains; create a new domain only when ownership genuinely changes.

---

## Anti-Patterns (forbidden)

* Hardcoded progression or unlock rules (use configurable rules)
* Step or lesson results driving world progress directly (use Result → Mastery → Progression)
* Client-side grading, or answers present in client payloads
* UI-driven logic (frontend calculating mastery, XP, readiness)
* Duplicate ownership (Progress updating companions; emit an event, Companion reacts)
* A new service per feature
* Company-specific data or logic anywhere
* Content in application code

## Decision Tree

* Need skill or lesson content → Content domain
* Need journey state, attempts, or grading → Journey domain
* Need mastery, readiness, XP → Progress domain
* Need companion data → Companion domain
* Need profile or preferences → User domain
* Unsure → STOP. Do not create new ownership.

## Success Criteria

The architecture supports new lessons, worlds, skills, topics, and companions **without code changes or rewrites**. The goal is not perfection; it is avoiding future rewrites.
