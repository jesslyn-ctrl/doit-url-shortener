package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_httpHealth "github.com/jesslyn-ctrl/doit-url-shortener/internal/http/health"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	// Use built-in middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	// Register health routes
	_httpHealth.RegisterHealthRoutes(r)

	return r
}
