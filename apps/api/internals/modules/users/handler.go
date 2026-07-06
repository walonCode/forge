package users

import (
	"api/internals/middleware"
	"api/internals/modules/auth"
	"api/pkg/utils"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
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

// GetProfile godoc
//
//	@Summary		Get profile
//	@Description	Get the authenticated user's profile
//	@Tags			users
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	UserResponse
//	@Failure		401	{object}	utils.Error
//	@Failure		404	{object}	utils.Error
//	@Failure		500	{object}	utils.Error
//	@Router			/user/profile [get]
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	correlationID := middleware.GetCorrelationID(r.Context())
	userId, _ := r.Context().Value(auth.UserIdKey).(string)

	profile, err := h.service.GetProfile(r.Context(), userId)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			h.logger.Warn("GetProfile: user not found", slog.String("correlation_id", correlationID), slog.String("user_id", userId))
			utils.ErrorResponse(w, http.StatusNotFound, "user not found")
			return
		}
		h.logger.Error("GetProfile: service error", slog.String("correlation_id", correlationID), slog.String("user_id", userId), slog.String("error", err.Error()))
		utils.ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, UserResponse{
		Message: "user profile retrieved successfully",
		Data:    profile,
	})
}

// UpdateProfile godoc
//
//	@Summary		Update profile
//	@Description	Update the authenticated user's name and/or username
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		UpdateProfileRequest	true	"Fields to update"
//	@Success		200		{object}	UserResponse
//	@Failure		400		{object}	utils.Error
//	@Failure		401		{object}	utils.Error
//	@Failure		404		{object}	utils.Error
//	@Failure		409		{object}	utils.Error
//	@Failure		500		{object}	utils.Error
//	@Router			/user/profile [patch]
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	correlationID := middleware.GetCorrelationID(r.Context())
	userId, _ := r.Context().Value(auth.UserIdKey).(string)

	var request UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("UpdateProfile: failed to decode request body", slog.String("correlation_id", correlationID), slog.String("error", err.Error()))
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := utils.ValidateStruct(request); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "name and username must be at least 2 characters when provided")
		return
	}

	profile, err := h.service.UpdateProfile(r.Context(), userId, request)
	if err != nil {
		switch {
		case errors.Is(err, ErrNoFieldsToUpdate):
			utils.ErrorResponse(w, http.StatusBadRequest, "no fields to update")
		case errors.Is(err, ErrUsernameTaken):
			h.logger.Warn("UpdateProfile: username taken", slog.String("correlation_id", correlationID), slog.String("user_id", userId), slog.String("username", request.Username))
			utils.ErrorResponse(w, http.StatusConflict, "username already taken")
		case errors.Is(err, ErrUserNotFound):
			utils.ErrorResponse(w, http.StatusNotFound, "user not found")
		default:
			h.logger.Error("UpdateProfile: service error", slog.String("correlation_id", correlationID), slog.String("user_id", userId), slog.String("error", err.Error()))
			utils.ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		}
		return
	}

	h.logger.Info("UpdateProfile: profile updated", slog.String("correlation_id", correlationID), slog.String("user_id", userId))

	utils.SuccessResponse(w, http.StatusOK, UserResponse{
		Message: "profile updated successfully",
		Data:    profile,
	})
}

// UpdatePassword godoc
//
//	@Summary		Update password
//	@Description	Change the authenticated user's password (requires the current password)
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			body	body		UpdatePasswordRequest	true	"Current and new password"
//	@Success		200		{object}	UserResponse
//	@Failure		400		{object}	utils.Error
//	@Failure		401		{object}	utils.Error
//	@Failure		404		{object}	utils.Error
//	@Failure		500		{object}	utils.Error
//	@Router			/user/password [patch]
func (h *Handler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	correlationID := middleware.GetCorrelationID(r.Context())
	userId, _ := r.Context().Value(auth.UserIdKey).(string)

	var request UpdatePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("UpdatePassword: failed to decode request body", slog.String("correlation_id", correlationID), slog.String("error", err.Error()))
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := utils.ValidateStruct(request); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "current and new password must be at least 8 characters")
		return
	}
	if request.CurrentPassword == request.NewPassword {
		utils.ErrorResponse(w, http.StatusBadRequest, "new password must be different from the current password")
		return
	}

	if err := h.service.UpdatePassword(r.Context(), userId, request); err != nil {
		switch {
		case errors.Is(err, ErrInvalidPassword):
			h.logger.Warn("UpdatePassword: invalid current password", slog.String("correlation_id", correlationID), slog.String("user_id", userId))
			utils.ErrorResponse(w, http.StatusBadRequest, "current password is incorrect")
		case errors.Is(err, ErrUserNotFound):
			utils.ErrorResponse(w, http.StatusNotFound, "user not found")
		default:
			h.logger.Error("UpdatePassword: service error", slog.String("correlation_id", correlationID), slog.String("user_id", userId), slog.String("error", err.Error()))
			utils.ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		}
		return
	}

	h.logger.Info("UpdatePassword: password updated", slog.String("correlation_id", correlationID), slog.String("user_id", userId))

	utils.SuccessResponse(w, http.StatusOK, UserResponse{
		Message: "password updated successfully",
	})
}

// DeleteAccount godoc
//
//	@Summary		Delete account
//	@Description	Permanently delete the authenticated user's account and all associated data
//	@Tags			users
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	UserResponse
//	@Failure		401	{object}	utils.Error
//	@Failure		500	{object}	utils.Error
//	@Router			/user [delete]
func (h *Handler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	correlationID := middleware.GetCorrelationID(r.Context())
	userId, _ := r.Context().Value(auth.UserIdKey).(string)

	if err := h.service.DeleteAccount(r.Context(), userId); err != nil {
		h.logger.Error("DeleteAccount: service error", slog.String("correlation_id", correlationID), slog.String("user_id", userId), slog.String("error", err.Error()))
		utils.ErrorResponse(w, http.StatusInternalServerError, "failed to delete account")
		return
	}

	h.logger.Info("DeleteAccount: account deleted", slog.String("correlation_id", correlationID), slog.String("user_id", userId))

	utils.SuccessResponse(w, http.StatusOK, UserResponse{
		Message: "account deleted successfully",
	})
}
