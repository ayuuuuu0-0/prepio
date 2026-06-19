# Content Audit — Prepio Codebase vs `.ai/CONTENT_SYSTEM.md`

**Date:** 2026-06-09  
**Scope:** Current implementation as of migration `000024`  
**Reference:** `.ai/PRODUCT.md`, `.ai/ARCHITECTURE.md`, `.ai/CONTENT_SYSTEM.md`, `.ai/EXECUTION.md`

---

## Executive Summary

The current content model is a **flat question bank** with company tags and a **daily paper** delivery mechanism. Journey exists visually but is **not content-driven** — nodes are cosmetic labels mapped by array index to today's daily paper questions. There is **no skill graph**, **no question pools**, and **no readiness-weight metadata** on questions. Content is primarily **seeded via SQL migrations**, which conflicts with Phase A goal A5 (content management without code changes).

The gap between current state and target hierarchy:

```
Current:  Question → Daily Paper → (index) → Journey Node label
Target:   World → Node → Skill → Question Pool → Question → Skill Mastery → Readiness
```

---

## 1. Current Question Schema

### Database: `questions` (migration `000006`)

| Column | Type | Notes |
|--------|------|-------|
| `id` | UUID | Primary key |
| `body` | TEXT | Question text |
| `round_type` | ENUM-like TEXT | `dsa`, `system_design`, `lld`, `aptitude`, `fundamentals`, `behavioral` |
| `difficulty` | TEXT | `easy`, `medium`, `hard` |
| `answer_guide` | TEXT | Evaluator rubric; some rows use `concepts:term1\|term2\|...` prefix (migration `000022`) |
| `status` | TEXT | `pending`, `approved`, `retired` |
| `is_weekend` | BOOLEAN | Weekend-only flag |
| `source` | TEXT | `manual`, `ai_generated`, `scraped` |

### Related: `question_tags` (migration `000007`)

| Column | Type | Notes |
|--------|------|-------|
| `question_id` | UUID | FK → questions |
| `company` | TEXT | Company slug only; no weighting |

### Application struct: `store.Question`

```go
// services/question/internal/store/questions.go
ID, Body, RoundType, Difficulty, AnswerGuide, Status, IsWeekend, CompanyTags[]
```

### Missing vs `.ai/CONTENT_SYSTEM.md` required fields

| Required Field | Present? |
|----------------|----------|
| Skill | **No** |
| Subskill | **No** |
| Readiness Weight | **No** |
| Estimated Time | **No** |
| Explanation (user-facing) | **No** (only `answer_guide` for evaluator) |
| Hints | **No** |
| Solution | **No** (embedded in `answer_guide` informally) |
| Evaluation Type | **Implicit** (free-text only; no enum column) |
| Question Variants | **No** |
| Question Pool membership | **No** |

---

## 2. Current Question-Related Tables

| Table | Purpose | Domain alignment |
|-------|---------|------------------|
| `questions` | Question bank | Content (partial) |
| `question_tags` | Company tags per question | Content (partial) |
| `daily_papers` | Per-user per-day session | **Misplaced** — session/delivery, not content |
| `daily_paper_questions` | Questions in a daily paper | **Misplaced** — couples progression to daily batch |
| `user_question_history` | Answer records | Progress/Analytics hybrid |
| `worlds` | Journey worlds | Journey |
| `journey_nodes` | Nodes within a world | Journey (no skill/pool FK) |
| `user_journey_progress` | Per-user node status | Journey (underused) |

**Tables that do NOT exist but are required by `.ai/` docs:**

- `skills`, `subskills`, `skill_categories`
- `question_skills` (M:N mapping with weights)
- `question_pools`, `pool_questions`
- `node_skills`, `node_pools` (node → skill/pool bindings)
- `user_skill_scores` / `user_skill_mastery`
- `company_skill_weights` (company readiness weighting)
- `question_hints`, `question_solutions` (or JSONB columns)
- Content authoring / review queue tables beyond `status` enum

