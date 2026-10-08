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

---

### Task 2 · Self-host the web fonts — waiting for your OK

**Why:** `next build` downloads the three Google Fonts (Sora, Manrope, IBM Plex Mono) at build time; once today it failed until retried. Shipping the font files in the repo (`next/font/local`) makes builds independent of Google.

**Why it's paused:** it needs the font files downloaded once (about 10 `.woff2` files, roughly 300 KB, from Google Fonts or the npm `@fontsource` packages). I don't download files without your go-ahead. Say "OK, download the fonts" and I'll do it as its own task. Until then the build still works; it just depends on Google being reachable.

---

### Task 3 · Focus topics order the journey (L11) — done (`3af9771`)

**Why:** onboarding (PRODUCT.MD) promises that the 1–3 topics you pick "reorder and highlight the path". Until now they only reordered the dashboard's topic cards: worlds didn't know their topic, so the journey and the Continue button always followed the fixed world order.

**What changed** (recorded as a Change Process entry in `.ai/EXECUTION.MD`, L11):

* Each world names its topic (`topic:` in the world file → `worlds.topic_id`, migration 000042). Content-sync checks it against the topic catalog.
* `GET /path?focus=dsa-refresher,system-design` puts those topics' worlds first, in your priority order, and marks them `focused`. This happens *before* node status is worked out, so the one "current" node (what Continue opens) follows your focus.
* The dashboard asks for the path with your profile's focus topics, so Continue and the journey always agree.
* The journey shows a "Your focus" label on those worlds.
* **Nothing is ever locked by focus:** unlocks still depend only on finishing the previous node. Unfocused worlds stay open ("available"). This is tested.

**Tested:**
* Unit: ordering (no focus, one, two, unknown topic; nodes inside a world never reorder) and `focus` input checks (max 3, valid slugs, no duplicates).
* Gateway: the dashboard really sends your focus to the path.
* Integration on the real 33-lesson content: every world has a topic; with focus the focused world is first and current, others available, none hidden.
* Migrations: the rollback test used a hard-coded list and broke as soon as a newer migration depended on topics. It now rolls back every migration from 000035 up, read from the folder, so future migrations are always covered. 042 → 035 down and back up passes.
* Live, full stack: a user focused on DSA then System Design gets Continue → *Algorithm Gardens / The Right Lookup*, and the journey screenshot shows Algorithm Gardens then First Ascent labelled "Your focus". Bad `focus` input → 400. No console errors.

**A slip, caught and checked:** while verifying, I ran `content-sync` once without pointing it at the test database, so it defaulted to `localhost:5432` (your installed Postgres). It was refused at login (`password authentication failed`) before doing anything, and sync runs in one transaction anyway, so nothing was written. I now always pass the test database explicitly.

**Clean-up:** stack and web server stopped; ports 3000, 8080–8085, 55432, 56379 confirmed free.

---

### Task 4 · Written answers (prose) in the player (L6.1) — done (`1d3596e`)

**Why:** PRODUCT.MD lists "rare rubric-graded prose" as an exercise type, and the server has graded it since L2, but the player showed "this exercise isn't available" instead. Content was barred from using it.

**What changed:**
* A **text area** for written answers with a live counter ("62 more characters to go" → "✓ Long enough to check"); Check stays disabled until the minimum length is met.
* **Keyboard:** Enter adds a new line (you're writing a paragraph); Ctrl+Enter (Cmd+Enter on Mac) checks.
* The **feedback tray** for written answers shows the rubric score, a "You covered" list, a "Worth adding" list, and "A strong answer:" (the worked answer).
* The server's low-score message used to say "focus on the approach, complexity, and tradeoffs" (left over from the old coding-interview grader). It now reads "Keep going: the ideas listed below are what a strong answer adds."
* `content/PLAN.md` now allows prose and explains how to write a rubric that grades fairly (list real synonyms, mark only essentials as required, test a weak and a strong answer).

**Tested:**
* Unit: prose drafts (empty start, minimum length ignores surrounding spaces, the wire answer, the mistakes-review prompt). 68 web tests pass; lint, typecheck, build clean; full question-service suite passes.
* Live, full stack, with a **test-only** prose lesson synced into the throwaway database (never into `content/`): too short → rejected with "write at least 100 characters"; a weak answer → score 33, covered "cache", worth adding "stale", "latency"; a strong answer → 100. The rubric itself never reaches the browser.
* Real browser at phone size: Check disabled while short; Enter doesn't submit; Ctrl+Enter does; the tray reads clearly (I added the visible "You covered" / "Worth adding" labels after seeing the first screenshot). No console errors.

**Not done (a content decision for you):** no real lesson uses a written answer yet. A few in the boss lessons would fit "rare", but each needs a carefully tested rubric; the plan explains how.

**Clean-up:** stack and web server stopped; all ports confirmed free.
