# CONTENT_SYSTEM.md

## AGENT INSTRUCTIONS

This document defines how learning, progression, mastery, and content work inside Prepio.

Before creating or changing lessons, skills, worlds, nodes, mastery logic, or grading logic: read this document completely, then PRODUCT.MD and ARCHITECTURE.MD. Never bypass the structures defined here.

---

## Purpose

The content system is the core intelligence of Prepio. UI and backend can change; this should stay stable.

## Core Principle

**Skills are the product. Lessons teach and measure them.** Users do not need a thousand exercises; they need mastery of skills. Depth beats volume.

## Learning Hierarchy

```
World → Node → Lesson → Step
                 └─ maps (weighted) to Skills
Step result → Skill mastery → Topic readiness
```

Never bypass this structure. Lesson completion does not unlock worlds directly: it moves skill mastery and node progress, and unlock rules (configurable) decide the rest.

## Paths, Worlds, Nodes

* A **path** is an ordered set of worlds. Phase L ships one path: **Working Professionals**.
* Worlds are destinations, not categories ("System Design" is a topic; a world has an identity).
* A node is one milestone and binds exactly one lesson.
* The world/node outline for the Working Professionals path is **designed jointly with the product owner and lives in `content/`**, not in these docs.
* There are no company worlds and no company content.

## Topics, Skills, Subskills

* **Topic** = a user-facing skill category. Initial topics: **System Design, Backend & Production, Low-Level Design, DSA Refresher.** Topics are rows in `topics`; each skill category belongs to at most one topic (`skill_categories.topic_id`). System Design, Low-Level Design, and Backend & Production map one category each; DSA Refresher groups Data Structures and Algorithms. Categories with no topic (e.g. Behavioral) do not contribute to readiness. Backend & Production starts with four skills: API Design, Databases & Transactions, Concurrency & Async, Observability & Operations.
* **Skill** = the unit of mastery, 0–100. **Subskills** refine skills (e.g. Caching under Scaling & Performance).
* Every lesson maps to ≥ 1 skill with weights summing to 1.0.

## Mastery

* Each user has a mastery value (0–100) per skill.
* A completed lesson moves mastery by a **gradual, capped, smoothed** amount derived from first-attempt accuracy, difficulty, and the lesson's skill weights. Mastery never jumps dramatically. Parameters live in config (`config/lessons.go`), not code.
* Per skill: `delta = round(gain × first-attempt accuracy × skill weight × difficulty multiplier × headroom)`, where `headroom = (100 − current) / 100` (so gains shrink near 100), boss lessons scale gain up, and a hard cap limits any one lesson's change to a skill. Mastery and XP are granted only on a user's **first completion** of a lesson; replays grant none.
* **Every mastery change is written to a ledger** (user, skill, delta, lesson attempt, time). This is what makes readiness explainable.

## Readiness (per topic)

* **Topic readiness = weighted mean of skill mastery over the skills in that topic.** Skills not yet attempted show as "not started", not as failures.
* Readiness answers: *how strong is this user in this topic, and why did it change?*
* The formula, exactly: a topic's readiness is the mean mastery of the skills the learner has **started** (attempted at least once), every started skill weighted equally, rounded to a whole number. A topic with nothing started has no number ("Not started"), never zero, and the dashboard also says how many of its skills are started. Computed by Progress (`GET /progress/topics`) from `user_skill_scores`; the client only displays it.
* Transparency is required. Bad: "System Design +2". Good: "System Design +2 · Caching +1 · Load Balancing +1".
* Never derived from XP, levels, or streaks. Never company-specific.

---

## Lesson Model

| Field | Rules |
|---|---|
| `slug` | Stable unique id; never reused |
| `title`, `summary[]` | Summary = 2–4 takeaway bullets (shown in preview and summary card) |
| `node` | Slug of the node that binds this lesson |
| `kind` | `lesson` or `boss` |
| `difficulty` | `easy`, `medium`, `hard` |
| `est_minutes` | Target ≤ 5 (boss ≤ 8) |
| `skills[]` | `{skill, weight}`; ≥ 1; weights sum to 1.0 |
| `status` | `draft` → `review` → `published` → `deprecated`. Only `published` is served. |
| `steps[]` | Ordered; exactly one `intro` first (optional), then 3–5 graded exercises (boss: up to 8) |

### Step types

| `type` | Payload | Graded by |
|---|---|---|
| `intro` | `beats[]` (`text`, `emphasis?`, `visual?`, `durationMs`), optional `mediaUrl` | not graded (marked viewed) |

