package models

type FeedsResponse struct {
	Post []PostResponse `json:"post"`
}

type PostResponse struct {
	ID        int64             `json:"id"`
	Title     string            `json:"title"`
	Content   string            `json:"content"`
	Tags      []string          `json:"tags"`
	Images    []ImageResponse   `json:"images"`
	IsEdited  bool              `json:"is_edited"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
	Comments  []CommentResponse `json:"comments"`
	User      UserFeedResponse  `json:"user"`
}

type ImageResponse struct {
	ImageUrl  string `json:"image_url"`
	ImageName string `json:"image_name"`
}

type CommentResponse struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	IsEdited  bool   `json:"is_edited"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserFeedResponse struct {
	Username string `json:"username"`
	UserID   int64  `json:"user_id"`
}
