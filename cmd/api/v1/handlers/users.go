package handlers

import (
	"net/http"

	"github.com/ArdiSasongko/SocialNetwork/cmd/api/v1/middlewares"
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

func getUserfromCtx(r *http.Request) *postgresql.User {
	user, _ := r.Context().Value(middlewares.UserCtx).(*postgresql.User)
	return user
}
