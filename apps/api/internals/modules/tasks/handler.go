package tasks

import (
	"api/internals/modules/auth"
	"api/pkg/utils"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func newHandler(s *Service)*Handler {
	return &Handler{
		service:s,
	}
}

func (h *Handler)CreateTask(w http.ResponseWriter, r *http.Request){
	var request CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.ErrorResponse(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	//getting the userId
	userId := r.Context().Value(auth.UserIdKey)

	//create Create 
	value, err := h.service.CreateTask(r.Context(), request,userId.(string))
	if err != nil {
		utils.ErrorResponse(
			w,
			http.StatusInternalServerError,
			"something went wrong",
		)
		return
	}

	data := TaskResponse{
		Message: "task created successfully",
		Data: value,
	}

	utils.SuccessResponse(
		w,
		http.StatusCreated,
		data,
	)
}


func (h *Handler)GetTask(w http.ResponseWriter, r *http.Request){
	userId := r.Context().Value(auth.UserIdKey)
	if strings.TrimSpace(userId.(string)) == "" {
		utils.ErrorResponse(
			w,
			http.StatusUnauthorized,
			"user not authenticated",
		)
		return 
	}

	task, err := h.service.GetTask(r.Context(), userId.(string))
	if err != nil {
		utils.ErrorResponse(
			w,
			http.StatusInternalServerError,
			"something went wrong",
		)
		return
	}

	data := TaskResponse{
		Message: "all user tasks",
		Data: &task,
	}

	utils.SuccessResponse(
		w,
		http.StatusOK,
		data,
	)
}


func (h *Handler)DeleteTask(w http.ResponseWriter, r *http.Request){
	taskId := chi.URLParam(r, "id")
	if strings.TrimSpace(taskId) == "" {
		utils.ErrorResponse(
			w,
			http.StatusBadGateway,
			"invalid url parameter",
		)
		return
	}

	userId := r.Context().Value(auth.UserIdKey)
	if strings.TrimSpace(userId.(string)) == "" {
		utils.ErrorResponse(
			w,
			http.StatusUnauthorized,
			"user not authenticated",
		)
		return 
	}

	if err := h.service.DeleteTask(r.Context(), taskId, userId.(string)); err != nil {
		utils.ErrorResponse(
			w,
			http.StatusInternalServerError,
			"failed to deleted tasks",
		)
		return
	}

	utils.SuccessResponse(
		w,
		http.StatusNoContent,
		TaskResponse{
			Message: "task delete successfully",
		},
	)
}
