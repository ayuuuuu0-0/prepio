-- Topics are the user-facing grouping of skills; readiness is computed per topic.
CREATE TABLE topics (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug        TEXT NOT NULL UNIQUE,
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    sort_order  INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TRIGGER topics_set_updated_at
    BEFORE UPDATE ON topics
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

INSERT INTO topics (id, slug, name, description, sort_order) VALUES
    ('c4000000-0000-4000-8000-000000000001', 'system-design', 'System Design', 'Design scalable, reliable systems and reason about their trade-offs.', 1),
    ('c4000000-0000-4000-8000-000000000002', 'backend-production', 'Backend & Production', 'APIs, databases, concurrency, and running services in production.', 2),
    ('c4000000-0000-4000-8000-000000000003', 'low-level-design', 'Low-Level Design', 'Object-oriented design, patterns, and clean code structure.', 3),
    ('c4000000-0000-4000-8000-000000000004', 'dsa-refresher', 'DSA Refresher', 'Data structures and algorithms, refreshed for working engineers.', 4);

-- A skill category belongs to at most one topic. Categories without a topic do not contribute to readiness.
ALTER TABLE skill_categories
    ADD COLUMN topic_id UUID REFERENCES topics(id);

CREATE INDEX skill_categories_topic_id_idx ON skill_categories (topic_id);

-- New category for the Backend & Production topic.
INSERT INTO skill_categories (id, slug, name, sort_order) VALUES
    ('c1000000-0000-4000-8000-000000000009', 'backend-production', 'Backend & Production', 9);

INSERT INTO skills (id, category_id, slug, name, sort_order) VALUES
    ('b2000001-0000-4000-8000-000000000030', 'c1000000-0000-4000-8000-000000000009', 'backend-api-design', 'API Design', 1),
    ('b2000001-0000-4000-8000-000000000031', 'c1000000-0000-4000-8000-000000000009', 'backend-databases', 'Databases & Transactions', 2),
    ('b2000001-0000-4000-8000-000000000032', 'c1000000-0000-4000-8000-000000000009', 'backend-concurrency', 'Concurrency & Async', 3),
    ('b2000001-0000-4000-8000-000000000033', 'c1000000-0000-4000-8000-000000000009', 'backend-observability', 'Observability & Operations', 4);

UPDATE skill_categories SET topic_id = 'c4000000-0000-4000-8000-000000000001' WHERE slug = 'system-design';
UPDATE skill_categories SET topic_id = 'c4000000-0000-4000-8000-000000000002' WHERE slug = 'backend-production';
UPDATE skill_categories SET topic_id = 'c4000000-0000-4000-8000-000000000003' WHERE slug = 'low-level-design';
UPDATE skill_categories SET topic_id = 'c4000000-0000-4000-8000-000000000004' WHERE slug IN ('data-structures', 'algorithms');
