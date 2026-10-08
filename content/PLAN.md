# Working Professionals path: content plan (Phase L8)

The authoring plan for the first path. The rules live in `.ai/CONTENT_SYSTEM.MD`; this file is the outline those rules ask for ("the world/node outline lives in `content/`"). The YAML files beside it are the source of truth for what ships; this plan says what to write and why.

## Goal

Forty lessons, one per day (33 in the first release, 7 more in World 5), across the four topics, so the path supports the launch criterion of 30+ days of daily use. Every lesson teaches one idea a working engineer uses in design reviews, on call, or in interviews, and measures it with 3–5 server-graded exercises.

## Constraints (from the docs and the validator)

* World → Node → Lesson → Step. A node binds exactly one lesson; node slug = lesson slug.
* Lessons: ≤ 5 min (boss ≤ 8), `summary` 2–4 takeaways, one optional `intro` first, then 3–5 graded steps (boss 3–8). Skill weights sum to 1.0 and use **existing** skill slugs only (the catalog changes by migration, not here).
* Intro: 2–4 beats, each 1.5–12 s, whole intro ≤ 45 s; `visual` only `bars`, `flow`, or `pulse`.
* Exercise types the player renders: `mcq`, `true_false`, `fill_blank`, `arrange`, and (since L6.1) `prose`. Prose stays rare (PRODUCT.MD) and needs a careful rubric: see the authoring standard.
* `mcq`: 3–4 options, a "why not" for every wrong option. `arrange`: ≥ 3 items, authored scrambled. `fill_blank`: numbered blanks, the bank holds every answer plus distractors.
* No company content. Professional tone: encouraging, never shaming; explanations say *why*, not just *what*.
* The existing lesson `why-caches-exist` keeps its slug and node so learners keep their progress; its options are reordered only (see Plan review).

## Path structure and unlock rules

Five worlds: one per topic, seven or eight lessons then a boss, plus a second DSA world (40 nodes in all).

* **Within a world** each node requires the previous one, so the world reads as a climb. The boss requires the last lesson (and therefore all of them).
* **Across topics** nothing is locked: the first world of each topic has no `requires` on its first node. A *sequel* world inside the same topic (World 5) may require the previous world's boss, because that only orders one topic's own material. Topics "never lock content" (PRODUCT.MD), so a learner who focuses on DSA is never stuck behind System Design. The path still has a single *current* node (the first unlocked, unfinished node in world order); the other open nodes are *available*.
* World order: System Design, Backend & Production, Low-Level Design, DSA Refresher, matching the topic catalog order.

## Outline

Skill slugs are from the seeded catalog (migrations 000026 and 000035). Weights are per lesson.

### World 1 · First Ascent (System Design) · theme `summit`
"Your first climb: how real systems stay fast under load."

Skill slugs: `system-design-scaling`, `system-design-data`, `system-design-reliability`, `system-design-fundamentals`, plus `backend-databases`.

| # | Node / lesson slug | Title | Skills | Teaches |
|---|---|---|---|---|
| 1 | `why-caches-exist` | Why Caches Exist | scaling 1.0 | *(exists)* cache-aside, staleness, what not to cache |
| 2 | `keeping-caches-honest` | Keeping Caches Honest | scaling 1.0 | TTL vs invalidation, write-through vs write-around, cache stampede |
| 3 | `spreading-the-load` | Spreading the Load | scaling 1.0 | load balancer job, round robin vs least connections, health checks, stateless servers |
| 4 | `finding-rows-fast` | Finding Rows Fast | data 0.7, backend-databases 0.3 | B-tree index, composite index left-prefix, write cost of indexes |
| 5 | `copies-that-lag` | Copies That Lag | reliability 0.6, data 0.4 | leader/follower replication, replication lag, read-your-writes |
| 6 | `splitting-the-data` | Splitting the Data | scaling 0.6, data 0.4 | sharding, shard key choice, hot partitions, consistent hashing |
| 7 | `work-that-can-wait` | Work That Can Wait | reliability 0.6, scaling 0.4 | queues, async work, at-least-once delivery, idempotent consumers |
| 8 | `back-of-the-envelope` | Back of the Envelope | fundamentals 1.0 | requirements first, QPS and storage estimates, read/write ratio |
| 9 | `boss-read-heavy-feed` | Boss: Scale a Read-Heavy Feed | scaling 0.4, data 0.2, reliability 0.2, fundamentals 0.2 | mixed review across the world (7 steps) |

### World 2 · The Engine Room (Backend & Production) · theme `engine`
"Where requests become work: APIs, data, and keeping services alive."

