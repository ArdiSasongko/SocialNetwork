package service

import (
	"context"

	"github.com/ArdiSasongko/SocialNetwork/internal/auth"
	"github.com/ArdiSasongko/SocialNetwork/internal/models"
	"github.com/ArdiSasongko/SocialNetwork/internal/storage/postgresql"
)

type UserService struct {
	storage *postgresql.Storage
	auth    auth.Authenticator
}

func (s *UserService) GetProfileByID(ctx context.Context, userID int64) (*models.UserResponse, error) {
	user, err := s.storage.Users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	userResponse := new(models.UserResponse)

	userResponse = &models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Fullname: user.Fullname,
		Email:    user.Email,
		ImageProfile: models.ImageUserResponse{
			ImageURL: user.ImgURL.ImageURL,
		},
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return userResponse, nil
}
