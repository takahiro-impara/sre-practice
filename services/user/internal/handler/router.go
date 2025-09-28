package handler

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(h *UserHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/healthz", h.HealthCheck)
	r.Get("/readyz", h.ReadinessCheck)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Get("/", h.ListUsers)
			r.Post("/", h.CreateUser)
			r.Get("/{userID}", h.GetUserByID)
			// r.Get("/{email}", h.GetUserByEmail)
			r.Put("/{userID}", h.UpdateUser)
			r.Delete("/{userID}", h.DeleteUser)
			r.Post("/authenticate", h.AuthenticateUser)
		})
	})

	return r
}