---

## 3. Current Content Flow

```
┌─────────────────────────────────────────────────────────────────┐
│ SEEDING (migrations 000019, 000022, 000024)                     │
│   INSERT questions + question_tags + worlds + journey_nodes     │
└────────────────────────────┬────────────────────────────────────┘
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│ DAILY PAPER GENERATION (Question Service)                       │
│   GET /api/v1/questions/daily                                   │
│   1. Check existing paper for user+date                         │
│   2. Select by difficulty (from user level) + unseen + random   │
│   3. Create daily_papers + daily_paper_questions                  │
│   Note: target company priority documented but NOT implemented  │
└────────────────────────────┬────────────────────────────────────┘
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│ CHALLENGE / SUBMIT                                              │
│   POST /api/v1/questions/{id}/submit                            │
│   1. Validate question in session                               │
│   2. Evaluate via answer_guide concept matching                 │
│   3. Insert user_question_history (correct, score)              │
│   4. Emit question.answered → streak, progress consumers        │
└────────────────────────────┬────────────────────────────────────┘
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│ READINESS (split across Question + Gateway — ownership violation)│
│   Question: CompanyPerformanceByUser, AvgScoreByUser            │
│   Gateway: computeReadiness() for dashboard cards               │
└─────────────────────────────────────────────────────────────────┘
```

**Key observation:** Progression is driven by **daily paper batch selection**, not by journey nodes or skills. Journey is a **read-model overlay** on the daily paper.

---

## 4. Current Journey Content Mapping

### Schema (`000024`)

- One world seeded: `foundation-forest` / "Foundation Forest"
- Five nodes: Arrays Basics, String Patterns, Hash Maps, Tree Traversal, Forest Boss
- Nodes have: `label`, `node_type` (`lesson` | `boss`), `sort_order`
- **No FK** to skills, question pools, or individual questions

### Runtime mapping (`services/question/internal/service/journey.go`)

```text
node[i].status  ←→  daily_paper.questions[i]  (by array index)
node[i].question_id = paper.questions[i].id
```

| Journey Node Label | Mapped To |
|--------------------|-----------|
| Arrays Basics | Daily paper question index 0 |
| String Patterns | Daily paper question index 1 |
| Hash Maps | Daily paper question index 2 |
| Tree Traversal | Daily paper question index 3 |
| Forest Boss | Unlocked when all daily questions answered |

**Problems:**

1. Node labels are **cosmetic** — they do not filter or select questions by topic.
2. A random daily paper can assign a system design question to "Arrays Basics".
3. Only `foundation-forest` is loaded (`const foundationWorldSlug = "foundation-forest"`).
4. `user_journey_progress` is written on completion but **not read** for unlock logic.

---

## 5. Existing Content Management Capabilities

| Capability | Status |
|------------|--------|
| Create question via API | **No** — read/submit only |
| Create question via migration | **Yes** — primary method today |
| Review queue (pending → approved) | **Schema only** — no admin UI or API |
| Edit question content | **No** |
| Manage worlds/nodes | **Migration only** |
| Manage skills/pools | **Does not exist** |
| AI content pipeline | **Not implemented** (AGENTS.md rules exist; no service) |
| Content versioning | **No** |
| Bulk import | **No** |

The only operational content path is: **edit SQL migration → run migrate → restart**.

---

## 6. Hardcoded Content Found in Code

