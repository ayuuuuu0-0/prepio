DROP TABLE IF EXISTS user_targets;

ALTER TABLE users
    DROP COLUMN IF EXISTS onboarding_completed,
    DROP COLUMN IF EXISTS experience_level;
