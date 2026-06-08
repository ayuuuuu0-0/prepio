INSERT INTO characters (id, name, species, gem_cost, is_default) VALUES
    ('b0000000-0000-4000-8000-000000000001', 'Byte', 'capybara', 0, true),
    ('b0000000-0000-4000-8000-000000000002', 'Pip', 'red_panda', 0, false),
    ('b0000000-0000-4000-8000-000000000003', 'Nova', 'pangolin', 0, false),
    ('b0000000-0000-4000-8000-000000000004', 'Kodo', 'axolotl', 0, false),
    ('b0000000-0000-4000-8000-000000000005', 'Zara', 'snow_leopard', 0, false)
ON CONFLICT (id) DO NOTHING;

INSERT INTO character_dialogues (character_id, notification_type, dialogue_line) VALUES
    ('b0000000-0000-4000-8000-000000000001', 'streak_reminder', 'Your streak is waiting — let''s level up your career today.'),
    ('b0000000-0000-4000-8000-000000000001', 'level_up', 'Level up! You''re getting interview-ready.'),
    ('b0000000-0000-4000-8000-000000000002', 'streak_reminder', 'Pip believes in you. One challenge keeps the momentum going.'),
    ('b0000000-0000-4000-8000-000000000003', 'level_up', 'Nova says: readiness unlocked!'),
    ('b0000000-0000-4000-8000-000000000004', 'streak_reminder', 'Kodo is cheering — don''t break the streak!'),
    ('b0000000-0000-4000-8000-000000000005', 'level_up', 'Zara roars — you leveled up!');
