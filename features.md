# FEATURES.md

# Purpose

This document defines the implementation order of Prepio.

Agents MUST follow this order.

Do not skip phases.

Do not build future phases before prerequisites exist.

Each phase creates the foundation for the next phase.

---

# Build Philosophy

Build in this order:

1. Foundation
2. Core Progression
3. Retention
4. Competition
5. Collection
6. Premium Content
7. Multiplayer
8. AI Features

The user should be able to use the product after every phase.

---

# PHASE 0 — Platform Foundation

Goal:

Users can authenticate and access the platform.

---

## Backend

Implement:

* Gateway
* User Service
* Authentication
* JWT
* Refresh Tokens
* PostgreSQL
* Redis
* Kafka Setup

---

## Flutter

Implement:

* Splash Screen
* Login
* Register
* Session Persistence

---

## Web

Implement:

* Login
* Register
* Session Management

---

## Completion Criteria

User can:

* Register
* Login
* Logout
* Persist session

Only then move forward.

---

# PHASE 1 — Onboarding

Goal:

Personalize the experience.

---

## Backend

Create:

user_targets

Store:

* Target Companies
* Experience Level

Endpoints:

POST /users/onboarding

GET /users/profile

---

## Flutter

Build:

Welcome Flow

Step 1

Choose Companies

Google

Amazon

Meta

Uber

Atlassian

Step 2

Choose Experience

Fresher

Junior

Mid

Senior

Step 3

Choose Companion

---

## Web

Same flow.

---

## Completion Criteria

Every user has:

company targets

experience level

active companion

---

# PHASE 2 — Core Progression Engine

Goal:

Users can answer questions and gain progress.

---

## Backend

Implement:

Question Service

Question Retrieval

Question Submission

Evaluation

Progress Service

XP

Levels

Gems

Level Formula

Reward Config

---

Kafka Events

question.answered

progress.updated

---

## Flutter

Question Screen

Submit Screen

Result Screen

XP Animation

Gem Animation

---

## Web

Question Page

Result Page

---

## Completion Criteria

User answers question.

Receives:

XP

Gems

Level Progress

---

# PHASE 3 — Dashboard

Goal:

Create the first dopamine loop.

---

## Backend

Create:

GET /dashboard/home

Return:

XP

Level

Streak

Gems

Readiness

Companion

League

Daily Quests

---

## Flutter

Dashboard

Components:

Top Bar

Readiness Card

Companion Card

Daily Quest Card

Continue Journey Button

---

## Web

Same structure.

---

## Completion Criteria

User immediately understands:

Current Progress

Next Goal

Current Status

---

# PHASE 4 — Journey System

Goal:

Replace question lists with progression.

---

## Backend

Create:

worlds

journey_nodes

user_journey_progress

world_unlocks

boss_nodes

---

## Flutter

Journey Map

Animated Nodes

Treasure Nodes

Boss Nodes

World Backgrounds

---

## Web

Zoomable SVG Map

---

## Worlds

Foundation

Problem Solving

System Design

Company Worlds

---

## Completion Criteria

Users progress through worlds.

Not question lists.

---

# PHASE 5 — Companion System

Goal:

Create emotional attachment.

---

## Backend

Create:

companions

companion_stages

companion_dialogues

user_companions

---

Evolution System

Level 1

Baby

Level 10

Student

Level 25

Engineer

Level 50

Senior

Level 75

Staff

Level 100

Legend

---

## Flutter

Rive Animations

Idle

Happy

Sad

Victory

Level Up

---

## Web

Animated Companion Cards

---

## Completion Criteria

Companions evolve.

Users can equip companions.

Companions react.

---

# PHASE 6 — Streak System

Goal:

Increase retention.

---

## Backend

Create:

streak_service

streak_tracking

streak_freezes

---

Events:

question.answered

streak.updated

---

## Flutter

Flame Indicator

Streak Screen

Freeze Purchase

---

## Web

Same.

---

## Completion Criteria

Users maintain streaks.

Streak breaks correctly.

Freezes work.

---

