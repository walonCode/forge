package auth

import (
	"api/pkg/utils"
	"context"
	"net/http"
	"strings"
)

type contextKey string

const UserIdKey contextKey = "userId"

// AuthMiddleware returns middleware that authenticates requests using the
// provided JWT secret and injects the caller's user id into the context.
func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if strings.TrimSpace(authHeader) == "" {
				utils.ErrorResponse(w, http.StatusUnauthorized, "user not authenticated")
				return
			}

			authValue := strings.Split(authHeader, " ")
			if len(authValue) < 2 || authValue[1] == "" {
				utils.ErrorResponse(w, http.StatusUnauthorized, "user not authenticated")
				return
			}

			tokenValue, err := utils.VerifyToken(authValue[1], secret)
			if err != nil {
				utils.ErrorResponse(w, http.StatusUnauthorized, "user not authenticated or invalid token")
				return
			}

			// a refresh token must not be accepted as an access token
			if tokenValue.Type != string(utils.AccessToken) {
				utils.ErrorResponse(w, http.StatusUnauthorized, "user not authenticated or invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), UserIdKey, tokenValue.UserId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
