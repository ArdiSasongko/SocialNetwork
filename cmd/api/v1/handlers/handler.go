package handlers

import (
	"database/sql"
	"net/http"

	"github.com/ArdiSasongko/SocialNetwork/internal/auth"
	"github.com/ArdiSasongko/SocialNetwork/internal/service"
	"github.com/ArdiSasongko/SocialNetwork/utils"
)

type Handler struct {
	Health interface {
		Get(w http.ResponseWriter, r *http.Request)
	}
	Users interface {
		GetProfile(w http.ResponseWriter, r *http.Request)
	}
	Auth interface {
		RegisterUser(w http.ResponseWriter, r *http.Request)
		LoginUser(w http.ResponseWriter, r *http.Request)
	}
}

func NewHandler(db *sql.DB, auth auth.Authenticator) Handler {
	service := service.NewService(db, auth)
	json := utils.NewJsonUtils()
	error := utils.NewErrorUtils()
	return Handler{
		Health: &healthHandler{
			json:  json,
			error: error,
		},
		Users: &UserHandler{
			service: service,
			json:    json,
			error:   error,
		},
		Auth: &AuthHandler{
			service: service,
			json:    json,
			error:   error,
		},
	}
}
