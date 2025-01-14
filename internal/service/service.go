package service

import (
	"context"
	"database/sql"

	"github.com/ArdiSasongko/SocialNetwork/internal/auth"
	"github.com/ArdiSasongko/SocialNetwork/internal/models"
	"github.com/ArdiSasongko/SocialNetwork/internal/storage/postgresql"
)

type Service struct {
	Users interface {
		RegisterUser(context.Context, *models.UserPayload) error
		LoginUser(context.Context, *models.LoginPayload) (string, error)
	}
}

func NewService(db *sql.DB, auth auth.Authenticator) Service {
	storage := postgresql.NewStorage(db)
	return Service{
		Users: &UserService{
			storage: &storage,
			auth:    auth,
		},
	}
}
