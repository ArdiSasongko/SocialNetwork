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

	userResp, err := h.service.Users.GetProfileByID(r.Context(), user.ID)
	if err != nil {
		h.error.BadRequestError(w, r, err)
		return
	}
	if err := h.json.JsonResponse(w, http.StatusOK, userResp); err != nil {
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

func (h *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	user := getUserProfileCtx(r)

	userResp, err := h.service.Users.GetProfileByID(r.Context(), user.ID)
	if err != nil {
		h.error.BadRequestError(w, r, err)
		return
	}
	if err := h.json.JsonResponse(w, http.StatusOK, userResp); err != nil {
		h.error.InternalServerError(w, r, err)
		return
	}
}

func (h *UserHandler) FollowUser(w http.ResponseWriter, r *http.Request) {
	user := getUserfromCtx(r)
	toFollow := getUserProfileCtx(r)

	if err := h.service.Users.FollowUser(r.Context(), toFollow.ID, user.ID); err != nil {
		h.error.BadRequestError(w, r, err)
		return
	}

	if err := h.json.JsonResponse(w, http.StatusCreated, nil); err != nil {
		h.error.InternalServerError(w, r, err)
		return
	}
}

func (h *UserHandler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	user := getUserfromCtx(r)
	toUnfollow := getUserProfileCtx(r)

	if err := h.service.Users.UnfollowUser(r.Context(), toUnfollow.ID, user.ID); err != nil {
		h.error.BadRequestError(w, r, err)
		return
	}

	if err := h.json.JsonResponse(w, http.StatusCreated, nil); err != nil {
		h.error.InternalServerError(w, r, err)
		return
	}
}

func getUserfromCtx(r *http.Request) *postgresql.User {
	user, _ := r.Context().Value(middlewares.UserCtx).(*postgresql.User)
	return user
}

func getUserProfileCtx(r *http.Request) *postgresql.User {
	user, _ := r.Context().Value(middlewares.UserProfileCtx).(*postgresql.User)
	return user
}
