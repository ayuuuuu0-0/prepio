# Skill Graph Proposal — Prepio Phase A Design

**Date:** 2026-06-09  
**Status:** Design only — do not implement  
**Reference:** `.ai/CONTENT_SYSTEM.md`, `.ai/EXECUTION.md` (A1)

---

## 1. Purpose

Define the complete skill graph that underpins journey progression, question mapping, and readiness calculation. Skills are the **canonical unit of learning** — not questions, not round types, not company tags.

---

## 2. Design Principles

1. **Skills are stable; questions are ephemeral** — add questions without restructuring the graph.
2. **Subskills enable granularity** — readiness explainability requires subskill-level tracking.
3. **Categories group for UI only** — categories do not affect readiness formulas directly.
4. **Prerequisites form a DAG** — no circular dependencies.
5. **Company readiness reads skills, never questions** — weights live in `company_skill_weights`.
6. **Phase A scope** — full taxonomy defined; only Foundation Forest skills need pools populated in Phase B.

---

## 3. Skill Categories

Categories organize the skill tree for navigation and world theming. A skill belongs to exactly one category.

| Slug | Name | Description |
|------|------|-------------|
| `programming-fundamentals` | Programming Fundamentals | Core CS building blocks |
| `data-structures` | Data Structures | Classic DS patterns |
| `algorithms` | Algorithms | Algorithmic paradigms |
| `system-design` | System Design | Distributed systems & architecture |
| `low-level-design` | Low Level Design | OOP, patterns, class design |
| `behavioral` | Behavioral | STAR, leadership, situational |
| `communication` | Communication | Clarity, structure, technical writing |
| `problem-solving` | Problem Solving | Meta-skill: approach, debugging, tradeoffs |

---

## 4. Skills (Top Level)

Aligned with `.ai/CONTENT_SYSTEM.md` initial skill tree.

### Programming Fundamentals (`programming-fundamentals`)

| Slug | Name | Prerequisites |
|------|------|---------------|
| `programming-fundamentals` | Programming Fundamentals | — |

### Data Structures (`data-structures`)

| Slug | Name | Prerequisites |
|------|------|---------------|
| `arrays` | Arrays | `programming-fundamentals` |
| `strings` | Strings | `programming-fundamentals` |
| `hash-maps` | Hash Maps | `arrays` |
| `linked-lists` | Linked Lists | `arrays` |
| `stacks` | Stacks | `arrays`, `linked-lists` |
| `queues` | Queues | `arrays`, `linked-lists` |
| `trees` | Trees | `arrays`, `recursion` |
| `graphs` | Graphs | `trees` |
| `heaps` | Heaps | `trees`, `arrays` |

### Algorithms (`algorithms`)

| Slug | Name | Prerequisites |
|------|------|---------------|
| `recursion` | Recursion | `programming-fundamentals` |
| `binary-search` | Binary Search | `arrays`, `recursion` |
| `greedy` | Greedy | `arrays`, `problem-solving` |
| `dynamic-programming` | Dynamic Programming | `recursion`, `arrays`, `problem-solving` |
| `sliding-window` | Sliding Window | `arrays`, `strings`, `two-pointers` |

### System Design (`system-design`)

| Slug | Name | Prerequisites |
|------|------|---------------|
| `system-design-fundamentals` | System Design Fundamentals | `programming-fundamentals` |
| `system-design-scaling` | Scaling & Performance | `system-design-fundamentals` |
| `system-design-data` | Data Modeling & Storage | `system-design-fundamentals`, `hash-maps` |
| `system-design-reliability` | Reliability & Fault Tolerance | `system-design-fundamentals` |

### Low Level Design (`low-level-design`)

| Slug | Name | Prerequisites |
|------|------|---------------|
| `lld-fundamentals` | LLD Fundamentals | `programming-fundamentals` |
| `lld-patterns` | Design Patterns | `lld-fundamentals` |
| `lld-oop` | OOP & SOLID | `lld-fundamentals` |

### Behavioral (`behavioral`)

