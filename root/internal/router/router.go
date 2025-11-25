package router

import (
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/cors"
    "github.com/yourname/taskify/internal/handler"
    "github.com/yourname/taskify/pkg/middleware"
)

func NewRouter(authHandler *handler.AuthHandler, taskHandler *handler.TaskHandler) http.Handler {
    r := chi.NewRouter()
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins: []string{"*"},
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
        AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
    }))

    r.Route("/api", func(r chi.Router) {
        r.Post("/auth/register", authHandler.Register)
        r.Post("/auth/login", authHandler.Login)

        r.Group(func(r chi.Router) {
            r.Use(middleware.AuthMiddleware)
            r.Get("/tasks", taskHandler.GetAll)
            r.Post("/tasks", taskHandler.Create)
            r.Get("/tasks/{id}", taskHandler.GetOne)
            r.Put("/tasks/{id}", taskHandler.Update)
            r.Delete("/tasks/{id}", taskHandler.Delete)
            r.Patch("/tasks/{id}/toggle", taskHandler.Toggle)
        })
    })

    // serve frontend static files (if ./frontend exists)
    r.Handle("/", http.FileServer(http.Dir("./frontend")))
    return r
}
