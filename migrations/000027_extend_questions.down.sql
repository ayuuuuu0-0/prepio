ALTER TABLE questions
    DROP COLUMN IF EXISTS estimated_time,
    DROP COLUMN IF EXISTS readiness_weight,
    DROP COLUMN IF EXISTS solution,
    DROP COLUMN IF EXISTS hints,
    DROP COLUMN IF EXISTS explanation,
    DROP COLUMN IF EXISTS evaluation_type;
