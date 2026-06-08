CREATE TABLE user_progress (
    user_id         UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    total_xp        INT NOT NULL DEFAULT 0,
    current_level   INT NOT NULL DEFAULT 1,
    gem_balance     INT NOT NULL DEFAULT 0,
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TRIGGER user_progress_set_updated_at
    BEFORE UPDATE ON user_progress
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
