CREATE TABLE user_question_history (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    question_id     UUID NOT NULL REFERENCES questions(id),
    correct         BOOLEAN NOT NULL,
    submitted_at    TIMESTAMPTZ NOT NULL,
    received_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    session_id      UUID NOT NULL,
    UNIQUE (user_id, question_id, session_id)
);
