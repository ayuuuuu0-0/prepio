-- Backfill user_skill_scores from existing answer history and question_skills mappings.
-- Mastery formula matches config/readiness.go: contribution per answer, smoothed average capped at 100.
INSERT INTO user_skill_scores (user_id, skill_id, mastery, attempts, last_practiced_at, source)
SELECT
    agg.user_id,
    agg.skill_id,
    LEAST(config.max_mastery, GREATEST(0, ROUND(agg.avg_contribution * 100)))::int,
    agg.attempts,
    agg.last_practiced_at,
    'backfill'
FROM (
    SELECT
        h.user_id,
        qs.skill_id,
        COUNT(*)::int AS attempts,
        MAX(h.submitted_at) AS last_practiced_at,
        AVG(
            (h.score::numeric / 100.0)
            * q.readiness_weight
            * qs.weight
            * CASE q.difficulty
                WHEN 'easy' THEN 0.90
                WHEN 'hard' THEN 1.10
                ELSE 1.00
              END
        ) AS avg_contribution
    FROM user_question_history h
    JOIN question_skills qs ON qs.question_id = h.question_id
    JOIN questions q ON q.id = h.question_id
    GROUP BY h.user_id, qs.skill_id
) agg
CROSS JOIN (SELECT 100::numeric AS max_mastery) config
ON CONFLICT (user_id, skill_id) DO UPDATE SET
    mastery = EXCLUDED.mastery,
    attempts = EXCLUDED.attempts,
    last_practiced_at = EXCLUDED.last_practiced_at,
    source = EXCLUDED.source,
    updated_at = now();