# PHASE 7 — Readiness Engine

Goal:

Show meaningful progress.

---

## Backend

Create:

skill_readiness

company_readiness

---

Track:

Arrays

Trees

Graphs

DP

LLD

HLD

Behavioral

---

Generate:

Google Readiness

Amazon Readiness

Meta Readiness

Uber Readiness

---

## Flutter

Readiness Dashboard

Progress Charts

---

## Web

Advanced Analytics

---

## Completion Criteria

Users know:

How interview-ready they are.

---

# PHASE 8 — Quests & Achievements

Goal:

Create short-term goals.

---

## Backend

quests

user_quests

achievements

user_achievements

---

## Flutter

Quest Screen

Achievement Book

Reward Animations

---

## Web

Same.

---

## Completion Criteria

Users earn achievements.

Users complete quests.

---

# PHASE 9 — League System

Goal:

Create competition.

---

## Backend

leaderboard_service

league_tiers

weekly_ranking

season_reset

---

Tiers

Bronze

Silver

Gold

Sapphire

Ruby

Emerald

Diamond

Legend

---

## Flutter

League Screen

Leaderboard Screen

Promotion Animation

---

## Web

Advanced Leaderboards

---

## Completion Criteria

Weekly competition works.

Promotions work.

Demotions work.

---

# PHASE 10 — Company Worlds

Goal:

Create destinations.

---

Google Mountain

Amazon Jungle

Meta City

Uber Highway

Netflix Studios

Atlassian Fortress

---

Each World Includes

Questions

Achievements

Boss

Theme

Rewards

---

## Completion Criteria

Company preparation becomes structured.

---

# PHASE 11 — Boss Interviews

Goal:

Create memorable moments.

---

## Backend

boss_rounds

boss_attempts

boss_rewards

---

Examples

Google Hiring Committee

Amazon Bar Raiser

Meta Product Round

Uber Design Review

---

## Flutter

Special Boss UI

Boss Intro

Boss Victory

Boss Failure

---

## Web

Same.

---

## Completion Criteria

Bosses feel unique.

---

# PHASE 12 — Premium Economy

Goal:

Give value to gems.

---

## Backend

premium_nodes

premium_unlocks

cosmetics

skins

---

Users Spend Gems On

Premium Worlds

Companion Skins

Boss Retries

Titles

Cosmetics

---

## Completion Criteria

Gems become meaningful.

---

# PHASE 13 — Arena

Goal:

Create multiplayer engagement.

---

## Backend

arena_service

matchmaking

arena_rating

arena_matches

---

Flow

Queue

Match

Challenge

Evaluate

Reward

---

## Flutter

Arena UI

Matchmaking

Result Screen

---

## Web

Arena Dashboard

---

## Completion Criteria

Players compete live.

---

# PHASE 14 — Notifications

Goal:

Bring users back.

---

## Backend

notification_service

FCM

Email

Scheduled Reminders

---

Triggers

Streak Risk

Quest Completion

Promotion

Boss Available

Companion Evolution

---

## Completion Criteria

Retention notifications work.

---

# PHASE 15 — Seasons

Goal:

Long-term retention.

---

Season Duration

90 Days

---

Rewards

Exclusive Skins

Exclusive Titles

Exclusive Badges

---

Season Reset

League Reset

Arena Reset

---

## Completion Criteria

Seasonal progression works.

---

# PHASE 16 — AI Features

Goal:

Become best-in-class.

---

AI Interviewer

Voice Interview

Resume Review

Company Readiness Analysis

Interview Feedback

Behavioral Evaluation

---

## Completion Criteria

AI becomes enhancement.

Not foundation.

---

# Strict Agent Rules

Never build Phase N+1 before Phase N is complete.

Never introduce new services without architecture review.

Never bypass event ownership.

Never duplicate business logic.

Never store progression outside Progress Service.

Never store companion state outside Companion Service.

Never calculate readiness outside Readiness Service.

Always update:

Architecture

Migrations

API Contracts

Frontend Types

Documentation

before marking a feature complete.

The order in this file is mandatory.
