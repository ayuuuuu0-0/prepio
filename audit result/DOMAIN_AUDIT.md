# Domain Audit — Prepio Backend vs `.ai/ARCHITECTURE.md`

**Date:** 2026-06-09  
**Scope:** All backend services, gateway, shared packages, migrations  
**Reference:** `.ai/ARCHITECTURE.md`, `AGENTS.md` (legacy 5-service model)

---

## Executive Summary

The codebase implements a **5-service monolith-style split** (User, Question, Streak, Progress, Notification) plus Gateway, defined in `AGENTS.md`. The `.ai/ARCHITECTURE.md` target defines **6 domains**: User, Content, Journey, Progress, Companion, Analytics.

**Critical finding:** Question service is a **god service** owning questions, daily papers, evaluation, history, journey, and partial readiness. Journey and readiness belong elsewhere. Content, Companion, and Analytics domains are **missing or fragmented**.

---

## 1. Domain Mapping — Current Code

### User Domain

| Component | Location | Responsibility |
|-----------|----------|----------------|
| Auth (register/login/refresh) | `services/user/` | ✅ Correct |
| Profile, timezone | `services/user/` | ✅ Correct |
| Onboarding, experience level | `services/user/` | ✅ Correct |
| Target companies (`user_targets`) | `services/user/` | ✅ Correct |
| Character unlocks | `services/user/` | ⚠️ Partial — companions split with notification |
| FCM token storage | `services/user/` + notification store | ⚠️ Split |

**Tables:** `users`, `user_targets`, `user_characters`, `characters`, `character_dialogues`

---

### Content Domain (Target: owns question bank, pools, skills, authoring)

| Component | Location | Should Be |
|-----------|----------|-----------|
| Question CRUD/read | `services/question/` | Content |
| Question selection | `services/question/internal/service/question.go` | Content (pools) |
| Daily paper generation | `services/question/` | **Journey or Content delivery** |
| Evaluation/rubric | `services/question/internal/service/evaluator.go` | Content (evaluation rules) |
| Seed data | `migrations/` | Content |
| Skills, pools | **Missing** | Content |

**Verdict:** Content domain **does not exist**. Question service is a stand-in.

---

### Journey Domain (Target: worlds, nodes, unlock rules, progression path)

| Component | Location | Should Be |
|-----------|----------|-----------|
| Worlds, nodes schema | `migrations/000024` | Journey |
| GetJourney API | `services/question/internal/service/journey.go` | **Journey service** |
| user_journey_progress writes | `services/question/` (on GET journey) | **Journey service** |
| Node unlock logic | Derived from daily paper index | **Journey service** |
| Daily paper as session | `services/question/` | Overlaps Journey + Content |

**Verdict:** Journey is **embedded in Question service** with incorrect dependencies.

---

### Progress Domain (Target: XP, gems, levels, skill mastery, readiness)

| Component | Location | Should Be |
|-----------|----------|-----------|
| XP, gems, levels | `services/progress/` | ✅ Correct |
| Kafka consumer (question.answered, streak.updated) | `services/progress/` | ✅ Correct |
| Readiness computation | `services/gateway/internal/dashboard/` + `services/question/` | **Progress** |
| readiness_delta on submit | `services/question/internal/service/question.go` | **Progress** |
| Company performance stats | `services/question/internal/store/history.go` | **Progress or Analytics** |
| user_question_history writes | `services/question/` | **Progress or Analytics** |
| Skill mastery | **Missing** | Progress |

**Verdict:** Progress owns gamification but **not readiness or skill scores**.

---

### Companion Domain (Target: character selection, dialogue, notification personality)

| Component | Location | Should Be |
|-----------|----------|-----------|
| Character catalog | `services/user/` | Companion |
| Character unlock (gem check) | User + Progress sync call | Split (acceptable) |
| Dialogue lines | DB `character_dialogues` | Companion |
| Dialogue selection for notifications | **Stub** — notification service logs only | Companion |
| Dashboard companion message | `services/gateway/internal/dashboard/service.go` | **Companion** |

**Verdict:** Companion domain **fragmented** across User, Gateway, Notification.

---

### Analytics Domain (Target: aggregates, league, insights, explainability)

| Component | Location | Should Be |
|-----------|----------|-----------|
| Question history queries | `services/question/` | Analytics |
| Readiness stats API | `services/question/` | Analytics (or Progress read API) |
| Leaderboard | **Placeholder** in gateway | Analytics |
| League position notifications | **Not implemented** | Analytics + Notification |
| Skill gap analysis | **Missing** | Analytics |

**Verdict:** Analytics domain **does not exist**.

---

## 2. Service Topology — Current vs Target

