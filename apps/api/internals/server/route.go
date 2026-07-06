package server

import (
	"api/internals/modules/auth"
	"api/internals/modules/health"
	"api/internals/modules/tasks"
	"api/internals/modules/users"

	"github.com/go-chi/chi/v5"
)

func (s *Server) mountModules() {
	health.New(s.db).Register(s.router)
	auth.New(s.db, s.logger).Register(s.router)

	//protect route
	s.router.Group(func(r chi.Router) {
		r.Use(auth.AuthMiddlware)
		tasks.New(s.db, s.logger).Register(r)
		users.New(s.db, s.logger).Register(r)
	})
}
