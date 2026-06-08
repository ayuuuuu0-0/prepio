# IMPLEMENTATION_RULES.md

# Purpose

This document is the highest-priority implementation guide for all AI agents and contributors working on Prepio.

The goal is not simply to make features work.

The goal is to make features work WITHOUT breaking:

* Product Vision
* Service Boundaries
* User Experience
* Event Ownership
* Data Ownership

Every implementation must follow these rules.

Violation of these rules is considered a bug even if the code compiles.

---

# Rule 1: Product First

Always remember:

Prepio is NOT an interview question platform.

Prepio is a progression game where interview preparation is the progression mechanic.

When implementing features:

DO NOT ask:

"How do I show questions?"

Ask:

"How does this increase user progression?"

Every feature should improve one of:

* Progress
* Motivation
* Retention
* Competition
* Collection
* Achievement

If it does none of these, question whether the feature should exist.

---

# Rule 2: Dashboard Is The Most Important Screen

Whenever adding new functionality:

First determine:

How does this appear on the dashboard?

The dashboard is the user's home.

Users should never need to dig through menus to understand:

* Their streak
* Their progress
* Their readiness
* Their companion
* Their league

If a feature does not surface meaningful progress on the dashboard, reconsider the feature.

---

# Rule 3: Service Ownership Is Absolute

Never violate service ownership.

Never.

---

## User Service Owns

* Authentication
* User Profiles
* User Preferences
* Company Targets
* Account Settings

No other service may write user profile data.

---

## Companion Service Owns

* Companions
* Companion Evolution
* Companion Cosmetics
* Companion Dialogues

No other service may mutate companion state.

---

## Journey Service Owns

* Worlds
* Journey Nodes
* Unlock Logic
* World Progression

No other service may modify journey state.

---

## Question Service Owns

* Questions
* Answers
* Evaluations
* Challenge Content

No other service stores questions.

---

## Progress Service Owns

* XP
* Gems
* Levels
* Readiness
* Achievements
* Quests

No other service writes XP.

No other service writes gems.

No other service calculates readiness.

---

## Leaderboard Service Owns

* Leagues
* Rankings
* Seasons

No other service calculates rankings.

---

## Arena Service Owns

* Matchmaking
* Arena Rating
* Battle Results

No other service decides arena outcomes.

---

## Notification Service Owns

* Push Notifications
* Character Messaging
* Engagement Notifications

No service may directly send FCM notifications.

EVERY notification goes through Notification Service.

---

# Rule 4: Event Ownership

All progression must be event driven.

Never directly mutate another service's state.

Bad:

Question Service updates XP.

Good:

Question Service emits:

question.answered

Progress Service consumes:

question.answered

Progress Service updates XP.

---

Approved Events

question.answered

quest.completed

achievement.unlocked

boss.completed

world.unlocked

league.promoted

league.demoted

readiness.updated

companion.evolved

arena.finished

streak.updated

Only introduce new events if absolutely necessary.

---

# Rule 5: Readiness Is The Core Metric

XP is not the product.

Gems are not the product.

Readiness is the product.

Whenever implementing analytics:

Prioritize:

Google Readiness

Amazon Readiness

Meta Readiness

Uber Readiness

Atlassian Readiness

Not:

Questions Solved

Users care about outcomes.

Not activity.

---

# Rule 6: Companions Are Sacred

Companions are not profile pictures.

Companions are emotional progression systems.

Every companion feature must answer:

Why would a user become attached to this companion?

Bad:

Static image.

Good:

Growth
Reactions
Evolution
Dialogue

Every companion should feel alive.

---

# Rule 7: Every Action Must Produce Feedback

Silent actions are forbidden.

Bad UX:

User completes challenge.

Nothing happens.

Good UX:

Challenge Complete
+50 XP
+10 Gems
Readiness +1%
Companion Reaction
Progress Animation

Every action should create visible progress.

---

# Rule 8: Progress Must Be Visible

Users should never wonder:

"Am I improving?"

Every screen should show progression.

Examples:

Progress Bars

Readiness Score

Journey Completion

Companion Growth

League Rank

Achievement Collection

Always expose progress.

---

# Rule 9: Collection Drives Retention

Whenever adding content:

Ask:

Can this be collected?

Examples:

Badges

Titles

Achievements

Companion Skins

League Rewards

World Trophies

Collection increases retention.

---

# Rule 10: Company Worlds Are Products

Google World

Amazon World

Meta World

Uber World

These are not question categories.

These are destinations.

Every world should have:

* Theme
* Visual Identity
* Boss
* Achievements
* Rewards

Treat them like mini-products.

---

# Rule 11: Bosses Must Feel Special

Bosses are rare.

Bosses should never feel like normal questions.

Bosses require:

Unique UI

Special Animations

Large Rewards

Unique Achievements

Bosses create memorable moments.

---

# Rule 12: Mobile And Web Have Different Purposes

Mobile:

Engagement

Habit Building

Quick Sessions

Companions

Journey

Streaks

Web:

Deep Learning

Analytics

Mock Interviews

Reports

Company Preparation

Never blindly duplicate screens.

---

# Rule 13: Flutter Rules

State Management:

Riverpod Only

Forbidden:

Provider
GetX
Bloc

Without architecture review.

---

Animations

Every major feature requires animation.

Use:

Rive

Lottie

Hero

Implicit Animations

No static UX.

---

Networking

Only repositories communicate with APIs.

Widgets never call APIs.

---

# Rule 14: Next.js Rules

Use App Router.

Server Components where possible.

Client Components only when necessary.

Use SWR or React Query.

Never fetch directly inside random components.

Create reusable hooks.

---

# Rule 15: Go Rules

No business logic inside handlers.

Handlers:

Validate

Authenticate

Call Service Layer

Return Response

Business logic belongs in services.

Persistence belongs in repositories.

---

Required Layers

handler

service

repository

model

No shortcuts.

---

# Rule 16: Database Rules

Every feature requires:

Migration

Indexes

Rollback Path

Never modify production tables manually.

Always create migrations.

---

# Rule 17: UX Rules

Users must always know:

Where they are

What they achieved

What they unlock next

If any screen fails to answer these questions:

Redesign it.

---

# Rule 18: Feature Acceptance Checklist

Before merging any feature:

Does it improve progression?

Does it improve retention?

Does it respect service ownership?

Does it respect event ownership?

Does it surface progress?

Does it fit the product vision?

Does it feel rewarding?

If any answer is NO:

Do not merge.

---

# Final Principle

Every implementation decision should move Prepio closer to this feeling:

"I am becoming interview ready."

Not:

"I am answering more questions."

If a feature strengthens that feeling, build it.

If it weakens that feeling, reject it.
