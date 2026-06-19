# A2 Content Architecture Summary

Phase A step 2 connects the journey map to the skill graph and question pools without changing daily paper selection, journey progression logic, or readiness scoring.

---

## New Schema

### `question_pools`

Curated question sets scoped to a single skill.

| Column        | Type    | Notes                                      |
|---------------|---------|--------------------------------------------|
| id            | UUID    | Primary key                                |
| skill_id      | UUID    | FK → `skills(id)`                          |
| slug          | TEXT    | Unique, stable identifier                  |
| name          | TEXT    | Display name                               |
| description   | TEXT    | Optional context for editors               |
| sort_order    | INT     | Ordering within a skill                    |
| created_at    | TIMESTAMPTZ | Auto-set                               |
| updated_at    | TIMESTAMPTZ | Auto-set via trigger                   |

### `pool_questions`

Many-to-many join between pools and approved questions.

| Column      | Type | Notes                              |
|-------------|------|------------------------------------|
| pool_id     | UUID | FK → `question_pools(id)`          |
| question_id | UUID | FK → `questions(id)`             |
| sort_order  | INT  | Selection order within the pool    |

Primary key: `(pool_id, question_id)`.

### `node_skills`

Skills taught or evaluated at a journey node.

| Column     | Type    | Notes                                      |
|------------|---------|--------------------------------------------|
| node_id    | UUID    | FK → `journey_nodes(id)`                   |
| skill_id   | UUID    | FK → `skills(id)`                          |
| is_primary | BOOLEAN | One primary skill per node in most worlds  |

Primary key: `(node_id, skill_id)`.

### `node_pools`

Question pools bound to a journey node with selection metadata.

| Column              | Type | Notes                                              |
|---------------------|------|----------------------------------------------------|
| node_id             | UUID | FK → `journey_nodes(id)`                           |
| pool_id             | UUID | FK → `question_pools(id)`                          |
| selection_strategy  | TEXT | `random_unseen`, `sequential`, or `boss_mixed`     |
| questions_required  | INT  | How many questions must be answered to clear node |

Primary key: `(node_id, pool_id)`.

### `journey_nodes.slug` (additive column)

Nullable slug added for stable content binding. Existing nodes without slugs continue to work; Foundation Forest nodes are backfilled.

---

## Relationships

```
worlds
  └── journey_nodes
        ├── node_skills ──→ skills
        │                      └── subskills
        │                      └── question_skills ← questions
        └── node_pools ──→ question_pools
                                └── pool_questions ──→ questions
```

**Content flow (target state):**

```
Journey Node → Skill(s) → Question Pool(s) → Question(s)
```

**Current runtime flow (unchanged):**

```
Journey Node → Daily Paper index → Question at position i
```

The new tables describe *what content belongs to a node*. Daily paper and journey status still use positional index mapping until Journey V2 is enabled.

---

## Migrations

| Migration | Purpose |
|-----------|---------|
| `000032_create_question_pools.up.sql` | Creates four tables; adds `journey_nodes.slug` |
| `000033_backfill_node_pools.up.sql` | Backfills Foundation Forest |

---

## Backfill Strategy

Foundation Forest (`foundation-forest`) is backfilled explicitly with stable UUIDs — no label-heuristic matching.

### Node slugs

| Node | Slug |
|------|------|
| Arrays Basics | `arrays-basics` |
| String Patterns | `string-patterns` |
| Hash Maps | `hash-maps` |
| Tree Traversal | `tree-traversal` |
| Forest Boss | `forest-boss` |

### Skill bindings

Each regular node maps to one primary skill. The boss node maps to three skills (arrays, hash-maps, trees) with `is_primary = false` on secondary skills.

### Pool bindings

Four primary pools cover the four regular nodes. Eight supplementary pools cover remaining seed-question skills for graph completeness. The boss node binds three pools with `boss_mixed` strategy.

All 12 approved seed questions are assigned to at least one pool via `pool_questions`.

### Future worlds

New worlds should follow the same pattern:

1. Assign slugs to `journey_nodes`
2. Create `question_pools` per skill segment
3. Populate `pool_questions` from approved questions tagged to those skills
4. Insert `node_skills` and `node_pools` rows

---

## Migration Strategy: Index → Skill → Pool → Question

