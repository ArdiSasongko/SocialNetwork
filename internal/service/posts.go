package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ArdiSasongko/SocialNetwork/internal/models"
	"github.com/ArdiSasongko/SocialNetwork/internal/storage/cldnary"
	"github.com/ArdiSasongko/SocialNetwork/internal/storage/postgresql"
)

const folderPost = "Posts"

type PostService struct {
	storage    *postgresql.Storage
	cloudinary cldnary.ClientCloudinary
}

func (s *PostService) CreatePost(ctx context.Context, payload *models.PostPayload) error {
	posts := postgresql.Post{
		UserID:  payload.UserID,
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
	}

	imagesPayloads := []postgresql.ImagePost{}
	publicIDs := []string{}

	// mapping to imageurl post
	for i, image := range payload.Images {
		imgUrl, publicID, err := s.cloudinary.UploadImage(ctx, image, folderPost)
		if err != nil {
			return err
		}

		filename := generateFilename(payload.Title, i+1)

		imagePayload := postgresql.ImagePost{
			ImageURL:  imgUrl,
			ImageName: filename,
		}

		publicIDs = append(publicIDs, publicID)
		imagesPayloads = append(imagesPayloads, imagePayload)
	}

	if err := s.storage.Posts.CreatePost(ctx, &posts, imagesPayloads); err != nil {
		var errRollback error
		for _, id := range publicIDs {
			if err := s.cloudinary.DeleteImage(ctx, id); err != nil {
				errRollback = err
			}
		}
		if errRollback != nil {
			return fmt.Errorf("rollback failed: %w", errRollback)
		}
		return err
	}

	return nil
}

func (s *PostService) UpdatePost(ctx context.Context, post *postgresql.Post, payload *models.PostUpdatePayload) error {
	if payload.Title != nil {
		post.Title = *payload.Title
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}

	if payload.Tags != nil {
		post.Tags = *payload.Tags
	}

	if err := s.storage.Posts.UpdatePost(ctx, post); err != nil {
		return err
	}

	return nil
}

func (s *PostService) DeletePost(ctx context.Context, postID int64) error {
	if err := s.storage.Posts.DeletePost(ctx, postID); err != nil {
		return err
	}

	return nil
}

func generateFilename(title string, index int) string {
	parts := strings.Split(title, " ")
	if len(parts) == 0 {
		return fmt.Sprintf("file-%d-%d", index, time.Now().Unix())
	}
	return fmt.Sprintf("%s-%d-%d", parts[0], index, time.Now().Unix())
}