| Slug | Name | Prerequisites |
|------|------|---------------|
| `behavioral-star` | STAR Framework | — |
| `behavioral-leadership` | Leadership & Ownership | `behavioral-star` |
| `behavioral-conflict` | Conflict & Failure Stories | `behavioral-star` |

### Communication (`communication`)

| Slug | Name | Prerequisites |
|------|------|---------------|
| `communication-clarity` | Technical Clarity | — |
| `communication-structure` | Structured Responses | `communication-clarity` |

### Problem Solving (`problem-solving`)

| Slug | Name | Prerequisites |
|------|------|---------------|
| `problem-solving` | Problem Solving | `programming-fundamentals` |
| `debugging` | Debugging | `problem-solving` |

**Total top-level skills:** 28 (expandable)

---

## 5. Subskills

Subskills are the **atomic mastery unit**. User mastery is tracked at subskill level; skill mastery rolls up.

### Arrays (`arrays`)

| Slug | Name |
|------|------|
| `arrays-traversal` | Traversal |
| `arrays-two-pointers` | Two Pointers |
| `arrays-sliding-window` | Sliding Window |
| `arrays-prefix-sum` | Prefix Sum |
| `arrays-hash-lookup` | Hash-Based Lookup |

### Strings (`strings`)

| Slug | Name |
|------|------|
| `strings-manipulation` | Manipulation |
| `strings-pattern-matching` | Pattern Matching |
| `strings-two-pointers` | Two Pointers |
| `strings-sliding-window` | Sliding Window |

### Hash Maps (`hash-maps`)

| Slug | Name |
|------|------|
| `hash-maps-frequency` | Frequency Counting |
| `hash-maps-lookup` | Lookup & Deduplication |
| `hash-maps-nested` | Nested / Composite Keys |

### Trees (`trees`)

| Slug | Name |
|------|------|
| `trees-traversal-dfs` | DFS Traversal |
| `trees-traversal-bfs` | BFS Traversal |
| `trees-bst` | BST Operations |
| `trees-recursion` | Recursive Tree Problems |

### Graphs (`graphs`)

| Slug | Name |
|------|------|
| `graphs-bfs-dfs` | BFS / DFS |
| `graphs-shortest-path` | Shortest Path |
| `graphs-topological` | Topological Sort |
| `graphs-union-find` | Union Find |

### Dynamic Programming (`dynamic-programming`)

| Slug | Name |
|------|------|
| `dp-1d` | 1D DP |
| `dp-2d` | 2D DP |
| `dp-state-machine` | State Machine |
| `dp-interval` | Interval DP |

### System Design (`system-design-fundamentals` and children)

| Skill | Subskills |
|-------|-----------|
| `system-design-fundamentals` | Requirements, API Design, Estimation |
| `system-design-scaling` | Caching, Sharding, Load Balancing |
| `system-design-data` | SQL vs NoSQL, Indexing, Consistency |
| `system-design-reliability` | Replication, Circuit Breakers, Idempotency |

### Behavioral

| Skill | Subskills |
|-------|-----------|
| `behavioral-star` | Situation, Task, Action, Result |
| `behavioral-leadership` | Ownership, Influence, Mentorship |
| `behavioral-conflict` | Disagreement, Failure, Recovery |

### LLD

| Skill | Subskills |
|-------|-----------|
| `lld-oop` | Encapsulation, Inheritance, Composition |
| `lld-patterns` | Factory, Observer, Strategy, Singleton |

**Phase A seed:** Full subskill tree defined in DB; **content populated for Foundation Forest only** (arrays, strings, hash-maps, trees subskills).

---

## 6. Relationships

### 6.1 Skill → Subskill (1:N)

Every subskill belongs to exactly one parent skill.

### 6.2 Skill → Skill Prerequisites (M:N)

Junction table `skill_prerequisites(skill_id, prerequisite_skill_id)`.

Used for:
- Journey node unlock ordering
- UI "recommended path"
- **Not** used in readiness formula directly

### 6.3 Question → Skill / Subskill (M:N)

`question_skills(question_id, skill_id, subskill_id, weight)`

