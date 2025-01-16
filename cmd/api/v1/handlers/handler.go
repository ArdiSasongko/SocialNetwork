package handlers

import (
	"database/sql"
	"net/http"

	"github.com/ArdiSasongko/SocialNetwork/internal/auth"
	"github.com/ArdiSasongko/SocialNetwork/internal/service"
	"github.com/ArdiSasongko/SocialNetwork/internal/storage/cldnary"
	"github.com/ArdiSasongko/SocialNetwork/utils"
)

type Handler struct {
	Health interface {
		Get(w http.ResponseWriter, r *http.Request)
	}
	Users interface {
		GetProfile(w http.ResponseWriter, r *http.Request)
		UpdateImages(w http.ResponseWriter, r *http.Request)
		UpdateUser(w http.ResponseWriter, r *http.Request)
		GetUserProfile(w http.ResponseWriter, r *http.Request)
		FollowUser(w http.ResponseWriter, r *http.Request)
		UnfollowUser(w http.ResponseWriter, r *http.Request)
	}
	Auth interface {
		RegisterUser(w http.ResponseWriter, r *http.Request)
		LoginUser(w http.ResponseWriter, r *http.Request)
	}
	Post interface {
		CreatePost(w http.ResponseWriter, r *http.Request)
		GetPostByID(w http.ResponseWriter, r *http.Request)
		UpdatePost(w http.ResponseWriter, r *http.Request)
		DeletePost(w http.ResponseWriter, r *http.Request)
		CheckOwnerPost(allowRole string, next http.HandlerFunc) http.HandlerFunc
		GetPostByUser(w http.ResponseWriter, r *http.Request)
	}
	Feed interface {
		GetFeeds(w http.ResponseWriter, r *http.Request)
	}
}

func NewHandler(db *sql.DB, auth auth.Authenticator, cld cldnary.ClientCloudinary) Handler {
	service := service.NewService(db, auth, cld)
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
		Post: &PostHandler{
			service: service,
			json:    json,
			error:   error,
		},
		Feed: &FeedHandler{
			service: service,
			json:    json,
			error:   error,
		},
	}
}
