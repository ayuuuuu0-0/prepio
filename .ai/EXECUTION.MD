# EXECUTION.md

# AGENT INSTRUCTIONS

This document is the current source of truth for implementation priorities.

Before implementing any feature:

Read:

1. PRODUCT.md
2. ARCHITECTURE.md
3. CONTENT_SYSTEM.md
4. EXECUTION.md

If a feature is not listed in the current phase:

Do not build it.

If uncertain:

Ask for clarification.

---

# Purpose

This document controls execution.

It defines:

* Current phase
* Current priorities
* Build order
* Out-of-scope work
* Definition of done
* Launch criteria

This is the only document expected to change frequently.

---

# Current Project Status

The product has undergone a major redesign.

The old version proved that:

* Authentication works
* Questions work
* Progression works
* Journey concept works

The redesign identified weaknesses in:

* Content architecture
* Readiness architecture
* Progression architecture
* Product identity

The goal is not adding features.

The goal is building the correct foundation.

---

# CURRENT PHASE

## Phase A — Foundation Rebuild

Status:

ACTIVE

Priority:

CRITICAL

No major feature work should begin until this phase is complete.

---

# Phase A Goal

Build the foundational systems correctly.

Avoid temporary solutions.

Avoid future rewrites.

Build systems that can support:

* Thousands of questions
* Hundreds of skills
* Many companies
* Many worlds

without architectural changes.

---

# PHASE A TASKS

---

## A1 Skill Graph

Priority:

P0

Goal:

Create the learning foundation.

---

Required:

skills

subskills

question_skills

skill_categories

---

Example:

Arrays

├ Two Pointers

├ Sliding Window

├ Prefix Sum

└ Hash Maps

---

Completion Criteria:

Every question maps to skills.

Every skill maps to a category.

---

## A2 Content Architecture

Priority:

P0

Goal:

Decouple progression from questions.

---

Current:

Question → Progress

---

Target:

World

↓

Node

↓

Skill

↓

Question Pool

↓

Question

---

Completion Criteria:

Nodes reference skills.

Not questions.

---

## A3 Question Schema Upgrade

Priority:

P0

Goal:

Create scalable question structure.

---

Every question must support:

Skill

Subskill

Difficulty

Companies

Readiness Weight

Estimated Time

Explanation

Hints

Status

Evaluation Type

---

Completion Criteria:

Question model supports future growth.

---

## A4 Readiness Engine V2

Priority:

P0

Goal:

Make readiness meaningful.

---

Current readiness is insufficient.

---

Target:

Question

↓

Skill Score

↓

Company Readiness

↓

Overall Readiness

---

Completion Criteria:

Every readiness score is explainable.

---

## A5 Content Management Foundation

Priority:

P1

Goal:

Remove dependency on migrations for content.

---

Build internal content management capabilities.

---

Must Support:

Create Skills

Create Questions

Create Pools

Create Nodes

Assign Companies

Assign Difficulty

---

Completion Criteria:

Content can be managed without code changes.

---

# PHASE A SUCCESS CRITERIA

Phase A is complete when:

✓ Skill graph exists

✓ Question schema upgraded

✓ Readiness V2 works

✓ Content architecture works

✓ Content management foundation exists

Only then proceed.

---

# NEXT PHASE

## Phase B — Content Expansion

Status:

LOCKED

Do not begin until Phase A is complete.

---

Goal:

Create enough content to validate the product.

Not enough content to scale.

---

Target:

100–150 high-quality questions

---

Coverage:

Programming Fundamentals

Arrays

Strings

Hash Maps

Trees

Graphs

DP

LLD

System Design

Behavioral

---

Important:

Quality over quantity.

Depth over breadth.

---

# PHASE C

Status:

LOCKED

---

Goal:

Create retention.

Build emotional systems.

---

Features:

Companion Evolution

Session Summary

Daily Quests

Weekly Quests

World Expansion

Achievement System

---

Reason:

Retention should sit on top of progression.

Not replace it.

---

# PHASE D

Status:

LOCKED

---

Goal:

Launch readiness.

---

Features:

Analytics

Reliability

Performance

UX Polish

Beta Testing

---

Target:

20–50 real users

---

Questions To Answer:

Do users understand readiness?

Do users understand progression?

Do companions matter?

Do worlds matter?

Do users return?

---

# OUT OF SCOPE

The following are explicitly forbidden before launch.

---

Arena

Multiplayer

Guilds

Voice Interviews

AI Interviewers

AI Resume Review

Recruiter Marketplace

Season Pass

Advanced Recommendation Engine

Community Content

Complex Social Features

---

Reason:

These features solve scale.

Not validation.

---

# DEFINITION OF DONE

A task is NOT complete because:

* Code compiles
* UI renders
* API responds

A task is complete when:

Architecture respected

Content model respected

Documentation updated

Tests pass

UI matches PRODUCT.md

No hardcoded assumptions

No ownership violations

---

# AGENT DECISION RULES

When making decisions:

1. Prefer configurable systems.

2. Prefer content-driven systems.

3. Avoid hardcoding.

4. Avoid shortcuts requiring rewrites.

5. Avoid speculative scale work.

6. Avoid creating new services.

7. Follow existing ownership.

---

# IF YOU DISCOVER A BETTER DESIGN

Do not implement immediately.

Instead:

Document:

Problem

Proposed Solution

Benefits

Risks

Migration Impact

Then request approval.

---

# LAUNCH CRITERIA

The product is ready for launch when:

Users can progress through worlds.

Users understand readiness.

Users understand skills.

Users feel progression.

Users receive meaningful feedback.

Users return voluntarily.

Content supports at least 30 days of usage.

The goal is not feature completeness.

The goal is validation.

---

# FINAL RULE

Do not optimize for:

100,000 users.

Optimize for:

20 users who genuinely love the product.

If 20 users love it:

Scaling can be solved later.

If 20 users do not love it:

Nothing else matters.