- A question may map to multiple skills (e.g., DP + arrays)
- `weight` (0.0–1.0) splits contribution across skills; must sum to 1.0 per question
- `subskill_id` required for approved questions

### 6.4 Skill → Question Pool (1:N)

`question_pools(skill_id, subskill_id nullable, difficulty, slug)`

A pool is scoped to a skill (and optionally subskill + difficulty band).

### 6.5 Company → Skill (M:N weighted)

`company_skill_weights(company, skill_id, weight)`

Weights are **relative importance** (0–100), normalized per company at computation time.

### 6.6 Node → Skill (M:N)

`node_skills(node_id, skill_id, is_primary)`

Journey nodes declare which skills they teach or evaluate.

---

## 7. Suggested Database Schema

```sql
-- Categories
CREATE TABLE skill_categories (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug        TEXT NOT NULL UNIQUE,
    name        TEXT NOT NULL,
    sort_order  INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Skills
CREATE TABLE skills (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id UUID NOT NULL REFERENCES skill_categories(id),
    slug        TEXT NOT NULL UNIQUE,
    name        TEXT NOT NULL,
    description TEXT,
    sort_order  INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Subskills
CREATE TABLE subskills (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    skill_id    UUID NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    slug        TEXT NOT NULL,
    name        TEXT NOT NULL,
    sort_order  INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (skill_id, slug)
);

-- Prerequisites
CREATE TABLE skill_prerequisites (
    skill_id              UUID NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    prerequisite_skill_id UUID NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    PRIMARY KEY (skill_id, prerequisite_skill_id),
    CHECK (skill_id != prerequisite_skill_id)
);

-- Question mapping
CREATE TABLE question_skills (
    question_id  UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    skill_id     UUID NOT NULL REFERENCES skills(id),
    subskill_id  UUID NOT NULL REFERENCES subskills(id),
    weight       NUMERIC(4,3) NOT NULL DEFAULT 1.0,
    PRIMARY KEY (question_id, skill_id, subskill_id)
);

-- User mastery (Progress domain)
CREATE TABLE user_subskill_scores (
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    subskill_id  UUID NOT NULL REFERENCES subskills(id),
    mastery      INT NOT NULL DEFAULT 0 CHECK (mastery >= 0 AND mastery <= 100),
    attempts     INT NOT NULL DEFAULT 0,
    last_practiced_at TIMESTAMPTZ,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, subskill_id)
);

-- Rolled-up skill mastery (materialized or computed view)
-- user_skill_scores(user_id, skill_id, mastery) — avg of subskills weighted equally

-- Company weights
CREATE TABLE company_skill_weights (
    company   TEXT NOT NULL,
    skill_id  UUID NOT NULL REFERENCES skills(id),
    weight    INT NOT NULL CHECK (weight >= 0 AND weight <= 100),
    PRIMARY KEY (company, skill_id)
);
```

---

## 8. Suggested APIs

All under `/api/v1/`. Read APIs are public (authenticated). Write APIs internal-only in Phase A.

### Read — Content / Progress

| Method | Path | Description |
|--------|------|-------------|
| GET | `/skills` | List categories with skills (tree) |
| GET | `/skills/{slug}` | Skill detail + subskills |
| GET | `/skills/{slug}/pools` | Question pools for skill |
| GET | `/progress/skills` | User mastery for all skills |
| GET | `/progress/skills/{slug}` | User mastery for skill + subskills |
| GET | `/progress/readiness` | Overall + per-target-company |
| GET | `/progress/readiness/{company}/gaps` | Weakest weighted skills |

### Write — Internal Content Management (A5)

| Method | Path | Description |
|--------|------|-------------|
| POST | `/internal/skills` | Create skill |
| POST | `/internal/subskills` | Create subskill |
| POST | `/internal/questions/{id}/skills` | Attach skill mapping |
| PUT | `/internal/company-skill-weights/{company}` | Bulk update weights |

---

## 9. Question-to-Skill Mapping Strategy

### 9.1 Authoring Workflow

1. Author selects **primary skill** + **primary subskill**
2. Optionally adds secondary skill mappings with weights
3. System validates: weights sum to 1.0, subskill belongs to skill
4. Question cannot move to `approved` without valid mapping

