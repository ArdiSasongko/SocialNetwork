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

const (
	folder = "Posts"
)

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
		imgUrl, publicID, err := s.cloudinary.UploadImage(ctx, image, folder)
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

func generateFilename(title string, index int) string {
	parts := strings.Split(title, " ")
	if len(parts) == 0 {
		return fmt.Sprintf("file-%d-%d", index, time.Now().Unix())
	}
	return fmt.Sprintf("%s-%d-%d", parts[0], index, time.Now().Unix())
}
