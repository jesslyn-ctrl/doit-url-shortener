package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_domainUrl "github.com/jesslyn-ctrl/doit-url-shortener/internal/domain/url"
	_httpHealth "github.com/jesslyn-ctrl/doit-url-shortener/internal/http/health"
	"github.com/jesslyn-ctrl/doit-url-shortener/internal/http/url"
)

func NewRouter(
	urlSvc *_domainUrl.Service,
	defaultTTL time.Duration,
) http.Handler {
	r := chi.NewRouter()

	// Use built-in middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	// Custom middleware by adding header X-Processing-Time-Micros
	r.Use(ProcessingTimeMiddleware)

	// Register health routes
	_httpHealth.RegisterHealthRoutes(r)
	// Register url routes
	url.RegisterUrlRoutes(r, urlSvc, defaultTTL)

	return r
}
