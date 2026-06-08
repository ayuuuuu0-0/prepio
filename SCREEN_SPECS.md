# SCREEN_SPECS.md

# Purpose

This document defines exactly how every major screen in Prepio should look, feel, animate, and behave.

This document exists because AI agents are naturally biased toward:

* Dashboards
* Tables
* Cards
* Forms
* Analytics

Prepio must NOT look like:

* Jira
* Notion
* Salesforce
* SAP
* Government Portals
* Enterprise Admin Dashboards

Prepio should feel like:

* Duolingo
* Pokemon
* Clash Royale
* Brawl Stars

The user should feel excitement before they solve a single question.

---

# GLOBAL DESIGN RULES

## Visual Energy Rule

Every screen must contain:

* Color
* Motion
* Character
* Progress
* Reward

No screen should feel static.

---

## Companion Visibility Rule

Companion must appear on:

Dashboard

Journey

Results

Quest Screen

Leaderboard

Character Screen

Profile

Companion should never disappear.

---

## Reward Visibility Rule

Users should always see:

XP

Gems

Level

Readiness

League

Current Goal

No hidden progression.

---

# SCREEN 1 — DASHBOARD

Most important screen.

Not Journey.

Not Questions.

Dashboard.

Users will see this 1000+ times.

---

# Goal

When opening dashboard:

User immediately understands:

* How far they have come
* What to do next
* What reward is waiting

---

# Layout

Top Section

Large Companion

Speech Bubble

Example:

Byte:

"Only 2 challenges until Level 18."

---

Middle Section

Readiness Cards

Google Readiness

72%

Amazon Readiness

65%

Meta Readiness

58%

Displayed as colorful circular progress rings.

---

League Card

🏆 Sapphire League

Rank #4

Animated rank movement indicator.

---

Daily Quest Card

⚡ Solve 3 DSA Questions

Progress bar

Reward preview

---

Bottom Section

Huge Continue Journey Button

Primary CTA.

Must be impossible to miss.

---

# Animation

Companion Idle Animation

Every 5 seconds:

Blink

Wave

Small movement

---

Readiness Bars

Animate on page load.

---

Quest Progress

Smooth fill animation.

---

# SCREEN 2 — JOURNEY MAP

Hero feature.

70% of design effort goes here.

---

# Goal

Feel like adventure.

Not curriculum.

---

# Structure

World Header

World Background

Journey Nodes

Companion

Rewards

Bosses

---

# World Example

Foundation Forest

Background:

Trees

Floating particles

Soft green environment

---

Google Mountain

Blue mountain range

Clouds

Tech aesthetic

---

Amazon Jungle

Orange jungle

Vines

Ancient ruins

---

# Node Types

Normal

Green circle

---

Challenge

Blue circle

---

Treasure

Golden chest

---

Premium

Purple crystal

---

Boss

Large animated crown node

---

# Node Animations

Completed

Bounce once

Show checkmark

---

Current

Pulse continuously

---

Locked

Dimmed

---

# Companion

Companion physically sits on active node.

Not floating randomly.

User should instantly know:

"This is where I am."

---

# SCREEN 3 — QUESTION SCREEN

Focus mode.

Less decoration.

Still alive.

---

# Layout

Top

Progress

Question Number

Timer

Companion Small Avatar

---

Middle

Question

Examples

Code Editor

---

Bottom

Submit Button

Large

Rounded

Bright

---

# Interaction

Submit button should feel powerful.

Not like a form.

---

# SCREEN 4 — RESULT SCREEN

Most emotional screen.

Critical for retention.

---

# Correct Answer

Screen Color

Green

---

Companion

Victory animation

---

Effects

Confetti

XP Count Up

Gem Count Up

Readiness Increase

Achievement Popups

---

Display

+50 XP

+10 Gems

Google Readiness +2%

---

# Wrong Answer

Red is forbidden.

Never punish aggressively.

---

Use:

Orange

Purple

Soft warning colors

---

Companion says:

