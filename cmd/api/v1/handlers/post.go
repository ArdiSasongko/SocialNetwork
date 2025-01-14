package handlers

import (
	"log"
	"mime/multipart"
	"net/http"

	"github.com/ArdiSasongko/SocialNetwork/cmd/api/v1/middlewares"
	"github.com/ArdiSasongko/SocialNetwork/internal/models"
	"github.com/ArdiSasongko/SocialNetwork/internal/service"
	"github.com/ArdiSasongko/SocialNetwork/internal/storage/postgresql"
	"github.com/ArdiSasongko/SocialNetwork/utils"
)

type PostHandler struct {
	service service.Service
	json    utils.JsonUtils
	error   utils.ErrorUtils
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	user := getUserfromCtx(r)
	payload := new(models.PostPayload)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.error.BadRequestError(w, r, err)
		return
	}

	files, err := extractFiles(r, "images")
	if err != nil {
		h.error.BadRequestError(w, r, err)
		return
	}

	log.Println(user.ID)
	payload.Images = files
	payload.UserID = user.ID
	payload.Content = r.FormValue("content")
	payload.Title = r.FormValue("title")
	payload.Tags = r.Form["tags"]

	if err := payload.Validate(); err != nil {
		h.error.BadRequestError(w, r, err)
		return
	}

	if err := h.service.Post.CreatePost(r.Context(), payload); err != nil {
		h.error.InternalServerError(w, r, err)
		return
	}

	if err := h.json.JsonResponse(w, http.StatusCreated, nil); err != nil {
		h.error.InternalServerError(w, r, err)
		return
	}
}

func (h *PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	post := getPostfromCtx(r)

	if err := h.json.JsonResponse(w, http.StatusOK, post); err != nil {
		h.error.InternalServerError(w, r, err)
		return
	}
}

func extractFiles(r *http.Request, fieldName string) ([]*multipart.FileHeader, error) {
	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		return nil, nil
	}

	files, ok := r.MultipartForm.File[fieldName]
	if !ok || len(files) == 0 {
		return nil, nil
	}

	return files, nil
}

func getPostfromCtx(r *http.Request) *postgresql.Post {
	user, _ := r.Context().Value(middlewares.PostCtx).(*postgresql.Post)
	return user
}
