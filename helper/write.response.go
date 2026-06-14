package helper

import (
	"encoding/json"
	"net/http"
	"rompelago/dto"
)

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, dto.ApiError{Status: status, Error: message})
}
