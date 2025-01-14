package service

import (
	"context"
	"time"

	"github.com/ArdiSasongko/SocialNetwork/internal/auth"
	"github.com/ArdiSasongko/SocialNetwork/internal/env"
	"github.com/ArdiSasongko/SocialNetwork/internal/models"
	"github.com/ArdiSasongko/SocialNetwork/internal/storage/postgresql"
	"github.com/golang-jwt/jwt/v5"
)

type UserService struct {
	storage *postgresql.Storage
	auth    auth.Authenticator
}

func (s *UserService) RegisterUser(ctx context.Context, payload *models.UserPayload) error {
	user := postgresql.User{
		Username: payload.Username,
		Fullname: payload.Fullname,
		Email:    payload.Email,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		return err
	}

	// todo update image upload handler
	randomImg := "https://www.freepik.com/icon/man_11043660#fromView=search&page=1&position=25&uuid=ce1b8a3b-1eef-400a-9380-e32a9971f96c"
	img := postgresql.ImgURL{
		ImageURL: randomImg,
	}

	if err := s.storage.Users.CreateUser(ctx, &user, &img); err != nil {
		return err
	}

	return nil
}

func (s *UserService) LoginUser(ctx context.Context, payload *models.LoginPayload) (string, error) {
	user, err := s.storage.Users.GetByEmail(ctx, payload.Email)
	if err != nil {
		return "", err
	}

	if err := user.Password.Compared(payload.Password); err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 3).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": env.GetString("JWT_ISS", "SocialNetwork"),
		"aud": env.GetString("JWT_ISS", "SocialNetwork"),
	}

	token, err := s.auth.GenerateToken(claims)
	if err != nil {
		return "", err
	}

	return token, nil
}
