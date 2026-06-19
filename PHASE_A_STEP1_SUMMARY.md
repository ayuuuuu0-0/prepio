# Phase A Step 1 Summary — A1 Skill Graph + A3 Question Schema Upgrade

**Date:** 2026-06-09  
**Scope:** Foundation only (no pools, journey changes, readiness V2, or CMS)

---

## What Was Completed

### A1 — Skill Graph

- Created `skill_categories`, `skills`, `subskills`, and `question_skills` tables
- Seeded 8 categories, 29 skills, and 68 subskills from `SKILL_GRAPH_PROPOSAL.md`
- Backfilled all 12 seed questions with `question_skills` mappings and weights summing to 1.0
- Added read APIs:
  - `GET /api/v1/skills` — full skill tree (categories → skills → subskills)
  - `GET /api/v1/skills/{slug}` — single skill detail
  - `GET /api/v1/questions/{id}/skills` — question-to-skill mappings
- Gateway routes `/skills` and `/skills/*` to question service

### A3 — Question Schema Upgrade

Added columns to `questions` (all existing columns preserved):

| Column | Type | Default |
|--------|------|---------|
| `evaluation_type` | TEXT | backfilled from `round_type` |
| `explanation` | TEXT | nullable |
| `hints` | JSONB | `[]` |
| `solution` | TEXT | nullable |
| `readiness_weight` | NUMERIC(3,2) | 1.00 (by difficulty: 0.80/1.00/1.20) |
| `estimated_time` | INT (minutes) | 10 (by difficulty: 8/15/25) |

Seed questions enriched with explanation, hints, and solution in migration `000028`.

---

## Files Changed

### Migrations (new)

| File | Purpose |
|------|---------|
| `migrations/000025_create_skill_graph.up.sql` | skill_categories, skills, subskills, question_skills |
| `migrations/000025_create_skill_graph.down.sql` | Rollback skill graph tables |
| `migrations/000026_seed_skill_graph.up.sql` | Seed categories, skills, subskills |
| `migrations/000026_seed_skill_graph.down.sql` | Clear seed data |
| `migrations/000027_extend_questions.up.sql` | Add A3 columns + backfill evaluation_type/defaults |
| `migrations/000027_extend_questions.down.sql` | Drop A3 columns |
| `migrations/000028_backfill_question_skills.up.sql` | Map 12 seed questions to skills + enrich metadata |
| `migrations/000028_backfill_question_skills.down.sql` | Remove mappings and metadata |

### Go — Constants & Config

| File | Purpose |
|------|---------|
| `constants/skills.go` | Evaluation type constants, defaults |
| `constants/errors.go` | Added `ErrSkillNotFound` |
| `config/readiness.go` | Default readiness weights and estimated times by difficulty |

### Go — Question Service

| File | Purpose |
|------|---------|
| `services/question/internal/store/skills.go` | SkillStore repository |
| `services/question/internal/store/skills_test.go` | Store integration tests |
| `services/question/internal/service/skill.go` | SkillService |
| `services/question/internal/service/skill_test.go` | Service tests |
| `services/question/internal/service/errors.go` | Added `ErrSkillNotFound` |
| `services/question/internal/dto/skill.go` | Skill API DTOs |
| `services/question/internal/handler/skill.go` | Skill HTTP handlers |
| `services/question/cmd/main.go` | Wire skill routes |
| `services/question/smoke/skills_test.go` | HTTP smoke tests |

### Go — Gateway

| File | Purpose |
|------|---------|
| `services/gateway/cmd/main.go` | Proxy `/skills` routes to question service |

---

## New Schema

```
skill_categories
  id, slug, name, sort_order, created_at, updated_at

skills
  id, category_id → skill_categories, slug, name, description,
  sort_order, created_at, updated_at

subskills
  id, skill_id → skills, slug, name, sort_order, created_at, updated_at
  UNIQUE (skill_id, slug)

question_skills
  question_id → questions, skill_id → skills, subskill_id → subskills,
  weight (0–1], PRIMARY KEY (question_id, skill_id, subskill_id)

questions (extended)
  + evaluation_type, explanation, hints, solution,
    readiness_weight, estimated_time
```

