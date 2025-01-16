package service

import (
	"context"
	"database/sql"

	"github.com/ArdiSasongko/SocialNetwork/internal/auth"
	"github.com/ArdiSasongko/SocialNetwork/internal/models"
	"github.com/ArdiSasongko/SocialNetwork/internal/storage/cldnary"
	"github.com/ArdiSasongko/SocialNetwork/internal/storage/postgresql"
)

type Service struct {
	Users interface {
		GetProfileByID(context.Context, int64) (*models.UserResponse, error)
		UpdateProfile(context.Context, *models.UpdateImagePayload) error
		UpdateUser(context.Context, *postgresql.User, *models.UserUpdatePayload) error
		FollowUser(context.Context, int64, int64) error
		UnfollowUser(context.Context, int64, int64) error
	}
	Auth interface {
		RegisterUser(context.Context, *models.UserPayload) error
		LoginUser(context.Context, *models.LoginPayload) (string, error)
	}
	Post interface {
		CreatePost(context.Context, *models.PostPayload) error
		UpdatePost(context.Context, *postgresql.Post, *models.PostUpdatePayload) error
		DeletePost(ctx context.Context, postID int64) error
	}
	Role interface {
		GetRole(context.Context, string) (*postgresql.Role, error)
	}
	Feeds interface {
		GetFeeds(context.Context, int64, postgresql.Pagination) (models.FeedsResponse, error)
	}
}

func NewService(db *sql.DB, auth auth.Authenticator, cloudinary cldnary.ClientCloudinary) Service {
	storage := postgresql.NewStorage(db)
	return Service{
		Users: &UserService{
			storage:    &storage,
			auth:       auth,
			cloudinary: cloudinary,
		},
		Auth: &AuthService{
			storage:    &storage,
			auth:       auth,
			cloudinary: cloudinary,
		},
		Post: &PostService{
			storage:    &storage,
			cloudinary: cloudinary,
		},
		Role: &RoleService{
			storage: &storage,
		},
		Feeds: &FeedService{
			storage: &storage,
		},
	}
}
