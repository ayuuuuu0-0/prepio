-- Progress domain: every mastery change is recorded so readiness is explainable.
CREATE TABLE mastery_ledger (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id        UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    skill_id       UUID NOT NULL REFERENCES skills(id),
    attempt_id     UUID NOT NULL,
    lesson_id      UUID NOT NULL,
    delta          INT NOT NULL,
    mastery_before INT NOT NULL CHECK (mastery_before >= 0 AND mastery_before <= 100),
    mastery_after  INT NOT NULL CHECK (mastery_after >= 0 AND mastery_after <= 100),
    accuracy       NUMERIC(4, 3) NOT NULL CHECK (accuracy >= 0 AND accuracy <= 1),
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (attempt_id, skill_id)
);

CREATE INDEX mastery_ledger_user_created_idx ON mastery_ledger (user_id, created_at DESC);
CREATE INDEX mastery_ledger_user_skill_idx ON mastery_ledger (user_id, skill_id);

-- One row per completed attempt: makes lesson.completed handling idempotent
-- and lets the completion screen look up what the attempt earned.
CREATE TABLE lesson_rewards (
    attempt_id       UUID PRIMARY KEY,
    user_id          UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    lesson_id        UUID NOT NULL,
    first_completion BOOLEAN NOT NULL,
    xp_awarded       INT NOT NULL DEFAULT 0,
    gems_awarded     INT NOT NULL DEFAULT 0,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX lesson_rewards_user_lesson_idx ON lesson_rewards (user_id, lesson_id);
