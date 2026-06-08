CREATE TABLE gem_ledger (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount          INT NOT NULL,
    reason          TEXT NOT NULL,
    source_event_id UUID NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
