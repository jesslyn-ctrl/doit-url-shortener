package health

import chiRouter "github.com/go-chi/chi/v5"

// RegisterHealthRoutes registers Health-related HTTP routes
func RegisterHealthRoutes(r chiRouter.Router) {
	r.Get("/health", HealthHandler())
}
