# Prepio

Prepio is a lesson-based progression platform for working engineers: short, instantly graded lessons in system design, backend and production, low-level design, and DSA, with per-topic mastery. The repository contains a Go backend, a Next.js web client, and a Flutter mobile client (paused; see `.ai/EXECUTION.MD`).

## Repository layout

- `services/` - Go services for gateway, user, question, streak, progress, and notification domains
- `shared/` - shared packages used across services
- `migrations/` - Postgres schema migrations
- `web/` - Next.js application
- `mobile/` - Flutter application (paused)
- `scripts/` - local development, deployment, and validation scripts
- `content/` - authored worlds and lessons (content-as-code), loaded by `content-sync`
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
make test-docker  # run the whole Go suite inside a Linux container (use this on Windows with Smart App Control)
make vet         # run go vet across the repo
make migrate-up  # apply database migrations
make e2e         # drive the lesson flow over HTTP against a running gateway (BASE_URL, default :8080)
make content-validate  # validate authored lessons without writing
make content-sync      # load worlds and lessons into the database (idempotent; runs on every deploy)
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
| progress | 8084 | XP, levels, skill mastery, topic readiness, weekly leagues |
| notification | 8085 | Event-driven notifications |

Supporting infrastructure is Postgres, Redis, and Kafka.

## Production deployment

Production runs on a single host with Docker Compose (`docker-compose.prod.yml`, `Dockerfile`, `Caddyfile`, `scripts/aws-setup.sh`). Copy `.env.example` to `.env` and set strong secrets.

## Future: landing page

Not part of the current phase: a marketing or landing page is out of scope until the lesson path is validated (see `OUT OF SCOPE` in `.ai/EXECUTION.MD`). When it is time, it should feel like Duolingo's: one clear promise, a friendly mascot (our companions), and a single obvious call to action. And it needs to look really cool, so take its visual lead from [zeptap.com](https://zeptap.com/). The components worth borrowing:

- **A window with several phone screens playing at once.** A macOS-style app window holding multiple phones, each showing live motion (moving bars, progress, a lesson being answered). For Prepio, show real lesson screens: the journey, a question with the feedback tray, the celebration with mastery moving.
- **Soft, tinted cards with a live illustration inside.** Each feature card has its own pastel gradient and a small product-style visual (a phone, a lock, a receipt). Use one card per topic (System Design, Backend & Production, Low-Level Design, DSA Refresher), each with its own tint and a mastery ring as the illustration.
- **Big, tight headlines with the second half muted.** For example "If you can do it on your iPhone, *your agents can too.*" Prepio's version is a bold claim in white followed by the supporting half in grey.
- **Floating sticker icons around the hero.** Chunky 3D-style icons drifting at the edges. Ours would be the companions, gems, streak flame, and badges.
- **A scrolling strip of example chips.** Small pills of example tasks that slide past, here the skills and lessons learners practice.
- **A floating pill navigation and an announcement pill above the headline.**

Keep Prepio's identity while borrowing the patterns: the dark theme with purple accents from `.ai/PRODUCT.MD` stays (Zeptap is light sky-blue; copy the structure, not the palette). Respect `prefers-reduced-motion`, and build it from real components and real screens, never mock data.

## Engineering context

Product, architecture, content, and execution rules live in [`.ai/`](.ai/). Start with [`agent/README.md`](agent/README.md) for the read order.
