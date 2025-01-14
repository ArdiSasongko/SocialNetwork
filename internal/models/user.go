package models

import "github.com/go-playground/validator/v10"

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

type UserPayload struct {
	Username string `json:"username" validate:"required,min=3,max=255"`
	Fullname string `json:"fullname" validate:"required,min=3,max=255"`
	Email    string `json:"email" validate:"required,email,max=72"`
	Password string `json:"password" validate:"required,min=5,max=72"`
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
