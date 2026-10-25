DROP INDEX IF EXISTS character_dialogues_lookup_idx;

DELETE FROM character_dialogues 
WHERE notification_type IN ('correct_answer', 'wrong_answer', 'lesson_complete_high', 'lesson_complete_low', 'streak_extended');

ALTER TABLE character_dialogues DROP COLUMN IF EXISTS expression;
