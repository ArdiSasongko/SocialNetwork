package postgresql

import (
	"context"
	"database/sql"
	"fmt"
)

type Activities struct {
	ID         int64  `json:"id"`
	UserID     int64  `json:"user_id"`
	PostID     int64  `json:"post_id"`
	IsLiked    bool   `json:"is_liked"`
	IsDisliked bool   `json:"is_disliked"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type UserActivities struct {
	db *sql.DB
}

func (s *UserActivities) ToggleLikePost(ctx context.Context, us *Activities) error {
	query := `
        INSERT INTO user_activities (user_id, post_id, is_liked, is_disliked)
        VALUES ($1, $2, TRUE, FALSE)
        ON CONFLICT (user_id, post_id)
        DO UPDATE SET 
            is_liked = CASE 
                WHEN user_activities.is_liked = TRUE THEN FALSE 
                ELSE TRUE 
            END,
            is_disliked = FALSE,
            updated_at = CURRENT_TIMESTAMP
        RETURNING is_liked
    `

	ctx, cancel := context.WithTimeout(ctx, TimeoutCtx)
	defer cancel()

	var isLiked bool
	err := s.db.QueryRowContext(ctx, query, us.UserID, us.PostID).Scan(&isLiked)
	if err != nil {
		return fmt.Errorf("toggle like failed: %v", err)
	}

	us.IsLiked = isLiked
	us.IsDisliked = false
	return nil
}

func (s *UserActivities) ToggleDislikePost(ctx context.Context, us *Activities) error {
	query := `
        INSERT INTO user_activities (user_id, post_id, is_liked, is_disliked)
        VALUES ($1, $2, FALSE, TRUE)
        ON CONFLICT (user_id, post_id)
        DO UPDATE SET 
            is_disliked = CASE 
                WHEN user_activities.is_disliked = TRUE THEN FALSE 
                ELSE TRUE 
            END,
            is_liked = FALSE,
            updated_at = CURRENT_TIMESTAMP
        RETURNING is_disliked
    `

	ctx, cancel := context.WithTimeout(ctx, TimeoutCtx)
	defer cancel()

	var isDisliked bool
	err := s.db.QueryRowContext(ctx, query, us.UserID, us.PostID).Scan(&isDisliked)
	if err != nil {
		return fmt.Errorf("toggle dislike failed: %v", err)
	}

	us.IsDisliked = isDisliked
	us.IsLiked = false
	return nil
}

func (s *UserActivities) GetLikesByPost(ctx context.Context, postID int64) (int64, error) {
	query := `
        SELECT COUNT(*)
        FROM user_activities
        WHERE post_id = $1 AND is_liked = true
    `

	ctx, cancel := context.WithTimeout(ctx, TimeoutCtx)
	defer cancel()

	var count int64
	if err := s.db.QueryRowContext(ctx, query, postID).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (s *UserActivities) GetDislikesByPost(ctx context.Context, postID int64) (int64, error) {
	query := `
        SELECT COUNT(*)
        FROM user_activities
        WHERE post_id = $1 AND is_disliked = true
    `

	ctx, cancel := context.WithTimeout(ctx, TimeoutCtx)
	defer cancel()

	var count int64
	if err := s.db.QueryRowContext(ctx, query, postID).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}