### 9.2 Backfill Heuristic (Existing 12 Questions)

| round_type | Provisional Skill | Notes |
|------------|-------------------|-------|
| `dsa` | Infer from answer_guide keywords | e.g., "two heap" → `heaps`, "median" → `heaps` |
| `system_design` | `system-design-fundamentals` | Split subskill by topic |
| `behavioral` | `behavioral-star` | |
| `lld` | `lld-fundamentals` | |
| `fundamentals` | `programming-fundamentals` | |
| `aptitude` | `problem-solving` | |

Manual review required for all backfilled mappings before Phase A sign-off.

### 9.3 Evaluation → Mastery Update

On `question.completed`:

```
for each (skill, subskill, weight) in question_skills:
    delta = (score / 100) * readiness_weight * weight * difficulty_multiplier
    user_subskill_scores[subskill].mastery += delta  // capped, smoothed
    user_subskill_scores[subskill].attempts += 1
```

Skill mastery = weighted average of subskill masteries.

### 9.4 Coverage Tracking (Analytics — Phase A late)

Query: skills with `< N` approved questions at each difficulty → content gap report.

---

## 10. Company Skill Weight Profiles (Seed Data)

Normalized weights (sum = 100 per company). Stored in DB, not code.

### Google

| Skill | Weight |
|-------|--------|
| arrays | 15 |
| trees | 15 |
| graphs | 15 |
| dynamic-programming | 20 |
| system-design-scaling | 15 |
| behavioral-star | 10 |
| communication-clarity | 10 |

### Amazon

| Skill | Weight |
|-------|--------|
| arrays | 15 |
| behavioral-leadership | 20 |
| behavioral-star | 15 |
| lld-patterns | 15 |
| problem-solving | 20 |
| system-design-data | 15 |

### Meta

| Skill | Weight |
|-------|--------|
| arrays | 15 |
| dynamic-programming | 15 |
| system-design-scaling | 20 |
| behavioral-star | 15 |
| communication-structure | 15 |
| graphs | 20 |

### Uber

| Skill | Weight |
|-------|--------|
| system-design-scaling | 25 |
| graphs | 20 |
| arrays | 15 |
| behavioral-star | 15 |
| problem-solving | 25 |

### Atlassian

| Skill | Weight |
|-------|--------|
| behavioral-star | 25 |
| lld-oop | 25 |
| arrays | 15 |
| system-design-fundamentals | 20 |
| communication-clarity | 15 |

---

## 11. Foundation Forest — Initial Skill Bindings

Maps existing journey nodes to skills (for A2).

| Node | Primary Skill | Subskill Focus |
|------|---------------|----------------|
| Arrays Basics | `arrays` | `arrays-traversal`, `arrays-two-pointers` |
| String Patterns | `strings` | `strings-pattern-matching` |
| Hash Maps | `hash-maps` | `hash-maps-frequency` |
| Tree Traversal | `trees` | `trees-traversal-dfs`, `trees-traversal-bfs` |
| Forest Boss | `arrays`, `hash-maps`, `trees` | Mixed assessment |

---

## 12. Phase A vs Phase B Scope

| Item | Phase A | Phase B |
|------|---------|---------|
| Full skill taxonomy in DB | ✅ | |
| Subskills for all skills | ✅ schema | Content |
| question_skills for seed questions | ✅ | 100–150 questions |
| company_skill_weights seed | ✅ | Tune from data |
| user_subskill_scores | ✅ | |
| Skill gap analytics | Basic API | Rich UI |

---

## 13. Open Decisions (For Review)

1. **Mastery at subskill vs skill level only?** — Recommend subskill storage, skill rollup.
2. **Decay/recency in Phase A?** — Recommend no decay in A; add in Phase C if needed.
3. **Separate skill service vs Content package?** — Recommend Content owns graph; Progress owns scores.
4. **`sliding-window` as skill vs subskill?** — Recommend subskill under arrays/strings; remove top-level `sliding-window` skill to avoid duplication.

---

*Design document only. No implementation.*
