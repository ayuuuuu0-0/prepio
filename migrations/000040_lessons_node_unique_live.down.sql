-- Restores one lesson per node, counting deprecated lessons too. This fails if a
-- node already holds a deprecated lesson and its replacement; resolve those rows
-- (they hold attempts, so they cannot simply be deleted) before rolling back.
ALTER TABLE lessons DROP CONSTRAINT lessons_node_id_live_excl;

ALTER TABLE lessons ADD CONSTRAINT lessons_node_id_key UNIQUE (node_id);
