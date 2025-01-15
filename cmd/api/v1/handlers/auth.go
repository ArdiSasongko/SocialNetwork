package handlers

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/ArdiSasongko/SocialNetwork/internal/models"
	"github.com/ArdiSasongko/SocialNetwork/internal/service"
	"github.com/ArdiSasongko/SocialNetwork/internal/storage/postgresql"
	"github.com/ArdiSasongko/SocialNetwork/utils"
)

type AuthHandler struct {
	service service.Service
	json    utils.JsonUtils
	error   utils.ErrorUtils
}

func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	payload := new(models.UserPayload)

	if err := h.json.ReadJSON(w, r, payload); err != nil {
		h.error.BadRequestError(w, r, err)
		return
	}

	if err := payload.Validate(); err != nil {
		h.error.BadRequestError(w, r, err)
		return
	}

	if err := h.service.Auth.RegisterUser(r.Context(), payload); err != nil {
		switch {
		case errors.Is(err, postgresql.ErrDuplicateEmail):
			h.error.BadRequestError(w, r, err)
		case errors.Is(err, postgresql.ErrDuplicateUsername):
			h.error.BadRequestError(w, r, err)
		default:
			h.error.InternalServerError(w, r, err)
		}
		return
	}

	if err := h.json.JsonResponse(w, http.StatusCreated, nil); err != nil {
		h.error.InternalServerError(w, r, err)
		return
	}
}

func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	payload := new(models.LoginPayload)

	if err := h.json.ReadJSON(w, r, payload); err != nil {
		h.error.BadRequestError(w, r, err)
		return
	}

	if err := payload.Validate(); err != nil {
		h.error.BadRequestError(w, r, err)
		return
	}

	token, err := h.service.Auth.LoginUser(r.Context(), payload)
	if err != nil {
		switch {
		case errors.Is(err, postgresql.ErrNotFound):
			h.error.NotFoundError(w, r, err)
		default:
			h.error.InternalServerError(w, r, err)
		}
		return
	}

	if err := h.json.JsonResponse(w, http.StatusOK, token); err != nil {
		h.error.InternalServerError(w, r, err)
		return
	}
}

func extractFile(r *http.Request, fieldName string) (*multipart.FileHeader, error) {
	file, fileheader, err := r.FormFile(fieldName)
	if err != nil && err != http.ErrMissingFile {
		return nil, err
	}

	if file == nil {
		return nil, nil
	}

	defer file.Close()

	if fileheader.Size > 2*1024*1024 {
		return nil, fmt.Errorf("file to large (max 2mb)")
	}

	contentType := fileheader.Header.Get("Content-Type")
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}
	if !allowedTypes[contentType] {
		return nil, fmt.Errorf("invalid file type")
	}

	ext := strings.ToLower(filepath.Ext(fileheader.Filename))
	allowedExt := map[string]bool{
		".jpg":  true,
		".png":  true,
		".jpeg": true,
	}

	if !allowedExt[ext] {
		return nil, fmt.Errorf("invalid file type only (jpg,png,jpeg)")
	}

	return fileheader, nil
}
