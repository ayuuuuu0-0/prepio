-- Progress domain: achievements a learner has earned. The catalog lives in config
-- (config/achievements.go); only earned slugs are stored. attempt_id records the lesson
-- attempt that earned it (null for streak achievements), so a lesson's rewards can list it.
CREATE TABLE user_achievements (
    user_id          UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    achievement_slug TEXT NOT NULL,
    unlocked_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    attempt_id       UUID,
    PRIMARY KEY (user_id, achievement_slug)
);

CREATE INDEX user_achievements_attempt_idx ON user_achievements (attempt_id) WHERE attempt_id IS NOT NULL;
