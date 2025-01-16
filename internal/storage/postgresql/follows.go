package postgresql

import (
	"context"
	"database/sql"
	"errors"
)

type FollowStore struct {
	db *sql.DB
}

func (s *FollowStore) FollowUser(ctx context.Context, userID, toFollow int64) error {
	query := `
		INSERT INTO follows (user_id, follower_id)
		VALUES ($1, $2)
	`

	ctx, cancel := context.WithTimeout(ctx, TimeoutCtx)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, query, toFollow, userID); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		case err.Error() == `pq: duplicate key value violates unique constraint "follows_pkey"`:
			return ErrConflict
		default:
			return err
		}
	}

	return nil
}

func (s *FollowStore) UnfollowUser(ctx context.Context, userID, toUnfollow int64) error {
	query := `
		DELETE FROM follows
		WHERE user_id = $1 AND follower_id =$2
	`
	ctx, cancel := context.WithTimeout(ctx, TimeoutCtx)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, toUnfollow, userID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