| Location | Hardcoded Content | Severity |
|----------|-------------------|----------|
| `migrations/000019_seed_questions.up.sql` | 12 questions + company tags | Expected for seed; should move to CMS |
| `migrations/000024_create_journey.up.sql` | 1 world, 5 nodes | Expected for seed |
| `migrations/000021_seed_starter_companions.up.sql` | Companion dialogues | Content in migration |
| `services/question/internal/service/journey.go` | `foundationWorldSlug = "foundation-forest"` | **High** — blocks multi-world |
| `constants/companies.go` | `TargetCompanies` list | **Medium** — onboarding whitelist |
| `web/src/lib/api.ts` | `TARGET_COMPANIES` duplicate | **Medium** — frontend copy of backend |
| `mobile/lib/core/config/constants.dart` | `targetCompanies` duplicate | **Medium** |
| `config/rewards.go` | XP/gems by difficulty, top-tier companies | Config (acceptable) but affects content weighting |
| `services/gateway/internal/dashboard/service.go` | `comingSoonQuests()` titles | **Low** — placeholder quests |
| `services/question/internal/service/question.go` | Comment: "priorities 1-2 require target companies; skipped" | **High** — documented gap in selection |

**AGENTS.md conflict:** Rule says "Never write question content inline in code" — satisfied for application code, but **migrations contain full question strings**, which is acceptable for seeds but not for ongoing content ops.

---

## 7. Content Architecture Problems

### P0 — Structural

1. **No skill layer** — cannot answer "what is this user good at?" per CONTENT_SYSTEM.md.
2. **Journey decoupled from content** — nodes are labels, not progression units.
3. **Daily paper is the progression engine** — violates "questions should not directly drive progression."
4. **Question selection ignores target companies** — onboarding collects targets but selection uses difficulty + random only.

### P1 — Metadata

5. **round_type used as proxy for skill** — coarse substitute, not mastery-tracked.
6. **answer_guide doubles as rubric + concepts** — no separate explanation/hints/solution.
7. **No readiness_weight on questions** — company readiness cannot be weighted per skill.

### P2 — Operations

8. **Content lives in migrations** — every new question requires deploy.
9. **No review workflow API** — `pending` status unused in practice.
10. **Duplicate company lists** across constants, web, mobile.

### P3 — Lifecycle

11. **No deprecated content path** — `retired` status exists but no consumer logic.
12. **No question variants** — users may memorize answers.

---

## 8. Recommended Migrations

Each recommendation includes **why**, **impact**, and **migration strategy**.

---

### R1 — Add skill graph tables (Phase A1)

**Why:** Skills are the foundation of readiness and progression per `.ai/ARCHITECTURE.md`. Without them, Readiness V2 is impossible.

**Impact:** Additive only. Existing questions continue to work. New tables: `skill_categories`, `skills`, `subskills`.

**Migration strategy:**
1. New migration `000025_create_skills.up.sql` — create tables, seed initial skill tree from CONTENT_SYSTEM.md.
2. No changes to existing queries until R2.
3. Backfill script (one-time) maps `round_type` → provisional skill IDs for legacy questions.

**Difficulty:** Medium | **Risk:** Low (additive)

---

### R2 — Add question_skills and extend questions (Phase A3)

**Why:** Every question must map to ≥1 skill + subskill per CONTENT_SYSTEM.md authoring rules.

**Impact:** Add columns to `questions`: `estimated_minutes`, `readiness_weight`, `evaluation_type`, `hints` (JSONB), `explanation`, `solution`. Add `question_skills(question_id, skill_id, subskill_id, weight)`.

**Migration strategy:**
1. Add nullable columns first (no breaking change).
2. Backfill from seed data + manual mapping spreadsheet.
3. Add CHECK constraint: approved questions must have ≥1 question_skill row (enforce via app before DB constraint).
4. Split `answer_guide` — keep for evaluator; move user-facing content to new columns gradually.

**Difficulty:** Medium | **Risk:** Low-Medium

---

### R3 — Add question pools (Phase A2)

**Why:** Nodes must reference pools, not individual questions. Enables content expansion without changing progression structure.

**Impact:** New tables: `question_pools`, `pool_questions`, `node_pools` or `node_skills`.

**Migration strategy:**
1. Create pools per initial skill (e.g., `arrays-pool-easy`).
2. Assign existing 12 seed questions to pools based on inferred skill.
3. Do NOT wire pools to daily paper yet — parallel path first.

**Difficulty:** Medium | **Risk:** Low

---

