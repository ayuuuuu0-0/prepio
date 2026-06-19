# Journey V2 Summary

Journey V2 replaces index-based node→question mapping with pool-driven selection behind the `JOURNEY_POOL_SELECTION` feature flag. Daily paper generation, readiness, companions, quests, and UI are unchanged.

---

## Architecture

```
GET /api/v1/journey
        │
        ▼
  JOURNEY_POOL_SELECTION?
        │
   ┌────┴────┐
   OFF       ON
   │         │
   ▼         ▼
 Index     Pool
 Journey   Journey
 (V1)      (V2)
```

**V1 (flag off):** Node `i` maps to `daily_paper.questions[i]`. Boss unlocks when all daily paper questions are answered in the session. Unchanged from pre-V2 behavior.

**V2 (flag on):** Each node reads `node_pools` bindings, selects questions from `pool_questions` using the configured strategy, and tracks progress against session answers.

### Components

| File | Role |
|------|------|
| `service/journey.go` | Routes to V1 or V2 builder |
| `service/journey_selection.go` | Pool strategies, fallback, boss mixed |
| `store/content.go` | `node_pools`, `pool_questions` queries |
| `store/history.go` | `user_question_history` for unseen tracking |
| `config/features.go` | `JourneyPoolSelectionEnabled()` |

### Data flow (V2)

```
journey_node
    └── node_pools (selection_strategy, questions_required)
            └── question_pools
                    └── pool_questions → questions
```

Unseen state comes from `user_question_history` (all-time, not session-scoped).

---

## Selection Flow

### Per-node resolution

1. Load `node_pools` for the node
2. If no pools → **index fallback**
3. If `boss_mixed` or multiple pools → select one question per pool
4. Otherwise use the first pool's strategy on its question list
5. If pool is empty or selection fails → **index fallback**

### Strategies

| Strategy | Behavior |
|----------|----------|
| `random_unseen` | Pick randomly from pool questions not in `user_question_history`. If all seen, pick randomly from full pool. Seed = hash(userID + sessionID + nodeID) for stability within a session. |
| `sequential` | First pool question (by `sort_order`) not in history. If all seen, return first in order. |
| `boss_mixed` | One question per bound pool, each using `random_unseen` unless the pool row specifies another strategy. |

### Status computation (V2)

| Status | Condition |
|--------|-----------|
| `done` | All assigned question IDs answered in current session |
| `current` | First unlocked node with unanswered assigned questions |
| `locked` | Prior nodes not done, or no questions assigned |

**Display `question_id`:** First unanswered assigned question; if all answered, last assigned question.

### Boss node (V2)

Boss node binds multiple pools (Foundation Forest: arrays, hash-maps, trees). Requires all three pool-selected questions to be answered in the session before status becomes `done`. Unlocks when prior nodes are done (sequential progression).

---

## Fallback Behavior

Pool selection is skipped and index mapping is used when:

| Condition | Result |
|-----------|--------|
| `JOURNEY_POOL_SELECTION` not `true` | Full V1 index journey |
| Node has no `node_pools` rows | `paper.questions[i]` |
| Pool has zero approved questions | Index fallback |
| Selection returns empty | Index fallback |
| `ListNodePools` / `ListPoolQuestionIDs` errors | Index fallback |

Index fallback preserves backward compatibility for worlds not yet backfilled.

---

## Migration Strategy

### Current state

- Foundation Forest is backfilled (migration `000033`)
- Flag defaults **off** in all environments
- V1 and V2 run side-by-side; switch via env var

### Rollout steps

1. **Deploy** with flag off — zero behavior change
2. **Validate** in staging with `JOURNEY_POOL_SELECTION=true`:
   - Confirm node `question_id` values match pool content
   - Confirm status transitions on session answers
   - Confirm boss node requires all pool questions
3. **Enable** per environment when content ops confirms pool coverage
4. **Future:** Decouple daily paper from journey (Phase 4 in `CONTENT_ARCHITECTURE_SUMMARY.md`)

### No new migrations required

Journey V2 uses tables created in A2 (`000032`, `000033`).

---

## Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| Pool question not in daily paper session | User sees journey question but submit rejects (session check) | Pool backfill uses seed questions; document overlap requirement until daily paper decoupling |
| Divergence from V1 progression | Boss unlock rules differ between V1 (all paper done) and V2 (prior nodes done) | Flag off preserves V1 exactly; document V2 boss behavior |
| All pool questions seen | `random_unseen` re-serves seen questions | Acceptable for small seed bank; expand pools as content grows |
| Non-deterministic random across restarts | Same user may get different question if session changes | Seed includes sessionID — stable within a daily session |
| Worlds without backfill | Empty pools trigger index fallback | Safe degradation; no broken journeys |

---

## Enable

```bash
export JOURNEY_POOL_SELECTION=true
```

Restart the question service (and gateway if proxying).

---

## Verification

```bash
go build ./...
go test ./services/question/internal/service/...
```

Tests:

- `journey_selection_test.go` — unit tests for all three strategies
- `journey_test.go` — integration test with Foundation Forest backfill

---

## Out of Scope

- Daily paper generation changes
- Readiness, companions, quests, UI
- Recency decay, subskill mastery, readiness caching
- Additional company profiles or worlds beyond existing backfill