```
CURRENT (AGENTS.md)                    TARGET (.ai/ARCHITECTURE.md)
─────────────────────                  ─────────────────────────────
Gateway (8080)                         API Gateway
  ├─ auth routing                        ├─ auth routing
  ├─ dashboard aggregation               └─ routing only (no business logic)
  └─ readiness computation ❌
User (8081)                            User Service
Question (8082) ❌ god service         Content Service
  ├─ questions                           Journey Service
  ├─ daily papers                        Progress Service
  ├─ evaluation                          Companion Service
  ├─ history                             Analytics Service
  ├─ journey
  └─ readiness stats
Streak (8083)                          (unchanged — streak is cross-cutting)
Progress (8084)                        (expand: readiness, skill scores)
Notification (8085)                    (unchanged — dispatch only)
```

---

## 3. Ownership Violations

| # | Violation | Severity | Current | Owner Should Be |
|---|-----------|----------|---------|-----------------|
| V1 | Journey logic in Question service | **Critical** | `journey.go`, `000024` handlers | Journey |
| V2 | Readiness computed in Gateway | **Critical** | `dashboard/service.go:computeReadiness` | Progress |
| V3 | Readiness stats in Question service | **High** | `history.go`, `/stats/readiness` | Progress (compute) / Analytics (expose) |
| V4 | readiness_delta on submit in Question | **High** | `question.go` submit handler | Progress |
| V5 | Journey progress written during GET | **High** | Side effect in read path | Journey |
| V6 | user_question_history owned by Question | **Medium** | Question service writes | Progress or Analytics (Question emits event only) |
| V7 | Companion microcopy in Gateway | **Medium** | `companionMessage()` hardcoded | Companion |
| V8 | XP/gems computed in Question on submit | **Medium** | `rewards.go` in question service | Progress (should consume event and award) |
| V9 | Streak eligibility uses question.answered | **Low** | Correct per AGENTS.md | OK for now |
| V10 | Target company list in Question selection comment only | **High** | Not implemented | Content (pool filtering) |

---

## 4. Duplicate Responsibilities

| Responsibility | Locations | Problem |
|----------------|-----------|---------|
| Readiness | Gateway + Question submit + Question stats API | Three formulas/paths, none skill-based |
| Session/progression | Daily paper + Journey overlay | Two parallel progression models |
| Company targeting | User (stores targets), Question (ignores), Gateway (displays) | Data collected but unused in selection |
| Reward calculation | Question service (sync on submit) + Progress service (Kafka) | Question returns xp/gems in response; Progress also awards via events — **potential double-award risk** |
| Character/companion UX | User service, Gateway dashboard, Notification stub | No single companion voice |

---

## 5. Business Logic in Wrong Places

### Gateway (`services/gateway/`)

- `computeReadiness()` — business formula
- `companionMessage()` — companion dialogue selection
- `comingSoonQuests()` — quest content
- Dashboard aggregation joins multiple services

**AGENTS.md says:** "API Gateway authenticates requests and routes them. It never contains business logic."

**`.ai/ARCHITECTURE.md` says:** Gateway routes only.

### Question Service

- Journey node status computation
- Reward estimation (`computeRewards`)
- Readiness delta calculation
- Daily paper = progression session

### Progress Service

- Only consumes Kafka; does not expose readiness API
- Does not own answer history

---

## 6. Services Doing Too Much

### Question Service — **Critical**

**Owns today:**
- Question bank reads
- Daily paper lifecycle
- Answer evaluation
- Answer submission + history
- Journey worlds/nodes/progress
- Readiness statistics
- Reward preview in response

**Should own (Content slice only):**
- Question/pool/skill reads
- Evaluation (or shared evaluator lib)
- Emit `question.answered` / `question.completed` events

**Estimated extraction:** 60% of question service code moves to Journey, Progress, or Analytics.

---

### Gateway — **High**

**Owns today:**
- Auth proxy
- Dashboard BFF with readiness, quests, league placeholders, companion copy

**Should own:**
- Auth proxy + routing only
- Optional thin BFF that **calls domain APIs without computing**

---

## 7. Missing Domains

| Domain | Status | Phase A Action |
|--------|--------|----------------|
| Content | Missing as service | Extract from Question; add skills/pools |
| Journey | Missing as service | Extract journey.go + progress tables |
| Companion | Missing as service | Extract from User + Gateway; wire notification dialogue |
| Analytics | Missing | New read-only service or Progress read APIs |

**Note:** EXECUTION.md Phase A does **not** require splitting into 6 microservices immediately. Domain boundaries can be **packages within monorepo** first, then extract services when stable.

---

