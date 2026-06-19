# Readiness Validation Report — Phase A4

**Date:** 2026-06-10  
**Scope:** Parallel V1/V2 validation before content architecture work  
**Status:** Validation infrastructure ready — production UI unchanged

---

## Purpose

Readiness is Prepio's moat metric. Before building question pools, journey refactoring, or companion systems, we must validate that **Readiness V2** produces useful, explainable progression compared to legacy V1.

This document defines both formulas, example outputs, expected differences, and tuning guidance for the validation period.

---

## Feature Flag

| Variable | Default | Effect |
|----------|---------|--------|
| `READINESS_V2` | `false` | When `true`, adds `readiness_v2` array to `GET /dashboard/home` alongside unchanged `readiness` (V1) |

**Production UI is not switched.** V1 cards remain the primary display. V2 appears only as an additive field when the flag is enabled.

Validation APIs work **regardless of flag state** — they always return both V1 and V2 for comparison.

```bash
# Enable V2 visibility on dashboard home (optional, for testers)
export READINESS_V2=true
```

---

## V1 Formula (Legacy)

**Source:** Question service answer history + company tags  
**Used by:** `GET /dashboard/home` → `readiness[]`, submit `readiness_delta` (partial)

```
correct_rate = (correct_answers / total_answers) × 100
company_readiness = (correct_rate + avg_score) / 2
cap at 95
overall = average across user_targets
```

**Characteristics:**
- Based on **company tag performance**, not skills
- Ignores question difficulty, readiness_weight, skill mapping
- Rewards answering many tagged questions correctly
- Not explainable by skill gap
- Can inflate when user answers easy tagged questions well

**Implementation:** `shared/readiness/v1.go` → `ComputeV1CompanyScore()`

---

## V2 Formula (Foundation)

**Source:** `user_skill_scores` × `company_skill_weights`  
**Used by:** Progress service readiness APIs

### Skill mastery (per answer)

```
contribution = (score/100) × readiness_weight × skill_weight × difficulty_multiplier
new_mastery  = min(100, current + contribution × 100 × 0.15)
```

### Company readiness

```
company_readiness = Σ(mastery[skill] × weight[company][skill]) / Σ(weight)
cap at 95
```

### Skill gap score

```
gap_score = weight × (100 - mastery) / 100
```

Skills with mastery ≥ 70 are not flagged as gaps.

**Characteristics:**
- Based on **weighted skill mastery**
- Explainable per company and per skill
- Penalizes missing practice on high-weight skills (shows as 0 mastery)
- Grows slowly due to smoothing (0.15 factor)
- Requires `question_skills` mapping to update

**Implementation:** `services/progress/internal/service/readiness.go`

---

## Validation APIs

| Endpoint | Service | Purpose |
|----------|---------|---------|
| `GET /api/v1/dashboard/readiness` | Gateway | Full V1/V2 comparison + gaps + explanations |
| `GET /api/v1/internal/readiness/compare` | Gateway | Same payload for internal tooling |
| `GET /api/v1/readiness/dashboard` | Progress | V2 skill + company dashboard |
| `GET /api/v1/skills/readiness` | Progress | Skill mastery + top/weakest/gaps |
| `GET /api/v1/companies/readiness` | Progress | Company readiness + breakdown |

All require authentication. **Not wired to production UI components.**

### Response highlights

Every V2 readiness response includes:

- **`top_skills`** — highest mastery among practiced skills
- **`weakest_skills`** — lowest mastery among practiced skills
- **`skill_gaps`** — weighted skills below 70 mastery blocking company readiness
- **`explanations`** — human-readable summaries

Validation response adds:

- **`v1`** — legacy snapshot with formula
- **`v2`** — full dashboard
- **`comparison`** — per-company and overall deltas

---

## Example Users

### User A — "Tag Grinder" (V1 inflated)

**Profile:** Target Google. Answered 8 easy Google-tagged questions with 88% avg score.

| Metric | V1 | V2 |
|--------|----|----|
| Google readiness | **82** | **28** |
| Overall | 82 | 28 |

