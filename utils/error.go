package utils

import (
	"log"
	"net/http"
)

type ErrorUtils struct {
	json JsonUtils
}

func NewErrorUtils() ErrorUtils {
	json := NewJsonUtils()
	return ErrorUtils{
		json: json,
	}
}

func (e *ErrorUtils) InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error, method: %v, path :%v, message: %v", r.Method, r.URL.Path, err.Error())
	e.json.WriteJSONError(w, http.StatusInternalServerError, err.Error())
}

func (e *ErrorUtils) BadRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error, method: %v, path :%v, message: %v", r.Method, r.URL.Path, err.Error())
	e.json.WriteJSONError(w, http.StatusBadRequest, err.Error())
}

func (e *ErrorUtils) NotFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error, method: %v, path :%v, message: %v", r.Method, r.URL.Path, err.Error())
	e.json.WriteJSONError(w, http.StatusNotFound, err.Error())
}
