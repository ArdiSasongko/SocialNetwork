package service

import (
	"context"
	"fmt"

	"github.com/ArdiSasongko/SocialNetwork/internal/auth"
	"github.com/ArdiSasongko/SocialNetwork/internal/models"
	"github.com/ArdiSasongko/SocialNetwork/internal/storage/cldnary"
	"github.com/ArdiSasongko/SocialNetwork/internal/storage/postgresql"
)

type UserService struct {
	storage    *postgresql.Storage
	auth       auth.Authenticator
	cloudinary cldnary.ClientCloudinary
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

func (s *UserService) UpdateProfile(ctx context.Context, payload *models.UpdateImagePayload) error {
	imgUrl, publicID, err := s.cloudinary.UploadImage(ctx, payload.Image, folderProfile)
	if err != nil {
		return err
	}

	imageUpdate := postgresql.ImgURL{
		ImageURL: imgUrl,
		UserID:   payload.UserID,
	}

	if err := s.storage.Users.UpdateProfile(ctx, &imageUpdate); err != nil {
		if err := s.cloudinary.DeleteImage(ctx, publicID); err != nil {
			return err
		}
		return err
	}

	return nil
}

func (s *UserService) UpdateUser(ctx context.Context, user *postgresql.User, payload *models.UserUpdatePayload) error {
	if payload.Username == &user.Username {
		return fmt.Errorf("your username already %s, please try another username", user.Username)
	}

	if payload.Fullname == &user.Fullname {
		return fmt.Errorf("your fullname already %s, please try another username", user.Fullname)
	}

	if payload.Fullname != nil {
		user.Fullname = *payload.Fullname
	}

	if payload.Username != nil {
		user.Username = *payload.Username
	}

	if err := s.storage.Users.UpdateUser(ctx, user); err != nil {
		return err
	}

	return nil
}
