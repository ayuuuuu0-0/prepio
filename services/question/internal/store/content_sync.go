package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prepio/prepio/services/question/internal/lesson"
)

// SyncReport counts what one content sync changed.
type SyncReport struct {
	Worlds            int
	Nodes             int
	LessonsCreated    int
	LessonsNewVersion int
	LessonsUnchanged  int
	LessonsDeprecated int
	WorldsDeprecated  int
	NodesDeprecated   int
}

// ContentSyncStore loads authored content into the database.
type ContentSyncStore struct {
	pool *pgxpool.Pool
}

// NewContentSyncStore creates a ContentSyncStore.
func NewContentSyncStore(pool *pgxpool.Pool) *ContentSyncStore {
	return &ContentSyncStore{pool: pool}
}

// SkillSlugs returns every skill slug in the catalog (for validation).
func (s *ContentSyncStore) SkillSlugs(ctx context.Context) (map[string]bool, error) {
	rows, err := s.pool.Query(ctx, `SELECT slug FROM skills`)
	if err != nil {
		return nil, fmt.Errorf("list skills: %w", err)
	}
	defer rows.Close()
	slugs := map[string]bool{}
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			return nil, fmt.Errorf("scan skill slug: %w", err)
		}
		slugs[slug] = true
	}
	return slugs, rows.Err()
}

// Sync upserts all content by stable slug in one transaction. It is idempotent:
// running it twice changes nothing. It never deletes: content removed from the
// repo is deprecated, and a lesson whose content changed gets a new version while
// older versions stay available to attempts already in flight.
func (s *ContentSyncStore) Sync(ctx context.Context, c *lesson.Content) (*SyncReport, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin sync: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	report := &SyncReport{}
	skillIDs, err := loadSkillIDs(ctx, tx)
	if err != nil {
		return nil, err
	}

	nodeIDs, err := syncWorlds(ctx, tx, c, report)
	if err != nil {
		return nil, err
	}
	if err := syncPrerequisites(ctx, tx, c, nodeIDs); err != nil {
		return nil, err
	}
	if err := syncLessons(ctx, tx, c, nodeIDs, skillIDs, report); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit sync: %w", err)
	}
	return report, nil
}

func loadSkillIDs(ctx context.Context, tx pgx.Tx) (map[string]string, error) {
	rows, err := tx.Query(ctx, `SELECT slug, id FROM skills`)
	if err != nil {
		return nil, fmt.Errorf("load skills: %w", err)
	}
	defer rows.Close()
	ids := map[string]string{}
	for rows.Next() {
		var slug, id string
		if err := rows.Scan(&slug, &id); err != nil {
			return nil, fmt.Errorf("scan skill: %w", err)
		}
		ids[slug] = id
	}
	return ids, rows.Err()
}

// syncWorlds upserts worlds and nodes and returns node slug -> node id.
func syncWorlds(ctx context.Context, tx pgx.Tx, c *lesson.Content, report *SyncReport) (map[string]string, error) {
	nodeIDs := map[string]string{}
	worldSlugs := make([]string, 0, len(c.Worlds))
	nodeSlugs := []string{}

	for _, w := range c.Worlds {
		var worldID string
		err := tx.QueryRow(ctx, `
			INSERT INTO worlds (slug, name, description, theme, sort_order, status)
			VALUES ($1, $2, $3, $4, $5, 'published')
			ON CONFLICT (slug) DO UPDATE
			SET name = EXCLUDED.name, description = EXCLUDED.description,
			    theme = EXCLUDED.theme, sort_order = EXCLUDED.sort_order, status = 'published'
			RETURNING id`,
			w.Slug, w.Name, w.Description, themeOrDefault(w.Theme), w.Order).Scan(&worldID)
		if err != nil {
			return nil, fmt.Errorf("upsert world %s: %w", w.Slug, err)
		}
		worldSlugs = append(worldSlugs, w.Slug)
		report.Worlds++

		for i, n := range w.Nodes {
			var nodeID string
			err := tx.QueryRow(ctx, `
				INSERT INTO journey_nodes (world_id, slug, label, node_type, sort_order, status)
				VALUES ($1, $2, $3, $4, $5, 'published')
				ON CONFLICT (slug) WHERE slug IS NOT NULL DO UPDATE
				SET world_id = EXCLUDED.world_id, label = EXCLUDED.label,
				    node_type = EXCLUDED.node_type, sort_order = EXCLUDED.sort_order, status = 'published'
				RETURNING id`,
				worldID, n.Slug, n.Label, n.Type, i+1).Scan(&nodeID)
			if err != nil {
				return nil, fmt.Errorf("upsert node %s: %w", n.Slug, err)
			}
			nodeIDs[n.Slug] = nodeID
			nodeSlugs = append(nodeSlugs, n.Slug)
			report.Nodes++
		}
	}

	tag, err := tx.Exec(ctx, `UPDATE worlds SET status = 'deprecated' WHERE status = 'published' AND slug <> ALL($1)`, worldSlugs)
	if err != nil {
		return nil, fmt.Errorf("deprecate worlds: %w", err)
	}
	report.WorldsDeprecated = int(tag.RowsAffected())

	tag, err = tx.Exec(ctx, `UPDATE journey_nodes SET status = 'deprecated' WHERE status = 'published' AND (slug IS NULL OR slug <> ALL($1))`, nodeSlugs)
	if err != nil {
		return nil, fmt.Errorf("deprecate nodes: %w", err)
	}
	report.NodesDeprecated = int(tag.RowsAffected())
	return nodeIDs, nil
}

