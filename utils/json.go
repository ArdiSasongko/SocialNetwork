package utils

import (
	"encoding/json"
	"net/http"
)

type JsonUtils struct{}

func NewJsonUtils() JsonUtils {
	return JsonUtils{}
}

func (j *JsonUtils) writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func (j *JsonUtils) ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxByte := 1_048_578 //1mb
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxByte))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func (j *JsonUtils) WriteJSONError(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}

	w.Header().Set("Content-Type", "application/json")
	return j.writeJSON(w, status, &envelope{Error: message})
}

func (j *JsonUtils) JsonResponse(w http.ResponseWriter, status int, data any) error {
	type envelope struct {
		Status int `json:"status"`
		Data   any `json:"data"`
	}

	return j.writeJSON(w, status, &envelope{
		Status: status,
		Data:   data,
	})
}
