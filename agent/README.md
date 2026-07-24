# Agent Context

This file is the consolidated working context for the Prepio codebase. Read this before making product, architecture, content, readiness, or execution changes.

## Read order

1. This file
2. The code or migration files directly involved in the change
3. The smallest nearby tests or call sites that can verify the behavior

If a proposed change conflicts with the boundaries below, stop and explain the conflict before editing code.

## Current project state

Prepio is a progression game for interview preparation. The product is built around visible career progress rather than a raw question bank.

The foundation work already completed covers:

- A1 skill graph schema and mappings
- A2 content architecture and journey bindings
- A3 question schema expansion
- A4 readiness V2 foundation and validation
- Journey V2 pool selection behind a feature flag

## Product model

The product is not a generic interview prep tool. The platform exists to make career progression feel visible, explainable, and motivating.

Core product rules:

- Readiness is the primary metric
- Skills are the foundation, not questions
- Questions are content used to evaluate and improve skills
- Journey is the main experience, not a flat list of questions
- Companions, rewards, worlds, and collection exist to strengthen progression and retention

The experience should feel deliberate and engineered, not like an admin console or a template site.

## Architecture boundaries

Prepio uses strict domain ownership.

- User service owns authentication, profiles, preferences, target companies, and account settings
- Content domain owns skills, subskills, questions, question pools, company mapping, difficulty, and content metadata
- Journey domain owns worlds, nodes, unlock rules, and journey progression
- Progress domain owns readiness computation

Rules:

- Do not mutate data owned by another service
- Do not bypass the World -> Node -> Skill -> Question Pool -> Question hierarchy
- Questions should not directly own progression logic
- Journey nodes describe content bindings; they do not own the questions themselves
- Keep additive changes backward compatible unless a migration explicitly changes behavior

## Content system

The content hierarchy is:

World -> Node -> Skill -> Question Pool -> Question -> Skill Mastery -> Readiness

Important model notes:

- Worlds are progression destinations, not categories
- Nodes are milestones within a world
- Skills are the unit of mastery
- Question pools let the content bank grow without changing progression structure
- Readiness should be explainable from weighted skill mastery

## Execution model

The active work is still foundation-first.

Working rules:

- Do not build features outside the current phase unless explicitly requested
- Keep changes small and additive where possible
- Use feature flags for behavior changes such as readiness V2 and journey pool selection
- If a requested change conflicts with product, architecture, or content boundaries, stop and explain why

## Practical rules for edits

- Start from the smallest concrete file or symbol that owns the behavior
- Prefer the nearest test or call site that can disconfirm the current hypothesis
- After the first substantive edit, validate the touched slice before broadening scope
- Keep the root README factual and concise
- Keep this file as the single consolidated AI-facing reference for the repo