"Almost there."

"Good attempt."

"Let's improve this."

---

Show:

What went wrong

What to learn

Next action

---

# SCREEN 5 — COMPANION PAGE

Pokemon screen.

Not settings screen.

---

# Layout

Large Companion

Evolution Timeline

Unlocked Skins

Animations

Dialogues

---

# Evolution Section

Baby

↓

Student

↓

Engineer

↓

Senior

↓

Staff

↓

Legend

---

Locked stages visible.

Users should want them.

---

# SCREEN 6 — LEAGUE SCREEN

Competition.

---

# Hero Card

Current League

Sapphire

Rank #4

---

# Top Players

Large avatars

Companions visible

XP visible

---

# Promotion Zone

Green area

---

# Demotion Zone

Red area

---

# User Row

Highlighted.

Always visible.

---

# SCREEN 7 — QUEST SCREEN

Quest board.

Not task manager.

---

# Layout

Quest Cards

Big icons

Progress

Rewards

---

Example

⚡ Solve 3 DSA Questions

Reward:

+20 Gems

---

Completed Quest

Confetti

Checkmark

Claim Animation

---

# SCREEN 8 — ACHIEVEMENTS

Pokemon Badge Book.

---

# Categories

DSA

System Design

Companies

Arena

Streaks

Special Events

---

# Locked Achievements

Visible.

Greyed out.

---

Users should feel:

"I want that."

---

# SCREEN 9 — SHOP

Not ecommerce.

Reward shop.

---

# Categories

Companions

Skins

Titles

World Themes

Special Effects

---

# Display

Large visual cards

Gem cost

Preview animation

---

# SCREEN 10 — PROFILE

Career identity.

---

Display

Level

Companion

Achievements

Readiness

League

Journey Completion

---

# Visual Rule

No boring stats tables.

Everything visual.

---

# ANIMATION RULES

Every reward animates.

Every progression animates.

Every level up animates.

Every achievement animates.

Every evolution animates.

No silent success.

---

# SOUND PREPARATION

Even if sound is not implemented:

Design every feature assuming future sound effects.

Level Up

Quest Complete

Boss Victory

League Promotion

Achievement Unlock

Companion Evolution

---

# FINAL TEST

Before shipping any screen ask:

Does this feel like:

A government website?

An admin dashboard?

A SaaS analytics tool?

If yes:

Redesign.

Then ask:

Would a Duolingo designer be proud of this screen?

If not:

Redesign again.

That standard is mandatory for all Prepio interfaces.

---

# Voice Guidelines

Every string visible to the user should pass this test:
"Would a senior engineer say this to a junior colleague?"

WRONG:  "You've got this! Take your time and think it through."
RIGHT:  "Think about hash maps for this one."

WRONG:  "Amazing! You're on a roll!"
RIGHT:  "Solid." or "Good — time complexity covered."

WRONG:  "Almost there! Keep going!"
RIGHT:  "Not quite. Here's what was missing: [specific bullet points]"

WRONG:  "LET'S GO!" (button copy)
RIGHT:  "Continue Prep" or "Submit" or "Start Session"

WRONG:  "New adventurer? Start your journey"
RIGHT:  "No account? Join 12k engineers in prep"

The companion can be warm and encouraging.
The companion cannot be condescending or childish.
There is a difference between "I believe in you!" and "One challenge from Level 12."
One is generic. One is specific and motivating.

---

# Input Styling

Dark background inputs, not white:
background: #1A1D27
border: 1px solid #2E3347
focus border: #7C6EF5
text: #E8EAED

No white or grey form fields anywhere.

---

# Result Screen Rules

Correct answer: Do NOT say "Amazing!" Say "Solid."
Wrong answer:   Do NOT say "Almost there!" Show what was missed.
Wrong answer:   Do NOT use red. Use orange (#FF6B35) or amber (#F5B942).
XP and gems:    Always animate with count-up. Never appear statically.
Feedback:       At least 2 bullet points. Never a single vague sentence.
