# Prepio — Architecture

> This document is the single source of truth for system design, data ownership,
> API contracts, and infrastructure decisions. Read it before touching any service.
> When this document and the code disagree, fix the code.

---

## System Overview

```
                          ┌─────────────────────────────────┐
                          │           Clients                │
                          │  Flutter App   │   Next.js Web   │
                          └──────┬─────────────────┬─────────┘
                                 │                 │
                          ┌──────▼─────────────────▼─────────┐
                          │           API Gateway             │
                          │   Auth · Rate Limit · Routing     │
                          │         (Go, port 8080)           │
                          └──┬────┬─────┬──────┬─────┬───────┘
                             │    │     │      │     │
               ┌─────────────┘    │     │      │     └──────────────┐
               │            ┌─────┘     └──┐   └───────┐            │
               ▼            ▼             ▼            ▼            ▼
         ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐
         │  User    │ │ Question │ │  Streak  │ │ Progress │ │  Notif   │
         │ Service  │ │ Service  │ │ Service  │ │ Service  │ │ Service  │
         │ :8081    │ │ :8082    │ │ :8083    │ │ :8084    │ │ :8085    │
         └────┬─────┘ └────┬─────┘ └────┬─────┘ └────┬─────┘ └────┬─────┘
              │             │            │             │             │
              └─────────────┴────────────┴─────────────┴─────────────┘
                                         │
                              ┌──────────▼──────────┐
                              │        Kafka         │
                              │  (event bus)         │
                              └──────────────────────┘
                                         │
                    ┌────────────────────┼────────────────────┐
                    ▼                    ▼                     ▼
             ┌────────────┐     ┌──────────────┐     ┌──────────────┐
             │ PostgreSQL │     │    Redis      │     │  Object      │
             │ (primary)  │     │ (cache/sets)  │     │  Storage     │
             └────────────┘     └──────────────┘      └──────────────┘
```

---

## Services

### API Gateway (`services/gateway/`)

Responsibilities: JWT validation, rate limiting per user ID, request routing
to the correct upstream service, response passthrough.

Does not contain business logic. No database access except a Redis lookup
for token blacklisting on logout.

Rate limits:
- Authenticated: 300 requests per minute per user
- Unauthenticated: 20 requests per minute per IP

### User Service (`services/user/`)

Owns: `users`, `user_devices`, `character_unlocks`, `user_characters`

Responsibilities: registration, login, JWT issuance, profile management,
FCM token storage, character unlock transactions.

Character unlock is a synchronous operation: the user service calls the
progress service internal API to verify and deduct gem balance atomically
before writing the unlock record. This is the only synchronous inter-service
call in the system.

### Question Service (`services/question/`)

Owns: `questions`, `question_tags`, `user_question_history`, `daily_papers`

Responsibilities: question bank management, daily paper generation per user,
answer submission, answer evaluation, weekend challenge content.

On every answer submission, produces a `question.answered` Kafka event.
Does not award XP or update streaks.

### Streak Service (`services/streak/`)

Owns: `user_streaks`, `streak_freeze_inventory`

Responsibilities: consuming `question.answered` events, determining streak
eligibility by user local timezone, updating streak counts, consuming freeze
purchases, producing `streak.updated` events.

Streak state is cached in Redis with a 36-hour TTL. PostgreSQL is the
write-through source of truth.

### Progress Service (`services/progress/`)

Owns: `user_progress`, `xp_ledger`, `gem_ledger`, `user_levels`

Responsibilities: consuming `question.answered` and `streak.updated` events,
awarding XP and gems per config rules, level-up detection, exposing gem
balance for the user service's character unlock call.

Internal API (not exposed via gateway): `GET /internal/progress/{userID}/gems`
and `POST /internal/progress/{userID}/gems/deduct`. Used only by the user
service for character unlock.

### Notification Service (`services/notification/`)

Owns: `notification_log`, `notification_preferences`

Responsibilities: consuming all Kafka events, selecting the appropriate
character dialogue, enforcing the 3-per-day per-user cap, dispatching via
FCM for mobile and email for web, logging every sent notification.

Never produces Kafka events. Terminal consumer only.

---

## Kafka Topics

All topics use JSON payloads. Schema is documented below each topic.
Consumers are in the same consumer group as their service name.

### `question.answered`

Produced by: question service  
Consumed by: streak service, progress service

