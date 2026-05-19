package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type responseWriterInterceptor struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriterInterceptor) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriterInterceptor) Write(b []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	return w.ResponseWriter.Write(b)
}

func LoggerMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			// Wrap the response writer to spy on the status code
			wrappedWriter := &responseWriterInterceptor{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(wrappedWriter, r)

			latency := time.Since(startTime)

			correlationID := GetCorrelationID(r.Context())
			if correlationID == "" {
				correlationID = wrappedWriter.Header().Get("X-Correlation-ID")
			}

			logger.Info("HTTP Request",
				slog.Int("status", wrappedWriter.statusCode),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Duration("latency", latency),
				slog.String("correlation_id", correlationID),
				slog.String("user_agent", r.UserAgent()),
			)
		})
	}
}
