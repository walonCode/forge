package auth

import (
	"api/pkg/utils"
	"context"
	"net/http"
	"strings"
)

type contextKey string
const UserIdKey contextKey = "userId"

func AuthMiddlware(next http.Handler)http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//getting the authorization header
		authHeader := r.Header.Get("Authorization")
		if strings.TrimSpace(authHeader) == ""{
			utils.ErrorResponse(
				w,
				401,
				"user not authenticated",
			)
			return
		}

		//spliting the token Bearer token
		authValue := strings.Split(authHeader, " ")
		if authValue[1] == "" || len(authValue) < 2 {
			utils.ErrorResponse(
				w,
				401,
				"user not authenticated",
			)
			return
		}

		//verify the token
		tokenValue, err := utils.VerifyToken(authValue[1])
		if err != nil {
			utils.ErrorResponse(
				w,
				401,
				"user not authenticated or invalid token",
			)
			return
		}

		//getting the userId
		userId := tokenValue.UserId

		//sending the userId in the context
		ctx := context.WithValue(r.Context(), UserIdKey, userId)
			
		next.ServeHTTP(w,r.WithContext(ctx))
	})
}