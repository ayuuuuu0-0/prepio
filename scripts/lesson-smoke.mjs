#!/usr/bin/env node
// Drives the lesson platform through the gateway the way the web client will:
// register -> path -> start attempt -> answer (wrong, replay, right) -> complete -> rewards.
//
//   BASE_URL=http://localhost:8080 node scripts/lesson-smoke.mjs
//
// It discovers answers through the API (a wrong try reveals the correct answer),
// so it needs no copy of the lesson content.
import assert from "node:assert/strict";

const BASE = (process.env.BASE_URL ?? "http://localhost:8080") + "/api/v1";
let failures = 0;

async function call(method, path, { token, body, expect } = {}) {
  const res = await fetch(BASE + path, {
    method,
    headers: { "Content-Type": "application/json", ...(token ? { Authorization: `Bearer ${token}` } : {}) },
    body: body === undefined ? undefined : JSON.stringify(body),
  });
  const text = await res.text();
  let json;
  try {
    json = JSON.parse(text);
  } catch {
    json = undefined;
  }
  if (expect !== undefined) {
    assert.equal(res.status, expect, `${method} ${path} -> ${res.status}: ${text.slice(0, 300)}`);
  }
  return { status: res.status, json, text };
}

async function step(name, fn) {
  try {
    await fn();
    console.log(`  ok   ${name}`);
  } catch (err) {
    failures++;
    console.log(`  FAIL ${name}\n       ${String(err.message).split("\n")[0]}`);
  }
}

const suffix = Math.random().toString(36).slice(2, 10);
const creds = { email: `smoke-${suffix}@example.test`, username: `smoke_${suffix}`, password: "password123" };
let token;
let other;
let node;
let attempt;
let completion;

console.log(`lesson smoke test against ${BASE}`);

await step("register returns tokens", async () => {
  const r = await call("POST", "/auth/register", { body: creds, expect: 201 });
  token = r.json.data.access_token;
  assert.ok(token);
  const o = await call("POST", "/auth/register", {
    body: { email: `other-${suffix}@example.test`, username: `other_${suffix}`, password: "password123" },
    expect: 201,
  });
  other = o.json.data.access_token;
});

await step("the topic catalog lists the four topics", async () => {
  const r = await call("GET", "/topics", { token, expect: 200 });
  assert.deepEqual(r.json.data.map((t) => t.slug), ["system-design", "backend-production", "low-level-design", "dsa-refresher"]);
});

await step("onboarding stores 1-3 focus topics and rejects bad input", async () => {
  const companions = (await call("GET", "/companions", { expect: 200 })).json.data;
  const body = (topics) => ({ experience_level: "mid", companion_id: companions[0].id, focus_topics: topics });
  await call("POST", "/users/onboarding", { token, body: body([]), expect: 400 });
  await call("POST", "/users/onboarding", { token, body: body(["system-design", "backend-production", "low-level-design", "dsa-refresher"]), expect: 400 });
  await call("POST", "/users/onboarding", { token, body: body(["astrology"]), expect: 400 });
  const ok = await call("POST", "/users/onboarding", { token, body: body(["dsa-refresher", "backend-production"]), expect: 200 });
  assert.deepEqual(ok.json.data.focus_topics, ["dsa-refresher", "backend-production"]);
});

await step("the dashboard shows focus topics first, not started, and a Continue target", async () => {
  const home = (await call("GET", "/dashboard/home", { token, expect: 200 })).json.data;
  assert.equal(home.onboarding_needed, false);
  assert.deepEqual(home.topics.map((t) => t.slug), ["dsa-refresher", "backend-production", "system-design", "low-level-design"]);
  assert.equal(home.topics[0].focused, true);
  assert.equal(home.topics[2].focused, false);
  for (const t of home.topics) {
    assert.equal(t.mastery, null, t.slug + " has no mastery before any lesson");
    assert.equal(t.skills_started, 0);
    assert.ok(t.skills_total > 0);
  }
  assert.ok(home.next_lesson && home.next_lesson.lesson_id && home.next_lesson.title);
  assert.equal("league" in home, false);
  assert.equal("daily_quests" in home, false);
});

await step("lesson endpoints require authentication", async () => {
  await call("GET", "/path", { expect: 401 });
  await call("POST", "/lessons/00000000-0000-4000-8000-000000000000/attempts", { expect: 401 });
});

await step("path shows a current node with a preview", async () => {
  const r = await call("GET", "/path", { token, expect: 200 });
  const worlds = r.json.data.worlds;
  assert.ok(worlds.length >= 1, "at least one world");
  node = worlds[0].nodes[0];
  assert.equal(node.status, "current");
  assert.ok(node.title && node.takeaways.length >= 2 && node.est_minutes > 0 && node.xp_preview > 0);
});

await step("unknown lesson is 404", async () => {
  await call("POST", "/lessons/00000000-0000-4000-8000-000000000000/attempts", { token, expect: 404 });
});