| # | Slug | Title | Skills | Teaches |
|---|---|---|---|---|
| 1 | `apis-that-make-sense` | APIs That Make Sense | api 1.0 | resources and verbs, status codes (4xx vs 5xx), safe vs idempotent methods |
| 2 | `retry-without-fear` | Retry Without Fear | api 0.6, reliability 0.4 | idempotency keys, why POST retries double-charge, server-side dedup |
| 3 | `page-by-page` | Page by Page | api 0.7, databases 0.3 | offset vs cursor pagination, stable ordering, deep-page cost |
| 4 | `all-or-nothing` | All or Nothing | databases 1.0 | transactions, ACID, isolation levels and the anomalies they allow |
| 5 | `the-lost-update` | The Lost Update | concurrency 0.6, databases 0.4 | read-modify-write races, row locks, optimistic concurrency with versions |
| 6 | `timeouts-and-backoff` | Timeouts and Backoff | concurrency 0.6, observability 0.4 | timeouts, exponential backoff with jitter, retry storms |
| 7 | `seeing-inside` | Seeing Inside | observability 1.0 | logs vs metrics vs traces, golden signals, useful alerts |
| 8 | `boss-incident-at-3am` | Boss: Incident at 3 a.m. | observability 0.3, concurrency 0.25, databases 0.25, api 0.2 | mixed review as an incident walkthrough (7 steps) |

Skill slugs: `backend-api-design`, `backend-databases`, `backend-concurrency`, `backend-observability`, plus `system-design-reliability`.

### World 3 · Blueprint Workshop (Low-Level Design) · theme `workshop`
"Classes, boundaries, and patterns that survive the next feature request."

| # | Slug | Title | Skills | Teaches |
|---|---|---|---|---|
| 1 | `guard-your-state` | Guard Your State | oop 0.6, fundamentals 0.4 | encapsulation, invariants, why public setters leak rules |
| 2 | `one-reason-to-change` | One Reason to Change | oop 1.0 | single responsibility, spotting god classes |
| 3 | `open-for-extension` | Open for Extension | patterns 0.5, oop 0.5 | open/closed with the strategy pattern, replacing type switches |
| 4 | `compose-dont-inherit` | Compose, Don't Inherit | oop 1.0 | composition over inheritance, fragile base classes, Liskov |
| 5 | `depend-on-abstractions` | Depend on Abstractions | oop 0.6, fundamentals 0.4 | dependency inversion, injection, testable seams |
| 6 | `building-objects` | Building Objects | patterns 1.0 | factory vs builder, when a constructor is enough |
| 7 | `when-things-happen` | When Things Happen | patterns 1.0 | observer / events, decoupling producers from consumers |
| 8 | `boss-parking-lot` | Boss: Design a Parking Lot | fundamentals 0.4, oop 0.3, patterns 0.3 | mixed review on one classic design (7 steps) |

Skill slugs: `lld-oop`, `lld-patterns`, `lld-fundamentals`.

### World 4 · Algorithm Gardens (DSA Refresher) · theme `garden`
"The structures and techniques you actually reach for, without the grind."

| # | Slug | Title | Skills | Teaches |
|---|---|---|---|---|
| 1 | `the-right-lookup` | The Right Lookup | hash-maps 1.0 | O(1) lookup, counting, trading memory for time |
| 2 | `two-pointers` | Two Pointers | arrays 1.0 | sorted-array pair search, in-place partitioning |
| 3 | `the-sliding-window` | The Sliding Window | sliding-window 1.0 | fixed and variable windows, when a window applies |
| 4 | `stacks-and-queues` | Stacks and Queues | stacks 0.5, queues 0.5 | LIFO vs FIFO, matching brackets, BFS order |
| 5 | `halving-the-search` | Halving the Search | binary-search 1.0 | binary search, off-by-one boundaries, searching on the answer |
| 6 | `walking-trees` | Walking Trees | trees 0.6, recursion 0.4 | DFS orders, BFS by level, recursion depth |
| 7 | `finding-the-way` | Finding the Way | graphs 1.0 | BFS shortest path in unweighted graphs, visited sets, cycles |
| 8 | `boss-pick-the-structure` | Boss: Pick the Structure | hash-maps 0.3, arrays 0.2, trees 0.25, graphs 0.25 | choose the right tool for seven scenarios (7 steps) |

### World 5 · Deep Roots (DSA Refresher, part 2) · theme `grove`
"The techniques that make hard problems small: strings, lists, heaps, greedy choices, and dynamic programming."

Added after the first four worlds shipped, to cover the five DSA skills they left untouched (`strings`, `linked-lists`, `heaps`, `greedy`, `dynamic-programming`). It is a sequel inside one topic, so its first node **requires the Algorithm Gardens boss**; that never locks another topic. `topic: dsa-refresher`, `order: 5`.

| # | Slug | Title | Skills | Teaches |
|---|---|---|---|---|
| 1 | `strings-without-surprises` | Strings Without Surprises | strings 1.0 | immutable strings, why repeated concatenation is quadratic, join/builders, counting characters |
| 2 | `following-the-links` | Following the Links | linked-lists 1.0 | when lists beat arrays, fast/slow pointers for cycles and middles, reversing in place |
| 3 | `always-the-smallest` | Always the Smallest | heaps 1.0 | priority queues, O(log n) push/pop, top-k with a size-k heap in O(n log k) |
| 4 | `take-the-best-now` | Take the Best Now | greedy 1.0 | interval scheduling by earliest finish, when greedy fails (coin systems), proving a choice is safe |
| 5 | `remember-the-answers` | Remember the Answers | dynamic-programming 0.7, recursion 0.3 | overlapping subproblems, memoization, exponential vs linear |
| 6 | `fill-the-table` | Fill the Table | dynamic-programming 1.0 | bottom-up tables, defining the state, the recurrence, and the base case |
| 7 | `boss-choose-the-technique` | Boss: Choose the Technique | strings 0.2, linked-lists 0.2, heaps 0.2, greedy 0.2, dynamic-programming 0.2 | seven scenarios, one technique each (7 steps) |

