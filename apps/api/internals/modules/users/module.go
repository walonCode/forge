package users

import (
	"database/sql"
	"log/slog"

	"github.com/go-chi/chi/v5"
)

type Module struct {
	handler *Handler
}

func New(db *sql.DB, logger *slog.Logger) *Module {
	repo := newRepository(db)
	service := newService(repo)
	handler := newHandler(service, logger)

	return &Module{
		handler: handler,
	}
}

func (m *Module) Register(r chi.Router) {
	r.Get("/user/profile", m.handler.GetProfile)
	r.Patch("/user/profile", m.handler.UpdateProfile)
	r.Patch("/user/password", m.handler.UpdatePassword)
	r.Delete("/user", m.handler.DeleteAccount)
}
