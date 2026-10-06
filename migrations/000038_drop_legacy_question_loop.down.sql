-- Restores the legacy question-loop STRUCTURE only. The seeded questions, pools, and user
-- history are not recoverable from a rollback; re-run the historical seed migrations on a
-- fresh database if that content is needed.
CREATE TABLE questions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    body            TEXT NOT NULL,
    round_type      TEXT NOT NULL CHECK (round_type IN (
                        'dsa', 'system_design', 'lld',
                        'aptitude', 'fundamentals', 'behavioral'
                    )),
    difficulty      TEXT NOT NULL CHECK (difficulty IN ('easy', 'medium', 'hard')),
    answer_guide    TEXT NOT NULL,
    status          TEXT NOT NULL DEFAULT 'pending'
                        CHECK (status IN ('pending', 'approved', 'retired')),
    is_weekend      BOOLEAN NOT NULL DEFAULT false,
    source          TEXT NOT NULL CHECK (source IN ('manual', 'ai_generated', 'scraped')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    evaluation_type TEXT
        CHECK (evaluation_type IS NULL OR evaluation_type IN (
            'multiple_choice', 'coding', 'system_design', 'behavioral'
        )),
    explanation     TEXT,
    hints           JSONB NOT NULL DEFAULT '[]'::jsonb,
    solution        TEXT,
    readiness_weight NUMERIC(3, 2) NOT NULL DEFAULT 1.00
        CHECK (readiness_weight > 0 AND readiness_weight <= 2.0),
    estimated_time  INT NOT NULL DEFAULT 10
        CHECK (estimated_time > 0)
);

CREATE TRIGGER questions_set_updated_at
    BEFORE UPDATE ON questions
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE user_journey_progress (
    user_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    node_id       UUID NOT NULL REFERENCES journey_nodes(id) ON DELETE CASCADE,
    status        TEXT NOT NULL DEFAULT 'locked',
    completed_at  TIMESTAMPTZ,
    PRIMARY KEY (user_id, node_id)
);

CREATE TABLE question_skills (
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    skill_id    UUID NOT NULL REFERENCES skills(id),
    subskill_id UUID NOT NULL REFERENCES subskills(id),
    weight      NUMERIC(4, 3) NOT NULL DEFAULT 1.000
        CHECK (weight > 0 AND weight <= 1),
    PRIMARY KEY (question_id, skill_id, subskill_id)
);

CREATE INDEX question_skills_skill_id_idx ON question_skills (skill_id);
CREATE INDEX question_skills_subskill_id_idx ON question_skills (subskill_id);

CREATE TABLE question_pools (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    skill_id    UUID NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    slug        TEXT NOT NULL UNIQUE,
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    sort_order  INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX question_pools_skill_id_idx ON question_pools (skill_id);

CREATE TRIGGER question_pools_set_updated_at
    BEFORE UPDATE ON question_pools
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE node_skills (
    node_id    UUID NOT NULL REFERENCES journey_nodes(id) ON DELETE CASCADE,
    skill_id   UUID NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    is_primary BOOLEAN NOT NULL DEFAULT true,
    PRIMARY KEY (node_id, skill_id)
);

CREATE INDEX node_skills_skill_id_idx ON node_skills (skill_id);

CREATE TABLE node_pools (
    node_id             UUID NOT NULL REFERENCES journey_nodes(id) ON DELETE CASCADE,
    pool_id             UUID NOT NULL REFERENCES question_pools(id) ON DELETE CASCADE,
    selection_strategy  TEXT NOT NULL DEFAULT 'random_unseen',
    questions_required  INT NOT NULL DEFAULT 1,
    PRIMARY KEY (node_id, pool_id)
);

CREATE INDEX node_pools_pool_id_idx ON node_pools (pool_id);

CREATE TABLE pool_questions (
    pool_id     UUID NOT NULL REFERENCES question_pools(id) ON DELETE CASCADE,
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    sort_order  INT NOT NULL DEFAULT 0,
    PRIMARY KEY (pool_id, question_id)
);

CREATE INDEX pool_questions_question_id_idx ON pool_questions (question_id);

CREATE TABLE user_question_history (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    question_id     UUID NOT NULL REFERENCES questions(id),
    correct         BOOLEAN NOT NULL,
    submitted_at    TIMESTAMPTZ NOT NULL,
    received_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    session_id      UUID NOT NULL,
    score           INT NOT NULL DEFAULT 0,
    UNIQUE (user_id, question_id, session_id)
);

CREATE TABLE daily_papers (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_id  UUID NOT NULL UNIQUE,
    paper_date  DATE NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, paper_date)
);

CREATE TRIGGER daily_papers_set_updated_at
    BEFORE UPDATE ON daily_papers
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE daily_paper_questions (
    daily_paper_id  UUID NOT NULL REFERENCES daily_papers(id) ON DELETE CASCADE,
    question_id     UUID NOT NULL REFERENCES questions(id),
    position        INT NOT NULL,
    PRIMARY KEY (daily_paper_id, question_id)
);
