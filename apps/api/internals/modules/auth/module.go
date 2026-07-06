package auth

import (
	"database/sql"
	"log/slog"

	"github.com/go-chi/chi/v5"
)

type Module struct {
	handler *Handler
}

func New(db *sql.DB, logger *slog.Logger, jwtSecret string) *Module {
	repo := newRepository(db)
	service := newService(repo)
	handler := newHandler(service, logger, jwtSecret)

	return &Module{
		handler: handler,
	}
}

func (m *Module) Register(r chi.Router) {
	r.Post("/auth/login", m.handler.LoginHandler)
	r.Post("/auth/signup", m.handler.SignupHandler)
	r.Post("/auth/refresh", m.handler.RefreshHandler)
	r.Post("/auth/logout", m.handler.LogoutHandler)
}
