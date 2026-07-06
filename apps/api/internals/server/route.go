package server

import (
	"net"
	"net/http"
	"time"

	"api/internals/modules/auth"
	"api/internals/modules/health"
	"api/internals/modules/tasks"
	"api/internals/modules/users"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
)

// keyByClientIP buckets rate limits by client IP. The RealIP middleware has
// already resolved r.RemoteAddr to the caller's address by the time this runs.
func keyByClientIP(r *http.Request) (string, error) {
	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return ip, nil
	}
	return r.RemoteAddr, nil
}

func (s *Server) mountModules() {
	health.New(s.db, s.config.AppVersion).Register(s.router)

	// auth endpoints are rate limited per IP to blunt credential brute-forcing
	s.router.Group(func(r chi.Router) {
		r.Use(httprate.LimitBy(20, time.Minute, keyByClientIP))
		auth.New(s.db, s.logger, s.config.JwtSecret).Register(r)
	})

	//protect route
	s.router.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddleware(s.config.JwtSecret))
		tasks.New(s.db, s.logger).Register(r)
		users.New(s.db, s.logger).Register(r)
	})
}
