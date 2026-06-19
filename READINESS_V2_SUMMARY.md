# Readiness V2 Foundation Summary

**Date:** 2026-06-09  
**Scope:** A4 Readiness Foundation — parallel to legacy readiness (no UI switch)  
**Reference:** `.ai/CONTENT_SYSTEM.md`, `.ai/ARCHITECTURE.md`, `.ai/EXECUTION.md` (A4)

---

## What Was Implemented

Skill-based readiness foundation in the **Progress service** (correct domain owner per `.ai/ARCHITECTURE.md`):

| Component | Status |
|-----------|--------|
| `user_skill_scores` table | ✅ |
| `company_skill_weights` table | ✅ |
| Seed weights (Google, Amazon, Meta, Uber) | ✅ |
| Live mastery updates on `question.answered` | ✅ |
| Historical backfill from answer history | ✅ |
| `GET /api/v1/skills/readiness` | ✅ |
| `GET /api/v1/companies/readiness` | ✅ |
| Legacy readiness (`/questions/stats/readiness`, dashboard) | ✅ Unchanged |

**Not implemented (deferred):** UI integration, gateway dashboard switch, `readiness.updated` events, subskill-level tracking, recency decay, Atlassian weights.

---

## Schema

### `user_skill_scores`

| Column | Type | Description |
|--------|------|-------------|
| `user_id` | UUID FK → users | User |
| `skill_id` | UUID FK → skills | Skill |
| `mastery` | INT 0–100 | Weighted skill mastery |
| `attempts` | INT | Answer count contributing to skill |
| `last_practiced_at` | TIMESTAMPTZ | Most recent answer timestamp |
| `source` | TEXT | `live` or `backfill` |
| `updated_at` | TIMESTAMPTZ | Auto-updated |

Primary key: `(user_id, skill_id)`

### `company_skill_weights`

| Column | Type | Description |
|--------|------|-------------|
| `company` | TEXT | Company slug |
| `skill_id` | UUID FK → skills | Weighted skill |
| `weight` | INT 0–100 | Relative importance |

Primary key: `(company, skill_id)`  
Constraint: weights per company sum to **100** (seed data).

---

## Formula

### Layer 1 — Answer → Skill Mastery (live updates)

For each `question_skills` mapping on an answered question:

```
contribution = (score / 100)
             × question.readiness_weight
             × question_skills.weight
             × difficulty_multiplier

new_mastery = min(100, round(current_mastery + contribution × 100 × MasterySmoothingFactor))
```

**Constants** (`config/readiness.go`):

| Constant | Value |
|----------|-------|
| `MasterySmoothingFactor` | 0.15 |
| `DifficultyMultiplier` easy / medium / hard | 0.90 / 1.00 / 1.10 |
| `MaxSkillMastery` | 100 |
| `MaxCompanyReadiness` | 95 |

### Layer 2 — Skill Mastery → Company Readiness

```
company_readiness = Σ(mastery[skill] × weight[company][skill]) / Σ(weight[company][skill])
```

Capped at `MaxCompanyReadiness` (95).

Skills without a user score contribute **0** mastery.

### Layer 3 — Company Readiness → Overall (API)

```
overall = average(company_readiness for user_targets)
```

Only companies in `user_targets` are included in `GET /companies/readiness`.

### Backfill formula (migration 000031)

For historical answers, per `(user_id, skill_id)`:

```
mastery = min(100, round(avg(contribution_per_answer) × 100))
```

Where `contribution_per_answer` uses the same factors as live updates.  
Source tagged `backfill`; live answers overwrite with `live` source.

---

## Company Weights (Seeded)

### Google (sum = 100)

| Skill | Weight |
|-------|--------|
| arrays | 15 |
| trees | 15 |
| graphs | 15 |
| dynamic-programming | 20 |
| system-design-scaling | 15 |
| behavioral-star | 10 |
| communication-clarity | 10 |

### Amazon (sum = 100)

| Skill | Weight |
|-------|--------|
| arrays | 15 |
| behavioral-leadership | 20 |
| behavioral-star | 15 |
| lld-patterns | 15 |
| problem-solving | 20 |
| system-design-data | 15 |

### Meta (sum = 100)

| Skill | Weight |
|-------|--------|
| arrays | 15 |
| dynamic-programming | 15 |
| system-design-scaling | 20 |
| behavioral-star | 15 |
| communication-structure | 15 |
| graphs | 20 |

