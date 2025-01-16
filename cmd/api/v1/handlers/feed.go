package handlers

import (
	"net/http"

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
