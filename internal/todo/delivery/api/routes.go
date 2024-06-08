package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterRoutes(r *chi.Mux, handler *Handler) {
	r.Use(middleware.Logger)

	r.Post("/tasks", handler.CreateTask)

	r.Get("/tasks", handler.ListTasks)

	r.Get("/tasks/{id}", handler.GetTask)

	r.Put("/tasks/{id}", handler.UpdateTask)

	r.Delete("/tasks/{id}", handler.DeleteTask)

	r.Get("/swagger/*", httpSwagger.WrapHandler)
}
