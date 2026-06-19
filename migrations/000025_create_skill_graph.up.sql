CREATE TABLE skill_categories (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug        TEXT NOT NULL UNIQUE,
    name        TEXT NOT NULL,
    sort_order  INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TRIGGER skill_categories_set_updated_at
    BEFORE UPDATE ON skill_categories
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE skills (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id UUID NOT NULL REFERENCES skill_categories(id),
    slug        TEXT NOT NULL UNIQUE,
    name        TEXT NOT NULL,
    description TEXT,
    sort_order  INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX skills_category_id_idx ON skills (category_id);

CREATE TRIGGER skills_set_updated_at
    BEFORE UPDATE ON skills
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE subskills (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    skill_id    UUID NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    slug        TEXT NOT NULL,
    name        TEXT NOT NULL,
    sort_order  INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (skill_id, slug)
);

CREATE INDEX subskills_skill_id_idx ON subskills (skill_id);

CREATE TRIGGER subskills_set_updated_at
    BEFORE UPDATE ON subskills
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

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
