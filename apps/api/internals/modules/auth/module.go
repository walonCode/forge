package auth

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
)

type Module struct {
	handler *Handler
}

func New(db *sql.DB) *Module {
	repo := newRepository(db)
	service := newService(repo)
	handler := newHandler(service)

	return &Module{
		handler: handler,
	}
}

func (m *Module) Register(r chi.Router) {
	r.Post("/auth/login", m.handler.LoginHandler)
	r.Post("/auth/signup", m.handler.SignupHandler)
	r.Post("/auth/logout", nil)
	r.Post("/auth/update_password", nil)
}
