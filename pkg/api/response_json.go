package api

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func ErrorJSON(w http.ResponseWriter, status int, err error) error {

	errorResponse := map[string]string{
		"error": err.Error(),
	}
	return WriteJSON(w, status, errorResponse)
}
