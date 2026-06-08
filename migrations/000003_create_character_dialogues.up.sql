CREATE TABLE character_dialogues (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    character_id        UUID NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    notification_type   TEXT NOT NULL,
    dialogue_line       TEXT NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TRIGGER character_dialogues_set_updated_at
    BEFORE UPDATE ON character_dialogues
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
