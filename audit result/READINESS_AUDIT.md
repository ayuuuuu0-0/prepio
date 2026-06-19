# Readiness Audit — Prepio vs Target Model

**Date:** 2026-06-09  
**Reference:** `.ai/CONTENT_SYSTEM.md`, `.ai/ARCHITECTURE.md`, `.ai/EXECUTION.md` (A4)

---

## Executive Summary

Readiness today is a **coarse proxy** derived from **company-tagged question performance** (correct rate + average score). It is **not skill-based**, **not explainable** per skill gap, and **split across three code paths** (Question stats, Gateway dashboard, Submit delta). No tables store readiness scores — they are **computed on read** from `user_question_history` + `question_tags`.

The target model requires:

```
Question → Skill → Company Readiness → Overall Readiness
```

None of the middle layers exist today.

---

## 1. How Readiness Is Currently Calculated

### Path A — Gateway Dashboard (`services/gateway/internal/dashboard/service.go`)

```go
// computeReadiness(companyStats []CompanyPerformance)
// For each target company:
//   readiness = (correctRate * 100 + avgScore) / 2
//   cap at 95
// Overall = average of company readiness values
```

**Inputs:** `GET /api/v1/questions/stats/readiness` from Question service

**Output:** Per-company cards on dashboard + overall percentage

---

### Path B — Question Stats API (`services/question/internal/store/history.go`)

```sql
-- CompanyPerformanceByUser
-- JOIN user_question_history + question_tags
-- GROUP BY company
-- Returns: company, answered, correct, avg_score
```

**Formula:** None at store level — raw aggregates only.

---

### Path C — Submit Response (`services/question/internal/service/question.go`)

```go
// readiness_delta = newOverallAvgScore - oldOverallAvgScore
// Where overall = AvgScoreByUser (all history, not company-specific)
```

**Displayed to user:** Delta on result card after answering.

**Inconsistency:** Dashboard uses company-weighted formula; submit delta uses global average score only.

---

## 2. Tables That Store Readiness

| Table | Stores Readiness? | Notes |
|-------|-------------------|-------|
| `user_question_history` | **Indirect** | Raw answers — source for computation |
| `question_tags` | **Indirect** | Company association only |
| `user_targets` | **No** | Target companies list, not scores |
| `users` | **No** | No readiness column |
| `user_progress` / progress tables | **No** | XP, gems, level only |
| `user_skill_scores` | **Does not exist** | Required by target architecture |
| `user_readiness` | **Does not exist** | Could cache company/overall scores |

**Verdict:** Readiness is **ephemeral** — computed at request time, never persisted.

---

## 3. Services That Calculate Readiness

| Service | Role | Should Own? |
|---------|------|-------------|
| **Gateway** | `computeReadiness()` formula | ❌ No — violates gateway rules |
| **Question** | Aggregates history by company; avg score; submit delta | ❌ Partial — should emit events only |
| **Progress** | Does not calculate readiness | ✅ Should own per `.ai/ARCHITECTURE.md` |
| **Analytics** | N/A | Could expose read APIs |

---

## 4. Inputs That Affect Readiness Today

| Input | Affects Readiness? | Notes |
|-------|---------------------|-------|
| Answer correctness | Yes | Via correct rate |
| Answer score (0-100) | Yes | Via avg_score |
| Company tag on question | Yes | Groups stats per company |
| User target companies | Yes | Dashboard filters to targets |
| Question difficulty | **No** | Not weighted |
| Skill | **No** | Does not exist |
| Recency / decay | **No** | All history weighted equally |
| Question readiness_weight | **No** | Column does not exist |
| Streak / level | **No** | |
| Time since last practice per skill | **No** | |

---

## 5. Target Readiness Model (`.ai/CONTENT_SYSTEM.md`)

### Skill Mastery

Each answer updates mastery for mapped skill(s):

```
mastery_delta = f(score, readiness_weight, difficulty, recency)
```

### Company Readiness

```
company_readiness = Σ (skill_mastery[skill] × company_skill_weight[company][skill])
                    ─────────────────────────────────────────────────────────
                                    Σ weights
```

### Overall Readiness

```
overall_readiness = weighted average across target companies
                  (or max/min depending on product rule — recommend average)
```

### Explainability Requirements

User must see:
- "You're 72% ready for Google"
- "Weakest skills: Dynamic Programming (41%), System Design Scaling (55%)"
- "Answer 3 more Hash Map questions to improve Arrays readiness"

**Current UI shows a single percentage with no skill breakdown.**

---

## 6. Gap Analysis

| ID | Gap | Severity |
|----|-----|----------|
| R1 | No skill mastery layer | **Critical** |
| R2 | No company_skill_weights | **Critical** |
| R3 | Readiness in Gateway | **Critical** |
| R4 | Inconsistent formulas (dashboard vs submit delta) | **High** |
| R5 | No persistence — cannot trend over time | **High** |
| R6 | No recency decay — old answers count equally | **Medium** |
| R7 | No readiness_weight on questions | **High** |
| R8 | Company tags ≠ company readiness model | **Medium** — tags are question metadata, not skill weights |
| R9 | No `readiness.updated` event | **Medium** |
| R10 | No explainability API | **High** |

