package handlers

import (
	"net/http"

	"github.com/ArdiSasongko/SocialNetwork/utils"
)

type healthHandler struct {
	json  utils.JsonUtils
	error utils.ErrorUtils
}

func (h *healthHandler) Get(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "ok",
	}

	if err := h.json.JsonResponse(w, http.StatusOK, data); err != nil {
		h.error.InternalServerError(w, r, err)
		return
	}
}
