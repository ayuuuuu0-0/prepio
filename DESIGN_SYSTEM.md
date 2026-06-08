# DESIGN_SYSTEM.md

# Purpose

Prepio must NEVER look like:

* Government software
* Admin dashboard
* Corporate HR portal
* Analytics platform
* SaaS CRM

If a screen resembles any of these:

It is wrong.

---

# Design Philosophy

Prepio is:

A progression game.

Not a dashboard.

Not a productivity tool.

Not an admin panel.

Users should feel:

Excited

Motivated

Rewarded

Curious

Proud

---

# Inspiration Sources

Primary:

Duolingo

Secondary:

Pokemon

Clash Royale

Brawl Stars

Supercell Products

Discord

Linear

---

# Visual Personality

Words that describe the UI:

Playful

Energetic

Friendly

Colorful

Rewarding

Alive

Expressive

Never:

Corporate

Minimalist

Enterprise

Boring

Sterile

---

# Companion First Design

Every major screen must contain:

A companion.

Not optional.

---

Dashboard

Companion visible

Journey

Companion visible

Result Screen

Companion visible

Quest Screen

Companion visible

League Screen

Companion visible

---

Bad:

Small avatar in corner.

Good:

Large animated companion.

---

# Color System

Forbidden:

White background

Grey cards

Enterprise blue

Bootstrap appearance

Material default appearance

---

Primary Colors

Green

Purple

Orange

Blue

Gold

Pink

---

World Colors

Foundation World

Green

Google World

Blue

Amazon World

Orange

Meta World

Purple

Uber World

Black + Green

Netflix World

Red

---

Every world must feel unique.

---

# Shapes

Nothing should feel sharp.

Use:

24px radius

32px radius

Pill buttons

Rounded cards

Rounded modals

Rounded progress bars

Rounded achievement badges

---

# Animations

Animations are mandatory.

Every user action produces feedback.

---

Correct Answer

Character jumps

XP counts up

Gems count up

Confetti

Screen bounce

---

Level Up

Full screen celebration

Companion evolves

Particle effects

Sound effect ready

---

Quest Complete

Reward animation

Progress fill

Character reaction

---

# Dashboard Rules

Dashboard is NOT analytics.

Dashboard is emotion.

---

Top Section

Companion

Speech Bubble

Current Goal

---

Middle Section

Readiness

Journey Progress

League Position

---

Bottom Section

Continue Journey

Large button

Primary action

---

# Cards

Every card should have:

Icon

Color

Illustration

Animation

Progress

Reward

---

Never:

Text-only cards

---

# Journey Map Rules

Journey is the hero screen.

It should occupy 70% of design effort.

---

Nodes

Large

Colorful

Animated

Touchable

Rewarding

---

Node Types

Normal

Treasure

Boss

Premium

Challenge

---

Every node visually distinct.

---

# Character Rules

Companions are products.

Not decorations.

---

Every companion requires:

Idle animation

Happy animation

Sad animation

Thinking animation

Victory animation

---

Companions must react.

Always.

---

# Typography

Friendly.

Rounded.

Modern.

Avoid:

Corporate fonts

Monospace-heavy UI

Dense text

---

# Reward Rules

Never silently reward.

Always animate rewards.

---

XP

Animated

---

Gems

Animated

---

Achievements

Animated

---

League Promotions

Animated

---

Boss Victories

Animated

---

# Agent Rules

If a screen resembles:

Admin Panel

Analytics Dashboard

Government Website

Enterprise SaaS

Then redesign it.

Before implementation ask:

Would Duolingo ship this?

If the answer is no:

Redesign.

The UI should feel like a game first.

Everything else second.

---

# Target Audience

Users are 19–28 year old engineers and CS students.
They are chasing real career outcomes — job offers, salary bumps, FAANG placements.
They respect competence. They distrust condescension.

The UI should feel like:
A senior engineer built a game to help you get their job.

Not like:
A teacher praising a child for finishing homework.

---

# Dark Theme Mandate

Primary background: #0F1117
Primary surface:    #1A1D27
Primary accent:     #7C6EF5 (electric purple)

NEVER use:
- White or light grey backgrounds
- Duolingo green (#58CC02) as the primary CTA color
- Fredoka font outside companion dialogue

ALWAYS use:
- JetBrains Mono for all numbers (XP, gems, streak, levels, %scores)
- Plus Jakarta Sans for display headings
- #E8EAED for primary text on dark backgrounds

---

# What is and is not Duolingo

Prepio is in the same genre: gamified daily-habit learning.
Prepio is NOT the same product: children's language app.

Structural patterns we share (intentionally):
- Bottom tab navigation
- Companion with speech bubble
- XP, gems, streak mechanics
- Progress rings and bars

Patterns we do NOT share:
- Lime green as primary CTA
- White/pastel backgrounds
- Childish rounded display fonts
- Overly enthusiastic copy ("AMAZING! YOU DID IT!")
- Single-column centered form layouts
