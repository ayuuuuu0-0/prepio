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
    let r = await call("POST", path, { token, body: { try: 1, answer: { choice: 0 } }, expect: 200 });
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
  assert.ok(completion.accuracy > 0 && completion.accuracy <= 1);
  assert.equal(completion.rewards_pending, false, "rewards should be ready (dev-sync events are synchronous)");
  assert.ok(completion.rewards.xp_awarded > 0, "XP awarded");
  assert.ok(completion.rewards.first_completion === true);
  assert.ok(completion.rewards.mastery_changes.length >= 1);
  const m = completion.rewards.mastery_changes[0];
  assert.ok(m.skill_name && m.topic_name && m.delta > 0 && m.after === m.before + m.delta);
});

await step("completing again is idempotent and does not double-award", async () => {
  const again = await call("POST", `/attempts/${attempt.attempt_id}/complete`, { token, expect: 200 });
  assert.equal(again.json.data.rewards.xp_awarded, completion.rewards.xp_awarded);
  const home = await call("GET", "/dashboard/home", { token, expect: 200 });
  assert.equal(home.json.data.progress.total_xp, completion.rewards.xp_awarded, "XP granted exactly once");
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
    let a = await call("POST", path, { token, body: { try: 1, answer: { choice: 0 } }, expect: 200 });
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
