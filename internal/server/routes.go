package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	webDir := "cmd/web"

	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir(webDir))))

	r.Get("/hello", s.HelloWorldHandler)
	r.Get("/health", s.healthHandler)
	r.Get("/api/orders", s.HandleGetOrderById)
	r.Post("/api/orders", s.HandlePostOrder)
	r.Delete("/api/orders", s.HandleDeleteOrderById)

	return r
}
