CREATE TABLE question_tags (
    question_id     UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    company         TEXT NOT NULL,
    PRIMARY KEY (question_id, company)
);
