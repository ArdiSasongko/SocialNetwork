package service

import (
	"context"
	"sync"

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

	var (
		wg                                      sync.WaitGroup
		commentCount, likesCount, disLikedCount int64
		posts                                   []models.PostsResponse
	)

	for _, p := range respPost {
		var images []models.ImageResponse
		errChan := make(chan error, 3)
		for _, i := range p.Post.Images {
			images = append(images, models.ImageResponse{
				ImageUrl:  i.ImageURL,
				ImageName: i.ImageName,
			})
		}

		wg.Add(3)
		go func() {
			defer wg.Done()
			count, err := s.storage.Comments.GetCommentCountByPost(ctx, p.Post.ID)
			if err != nil {
				errChan <- err
				return
			}

			commentCount = count
		}()

		go func() {
			defer wg.Done()
			count, err := s.storage.Activities.GetLikesByPost(ctx, p.Post.ID)
			if err != nil {
				errChan <- err
				return
			}

			likesCount = count
		}()

		go func() {
			defer wg.Done()
			count, err := s.storage.Activities.GetDislikesByPost(ctx, p.Post.ID)
			if err != nil {
				errChan <- err
				return
			}

			disLikedCount = count
		}()

		wg.Wait()
		close(errChan)

		if err := <-errChan; err != nil {
			return models.FeedsResponse{}, err
		}

		post := models.PostsResponse{
			Username: p.Post.User.Username,
			Title:    p.Post.Title,
			Content:  p.Post.Content,
			Tags:     p.Post.Tags,
			Images:   images,
			MetaData: models.MetaData{
				CommentCount: commentCount,
				LikeCount:    likesCount,
				DislikeCount: disLikedCount,
			},
		}

		posts = append(posts, post)
	}

	return models.FeedsResponse{
		Posts: posts,
	}, nil
}

func (s *FeedService) GetFeed(ctx context.Context, postID int64) (models.PostResponse, error) {
	respPost, err := s.storage.Posts.GetPostByID(ctx, postID)
	if err != nil {
		return models.PostResponse{}, err
	}

	var (
		wg                                      sync.WaitGroup
		commentCount, likesCount, disLikedCount int64
		allComments                             []models.CommentResponse
		images                                  []models.ImageResponse
	)

	errChan := make(chan error, 4)
	wg.Add(4)
	go func() {
		defer wg.Done()
		count, err := s.storage.Comments.GetCommentCountByPost(ctx, postID)
		if err != nil {
			errChan <- err
			return
		}

		commentCount = count
	}()

	go func() {
		defer wg.Done()
		count, err := s.storage.Activities.GetLikesByPost(ctx, postID)
		if err != nil {
			errChan <- err
			return
		}

		likesCount = count
	}()

	go func() {
		defer wg.Done()
		count, err := s.storage.Activities.GetDislikesByPost(ctx, postID)
		if err != nil {
			errChan <- err
			return
		}

		disLikedCount = count
	}()

	go func() {
		defer wg.Done()
		comments, err := s.getComments(ctx, postID)
		if err != nil {
			errChan <- err
			return
		}
		allComments = comments
	}()
	wg.Wait()
	close(errChan)

	if err := <-errChan; err != nil {
		return models.PostResponse{}, err
	}

	for _, i := range respPost.Images {
		image := models.ImageResponse{
			ImageUrl:  i.ImageURL,
			ImageName: i.ImageName,
		}

		images = append(images, image)
	}

	return models.PostResponse{
		ID:        respPost.ID,
		Title:     respPost.Title,
		Content:   respPost.Content,
		Tags:      respPost.Tags,
		Images:    images,
		IsEdited:  respPost.IsEdited,
		CreatedAt: respPost.CreatedAt,
		UpdatedAt: respPost.UpdatedAt,
		User: models.UserFeedResponse{
			Username: respPost.User.Username,
			UserID:   respPost.ID,
		},
		Comments: allComments,
		MetaData: models.MetaData{
			CommentCount: commentCount,
			LikeCount:    likesCount,
			DislikeCount: disLikedCount,
		},
	}, nil
}

func (s *FeedService) getComments(ctx context.Context, postID int64) ([]models.CommentResponse, error) {
	commentResp, err := s.storage.Comments.GetCommentsByPost(ctx, postID)
	if err != nil {
		return nil, err
	}

	var (
		comments []models.CommentResponse
		comment  models.CommentResponse
	)

	for _, c := range commentResp {
		comment = models.CommentResponse{
			UserID:    c.ID,
			Content:   c.Content,
			IsEdited:  c.IsEdited,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

func (s *FeedService) LikePost(ctx context.Context, p *models.UserActivitiesPayload) error {
	liked := postgresql.Activities{
		UserID: p.UserID,
		PostID: p.PostID,
	}

	return s.storage.Activities.ToggleLikePost(ctx, &liked)
}

func (s *FeedService) DislikePost(ctx context.Context, p *models.UserActivitiesPayload) error {
	disliked := postgresql.Activities{
		UserID: p.UserID,
		PostID: p.PostID,
	}

	return s.storage.Activities.ToggleDislikePost(ctx, &disliked)
}

func (s *FeedService) CreateCommentPost(ctx context.Context, p *models.CommentPayload) error {
	comment := postgresql.Comment{
		UserID:  p.UserID,
		PostID:  p.PostID,
		Content: p.Content,
	}

	return s.storage.Comments.CreateComments(ctx, &comment)
}
