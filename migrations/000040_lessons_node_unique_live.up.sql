-- A node binds exactly one live lesson. The original column-level UNIQUE also counted
-- deprecated lessons, so content-sync could never replace the lesson on a node (or
-- move a lesson onto a node a deprecated lesson still held).
--
-- The replacement only counts lessons that are not deprecated. It is an exclusion
-- constraint (equivalent to a partial unique index on node_id) so it can be deferred
-- to commit: content-sync rebinds every lesson in one transaction, and swapping the
-- lessons of two nodes is only valid once both rows have moved.
ALTER TABLE lessons DROP CONSTRAINT lessons_node_id_key;

ALTER TABLE lessons
    ADD CONSTRAINT lessons_node_id_live_excl
    EXCLUDE USING btree (node_id WITH =) WHERE (status <> 'deprecated')
    DEFERRABLE INITIALLY DEFERRED;
