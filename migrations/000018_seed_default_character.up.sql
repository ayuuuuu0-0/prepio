INSERT INTO characters (id, name, species, gem_cost, is_default)
VALUES (
    'a0000000-0000-4000-8000-000000000001',
    'Prep',
    'owl',
    0,
    true
);

INSERT INTO character_dialogues (character_id, notification_type, dialogue_line)
VALUES
    ('a0000000-0000-4000-8000-000000000001', 'streak_reminder', 'Your streak is waiting. A few minutes of practice keeps the momentum going.'),
    ('a0000000-0000-4000-8000-000000000001', 'streak_reminder', 'No practice yet today. Your league mates are already ahead.'),
    ('a0000000-0000-4000-8000-000000000001', 'streak_broken', 'The streak ended, but every expert was once a beginner who showed up again.'),
    ('a0000000-0000-4000-8000-000000000001', 'level_up', 'Level up! Your consistency is paying off.'),
    ('a0000000-0000-4000-8000-000000000001', 'weekend_challenge_available', 'Weekend challenge is live. Bigger problem, bigger reward.');
