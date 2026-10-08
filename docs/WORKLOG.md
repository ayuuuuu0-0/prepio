# Work log

A plain-language record of what was built, why, how it was tested, and what is still open. Newest work at the bottom of each section. The rules the project follows live in `.ai/`; this file only explains the work.

Branch: `feat/l9-leagues-landing` (not pushed yet).

---

## Done before this log started

| Commit | What | In one line |
|---|---|---|
| `7b4d384` | L5.1 hardening | Fixed review findings in the merged L2–L5 code: logouts on reload, lost XP under concurrent events, content-sync crash loop, streak resets, dead Kafka consumers, swallowed migration failures, dead code. |
| `cb961e2` | L9 leagues | Weekly leagues owned by Progress: cohorts of 30, promotion/demotion at week end, `/league` screen and nav tab. |
| `d247e21` | L10 landing page | Public `/` page in Duolingo's structure with Zeptap-style components, dark theme. |
| `59fdf69` | L6 exercise types | Player renders true/false, fill-in-the-blank, and arrange, tap-only. |
| `3d4e786` | L7 retention polish | Mistakes review, sound (off by default), inert exit dialog, Back asks before leaving. |
| `e4dc964` | L8 content | 33 lessons in four worlds, planned in `content/PLAN.md`, reviewed for accuracy and answer tells. |

---

## Session: further features (2026-10-08)

How I work through this list: one task at a time. Each task is implemented, reviewed, tested, documented here, and committed on its own before the next starts.

### Task queue

1. **Live end-to-end run** of the full stack, so the new exercise types, the 33 lessons, and leagues are proven against real services, not just unit tests. *(status below)*
2. **Self-host the web fonts**, so `next build` stops failing when Google Fonts is slow.
3. **Focus topics order the journey**: the path shows your focus topics' worlds first (needs worlds to know their topic).
4. **Prose exercises in the player**, so rubric-graded written answers (already graded by the server) can be used in lessons.
5. **League card on the dashboard**, so "what reward is next" includes your weekly rank.

Anything outside `.ai/EXECUTION.MD`'s scope goes through its Change Process and is recorded here.

---

### Task 1 · Live end-to-end run — done (`ff3cdcb`)

**Why:** everything up to L8 was proven by unit, integration, and screenshot tests, but never with all services running together on the real content.

**How:** a throwaway local stack (Postgres on 55432, Redis on 56379, all six Go services, the web app on 3000) with a fresh database, all 41 migrations, and the real `content/`. Nothing in the repo or your installed Postgres/Redis services is touched. It is torn down at the end of every task that starts it, and the ports are checked free afterwards.

**What ran and what it showed**

1. `scripts/lesson-smoke.mjs` (the HTTP smoke test). First run: 2 of 21 checks failed. Both were the script's own assumptions from when the path had one lesson with every answer at option 0 (it expected mastery to always rise and "nothing to continue" after one lesson). Not product bugs. Fixed and committed; now 21/21 pass.
2. A full-path playthrough over the API as a fresh user: all **33 lessons, 144 exercises** (76 multiple choice, 35 true/false, 18 fill-in, 15 arrange). Every answer the server revealed after a miss was accepted on retry, so no lesson's answer key contradicts its grader. 29 unlocks fired (one per chained node); dashboard XP = league weekly XP = 879; focus topics listed first; all four topics gained mastery; skills touched: System Design 4/4, Backend 4/4, LLD 3/3, DSA 9/14 (as planned).
3. The real web app in a headless browser at phone size, signed in through the real refresh cookie: dashboard, journey, league with live standings, and a true/false, a fill-in, and an arrange exercise played with taps through to server-graded feedback. **No console errors.** Screenshots looked right.

**Things learned**

* The gateway rate limit (300 requests/minute per user) stopped the automated playthrough once; a person can't reach it, so the script waits and retries instead.
* **A mistake I made and corrected:** earlier today I said nothing was running in the background. In fact a stack I had started (the command you rejected had already begun) kept running from 07:53 until this task. It was idle, but I told you otherwise. It is stopped, and the teardown script now stops only what listens on its own ports, never anything by process name.
