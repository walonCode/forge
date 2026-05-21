package tasks

import (
	"api/internals/middleware"
	"api/internals/modules/auth"
	"api/pkg/utils"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
	logger  *slog.Logger
}

func newHandler(s *Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: s,
		logger:  logger,
	}
}

// CreateTask godoc
//
//	@Summary		Create a task
//	@Description	Create a new task for the authenticated user
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		CreateTaskRequest	true	"Task details"
//	@Success		201		{object}	TaskResponse
//	@Failure		400		{object}	utils.Error
//	@Failure		500		{object}	utils.Error
//	@Router			/task [post]
func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	correlationID := middleware.GetCorrelationID(r.Context())

	var request CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("CreateTask: failed to decode request body", slog.String("correlation_id", correlationID), slog.String("error", err.Error()))
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userId, _ := r.Context().Value(auth.UserIdKey).(string)

	value, err := h.service.CreateTask(r.Context(), request, userId)
	if err != nil {
		h.logger.Error("CreateTask: service error", slog.String("correlation_id", correlationID), slog.String("user_id", userId), slog.String("error", err.Error()))
		utils.ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	h.logger.Info("CreateTask: task created", slog.String("correlation_id", correlationID), slog.String("user_id", userId))

	utils.SuccessResponse(w, http.StatusCreated, TaskResponse{
		Message: "task created successfully",
		Data:    value,
	})
}

// GetTask godoc
//
//	@Summary		Get all tasks
//	@Description	Retrieve all tasks for the authenticated user
//	@Tags			tasks
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	TaskResponse
//	@Failure		401	{object}	utils.Error
//	@Failure		500	{object}	utils.Error
//	@Router			/tasks [get]
func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	correlationID := middleware.GetCorrelationID(r.Context())
	userId, _ := r.Context().Value(auth.UserIdKey).(string)

	task, err := h.service.GetTasks(r.Context(), userId)
	if err != nil {
		h.logger.Error("GetTasks: service error", slog.String("correlation_id", correlationID), slog.String("user_id", userId), slog.String("error", err.Error()))
		utils.ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	h.logger.Info("GetTasks: fetched tasks", slog.String("correlation_id", correlationID), slog.String("user_id", userId))

	utils.SuccessResponse(w, http.StatusOK, TaskResponse{
		Message: "all user tasks",
		Data:    &task,
	})
}

// DeleteTask godoc
//
//	@Summary		Delete a task
//	@Description	Delete a task by ID for the authenticated user
//	@Tags			tasks
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Task ID"
//	@Success		204	{object}	TaskResponse
//	@Failure		401	{object}	utils.Error
//	@Failure		500	{object}	utils.Error
//	@Router			/task/{id} [delete]
func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	correlationID := middleware.GetCorrelationID(r.Context())

	taskId := chi.URLParam(r, "id")
	if strings.TrimSpace(taskId) == "" {
		h.logger.Warn("DeleteTask: missing task id", slog.String("correlation_id", correlationID))
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid url parameter")
		return
	}

	userId, _ := r.Context().Value(auth.UserIdKey).(string)

	if err := h.service.DeleteTask(r.Context(), taskId, userId); err != nil {
		h.logger.Error("DeleteTask: service error", slog.String("correlation_id", correlationID), slog.String("user_id", userId), slog.String("task_id", taskId), slog.String("error", err.Error()))
		utils.ErrorResponse(w, http.StatusInternalServerError, "failed to delete task")
		return
	}

	h.logger.Info("DeleteTask: task deleted", slog.String("correlation_id", correlationID), slog.String("user_id", userId), slog.String("task_id", taskId))

	utils.SuccessResponse(w, http.StatusNoContent, TaskResponse{
		Message: "task delete successfully",
	})
}

// GetTask godoc
//
//	@Summary		Get a single task details
//	@Description	Get a task by ID for the authenticated user
//	@Tags			tasks
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Task ID"
//	@Success		200	{object}	TaskResponse
//	@Failure		401	{object}	utils.Error
//	@Failure		500	{object}	utils.Error
//	@Router			/task/{id} [get]
func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	correlationID := middleware.GetCorrelationID(r.Context())

	taskId := chi.URLParam(r, "id")
	if strings.TrimSpace(taskId) == "" {
		h.logger.Warn("GetTask: missing task id", slog.String("correlation_id", correlationID))
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid url param for taskId")
		return
	}

	userId, _ := r.Context().Value(auth.UserIdKey).(string)

	task, err := h.service.GetTask(r.Context(), userId, taskId)
	if err != nil {
		h.logger.Error("GetTask: service error", slog.String("correlation_id", correlationID), slog.String("user_id", userId), slog.String("task_id", taskId), slog.String("error", err.Error()))
		utils.ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	h.logger.Info("GetTask: fetched task", slog.String("correlation_id", correlationID), slog.String("user_id", userId), slog.String("task_id", taskId))

	utils.SuccessResponse(w, http.StatusOK, TaskResponse{
		Message: "task details",
		Data:    task,
	})
}

// UpdateTask godoc
//
//	@Summary		Update a single task details
//	@Description	Change the isCompleted field for an authenticated user
//	@Tags			tasks
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id		path		string				true	"Task ID"
//	@Param			body	body		UpdateTaskRequest	true	"Update payload"
//	@Success		200		{object}	TaskResponse
//	@Failure		400		{object}	utils.Error
//	@Failure		401		{object}	utils.Error
//	@Failure		500		{object}	utils.Error
//	@Router			/task/{id} [patch]
func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	correlationID := middleware.GetCorrelationID(r.Context())

	taskId := chi.URLParam(r, "id")
	if strings.TrimSpace(taskId) == "" {
		h.logger.Warn("UpdateTask: missing task id", slog.String("correlation_id", correlationID))
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid url param for taskId")
		return
	}

	userId, _ := r.Context().Value(auth.UserIdKey).(string)

	var isCompleted bool
	if err := json.NewDecoder(r.Body).Decode(&isCompleted); err != nil {
		h.logger.Error("UpdateTask: failed to decode request body", slog.String("correlation_id", correlationID), slog.String("error", err.Error()))
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	task, err := h.service.UpdateTask(r.Context(), userId, taskId, isCompleted)
	if err != nil {
		h.logger.Error("UpdateTask: service error", slog.String("correlation_id", correlationID), slog.String("user_id", userId), slog.String("task_id", taskId), slog.String("error", err.Error()))
		utils.ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	h.logger.Info("UpdateTask: task updated", slog.String("correlation_id", correlationID), slog.String("user_id", userId), slog.String("task_id", taskId))

	utils.SuccessResponse(w, http.StatusOK, TaskResponse{
		Message: "task updated successfully",
		Data:    task,
	})
}
