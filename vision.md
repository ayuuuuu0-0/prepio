# PREPIO 2.0 MASTER PRODUCT & ARCHITECTURE DOCUMENT

# Vision

Prepio is not an interview preparation platform.

Prepio is a progression game where the reward is becoming interview-ready.

Users do not wake up wanting to solve DSA questions.

Users wake up wanting:

* Better jobs
* Better companies
* Better salaries
* Career growth

The purpose of every system in Prepio is to make users feel closer to their career goals every single day.

---

# Product Positioning

Current Market

* LeetCode → Question Bank
* InterviewBit → Practice Platform
* GeeksForGeeks → Learning Platform
* Duolingo → Habit Building Platform

Prepio Goal

* Duolingo + Pokémon + Career Growth Platform

The user should feel:

"I am leveling up my career."

Not:

"I am solving questions."

---

# Core User Loop

Open App
↓
Companion Greets User
↓
Check Readiness Score
↓
Continue Journey
↓
Complete Challenge
↓
Earn XP
↓
Earn Gems
↓
Improve Readiness
↓
Companion Evolves
↓
League Rank Improves
↓
Unlock New World
↓
Return Tomorrow

Every feature must strengthen this loop.

---

# Pillar 1: Journey System

Journey becomes the primary interface.

Questions are no longer the product.

The journey is the product.

---

## World Structure

Foundation World

* Arrays
* Strings
* Hashing
* Sorting
* Foundation Boss

Problem Solving World

* Trees
* Graphs
* Greedy
* Dynamic Programming
* Problem Solving Boss

System Design World

* Scalability
* Databases
* Caching
* Messaging
* Design Boss

Company Worlds

Google Mountain

* OA
* DSA
* LLD
* HLD
* Hiring Committee Boss

Amazon Jungle

* OA
* DSA
* Leadership Principles
* LLD
* Bar Raiser Boss

Meta City

* Coding
* Product Sense
* Behavioral
* Architecture
* Meta Panel Boss

Uber Highway

* Coding
* LLD
* HLD
* Debugging
* Uber Technical Committee

---

# Journey UX

Flutter

* Full screen vertical path
* Parallax world backgrounds
* Animated node transitions
* Companion floating near active node
* Treasure nodes
* Boss nodes
* Premium nodes

Web

* SVG-based zoomable map
* World transitions
* Animated completion effects
* Scrollable progression path

---

# Pillar 2: Companion System

Companions are emotional progression systems.

Not profile avatars.

---

## Starter Companions

Byte – Capybara

Pip – Red Panda

Nova – Pangolin

Kodo – Axolotl

Zara – Snow Leopard

---

## Evolution System

Level 1

Baby Companion

Level 10

Student Companion

Level 25

Engineer Companion

Level 50

Senior Companion

Level 75

Staff Companion

Level 100

Legendary Companion

Every stage unlocks:

* New artwork
* New animations
* New reactions
* New dialogue
* New profile badge

---

## Companion Reactions

Correct Answer

Victory animation

Wrong Answer

Encouraging animation

Streak Milestone

Celebration animation

League Promotion

Special animation

Boss Victory

Legendary animation

---

# Pillar 3: Readiness Engine

This is the true IP.

XP is not the goal.

Readiness is the goal.

---

## User Skills

Arrays

Strings

Hashing

Trees

Graphs

Dynamic Programming

LLD

HLD

Behavioral

Problem Solving

System Design

Communication

---

## Readiness Formula

Readiness Score =

Performance
+
Difficulty Weight
+
Consistency
+
Recent Activity
+
Boss Completion

---

## Company Readiness

Google Readiness

Amazon Readiness

Meta Readiness

Uber Readiness

Atlassian Readiness

Netflix Readiness

Displayed on dashboard at all times.

---

# Pillar 4: Leagues

Weekly competition.

Bronze

Silver

Gold

Sapphire

Ruby

Emerald

Diamond

Legend

---

## League Rewards

Promotion Chest

Contains:

* Gems
* Titles
* Skins
* Companion Cosmetics

---

# Pillar 5: Collection System

Humans love collecting.

---

## Achievement Categories

Questions Solved

Streaks

Boss Victories

Company Completion

Arena Wins

System Design

Behavioral

Special Events

---

## Achievement Book

Pokémon-style collection page.

Users continuously unlock:

* Badges
* Titles
* Trophies
* Skins

---

