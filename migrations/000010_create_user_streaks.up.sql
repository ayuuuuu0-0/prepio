CREATE TABLE user_streaks (
    user_id             UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    current_streak      INT NOT NULL DEFAULT 0,
    longest_streak      INT NOT NULL DEFAULT 0,
    last_activity_date  DATE,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TRIGGER user_streaks_set_updated_at
    BEFORE UPDATE ON user_streaks
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