### R4 — Decouple journey nodes from daily paper index (Phase A2)

**Why:** Current index-based mapping produces incorrect UX and blocks real progression.

**Impact:** Add `journey_nodes.skill_id` or `journey_nodes.pool_id`. Change `GetJourney` to derive status from skill mastery / pool completion, not daily paper order.

**Migration strategy:**
1. Add nullable FK columns to `journey_nodes`.
2. Backfill Foundation Forest nodes with skill mappings (Arrays → arrays skill, etc.).
3. Feature flag: `JOURNEY_V2=true` switches read path.
4. Keep daily paper for "Play" tab until node-driven selection replaces it.

**Difficulty:** High | **Risk:** Medium (UX change)

---

### R5 — Add user_skill_scores (Phase A4 prerequisite)

**Why:** Readiness must derive from skill mastery, not company tag hit rate.

**Impact:** New table `user_skill_scores(user_id, skill_id, mastery INT, updated_at)`. Written by Progress service on `question.completed` events.

**Migration strategy:**
1. Add table + indexes.
2. Backfill from `user_question_history` using question_skills mapping (rough initial mastery).
3. Switch readiness computation to read from this table (Progress service).

**Difficulty:** Medium | **Risk:** Medium (readiness numbers will change)

---

### R6 — Add company_skill_weights (Phase A4)

**Why:** Company readiness must be configurable weighting over skills, not hardcoded.

**Impact:** New table `company_skill_weights(company, skill_id, weight)`. Seed Google/Amazon/Meta profiles from CONTENT_SYSTEM.md.

**Migration strategy:**
1. Seed weights via migration (config data, not code).
2. Move readiness computation from Gateway → Progress service.
3. Gateway dashboard calls Progress readiness API.

**Difficulty:** Medium | **Risk:** Low

---

### R7 — Content management API foundation (Phase A5)

**Why:** Stop requiring migrations for every question. EXECUTION.md A5 completion criteria.

**Impact:** Internal admin endpoints or CLI under Content domain (could extend Question service initially with clear boundary).

**Migration strategy:**
1. CRUD for skills, pools, questions (internal auth only).
2. Review queue: list pending, approve, retire.
3. Export/import JSON for bulk content ops.
4. Do NOT build public admin UI in Phase A — API/CLI sufficient.

**Difficulty:** High | **Risk:** Low (internal only)

---

### R8 — Deprecate migration-based content additions

**Why:** Long-term content ops must not require deploys.

**Impact:** Process change. Keep migrations for schema; move question seeds to seed CLI or JSON import.

**Migration strategy:**
1. After R7, freeze new question INSERT migrations.
2. Convert `000019` to importable JSON fixture for dev/test only.
3. Document content authoring workflow in CONTENT_SYSTEM.md.

**Difficulty:** Low | **Risk:** Low

---

## 9. Backward Compatibility Notes

| Change | Breaks existing users? |
|--------|------------------------|
| Skill tables | No |
| question_skills backfill | No — additive metadata |
| Journey V2 logic | **UX change** — node status may differ from today |
| Readiness V2 | **Numbers change** — explain to users via transparency UI |
| Daily paper removal | **Yes** — defer until node/pool selection works |

**Recommendation:** Run daily paper and journey-v2 in parallel during transition. "Play" tab uses daily paper; "Journey" tab uses skill-based status when flag enabled.

---

## 10. Priority Order (maps to EXECUTION.md Phase A)

| Order | Item | EXECUTION task |
|-------|------|----------------|
| 1 | R1 Skill graph | A1 |
| 2 | R2 Question schema upgrade | A3 |
| 3 | R3 Question pools | A2 |
| 4 | R5 user_skill_scores | A4 |
| 5 | R6 company_skill_weights | A4 |
| 6 | R4 Journey decoupling | A2 |
| 7 | R7 Content management API | A5 |
| 8 | R8 Stop migration content | A5 |

---

*This document is analysis only. No code, migrations, or schemas were modified.*
