package health

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
)

type Module struct {
	handler *Handler
}

func New(db *sql.DB)*Module{
	service := newService(db)
	handler := newHandler(service)

	return &Module{
		handler: handler,
	}
}

func(m *Module)Register(r chi.Router){
	r.Get("/health", m.handler.Health)
	r.Get("/health/ready",m.handler.Ready)
}