package auth

import (
	"api/pkg/utils"
	"context"
	"fmt"
	"net/http"
	"strings"
)

type contextKey string
const UserIdKey contextKey = "userId"

func AuthMiddlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		//logger
		// fmt.Println(authHeader)
		if strings.TrimSpace(authHeader) == "" {
			utils.ErrorResponse(w, 401, "user not authenticated")
			return
		}

		authValue := strings.Split(authHeader, " ")
		if len(authValue) < 2 || authValue[1] == "" {
			utils.ErrorResponse(w, 401, "user not authenticated")
			return
		}
		// fmt.Println(authValue[1])
		tokenValue, err := utils.VerifyToken(authValue[1])
		fmt.Println(tokenValue)
		fmt.Println(err)
		if err != nil {
			utils.ErrorResponse(w, 401, "user not authenticated or invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), UserIdKey, tokenValue.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}