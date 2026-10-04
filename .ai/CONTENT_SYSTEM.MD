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

* **Topic** = a user-facing skill category. Initial topics: **System Design, Backend & Production, Low-Level Design, DSA Refresher.** (Catalog changes needed to express these, such as a Backend & Production category and grouping data structures + algorithms under DSA Refresher, are made by migration in Phase L2.)
* **Skill** = the unit of mastery, 0–100. **Subskills** refine skills (e.g. Caching under Scaling & Performance).
* Every lesson maps to ≥ 1 skill with weights summing to 1.0.

## Mastery

* Each user has a mastery value (0–100) per skill.
* A completed lesson moves mastery by a **gradual, capped, smoothed** amount derived from first-attempt accuracy, difficulty, and the lesson's skill weights. Mastery never jumps dramatically. Parameters live in config, not code.
* **Every mastery change is written to a ledger** (user, skill, delta, lesson attempt, time). This is what makes readiness explainable.

## Readiness (per topic)

* **Topic readiness = weighted mean of skill mastery over the skills in that topic.** Skills not yet attempted show as "not started", not as failures.
* Readiness answers: *how strong is this user in this topic, and why did it change?*
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
| `mcq` | `prompt`, `options[]` (3–4), `answer` (index), `explanations.correct`, `explanations.wrong[index]` | deterministic |
| `true_false` | `statement`, `answer` (bool), `explanation` | deterministic |
| `fill_blank` | `code` with numbered blanks, `bank[]` (includes distractors), `answers[]` (per blank), `explanation` | deterministic |
| `arrange` | `prompt`, `items[]`, `order[]` (correct sequence as indices), `explanation` | deterministic |
| `prose` | `prompt`, `rubric` (concept list with aliases and required flags), `minChars` | rubric evaluator (LLM stage later) |

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
* `content-sync --validate` runs in CI. `content-sync` upserts by stable slug after migrations on every deploy, is idempotent, and never deletes user data (it deprecates).
* Changing a published lesson increments its `version`; attempts record the version they ran against.

### Validation rules (CI blocks on failure)

* All references resolve (node, skills). Slugs and step ids are unique.
* Every lesson maps to ≥ 1 skill; weights sum to 1.0.
* ≥ 1 graded step; **every graded step has an explanation**.
* `mcq`: 3–4 options, exactly one correct, an explanation for every wrong option.
* `true_false`: explanation present.
* `fill_blank`: blank count = `answers` count; `bank` contains every answer.
* `arrange`: ≥ 3 items; `order` is a permutation of the item indices.
* `prose`: rubric has ≥ 1 required concept.
* `est_minutes` within limits for the lesson `kind`.

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
