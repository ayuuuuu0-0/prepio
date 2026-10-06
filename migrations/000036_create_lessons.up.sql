-- Lesson platform: worlds/nodes become content-as-code, lessons are server-graded.

-- Worlds and nodes are synced from content/ and are deprecated, never deleted.
ALTER TABLE worlds
    ADD COLUMN status TEXT NOT NULL DEFAULT 'published'
        CHECK (status IN ('published', 'deprecated'));

ALTER TABLE journey_nodes
    ADD COLUMN status TEXT NOT NULL DEFAULT 'published'
        CHECK (status IN ('published', 'deprecated'));

-- The legacy question-driven world has no lessons and is superseded by authored worlds.
UPDATE worlds SET status = 'deprecated' WHERE slug = 'foundation-forest';
UPDATE journey_nodes SET status = 'deprecated'
WHERE world_id = 'c0000000-0000-4000-8000-000000000001';

-- Node slugs are the stable key lessons use to bind to a node.
CREATE UNIQUE INDEX journey_nodes_slug_idx ON journey_nodes (slug) WHERE slug IS NOT NULL;

-- Configurable unlock rules: a node unlocks once all of its prerequisites are done.
CREATE TABLE node_prerequisites (
    node_id          UUID NOT NULL REFERENCES journey_nodes(id) ON DELETE CASCADE,
    requires_node_id UUID NOT NULL REFERENCES journey_nodes(id) ON DELETE CASCADE,
    PRIMARY KEY (node_id, requires_node_id),
    CHECK (node_id <> requires_node_id)
);

CREATE INDEX node_prerequisites_requires_idx ON node_prerequisites (requires_node_id);

CREATE TABLE lessons (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug          TEXT NOT NULL UNIQUE,
    node_id       UUID NOT NULL UNIQUE REFERENCES journey_nodes(id),
    title         TEXT NOT NULL,
    summary       JSONB NOT NULL DEFAULT '[]'::jsonb,
    kind          TEXT NOT NULL DEFAULT 'lesson' CHECK (kind IN ('lesson', 'boss')),
    difficulty    TEXT NOT NULL CHECK (difficulty IN ('easy', 'medium', 'hard')),
    est_minutes   INT NOT NULL CHECK (est_minutes > 0),
    status        TEXT NOT NULL DEFAULT 'draft'
        CHECK (status IN ('draft', 'review', 'published', 'deprecated')),
    version       INT NOT NULL DEFAULT 1 CHECK (version >= 1),
    content_hash  TEXT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TRIGGER lessons_set_updated_at
    BEFORE UPDATE ON lessons
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Steps and skill weights are kept per lesson version so in-flight attempts
-- keep grading against the content they started with.
CREATE TABLE lesson_steps (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    lesson_id      UUID NOT NULL REFERENCES lessons(id),
    lesson_version INT NOT NULL,
    step_key       TEXT NOT NULL,
    position       INT NOT NULL CHECK (position >= 0),
    type           TEXT NOT NULL
        CHECK (type IN ('intro', 'mcq', 'true_false', 'fill_blank', 'arrange', 'prose')),
    payload        JSONB NOT NULL,
    UNIQUE (lesson_id, lesson_version, step_key),
    UNIQUE (lesson_id, lesson_version, position)
);

CREATE TABLE lesson_skills (
    lesson_id      UUID NOT NULL REFERENCES lessons(id),
    lesson_version INT NOT NULL,
    skill_id       UUID NOT NULL REFERENCES skills(id),
    weight         NUMERIC(4, 3) NOT NULL CHECK (weight > 0 AND weight <= 1),
    PRIMARY KEY (lesson_id, lesson_version, skill_id)
);

CREATE INDEX lesson_skills_skill_id_idx ON lesson_skills (skill_id);

CREATE TABLE lesson_attempts (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    lesson_id           UUID NOT NULL REFERENCES lessons(id),
    lesson_version      INT NOT NULL,
    status              TEXT NOT NULL DEFAULT 'in_progress'
        CHECK (status IN ('in_progress', 'completed')),
    graded_steps        INT,
    first_try_correct   INT,
    total_tries         INT,
    completion_event_id UUID,
    event_published_at  TIMESTAMPTZ,
    started_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at        TIMESTAMPTZ
);

-- One in-progress attempt per user and lesson; it is resumable.
CREATE UNIQUE INDEX lesson_attempts_one_active_idx
    ON lesson_attempts (user_id, lesson_id) WHERE status = 'in_progress';
CREATE INDEX lesson_attempts_user_lesson_idx ON lesson_attempts (user_id, lesson_id, status);

-- Step answers are idempotent per (attempt, step, try).
CREATE TABLE lesson_step_results (
    attempt_id UUID NOT NULL REFERENCES lesson_attempts(id) ON DELETE CASCADE,
    step_id    UUID NOT NULL REFERENCES lesson_steps(id),
    try_no     INT NOT NULL CHECK (try_no >= 1),
    answer     JSONB NOT NULL,
    correct    BOOLEAN NOT NULL,
    feedback   JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (attempt_id, step_id, try_no)
);
