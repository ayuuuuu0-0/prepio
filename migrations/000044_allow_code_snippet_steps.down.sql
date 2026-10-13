ALTER TABLE lesson_steps DROP CONSTRAINT lesson_steps_type_check;
ALTER TABLE lesson_steps ADD CONSTRAINT lesson_steps_type_check
    CHECK (type IN ('intro', 'mcq', 'true_false', 'fill_blank', 'arrange', 'prose'));
