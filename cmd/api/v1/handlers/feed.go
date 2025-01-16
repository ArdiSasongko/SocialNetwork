package handlers

import (
	"net/http"

	"github.com/ArdiSasongko/SocialNetwork/internal/models"
	"github.com/ArdiSasongko/SocialNetwork/internal/service"
	"github.com/ArdiSasongko/SocialNetwork/internal/storage/postgresql"
	"github.com/ArdiSasongko/SocialNetwork/utils"
)

type FeedHandler struct {
	service service.Service
	json    utils.JsonUtils
	error   utils.ErrorUtils
}

func (h *FeedHandler) GetFeeds(w http.ResponseWriter, r *http.Request) {
	user := getUserfromCtx(r)
	pf := postgresql.Pagination{
		Limit:  10,
		Offset: 0,
		Sort:   "asc",
	}

	pf, err := pf.Parse(r)
	if err != nil {
		h.error.BadRequestError(w, r, err)
		return
	}

	if err := pf.Validate(); err != nil {
		h.error.BadRequestError(w, r, err)
		return
	}

	feeds, err := h.service.Feeds.GetFeeds(r.Context(), user.ID, pf)
	if err != nil {
		h.error.InternalServerError(w, r, err)
		return
	}

	if err := h.json.JsonResponse(w, http.StatusOK, feeds); err != nil {
		h.error.InternalServerError(w, r, err)
		return
	}
}

func (h *FeedHandler) GetFeed(w http.ResponseWriter, r *http.Request) {
	post := getPostfromCtx(r)

	feed, err := h.service.Feeds.GetFeed(r.Context(), post.ID)
	if err != nil {
		h.error.InternalServerError(w, r, err)
		return
	}

	if err := h.json.JsonResponse(w, http.StatusOK, feed); err != nil {
		h.error.InternalServerError(w, r, err)
		return
	}
}

func (h *FeedHandler) LikedFeed(w http.ResponseWriter, r *http.Request) {
	post := getPostfromCtx(r)
	user := getUserfromCtx(r)

	payload := new(models.UserActivitiesPayload)
	payload.PostID = post.ID
	payload.UserID = user.ID

	if err := h.service.Feeds.LikePost(r.Context(), payload); err != nil {
		h.error.InternalServerError(w, r, err)
		return
	}

	if err := h.json.JsonResponse(w, http.StatusOK, nil); err != nil {
		h.error.InternalServerError(w, r, err)
		return
	}
}

func (h *FeedHandler) DisikedFeed(w http.ResponseWriter, r *http.Request) {
	post := getPostfromCtx(r)
	user := getUserfromCtx(r)

	payload := new(models.UserActivitiesPayload)
	payload.PostID = post.ID
	payload.UserID = user.ID

	if err := h.service.Feeds.DislikePost(r.Context(), payload); err != nil {
		h.error.InternalServerError(w, r, err)
		return
	}

	if err := h.json.JsonResponse(w, http.StatusOK, nil); err != nil {
		h.error.InternalServerError(w, r, err)
		return
	}
}

func (h *FeedHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	post := getPostfromCtx(r)
	user := getUserfromCtx(r)

	payload := new(models.CommentPayload)

	if err := h.json.ReadJSON(w, r, payload); err != nil {
		h.error.BadRequestError(w, r, err)
		return
	}

	if err := payload.Validate(); err != nil {
		h.error.BadRequestError(w, r, err)
		return
	}

	payload.UserID = user.ID
	payload.PostID = post.ID

	if err := h.service.Feeds.CreateCommentPost(r.Context(), payload); err != nil {
		h.error.InternalServerError(w, r, err)
	}

	if err := h.json.JsonResponse(w, http.StatusOK, nil); err != nil {
		h.error.InternalServerError(w, r, err)
		return
	}
}