**Stable seed UUID prefixes:**
- Categories: `c1000000-...`
- Skills: `b2000001-...`
- Subskills: `c3000001-...` through `c3000006-...`

---

## Migration Strategy

### Apply order

1. `000025` — Create tables (additive, no impact on running app)
2. `000026` — Seed skill graph (additive)
3. `000027` — Extend questions with nullable/defaulted columns (non-breaking)
4. `000028` — Backfill question_skills + enrich seed question metadata

### Backfill approach for existing questions

| Question | Skills Mapped |
|----------|---------------|
| Two Sum | arrays (0.6) + hash-maps (0.4) |
| Linked list cycle | linked-lists |
| URL shortener | system-design-fundamentals |
| Max depth tree | trees (DFS) |
| Parking lot | lld-fundamentals |
| Longest substring | strings (sliding window) |
| Aptitude widgets | problem-solving |
| Process vs thread | programming-fundamentals |
| Behavioral deadline | behavioral-star |
| Rate limiter | system-design-scaling + fundamentals |
| Median stream | heaps |
| Food delivery | system-design-scaling + fundamentals |

**Heuristic for future questions:** map by `round_type` first, refine by answer_guide keywords. New questions created at runtime (tests) are not auto-mapped — CMS validation comes in Step 2.

### Rollback

Run down migrations in reverse order. Down `000025` drops skill tables (cascades `question_skills`). Down `000027` removes new question columns.

---

## What Was NOT Changed (by design)

- Daily paper selection logic
- Journey node unlock / index mapping
- Readiness computation (gateway + question stats)
- Submit flow / evaluator / rewards
- Kafka events (no skill fields yet)
- Question pools (A2)
- Content management API (A5)
- `user_subskill_scores` / company_skill_weights (A4)

Existing API response shapes unchanged — new fields are not exposed on daily paper or submit responses yet.

---

## Tests Added

| Test | Coverage |
|------|----------|
| `store/skills_test.go` | Categories, skill lookup, question_skills, content metadata |
| `service/skill_test.go` | Skill tree, not-found, question skill list |
| `smoke/skills_test.go` | GET /skills, /skills/arrays, /questions/{id}/skills |

All pass with embedded Postgres via `testdb.Migrate`.

---

## Risks

| Risk | Severity | Mitigation |
|------|----------|------------|
| New questions without skill mappings | Medium | Only affects unmapped questions; existing flows ignore `question_skills`. Step 2 should enforce mapping on approve. |
| `sliding-window` exists as top-level skill and as subskill | Low | Documented open decision; subskills under arrays/strings are canonical for mastery. |
| Skill prerequisites not seeded | Low | `skill_prerequisites` table deferred; unlock rules use node bindings in A2. |
| Gateway must restart for /skills routes | Low | Standard deploy step |
| Backfill mappings are heuristic | Medium | All 12 seed questions manually mapped; review before readiness V2 |

---

## Remaining Work Before Step 2 (A2 — Content Architecture)

Per `PHASE_A_IMPLEMENTATION_PLAN.md`:

1. **Question pools** — `question_pools`, `pool_questions`, `user_pool_progress`
2. **Journey node bindings** — `node_skills`, `node_pools`; backfill Foundation Forest
3. **Pool selection API** — `GET /journey/nodes/{id}/question` behind feature flag
4. **Decouple journey from daily paper index** — Journey V2 (flag-gated)
5. **Extend Kafka events** — include skill IDs on `question.answered` (prep for A4)
6. **Validation** — block `approved` status without `question_skills` row (app-level)
7. **Optional:** `skill_prerequisites` table for unlock ordering UI

### Also before A4 (Readiness V2)

- `company_skill_weights` table + seed
- `user_subskill_scores` table
- Progress service readiness engine
- Remove readiness from gateway

---

## Verification Checklist

- [x] `go build ./...` passes
- [x] `go vet ./...` clean
- [x] `go test ./...` passes
- [x] Migrations 000025–000028 apply on fresh DB
- [x] All 12 approved seed questions have `question_skills` rows
- [x] Every skill has a category
- [x] Existing daily paper + submit smoke tests pass
- [x] No readiness or journey logic modified

---

*Step 1 complete. Proceed to Step 2 (question pools + journey bindings) after review.*
