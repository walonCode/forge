package auth

import (
	"api/pkg/utils"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	service *Service
}

func newHandler(s *Service) *Handler {
	return &Handler{
		service: s,
	}
}

// login
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var request LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.ErrorResponse(
			w,
			http.StatusBadRequest,
			"invalid body",
		)
		return
	}

	//see if the user exist
	user, err := h.service.FindUserByUsername(r.Context(), request.Username)
	if err != nil {
		utils.ErrorResponse(
			w,
			http.StatusUnauthorized,
			"invalid username or password",
		)
		return
	}

	//verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		utils.ErrorResponse(
			w,
			http.StatusUnauthorized,
			"invalid username or password",
		)
		return
	}

	accessToken, err := utils.CreateToken(user.ID)
	if err != nil {
		utils.ErrorResponse(
			w,
			http.StatusInternalServerError,
			"something went wrong",
		)

		return
	}

	refreshToken, err := utils.CreateToken(user.ID)
	if err != nil {
		utils.ErrorResponse(
			w,
			http.StatusInternalServerError,
			"something went wrong",
		)

		return
	}

	userLogin := AuthResponse{
		Message: "user login succcessfully",
		Data: AuthResponseData{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}

	utils.SuccessResponse(
		w,
		http.StatusOK,
		userLogin,
	)
}

// signup
func (h *Handler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	var body SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.ErrorResponse(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	//check if user already exist
	user, err := h.service.FindUserByUsername(r.Context(), body.Username)
	if err != nil {
		utils.ErrorResponse(
			w,
			http.StatusInternalServerError,
			"something went wrong",
		)
		return
	}

	if user != nil {
		utils.ErrorResponse(
			w,
			http.StatusConflict,
			"user with username already exist",
		)
		return
	}

	//create user, and create the access and refresh tokens
	userId, err := h.service.CreateUser(r.Context(), body)
	if err != nil {
		utils.ErrorResponse(
			w,
			http.StatusInternalServerError,
			"something went wrong",
		)
		return
	}

	accessToken, err := utils.CreateToken(userId)
	if err != nil {
		utils.ErrorResponse(
			w,
			http.StatusInternalServerError,
			"something went wrong",
		)

		return
	}

	refreshToken, err := utils.CreateToken(userId)
	if err != nil {
		utils.ErrorResponse(
			w,
			http.StatusInternalServerError,
			"something went wrong",
		)

		return
	}

	userLogin := AuthResponse{
		Message: "user login succcessfully",
		Data: AuthResponseData{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}

	utils.SuccessResponse(
		w,
		http.StatusCreated,
		userLogin,
	)
}

//forget password
//logout
