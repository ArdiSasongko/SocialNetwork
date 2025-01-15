package handlers

import (
	"net/http"

	"github.com/ArdiSasongko/SocialNetwork/cmd/api/v1/middlewares"
	"github.com/ArdiSasongko/SocialNetwork/internal/models"
	"github.com/ArdiSasongko/SocialNetwork/internal/service"
	"github.com/ArdiSasongko/SocialNetwork/internal/storage/postgresql"
	"github.com/ArdiSasongko/SocialNetwork/utils"
)

type UserHandler struct {
	service service.Service
	json    utils.JsonUtils
	error   utils.ErrorUtils
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	user := getUserfromCtx(r)

	if err := h.json.JsonResponse(w, http.StatusOK, user); err != nil {
		h.error.InternalServerError(w, r, err)
		return
	}
}

func (h *UserHandler) UpdateImages(w http.ResponseWriter, r *http.Request) {
	user := getUserfromCtx(r)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.error.BadRequestError(w, r, err)
		return
	}

	file, err := extractFile(r, "image")
	if err != nil {
		h.error.BadRequestError(w, r, err)
		return
	}

	payload := models.UpdateImagePayload{
		Image: file,
	}

	payload.UserID = user.ID

	if err := payload.Validate(); err != nil {
		h.error.BadRequestError(w, r, err)
		return
	}

	if err := h.service.Users.UpdateProfile(r.Context(), &payload); err != nil {
		h.error.InternalServerError(w, r, err)
		return
	}

	if err := h.json.JsonResponse(w, http.StatusOK, nil); err != nil {
		h.error.InternalServerError(w, r, err)
		return
	}
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := getUserfromCtx(r)
	payload := new(models.UserUpdatePayload)

	if err := h.json.ReadJSON(w, r, payload); err != nil {
		h.error.BadRequestError(w, r, err)
		return
	}

	if err := payload.Validate(); err != nil {
		h.error.BadRequestError(w, r, err)
		return
	}

	//log.Printf("payload : %v, user : %v", *payload.Username, user.Username)

	if err := h.service.Users.UpdateUser(r.Context(), user, payload); err != nil {
		h.error.InternalServerError(w, r, err)
		return
	}

	if err := h.json.JsonResponse(w, http.StatusOK, nil); err != nil {
		h.error.InternalServerError(w, r, err)
		return
	}
}

func getUserfromCtx(r *http.Request) *postgresql.User {
	user, _ := r.Context().Value(middlewares.UserCtx).(*postgresql.User)
	return user
}