**Why they diverge:**
- V1 sees high correct rate + score on Google tags
- V2 sees only 2–3 skills practiced (arrays, hash-maps) — trees, graphs, DP, system design all **0**
- Google weights DP (20%) and graphs (15%) heavily

**V2 skill gaps (Google):**
1. Dynamic Programming — 0 mastery, weight 20, gap_score 20
2. Graphs — 0 mastery, weight 15, gap_score 15
3. System Design Scaling — 0 mastery, weight 15, gap_score 15

**Validation insight:** V2 correctly identifies this user is **not** Google-ready despite good tag performance.

---

### User B — "Balanced Prepper" (V1 and V2 aligned)

**Profile:** Target Google + Amazon. Mixed practice across arrays, trees, behavioral.

| Skill | Mastery |
|-------|---------|
| arrays | 78 |
| trees | 65 |
| behavioral-star | 72 |
| dynamic-programming | 45 |
| graphs | 38 |

| Metric | V1 Google | V2 Google |
|--------|-----------|-----------|
| Readiness | **58** | **52** |

**Delta:** -6 (V2 slightly lower due to zero-weight unpracticed skills counting in denominator)

**Validation insight:** Scores converge when practice matches company weight profile. V2 adds explainability V1 lacks.

---

### User C — "New User" (both zero)

**Profile:** Just onboarded, no answers.

| Metric | V1 | V2 |
|--------|----|----|
| Google | 0 | 0 |
| Overall | 0 | 0 |

**Validation insight:** Both agree on cold start. No false confidence.

---

### User D — "Single Boss Question" (post backfill)

**Profile:** Answered food delivery system design (hard) once with score 75.

| Skill | Mastery (backfill) |
|-------|-------------------|
| system-design-scaling | ~12 |
| system-design-fundamentals | ~12 |

| Metric | V1 (zepto tag) | V2 Google |
|--------|----------------|-----------|
| Readiness | 75 | 4 |

**Validation insight:** V2 reflects narrow skill exposure. V1 may show high if tagged company matches.

---

## Example API Output (abbreviated)

### `GET /api/v1/dashboard/readiness`

```json
{
  "data": {
    "readiness_v2_enabled": false,
    "v1": {
      "version": "v1",
      "formula": "(correct_rate * 100 + avg_score) / 2, capped at 95",
      "overall": 77,
      "companies": [{"company": "google", "score": 77}]
    },
    "v2": {
      "skill_mastery": {
        "overall": 41,
        "top_skills": [{"skill_slug": "arrays", "mastery": 85, "attempts": 3}],
        "weakest_skills": [{"skill_slug": "graphs", "mastery": 0, "attempts": 0}],
        "skill_gaps": [
          {
            "company": "google",
            "skill_slug": "dynamic-programming",
            "mastery": 0,
            "weight": 20,
            "gap_score": 20,
            "explanation": "Dynamic Programming mastery is 0 — carries 20% weight for google readiness"
          }
        ],
        "explanations": [
          {
            "scope": "skills",
            "summary": "Overall skill mastery average is 41 across practiced skills"
          }
        ]
      },
      "company_readiness": {
        "overall": 33,
        "companies": [
          {
            "company": "google",
            "readiness": 33,
            "explanation": {
              "scope": "google",
              "summary": "google readiness is 33 — weakest weighted skill is Dynamic Programming (0 mastery)"
            }
          }
        ]
      }
    },
    "comparison": {
      "overall_v1": 77,
      "overall_v2": 33,
      "overall_delta": -44,
      "v1_formula": "(correct_rate * 100 + avg_score) / 2, capped at 95",
      "v2_formula": "sum(skill_mastery * company_weight) / sum(company_weight), capped at 95",
      "by_company": [
        {"company": "google", "v1_score": 77, "v2_score": 33, "delta": -44}
      ]
    }
  }
}
```

---

## Expected Differences

| Scenario | V1 vs V2 | Expected |
|----------|----------|----------|
| Few answers, high scores | V1 > V2 | Yes — V2 penalizes unpracticed weighted skills |
| Broad skill practice | V1 ≈ V2 | Yes — convergence |
| Wrong company tags on questions | V1 misleading | V2 unaffected by tags |
| Hard questions | V1 ignores difficulty | V2 applies difficulty multiplier |
| New user | Both 0 | Yes |
| Many easy streak answers | V1 climbs fast | V2 climbs slowly (smoothing) |

