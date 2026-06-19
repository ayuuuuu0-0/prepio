ALTER TABLE questions
    ADD COLUMN evaluation_type TEXT
        CHECK (evaluation_type IS NULL OR evaluation_type IN (
            'multiple_choice', 'coding', 'system_design', 'behavioral'
        )),
    ADD COLUMN explanation TEXT,
    ADD COLUMN hints JSONB NOT NULL DEFAULT '[]'::jsonb,
    ADD COLUMN solution TEXT,
    ADD COLUMN readiness_weight NUMERIC(3, 2) NOT NULL DEFAULT 1.00
        CHECK (readiness_weight > 0 AND readiness_weight <= 2.0),
    ADD COLUMN estimated_time INT NOT NULL DEFAULT 10
        CHECK (estimated_time > 0);

-- Map legacy round_type values to evaluation_type for existing rows.
UPDATE questions SET evaluation_type = 'coding'
WHERE round_type IN ('dsa', 'lld') AND evaluation_type IS NULL;

UPDATE questions SET evaluation_type = 'system_design'
WHERE round_type = 'system_design' AND evaluation_type IS NULL;

UPDATE questions SET evaluation_type = 'behavioral'
WHERE round_type = 'behavioral' AND evaluation_type IS NULL;

UPDATE questions SET evaluation_type = 'multiple_choice'
WHERE round_type IN ('aptitude', 'fundamentals') AND evaluation_type IS NULL;

-- Default readiness_weight and estimated_time by difficulty.
UPDATE questions SET readiness_weight = 0.80, estimated_time = 8
WHERE difficulty = 'easy';

UPDATE questions SET readiness_weight = 1.00, estimated_time = 15
WHERE difficulty = 'medium';

UPDATE questions SET readiness_weight = 1.20, estimated_time = 25
WHERE difficulty = 'hard';
