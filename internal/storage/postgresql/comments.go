package postgresql

import (
	"context"
	"database/sql"
)

type Comment struct {
	ID           int64  `json:"id"`
	UserID       int64  `json:"user_id"`
	PostID       int64  `json:"post_id"`
	Content      string `json:"content"`
	IsEdited     bool   `json:"is_edited"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	CommentCount int64  `json:"comment_count"`
}

type CommentStore struct {
	db *sql.DB
}

func (s *CommentStore) CreateComments(ctx context.Context, c *Comment) error {
	query := `
		INSERT INTO comments (user_id, post_id, content)
		VALUES ($1, $2, $3)
	`

	ctx, cancel := context.WithTimeout(ctx, TimeoutCtx)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, query, c.UserID, c.PostID, c.Content); err != nil {
		return err
	}

	return nil
}

func (s *CommentStore) GetCommentsByPost(ctx context.Context, postID int64) ([]Comment, error) {
	query := `
        SELECT 
            id, user_id, post_id, content, created_at, updated_at, is_edited
        FROM comments
        WHERE post_id = $1
        ORDER BY created_at DESC
    `

	ctx, cancel := context.WithTimeout(ctx, TimeoutCtx)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(
			&c.ID,
			&c.UserID,
			&c.PostID,
			&c.Content,
			&c.CreatedAt,
			&c.UpdatedAt,
			&c.IsEdited,
		); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *CommentStore) GetCommentCountByPost(ctx context.Context, postID int64) (int64, error) {
	query := `
        SELECT COUNT(*)
        FROM comments
        WHERE post_id = $1
    `

	ctx, cancel := context.WithTimeout(ctx, TimeoutCtx)
	defer cancel()

	var count int64
	if err := s.db.QueryRowContext(ctx, query, postID).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}
