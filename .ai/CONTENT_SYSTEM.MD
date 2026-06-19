# CONTENT_SYSTEM.md

# AGENT INSTRUCTIONS

This document defines how learning, progression, readiness, skills, and content work inside Prepio.

Before creating:

* Questions
* Skills
* Worlds
* Nodes
* Readiness logic
* Evaluation logic

Read this document completely.

Never bypass the structures defined here.

---

# Purpose

The content system is the core intelligence of Prepio.

UI can change.

Backend can change.

Architecture can evolve.

The content system should remain stable.

This document defines:

* Learning structure
* Skill structure
* Progression structure
* Readiness structure
* Company preparation structure

---

# Core Principle

Questions are not the product.

Skills are the product.

Users do not need:

1000 questions.

Users need:

Mastery of skills.

Questions exist only to measure and improve skills.

---

# Learning Hierarchy

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

Skill Mastery

↓

Readiness

Never bypass this structure.

---

# Example

Foundation Forest

↓

Arrays Basics

↓

Arrays

↓

Question Pool

↓

Two Sum

Contains Duplicate

Best Time To Buy Stock

Move Zeroes

↓

Skill Score

↓

Google Readiness

---

# Worlds

Worlds are progression destinations.

Worlds should feel meaningful.

Worlds should not be categories.

Bad:

DSA

Good:

Foundation Forest

---

# Initial Worlds

Foundation Forest

Problem Solving Canyon

System Design City

Google Mountain

Amazon Jungle

Meta City

Uber Highway

Behavioral Valley

---

# Nodes

Nodes represent progression milestones.

A node teaches or evaluates one or more skills.

Examples:

Arrays Basics

Hash Maps

Sliding Window

Graphs

Behavioral Communication

Node progression should feel meaningful.

---

# Skills

Skills are the foundation of progression.

Everything eventually maps to skills.

---

# Initial Skill Tree

Programming Fundamentals

Arrays

Strings

Hash Maps

Linked Lists

Stacks

Queues

Trees

Graphs

Heaps

Recursion

Greedy

Dynamic Programming

Binary Search

System Design

Low Level Design

Behavioral

Communication

Debugging

Problem Solving

---

# Subskills

Skills may contain subskills.

Example:

Arrays

├── Traversal

├── Two Pointers

├── Sliding Window

├── Prefix Sum

└── Hash Based Lookup

---

# Skill Mastery

Every user has mastery values.

Range:

0 - 100

Example:

Arrays 84

Trees 63

Graphs 41

DP 22

---

Mastery should improve gradually.

Never jump dramatically.

---

# Question Model

Every question must contain:

Question

Difficulty

Skill

Subskill

Companies

Estimated Time

Readiness Weight

Explanation

Hints

Solution

Evaluation Type

---

# Difficulty

Current levels:

Easy

Medium

Hard

---

Future:

Easy 1-5

Medium 1-5

Hard 1-5

Do not implement until needed.

---

# Question Pools

Questions belong to pools.

Example:

Arrays Pool

Contains:

Two Sum

Contains Duplicate

Move Zeroes

Best Time To Buy Stock

Merge Sorted Array

---

Nodes reference pools.

Not individual questions.

This prevents future rewrites.

---

# Question Variants

Question variants should exist.

Example:

Two Sum

↓

Variant A

Variant B

Variant C

Users should practice concepts.

Not memorize answers.

---

# Evaluation Types

Current:

Multiple Choice

Coding

System Design

Behavioral

---

Future:

Voice

Architecture Review

Mock Interview

Do not implement yet.

---

# Readiness System

Readiness is the primary metric.

Not XP.

Not levels.

Not streaks.

---

Readiness answers:

Can this user likely pass interviews?

---

# Readiness Layers

Question

↓

Skill Score

↓

Company Readiness

↓

Overall Readiness

---

# Example

Question Result

↓

Arrays +2

Hash Maps +1

↓

Arrays Mastery = 85

↓

Google Readiness = 74

---

# Company Readiness

Company readiness is derived from skills.

Not question counts.

---

Example

Google

Arrays

Trees

Graphs

DP

System Design

Behavioral

---

Amazon

Arrays

Behavioral

Leadership

LLD

Problem Solving

---

Meta

Coding

Communication

System Design

Behavioral

---

Each company should have a configurable weighting system.

Never hardcode company readiness.

---

# Readiness Transparency

Users should understand:

Why readiness changed.

Bad:

Google Readiness +2

Good:

Google Readiness +2

Arrays +1

Graphs +1

---

# Content Authoring Rules

Every question must:

Map to at least one skill

Map to at least one subskill

Have a difficulty

Have a solution

Have an explanation

Have readiness weights

Have company tags

Questions missing metadata should not be published.

---

# Content Lifecycle

Draft

↓

Review

↓

Published

↓

Deprecated

Questions should not go directly to production.

---

# Future AI Content

AI content generation is not currently a priority.

When introduced:

AI should generate:

Draft Questions

Draft Variants

Draft Explanations

Human review remains required.

---

# Progression Rules

Question completion should never directly unlock worlds.

Bad:

Question Complete

↓

World Unlock

Good:

Question Complete

↓

Skill Mastery

↓

Node Progress

↓

World Progress

↓

Unlock

Progression should feel earned.

---

# Content Expansion Rule

Before adding new questions ask:

Does this skill already have enough coverage?

If no:

Improve coverage.

If yes:

Add new content.

Depth is more valuable than volume.

---

# Anti Patterns

Forbidden:

Questions without skills

Questions without explanations

Questions without company mappings

Question-based progression

Hardcoded readiness

Hardcoded unlocks

Question count as primary metric

XP as readiness

---

# Long-Term Goal

The final content system should be able to answer:

What is this user good at?

What is this user weak at?

Why is their readiness changing?

What should they practice next?

Which company are they closest to?

What skill is preventing them from progressing?

If the system cannot answer these questions:

The implementation is incomplete.
