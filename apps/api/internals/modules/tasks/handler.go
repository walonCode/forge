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
func (h *Handler)GetTasks(w http.ResponseWriter, r *http.Request){
	userId := r.Context().Value(auth.UserIdKey)
	if strings.TrimSpace(userId.(string)) == "" {
		utils.ErrorResponse(
			w,
			http.StatusUnauthorized,
			"user not authenticated",
		)
		return 
	}

	task, err := h.service.GetTasks(r.Context(), userId.(string))
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
func (h *Handler)GetTask(w http.ResponseWriter, r *http.Request){
	//taskId
	taskId := chi.URLParam(r, "id")
	if strings.TrimSpace(taskId) == ""{
		utils.ErrorResponse(
			w,
			http.StatusBadRequest,
			"invalid url param for taskId",
		)
		return 
	}

	//userId
	userId := r.Context().Value(auth.UserIdKey)
	if strings.TrimSpace(userId.(string)) == "" {
		utils.ErrorResponse(
			w,
			http.StatusUnauthorized,
			"user not authenticated",
		)
		return 
	}

	//get the tasks 
	task, err := h.service.GetTask(r.Context(), userId.(string), taskId)
	if err != nil {
		utils.ErrorResponse(
			w,
			http.StatusInternalServerError,
			"something went wrong",
		)
		return
	}

	utils.SuccessResponse(
		w,
		http.StatusOK,
		TaskResponse{
			Message: "task details",
			Data: task,
		},
	)
}


// UpdateTask godoc
//
//	@Summary		Update a single task details
//	@Description	Change the isCompleted field for an authenticated user
//	@Tags			tasks
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Task ID"
//	@Success		200	{object}	TaskResponse
//	@Failure		401	{object}	utils.Error
//	@Failure		500	{object}	utils.Error
//	@Router			/task/{id} [patch]
func (h *Handler)UpdateTask(w http.ResponseWriter, r *http.Request){
	//taskId
	taskId := chi.URLParam(r, "id")
	if strings.TrimSpace(taskId) == ""{
		utils.ErrorResponse(
			w,
			http.StatusBadRequest,
			"invalid url param for taskId",
		)
		return 
	}

	//userId
	userId := r.Context().Value(auth.UserIdKey)
	if strings.TrimSpace(userId.(string)) == "" {
		utils.ErrorResponse(
			w,
			http.StatusUnauthorized,
			"user not authenticated",
		)
		return 
	}

	var isCompleted bool
	if err := json.NewDecoder(r.Body).Decode(&isCompleted); err != nil {
		utils.ErrorResponse(
			w,
			http.StatusBadGateway,
			"invalid request body",
		)
		return 
	}

	task, err := h.service.UpdateTask(r.Context(), userId.(string), taskId, isCompleted)
	if err != nil {
		utils.ErrorResponse(
			w,
			http.StatusInternalServerError,
			"something went wrong",
		)
		return 
	}

	utils.SuccessResponse(
		w,
		http.StatusOK,
		TaskResponse{
			Message: "task updated successfully",
			Data: task,
		},
	)
}