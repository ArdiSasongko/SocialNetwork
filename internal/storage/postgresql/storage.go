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
		UpdateProfile(context.Context, *ImgURL) error
		UpdateUser(context.Context, *User) error
	}
	Posts interface {
		CreatePost(context.Context, *Post, []ImagePost) error
		UpdatePost(context.Context, *Post) error
		GetPostByID(context.Context, int64) (*Post, error)
		GetByID(context.Context, *sql.Tx, int64) (*Post, error)
		DeletePost(context.Context, int64) error
		GetByUser(context.Context, int64) (*[]Post, error)
		GetFeeds(context.Context, int64, Pagination) ([]PostWithMetaData, error)
	}
	Roles interface {
		GetByName(context.Context, string) (*Role, error)
	}
	Follows interface {
		FollowUser(context.Context, int64, int64) error
		UnfollowUser(context.Context, int64, int64) error
	}
	Activities interface {
		ToggleLikePost(context.Context, *Activities) error
		ToggleDislikePost(context.Context, *Activities) error
		GetLikesByPost(context.Context, int64) (int64, error)
		GetDislikesByPost(context.Context, int64) (int64, error)
	}
	Comments interface {
		CreateComments(context.Context, *Comment) error
		GetCommentsByPost(context.Context, int64) ([]Comment, error)
		GetCommentCountByPost(context.Context, int64) (int64, error)
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
		Roles: &RoleStore{
			db: db,
		},
		Follows: &FollowStore{
			db: db,
		},
		Comments: &CommentStore{
			db: db,
		},
		Activities: &UserActivities{
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
