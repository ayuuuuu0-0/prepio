-- Backfill Foundation Forest content architecture with explicit skill/pool/node bindings.

-- Node slugs
UPDATE journey_nodes SET slug = 'arrays-basics' WHERE id = 'd0000000-0000-4000-8000-000000000001';
UPDATE journey_nodes SET slug = 'string-patterns' WHERE id = 'd0000000-0000-4000-8000-000000000002';
UPDATE journey_nodes SET slug = 'hash-maps' WHERE id = 'd0000000-0000-4000-8000-000000000003';
UPDATE journey_nodes SET slug = 'tree-traversal' WHERE id = 'd0000000-0000-4000-8000-000000000004';
UPDATE journey_nodes SET slug = 'forest-boss' WHERE id = 'd0000000-0000-4000-8000-000000000005';

-- Foundation Forest question pools
INSERT INTO question_pools (id, skill_id, slug, name, description, sort_order) VALUES
    ('e0000001-0000-4000-8000-000000000001', 'b2000001-0000-4000-8000-000000000002', 'foundation-arrays-beginner', 'Foundation Arrays', 'Beginner array questions for Foundation Forest', 1),
    ('e0000001-0000-4000-8000-000000000002', 'b2000001-0000-4000-8000-000000000003', 'foundation-strings-beginner', 'Foundation Strings', 'Beginner string questions for Foundation Forest', 1),
    ('e0000001-0000-4000-8000-000000000003', 'b2000001-0000-4000-8000-000000000004', 'foundation-hash-maps-beginner', 'Foundation Hash Maps', 'Beginner hash map questions for Foundation Forest', 1),
    ('e0000001-0000-4000-8000-000000000004', 'b2000001-0000-4000-8000-000000000008', 'foundation-trees-beginner', 'Foundation Trees', 'Beginner tree questions for Foundation Forest', 1);

-- Pool question assignments (approved seed questions)
INSERT INTO pool_questions (pool_id, question_id, sort_order) VALUES
    ('e0000001-0000-4000-8000-000000000001', 'b0000000-0000-4000-8000-000000000001', 1),
    ('e0000001-0000-4000-8000-000000000002', 'b0000000-0000-4000-8000-000000000006', 1),
    ('e0000001-0000-4000-8000-000000000003', 'b0000000-0000-4000-8000-000000000001', 1),
    ('e0000001-0000-4000-8000-000000000004', 'b0000000-0000-4000-8000-000000000004', 1);

-- Supplementary pools for skills with mapped questions (content graph completeness)
INSERT INTO question_pools (id, skill_id, slug, name, description, sort_order) VALUES
    ('e0000001-0000-4000-8000-000000000005', 'b2000001-0000-4000-8000-000000000005', 'foundation-linked-lists-beginner', 'Foundation Linked Lists', 'Linked list questions', 1),
    ('e0000001-0000-4000-8000-000000000006', 'b2000001-0000-4000-8000-000000000016', 'foundation-system-design-beginner', 'Foundation System Design', 'System design fundamentals', 1),
    ('e0000001-0000-4000-8000-000000000007', 'b2000001-0000-4000-8000-000000000020', 'foundation-lld-beginner', 'Foundation LLD', 'LLD fundamentals', 1),
    ('e0000001-0000-4000-8000-000000000008', 'b2000001-0000-4000-8000-000000000028', 'foundation-problem-solving-beginner', 'Foundation Problem Solving', 'Aptitude and reasoning', 1),
    ('e0000001-0000-4000-8000-000000000009', 'b2000001-0000-4000-8000-000000000001', 'foundation-programming-beginner', 'Foundation Programming', 'CS fundamentals', 1),
    ('e0000001-0000-4000-8000-000000000010', 'b2000001-0000-4000-8000-000000000023', 'foundation-behavioral-beginner', 'Foundation Behavioral', 'Behavioral interview questions', 1),
    ('e0000001-0000-4000-8000-000000000011', 'b2000001-0000-4000-8000-000000000017', 'foundation-system-design-scaling', 'Foundation System Design Scaling', 'Scaling and performance', 1),
    ('e0000001-0000-4000-8000-000000000012', 'b2000001-0000-4000-8000-000000000010', 'foundation-heaps-beginner', 'Foundation Heaps', 'Heap-based questions', 1);

INSERT INTO pool_questions (pool_id, question_id, sort_order) VALUES
    ('e0000001-0000-4000-8000-000000000005', 'b0000000-0000-4000-8000-000000000002', 1),
    ('e0000001-0000-4000-8000-000000000006', 'b0000000-0000-4000-8000-000000000003', 1),
    ('e0000001-0000-4000-8000-000000000007', 'b0000000-0000-4000-8000-000000000005', 1),
    ('e0000001-0000-4000-8000-000000000008', 'b0000000-0000-4000-8000-000000000007', 1),
    ('e0000001-0000-4000-8000-000000000009', 'b0000000-0000-4000-8000-000000000008', 1),
    ('e0000001-0000-4000-8000-000000000010', 'b0000000-0000-4000-8000-000000000009', 1),
    ('e0000001-0000-4000-8000-000000000011', 'b0000000-0000-4000-8000-000000000010', 1),
    ('e0000001-0000-4000-8000-000000000011', 'b0000000-0000-4000-8000-000000000012', 2),
    ('e0000001-0000-4000-8000-000000000012', 'b0000000-0000-4000-8000-000000000011', 1);

-- Node → skill bindings (Foundation Forest)
INSERT INTO node_skills (node_id, skill_id, is_primary) VALUES
    ('d0000000-0000-4000-8000-000000000001', 'b2000001-0000-4000-8000-000000000002', true),
    ('d0000000-0000-4000-8000-000000000002', 'b2000001-0000-4000-8000-000000000003', true),
    ('d0000000-0000-4000-8000-000000000003', 'b2000001-0000-4000-8000-000000000004', true),
    ('d0000000-0000-4000-8000-000000000004', 'b2000001-0000-4000-8000-000000000008', true),
    ('d0000000-0000-4000-8000-000000000005', 'b2000001-0000-4000-8000-000000000002', true),
    ('d0000000-0000-4000-8000-000000000005', 'b2000001-0000-4000-8000-000000000004', false),
    ('d0000000-0000-4000-8000-000000000005', 'b2000001-0000-4000-8000-000000000008', false);

-- Node → pool bindings
INSERT INTO node_pools (node_id, pool_id, selection_strategy, questions_required) VALUES
    ('d0000000-0000-4000-8000-000000000001', 'e0000001-0000-4000-8000-000000000001', 'random_unseen', 1),
    ('d0000000-0000-4000-8000-000000000002', 'e0000001-0000-4000-8000-000000000002', 'random_unseen', 1),
    ('d0000000-0000-4000-8000-000000000003', 'e0000001-0000-4000-8000-000000000003', 'random_unseen', 1),
    ('d0000000-0000-4000-8000-000000000004', 'e0000001-0000-4000-8000-000000000004', 'random_unseen', 1),
    ('d0000000-0000-4000-8000-000000000005', 'e0000001-0000-4000-8000-000000000001', 'boss_mixed', 1),
    ('d0000000-0000-4000-8000-000000000005', 'e0000001-0000-4000-8000-000000000003', 'boss_mixed', 1),
    ('d0000000-0000-4000-8000-000000000005', 'e0000001-0000-4000-8000-000000000004', 'boss_mixed', 1);
