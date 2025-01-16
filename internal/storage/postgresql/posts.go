package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64       `json:"id"`
	UserID    int64       `json:"user_id"`
	Title     string      `json:"title"`
	Content   string      `json:"content"`
	Tags      []string    `json:"tags"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
	IsEdited  bool        `json:"is_edited"`
	Images    []ImagePost `json:"images"`
	User      User        `json:"user"`
}

type ImagePost struct {
	ImageName string `json:"image_name"`
	PostID    int64  `json:"post_id"`
	ImageURL  string `json:"image_url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type PostWithMetaData struct {
	Post
}
type PostStore struct {
	db *sql.DB
}

func (p *PostStore) insertPost(ctx context.Context, tx *sql.Tx, post *Post) (*Post, error) {
	query := `
		INSERT INTO posts (user_id, title, content, tags)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	if err := tx.QueryRowContext(
		ctx,
		query,
		post.UserID,
		post.Title,
		post.Content,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to insert post, error : %v", err)
	}

	return post, nil
}

func (s *PostStore) insertImage(ctx context.Context, tx *sql.Tx, postID int64, imagePost ImagePost) error {
	query := `
		INSERT INTO images_post (image_name, image_url, post_id)
		VALUES ($1, $2, $3)
	`
	_, err := tx.ExecContext(ctx, query, imagePost.ImageName, imagePost.ImageURL, postID)
	if err != nil {
		return fmt.Errorf("failed to insertt image, error : %v", err)
	}

	return nil
}

func (s *PostStore) CreatePost(ctx context.Context, p *Post, images []ImagePost) error {
	ctx, cancel := context.WithTimeout(ctx, TimeoutCtx)
	defer cancel()

	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		user, err := s.insertPost(ctx, tx, p)
		if err != nil {
			return err
		}

		for _, image := range images {
			if err := s.insertImage(ctx, tx, user.ID, image); err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *PostStore) GetByID(ctx context.Context, tx *sql.Tx, postID int64) (*Post, error) {
	query := `
		SELECT id, user_id, title, content, tags, created_at, updated_at, is_edited
		FROM posts
		WHERE id = $1
	`

	post := new(Post)
	if err := tx.QueryRowContext(
		ctx,
		query,
		postID,
	).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.IsEdited,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return post, nil
}

func (s *PostStore) getImageByID(ctx context.Context, tx *sql.Tx, postID int64) ([]ImagePost, error) {
	query := `
		SELECT image_name, image_url, post_id, created_at
		FROM images_post
		WHERE post_id = $1
	`

	rows, err := tx.QueryContext(
		ctx,
		query,
		postID,
	)

	if err != nil {
		return nil, err
	}

	var (
		images []ImagePost
		image  ImagePost
	)

	for rows.Next() {
		if err := rows.Scan(
			&image.ImageName,
			&image.ImageURL,
			&image.PostID,
			&image.CreatedAt,
		); err != nil {
			return nil, err
		}
		images = append(images, image)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return images, nil
}

func (s *PostStore) GetPostByID(ctx context.Context, postID int64) (*Post, error) {
	ctx, cancel := context.WithTimeout(ctx, TimeoutCtx)
	defer cancel()

	result := new(Post)

	return result, withTx(s.db, ctx, func(tx *sql.Tx) error {
		// fetch post
		post, err := s.GetByID(ctx, tx, postID)
		if err != nil {
			return err
		}
		result = post

		// fetch images
		images, err := s.getImageByID(ctx, tx, postID)
		if err != nil {
			return err
		}
		result.Images = images

		return nil
	})
}

func (s *PostStore) UpdatePost(ctx context.Context, p *Post) error {
	query := `
		UPDATE posts
		SET title = $1, content = $2, tags = $3, is_edited = true, updated_at = NOW()
		WHERE id = $4
	`

	ctx, cancel := context.WithTimeout(ctx, TimeoutCtx)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, p.Title, p.Content, pq.Array(p.Tags), p.ID)
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

func (s *PostStore) DeletePost(ctx context.Context, postID int64) error {
	query := `
		DELETE FROM posts
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, TimeoutCtx)
	defer cancel()

	res, err := s.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
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

func (s *PostStore) GetByUser(ctx context.Context, userID int64) (*[]Post, error) {
	query := `
	SELECT id, title, content, tags, is_edited
	FROM posts
	WHERE user_id = $1
	`

	var (
		post  Post
		posts []Post
	)

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			pq.Array(&post.Tags),
			&post.IsEdited,
		); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return &posts, nil

}

func (s *PostStore) GetFeeds(ctx context.Context, userID int64, pf Pagination) ([]PostWithMetaData, error) {
	queryBuilder := strings.Builder{}
	params := make([]interface{}, 0)
	paramsCount := 1

	// base query
	queryBuilder.WriteString(`
		SELECT 
			p.id, 
			p.user_id, 
			u.username, 
			p.title, 
			p.content, 
			p.tags, 
			p.is_edited, 
			p.created_at, 
			p.updated_at
		FROM posts p
		LEFT JOIN users u ON u.id = p.user_id
		LEFT JOIN follows f ON f.user_id = p.user_id AND f.follower_id = $1
			WHERE p.user_id = $1  -- Own posts
   		OR f.user_id IS NOT NULL
	`)

	params = append(params, userID)

	// condition query
	queryBuilder.WriteString(` GROUP BY p.id, u.username ORDER BY p.created_at ` + pf.Sort)
	queryBuilder.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramsCount+1, paramsCount+2))
	params = append(params, pf.Limit, pf.Offset)

	ctx, cancel := context.WithTimeout(ctx, TimeoutCtx)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, queryBuilder.String(), params...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query :%w", err)
	}
	defer rows.Close()

	var feeds []PostWithMetaData

	for rows.Next() {
		var feed PostWithMetaData
		if err := rows.Scan(
			&feed.Post.ID,
			&feed.Post.UserID,
			&feed.Post.User.Username,
			&feed.Post.Title,
			&feed.Post.Content,
			pq.Array(&feed.Post.Tags),
			&feed.Post.IsEdited,
			&feed.Post.CreatedAt,
			&feed.Post.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		tx, err := s.db.BeginTx(ctx, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to begin transaction: %w", err)
		}
		images, err := s.getImageByID(ctx, tx, feed.ID)
		if err != nil {
			return nil, err
		}

		feed.Images = images

		feeds = append(feeds, feed)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error for iterating rows :%w", err)
	}

	return feeds, nil
}
