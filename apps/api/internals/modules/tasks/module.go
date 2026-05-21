package tasks

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
)

type Module struct {
	handler *Handler
}

func New(db *sql.DB)*Module{
	repo := newRepository(db)
	service := newService(repo)
	handler := newHandler(service)

	return &Module{
		handler: handler,
	}
}

func (m *Module)Register(r chi.Router){
	r.Post("/task", m.handler.CreateTask)
	r.Get("/tasks", m.handler.GetTasks)
	r.Delete("/task/{id}", m.handler.DeleteTask)
	r.Get("/task/{id}", m.handler.GetTask)
	r.Patch("/task/{id}", m.handler.UpdateTask)
}