### Phase 1 — Schema + backfill (this step, complete)

- Tables exist and Foundation Forest is wired
- Journey API returns additive `slug`, `skills`, and `pools` fields
- Daily paper selection unchanged
- Journey status still derived from daily paper index

### Phase 2 — Dual-read validation

- Compare index-mapped question vs pool-selected question per node
- Log mismatches in shadow mode
- No user-facing behavior change

### Phase 3 — Feature-flagged pool selection (Journey V2)

Enable with `JOURNEY_POOL_SELECTION=true` (`config.JourneyPoolSelectionEnabled()`).

When enabled:

1. For each `current` node, read `node_pools` bindings
2. Apply `selection_strategy`:
   - `random_unseen`: pick unseen approved question from pool
   - `sequential`: walk `pool_questions.sort_order`
   - `boss_mixed`: one question from each bound pool
3. Fall back to index mapping if pool is empty or flag is off

### Phase 4 — Deprecate index mapping

- Daily paper becomes a separate practice mode, not the journey driver
- Journey nodes own their question selection via pools
- Remove positional hack in `GetJourney`

---

## API Changes (backward compatible)

All changes are additive. Existing clients ignore new JSON fields.

### `GET /api/v1/journey`

Each node now includes optional fields:

```json
{
  "id": "...",
  "slug": "arrays-basics",
  "label": "Arrays Basics",
  "node_type": "lesson",
  "status": "current",
  "question_id": "...",
  "sort_order": 1,
  "skills": [
    { "skill_slug": "arrays", "skill_name": "Arrays", "is_primary": true }
  ],
  "pools": [
    {
      "pool_slug": "foundation-arrays-beginner",
      "pool_name": "Foundation Arrays",
      "skill_slug": "arrays",
      "selection_strategy": "random_unseen",
      "questions_required": 1,
      "question_count": 1
    }
  ]
}
```

`question_id` and `status` logic is unchanged.

### `GET /api/v1/journey/nodes/{id}/content` (new)

Returns full skill and pool bindings for a single node. Proxied through the gateway at the same path.

---

## Go Layer

| Package | File | Responsibility |
|---------|------|----------------|
| `constants` | `content.go` | Pool selection strategies, feature env vars, world slug |
| `config` | `features.go` | `JourneyPoolSelectionEnabled()` flag |
| `store` | `content.go` | Pool and node binding queries |
| `service` | `content.go` | Node content assembly, journey enrichment |
| `handler` | `content.go` | `GET /journey/nodes/{id}/content` |
| `dto` | `content.go` | Response types |

Journey progression logic in `service/journey.go` is untouched except for returning `slug` and additive enrichment via the handler layer.

---

## Risks

| Risk | Mitigation |
|------|------------|
| Index mapping and pool content diverge | Phase 2 shadow comparison before enabling pool selection |
| Empty pools block node progression | Fallback to index mapping when flag is off; empty-pool alerts in Phase 2 |
| Boss node multi-pool complexity | Explicit `boss_mixed` strategy; tested in backfill |
| New worlds added without backfill | Nodes work without bindings; slugs nullable; migration checklist for new worlds |
| Performance: N+1 queries on journey enrichment | Acceptable at current scale (5 nodes); batch query optimization in Journey V2 |
| Supplementary pools share questions across skills | Intentional for seed content; real pools will be skill-pure as bank grows |

---

## Future Journey V2 Plan

1. **Pool-driven selection** behind `JOURNEY_POOL_SELECTION` flag
2. **Unseen tracking** uses existing `user_question_history`, not pool membership alone
3. **Boss nodes** pull one question per bound pool using `boss_mixed`
4. **Daily paper decoupling** — daily paper becomes independent practice; journey uses pools
5. **Content ops tooling** — admin API to manage pools, assign questions, bind nodes (out of scope for A2)
6. **Multi-world rollout** — repeat Foundation Forest backfill pattern per world slug

### Explicitly out of scope (per Phase A plan)

- Recency decay
- Subskill mastery
- Readiness caching
- Additional company profiles
- Frontend changes
- Companion, quest, or world visual changes

---

## Verification

```bash
go build ./...
go test ./...
```

Store tests in `services/question/internal/store/content_test.go` verify Foundation Forest node bindings, pool question assignments, and skill-scoped pool listing after full migration run.
