package auth

import (
	"api/internals/middleware"
	"api/pkg/database"
	"api/pkg/utils"
	"encoding/json"
	"log/slog"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	service *Service
	logger  *slog.Logger
	secret  string
}

func newHandler(s *Service, logger *slog.Logger, secret string) *Handler {
	return &Handler{
		service: s,
		logger:  logger,
		secret:  secret,
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

	if err := utils.ValidateStruct(request); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid username or password")
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

	accessToken, refreshToken, err := h.issueTokens(user.ID)
	if err != nil {
		h.logger.Error("LoginHandler: failed to create tokens", slog.String("correlation_id", correlationID), slog.String("error", err.Error()))
		utils.ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	h.logger.Info("LoginHandler: user logged in", slog.String("correlation_id", correlationID), slog.String("username", request.Username))

	utils.SuccessResponse(w, http.StatusOK, AuthResponse{
		Message: "user logged in successfully",
		Data: &AuthResponseData{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}

// issueTokens creates a fresh access/refresh token pair for the user.
func (h *Handler) issueTokens(userID string) (access, refresh string, err error) {
	access, err = utils.CreateToken(userID, h.secret, utils.AccessToken, utils.AccessTokenTTL)
	if err != nil {
		return "", "", err
	}
	refresh, err = utils.CreateToken(userID, h.secret, utils.RefreshToken, utils.RefreshTokenTTL)
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
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

	if err := utils.ValidateStruct(body); err != nil {
		h.logger.Warn("SignupHandler: validation failed", slog.String("correlation_id", correlationID), slog.String("error", err.Error()))
		utils.ErrorResponse(w, http.StatusBadRequest, "name and username must be at least 2 characters and password at least 8")
		return
	}

	userId, err := h.service.CreateUser(r.Context(), body)
	if err != nil {
		// the unique index on username is the source of truth; this also closes
		// the check-then-insert race two concurrent signups would otherwise hit
		if database.IsUniqueViolation(err) {
			h.logger.Warn("SignupHandler: username already taken", slog.String("correlation_id", correlationID), slog.String("username", body.Username))
			utils.ErrorResponse(w, http.StatusConflict, "user with username already exist")
			return
		}
		h.logger.Error("SignupHandler: failed to create user", slog.String("correlation_id", correlationID), slog.String("username", body.Username), slog.String("error", err.Error()))
		utils.ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	accessToken, refreshToken, err := h.issueTokens(userId)
	if err != nil {
		h.logger.Error("SignupHandler: failed to create tokens", slog.String("correlation_id", correlationID), slog.String("error", err.Error()))
		utils.ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	h.logger.Info("SignupHandler: user created", slog.String("correlation_id", correlationID), slog.String("username", body.Username))

	utils.SuccessResponse(w, http.StatusCreated, AuthResponse{
		Message: "user created successfully",
		Data: &AuthResponseData{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}

// LogoutHandler godoc
//
//	@Summary		Logout
//	@Description	Stateless logout. Tokens are not tracked server-side, so the client should discard them.
//	@Tags			auth
//	@Produce		json
//	@Success		200	{object}	AuthResponse
//	@Router			/auth/logout [post]
func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	correlationID := middleware.GetCorrelationID(r.Context())

	h.logger.Info("LogoutHandler: user logged out", slog.String("correlation_id", correlationID))

	utils.SuccessResponse(w, http.StatusOK, AuthResponse{
		Message: "user logged out successfully",
	})
}

// RefreshHandler godoc
//
//	@Summary		Refresh tokens
//	@Description	Exchange a valid refresh token for a new access/refresh token pair
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		RefreshRequest	true	"Refresh token"
//	@Success		200		{object}	AuthResponse
//	@Failure		400		{object}	utils.Error
//	@Failure		401		{object}	utils.Error
//	@Router			/auth/refresh [post]
func (h *Handler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	correlationID := middleware.GetCorrelationID(r.Context())

	var request RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("RefreshHandler: failed to decode request body", slog.String("correlation_id", correlationID), slog.String("error", err.Error()))
		utils.ErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := utils.ValidateStruct(request); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "refresh_token is required")
		return
	}

	claims, err := utils.VerifyToken(request.RefreshToken, h.secret)
	if err != nil || claims.Type != string(utils.RefreshToken) {
		h.logger.Warn("RefreshHandler: invalid refresh token", slog.String("correlation_id", correlationID))
		utils.ErrorResponse(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	accessToken, refreshToken, err := h.issueTokens(claims.UserId)
	if err != nil {
		h.logger.Error("RefreshHandler: failed to create tokens", slog.String("correlation_id", correlationID), slog.String("error", err.Error()))
		utils.ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	h.logger.Info("RefreshHandler: tokens refreshed", slog.String("correlation_id", correlationID), slog.String("user_id", claims.UserId))

	utils.SuccessResponse(w, http.StatusOK, AuthResponse{
		Message: "tokens refreshed successfully",
		Data: &AuthResponseData{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}

//forget password
