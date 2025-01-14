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
	}
	Auth interface {
		RegisterUser(context.Context, *models.UserPayload) error
		LoginUser(context.Context, *models.LoginPayload) (string, error)
	}
	Post interface {
		CreatePost(context.Context, *models.PostPayload) error
	}
}

func NewService(db *sql.DB, auth auth.Authenticator, cloudinary cldnary.ClientCloudinary) Service {
	storage := postgresql.NewStorage(db)
	return Service{
		Users: &UserService{
			storage: &storage,
			auth:    auth,
		},
		Auth: &AuthService{
			storage: &storage,
			auth:    auth,
		},
		Post: &PostService{
			storage:    &storage,
			cloudinary: cloudinary,
		},
	}
}
