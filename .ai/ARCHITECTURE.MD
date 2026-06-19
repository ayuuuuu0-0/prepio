# ARCHITECTURE.md

# AGENT INSTRUCTIONS

Before making any code changes:

Read:

1. PRODUCT.md
2. ARCHITECTURE.md
3. CONTENT_SYSTEM.md
4. EXECUTION.md

If a proposed implementation violates architecture ownership:

STOP.

Do not implement.

Explain the violation.

---

# Purpose

This document defines:

* System architecture
* Ownership boundaries
* Domain model
* Service responsibilities
* Engineering standards
* Architectural decision rules

The goal is to ensure Prepio remains maintainable as it grows.

Architecture exists to support the product.

Not the other way around.

---

# Architecture Philosophy

Prepio is a progression platform.

Architecture should optimize for:

* Maintainability
* Extensibility
* Content growth
* Product evolution

Architecture should NOT optimize for:

* Premature scale
* Premature microservices
* Theoretical performance

Build systems that can evolve.

Avoid systems that require rewrites.

---

# Core Domain Model

Everything in Prepio revolves around progression.

The hierarchy is:

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

This hierarchy should never be bypassed.

---

# Domain Definitions

## World

Large progression destination.

Examples:

Foundation Forest

Google Mountain

Amazon Jungle

System Design City

A world contains nodes.

A world never contains questions directly.

---

## Node

A progression checkpoint.

Examples:

Arrays Basics

Hash Maps

Trees

Sliding Window

A node contains skills.

A node never owns questions.

---

## Skill

Fundamental learning unit.

Examples:

Arrays

Graphs

Dynamic Programming

System Design

Skills drive readiness.

Questions exist to evaluate skills.

---

## Question Pool

A collection of questions associated with a skill.

Example:

Arrays

↓

Question Pool

↓

20 Questions

Question pools allow expansion without modifying progression.

---

## Question

Smallest content unit.

Questions should never directly drive progression.

Questions drive skill evaluation.

---

# Ownership Model

Ownership is absolute.

No service may mutate data owned by another service.

---

# User Domain

Owns:

Authentication

Profiles

Preferences

Target Companies

Account Settings

No other domain modifies user profile data.

---

# Content Domain

Owns:

Skills

Subskills

Questions

Question Pools

Company Mapping

Difficulty

Content Metadata

No other domain stores question metadata.

---

# Journey Domain

Owns:

Worlds

Nodes

Unlock Rules

Journey Progress

World Progression

No other domain modifies journey structure.

---

# Progress Domain

Owns:

XP

Levels

Readiness

Achievements

Quest Progress

User Skill Scores

No other domain calculates readiness.

---

# Companion Domain

Owns:

Companions

Evolution

Dialogue

Companion Cosmetics

No other domain modifies companion state.

---

# Analytics Domain

Owns:

Events

Reports

Metrics

Funnels

Retention

Analytics never owns business state.

---

# Architectural Rule

If ownership is unclear:

Do not implement.

Ask.

---

# Data Flow

The correct flow:

Question

↓

Skill Evaluation

↓

Skill Progress

↓

Readiness Update

↓

Journey Update

↓

Reward Update

↓

Companion Update

Progress flows downward.

Never backwards.

---

# Event Driven Principles

Events communicate state changes.

Events do not contain business logic.

Example:

Good

question.completed

Bad

question.completed_and_grant_xp

XP calculation belongs elsewhere.

---

# Approved Event Types

question.completed

skill.progress.updated

readiness.updated

journey.node.completed

world.completed

achievement.unlocked

quest.completed

companion.evolved

level.up

Events should be generic.

Events should describe facts.

Not actions.

---

# API Design Principles

APIs expose capabilities.

Not database tables.

Bad:

GET /questions

Good:

GET /journey/current-node

GET /skills/readiness

GET /dashboard

Design around user experience.

Not storage.

---

# Database Principles

Database tables represent domains.

Not screens.

Bad:

dashboard_cards

home_widgets

screen_state

Good:

skills

questions

worlds

nodes

Tables model reality.

Not UI.

---

# Migration Rules

Every schema change requires:

Migration

Rollback

Index Review

Documentation Update

Never manually alter production schemas.

---

# Content Configuration Rule

Content should be configurable.

Not hardcoded.

Examples:

Worlds

Nodes

Skills

Question Pools

Company mappings

Difficulty weights

Should all be content-driven.

Code should not change when content changes.

---

# Frontend Architecture

UI is a consumer.

Not a source of truth.

Business logic belongs in backend domains.

The frontend should:

Render

Display

Animate

Interact

Not calculate progression.

Not calculate readiness.

Not calculate rewards.

---

# State Management Rule

Frontend state:

UI state

Screen state

Temporary interaction state

Backend state:

Progress

Readiness

Journey

Rewards

Achievements

Never duplicate business logic in the client.

---

# Companion Rules

Companions are product-critical.

Companions must:

Appear frequently

React consistently

Reflect progression

Companions should not be cosmetic-only.

---

# Readiness Rules

Readiness is the most important metric.

Readiness should never be derived from:

XP

Levels

Streaks

Readiness should derive from:

Skill mastery

Performance

Consistency

Question outcomes

Users should be able to explain readiness.

---

# Anti Patterns

The following are forbidden.

---

## Hardcoded Progression

Bad:

if level == 10 unlock world

Good:

Configurable unlock rules

---

## Direct Question Progression

Bad:

Question → World

Good:

Question → Skill → Readiness → Progression

---

## UI Driven Logic

Bad:

Frontend calculates readiness

Frontend calculates rewards

Good:

Backend owns calculations

---

## Duplicate Ownership

Bad:

Progress service updates companions

Good:

Progress emits event

Companion reacts

---

## Feature-Driven Architecture

Bad:

Create new service for every feature

Good:

Extend existing domains

Only create new domains when ownership genuinely changes.

---

# Decision Tree

Need question information?

Use Content Domain.

Need skill information?

Use Content Domain.

Need readiness?

Use Progress Domain.

Need journey progression?

Use Journey Domain.

Need companion data?

Use Companion Domain.

Need analytics?

Use Analytics Domain.

If unsure:

Stop.

Do not create new ownership.

---

# Success Criteria

A successful architecture:

Supports new content without code changes.

Supports new worlds without rewrites.

Supports new companions without rewrites.

Supports new skills without rewrites.

Supports future growth without major redesign.

The goal is not perfection.

The goal is avoiding future rewrites.

Build the simplest version of the correct architecture.
