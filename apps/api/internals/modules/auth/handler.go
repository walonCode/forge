package auth

import (
	"api/internals/middleware"
	"api/pkg/utils"
	"encoding/json"
	"log/slog"
	"net/http"

	"golang.org/x/crypto/bcrypt"
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

// LoginHandler godoc
//
//	@Summary		Login
//	@Description	Authenticate a user and return access and refresh tokens
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		LoginRequest	true	"Login credentials"
//	@Success		200		{object}	AuthResponse
//	@Failure		400		{object}	utils.Error
//	@Failure		401		{object}	utils.Error
//	@Router			/auth/login [post]
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	correlationID := middleware.GetCorrelationID(r.Context())

	var request LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("LoginHandler: failed to decode request body", slog.String("correlation_id", correlationID), slog.String("error", err.Error()))
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid body")
		return
	}

	user, err := h.service.FindUserByUsername(r.Context(), request.Username)
	if err != nil {
		h.logger.Warn("LoginHandler: user not found", slog.String("correlation_id", correlationID), slog.String("username", request.Username))
		utils.ErrorResponse(w, http.StatusUnauthorized, "invalid username or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		h.logger.Warn("LoginHandler: invalid password", slog.String("correlation_id", correlationID), slog.String("username", request.Username))
		utils.ErrorResponse(w, http.StatusUnauthorized, "invalid username or password")
		return
	}

	accessToken, err := utils.CreateToken(user.ID)
	if err != nil {
		h.logger.Error("LoginHandler: failed to create access token", slog.String("correlation_id", correlationID), slog.String("error", err.Error()))
		utils.ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	refreshToken, err := utils.CreateToken(user.ID)
	if err != nil {
		h.logger.Error("LoginHandler: failed to create refresh token", slog.String("correlation_id", correlationID), slog.String("error", err.Error()))
		utils.ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	h.logger.Info("LoginHandler: user logged in", slog.String("correlation_id", correlationID), slog.String("username", request.Username))

	utils.SuccessResponse(w, http.StatusOK, AuthResponse{
		Message: "user login succcessfully",
		Data: AuthResponseData{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}

// SignupHandler godoc
//
//	@Summary		Signup
//	@Description	Register a new user and return access and refresh tokens
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		SignupRequest	true	"Signup details"
//	@Success		201		{object}	AuthResponse
//	@Failure		400		{object}	utils.Error
//	@Failure		409		{object}	utils.Error
//	@Failure		500		{object}	utils.Error
//	@Router			/auth/signup [post]
func (h *Handler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	correlationID := middleware.GetCorrelationID(r.Context())

	var body SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.logger.Error("SignupHandler: failed to decode request body", slog.String("correlation_id", correlationID), slog.String("error", err.Error()))
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.service.FindUserByUsername(r.Context(), body.Username)
	if err != nil {
		h.logger.Error("SignupHandler: error checking existing user", slog.String("correlation_id", correlationID), slog.String("username", body.Username), slog.String("error", err.Error()))
		utils.ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	if user != nil {
		h.logger.Warn("SignupHandler: username already taken", slog.String("correlation_id", correlationID), slog.String("username", body.Username))
		utils.ErrorResponse(w, http.StatusConflict, "user with username already exist")
		return
	}

	userId, err := h.service.CreateUser(r.Context(), body)
	if err != nil {
		h.logger.Error("SignupHandler: failed to create user", slog.String("correlation_id", correlationID), slog.String("username", body.Username), slog.String("error", err.Error()))
		utils.ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	accessToken, err := utils.CreateToken(userId)
	if err != nil {
		h.logger.Error("SignupHandler: failed to create access token", slog.String("correlation_id", correlationID), slog.String("error", err.Error()))
		utils.ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	refreshToken, err := utils.CreateToken(userId)
	if err != nil {
		h.logger.Error("SignupHandler: failed to create refresh token", slog.String("correlation_id", correlationID), slog.String("error", err.Error()))
		utils.ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	h.logger.Info("SignupHandler: user created", slog.String("correlation_id", correlationID), slog.String("username", body.Username))

	utils.SuccessResponse(w, http.StatusCreated, AuthResponse{
		Message: "user login succcessfully",
		Data: AuthResponseData{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}

//forget password
//logout