await step("starting an attempt returns steps with no answers", async () => {
  const r = await call("POST", `/lessons/${node.lesson_id}/attempts`, { token, expect: 201 });
  attempt = r.json.data;
  assert.ok(attempt.attempt_id && attempt.steps.length >= 4);
  for (const leaked of ['"answer"', '"explanation', '"rubric"', '"wrong"', '"correct"']) {
    assert.ok(!r.text.includes(leaked), `payload must not contain ${leaked}`);
  }
});

await step("starting again resumes the same attempt (200)", async () => {
  const r = await call("POST", `/lessons/${node.lesson_id}/attempts`, { token, expect: 200 });
  assert.equal(r.json.data.attempt_id, attempt.attempt_id);
  assert.equal(r.json.data.resumed, true);
});

await step("another user cannot read or answer the attempt", async () => {
  const first = attempt.steps.find((s) => s.type !== "intro");
  await call("POST", `/attempts/${attempt.attempt_id}/steps/${first.id}/answer`, {
    token: other, body: { try: 1, answer: { choice: 0 } }, expect: 404,
  });
  await call("POST", `/attempts/${attempt.attempt_id}/complete`, { token: other, expect: 404 });
});

const graded = () => attempt.steps.filter((s) => s.type !== "intro");

await step("malformed answers are rejected with 400", async () => {
  const first = graded()[0];
  await call("POST", `/attempts/${attempt.attempt_id}/steps/${first.id}/answer`, {
    token, body: { try: 1, answer: { choice: 99 } }, expect: 400,
  });
  await call("POST", `/attempts/${attempt.attempt_id}/steps/${first.id}/answer`, {
    token, body: { try: 1, answer: {} }, expect: 400,
  });
  await call("POST", `/attempts/${attempt.attempt_id}/steps/${first.id}/answer`, {
    token, body: "not-json", expect: 400,
  });
});

await step("completing early is rejected with 409", async () => {
  await call("POST", `/attempts/${attempt.attempt_id}/complete`, { token, expect: 409 });
});

// firstGuess builds a well-formed first answer for any deterministic step type, so the script
// works whatever exercises the first lesson uses; a miss reveals the correct answer to retry with.
function firstGuess(step) {
  switch (step.type) {
    case "mcq": return { choice: 0 };
    case "true_false": return { value: true };
    case "fill_blank": return { blanks: step.fill_blank.bank.slice(0, step.fill_blank.blanks) };
    case "arrange": return { order: step.arrange.items.map((_, i) => i) };
    default: throw new Error("smoke test cannot answer step type " + step.type);
  }
}

await step("a wrong first try reveals the answer; replay is idempotent; the right try passes", async () => {
  const first = graded()[0];
  const path = `/attempts/${attempt.attempt_id}/steps/${first.id}/answer`;
  const wrongChoice = 1; // a deliberate guess; if it happens to be right the wrong-path checks are skipped
  let r = await call("POST", path, { token, body: { try: 1, answer: { choice: wrongChoice } }, expect: 200 });
  let truth = r.json.data;
  if (truth.correct) return;
  assert.ok(truth.explanation && truth.correct_answer && typeof truth.correct_answer.choice === "number");
  const again = await call("POST", path, { token, body: { try: 1, answer: { choice: wrongChoice } }, expect: 200 });
  assert.equal(again.json.data.replayed, true);
  await call("POST", path, { token, body: { try: 1, answer: { choice: truth.correct_answer.choice } }, expect: 409 });
  const right = await call("POST", path, { token, body: { try: 2, answer: { choice: truth.correct_answer.choice } }, expect: 200 });
  assert.equal(right.json.data.correct, true);
  assert.equal(right.json.data.correct_answer, undefined, "no answer reveal on a correct try");
  await call("POST", path, { token, body: { try: 3, answer: { choice: truth.correct_answer.choice } }, expect: 409 });
});

await step("the remaining steps can be answered (answers discovered through feedback)", async () => {
  for (const s of graded().slice(1)) {
    const path = `/attempts/${attempt.attempt_id}/steps/${s.id}/answer`;
    let r = await call("POST", path, { token, body: { try: 1, answer: firstGuess(s) }, expect: 200 });
    if (!r.json.data.correct) {
      r = await call("POST", path, { token, body: { try: 2, answer: r.json.data.correct_answer }, expect: 200 });
      assert.equal(r.json.data.correct, true);
    }
  }
});