---

## 7. Architecture Changes Required

### 7.1 New Tables

| Table | Purpose |
|-------|---------|
| `skills` | Skill definitions |
| `company_skill_weights` | Per-company skill importance |
| `question_skills` | Question → skill mapping with weight |
| `user_skill_scores` | `(user_id, skill_id, mastery, attempts, last_practiced_at)` |
| `user_company_readiness` | Optional cache `(user_id, company, score, updated_at)` |
| `user_overall_readiness` | Optional cache `(user_id, score, updated_at)` |

### 7.2 Service Ownership

| Responsibility | Owner |
|----------------|-------|
| Update skill mastery on answer | **Progress** |
| Compute company readiness | **Progress** |
| Compute overall readiness | **Progress** |
| Expose readiness API | **Progress** (or Analytics read replica) |
| Store question→skill mapping | **Content** |
| Emit `skill.progress.updated` | **Progress** |
| Emit `readiness.updated` | **Progress** |

### 7.3 Event Flow

```
question.completed (skill_ids, score, readiness_weight)
    → Progress: update user_skill_scores
    → Progress: recompute company + overall readiness (or async job)
    → Progress: emit skill.progress.updated
    → Progress: emit readiness.updated (if delta > threshold)
    → Notification: "You're now 68% ready for Amazon"
    → Analytics: materialize for dashboards
```

### 7.4 API Changes

| Endpoint | Purpose |
|----------|---------|
| `GET /api/v1/progress/readiness` | Overall + per-company |
| `GET /api/v1/progress/readiness/skills` | Skill breakdown for a company |
| `GET /api/v1/progress/readiness/gaps` | Weakest skills vs target |

Deprecate: `GET /api/v1/questions/stats/readiness` (move to Progress)

---

## 8. Migration Strategy

### Phase 1 — Foundation (no user-visible change)

1. Create `skills`, `question_skills`, `company_skill_weights` tables
2. Backfill question→skill mappings from `round_type` heuristic
3. Seed company skill weights from CONTENT_SYSTEM.md profiles

**Risk:** None — no consumer yet

---

### Phase 2 — Skill Mastery (shadow mode)

1. Create `user_skill_scores`
2. On each submit, Progress consumer updates mastery **in addition to** existing flow
3. Log shadow readiness vs old formula for validation

**Risk:** Low — dual-write observation

---

### Phase 3 — Switch Readiness V2 (feature flag)

1. Implement `ReadinessEngine` in Progress service
2. Add `GET /progress/readiness` endpoints
3. Gateway dashboard calls new API behind `READINESS_V2=true`
4. Align submit `readiness_delta` to skill-based overall change

**Risk:** Medium — numbers change; communicate in UI

---

### Phase 4 — Persist and Trend

1. Write to `user_company_readiness` on each update
2. Enable readiness history graph (Analytics)
3. Remove Gateway `computeReadiness()`
4. Deprecate Question stats endpoint

**Risk:** Low

---

### Phase 5 — Explainability UI

1. Skill gap API
2. Dashboard "focus areas" widget
3. Companion messages reference weakest skill

**Risk:** Product/UX only

---

## 9. Backfill Strategy for Existing Users

For each user with `user_question_history`:

1. Join history → questions → question_skills (backfilled)
2. Compute initial mastery per skill using same formula as live updates
3. Compute initial company readiness from weights
4. Store in `user_skill_scores` + readiness cache

**Caution:** Historical data lacks skill mapping precision — mark backfilled scores with `source = 'backfill'` for potential recalibration.

---

## 10. Formula Recommendation (Design Only)

```
// Per answer, for each mapped skill:
mastery += (score / 100) * question.readiness_weight * difficulty_multiplier

// Rolling cap with decay (optional Phase B):
mastery = mastery * decay_factor + new_contribution

// Company readiness:
company_score = weighted_avg(skill_mastery for skills in company profile)

// Overall:
overall = avg(company_score for user_targets)
```

Constants in `config/readiness.go` (new file) — not hardcoded in handlers.

---

## 11. Checklist Answers

| Question | Answer |
|----------|--------|
| How is readiness currently calculated? | `(correctRate*100 + avgScore)/2` per company in Gateway; global avg score delta on submit |
| Which tables store readiness? | **None** — computed from `user_question_history` |
| Which services calculate readiness? | Gateway + Question (should be Progress only) |
| What inputs affect readiness? | Correctness, score, company tags, target companies |
| Architecture changes for Question→Skill→Company→Overall? | Skill graph, question_skills, user_skill_scores, company_skill_weights, Progress-owned engine, new events/APIs |

---

*This document is analysis only. No code was modified.*
