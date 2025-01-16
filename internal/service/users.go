package service

import (
	"context"
	"fmt"
	"sync"

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
	var (
		wg    sync.WaitGroup
		user  *postgresql.User
		posts []models.PostsByUserResponse
	)

	wg.Add(2)
	errChan := make(chan error, 2)

	// fetch user
	go func() {
		defer wg.Done()
		u, err := s.storage.Users.GetByID(ctx, userID)
		if err != nil {
			errChan <- fmt.Errorf("get user: %w", err)
			return
		}
		user = u
	}()

	// fetch posts
	go func() {
		defer wg.Done()
		p, err := s.getPostByUser(ctx, userID)
		if err != nil {
			errChan <- fmt.Errorf("get posts: %w", err)
			return
		}
		posts = p
	}()

	wg.Wait()
	close(errChan)

	// Check for errors
	if err := <-errChan; err != nil {
		return nil, err
	}

	// Jika user nil, berarti ada error yang terjadi
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	userResponse := &models.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Fullname: user.Fullname,
		Email:    user.Email,
		ImageProfile: models.ImageUserResponse{
			ImageURL: user.ImgURL.ImageURL,
		},
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Posts:     posts,
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

func (s *UserService) getPostByUser(ctx context.Context, userID int64) ([]models.PostsByUserResponse, error) {
	var (
		post  models.PostsByUserResponse
		posts []models.PostsByUserResponse
	)

	resp, err := s.storage.Posts.GetByUser(ctx, userID)
	if err != nil {
		return []models.PostsByUserResponse{}, err
	}

	for _, p := range *resp {
		post = models.PostsByUserResponse{
			ID:       p.ID,
			Title:    p.Title,
			Content:  p.Content,
			Tags:     p.Tags,
			IsEdited: p.IsEdited,
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (s *UserService) FollowUser(ctx context.Context, toFollow, userID int64) error {
	if toFollow == userID {
		return fmt.Errorf("invalid data")
	}

	return s.storage.Follows.FollowUser(ctx, userID, toFollow)
}

func (s *UserService) UnfollowUser(ctx context.Context, toUnfollow, userID int64) error {
	return s.storage.Follows.UnfollowUser(ctx, userID, toUnfollow)
}
