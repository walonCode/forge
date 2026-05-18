package health

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	service *Service
}

func newHandler(s *Service)*Handler{
	return &Handler{
		service: s,
	}
}

func(h *Handler)Health(w http.ResponseWriter, r *http.Request){
	resp := h.service.Health()

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler)Ready(w http.ResponseWriter, r *http.Request){
	resp, ok := h.service.Ready(r.Context())

	status := http.StatusOK
	if !ok {
		status = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}