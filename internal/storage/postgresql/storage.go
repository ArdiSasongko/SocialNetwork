package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("resource not found")
	ErrDuplicateEmail    = errors.New("email already exists")
	ErrDuplicateUsername = errors.New("username already exists")
	ErrConflict          = errors.New("resource already exists")
	TimeoutCtx           = time.Second * 5
)

type Storage struct {
	Users interface {
		GetByID(context.Context, int64) (*User, error)
		GetByEmail(context.Context, string) (*User, error)
		CreateUser(context.Context, *User, *ImgURL) error
	}
	Posts interface {
		CreatePost(context.Context, *Post, []ImagePost) error
		UpdatePost(context.Context, *Post) error
		GetPostByID(context.Context, int64) (*Post, error)
		GetByID(context.Context, *sql.Tx, int64) (*Post, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Users: &UserStorage{
			db: db,
		},
		Posts: &PostStore{
			db: db,
		},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
