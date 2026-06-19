CREATE TABLE user_skill_scores (
    user_id           UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    skill_id          UUID NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    mastery           INT NOT NULL DEFAULT 0 CHECK (mastery >= 0 AND mastery <= 100),
    attempts          INT NOT NULL DEFAULT 0,
    last_practiced_at TIMESTAMPTZ,
    source            TEXT NOT NULL DEFAULT 'live'
        CHECK (source IN ('live', 'backfill')),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, skill_id)
);

CREATE INDEX user_skill_scores_user_id_idx ON user_skill_scores (user_id);
CREATE INDEX user_skill_scores_skill_id_idx ON user_skill_scores (skill_id);

CREATE TRIGGER user_skill_scores_set_updated_at
    BEFORE UPDATE ON user_skill_scores
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE company_skill_weights (
    company  TEXT NOT NULL,
    skill_id UUID NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    weight   INT NOT NULL CHECK (weight >= 0 AND weight <= 100),
    PRIMARY KEY (company, skill_id)
);

CREATE INDEX company_skill_weights_company_idx ON company_skill_weights (company);