Intro motion: each beat plays for its `durationMs` (1.5-12 s; the whole intro at most 45 s): its words reveal, the `emphasis` phrase glows, and the next beat follows on its own. Tapping advances, "Skip intro" always works, it can be paused, and at the end it can be replayed (also from the player header). With reduced motion nothing auto-advances. `visual` is one of a fixed set of named illustrations the client draws (`bars`, `flow`, `pulse`); the validator rejects any other name, and adding one means adding it to both the web client and the validator. `mediaUrl` is optional, must be an `https://` URL or a site-relative path, and must end in a known image or video type (mp4, webm, png, jpg, gif, webp, avif, svg); it is decorative, so it must never carry information the text does not.
| `mcq` | `prompt`, `options[]` (3–4), `answer` (index), `explanations.correct`, `explanations.wrong[index]` | deterministic |
| `true_false` | `statement`, `answer` (bool), `explanation` | deterministic |
| `fill_blank` | `code` with blanks written `___1___`, `___2___`, … (each exactly once), `bank[]` (includes distractors), `answers[]` (per blank), `explanation` | deterministic |
| `arrange` | `prompt`, `items[]` (authored scrambled; the client receives them as written), `order[]` (correct sequence as indices, never the identity), `explanation` | deterministic |
| `prose` | `prompt`, `rubric` (concept list with aliases and required flags), `minChars` (≥ 100), `explanation` (the worked answer shown after grading) | rubric evaluator (LLM stage later) |

Later (not before its phase): `type_line`, a short free-text answer matched against accepted variants after whitespace/quote normalisation. No code-execution sandbox.

### Authoring format

```yaml
slug: <lesson-slug>
title: "<Title>"
node: <node-slug>
kind: lesson
status: published
difficulty: easy
est_minutes: 5
skills:
  - { skill: <existing-skill-slug>, weight: 1.0 }
summary:
  - "<takeaway>"
steps:
  - id: intro
    type: intro
    beats:
      - { text: "<beat>", emphasis: "<word>", durationMs: 4000 }
  - id: q1
    type: mcq
    prompt: "<question>"
    options: ["<a>", "<b>", "<c>"]
    answer: 1
    explanations:
      correct: "<why this works>"
      wrong:
        0: "<why not a>"
        2: "<why not c>"
```

## Content-as-Code

* `content/worlds/*.yaml` define worlds and their nodes. `content/lessons/**/*.yaml` define lessons.
* `content-sync --validate` validates without writing (structure always; skill references too when `DATABASE_URL` is set). The same checks run in `go test ./...` so CI blocks invalid content. `content-sync` upserts by stable slug after migrations on every deploy, is idempotent, and never deletes user data (it deprecates).
* Changing a lesson's content (anything but `status`) increments its `version`; its steps and skill weights are stored per version, attempts record the version they ran against, and old versions are kept so in-flight attempts finish on the content they started with.
* A world file is `slug`, `name`, `description`, `theme`, `order`, and `nodes[]`. A node is `slug` (unique across worlds), `label`, `type` (`lesson` or `boss`), and optional `requires[]` (node slugs that must be completed first; this is the configurable unlock rule). Nodes with no `requires` are open from the start.

### Validation rules (CI blocks on failure)

* All references resolve (node, skills). Slugs and step ids are unique.
* Every lesson maps to ≥ 1 skill; weights sum to 1.0.
* ≥ 1 graded step; **every graded step has an explanation**.
* `mcq`: 3–4 options, exactly one correct, an explanation for every wrong option.
* `true_false`: explanation present.
* `fill_blank`: blank count = `answers` count; `bank` contains every answer.
* `arrange`: ≥ 3 items; `order` is a permutation of the item indices.
* `prose`: rubric has ≥ 1 required concept; `minChars` ≥ 100.
* `est_minutes` within limits for the lesson `kind`.
* `summary` has 2–4 bullets; an `intro` is the first step and there is at most one; graded step count is 3–5 (boss 3–8).
* A lesson's `kind` matches its node's `type`; every node binds exactly one lesson; node `requires` resolve and never form a cycle; node and lesson slugs are lowercase kebab-case.
* `arrange`: `order` is not the identity permutation.
* `fill_blank`: blanks are numbered `1..n`, each appears once.

Lessons failing validation cannot be published.

## Lifecycle

Draft → Review → Published → Deprecated. Content never goes straight to production. AI-assisted drafting may be introduced later; human review stays required.

## Evaluation

* Deterministic types: pure Go functions, one package, table-tested.
* Prose: rubric evaluator (keyword filter + structural coverage; LLM stage when configured). The LLM never grades deterministic types.
* All grading is server-side (see ARCHITECTURE.MD).

## Expansion Rule

Before adding lessons, ask: does this skill already have enough coverage? If no, improve coverage. If yes, add content. Quality over quantity.

## Anti-Patterns (forbidden)

* Lessons without skills or explanations
* Wrong answers without an explanation of why
* Lesson-based unlocks that bypass mastery rules
* Hardcoded readiness or hardcoded unlocks
* Lesson count or XP as the primary metric
* Company tags, company readiness, company worlds
* Content in application code

## Long-Term Goal

The content system should be able to answer: What is this user good at? What are they weak at? Why is their readiness changing? What should they practice next? What skill is holding them back? If it cannot, the implementation is incomplete.
