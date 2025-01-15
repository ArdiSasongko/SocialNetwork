package handlers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

const maxSize = 2 * 1024 * 1024 // 2mb

func extractFiles(r *http.Request, fieldName string) ([]*multipart.FileHeader, error) {
	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		return nil, nil
	}

	files, ok := r.MultipartForm.File[fieldName]
	if !ok || len(files) == 0 {
		return nil, nil
	}

	validFiles := []*multipart.FileHeader{}
	for _, file := range files {
		validFile, err := validateFile(file)
		if err != nil {
			return nil, err
		}

		validFiles = append(validFiles, validFile)
	}

	return validFiles, nil
}

func extractFile(r *http.Request, fieldName string) (*multipart.FileHeader, error) {
	if r.MultipartForm == nil || r.MultipartForm.File == nil {
		return nil, nil
	}

	files, ok := r.MultipartForm.File[fieldName]
	if !ok || len(files) == 0 {
		return nil, nil
	}

	if len(files) > 1 {
		return nil, fmt.Errorf("only one image is allowed")
	}

	validFile, err := validateFile(files[0])
	if err != nil {
		return nil, err
	}

	return validFile, nil
}

func validateFile(file *multipart.FileHeader) (*multipart.FileHeader, error) {
	if file.Size > maxSize {
		return nil, fmt.Errorf("file to large (max 2mb)")
	}

	contentType := file.Header.Get("Content-Type")
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}
	if !allowedTypes[contentType] {
		return nil, fmt.Errorf("invalid file type")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExt := map[string]bool{
		".jpg":  true,
		".png":  true,
		".jpeg": true,
	}

	if !allowedExt[ext] {
		return nil, fmt.Errorf("invalid file type only (jpg,png,jpeg)")
	}

	return file, nil
}
