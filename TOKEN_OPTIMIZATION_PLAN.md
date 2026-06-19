# Token Optimization Plan

Analysis of Prepio documentation to reduce AI context consumption while preserving architectural decision quality.

**Scope:** Documentation only. No code or schema changes.

---

## Executive Summary

The repository carries **~191 KB of markdown** (~**48,000 tokens**) across **18 planning/audit/summary files**, plus **~3,000 tokens** of `AGENTS.md` rules injected on every Cursor conversation via workspace rules.

The dominant waste is:

1. **Reading the full `.ai/` canon (4 files, ~6,500 tokens) on every task** when only boundaries and current phase matter.
2. **Re-reading completed phase summaries (~12,000 tokens)** that duplicate `.ai/` specs and implemented migrations.
3. **Audit documents (~24,300 tokens)** describing pre-Phase-A state that no longer exists.
4. **Triple-stacked non-negotiables** across Cursor rules, `.ai/ARCHITECTURE.MD`, and missing/stale root copies.
5. **Stale root `PRODUCT.md` / `EXECUTION.md`** that conflict with `.ai/` versions.

**Estimated savings after restructure: 80–93%** on documentation tokens per task.

---

## Current Documentation Structure

### Tier 1 — Canonical `.ai/` (intended source of truth)

| File | Lines | ~Tokens | Role |
|------|------:|--------:|------|
| `.ai/PRODUCT.MD` | 653 | 1,750 | Product philosophy, pillars, UX anti-patterns |
| `.ai/ARCHITECTURE.MD` | 722 | 1,750 | Service ownership, events, API/DB rules |
| `.ai/CONTENT_SYSTEM.MD` | 708 | 1,500 | Learning hierarchy, skills, readiness model |
| `.ai/EXECUTION.MD` | 630 | 1,500 | Phase tasks, priorities (**status stale**) |
| **Subtotal** | **2,713** | **~6,500** | |

Each file repeats an **AGENT INSTRUCTIONS** preamble (read all 4 docs). ~80 lines × 4 = **~320 lines of duplicate preamble**.

### Tier 2 — Phase completion summaries (root)

| File | Lines | ~Tokens | Status |
|------|------:|--------:|--------|
| `PHASE_A_STEP1_SUMMARY.md` | 220 | 2,000 | A1 + A3 **complete** |
| `CONTENT_ARCHITECTURE_SUMMARY.md` | 271 | 2,300 | A2 **complete** |
| `READINESS_V2_SUMMARY.md` | 443 | 3,200 | A4 foundation **complete** |
| `READINESS_VALIDATION.md` | 396 | 3,100 | A4 validation **complete** |
| `JOURNEY_V2_SUMMARY.md` | 166 | 1,400 | Journey V2 **complete** |
| **Subtotal** | **1,496** | **~12,000** | |

### Tier 3 — Pre-implementation audits (`audit result/`)

| File | Lines | ~Tokens | Status |
|------|------:|--------:|--------|
| `CONTENT_AUDIT.md` | 387 | 4,300 | **Obsolete** — describes schema before 000025 |
| `DOMAIN_AUDIT.md` | 339 | 3,700 | **Obsolete** — pre skill/pool/journey V2 |
| `JOURNEY_AUDIT.md` | 250 | 2,200 | **Obsolete** — index-mapping since replaced |
| `READINESS_AUDIT.md` | 311 | 2,500 | **Obsolete** — V1-only readiness |
| `SKILL_GRAPH_PROPOSAL.md` | 492 | 3,900 | **Implemented** — migrations 000025–000028 |
| `CONTENT_MODEL_PROPOSAL.md` | 466 | 3,700 | **Implemented** — migrations 000032–000033 |
| `PHASE_A_IMPLEMENTATION_PLAN.md` | 499 | 3,900 | **Mostly done** — steps 1–11 complete |
| **Subtotal** | **2,744** | **~24,300** | |

Also: `audit-result.zip` (duplicate of above).

### Tier 4 — Stale root copies (conflicting)

| File | Lines | ~Tokens | Problem |
|------|------:|--------:|---------|
| `PRODUCT.md` | 333 | ~2,500 | Different content from `.ai/PRODUCT.MD` |
| `EXECUTION.md` | 341 | ~2,500 | Different content from `.ai/EXECUTION.MD` |

### Tier 5 — Always-on context (Cursor workspace rules)

