-- Add expression column to character_dialogues
ALTER TABLE character_dialogues ADD COLUMN IF NOT EXISTS expression TEXT NOT NULL DEFAULT 'neutral'
    CHECK (expression IN ('happy', 'disappointed', 'excited', 'thinking', 'neutral'));

-- Index character dialogues by character and event type for fast random selection
CREATE INDEX IF NOT EXISTS character_dialogues_lookup_idx ON character_dialogues (character_id, notification_type);

-- Seed dialogue lines for all 5 companions
-- Byte (Capybara - Calm, Zen mentor)
INSERT INTO character_dialogues (character_id, notification_type, expression, dialogue_line) VALUES
    ('b0000000-0000-4000-8000-000000000001', 'correct_answer', 'happy', 'Clean architecture. Byte sips water in calm approval.'),
    ('b0000000-0000-4000-8000-000000000001', 'correct_answer', 'happy', 'Zero stress, maximum correctness. That''s the capybara way.'),
    ('b0000000-0000-4000-8000-000000000001', 'correct_answer', 'happy', 'Look at that. You didn''t overcomplicate it. Byte is proud.'),
    ('b0000000-0000-4000-8000-000000000001', 'correct_answer', 'happy', 'Solid engineering. Even production wouldn''t panic over that.'),
    ('b0000000-0000-4000-8000-000000000001', 'wrong_answer', 'disappointed', 'Byte didn''t see that. Take a breath and let''s pretend that didn''t happen.'),
    ('b0000000-0000-4000-8000-000000000001', 'wrong_answer', 'disappointed', 'A minor distributed system incident. Nothing a coffee and another try can''t fix.'),
    ('b0000000-0000-4000-8000-000000000001', 'wrong_answer', 'disappointed', 'That answer caused a 500 error in my soul. But we iterate.'),
    ('b0000000-0000-4000-8000-000000000001', 'wrong_answer', 'disappointed', 'Bold choice! Statistically incorrect, but bold.'),
    ('b0000000-0000-4000-8000-000000000001', 'lesson_complete_high', 'excited', 'Flawless run. Byte awards you the honorary chill senior engineer badge.'),
    ('b0000000-0000-4000-8000-000000000001', 'lesson_complete_high', 'excited', 'Not a single microservice melted today. Exceptional work.'),
    ('b0000000-0000-4000-8000-000000000001', 'lesson_complete_low', 'happy', 'Lesson wrapped! Scars build senior engineers. Byte is still proud.'),
    ('b0000000-0000-4000-8000-000000000001', 'lesson_complete_low', 'happy', 'You fought production and survived. Tomorrow will be even smoother.'),
    ('b0000000-0000-4000-8000-000000000001', 'streak_extended', 'excited', 'Day streak extended. Byte approves of your consistency.');

-- Pip (Red Panda - Hyper-enthusiastic, caffeine-fueled)
INSERT INTO character_dialogues (character_id, notification_type, expression, dialogue_line) VALUES
    ('b0000000-0000-4000-8000-000000000002', 'correct_answer', 'happy', 'BOOM! Nailed it! Pip is literally doing backflips right now!'),
    ('b0000000-0000-4000-8000-000000000002', 'correct_answer', 'happy', 'YES! Ship it straight to prod! Who needs code review?!'),
    ('b0000000-0000-4000-8000-000000000002', 'correct_answer', 'happy', 'You''re a genius! Pip''s tail is wagging at maximum RPM!'),
    ('b0000000-0000-4000-8000-000000000002', 'correct_answer', 'happy', 'Look at those brain gains! 10x developer energy!'),
    ('b0000000-0000-4000-8000-000000000002', 'wrong_answer', 'disappointed', 'AAAH! The build is red! Pip is panicking! Quick, try again!'),
    ('b0000000-0000-4000-8000-000000000002', 'wrong_answer', 'disappointed', 'Oof, did the linter fail or did your cat walk on the keyboard?!'),
    ('b0000000-0000-4000-8000-000000000002', 'wrong_answer', 'disappointed', 'That was... certainly one of the answers of all time.'),
    ('b0000000-0000-4000-8000-000000000002', 'wrong_answer', 'disappointed', 'Pip winces dramatically! But hey, failure is just learning in disguise!'),
    ('b0000000-0000-4000-8000-000000000002', 'lesson_complete_high', 'excited', 'WE DID IT! 100% pure genius! Coffee is on Pip today!'),
    ('b0000000-0000-4000-8000-000000000002', 'lesson_complete_high', 'excited', 'That was incredible! Pip is never doubting you again!'),
    ('b0000000-0000-4000-8000-000000000002', 'lesson_complete_low', 'happy', 'PHEW! We made it through! That was wild, let''s crush the next one!'),
    ('b0000000-0000-4000-8000-000000000002', 'lesson_complete_low', 'happy', 'Finished! Pip''s heart rate was 180 bpm the whole time, but you did it!'),
    ('b0000000-0000-4000-8000-000000000002', 'streak_extended', 'excited', 'STREAK ALIVE! Pip does a victory dance across the desk!');

