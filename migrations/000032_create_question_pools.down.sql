DROP INDEX IF EXISTS journey_nodes_world_slug_idx;

ALTER TABLE journey_nodes
    DROP COLUMN IF EXISTS slug;

DROP TABLE IF EXISTS node_pools;
DROP TABLE IF EXISTS node_skills;
DROP TABLE IF EXISTS pool_questions;
DROP TABLE IF EXISTS question_pools;
