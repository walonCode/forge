package health

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	service *Service
}

func newHandler(s *Service) *Handler {
	return &Handler{
		service: s,
	}
}

// Health godoc
//
//	@Summary		Health check
//	@Description	Returns API liveness status
//	@Tags			health
//	@Produce		json
//	@Success		200	{object}	HealthResponse
//	@Router			/health [get]
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	resp := h.service.Health()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// Ready godoc
//
//	@Summary		Readiness check
//	@Description	Returns API and dependency readiness status
//	@Tags			health
//	@Produce		json
//	@Success		200		{object}	ReadinessResponse
//	@Failure		503		{object}	ReadinessResponse
//	@Router			/health/ready [get]
func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	resp, ok := h.service.Ready(r.Context())

	status := http.StatusOK
	if !ok {
		status = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}
