ALTER TABLE users
    ADD COLUMN experience_level TEXT,
    ADD COLUMN onboarding_completed BOOLEAN NOT NULL DEFAULT false;

CREATE TABLE user_targets (
    user_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company   TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, company)
);

CREATE INDEX user_targets_user_id_idx ON user_targets (user_id);
