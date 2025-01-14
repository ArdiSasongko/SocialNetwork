package models

import "mime/multipart"

type PostPayload struct {
	UserID  int64                   `json:"user_id" form:"user_id"`
	Title   string                  `json:"title" form:"title" validate:"required,min=5,max=255"`
	Content string                  `json:"content" form:"content" validate:"required,min=10"`
	Tags    []string                `json:"tags" form:"tags" validate:"omitempty"`
	Images  []*multipart.FileHeader `json:"images" form:"images" validate:"omitempty"`
}

func (u *PostPayload) Validate() error {
	return Validate.Struct(u)
}

type ImagePayload struct {
	ImageUrl  string `json:"image_url"`
	ImageName string `json:"image_name"`
}
