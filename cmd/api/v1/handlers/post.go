package handlers

import (
	"context"
	"fmt"
	"log"
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

func (h *PostHandler) GetPostByUser(w http.ResponseWriter, r *http.Request) {
	post := getPostfromCtx(r)
	user := getUserfromCtx(r)

	if post.UserID != user.ID {
		h.error.BadRequestError(w, r, fmt.Errorf("source not found in this user"))
		return
	}

	if err := h.json.JsonResponse(w, http.StatusOK, post); err != nil {
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

func (h *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	payload := new(models.PostUpdatePayload)
	post := getPostfromCtx(r)

	if err := h.json.ReadJSON(w, r, payload); err != nil {
		h.error.BadRequestError(w, r, err)
		return
	}

	if err := payload.Validate(); err != nil {
		h.error.BadRequestError(w, r, err)
		return
	}

	if err := h.service.Post.UpdatePost(r.Context(), post, payload); err != nil {
		h.error.InternalServerError(w, r, err)
		return
	}

	if err := h.json.JsonResponse(w, http.StatusOK, nil); err != nil {
		h.error.InternalServerError(w, r, err)
		return
	}
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	post := getPostfromCtx(r)

	if err := h.service.Post.DeletePost(r.Context(), post.ID); err != nil {
		h.error.InternalServerError(w, r, err)
		return
	}

	if err := h.json.JsonResponse(w, http.StatusNoContent, nil); err != nil {
		h.error.InternalServerError(w, r, err)
		return
	}
}

func getPostfromCtx(r *http.Request) *postgresql.Post {
	user, _ := r.Context().Value(middlewares.PostCtx).(*postgresql.Post)
	return user
}

func (h *PostHandler) CheckOwnerPost(allowRole string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		post := getPostfromCtx(r)
		user := getUserfromCtx(r)

		if user.ID == post.UserID {
			next.ServeHTTP(w, r)
			return
		}

		// check roles allowed
		allow, err := h.checkRoleAllowed(r.Context(), user, allowRole)
		if err != nil {
			h.error.InternalServerError(w, r, err)
			return
		}

		if !allow {
			h.error.ForbiddenError(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *PostHandler) checkRoleAllowed(ctx context.Context, user *postgresql.User, roleRequired string) (bool, error) {
	role, err := h.service.Role.GetRole(ctx, roleRequired)
	if err != nil {
		return false, err
	}

	return user.Role.Level >= role.Level, nil
}
