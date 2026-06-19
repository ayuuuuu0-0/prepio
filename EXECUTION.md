# EXECUTION.md

# Purpose

This document defines:

* Architecture
* Ownership
* Build Order
* Engineering Rules

This is the source of truth for implementation.

---

# Build Philosophy

Avoid both extremes.

Bad:

Temporary hacks.

Bad:

Infinite architecture.

Build the simplest version of the correct system.

---

# Current Priority

Build foundations before features.

The next goal is not:

Arena

AI Interviewers

Guilds

Seasons

The next goal is:

Content Architecture

Skill Graph

Readiness Engine

Journey Architecture

---

# Architecture

## User

Owns:

Authentication

Profiles

Preferences

Target Companies

---

## Journey

Owns:

Worlds

Nodes

Unlocks

Journey Progress

---

## Content

Owns:

Skills

Subskills

Question Pools

Questions

Company Mapping

Readiness Weights

---

## Progress

Owns:

XP

Levels

Gems

Achievements

Quests

Readiness Scores

---

## Companion

Owns:

Companions

Evolution

Dialogue

Cosmetics

---

# Core Domain Model

World

↓

Node

↓

Skill

↓

Question Pool

↓

Question

↓

Skill Progress

↓

Readiness

This hierarchy must never be broken.

---

# Skill Graph

Required Before Expansion.

Example:

Arrays

├── Two Pointers

├── Sliding Window

├── Prefix Sum

└── Hash Maps

Questions must map to skills.

Skills drive readiness.

---

# Readiness Engine

Question

↓

Skill Score

↓

Company Readiness

↓

Overall Readiness

Never calculate readiness directly from XP.

XP measures activity.

Readiness measures capability.

---

# Content Rules

Do not hardcode progression.

Do not hardcode worlds.

Do not hardcode company mappings.

Use content-driven systems.

Questions should be configurable.

Worlds should be configurable.

Nodes should be configurable.

---

# Build Order

## Phase A

Foundation

* Skill Graph
* Content Architecture
* Readiness Engine V2
* Question Schema Upgrade
* Content Management Foundation

---

## Phase B

Content Expansion

* 100–150 quality questions
* Company mappings
* Skill coverage
* Readiness validation

---

## Phase C

Retention

* Companion Evolution
* Session Summary
* Quests
* World Expansion

---

## Phase D

Launch

* Analytics
* Reliability
* UX Polish
* Beta Release

20–50 users

---

# Forbidden Before Launch

Do not build:

Arena

Guilds

Multiplayer

AI Interviewers

Voice Interviews

Recruiter Features

Season Passes

Advanced Recommendations

These solve scale problems.

Not validation problems.

---

# Engineering Rules

Every feature must:

1. Fit PRODUCT.md
2. Support progression
3. Support readiness
4. Respect ownership boundaries
5. Avoid rewrites later

---

# Launch Criteria

Launch when:

Users understand readiness.

Users understand progression.

Users complete journeys.

Users return daily.

Users request more content.

Not when every feature is complete.

Not when every dream feature exists.

A product is validated by users.

Not by architecture.