With World 5 every skill in all four topics is exercised by the path.

## Authoring standard (every lesson)

* **Intro:** 3 beats, 12–16 s total, one `emphasis` phrase per beat, visuals chosen to match the idea (`bars` for comparisons, `flow` for pipelines, `pulse` for timing and events).
* **Exercises:** 4 graded steps for a lesson (boss 7), using at least two types per lesson and every type across a world. Use `fill_blank` only where a short code or config line is natural; keep code to about 6 lines. A fill-in step has no prompt, so when the right answer depends on the goal, the first line is a comment stating it (without giving the answer away). `arrange` for ordered processes. Distractors must be plausible mistakes real engineers make, never jokes.
* **No answer by position or pattern:** options are sent in authored order, so spread each world's correct `mcq` answers across positions (never three in a row at the same index) and don't make the correct option systematically the longest; keep `true_false` answers roughly half `false`.
* **Plain text:** no markdown or backticks (the player renders text literally); an intro `emphasis` must match its beat text exactly, including case; each fill-in blank appears once.
* **Prose (rare):** the rubric grader matches concepts by name and aliases, so list the real synonyms a correct answer would use (e.g. stale: "out of date", "outdated"), mark only the essential ideas `required`, and write the `explanation` as the worked answer the learner sees after grading. Try a weak and a strong sample answer against it before publishing.
* **Explanations:** the correct explanation says why it works; each wrong-option explanation names the misconception kindly. Facts must be correct and current; avoid vendor-specific claims.
* **Summary:** usually 3 takeaways (the docs allow 2–4), each one sentence, each answerable from the lesson.
* **Option length:** the correct option must not be a tell: across a world it is the longest option in no more than about 45% of `mcq`s (chance is 33% with three options).
* **Difficulty:** worlds open with `easy`, middle lessons `medium`, bosses `medium` or `hard`.

## Delivery order and gates

1. **World 1** (already started with `why-caches-exist`). Done when it validates and passes review.
2. **Worlds 2–4**, authored against the World 1 template.
3. Gates for every world: `content-sync --validate`; `go test ./test/integration/...` (validates real content and resolves skill slugs against the migrated catalog); a full read-through review for accuracy, tone, ambiguity (exactly one defensible answer), and that no answer is guessable from wording alone.
4. `.ai/EXECUTION.MD` L8 status updated in the same change.

## Review checklist (applied to this plan and to every lesson)

* Every skill slug exists; weights sum to 1.0; every skill in the System Design, Backend & Production, and Low-Level Design topics is exercised by its world. DSA Refresher covered its core skills in World 4; World 5 adds `strings`, `linked-lists`, `heaps`, `greedy`, and `dynamic-programming`.
* Unlock chain: no cycles, each world's first node open, boss last.
* No `prose` steps; no company names; no hearts/lives language.
* One idea per lesson; no lesson depends on a later one.
* Bosses only review what their world taught.

## Plan review (done before authoring)

* **Skill coverage:** the first draft never exercised `system-design-fundamentals`, so a whole System Design skill could not move. Added *Back of the Envelope* as World 1 node 8 and gave the boss a fundamentals weight.
* **Unlocks:** checked against PRODUCT.MD (topics never lock content): worlds open in parallel, chains only inside a world.
* **Player support:** `prose` excluded because the player cannot render it yet; every other type is used.
* **Existing content:** `why-caches-exist` stays World 1 node 1 with the same slug, so learners who finished it keep their progress.
* **Validator limits** checked against `services/question/internal/lesson/validate.go`: beats 1.5–12 s, intro ≤ 45 s, 3–5 graded steps (boss ≤ 8), boss ≤ 8 min.
* **Answer patterns (found while writing World 1):** every first draft put the correct option first and most true/false answers were true, which a learner could exploit. World 1 was rebalanced, including `why-caches-exist` (which becomes version 2: attempts in progress finish on version 1, and completions are kept because they are recorded per lesson), and the authoring standard now forbids it.
* **Fill-in context (found in review):** `fill_blank` has no prompt field, so six snippets whose answers depended on an unstated goal (binary search bounds, BFS order, the cursor direction, the alert window, jitter, the two-pointer goal) now open with a goal comment.
* **Theme:** `theme` is a free-form label today (nothing renders it), so new themes need no code.

## Out of scope for L8 (follow-ups)

* ~~Focus topics reordering the journey~~: done in L11 (worlds now name their `topic`).
* ~~`prose` exercises in the player~~: done in L6.1. No lesson uses prose yet; adding a few to bosses is a content decision.
* A second path for other audiences.
