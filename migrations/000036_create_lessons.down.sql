DROP TABLE IF EXISTS lesson_step_results;
DROP TABLE IF EXISTS lesson_attempts;
DROP TABLE IF EXISTS lesson_skills;
DROP TABLE IF EXISTS lesson_steps;
DROP TABLE IF EXISTS lessons;
DROP TABLE IF EXISTS node_prerequisites;
DROP INDEX IF EXISTS journey_nodes_slug_idx;
ALTER TABLE journey_nodes DROP COLUMN IF EXISTS status;
ALTER TABLE worlds DROP COLUMN IF EXISTS status;