```json
{
  "event_id": "uuid",
  "user_id": "uuid",
  "question_id": "uuid",
  "round_type": "dsa",
  "difficulty": "medium",
  "company_tags": ["google", "meta"],
  "correct": true,
  "submitted_at": "2026-06-09T22:14:00+05:30",
  "session_id": "uuid"
}
```

### `streak.updated`

Produced by: streak service  
Consumed by: progress service, notification service

```json
{
  "event_id": "uuid",
  "user_id": "uuid",
  "previous_streak": 4,
  "current_streak": 5,
  "streak_broken": false,
  "freeze_consumed": false,
  "updated_at": "2026-06-09T22:14:01+05:30"
}
```

### `progress.updated`

Produced by: progress service  
Consumed by: notification service

```json
{
  "event_id": "uuid",
  "user_id": "uuid",
  "xp_awarded": 50,
  "gems_awarded": 10,
  "total_xp": 1240,
  "total_gems": 85,
  "level_before": 4,
  "level_after": 5,
  "leveled_up": true,
  "updated_at": "2026-06-09T22:14:02+05:30"
}
```

### `notifications.dispatch`

Produced by: any service that needs to trigger a notification  
Consumed by: notification service

```json
{
  "event_id": "uuid",
  "user_id": "uuid",
  "notification_type": "streak_reminder",
  "metadata": {},
  "triggered_at": "2026-06-09T21:50:00+05:30"
}
```

Notification types: `streak_reminder`, `streak_broken`, `level_up`,
`league_position_change`, `weekend_challenge_available`, `streak_freeze_low`.

---

## Database Schema

### PostgreSQL

