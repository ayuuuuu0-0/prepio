CREATE TABLE worlds (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug        TEXT NOT NULL UNIQUE,
    name        TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    theme       TEXT NOT NULL DEFAULT 'forest',
    sort_order  INT NOT NULL DEFAULT 0
);

CREATE TABLE journey_nodes (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    world_id    UUID NOT NULL REFERENCES worlds(id) ON DELETE CASCADE,
    label       TEXT NOT NULL,
    node_type   TEXT NOT NULL DEFAULT 'lesson',
    sort_order  INT NOT NULL DEFAULT 0
);

CREATE TABLE user_journey_progress (
    user_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    node_id       UUID NOT NULL REFERENCES journey_nodes(id) ON DELETE CASCADE,
    status        TEXT NOT NULL DEFAULT 'locked',
    completed_at  TIMESTAMPTZ,
    PRIMARY KEY (user_id, node_id)
);

INSERT INTO worlds (id, slug, name, description, theme, sort_order)
VALUES (
    'c0000000-0000-4000-8000-000000000001',
    'foundation-forest',
    'Foundation Forest',
    'World 1 · Beginner Journey',
    'forest',
    1
);

INSERT INTO journey_nodes (id, world_id, label, node_type, sort_order)
VALUES
    ('d0000000-0000-4000-8000-000000000001', 'c0000000-0000-4000-8000-000000000001', 'Arrays Basics', 'lesson', 1),
    ('d0000000-0000-4000-8000-000000000002', 'c0000000-0000-4000-8000-000000000001', 'String Patterns', 'lesson', 2),
    ('d0000000-0000-4000-8000-000000000003', 'c0000000-0000-4000-8000-000000000001', 'Hash Maps', 'lesson', 3),
    ('d0000000-0000-4000-8000-000000000004', 'c0000000-0000-4000-8000-000000000001', 'Tree Traversal', 'lesson', 4),
    ('d0000000-0000-4000-8000-000000000005', 'c0000000-0000-4000-8000-000000000001', 'Forest Boss', 'boss', 5);