await step("completion returns a summary with rewards from Progress", async () => {
  const r = await call("POST", `/attempts/${attempt.attempt_id}/complete`, { token, expect: 200 });
  completion = r.json.data;
  assert.equal(completion.graded_steps, graded().length);
  assert.ok(completion.accuracy >= 0 && completion.accuracy <= 1, "accuracy in [0,1]: " + completion.accuracy);
  assert.ok(completion.rewards, "rewards present: " + JSON.stringify(completion));
  assert.equal(completion.rewards_pending, false, "rewards should be ready (dev-sync events are synchronous)");
  assert.ok(completion.rewards.xp_awarded > 0, "XP awarded");
  assert.ok(completion.rewards.first_completion === true, "first completion: " + JSON.stringify(completion.rewards));
  assert.ok(completion.rewards.mastery_changes.length >= 1, "a mastery change per lesson skill");
  // Mastery follows first-try accuracy: the script guesses before it knows answers, so a change
  // can be 0. Every change must still add up, and some skill must move whenever a first try was right.
  for (const m of completion.rewards.mastery_changes) {
    assert.ok(m.skill_name && m.topic_name && m.delta >= 0 && m.after === m.before + m.delta, JSON.stringify(m));
  }
  if (completion.first_try_correct > 0) {
    assert.ok(completion.rewards.mastery_changes.some((m) => m.delta > 0), "a right first try moves mastery: " + JSON.stringify(completion.rewards.mastery_changes));
  }
});

await step("completing again is idempotent and does not double-award", async () => {
  const again = await call("POST", `/attempts/${attempt.attempt_id}/complete`, { token, expect: 200 });
  assert.equal(again.json.data.rewards.xp_awarded, completion.rewards.xp_awarded);
  const home = await call("GET", "/dashboard/home", { token, expect: 200 });
  assert.equal(home.json.data.progress.total_xp, completion.rewards.xp_awarded, "XP granted exactly once");
});

await step("mastery moved is explained by topic and skill, and topic readiness reflects it", async () => {
  const changes = completion.rewards.mastery_changes;
  assert.ok(changes.length >= 1);
  assert.ok(changes.every((c) => c.topic_name && c.skill_name && c.delta >= 0));
  const topics = (await call("GET", "/progress/topics", { token, expect: 200 })).json.data;
  const sd = topics.find((t) => t.slug === changes[0].topic_slug);
  assert.ok(sd.mastery !== null && sd.skills_started >= 1, "the topic is started: " + JSON.stringify(sd));
  assert.equal(sd.mastery, Math.round(sd.skills.filter((k) => k.mastery !== null).reduce((a, k) => a + k.mastery, 0) / sd.skills_started));
  const untouched = topics.filter((t) => t.slug !== sd.slug);
  assert.ok(untouched.every((t) => t.mastery === null), "other topics stay not started");
  const home = (await call("GET", "/dashboard/home", { token, expect: 200 })).json.data;
  const pathNow = (await call("GET", "/path", { token, expect: 200 })).json.data;
  const current = pathNow.worlds.flatMap((w) => w.nodes).find((n) => n.status === "current");
  if (current) {
    assert.equal(home.next_lesson?.lesson_id, current.lesson_id, "Continue opens the path's current node");
    assert.notEqual(current.lesson_id, node.lesson_id, "the finished lesson is no longer current");
  } else {
    assert.equal(home.next_lesson, null, "with nothing current there is nothing to continue");
  }
});

await step("streak counts the completed lesson", async () => {
  const home = await call("GET", "/dashboard/home", { token, expect: 200 });
  assert.equal(home.json.data.streak.current_streak, 1);
  assert.equal(home.json.data.streak.streak_active_today, true);
});

await step("path now shows the node done", async () => {
  const r = await call("GET", "/path", { token, expect: 200 });
  assert.equal(r.json.data.worlds[0].nodes[0].status, "done");
});

await step("rewards are private to the user who earned them", async () => {
  await call("GET", `/progress/attempts/${attempt.attempt_id}/rewards`, { token, expect: 200 });
  await call("GET", `/progress/attempts/${attempt.attempt_id}/rewards`, { token: other, expect: 404 });
});

await step("replaying a finished lesson grants nothing more", async () => {
  const r = await call("POST", `/lessons/${node.lesson_id}/attempts`, { token, expect: 201 });
  const replay = r.json.data;
  for (const s of replay.steps.filter((x) => x.type !== "intro")) {
    const path = `/attempts/${replay.attempt_id}/steps/${s.id}/answer`;
    let a = await call("POST", path, { token, body: { try: 1, answer: firstGuess(s) }, expect: 200 });
    if (!a.json.data.correct) {
      await call("POST", path, { token, body: { try: 2, answer: a.json.data.correct_answer }, expect: 200 });
    }
  }
  const done = await call("POST", `/attempts/${replay.attempt_id}/complete`, { token, expect: 200 });
  assert.equal(done.json.data.rewards.first_completion, false);
  assert.equal(done.json.data.rewards.xp_awarded, 0);
  assert.equal(done.json.data.rewards.mastery_changes.length, 0);
});

console.log(failures === 0 ? "\nall checks passed" : `\n${failures} check(s) failed`);
process.exit(failures === 0 ? 0 : 1);