#### `users`
```sql
CREATE TABLE users (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           TEXT NOT NULL UNIQUE,
    username        TEXT NOT NULL UNIQUE,
    password_hash   TEXT NOT NULL,
    timezone        TEXT NOT NULL DEFAULT 'Asia/Kolkata',
    active_char_id  UUID REFERENCES characters(id),
    reminder_time   TIME NOT NULL DEFAULT '21:50:00',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

#### `user_devices`
```sql
CREATE TABLE user_devices (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    fcm_token   TEXT NOT NULL,
    platform    TEXT NOT NULL CHECK (platform IN ('android', 'ios', 'web')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, platform)
);
```

#### `characters`
```sql
CREATE TABLE characters (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            TEXT NOT NULL,
    species         TEXT NOT NULL,
    gem_cost        INT NOT NULL DEFAULT 0,
    is_default      BOOLEAN NOT NULL DEFAULT false,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

#### `character_dialogues`
```sql
CREATE TABLE character_dialogues (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    character_id        UUID NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    notification_type   TEXT NOT NULL,
    dialogue_line       TEXT NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

#### `character_unlocks`
```sql
CREATE TABLE character_unlocks (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    character_id    UUID NOT NULL REFERENCES characters(id),
    unlocked_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, character_id)
);
```

#### `questions`
```sql
CREATE TABLE questions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    body            TEXT NOT NULL,
    round_type      TEXT NOT NULL CHECK (round_type IN (
                        'dsa', 'system_design', 'lld',
                        'aptitude', 'fundamentals', 'behavioral'
                    )),
    difficulty      TEXT NOT NULL CHECK (difficulty IN ('easy', 'medium', 'hard')),
    answer_guide    TEXT NOT NULL,
    status          TEXT NOT NULL DEFAULT 'pending'
                        CHECK (status IN ('pending', 'approved', 'retired')),
    is_weekend      BOOLEAN NOT NULL DEFAULT false,
    source          TEXT NOT NULL CHECK (source IN ('manual', 'ai_generated', 'scraped')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

#### `question_tags`
```sql
CREATE TABLE question_tags (
    question_id     UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    company         TEXT NOT NULL,
    PRIMARY KEY (question_id, company)
);
```

#### `user_question_history`
```sql
CREATE TABLE user_question_history (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    question_id     UUID NOT NULL REFERENCES questions(id),
    correct         BOOLEAN NOT NULL,
    submitted_at    TIMESTAMPTZ NOT NULL,
    received_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    session_id      UUID NOT NULL,
    UNIQUE (user_id, question_id, session_id)
);
```

#### `user_streaks`
```sql
CREATE TABLE user_streaks (
    user_id             UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    current_streak      INT NOT NULL DEFAULT 0,
    longest_streak      INT NOT NULL DEFAULT 0,
    last_activity_date  DATE,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

#### `streak_freeze_inventory`
```sql
CREATE TABLE streak_freeze_inventory (
    user_id     UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    count       INT NOT NULL DEFAULT 0,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

#### `user_progress`
```sql
CREATE TABLE user_progress (
    user_id         UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    total_xp        INT NOT NULL DEFAULT 0,
    current_level   INT NOT NULL DEFAULT 1,
    gem_balance     INT NOT NULL DEFAULT 0,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

#### `xp_ledger`
```sql
CREATE TABLE xp_ledger (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount          INT NOT NULL,
    reason          TEXT NOT NULL,
    source_event_id UUID NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

#### `gem_ledger`
```sql
CREATE TABLE gem_ledger (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount          INT NOT NULL,
    reason          TEXT NOT NULL,
    source_event_id UUID NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

#### `notification_log`
```sql
CREATE TABLE notification_log (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    notification_type   TEXT NOT NULL,
    channel             TEXT NOT NULL CHECK (channel IN ('fcm', 'email')),
    sent_at             TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

---

### Redis Key Schema

```
streak:{userID}                     → hash: current_streak, last_activity_date, freeze_count
                                      TTL: 36 hours

leaderboard:weekly:{weekISO}        → sorted set: score = weekly XP, member = userID
                                      TTL: 8 days

leaderboard:level:{levelBand}       → sorted set: score = weekly XP, member = userID
                                      TTL: 8 days

session:{sessionID}                 → string: userID
                                      TTL: 24 hours

token_blacklist:{jti}               → string: "1"
                                      TTL: token remaining lifetime

notif_cap:{userID}:{dateYYYYMMDD}   → string: count of notifications sent today
                                      TTL: 36 hours

gems:{userID}                       → string: gem balance (read cache only)
                                      TTL: 5 minutes
```

---

## API Contracts

All routes are prefixed `/api/v1/`. All responses follow this envelope:

**Success:**
```json
{ "data": { ... } }
```

**Error:**
```json
{ "error": { "code": "error_code_constant", "message": "lowercase description" } }
```

**Pagination (cursor-based):**
```json
{
  "data": [ ... ],
  "pagination": {
    "next_cursor": "opaque_string",
    "has_more": true
  }
}
```

---

### Auth

```
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/logout
POST   /api/v1/auth/refresh
```

`register` and `login` return:
```json
{
  "data": {
    "access_token": "jwt",
    "refresh_token": "jwt",
    "user": { "id": "uuid", "username": "string", "email": "string" }
  }
}
```

---

### Users

```
GET    /api/v1/users/me
PATCH  /api/v1/users/me
GET    /api/v1/users/me/characters
POST   /api/v1/users/me/characters/{characterID}/unlock
PUT    /api/v1/users/me/characters/{characterID}/activate
POST   /api/v1/users/me/devices
DELETE /api/v1/users/me/devices/{deviceID}
```

---

### Questions

```
GET    /api/v1/questions/daily          → today's paper for the authenticated user
POST   /api/v1/questions/{id}/submit    → submit an answer
GET    /api/v1/questions/history        → cursor-paginated answered history
GET    /api/v1/questions/companies      → list of available company tags
```

`daily` response:
```json
{
  "data": {
    "session_id": "uuid",
    "date": "2026-06-09",
    "questions": [
      {
        "id": "uuid",
        "body": "string",
        "round_type": "dsa",
        "difficulty": "medium",
        "company_tags": ["google"],
        "is_weekend": false
      }
    ],
    "minimum_to_streak": 1
  }
}
```

`submit` request:
```json
{
  "session_id": "uuid",
  "answer": "string",
  "time_spent_seconds": 840
}
```

`submit` response:
```json
{
  "data": {
    "correct": true,
    "xp_awarded": 50,
    "gems_awarded": 10,
    "streak_updated": true,
    "feedback": "string"
  }
}
```

---

### Streaks

```
GET    /api/v1/streaks/me
POST   /api/v1/streaks/me/freeze/purchase
```

`me` response:
```json
{
  "data": {
    "current_streak": 12,
    "longest_streak": 34,
    "freeze_count": 1,
    "last_activity_date": "2026-06-09",
    "streak_active_today": true
  }
}
```

---

### Progress

```
GET    /api/v1/progress/me
GET    /api/v1/progress/me/xp/history
GET    /api/v1/progress/me/gems/history
```

`me` response:
```json
{
  "data": {
    "total_xp": 1240,
    "current_level": 5,
    "gem_balance": 85,
    "xp_to_next_level": 260
  }
}
```

---

### Leaderboard

```
GET    /api/v1/leaderboard/weekly          → current week, user's level band
GET    /api/v1/leaderboard/weekly/global   → current week, all users
```

WebSocket: `WS /api/v1/leaderboard/live`  
Emits `leaderboard.update` events when the user's rank changes by 3 or more
positions. JSON frame:
```json
{
  "type": "leaderboard.update",
  "data": {
    "previous_rank": 8,
    "current_rank": 4,
    "weekly_xp": 340
  }
}
```

---

### Characters (public catalogue)

```
GET    /api/v1/characters
GET    /api/v1/characters/{id}
```

---

## Directory Structure

```
prepio/
├── services/
│   ├── gateway/
│   ├── user/
│   │   ├── cmd/main.go
│   │   ├── internal/
│   │   │   ├── handler/
│   │   │   ├── store/
│   │   │   └── dto/
│   │   └── Dockerfile
│   ├── question/
│   ├── streak/
│   ├── progress/
│   └── notification/
├── shared/
│   ├── kafka/           # producer and consumer helpers
│   ├── jwt/             # token generation and validation
│   ├── postgres/        # connection pool setup
│   ├── redis/           # client setup
│   └── events/          # Kafka event type definitions and schemas
├── migrations/
│   ├── 001_create_users.sql
│   ├── 002_create_characters.sql
│   ├── 003_create_questions.sql
│   └── ...
├── config/
│   ├── rewards.go       # XP amounts, gem awards, streak freeze cost
│   └── levels.go        # XP thresholds per level
├── constants/
│   └── errors.go        # all API error code constants
├── test/
│   ├── factories/
│   ├── smoke/
│   └── integration/
├── mobile/              # Flutter app
│   ├── lib/
│   │   ├── features/
│   │   │   ├── auth/
│   │   │   ├── streak/
│   │   │   ├── question/
│   │   │   ├── progress/
│   │   │   ├── leaderboard/
│   │   │   └── characters/
│   │   ├── core/
│   │   │   ├── api/
│   │   │   ├── storage/     # Hive offline queue
│   │   │   └── notifications/
│   │   └── l10n/
│   └── pubspec.yaml
├── web/                 # Next.js frontend
│   ├── app/
│   ├── components/
│   └── package.json
├── docker-compose.yml
└── AGENTS.md
```

Each service follows the same internal layout:
```
service/
├── cmd/
│   └── main.go          # wires dependencies, starts server
├── internal/
│   ├── handler/         # HTTP handlers, one file per route group
│   ├── store/           # all SQL queries, one file per table
│   ├── dto/             # request and response structs
│   └── service/         # business logic, one file per domain concept
└── Dockerfile
```

---

## Infrastructure

**Local development:** `docker-compose.yml` at repo root spins up PostgreSQL,
Redis, Kafka (single broker), and all five services with live reload via Air.

**Production:** Each service is a standalone Docker image. Deploy to any
container platform. Services communicate over internal network only except
the gateway which is public-facing.

**Secrets:** All credentials are environment variables. No secrets in code,
no secrets in `docker-compose.yml` committed to the repo. Use a `.env` file
locally, a secrets manager in production.

**Migrations:** Run via `migrate` CLI before any service starts.
`make migrate-up` in the Makefile handles this.

---

## Scalability Decisions and Why

**Kafka over direct HTTP between services:** Services can be restarted,
redeployed, or scaled independently without message loss. The streak service
can replay events if it goes down during a peak period.

**Cursor pagination on leaderboards:** The leaderboard sorted set in Redis
is updated continuously. Offset pagination would return inconsistent results
as users move up and down. Cursors are stable.

**Gem deduction as synchronous internal call:** Character unlocks cannot be
async. A user clicks unlock, expects immediate feedback, and cannot spend the
same gems twice. The internal call between user service and progress service
uses optimistic locking on the gem balance row.

**Offline answer queue in Flutter:** Submissions with a `submitted_at`
timestamp are accepted up to 48 hours after the fact. The question service
uses `submitted_at`, not `received_at`, for streak eligibility. This handles
subway commutes and bad connectivity without punishing users.

**Leaderboard scoped to level band:** Global leaderboards are demoralising
for new users. Competing against people at a similar level keeps the league
meaningful and competitive.