CREATE TABLE character_unlocks (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    character_id    UUID NOT NULL REFERENCES characters(id),
    unlocked_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, character_id)
);
