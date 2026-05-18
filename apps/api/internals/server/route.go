package server

import (
	"api/internals/modules/auth"
	"api/internals/modules/health"

	"github.com/go-chi/chi/v5"
)

func (s *Server)mountModules(){
	health.New(s.db).Register(s.router)
	auth.New(s.db).Register(s.router)

	//protect route
	s.router.Group(func(r chi.Router) {
		r.Use()
		
	})
}