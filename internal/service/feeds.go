package service

import (
	"context"

	"github.com/ArdiSasongko/SocialNetwork/internal/models"
	"github.com/ArdiSasongko/SocialNetwork/internal/storage/postgresql"
)

type FeedService struct {
	storage *postgresql.Storage
}

func (s *FeedService) GetFeeds(ctx context.Context, userID int64, pf postgresql.Pagination) (models.FeedsResponse, error) {
	respPost, err := s.storage.Posts.GetFeeds(ctx, userID, pf)
	if err != nil {
		return models.FeedsResponse{}, err
	}

	var posts []models.PostResponse
	for _, p := range respPost {
		var images []models.ImageResponse
		for _, i := range p.Post.Images {
			images = append(images, models.ImageResponse{
				ImageUrl:  i.ImageURL,
				ImageName: i.ImageName,
			})
		}

		post := models.PostResponse{
			ID:        p.Post.ID,
			Title:     p.Post.Title,
			Content:   p.Post.Content,
			Tags:      p.Post.Tags,
			IsEdited:  p.Post.IsEdited,
			CreatedAt: p.Post.CreatedAt,
			UpdatedAt: p.Post.UpdatedAt,
			Images:    images,
			User: models.UserFeedResponse{
				Username: p.Post.User.Username,
				UserID:   p.Post.UserID,
			},
		}

		posts = append(posts, post)
	}

	return models.FeedsResponse{
		Post: posts,
	}, nil
}