| Source | ~Tokens | Problem |
|--------|--------:|---------|
| `AGENTS.md` content in rules | ~3,000 | Injected every conversation; overlaps `.ai/ARCHITECTURE.MD` |
| Note: `AGENTS.md` file missing from repo | — | Rules reference a file that does not exist on disk |

### Tier 6 — Not in scope (do not archive)

| File | Reason |
|------|--------|
| `web/README.md`, `mobile/README.md` | Platform setup |
| `mobile/ios/.../README.md` | Asset placeholder |

---

## Duplicate Information Map

| Topic | Duplicated across | Keep one source |
|-------|-------------------|-----------------|
| **Non-negotiables** (streak ownership, no hardcoded rewards, Kafka notifications) | Cursor rules, `.ai/ARCHITECTURE.MD` §Anti Patterns, audit `DOMAIN_AUDIT` | `.ai/RULES.md` (new, condensed) |
| **Learning hierarchy** World→Node→Skill→Pool→Question | `.ai/CONTENT_SYSTEM.MD`, `CONTENT_MODEL_PROPOSAL`, `CONTENT_ARCHITECTURE_SUMMARY`, `.ai/ARCHITECTURE.MD` §Domain Model | `.ai/CONTENT.md` (trimmed) |
| **Service ownership** | `.ai/ARCHITECTURE.MD`, `DOMAIN_AUDIT`, Cursor rules | `.ai/ARCHITECTURE.md` (trimmed boundaries only) |
| **Phase A task definitions** | `.ai/EXECUTION.MD` §A1–A5, `PHASE_A_IMPLEMENTATION_PLAN`, all `*_SUMMARY.md` headers | `.ai/STATE.md` (status) + `.ai/ROADMAP.md` (future phases) |
| **Skill graph taxonomy** | `.ai/CONTENT_SYSTEM.MD`, `SKILL_GRAPH_PROPOSAL`, `PHASE_A_STEP1_SUMMARY`, migration seeds | DB seeds + `.ai/CONTENT.md` §Skills (reference slugs only) |
| **Readiness V2 formula** | `.ai/CONTENT_SYSTEM.MD` §Readiness, `READINESS_V2_SUMMARY`, `READINESS_VALIDATION`, `READINESS_AUDIT` | `docs/reference/READINESS.md` (on demand) |
| **Journey pool selection** | `CONTENT_ARCHITECTURE_SUMMARY` §Phase 3–4, `JOURNEY_V2_SUMMARY`, `JOURNEY_AUDIT` | `docs/reference/JOURNEY.md` (on demand) |
| **Schema for pools/nodes** | `CONTENT_MODEL_PROPOSAL` §5, `CONTENT_ARCHITECTURE_SUMMARY` §New Schema, migrations | Migrations (code truth) + `docs/reference/CONTENT_ARCHITECTURE.md` |
| **Agent read instructions** | Preamble in all 4 `.ai/` files | Single line in `.ai/STATE.md` |
| **"What is Prepio"** narrative | `.ai/PRODUCT.MD` (654 lines), root `PRODUCT.md`, `.ai/EXECUTION.MD` intro | `.ai/PRODUCT.md` (100-line trim, product tasks only) |
| **Pre-implementation gaps** | All 7 audit files | **Archive** — misleading for current codebase |

---

## Documents to Merge

| Merge into | Source files | Rationale |
|------------|--------------|-----------|
| `docs/reference/READINESS.md` | `READINESS_V2_SUMMARY.md` + `READINESS_VALIDATION.md` | Same domain; validation adds examples only |
| `docs/reference/JOURNEY.md` | `JOURNEY_V2_SUMMARY.md` + `CONTENT_ARCHITECTURE_SUMMARY.md` §Migration Phase 3–4 | Journey V2 is the operational doc; A2 summary's future plan is now past |
| `docs/reference/CONTENT_ARCHITECTURE.md` | `CONTENT_ARCHITECTURE_SUMMARY.md` §Schema + §Backfill | Keep schema/backfill UUIDs; drop completed migration narrative |
| `docs/archive/completed/A1_A3.md` | `PHASE_A_STEP1_SUMMARY.md` | Historical record only |
| `.ai/RULES.md` (new) | Cursor `AGENTS.md` rules + `.ai/ARCHITECTURE.MD` §Anti Patterns + §Architectural Rule | Single non-negotiables file (~120 lines) |
| `.ai/STATE.md` (new) | `.ai/EXECUTION.MD` §Current Phase + completion summaries' status sections | Only file that changes weekly |

