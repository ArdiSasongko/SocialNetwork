package models

type FeedsResponse struct {
	Posts []PostsResponse `json:"posts"`
}

type PostsResponse struct {
	Username string          `json:"username"`
	Title    string          `json:"title"`
	Content  string          `json:"content"`
	Tags     []string        `json:"tags"`
	Images   []ImageResponse `json:"images"`
	MetaData MetaData        `json:"meta_data"`
}

type MetaData struct {
	CommentCount int64 `json:"comments_count"`
	LikeCount    int64 `json:"dlike_count"`
	DislikeCount int64 `json:"dislike_count"`
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
	MetaData  MetaData          `json:"meta_data"`
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

type UserActivitiesPayload struct {
	UserID int64 `json:"user_id"`
	PostID int64 `json:"post_id"`
}

type CommentPayload struct {
	UserID  int64  `json:"user_id"`
	PostID  int64  `json:"post_id"`
	Content string `json:"content" validate:"required,max=255"`
}

func (u *CommentPayload) Validate() error {
	return Validate.Struct(u)
}