### Uber (sum = 100)

| Skill | Weight |
|-------|--------|
| system-design-scaling | 25 |
| graphs | 20 |
| arrays | 15 |
| behavioral-star | 15 |
| problem-solving | 25 |

Stored in `company_skill_weights` — **not hardcoded in application logic**.

---

## Sample Calculations

### Example A — Live answer update

User answers **Two Sum** (question `...001`) with **score = 90**.

Mappings: arrays (0.6), hash-maps (0.4). Question: `readiness_weight = 0.80`, difficulty = easy.

```
arrays contribution = 0.90 × 0.80 × 0.60 × 0.90 = 0.3888
arrays delta        = 0.3888 × 100 × 0.15 ≈ 5.8 → +6 mastery (from 0 → 6)
```

Separate row updated for hash-maps with weight 0.4.

### Example B — Google company readiness

User target: **google**. Skill masteries:

| Skill | Mastery | Weight |
|-------|---------|--------|
| arrays | 85 | 15 |
| trees | 63 | 15 |
| graphs | 41 | 15 |
| dynamic-programming | 22 | 20 |
| system-design-scaling | 0 | 15 |
| behavioral-star | 0 | 10 |
| communication-clarity | 0 | 10 |

```
weighted_sum = 85×15 + 63×15 + 41×15 + 22×20 + 0 + 0 + 0
             = 1275 + 945 + 615 + 440 = 3275

google_readiness = 3275 / 100 = 32.75 → 33
```

If the same user later improves DP to 60 and behavioral-star to 50:

```
weighted_sum = 1275 + 945 + 615 + 1200 + 0 + 500 + 0 = 4535
google_readiness = 4535 / 100 = 45.35 → 45
```

### Example C — Parallel systems comparison

Same user with 4 answered Google-tagged questions, 75% correct rate, avg score 80:

| System | Google Score | Basis |
|--------|--------------|-------|
| **Legacy V1** | `(75 + 80) / 2 = 77` | Company tag hit rate + avg score |
| **V2 Foundation** | `33` (example B) | Weighted skill masteries |

Scores **will differ** until V2 has sufficient skill coverage and the user practices weighted skills. This is expected during parallel operation.

---

## API Reference

### `GET /api/v1/skills/readiness`

**Service:** Progress (8084)  
**Gateway:** `/api/v1/skills/readiness` → progress (routed before `/skills/*`)

```json
{
  "data": {
    "skills": [
      {
        "skill_slug": "arrays",
        "skill_name": "Arrays",
        "mastery": 6,
        "attempts": 1,
        "last_practiced_at": "2026-06-09T12:00:00Z"
      }
    ],
    "overall": 6,
    "version": "v2"
  }
}
```

### `GET /api/v1/companies/readiness`

```json
{
  "data": {
    "companies": [
      {
        "company": "google",
        "readiness": 33,
        "skill_contributions": [
          {"skill_slug": "arrays", "skill_name": "Arrays", "mastery": 85, "weight": 15}
        ]
      }
    ],
    "overall": 33,
    "version": "v2"
  }
}
```

**Not wired to UI.** Dashboard continues using legacy V1 via `GET /questions/stats/readiness`.

---

## Files Changed

### Migrations

| File | Purpose |
|------|---------|
| `000029_create_readiness_foundation.up.sql` | Tables |
| `000030_seed_company_skill_weights.up.sql` | Google, Amazon, Meta, Uber weights |
| `000031_backfill_user_skill_scores.up.sql` | Historical mastery from answer history |

### Progress Service

| File | Purpose |
|------|---------|
| `internal/store/readiness.go` | ReadinessStore repository |
| `internal/store/readiness_test.go` | Store tests |
| `internal/service/readiness.go` | ReadinessService + formula |
| `internal/service/readiness_test.go` | Unit tests |
| `internal/service/readiness_integration_test.go` | API + company sample tests |
| `internal/dto/readiness.go` | Response DTOs |
| `internal/handler/readiness.go` | HTTP handlers |
| `internal/handler/progress.go` | Inject readiness into handler |
| `internal/handler/internal.go` | Dev-sync updates mastery |
| `internal/consumer/consumer.go` | Kafka updates mastery |
| `cmd/main.go` | Wire routes |

### Shared

| File | Purpose |
|------|---------|
| `config/readiness.go` | Formula constants |
| `constants/readiness.go` | Source labels, company list |
| `services/gateway/cmd/main.go` | Proxy readiness routes |

