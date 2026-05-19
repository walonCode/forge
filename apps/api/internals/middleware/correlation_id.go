package middleware

import (
	"api/pkg/utils"
	"context"
	"net/http"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type contextkey string

const CorrelationIDKey contextkey = "correlation_id"

func CreateCorrelationIdMiddleware(next http.Handler)http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := gonanoid.New(30)
		if err != nil {
			utils.ErrorResponse(
				w,
				500,
				"something went wrong",
			)
			return
		}

		w.Header().Set("X-Correlation-ID", id)

		//setting the context and passing it down
		ctx := context.WithValue(r.Context(), CorrelationIDKey, id)
		next.ServeHTTP(w,r.WithContext(ctx))
	})
}

func GetCorrelationID(ctx context.Context)string{
	if id, ok := ctx.Value(CorrelationIDKey).(string); ok {
		return id
	}

	return ""
}