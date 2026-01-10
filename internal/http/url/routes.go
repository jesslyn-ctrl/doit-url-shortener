package url

import (
	"time"

	"github.com/go-chi/chi/v5"
	_domainUrl "github.com/jesslyn-ctrl/doit-url-shortener/internal/domain/url"
)

// RegisterUrlRoutes registers URL-related HTTP routes
func RegisterUrlRoutes(
	r chi.Router,
	svc *_domainUrl.Service,
	defaultTTL time.Duration,
) {
	r.Post("/shorten", CreateShortURLHandler(svc, defaultTTL))
	r.Get("/s/{code}", ResolveHandler(svc))
}
