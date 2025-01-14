package utils

import (
	"encoding/json"
	"log"
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

func (j *JsonUtils) ReadFormData(w http.ResponseWriter, r *http.Request, data any) error {
	// Tentukan batas ukuran form
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		log.Println(err.Error())
		return j.WriteJSONError(w, http.StatusBadRequest, "Unable to parse form data")
	}

	// Ambil data dari r.Form dan decode ke dalam struct
	formData := r.Form
	formDataJSON, err := json.Marshal(formData)
	if err != nil {
		log.Println(err.Error())
		return j.WriteJSONError(w, http.StatusInternalServerError, "Error processing form data")
	}

	err = json.Unmarshal(formDataJSON, data)
	if err != nil {
		log.Println(err.Error())
		return j.WriteJSONError(w, http.StatusInternalServerError, "Error decoding form data")
	}

	return nil
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
