---
title: Prepio Backend
emoji: 🚀
colorFrom: indigo
colorTo: blue
sdk: docker
pinned: false
---

# 🎮 Prepio

> **Duolingo sells progress. Pokemon sells collection. Clash Royale sells advancement. Prepio sells career progression.**

Prepio is not an interview preparation platform. It is a progression game where becoming interview-ready is the reward. We believe developers don't wake up wanting to grind dry LeetCode questions; they wake up wanting better jobs, higher salaries, and real career growth. Prepio transforms that journey into a visual, addictive loop.

---

## 🔁 The Core Loop

```text
       Open App
          │
          ▼
   Companion Greets
          │
          ▼
    Check Readiness
          │
          ▼
    Continue Journey
          │
          ▼
   Complete Challenge
          │
          ▼
      Gain XP/Gems
          │
          ▼
   Evolve Companion
          │
          ▼
    Unlock Worlds
```

---

## 🕹️ Core Game Mechanics

* **Ready-Score™**: The ultimate metric. Not how many questions you solved, but how ready you are for specific targets (e.g., *Google: 74%*, *Amazon: 68%*). Every score is explainable and backed by real-time skill gaps.
* **The Skill Graph**: Questions are just content; Skills (e.g., *Two Pointers*, *Sliding Window*, *Prefix Sum*) are the foundations. Prepio maps your knowledge base directly.
* **Worlds**: Ditch categories. Progress through **Foundation Forest**, scale **Google Mountain**, or navigate **System Design City**.
* **Companions**: Emotional progression systems. They react, evolve, celebrate, and push you to stay consistent.
* **Collection**: Titles, badges, companion skins, and world trophies to make retention natural and fun.

---

## 🏗️ Architecture: The Engine Room

Prepio is built as a highly decoupled, event-driven microservices architecture in **Go**, paired with a **Next.js** web frontend and a **Flutter** mobile client.

```text
                     ┌──────────────────┐
                     │   Web / Mobile   │
                     └────────┬─────────┘
                              │ (REST)
                              ▼
                     ┌──────────────────┐
                     │     Gateway      │
                     └────────┬─────────┘
                              │
          ┌───────────┬───────┴───┬───────────┐
          ▼           ▼           ▼           ▼
      ┌───────┐   ┌───────┐   ┌───────┐   ┌───────┐
      │ User  │   │Quest'n│   │Streak │   │Progr's│
      └───────┘   └───────┘   └───────┘   └───────┘
          │           │           │           │
          └───────────┼───────────┼───────────┘
                      ▼ (Kafka Events)
                 ┌──────────┐
                 │ Notif.   │
                 └──────────┘
```

### Services Directory

| Service | Port | Description | DB / Cache / Message Queue |
| :--- | :--- | :--- | :--- |
| **`gateway`** | `8080` | Entrypoint & reverse proxy routing requests to appropriate services. | - |
| **`user`** | `8081` | Authentication, profile management, and Companion state. | Postgres, Redis |
| **`question`** | `8082` | Question bank, pools, skills, and daily paper generation. | Postgres, Kafka |
| **`streak`** | `8083` | Tracking daily check-ins and gamified streaks. | Postgres, Redis, Kafka |
| **`progress`** | `8084` | User readiness engine and node completion state. | Postgres, Kafka |
| **`notification`**| `8085` | Event-driven notification dispatch (streak freezes, levels).| Postgres, Redis, Kafka |

---

## 🚀 Booting the Game (Local Development)

### Prerequisites

* **Go** (1.21+)
* **Docker** & **Docker Compose**
* **Node.js** (for web)
* **Flutter** (for mobile)

### 1. Fire up Infrastructure
Spin up Postgres, Redis, and Kafka in the background:
```bash
make docker-up
```

### 2. Launch Backend Microservices
Run all 6 Go services concurrently (outputs are directed to `.run/*.log`):
```bash
make dev
```
To stop the services, just press `Ctrl + C` (it traps the signal and cleans up the running PIDs).

### 3. Run the Frontend (Next.js)
```bash
cd web
npm install
npm run dev
```

### 4. Run the Mobile App (Flutter)
```bash
cd mobile
flutter pub get
flutter run
```

---

## 🛠️ Developer Playbook

We use a simple `Makefile` to orchestrate common commands:

* **`make build-all`**: Compile all Go packages to verify build health.
* **`make test`**: Run the backend unit and integration test suite.
* **`make test-short`**: Run fast gateway and shared package tests.
* **`make migrate-up`**: Apply Postgres database migrations.
* **`make e2e`**: Execute the end-to-end service validation tests.

---

## 🎨 Visual Philosophy
Never build **government software**, **enterprise SaaS**, or **dry admin dashboards**. 

Prepio should feel like a polished indie game. Every interface must contain:
1. **Progress**: Visible growth bars or levels.
2. **Color & Motion**: Vibrant highlights, smooth transitions, no static stagnation.
3. **Companion Presence**: Your ally is always visible, reacting to your actions.
4. **Instant Feedback**: Every correct response or level up must trigger celebratory visual/audio cues.
