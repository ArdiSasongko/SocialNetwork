package service

import (
	"context"
	"math/rand/v2"
	"time"

	"github.com/ArdiSasongko/SocialNetwork/internal/auth"
	"github.com/ArdiSasongko/SocialNetwork/internal/env"
	"github.com/ArdiSasongko/SocialNetwork/internal/models"
	"github.com/ArdiSasongko/SocialNetwork/internal/storage/cldnary"
	"github.com/ArdiSasongko/SocialNetwork/internal/storage/postgresql"
	"github.com/golang-jwt/jwt/v5"
)

const folderProfile = "Profiles"

var defaultImage = []string{
	"https://res.cloudinary.com/drbxy46kq/image/upload/v1736916378/Minimalist_Avatar_1_ayc9ei.jpg",
	"https://res.cloudinary.com/drbxy46kq/image/upload/v1736916378/Minimalist_Avatar_3_lui2ue.jpg",
	"https://res.cloudinary.com/drbxy46kq/image/upload/v1736916378/Minimalist_Avatar_2_lbqoac.jpg",
}

type AuthService struct {
	storage    *postgresql.Storage
	auth       auth.Authenticator
	cloudinary cldnary.ClientCloudinary
}

func (s *AuthService) RegisterUser(ctx context.Context, payload *models.UserPayload) error {
	user := postgresql.User{
		Username: payload.Username,
		Fullname: payload.Fullname,
		Email:    payload.Email,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		return err
	}

	// todo update image upload handler
	imgUrl := defaultImage[rand.IntN(len(defaultImage))]

	img := postgresql.ImgURL{
		ImageURL: imgUrl,
	}

	if err := s.storage.Users.CreateUser(ctx, &user, &img); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) LoginUser(ctx context.Context, payload *models.LoginPayload) (string, error) {
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
