-- Content domain: a world can name the topic it builds. Focus topics put their worlds
-- first on the path (ordering only; unlock rules are unchanged, so nothing is locked).
ALTER TABLE worlds ADD COLUMN topic_id UUID REFERENCES topics(id);