---

## Documents to Archive

Move to `docs/archive/` — do not delete (git history + audit trail).

### `docs/archive/audits/` (never read for implementation)

- `CONTENT_AUDIT.md`
- `DOMAIN_AUDIT.md`
- `JOURNEY_AUDIT.md`
- `READINESS_AUDIT.md`
- `SKILL_GRAPH_PROPOSAL.md`
- `CONTENT_MODEL_PROPOSAL.md`
- `PHASE_A_IMPLEMENTATION_PLAN.md`
- `audit-result.zip`

### `docs/archive/completed/` (read only for migration UUID lookup)

- `PHASE_A_STEP1_SUMMARY.md`
- `CONTENT_ARCHITECTURE_SUMMARY.md`
- `READINESS_V2_SUMMARY.md`
- `READINESS_VALIDATION.md`
- `JOURNEY_V2_SUMMARY.md`

### `docs/archive/legacy/` (conflicting stale copies)

- `PRODUCT.md` (root)
- `EXECUTION.md` (root)

---

## Proposed Documentation Structure

```
.ai/
├── STATE.md              ~80 lines   ← READ EVERY TASK (only volatile doc)
├── RULES.md              ~120 lines  ← READ EVERY TASK (non-negotiables)
├── ARCHITECTURE.md       ~200 lines  ← backend/service tasks
├── CONTENT.md            ~180 lines  ← content/journey/readiness tasks
├── PRODUCT.md            ~100 lines  ← product/UI/UX tasks only
└── ROADMAP.md            ~80 lines   ← Phase B/C/D (locked, rarely read)

docs/reference/           ← on-demand deep dives
├── READINESS.md          ~150 lines  (merged A4 docs)
├── JOURNEY.md            ~120 lines  (pool selection, flags)
└── CONTENT_ARCHITECTURE.md ~100 lines (schema, UUIDs, backfill)

docs/archive/
├── audits/               (7 files + zip)
├── completed/            (5 phase summaries)
└── legacy/               (root PRODUCT.md, EXECUTION.md)

migrations/               ← schema source of truth (agents read specific files, not all)
config/                   ← reward/level constants source of truth
```

**Total active AI corpus:** ~760 lines ≈ **~9,500 tokens** (down from ~48,000).

---

## STATE.md Design

Lightweight file. Updated when phase status changes. **No philosophy, no schema, no audit history.**

```markdown
# STATE

Last updated: 2026-06-10

## Current Phase
Phase A — Foundation Rebuild (A5 remaining)

## Completed
- [x] A1 Skill Graph (migrations 000025–000026, APIs live)
- [x] A3 Question Schema Upgrade (000027–000028)
- [x] A4 Readiness V2 Foundation (000029–000031, `READINESS_V2` flag)
- [x] A4 Validation (approved — no further readiness work unless requested)
- [x] A2 Content Architecture (000032–000033, node/pool bindings)
- [x] Journey V2 pool selection (`JOURNEY_POOL_SELECTION` flag, default OFF)

## Current Priorities
1. A5 Content Management Foundation (P1)
2. Phase A definition-of-done checklist
3. Enable Journey V2 in staging after pool/session overlap validation

## Blocked
- Phase B (Content Expansion) — locked until A5 complete
- Readiness V2 UI switch — awaiting product decision; backend ready
- Daily paper decoupling from journey — blocked on session overlap (see docs/reference/JOURNEY.md)

## Next Tasks (A5)
- Internal APIs: create/list skills, questions, pools, nodes
- Admin auth gate (internal only)
- No migration-based content additions for new questions

## Feature Flags
| Flag | Default | Effect |
|------|---------|--------|
| `READINESS_V2` | OFF | Adds `readiness_v2` on dashboard home |
| `JOURNEY_POOL_SELECTION` | OFF | Pool-driven journey question selection |

## Out of Scope (do not implement unless explicitly requested)
- Recency decay, subskill mastery, readiness caching
- Additional company profiles beyond seeded four
- Frontend/UI/companion/quest/world visual changes
- Phase B content expansion (100+ questions)
- Phase C (Arena), Phase D (AI Interviewers)
- Re-reading audit documents

## Task Reading Guide
| Task type | Read |
|-----------|------|
| Any | `.ai/STATE.md`, `.ai/RULES.md` |
| Backend/service | + `.ai/ARCHITECTURE.md` |
| Content/journey/pools/readiness | + `.ai/CONTENT.md` + relevant `docs/reference/` |
| Product/UX | + `.ai/PRODUCT.md` |
| Future phase planning | + `.ai/ROADMAP.md` |

## Key Migration Range
000025–000033 (Phase A foundation). Do not re-read completed summaries; read migration files directly if schema detail needed.
```

