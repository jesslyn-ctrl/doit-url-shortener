package health

import "github.com/go-chi/chi/v5"

// RegisterHealthRoutes registers Health-related HTTP routes
func RegisterHealthRoutes(r chi.Router) {
	r.Get("/health", HealthHandler())
}
