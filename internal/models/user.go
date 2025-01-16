package models

import (
	"mime/multipart"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

type UserPayload struct {
	Username     string                `json:"username" form:"username" validate:"required,min=3,max=255"`
	Fullname     string                `json:"fullname" form:"fullname" validate:"required,min=3,max=255"`
	Email        string                `json:"email" form:"email" validate:"required,email,max=72"`
	Password     string                `json:"password" form:"password" validate:"required,min=5,max=72"`
	ImageProfile *multipart.FileHeader `json:"images" form:"images" validate:"omitempty"`
}

func (u *UserPayload) Validate() error {
	return Validate.Struct(u)
}

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email,max=72"`
	Password string `json:"password" validate:"required,min=5,max=72"`
}

func (u *LoginPayload) Validate() error {
	return Validate.Struct(u)
}

type UserResponse struct {
	ID           int64                 `json:"id"`
	Username     string                `json:"username"`
	Fullname     string                `json:"fullname"`
	Email        string                `json:"email"`
	ImageProfile ImageUserResponse     `json:"image_profile"`
	CreatedAt    string                `json:"created_at"`
	UpdatedAt    string                `json:"updated_at"`
	Posts        []PostsByUserResponse `json:"posts"`
}

type ImageUserResponse struct {
	ImageURL string `json:"image_url"`
}

type UpdateImagePayload struct {
	Image  *multipart.FileHeader `json:"image" form:"image" validate:"required"`
	UserID int64                 `json:"user_id"`
}

func (u *UpdateImagePayload) Validate() error {
	return Validate.Struct(u)
}

type UserUpdatePayload struct {
	Username *string `json:"username" form:"username" validate:"omitempty,min=3,max=255"`
	Fullname *string `json:"fullname" form:"fullname" validate:"omitempty,min=3,max=255"`
}

func (u *UserUpdatePayload) Validate() error {
	return Validate.Struct(u)
}
