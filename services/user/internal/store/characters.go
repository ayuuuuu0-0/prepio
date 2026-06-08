package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prepio/prepio/constants"
)

// Character is a row from the characters table.
type Character struct {
	ID        string
	Name      string
	Species   string
	GemCost   int
	IsDefault bool
}

// CharacterStore handles characters table queries.
type CharacterStore struct {
	pool *pgxpool.Pool
}

// NewCharacterStore creates a CharacterStore.
func NewCharacterStore(pool *pgxpool.Pool) *CharacterStore {
	return &CharacterStore{pool: pool}
}

// ListStarters returns onboarding-selectable companions.
func (s *CharacterStore) ListStarters(ctx context.Context) ([]Character, error) {
	const q = `
		SELECT id, name, species, gem_cost, is_default
		FROM characters
		WHERE id = ANY($1)
		ORDER BY name`

	rows, err := s.pool.Query(ctx, q, constants.StarterCompanionIDs)
	if err != nil {
		return nil, fmt.Errorf("list starter characters: %w", err)
	}
	defer rows.Close()

	var characters []Character
	for rows.Next() {
		var c Character
		if err := rows.Scan(&c.ID, &c.Name, &c.Species, &c.GemCost, &c.IsDefault); err != nil {
			return nil, fmt.Errorf("scan character: %w", err)
		}
		characters = append(characters, c)
	}
	return characters, rows.Err()
}

// GetByID returns a character by ID.
func (s *CharacterStore) GetByID(ctx context.Context, id string) (*Character, error) {
	const q = `
		SELECT id, name, species, gem_cost, is_default
		FROM characters WHERE id = $1`

	var c Character
	err := s.pool.QueryRow(ctx, q, id).Scan(&c.ID, &c.Name, &c.Species, &c.GemCost, &c.IsDefault)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get character: %w", err)
	}
	return &c, nil
}