**Large negative delta (V2 << V1) is often correct behavior** — it means the user looks ready by tag hit rate but has skill gaps.

---

## Tuning Recommendations

### During validation (now)

1. **Log comparison deltas** for every user session in staging
2. Review `skill_gaps` — do they match intuitive weak areas?
3. Check explanations — are they actionable?
4. Track whether `top_skills` / `weakest_skills` match user self-assessment

### If V2 moves too slowly

| Parameter | Current | Try |
|-----------|---------|-----|
| `MasterySmoothingFactor` | 0.15 | 0.20–0.25 |
| `readinessGapMasteryThreshold` | 70 | 65 |

Location: `config/readiness.go`, `readiness_analysis.go`

### If V2 punishes new users too harshly

- Option A: Default unpracticed weighted skills to `null` (exclude from denominator) instead of 0
- Option B: Show "insufficient data" below N attempts per company profile
- **Do not change before validation period completes**

### If gap list is too long

- Current: all skills below 70 mastery with weight > 0
- Consider: top 3 gaps by `gap_score` only in UI (API already sorts)

### Before UI switch

1. Enable `READINESS_V2=true` for internal testers
2. Compare dashboard `readiness` vs `readiness_v2` side by side
3. Require ≥2 weeks of session data
4. Sign-off criteria:
   - [ ] Explanations rated useful by team
   - [ ] Gap skills match content roadmap priorities
   - [ ] No user confusion in moderated testing
   - [ ] Delta direction makes sense for 80%+ test accounts

---

## Migration Path (unchanged from foundation)

| Phase | Action |
|-------|--------|
| **Now** | Validation APIs + parallel operation |
| Next | Shadow logging, `READINESS_V2=true` for testers |
| Then | Dashboard UI reads V2 when flag on |
| Then | Submit response includes skill deltas |
| Finally | Remove V1 formula from gateway |

---

## Files Added/Changed

| File | Change |
|------|--------|
| `config/features.go` | `READINESS_V2` flag |
| `shared/readiness/v1.go` | Extracted V1 formula |
| `shared/readiness/types.go` | Shared readiness DTOs |
| `services/progress/internal/service/readiness_analysis.go` | Top/weakest/gaps |
| `services/progress/internal/service/readiness.go` | Enhanced responses |
| `services/gateway/internal/dashboard/readiness_validation.go` | Comparison aggregator |
| `services/gateway/internal/dashboard/service.go` | Optional V2 on home |
| Gateway routes | `/dashboard/readiness`, `/internal/readiness/compare` |

---

## Risks

| Risk | Mitigation |
|------|------------|
| Team misreads low V2 as bug | Document expected V1 > V2 for tag grinders |
| Validation APIs used in prod UI accidentally | Not wired to web/mobile components |
| Gateway → Progress dependency for validation | Progress must be running; errors return 502 |
| Cross-service DTO drift | Shared types in `shared/readiness/types.go` |

---

## Remaining Work (after validation sign-off)

- [ ] Wire web/mobile validation panel (internal only)
- [ ] Shadow log V1/V2 delta on every answer submit
- [ ] Dashboard UI switch behind `READINESS_V2`
- [ ] Submit `skill_deltas[]` on result card
- [ ] `readiness.updated` Kafka event
- [ ] Atlassian company weights
- [ ] Question pools + journey (deferred until validation complete)

---

## How to Validate

```bash
# 1. Apply all migrations (through 000031)
# 2. Start stack
# 3. Login and answer questions
# 4. Compare:

curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/dashboard/readiness | jq .

curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/internal/readiness/compare | jq .

# 5. Optional: enable flag and check home adds readiness_v2
READINESS_V2=true # restart gateway
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/dashboard/home | jq '.data.readiness, .data.readiness_v2'
```

---

*Validation infrastructure complete. Do not proceed to question pools or journey until readiness sign-off.*