---

## Parallel Operation

```
                    ┌─────────────────────────────────────┐
                    │         question.answered           │
                    └──────────────┬──────────────────────┘
                                   │
              ┌────────────────────┼────────────────────┐
              ▼                    ▼                    ▼
        Streak Service      Progress Service     (unchanged)
              │                    │
              │         ┌──────────┴──────────┐
              │         ▼                     ▼
              │    XP / Gems (V1)      user_skill_scores (V2)
              │                              │
              ▼                              ▼
        streak.updated              GET /skills/readiness
                                              │
                                              ▼
                                    GET /companies/readiness

Legacy path (unchanged):
  user_question_history → GET /questions/stats/readiness
                        → gateway computeReadiness()
                        → dashboard cards (V1)
```

Both systems run simultaneously. No feature flag required yet — V2 APIs are additive.

---

## Migration Path: V1 → V2

### Phase 1 — Foundation (this step) ✅

- [x] Schema + seed weights
- [x] Live mastery updates
- [x] Backfill from history
- [x] V2 read APIs
- [x] Parallel operation

### Phase 2 — Validation (next)

1. Add `READINESS_V2=true` env flag on gateway
2. Log V1 vs V2 scores for target users (shadow mode)
3. Compare deltas after each answer session
4. Tune `MasterySmoothingFactor` if mastery moves too fast/slow

### Phase 3 — Dashboard switch

1. Gateway dashboard calls `GET /companies/readiness` when flag enabled
2. Submit response adds `skill_deltas[]` from V2 (additive JSON field)
3. Deprecate `computeReadiness()` in gateway
4. Keep `GET /questions/stats/readiness` for 1 release cycle

### Phase 4 — Events & notifications

1. Emit `readiness.updated` from Progress on meaningful delta
2. Notification service consumes for transparency pushes ("Arrays +2 → Google 45%")
3. Companion dialogue references weakest skill from V2 gaps

### Phase 5 — Cleanup

1. Remove V1 formula from gateway
2. Remove `GET /questions/stats/readiness` from question service
3. Align submit `readiness_delta` to V2 overall change

### Re-backfill trigger

Re-run backfill logic if:
- Question skill mappings change materially
- Formula constants change
- New skills added to company profiles

Use `source = 'backfill'` filter to identify rows for recalculation without touching `live` rows.

---

## Risks

| Risk | Severity | Mitigation |
|------|----------|------------|
| V1 and V2 scores diverge, confusing internal testing | Medium | APIs tagged `version: "v2"`; document parallel operation |
| Users with no `user_targets` get empty company readiness | Low | Returns empty `companies[]`; same as V1 with no targets |
| Questions without `question_skills` skip mastery update | Medium | Step 1 backfilled all 12 seed questions; CMS validation in future |
| Backfill vs live formula slight mismatch | Low | Backfill uses avg contribution; live uses smoothing — converges with more answers |
| Mastery grows slowly due to 0.15 smoothing | Low | Tunable in `config/readiness.go` before UI switch |
| Gateway route `/skills/readiness` must precede `/skills/*` | Medium | Explicit route order in gateway — verified |

---

## Remaining Work

### Before content architecture (A2)

- [ ] Shadow comparison logging (V1 vs V2)
- [ ] Seed Atlassian weights (optional, onboarding list includes atlassian)
- [ ] `skill_prerequisites` table (optional, for gap UX)

### Before UI switch

- [ ] Gateway feature flag `READINESS_V2`
- [ ] Dashboard readiness cards from V2
- [ ] Submit response skill deltas
- [ ] Skill gap widget ("Weakest: Dynamic Programming")

### Before full A4 completion

- [ ] `readiness.updated` Kafka event
- [ ] Persisted `user_company_readiness` cache (optional performance)
- [ ] Subskill-level mastery (`user_subskill_scores`) for finer explainability
- [ ] Recency decay on mastery

### Explicitly not in scope (later phases)

- Journey / pool / quest / companion changes
- UI changes
- Removing legacy readiness

---

## Verification

```bash
go build ./...
go vet ./...
go test ./...
```

Apply migrations `000029` → `000031` on existing databases with Step 1 schema present.

Test V2 APIs (authenticated):

```bash
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/skills/readiness
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/companies/readiness
```

---

*Readiness foundation only. Legacy system preserved. UI unchanged.*