func themeOrDefault(theme string) string {
	if theme == "" {
		return "default"
	}
	return theme
}

// syncPrerequisites replaces the unlock rules of every synced node. Unlock rules
// are authored configuration, not user data.
func syncPrerequisites(ctx context.Context, tx pgx.Tx, c *lesson.Content, nodeIDs map[string]string) error {
	ids := make([]string, 0, len(nodeIDs))
	for _, id := range nodeIDs {
		ids = append(ids, id)
	}
	if _, err := tx.Exec(ctx, `DELETE FROM node_prerequisites WHERE node_id = ANY($1::uuid[])`, ids); err != nil {
		return fmt.Errorf("clear prerequisites: %w", err)
	}
	for _, w := range c.Worlds {
		for _, n := range w.Nodes {
			for _, req := range n.Requires {
				if _, err := tx.Exec(ctx,
					`INSERT INTO node_prerequisites (node_id, requires_node_id) VALUES ($1, $2)`,
					nodeIDs[n.Slug], nodeIDs[req]); err != nil {
					return fmt.Errorf("insert prerequisite %s -> %s: %w", n.Slug, req, err)
				}
			}
		}
	}
	return nil
}

func syncLessons(ctx context.Context, tx pgx.Tx, c *lesson.Content, nodeIDs, skillIDs map[string]string, report *SyncReport) error {
	slugs := make([]string, 0, len(c.Lessons))
	for _, l := range c.Lessons {
		slugs = append(slugs, l.Slug)
		hash, err := l.Hash()
		if err != nil {
			return err
		}
		summary, err := json.Marshal(l.Summary)
		if err != nil {
			return fmt.Errorf("encode summary %s: %w", l.Slug, err)
		}
		nodeID := nodeIDs[l.Node]

		var lessonID, storedHash string
		var version int
		err = tx.QueryRow(ctx, `SELECT id, version, content_hash FROM lessons WHERE slug = $1`, l.Slug).Scan(&lessonID, &version, &storedHash)
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			version = 1
			err = tx.QueryRow(ctx, `
				INSERT INTO lessons (slug, node_id, title, summary, kind, difficulty, est_minutes, status, version, content_hash)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 1, $9)
				RETURNING id`,
				l.Slug, nodeID, l.Title, summary, l.Kind, l.Difficulty, l.EstMinutes, l.Status, hash).Scan(&lessonID)
			if err != nil {
				return fmt.Errorf("insert lesson %s: %w", l.Slug, err)
			}
			if err := insertLessonVersion(ctx, tx, lessonID, version, l, skillIDs); err != nil {
				return err
			}
			report.LessonsCreated++
		case err != nil:
			return fmt.Errorf("look up lesson %s: %w", l.Slug, err)
		case storedHash != hash:
			version++
			_, err = tx.Exec(ctx, `
				UPDATE lessons
				SET node_id = $2, title = $3, summary = $4, kind = $5, difficulty = $6,
				    est_minutes = $7, status = $8, version = $9, content_hash = $10
				WHERE id = $1`,
				lessonID, nodeID, l.Title, summary, l.Kind, l.Difficulty, l.EstMinutes, l.Status, version, hash)
			if err != nil {
				return fmt.Errorf("update lesson %s: %w", l.Slug, err)
			}
			if err := insertLessonVersion(ctx, tx, lessonID, version, l, skillIDs); err != nil {
				return err
			}
			report.LessonsNewVersion++
		default:
			// Same content: only the lifecycle status and node binding may differ.
			_, err = tx.Exec(ctx, `
				UPDATE lessons SET status = $2, node_id = $3
				WHERE id = $1 AND (status <> $2 OR node_id <> $3)`,
				lessonID, l.Status, nodeID)
			if err != nil {
				return fmt.Errorf("update lesson status %s: %w", l.Slug, err)
			}
			report.LessonsUnchanged++
		}
	}

	tag, err := tx.Exec(ctx, `UPDATE lessons SET status = 'deprecated' WHERE status <> 'deprecated' AND slug <> ALL($1)`, slugs)
	if err != nil {
		return fmt.Errorf("deprecate lessons: %w", err)
	}
	report.LessonsDeprecated = int(tag.RowsAffected())
	return nil
}

func insertLessonVersion(ctx context.Context, tx pgx.Tx, lessonID string, version int, l lesson.Lesson, skillIDs map[string]string) error {
	for i, step := range l.Steps {
		payload, err := step.StoredPayload()
		if err != nil {
			return fmt.Errorf("lesson %s: %w", l.Slug, err)
		}
		raw, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("encode step %s/%s: %w", l.Slug, step.ID, err)
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO lesson_steps (lesson_id, lesson_version, step_key, position, type, payload)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			lessonID, version, step.ID, i, step.Type, raw); err != nil {
			return fmt.Errorf("insert step %s/%s: %w", l.Slug, step.ID, err)
		}
	}
	for _, sk := range l.Skills {
		skillID, ok := skillIDs[sk.Skill]
		if !ok {
			return fmt.Errorf("lesson %s: skill %q not in catalog", l.Slug, sk.Skill)
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO lesson_skills (lesson_id, lesson_version, skill_id, weight)
			VALUES ($1, $2, $3, $4)`,
			lessonID, version, skillID, sk.Weight); err != nil {
			return fmt.Errorf("insert lesson skill %s/%s: %w", l.Slug, sk.Skill, err)
		}
	}
	return nil
}
