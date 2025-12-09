package router

import (
	"net/http"

	"github.com/FUADIKAMIL/taskify/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type AuthHandlers struct {
	Register http.HandlerFunc
	Login    http.HandlerFunc
}

type TaskHandlers struct {
	GetAll http.HandlerFunc
	Create http.HandlerFunc
	GetOne http.HandlerFunc
	Update http.HandlerFunc
	Delete http.HandlerFunc
	Toggle http.HandlerFunc
}

func NewRouter(auth AuthHandlers, task TaskHandlers) http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}))

	r.Route("/api", func(r chi.Router) {

		r.Post("/auth/register", auth.Register)
		r.Post("/auth/login", auth.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware)

			r.Get("/tasks", task.GetAll)
			r.Post("/tasks", task.Create)
			r.Get("/tasks/{id}", task.GetOne)
			r.Put("/tasks/{id}", task.Update)
			r.Delete("/tasks/{id}", task.Delete)
			r.Patch("/tasks/{id}/toggle", task.Toggle)
		})
	})

	r.Handle("/", http.FileServer(http.Dir("./frontend")))
	return r
}
