-- User domain: the 1-3 topics a learner wants to sharpen. Focus topics reorder and highlight
-- the path; they never lock content.
CREATE TABLE user_focus_topics (
    user_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    topic_id  UUID NOT NULL REFERENCES topics(id),
    position  INT NOT NULL CHECK (position BETWEEN 1 AND 3),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, topic_id),
    UNIQUE (user_id, position)
);
