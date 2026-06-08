CREATE TABLE streak_freeze_inventory (
    user_id     UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    count       INT NOT NULL DEFAULT 0,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TRIGGER streak_freeze_inventory_set_updated_at
    BEFORE UPDATE ON streak_freeze_inventory
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