# Pillar 6: Boss Interviews

Every world ends with a boss.

---

Examples

Google Hiring Committee

Amazon Bar Raiser

Meta Product Interview

Uber Architecture Review

---

Boss Rewards

Huge XP

Huge Gems

Unique Badge

World Unlock

---

# Pillar 7: Daily Quests

Daily Objectives

Examples

Solve 3 DSA Questions

Score Above 80%

Complete One System Design Challenge

Maintain Streak

Finish One Company Challenge

---

Rewards

XP

Gems

Cosmetics

Achievements

---

# Pillar 8: Arena

Future flagship feature.

---

Arena Flow

Join Queue
↓
Find Opponent
↓
Receive Same Question
↓
Timed Round
↓
Evaluation
↓
Winner Determined

---

Rewards

XP

Arena Rating

League Points

Trophies

---

# Dashboard Design

Most Important Screen

---

Top Section

🔥 27 Day Streak

🦫 Byte says:
"Only 2 challenges until Level 18."

---

Middle Section

Google Readiness 72%

Amazon Readiness 64%

Meta Readiness 58%

---

Current League

Sapphire League

Rank #4

---

Daily Quests

2/3 Completed

---

Continue Journey

Primary CTA

---

# Backend Architecture

Microservices

gateway

user-service

journey-service

question-service

progress-service

companion-service

leaderboard-service

notification-service

arena-service

analytics-service

---

# Event Driven Architecture

Kafka Events

question.answered

quest.completed

achievement.unlocked

league.promoted

boss.completed

world.unlocked

readiness.updated

companion.evolved

arena.finished

---

# Service Ownership

User Service

* Identity
* Profiles
* Authentication
* Preferences

Companion Service

* Companions
* Evolution
* Dialogues
* Cosmetics

Journey Service

* Worlds
* Nodes
* Progression

Question Service

* Questions
* Evaluation
* Challenges

Progress Service

* XP
* Gems
* Levels
* Readiness

Leaderboard Service

* Leagues
* Rankings
* Seasons

Arena Service

* Matchmaking
* Battles
* Ratings

Notification Service

* Push Notifications
* Character Messages

Analytics Service

* Performance Insights
* User Reports

---

# Flutter Architecture

features/

auth

dashboard

journey

companions

arena

leaderboard

quests

achievements

profile

analytics

---

State Management

Riverpod

---

Animations

Rive

Lottie

Hero Animations

Parallax Scrolling

---

# Next.js Architecture

app/

dashboard

journey

companies

arena

leaderboard

companions

achievements

analytics

profile

---

Web Focus

Detailed Progress

Interview Analytics

Mock Interviews

Company Preparation

Reports

---

# Design Principles

Every action should produce feedback.

No silent rewards.

No invisible progression.

Every session must produce:

* XP Gain
* Gem Gain
* Readiness Change
* Companion Interaction
* Progress Visualization

---

# Phase 1

Journey System

Dashboard

XP

Gems

Companions

Readiness Engine

---

# Phase 2

Achievements

Leagues

Company Worlds

Boss Interviews

---

# Phase 3

Arena

Seasons

Cosmetics

Companion Evolution

---

# Phase 4

AI Interviewers

Voice Interviews

Guilds

Mentorship

Career Marketplace

---

# Success Metric

A user should be able to open the app after 100 days and immediately feel:

"I have become significantly more interview-ready than when I started."

That feeling is the core product.

Everything else exists to support it.

---

# Audience Differentiation

Duolingo's user: Learning French for vacation, hobby, personal growth.
Prepio's user:   Trying to get their first job at Google or double their salary.

The stakes are categorically different.

Duolingo can be whimsical because nothing rides on a missed lesson.
Prepio users are under real pressure. The UI must honor that.

Playful: Yes — this is a game and games should feel rewarding.
Childish: Never — users are making career decisions that affect their income.

The emotional target for every session:
"I made measurable progress toward the job I want."

Not:
"I earned a fun star today!"

---

# Design Reference Points

Primary inspiration: Duolingo (genre, mechanics, companion)
Secondary inspiration: Linear, Vercel, Raycast (professionalism, dark theme)
Tertiary inspiration: Pokémon (collection, evolution, world progression)

The combination:
- Duolingo's daily habit loop
- Linear's dark professional aesthetic
- Pokémon's emotional attachment to companions

This combination does not currently exist in the market.
It is Prepio's moat.