## 8. Missing Abstractions

| Abstraction | Purpose | Status |
|-------------|---------|--------|
| `SkillRepository` | Skill graph reads | Missing |
| `QuestionPoolSelector` | Pool-based question pick | Missing (only difficulty+random) |
| `SkillMasteryCalculator` | Update mastery on answer | Missing |
| `ReadinessEngine` | Skill → company → overall | Missing |
| `JourneyProgressionEngine` | Node unlock rules | Missing (index hack) |
| `ContentAuthoringService` | CRUD + review queue | Missing |
| Event: `skill.progress.updated` | Progress → Analytics | Missing |
| Event: `readiness.updated` | Progress → Notification | Missing |
| Event: `journey.node.completed` | Journey → Progress/Notification | Missing |

**Existing events** (`shared/events/events.go`):
- `question.answered` ✅
- `streak.updated` ✅
- `progress.updated` ✅
- `notifications.dispatch` ✅

---

## 9. Cross-Cutting Concerns

| Concern | Current | Recommendation |
|---------|---------|----------------|
| Evaluator | Question service package | Shared lib `shared/evaluation/` — used by Content, called on submit |
| Rewards config | `config/rewards.go` ✅ | Keep centralized |
| Error codes | `constants/errors.go` ✅ | Extend for skill/journey/readiness |
| Auth | Gateway + User ✅ | Unchanged |
| Kafka | Shared producer/consumer ✅ | Add new event types |

---

## 10. Issue Register

| ID | Issue | Severity | Proposed Solution | Migration Difficulty |
|----|-------|----------|-------------------|----------------------|
| D1 | Question service god object | Critical | Split into Content + Journey packages; extract readiness to Progress | **Hard** — phased extraction |
| D2 | Readiness in Gateway | Critical | Move `computeReadiness` to Progress; Gateway calls `GET /progress/readiness` | **Medium** |
| D3 | Journey in Question | Critical | Move journey handlers/store to `services/journey/` or `question/internal/journey/` package with clear boundary | **Medium** |
| D4 | Daily paper drives progression | High | Introduce pool-based node selection; daily paper becomes optional "quick play" | **Hard** |
| D5 | Double reward path | High | Question returns preview only OR Progress is sole awarder; audit Kafka flow | **Medium** |
| D6 | History in Question | Medium | Question emits event; Progress/Analytics persists history | **Medium** |
| D7 | No skill abstraction | Critical | Implement A1 skill graph per EXECUTION.md | **Medium** |
| D8 | Companion in Gateway | Low | Move to Companion module; read active character + pick dialogue | **Easy** |
| D9 | No Analytics service | Medium | Start as Progress read APIs + materialized views | **Medium** |
| D10 | AGENTS.md vs .ai/ docs conflict | Medium | Update AGENTS.md after Phase A to reflect 6 domains | **Easy** (docs) |
| D11 | Target companies unused in selection | High | Wire user_targets into pool selection | **Medium** |
| D12 | Side effects on GET journey | High | Separate `POST /journey/sync` or compute read-only; persist on events | **Easy** |

---

## 11. Recommended Domain Boundary (Phase A — Package Level)

```
services/
├── gateway/          # Route only; thin dashboard aggregation
├── user/             # User domain
├── content/          # NEW — questions, pools, skills, evaluation (split from question)
├── journey/          # NEW — worlds, nodes, unlock, user journey progress
├── progress/         # XP, gems, levels, skill scores, readiness
├── companion/        # NEW — character voice, dialogue selection (optional Phase A late)
├── analytics/        # NEW — history queries, stats (optional Phase A late)
├── streak/           # Unchanged
└── notification/     # Unchanged
```

**Pragmatic Phase A approach:** Rename/refactor `services/question` → keep binary but **internal packages** mirror domains. Extract services in Phase B if needed.

---

## 12. Kafka Event Flow — Current vs Target

### Current

```
Question.Submit → question.answered → Streak, Progress
Streak → streak.updated → Progress
Progress → progress.updated → Notification
```

### Target (Phase A additions)

```
Question.Submit → question.completed (with skill_ids, score)
                → Progress updates skill mastery
                → Progress emits skill.progress.updated
                → Progress emits readiness.updated
                → Journey consumes → journey.node.completed (if threshold met)
                → Notification consumes readiness/journey events
```

---

## 13. Testing Implications

| Test | Location | Domain gap |
|------|----------|------------|
| Core loop integration | `test/integration/core_loop_test.go` | Tests old model only |
| Streak unit tests | streak service | OK |
| Skill/readiness/journey tests | **Missing** | Required in Phase A |

---

*This document is analysis only. No code was modified.*
