# Prepio

Prepio is a progression-focused interview preparation platform. The repository contains a Go backend, a Next.js web client, and a Flutter mobile client.

## Repository layout

- `services/` - Go services for gateway, user, question, streak, progress, and notification domains
- `shared/` - shared packages used across services
- `migrations/` - Postgres schema migrations
- `web/` - Next.js application
- `mobile/` - Flutter application
- `scripts/` - local development and validation scripts
- `config/` and `constants/` - runtime configuration and shared values
- `agent/` - consolidated AI context and working notes

## Local development

Prerequisites:

- Go
- Docker and Docker Compose
- Node.js for the web client
- Flutter for the mobile client

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

Run the mobile client:

```bash
cd mobile
flutter pub get
flutter run
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
| gateway | 8080 | Entry point and request routing |
| user | 8081 | Authentication, profile, companion state |
| question | 8082 | Question bank, skills, and daily paper generation |
| streak | 8083 | Daily check-ins and streak tracking |
| progress | 8084 | Readiness and progression state |
| notification | 8085 | Event-driven notifications |

Supporting infrastructure is Postgres, Redis, and Kafka.

## Engineering context

The product is built around progression, readiness, skills, and journey-based content. The canonical implementation and agent-facing notes live in [agent/README.md](agent/README.md).
