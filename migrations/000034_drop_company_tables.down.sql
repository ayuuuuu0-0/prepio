CREATE TABLE IF NOT EXISTS question_tags (
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    company VARCHAR(50) NOT NULL,
    PRIMARY KEY (question_id, company)
);
CREATE INDEX IF NOT EXISTS idx_question_tags_company ON question_tags(company);

CREATE TABLE IF NOT EXISTS user_targets (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, company)
);
CREATE INDEX IF NOT EXISTS user_targets_user_id_idx ON user_targets (user_id);

CREATE TABLE IF NOT EXISTS company_skill_weights (
    company VARCHAR(50) NOT NULL,
    skill_id UUID NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    weight INT NOT NULL CHECK (weight >= 0 AND weight <= 100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (company, skill_id)
);
CREATE INDEX IF NOT EXISTS company_skill_weights_company_idx ON company_skill_weights (company);
