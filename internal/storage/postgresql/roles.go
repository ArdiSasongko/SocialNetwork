package postgresql

import (
	"context"
	"database/sql"
	"errors"
)

type RoleStore struct {
	db *sql.DB
}

type Role struct {
	ID          int64  `json:"id"`
	Level       int64  `json:"level"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *RoleStore) GetByName(ctx context.Context, roleName string) (*Role, error) {
	query := `
		SELECT id, level, name, description
		FROM roles
		WHERE name = $1
	`

	role := new(Role)
	if err := s.db.QueryRowContext(ctx, query, roleName).Scan(
		&role.ID,
		&role.Level,
		&role.Name,
		&role.Description,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return role, nil
}
