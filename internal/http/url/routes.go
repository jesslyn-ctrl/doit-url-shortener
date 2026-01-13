package url

import (
	"time"

	chiRouter "github.com/go-chi/chi/v5"
	_domainUrl "github.com/jesslyn-ctrl/doit-url-shortener/internal/domain/url"
)

// RegisterUrlRoutes registers URL-related HTTP routes
func RegisterUrlRoutes(
	r chiRouter.Router,
	svc *_domainUrl.Service,
	defaultTTL time.Duration,
) {
	r.Post("/shorten", CreateShortURLHandler(svc, defaultTTL))
	r.Get("/s/{code}", ResolveHandler(svc))
	r.Get("/stats/{code}", GetStatsHandler(svc))
}