---

## Agent Reading Recommendations

### Read on every task (~1,250 tokens)

| File | Why |
|------|-----|
| `.ai/STATE.md` | Current phase, priorities, blockers, out-of-scope |
| `.ai/RULES.md` | Non-negotiables that prevent architectural violations |

Replace the current Cursor workspace rule blob with: *"Read `.ai/STATE.md` and `.ai/RULES.md` before any change."*

### Read for specific task types only

| Task type | Additional reads | ~Tokens |
|-----------|------------------|--------:|
| Go backend / API / Kafka | `.ai/ARCHITECTURE.md` | +1,000 |
| Skills, pools, journey, questions | `.ai/CONTENT.md` | +900 |
| Journey V2 / pool selection bug | `docs/reference/JOURNEY.md` | +600 |
| Readiness formula / V2 bug | `docs/reference/READINESS.md` | +750 |
| Schema / backfill UUIDs | `docs/reference/CONTENT_ARCHITECTURE.md` or specific `migrations/0000XX_*.up.sql` | +400–800 |
| Product / UX / gamification | `.ai/PRODUCT.md` | +500 |
| Phase B+ planning | `.ai/ROADMAP.md` | +400 |
| Flutter / Next.js | Platform README + `.ai/RULES.md` Flutter/frontend sections only | +300 |

### Never read again (unless explicitly investigating history)

| Category | Files | Why |
|----------|-------|-----|
| Pre-implementation audits | `docs/archive/audits/*` | Describe codebase that no longer exists; causes wrong decisions |
| Completed phase summaries | `docs/archive/completed/*` | Duplicates migrations + reference docs |
| Stale root copies | `docs/archive/legacy/*` | Conflicts with `.ai/` |
| Full `.ai/` canon (old 4-file read) | Old 2,713-line read pattern | Replaced by STATE + RULES + targeted slice |
| `PHASE_A_IMPLEMENTATION_PLAN.md` | Steps 1–11 done | Misleading remaining-work sections |
| Proposal docs | `SKILL_GRAPH_PROPOSAL`, `CONTENT_MODEL_PROPOSAL` | Implemented in migrations |
| Entire migration folder | All 33 migrations | Read only the 1–2 relevant files cited in STATE |

---

## Recommended Prompting Pattern

### User prompt template (copy-paste)

```
Read .ai/STATE.md and .ai/RULES.md only.

Task: [description]

If backend: also read .ai/ARCHITECTURE.md.
If content/journey/readiness: also read .ai/CONTENT.md and docs/reference/[RELEVANT].md.
Do not read docs/archive/ or completed summaries.
Do not read all four legacy .ai/ canon files.
```

### Cursor rule replacement (suggested)

Current rule injects ~3,000 tokens of full `AGENTS.md` on every message.

Replace with:

```
Before any code change, read .ai/STATE.md and .ai/RULES.md.
Non-negotiables in RULES.md override all other docs.
For schema truth, read migrations/ directly — not audit docs.
Never read docs/archive/ unless user explicitly asks for historical analysis.
```

Keep detailed Go/Flutter rules in `.ai/RULES.md` on disk (loaded on demand), not in always-on workspace rules.

### Agent behavior rules

1. **STATE.md is the only doc that changes phase status** — do not infer completion from EXECUTION.MD.
2. **Migrations beat markdown** for schema questions.
3. **Config beats markdown** for reward/XP/gem numbers (`config/rewards.go`, `config/levels.go`).
4. **One reference doc max** per task domain — do not chain-read all summaries.
5. **Stop reading when STATE.md + RULES.md + one domain doc answer the question.**

---

## Trim Targets for `.ai/` Canon (when executing this plan)

