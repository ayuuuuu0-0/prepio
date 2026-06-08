CREATE TABLE daily_papers (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_id  UUID NOT NULL UNIQUE,
    paper_date  DATE NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, paper_date)
);

CREATE TRIGGER daily_papers_set_updated_at
    BEFORE UPDATE ON daily_papers
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE daily_paper_questions (
    daily_paper_id  UUID NOT NULL REFERENCES daily_papers(id) ON DELETE CASCADE,
    question_id     UUID NOT NULL REFERENCES questions(id),
    position        INT NOT NULL,
    PRIMARY KEY (daily_paper_id, question_id)
);
