CREATE TABLE notification_log (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    notification_type   TEXT NOT NULL,
    channel             TEXT NOT NULL CHECK (channel IN ('fcm', 'email')),
    sent_at             TIMESTAMPTZ NOT NULL DEFAULT now()
);