-- Nova (Pangolin - Dry, paranoid security architect)
INSERT INTO character_dialogues (character_id, notification_type, expression, dialogue_line) VALUES
    ('b0000000-0000-4000-8000-000000000003', 'correct_answer', 'happy', 'Acceptable. My threat model detects zero vulnerabilities in that answer.'),
    ('b0000000-0000-4000-8000-000000000003', 'correct_answer', 'happy', 'Accurate. Even my paranoia couldn''t find a flaw there.'),
    ('b0000000-0000-4000-8000-000000000003', 'correct_answer', 'happy', 'Correct. Nova records this in your permanent file as a rare success.'),
    ('b0000000-0000-4000-8000-000000000003', 'correct_answer', 'happy', 'Optimal complexity achieved. Nova''s armor plates relax by 2%.'),
    ('b0000000-0000-4000-8000-000000000003', 'wrong_answer', 'disappointed', 'Critical vulnerability detected in your logic. Nova is rolling into a protective ball.'),
    ('b0000000-0000-4000-8000-000000000003', 'wrong_answer', 'disappointed', 'That answer would have failed penetration testing in 3 milliseconds.'),
    ('b0000000-0000-4000-8000-000000000003', 'wrong_answer', 'disappointed', 'Did you read the docs or are we doing stochastic guessing today?'),
    ('b0000000-0000-4000-8000-000000000003', 'wrong_answer', 'disappointed', 'Incorrect. If that was production, paging duty would like a word.'),
    ('b0000000-0000-4000-8000-000000000003', 'lesson_complete_high', 'excited', 'Audit complete. Zero security incidents reported. High praise indeed.'),
    ('b0000000-0000-4000-8000-000000000003', 'lesson_complete_high', 'excited', 'Remarkable accuracy. Nova certifies your system as production-grade.'),
    ('b0000000-0000-4000-8000-000000000003', 'lesson_complete_low', 'happy', 'Lesson concluded with several near-catastrophic CVEs. Review is advised.'),
    ('b0000000-0000-4000-8000-000000000003', 'lesson_complete_low', 'happy', 'You survived the security review. Barely. Nova is watching.'),
    ('b0000000-0000-4000-8000-000000000003', 'streak_extended', 'excited', 'Streak intact. Consistent authentication verified.');