| File | Current lines | Target lines | What to cut |
|------|-------------:|-------------:|-------------|
| `PRODUCT.MD` | 653 | 100 | Pillars summary table; cut psychology, visual identity essays |
| `ARCHITECTURE.MD` | 722 | 200 | Keep ownership, events, API, DB rules; cut philosophy, decision tree prose |
| `CONTENT_SYSTEM.MD` | 708 | 180 | Keep hierarchy, evaluation types, readiness layers; cut examples, future AI sections |
| `EXECUTION.MD` | 630 | 0 → `ROADMAP.md` 80 | Move completed A1–A4 to STATE; keep only Phase B/C/D locked specs |

---

## Token Savings Estimate

Assumptions: ~4 characters per token; typical task previously reads 4 `.ai/` files + 1–2 summaries + always-on AGENTS rules.

| Scenario | Before | After | Reduction |
|----------|-------:|------:|----------:|
| **Always-on context** (Cursor rules) | ~3,000 | ~750 | **75%** |
| **Typical backend task** (4 `.ai/` + RULES) | ~9,500 | ~2,000 | **79%** |
| **Phase implementation task** (4 `.ai/` + 2 summaries) | ~15,500 | ~2,600 | **83%** |
| **Over-cautious task** (all docs + audits) | ~47,500 | ~3,500 | **93%** |
| **Full corpus available** | ~48,000 | ~9,500 active | **80%** corpus retired to archive |

### Per-conversation savings

| Metric | Estimate |
|--------|----------|
| Average doc tokens saved per task | **~12,000 tokens** |
| Average always-on savings (rule slimming) | **~2,250 tokens per message** |
| 10-task session savings | **~120,000+ tokens** |

### Overall

| | |
|--|--|
| **Estimated token reduction** | **80–93%** on documentation context |
| **Decision quality preserved by** | STATE.md (fresh status), RULES.md (non-negotiables), migrations (schema truth), targeted reference docs |

---

## Implementation Checklist (documentation ops only)

This plan does not require code changes. Recommended execution order:

1. [ ] Create `.ai/STATE.md` from template above; update completion status.
2. [ ] Create `.ai/RULES.md` — extract non-negotiables from Cursor rules + ARCHITECTURE anti-patterns.
3. [ ] Trim `.ai/ARCHITECTURE.MD`, `CONTENT_SYSTEM.MD`, `PRODUCT.MD` per cut targets (or create new trimmed files, rename old to archive).
4. [ ] Split `.ai/EXECUTION.MD` → `.ai/ROADMAP.md` (Phase B+ only); archive task specs now in STATE.
5. [ ] Merge reference docs under `docs/reference/`.
6. [ ] Move audits, completed summaries, legacy root copies to `docs/archive/`.
7. [ ] Update Cursor workspace rules to point at STATE + RULES only (~750 tokens always-on).
8. [ ] Add `docs/archive/README.md` one-liner: *"Historical only — do not use for implementation decisions."*
9. [ ] Delete or redirect root `PRODUCT.md` / `EXECUTION.md` to prevent agent confusion.

---

## Risks of Optimization

| Risk | Mitigation |
|------|------------|
| Agents miss product philosophy nuance | `.ai/PRODUCT.md` trim retains pillars + anti-patterns; load for UX tasks |
| Lost migration UUID context | `docs/reference/CONTENT_ARCHITECTURE.md` keeps Foundation Forest UUIDs |
| Over-trimming ARCHITECTURE breaks ownership rules | RULES.md retains the 5 hard non-negotiables verbatim |
| STATE.md goes stale | Update STATE.md as part of every phase completion PR |
| Developers read archive by habit | Archive README + STATE.md "never read" list |

---

## Summary

| Item | Action |
|------|--------|
| **Single source of truth for status** | `.ai/STATE.md` (new) |
| **Single source for non-negotiables** | `.ai/RULES.md` (new, slim Cursor rules) |
| **Merge** | Readiness docs, Journey docs, A2 schema into `docs/reference/` |
| **Archive** | 7 audits, 5 completed summaries, 2 stale root copies |
| **Trim** | `.ai/` canon from 2,713 → ~760 active lines |
| **Never read again** | `audit result/*`, completed summaries, proposals, stale root docs |
| **Token reduction** | **80–93%** on documentation per task |

Architectural correctness is preserved by keeping **RULES.md** (hard constraints), **migrations** (schema truth), and **on-demand reference docs** (formulas, flags, UUIDs) — while eliminating repeated philosophy, obsolete gap analysis, and completed phase narratives from every agent session.
