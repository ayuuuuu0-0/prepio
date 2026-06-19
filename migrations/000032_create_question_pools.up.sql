-- Question pools and journey content bindings (A2)

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

CREATE TABLE pool_questions (
    pool_id     UUID NOT NULL REFERENCES question_pools(id) ON DELETE CASCADE,
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    sort_order  INT NOT NULL DEFAULT 0,
    PRIMARY KEY (pool_id, question_id)
);

CREATE INDEX pool_questions_question_id_idx ON pool_questions (question_id);

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

-- Node slugs for stable content binding (additive, nullable for backward compatibility)
ALTER TABLE journey_nodes
    ADD COLUMN slug TEXT;

CREATE UNIQUE INDEX journey_nodes_world_slug_idx ON journey_nodes (world_id, slug)
    WHERE slug IS NOT NULL;
