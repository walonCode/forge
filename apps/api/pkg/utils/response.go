package utils

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Error string `json:"error"`
	OK    bool   `json:"ok"`
}

func ErrorResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := Error{
		Error: message,
		OK:    false,
	}

	json.NewEncoder(w).Encode(resp)
}

func SuccessResponse(w http.ResponseWriter, code int, param any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(param)
}
