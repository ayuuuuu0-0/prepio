# Prepio

Prepio is a lesson-based progression platform for working engineers: short, instantly graded lessons in system design, backend and production, low-level design, and DSA, with per-topic mastery. The repository contains a Go backend, a Next.js web client, and a Flutter mobile client (paused; see `.ai/EXECUTION.MD`).

## Repository layout

- `services/` - Go services for gateway, user, question, streak, progress, and notification domains
- `shared/` - shared packages used across services
- `migrations/` - Postgres schema migrations
- `web/` - Next.js application
- `mobile/` - Flutter application (paused)
- `scripts/` - local development, deployment, and validation scripts
- `config/` and `constants/` - runtime configuration and shared values
- `.ai/` - product, architecture, content, and execution documents (the source of truth)
- `agent/` - index for AI agents

## Local development

Prerequisites:

- Go
- Docker and Docker Compose
- Node.js for the web client
- Flutter for the mobile client (optional while paused)

Start local infrastructure:

```bash
make docker-up
```

Run the backend services:

```bash
make dev
```

Run the web client:

```bash
cd web
npm install
npm run dev
```

## Common commands

```bash
make build-all   # build all Go packages
make test        # run the full Go test suite
make test-short  # run the fast gateway and shared tests
make vet         # run go vet across the repo
make migrate-up  # apply database migrations
make e2e         # run end-to-end service validation
make docker-up   # start postgres, redis, and kafka
make docker-up-all
make docker-down
```

## Services

| Service | Port | Purpose |
| --- | --- | --- |
| gateway | 8080 | Entry point, request routing, response aggregation |
| user | 8081 | Authentication, profile, preferences, companion state |
| question | 8082 | Content and journey: skills, lessons, attempts, grading (historical name) |
| streak | 8083 | Daily check-ins and streak tracking |
| progress | 8084 | XP, levels, skill mastery, topic readiness |
| notification | 8085 | Event-driven notifications |

Supporting infrastructure is Postgres, Redis, and Kafka.

## Production deployment

Production runs on a single host with Docker Compose (`docker-compose.prod.yml`, `Dockerfile`, `Caddyfile`, `scripts/aws-setup.sh`). Copy `.env.example` to `.env` and set strong secrets.

## Engineering context

Product, architecture, content, and execution rules live in [`.ai/`](.ai/). Start with [`agent/README.md`](agent/README.md) for the read order.