-- Kodo (Axolotl - Wholesome, regenerating, warm)
INSERT INTO character_dialogues (character_id, notification_type, expression, dialogue_line) VALUES
    ('b0000000-0000-4000-8000-000000000004', 'correct_answer', 'happy', 'Yay! Kodo''s gills are glowing! You did so wonderfully!'),
    ('b0000000-0000-4000-8000-000000000004', 'correct_answer', 'happy', 'Brilliant! Kodo is blowing happy bubbles for you! 🫧'),
    ('b0000000-0000-4000-8000-000000000004', 'correct_answer', 'happy', 'You make this look so easy! Keep shining!'),
    ('b0000000-0000-4000-8000-000000000004', 'correct_answer', 'happy', 'Kodo knew you could do it! Spot on!'),
    ('b0000000-0000-4000-8000-000000000004', 'wrong_answer', 'disappointed', 'Aww, don''t worry! Like my gills, your confidence will regenerate right away! 🫧'),
    ('b0000000-0000-4000-8000-000000000004', 'wrong_answer', 'disappointed', 'Mistakes are just stepping stones! Kodo still thinks you''re great!'),
    ('b0000000-0000-4000-8000-000000000004', 'wrong_answer', 'disappointed', 'Oopsie! That''s okay, even senior staff engineers trip on this one!'),
    ('b0000000-0000-4000-8000-000000000004', 'wrong_answer', 'disappointed', 'Big hugs from Kodo! Shake it off and let''s get the next one!'),
    ('b0000000-0000-4000-8000-000000000004', 'lesson_complete_high', 'excited', 'You were absolutely magnificent today! Kodo is floating with joy!'),
    ('b0000000-0000-4000-8000-000000000004', 'lesson_complete_high', 'excited', 'Full lesson complete! You''re officially Kodo''s favorite engineer!'),
    ('b0000000-0000-4000-8000-000000000004', 'lesson_complete_low', 'happy', 'You finished the whole lesson! Every little step makes you stronger!'),
    ('b0000000-0000-4000-8000-000000000004', 'lesson_complete_low', 'happy', 'Kodo is so proud of you for sticking it out to the end!'),
    ('b0000000-0000-4000-8000-000000000004', 'streak_extended', 'excited', 'Another day, another streak kept! Kodo is cheering for you!');

-- Zara (Snow Leopard - High-standards Staff architect, sharp wit)
INSERT INTO character_dialogues (character_id, notification_type, expression, dialogue_line) VALUES
    ('b0000000-0000-4000-8000-000000000005', 'correct_answer', 'happy', 'Naturally. I wouldn''t expect anything less from someone I''m mentoring.'),
    ('b0000000-0000-4000-8000-000000000005', 'correct_answer', 'happy', 'Clean. Elegant. Zara nods approvingly from high altitude.'),
    ('b0000000-0000-4000-8000-000000000005', 'correct_answer', 'happy', 'That''s how real systems are designed. Commendable.'),
    ('b0000000-0000-4000-8000-000000000005', 'correct_answer', 'happy', 'Spot on. You actually read between the lines.'),
    ('b0000000-0000-4000-8000-000000000005', 'wrong_answer', 'disappointed', 'Zara is going to look away and pretend that answer never happened.'),
    ('b0000000-0000-4000-8000-000000000005', 'wrong_answer', 'disappointed', 'Did you consult a Ouija board for that architecture decision?'),
    ('b0000000-0000-4000-8000-000000000005', 'wrong_answer', 'disappointed', 'Wrong. I''ve seen junior interns propose better designs, but do go on.'),
    ('b0000000-0000-4000-8000-000000000005', 'wrong_answer', 'disappointed', 'That answer just made Zara''s whiskers twitch in mild horror.'),
    ('b0000000-0000-4000-8000-000000000005', 'lesson_complete_high', 'excited', 'Hmph. I suppose you actually possess genuine talent. Zara approves.'),
    ('b0000000-0000-4000-8000-000000000005', 'lesson_complete_high', 'excited', 'A masterclass. You''ve earned your place at the architectural table.'),
    ('b0000000-0000-4000-8000-000000000005', 'lesson_complete_low', 'happy', 'A bit messy, darling, but you got through it. Refinement comes with practice.'),
    ('b0000000-0000-4000-8000-000000000005', 'lesson_complete_low', 'happy', 'Completed. We''ll work on your instincts tomorrow.'),
    ('b0000000-0000-4000-8000-000000000005', 'streak_extended', 'excited', 'Streak preserved. Discipline is what separates professionals from amateurs.');
